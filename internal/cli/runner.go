package cli

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

func buildJobs(sourceDir, targetDir, languages string) ([]translateJob, error) {
	jobs := []translateJob{}

	info, err := os.Stat(sourceDir)
	if err != nil {
		return nil, fmt.Errorf("source path does not exist: %s", sourceDir)
	}

	if info.IsDir() {
		err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() || !strings.HasSuffix(path, ".json") {
				return nil
			}

			relPath, _ := filepath.Rel(sourceDir, path)
			for lang := range strings.SplitSeq(languages, ",") {
				lang = strings.TrimSpace(lang)
				if lang == "" {
					continue
				}
				targetPath := filepath.Join(targetDir, lang, relPath)
				jobs = append(jobs, translateJob{
					sourcePath: path,
					targetPath: targetPath,
					relPath:    relPath,
					lang:       lang,
				})
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		if !strings.HasSuffix(sourceDir, ".json") {
			return nil, fmt.Errorf("source file must be a .json file: %s", sourceDir)
		}

		fileName := filepath.Base(sourceDir)
		for lang := range strings.SplitSeq(languages, ",") {
			lang = strings.TrimSpace(lang)
			if lang == "" {
				continue
			}
			targetPath := filepath.Join(targetDir, lang, fileName)
			jobs = append(jobs, translateJob{
				sourcePath: sourceDir,
				targetPath: targetPath,
				relPath:    fileName,
				lang:       lang,
			})
		}
	}

	return jobs, nil
}

func filterJobsByExisting(jobs []translateJob, skipExisting bool) ([]translateJob, int) {
	if !skipExisting {
		return jobs, 0
	}

	filtered := make([]translateJob, 0, len(jobs))
	skipped := 0
	for _, job := range jobs {
		_, err := os.Stat(job.targetPath)
		if err == nil {
			skipped++
			continue
		}
		if !os.IsNotExist(err) {
			log.Printf("Warning: failed to check target file %s: %v", job.targetPath, err)
		}
		filtered = append(filtered, job)
	}
	return filtered, skipped
}

func printPlannedJobs(jobs []translateJob, sourceLang string) {
	for _, job := range jobs {
		fmt.Printf("%s -> %s (%s -> %s)\n", job.sourcePath, job.targetPath, sourceLang, job.lang)
	}
}

func printRunSummary(stats runStats) {
	fmt.Println("Run summary:")
	fmt.Printf("  Total: %d\n", stats.Total)
	fmt.Printf("  Succeeded: %d\n", stats.Succeeded)
	fmt.Printf("  Failed: %d\n", stats.Failed)
	fmt.Printf("  Skipped: %d\n", stats.Skipped)
	fmt.Printf("  Retried: %d\n", stats.Retried)
	fmt.Printf("  Fail-fast triggered: %t\n", stats.FailFastTriggered)
	if stats.FailFastFirstError != "" {
		fmt.Printf("  First fail-fast error: %s\n", stats.FailFastFirstError)
	}
}

func shouldRetry(err error) bool {
	if err == nil {
		return false
	}

	var apiErr *openai.APIError
	if errors.As(err, &apiErr) {
		if apiErr.HTTPStatusCode == 429 || apiErr.HTTPStatusCode >= 500 {
			return true
		}
	}

	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}

	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "timeout") || strings.Contains(msg, "tempor") || strings.Contains(msg, "connection reset")
}

func retryDelayForAttempt(baseDelay time.Duration, attempt int) time.Duration {
	if baseDelay <= 0 {
		return 0
	}
	if attempt <= 1 {
		return baseDelay
	}

	delay := baseDelay
	for i := 1; i < attempt; i++ {
		if delay > 30*time.Second/2 {
			return 30 * time.Second
		}
		delay *= 2
	}
	return delay
}

func runJobs(translator *translator, jobs []translateJob, sourceLang string, concurrency int, opts runOptions) runStats {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobCh := make(chan translateJob)
	var wg sync.WaitGroup
	var succeededCount atomic.Int64
	var failedCount atomic.Int64
	var skippedCount atomic.Int64
	var retriedCount atomic.Int64
	var failFastTriggered atomic.Bool
	var failFastErr atomic.Value
	var failFastOnce sync.Once

	triggerFailFast := func(err error) {
		if !opts.failFast || err == nil {
			return
		}
		failFastOnce.Do(func() {
			failFastTriggered.Store(true)
			failFastErr.Store(err.Error())
			cancel()
		})
	}

	for range concurrency {
		wg.Go(func() {
			for {
				var job translateJob
				var ok bool
				select {
				case <-ctx.Done():
					return
				case job, ok = <-jobCh:
					if !ok {
						return
					}
				}

				if err := os.MkdirAll(filepath.Dir(job.targetPath), 0755); err != nil {
					log.Printf("Failed to create directory %s: %v", filepath.Dir(job.targetPath), err)
					failedCount.Add(1)
					triggerFailFast(err)
					continue
				}
				fmt.Printf("Translating %s to %s...\n", job.relPath, job.lang)

				var err error
				attempts := opts.retries + 1
				retried := false
				for attempt := 1; attempt <= attempts; attempt++ {
					err = translator.TranslateFile(job.sourcePath, job.targetPath, sourceLang, job.lang)
					if err == nil {
						break
					}
					if attempt == attempts || !shouldRetry(err) {
						break
					}

					delay := retryDelayForAttempt(opts.retryDelay, attempt)
					log.Printf("Retrying %s to %s (%d/%d) after %s: %v", job.relPath, job.lang, attempt+1, attempts, delay, err)
					retried = true
					if delay > 0 {
						time.Sleep(delay)
					}
				}

				if retried {
					retriedCount.Add(1)
				}

				if err != nil {
					log.Printf("Error translating %s to %s: %v", job.sourcePath, job.lang, err)
					failedCount.Add(1)
					triggerFailFast(err)
					continue
				}
				fmt.Printf("✓ Successfully translated %s to %s\n", job.relPath, job.lang)
				succeededCount.Add(1)
			}
		})
	}

	go func() {
		for idx, j := range jobs {
			select {
			case <-ctx.Done():
				skippedCount.Add(int64(len(jobs) - idx))
				close(jobCh)
				return
			case jobCh <- j:
			}
		}
		close(jobCh)
	}()

	wg.Wait()

	firstError := ""
	if raw := failFastErr.Load(); raw != nil {
		if msg, ok := raw.(string); ok {
			firstError = msg
		}
	}

	return runStats{
		Succeeded:          int(succeededCount.Load()),
		Failed:             int(failedCount.Load()),
		Skipped:            int(skippedCount.Load()),
		Retried:            int(retriedCount.Load()),
		FailFastTriggered:  failFastTriggered.Load(),
		FailFastFirstError: firstError,
	}
}
