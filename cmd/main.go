package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/wowblvck/i18n-translator/internal/providers"
)

type TranslationService interface {
	Translate(text, sourceLang, targetLang string) (string, error)
}

type I18nTranslator struct {
	service   TranslationService
	sourceDir string
	targetDir string
}

type translateJob struct {
	sourcePath string
	targetPath string
	relPath    string
	lang       string
}

func main() {
	var (
		concurrency = flag.Int("concurrency", 4, "Number of parallel workers")
		model       = flag.String("model", "gpt-4o-mini", "Model for chatgpt service (e.g., gpt-4o-mini)")
		sourceDir   = flag.String("source", "./locales/en", "Source language directory")
		targetDir   = flag.String("target", "./locales", "Target directory for translations")
		sourceLang  = flag.String("from", "en", "Source language code")
		targetLang  = flag.String("to", "es,fr,de", "Target language codes (comma-separated)")
		apiKey      = flag.String("api-key", "", "Translation API key (required for Google Translate)")
		service     = flag.String("service", "chatgpt", "Translation service (currently only 'chatgpt' is supported)")
		help        = flag.Bool("help", false, "Show help message")
		version     = flag.Bool("version", false, "Show version")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Automatic translation tool for i18n JSON files\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s --api-key=YOUR_API_KEY --from=en --to=ru,es,fr\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --source=./src/locales/en --target=./src/locales --api-key=YOUR_API_KEY\n", os.Args[0])
	}

	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	if *version {
		fmt.Println("i18n-translator v1.0.0")
		return
	}

	if *apiKey == "" {
		fmt.Fprintf(os.Stderr, "Error: API key is required for translation service\n")
		fmt.Fprintf(os.Stderr, "Use --api-key flag or set OPENAI_API_KEY environment variable\n\n")
		flag.Usage()
		os.Exit(1)
	}

	if _, err := os.Stat(*sourceDir); os.IsNotExist(err) {
		log.Fatalf("Source directory does not exist: %s", *sourceDir)
	}

	var translationService TranslationService
	switch *service {
	case "chatgpt":
		var err error
		translationService, err = providers.ChatGPTProvider(*apiKey, *model)
		if err != nil {
			log.Fatalf("Failed to initialize ChatGPT service: %v", err)
		}
	default:
		log.Fatalf("Unsupported translation service: %s", *service)
	}

	translator := &I18nTranslator{
		service:   translationService,
		sourceDir: *sourceDir,
		targetDir: *targetDir,
	}

	fmt.Printf("Starting translation from %s to languages: %s\n", *sourceLang, *targetLang)
	fmt.Printf("Source directory: %s\n", *sourceDir)
	fmt.Printf("Target directory: %s\n", *targetDir)

	jobs, err := buildJobs(*sourceDir, *targetDir, *targetLang)
	if err != nil {
		log.Fatalf("Failed to build jobs: %v", err)
	}
	if err := runJobs(translator, jobs, *sourceLang, *concurrency); err != nil {
		log.Fatalf("Error during translation: %v", err)
	}

	fmt.Println("Translation completed successfully!")
}

func (t *I18nTranslator) TranslateFile(sourceFile, targetFile, sourceLang, targetLang string) error {
	data, err := os.ReadFile(sourceFile)
	if err != nil {
		return fmt.Errorf("failed to read source file: %v", err)
	}

	translatedJSON, err := t.service.Translate(string(data), sourceLang, targetLang)
	if err != nil {
		return fmt.Errorf("failed to translate content: %v", err)
	}

	return os.WriteFile(targetFile, []byte(strings.TrimSpace(translatedJSON)), 0644)
}

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
			for _, lang := range strings.Split(languages, ",") {
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
		for _, lang := range strings.Split(languages, ",") {
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

func runJobs(translator *I18nTranslator, jobs []translateJob, sourceLang string, concurrency int) error {
	jobCh := make(chan translateJob)
	var wg sync.WaitGroup

	for range concurrency {
		wg.Go(func() {
			for job := range jobCh {
				if err := os.MkdirAll(filepath.Dir(job.targetPath), 0755); err != nil {
					log.Printf("Failed to create directory %s: %v", filepath.Dir(job.targetPath), err)
					continue
				}
				fmt.Printf("Translating %s to %s...\n", job.relPath, job.lang)
				if err := translator.TranslateFile(job.sourcePath, job.targetPath, sourceLang, job.lang); err != nil {
					log.Printf("Error translating %s to %s: %v", job.sourcePath, job.lang, err)
					continue
				}
				fmt.Printf("✓ Successfully translated %s to %s\n", job.relPath, job.lang)
			}
		})
	}

	go func() {
		for _, j := range jobs {
			jobCh <- j
		}
		close(jobCh)
	}()

	wg.Wait()
	return nil
}
