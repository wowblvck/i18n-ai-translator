package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

func runTranslateCommand(binaryName string, args []string) int {
	fs := flag.NewFlagSet(binaryName, flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var (
		concurrency  = fs.Int("concurrency", defaultConcurrency, "Number of parallel workers")
		model        = fs.String("model", "", "Model name (defaults: chatgpt=gpt-4o-mini, groq=llama-3.3-70b-versatile, gemini=gemini-2.0-flash, ollama=llama3.2)")
		sourceDir    = fs.String("source", defaultSourceDir, "Source language directory")
		targetDir    = fs.String("target", defaultTargetDir, "Target directory for translations")
		sourceLang   = fs.String("from", defaultSourceLang, "Source language code")
		targetLang   = fs.String("to", defaultTargetLang, "Target language codes (comma-separated)")
		apiKey       = fs.String("api-key", "", "API key (required for chatgpt, groq, and gemini)")
		url          = fs.String("url", "", "Base URL for ollama or lmstudio (e.g., http://localhost:11434/v1)")
		service      = fs.String("service", defaultService, "Translation service: chatgpt, groq, gemini, ollama, lmstudio")
		configFile   = fs.String("config", "", "Path to YAML config file (auto-loads .i18n-translator.yaml/.yml when omitted)")
		dryRun       = fs.Bool("dry-run", false, "Preview translation jobs without writing files")
		listFiles    = fs.Bool("list-files", false, "Print planned source/target files and exit")
		skipExisting = fs.Bool("skip-existing", false, "Skip jobs whose target file already exists")
		failFast     = fs.Bool("fail-fast", false, "Stop processing new jobs after first translation error")
		retries      = fs.Int("retries", defaultRetries, "Retry count for temporary provider errors (429/5xx/timeout)")
		retryDelay   = fs.Duration("retry-delay", defaultRetryDelay, "Base delay between retries (e.g., 500ms, 1s, 2s)")
		timeout      = fs.Duration("timeout", defaultTimeout, "Per-request timeout for provider API calls (e.g., 30s, 2m). 0 disables timeout")
		help         = fs.Bool("help", false, "Show help message")
		version      = fs.Bool("version", false, "Show version")
	)

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", binaryName)
		fmt.Fprintln(os.Stderr, "Automatic translation tool for i18n JSON files")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Commands:")
		fmt.Fprintln(os.Stderr, "  init    Generate a starter config file")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Options:")
		fs.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nExamples:")
		fmt.Fprintf(os.Stderr, "  %s --service=chatgpt --api-key=YOUR_KEY --from=en --to=ru,es,fr\n", binaryName)
		fmt.Fprintf(os.Stderr, "  %s --service=groq --api-key=YOUR_KEY --model=llama-3.3-70b-versatile --to=ru\n", binaryName)
		fmt.Fprintf(os.Stderr, "  %s --service=gemini --api-key=YOUR_KEY --model=gemini-2.0-flash --to=ru\n", binaryName)
		fmt.Fprintf(os.Stderr, "  %s --service=ollama --model=llama3.2 --to=ru,es\n", binaryName)
		fmt.Fprintf(os.Stderr, "  %s --service=lmstudio --model=your-model --to=ru,es\n", binaryName)
		fmt.Fprintf(os.Stderr, "  %s --config=.i18n-translator.yaml --to=ru,es\n", binaryName)
		fmt.Fprintf(os.Stderr, "  %s --dry-run --skip-existing --to=ru,es\n", binaryName)
		fmt.Fprintf(os.Stderr, "  %s --fail-fast --retries=1 --to=ru,es\n", binaryName)
		fmt.Fprintf(os.Stderr, "  %s --retries=3 --retry-delay=2s --timeout=60s --to=ru,es\n", binaryName)
		fmt.Fprintf(os.Stderr, "  %s init\n", binaryName)
		fmt.Fprintln(os.Stderr, "\nEnvironment variable fallbacks:")
		fmt.Fprintln(os.Stderr, "  I18N_SERVICE, I18N_API_KEY, I18N_MODEL, I18N_URL, I18N_FROM, I18N_TO, I18N_SOURCE, I18N_TARGET, I18N_CONCURRENCY, I18N_CONFIG, I18N_DRY_RUN, I18N_LIST_FILES, I18N_SKIP_EXISTING, I18N_FAIL_FAST, I18N_RETRIES, I18N_RETRY_DELAY, I18N_TIMEOUT")
	}

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	if *help {
		fs.Usage()
		return 0
	}

	if *version {
		fmt.Printf("i18n-translator %s\n", appVersion())
		return 0
	}

	explicitFlags := collectExplicitFlags(fs)
	configPath := resolveStringOption("config", *configFile, os.Getenv("I18N_CONFIG"), "", "", explicitFlags)
	cfg, loadedConfigPath, err := loadConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to load config: %v\n", err)
		return 1
	}
	if loadedConfigPath != "" {
		fmt.Printf("Using config file: %s\n", loadedConfigPath)
	}

	finalConcurrency, err := resolveIntOption("concurrency", "I18N_CONCURRENCY", *concurrency, os.Getenv("I18N_CONCURRENCY"), cfg.Concurrency, defaultConcurrency, explicitFlags)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to resolve concurrency: %v\n", err)
		return 1
	}
	if finalConcurrency < 1 {
		fmt.Fprintln(os.Stderr, "Error: concurrency must be greater than 0")
		return 1
	}

	finalRetries, err := resolveIntOption("retries", "I18N_RETRIES", *retries, os.Getenv("I18N_RETRIES"), cfg.Retries, defaultRetries, explicitFlags)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to resolve retries: %v\n", err)
		return 1
	}
	if finalRetries < 0 {
		fmt.Fprintln(os.Stderr, "Error: retries must be greater than or equal to 0")
		return 1
	}

	finalRetryDelay, err := resolveDurationOption("retry-delay", "I18N_RETRY_DELAY", *retryDelay, os.Getenv("I18N_RETRY_DELAY"), cfg.RetryDelay, defaultRetryDelay, explicitFlags)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to resolve retry-delay: %v\n", err)
		return 1
	}
	if finalRetryDelay < 0 {
		fmt.Fprintln(os.Stderr, "Error: retry-delay must be greater than or equal to 0")
		return 1
	}

	finalTimeout, err := resolveDurationOption("timeout", "I18N_TIMEOUT", *timeout, os.Getenv("I18N_TIMEOUT"), cfg.Timeout, defaultTimeout, explicitFlags)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to resolve timeout: %v\n", err)
		return 1
	}
	if finalTimeout < 0 {
		fmt.Fprintln(os.Stderr, "Error: timeout must be greater than or equal to 0")
		return 1
	}

	finalService := resolveStringOption("service", *service, os.Getenv("I18N_SERVICE"), cfg.Service, defaultService, explicitFlags)
	finalModel := resolveStringOption("model", *model, os.Getenv("I18N_MODEL"), cfg.Model, "", explicitFlags)
	finalSourceDir := resolveStringOption("source", *sourceDir, os.Getenv("I18N_SOURCE"), cfg.Source, defaultSourceDir, explicitFlags)
	finalTargetDir := resolveStringOption("target", *targetDir, os.Getenv("I18N_TARGET"), cfg.Target, defaultTargetDir, explicitFlags)
	finalSourceLang := resolveStringOption("from", *sourceLang, os.Getenv("I18N_FROM"), cfg.From, defaultSourceLang, explicitFlags)
	finalTargetLang := resolveStringOption("to", *targetLang, os.Getenv("I18N_TO"), cfg.To, defaultTargetLang, explicitFlags)
	finalAPIKey := resolveStringOption("api-key", *apiKey, os.Getenv("I18N_API_KEY"), cfg.apiKey(), "", explicitFlags)
	finalURL := resolveStringOption("url", *url, os.Getenv("I18N_URL"), cfg.URL, "", explicitFlags)
	finalDryRun, err := resolveBoolOption("dry-run", *dryRun, os.Getenv("I18N_DRY_RUN"), cfg.DryRun, false, explicitFlags)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to resolve dry-run: %v\n", err)
		return 1
	}
	finalListFiles, err := resolveBoolOption("list-files", *listFiles, os.Getenv("I18N_LIST_FILES"), cfg.ListFiles, false, explicitFlags)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to resolve list-files: %v\n", err)
		return 1
	}
	finalSkipExisting, err := resolveBoolOption("skip-existing", *skipExisting, os.Getenv("I18N_SKIP_EXISTING"), cfg.SkipExisting, false, explicitFlags)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to resolve skip-existing: %v\n", err)
		return 1
	}
	finalFailFast, err := resolveBoolOption("fail-fast", *failFast, os.Getenv("I18N_FAIL_FAST"), cfg.FailFast, false, explicitFlags)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to resolve fail-fast: %v\n", err)
		return 1
	}

	if _, err := os.Stat(finalSourceDir); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: source directory does not exist: %s\n", finalSourceDir)
		return 1
	}

	jobs, err := buildJobs(finalSourceDir, finalTargetDir, finalTargetLang)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to build jobs: %v\n", err)
		return 1
	}

	summary := runStats{Total: len(jobs)}
	jobs, skippedExistingCount := filterJobsByExisting(jobs, finalSkipExisting)
	summary.Skipped = skippedExistingCount
	if finalSkipExisting && skippedExistingCount > 0 {
		fmt.Printf("Skipping %d existing target files due to --skip-existing\n", skippedExistingCount)
	}

	if finalListFiles || finalDryRun {
		if finalListFiles {
			fmt.Println("Planned translation jobs:")
		}
		if finalDryRun {
			fmt.Println("Dry run mode: no files will be written.")
		}
		printPlannedJobs(jobs, finalSourceLang)
		fmt.Printf("Total jobs: %d\n", len(jobs))
		return 0
	}

	if len(jobs) == 0 {
		fmt.Println("No translation jobs to run.")
		printRunSummary(summary)
		return 0
	}

	needsAPIKey := finalService == "chatgpt" || finalService == "groq" || finalService == "gemini"
	if needsAPIKey && finalAPIKey == "" {
		fmt.Fprintf(os.Stderr, "Error: --api-key is required for %s service\n\n", finalService)
		fs.Usage()
		return 1
	}

	serviceImpl, err := createTranslationService(finalService, finalAPIKey, finalModel, finalURL, finalTimeout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to initialize %s service: %v\n", finalService, err)
		return 1
	}

	tr := &translator{service: serviceImpl}

	fmt.Printf("Starting translation from %s to languages: %s\n", finalSourceLang, finalTargetLang)
	fmt.Printf("Source directory: %s\n", finalSourceDir)
	fmt.Printf("Target directory: %s\n", finalTargetDir)

	result := runJobs(tr, jobs, finalSourceLang, finalConcurrency, runOptions{
		retries:    finalRetries,
		retryDelay: finalRetryDelay,
		failFast:   finalFailFast,
	})

	summary.Succeeded = result.Succeeded
	summary.Failed = result.Failed
	summary.Skipped += result.Skipped
	summary.Retried = result.Retried
	summary.FailFastTriggered = result.FailFastTriggered
	summary.FailFastFirstError = result.FailFastFirstError
	printRunSummary(summary)

	if summary.Failed > 0 {
		return 1
	}

	fmt.Println("Translation completed successfully!")
	return 0
}
