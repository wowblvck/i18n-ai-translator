package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"

	"github.com/wowblvck/i18n-translator/internal/providers"
)

const (
	defaultConcurrency = 4
	defaultService     = "chatgpt"
	defaultSourceDir   = "./locales/en"
	defaultTargetDir   = "./locales"
	defaultSourceLang  = "en"
	defaultTargetLang  = "es,fr,de"
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

type appConfig struct {
	Service      string `yaml:"service"`
	Model        string `yaml:"model"`
	Source       string `yaml:"source"`
	Target       string `yaml:"target"`
	From         string `yaml:"from"`
	To           string `yaml:"to"`
	APIKey       string `yaml:"api_key"`
	APIKeyAlt    string `yaml:"api-key"`
	URL          string `yaml:"url"`
	Concurrency  int    `yaml:"concurrency"`
	DryRun       bool   `yaml:"dry_run"`
	ListFiles    bool   `yaml:"list_files"`
	SkipExisting bool   `yaml:"skip_existing"`
}

func (c appConfig) apiKey() string {
	if strings.TrimSpace(c.APIKey) != "" {
		return strings.TrimSpace(c.APIKey)
	}
	return strings.TrimSpace(c.APIKeyAlt)
}

func loadConfig(configPath string) (appConfig, string, error) {
	configPath = strings.TrimSpace(configPath)
	paths := []string{".i18n-translator.yaml", ".i18n-translator.yml"}
	if configPath != "" {
		paths = []string{configPath}
	}

	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				if configPath != "" {
					return appConfig{}, "", fmt.Errorf("config file not found: %s", path)
				}
				continue
			}
			return appConfig{}, "", fmt.Errorf("failed to read config file %s: %w", path, err)
		}

		var cfg appConfig
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return appConfig{}, "", fmt.Errorf("failed to parse config file %s: %w", path, err)
		}
		return cfg, path, nil
	}

	return appConfig{}, "", nil
}

func collectExplicitFlags() map[string]bool {
	explicit := map[string]bool{}
	flag.Visit(func(f *flag.Flag) {
		explicit[f.Name] = true
	})
	return explicit
}

func resolveStringOption(flagName, flagValue, envValue, configValue, defaultValue string, explicitFlags map[string]bool) string {
	if explicitFlags[flagName] {
		return flagValue
	}
	if env := strings.TrimSpace(envValue); env != "" {
		return env
	}
	if cfg := strings.TrimSpace(configValue); cfg != "" {
		return cfg
	}
	return defaultValue
}

func resolveIntOption(flagName string, flagValue int, envValue string, configValue, defaultValue int, explicitFlags map[string]bool) (int, error) {
	if explicitFlags[flagName] {
		return flagValue, nil
	}

	if env := strings.TrimSpace(envValue); env != "" {
		v, err := strconv.Atoi(env)
		if err != nil {
			return 0, fmt.Errorf("I18N_CONCURRENCY must be a valid integer: %w", err)
		}
		return v, nil
	}

	if configValue != 0 {
		return configValue, nil
	}

	return defaultValue, nil
}

func resolveBoolOption(flagName string, flagValue bool, envValue string, configValue, defaultValue bool, explicitFlags map[string]bool) (bool, error) {
	if explicitFlags[flagName] {
		return flagValue, nil
	}

	if env := strings.TrimSpace(envValue); env != "" {
		v, err := strconv.ParseBool(env)
		if err != nil {
			return false, fmt.Errorf("%s must be a valid boolean: %w", flagName, err)
		}
		return v, nil
	}

	if configValue {
		return true, nil
	}

	return defaultValue, nil
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

func main() {
	var (
		concurrency  = flag.Int("concurrency", defaultConcurrency, "Number of parallel workers")
		model        = flag.String("model", "", "Model name (defaults: chatgpt=gpt-4o-mini, groq=llama-3.3-70b-versatile, gemini=gemini-2.0-flash, ollama=llama3.2)")
		sourceDir    = flag.String("source", defaultSourceDir, "Source language directory")
		targetDir    = flag.String("target", defaultTargetDir, "Target directory for translations")
		sourceLang   = flag.String("from", defaultSourceLang, "Source language code")
		targetLang   = flag.String("to", defaultTargetLang, "Target language codes (comma-separated)")
		apiKey       = flag.String("api-key", "", "API key (required for chatgpt, groq, and gemini)")
		url          = flag.String("url", "", "Base URL for ollama or lmstudio (e.g., http://localhost:11434/v1)")
		service      = flag.String("service", defaultService, "Translation service: chatgpt, groq, gemini, ollama, lmstudio")
		configFile   = flag.String("config", "", "Path to YAML config file (auto-loads .i18n-translator.yaml/.yml when omitted)")
		dryRun       = flag.Bool("dry-run", false, "Preview translation jobs without writing files")
		listFiles    = flag.Bool("list-files", false, "Print planned source/target files and exit")
		skipExisting = flag.Bool("skip-existing", false, "Skip jobs whose target file already exists")
		help         = flag.Bool("help", false, "Show help message")
		version      = flag.Bool("version", false, "Show version")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Automatic translation tool for i18n JSON files\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s --service=chatgpt --api-key=YOUR_KEY --from=en --to=ru,es,fr\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --service=groq --api-key=YOUR_KEY --model=llama-3.3-70b-versatile --to=ru\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --service=gemini --api-key=YOUR_KEY --model=gemini-2.0-flash --to=ru\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --service=ollama --model=llama3.2 --to=ru,es\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --service=lmstudio --model=your-model --to=ru,es\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --config=.i18n-translator.yaml --to=ru,es\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --dry-run --skip-existing --to=ru,es\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nEnvironment variable fallbacks:\n")
		fmt.Fprintf(os.Stderr, "  I18N_SERVICE, I18N_API_KEY, I18N_MODEL, I18N_URL, I18N_FROM, I18N_TO, I18N_SOURCE, I18N_TARGET, I18N_CONCURRENCY, I18N_CONFIG, I18N_DRY_RUN, I18N_LIST_FILES, I18N_SKIP_EXISTING\n")
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

	explicitFlags := collectExplicitFlags()
	configPath := resolveStringOption("config", *configFile, os.Getenv("I18N_CONFIG"), "", "", explicitFlags)
	cfg, loadedConfigPath, err := loadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	if loadedConfigPath != "" {
		fmt.Printf("Using config file: %s\n", loadedConfigPath)
	}

	finalConcurrency, err := resolveIntOption("concurrency", *concurrency, os.Getenv("I18N_CONCURRENCY"), cfg.Concurrency, defaultConcurrency, explicitFlags)
	if err != nil {
		log.Fatalf("Failed to resolve concurrency: %v", err)
	}
	if finalConcurrency < 1 {
		log.Fatalf("concurrency must be greater than 0")
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
		log.Fatalf("Failed to resolve dry-run: %v", err)
	}
	finalListFiles, err := resolveBoolOption("list-files", *listFiles, os.Getenv("I18N_LIST_FILES"), cfg.ListFiles, false, explicitFlags)
	if err != nil {
		log.Fatalf("Failed to resolve list-files: %v", err)
	}
	finalSkipExisting, err := resolveBoolOption("skip-existing", *skipExisting, os.Getenv("I18N_SKIP_EXISTING"), cfg.SkipExisting, false, explicitFlags)
	if err != nil {
		log.Fatalf("Failed to resolve skip-existing: %v", err)
	}

	if _, err := os.Stat(finalSourceDir); os.IsNotExist(err) {
		log.Fatalf("Source directory does not exist: %s", finalSourceDir)
	}

	jobs, err := buildJobs(finalSourceDir, finalTargetDir, finalTargetLang)
	if err != nil {
		log.Fatalf("Failed to build jobs: %v", err)
	}
	jobs, skippedExistingCount := filterJobsByExisting(jobs, finalSkipExisting)
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
		return
	}

	if len(jobs) == 0 {
		fmt.Println("No translation jobs to run.")
		return
	}

	needsAPIKey := finalService == "chatgpt" || finalService == "groq" || finalService == "gemini"
	if needsAPIKey && finalAPIKey == "" {
		fmt.Fprintf(os.Stderr, "Error: --api-key is required for %s service\n\n", finalService)
		flag.Usage()
		os.Exit(1)
	}

	var translationService TranslationService
	switch finalService {
	case "chatgpt":
		translationService, err = providers.ChatGPTProvider(finalAPIKey, finalModel)
		if err != nil {
			log.Fatalf("Failed to initialize ChatGPT service: %v", err)
		}
	case "groq":
		translationService, err = providers.GroqProvider(finalAPIKey, finalModel)
		if err != nil {
			log.Fatalf("Failed to initialize Groq service: %v", err)
		}
	case "gemini":
		translationService, err = providers.GeminiProvider(finalAPIKey, finalModel)
		if err != nil {
			log.Fatalf("Failed to initialize Gemini service: %v", err)
		}
	case "ollama":
		translationService, err = providers.OllamaProvider(finalURL, finalModel)
		if err != nil {
			log.Fatalf("Failed to initialize Ollama service: %v", err)
		}
	case "lmstudio":
		translationService, err = providers.LMStudioProvider(finalURL, finalModel)
		if err != nil {
			log.Fatalf("Failed to initialize LM Studio service: %v", err)
		}
	default:
		log.Fatalf("Unsupported translation service: %s (supported: chatgpt, groq, gemini, ollama, lmstudio)", finalService)
	}

	translator := &I18nTranslator{
		service:   translationService,
		sourceDir: finalSourceDir,
		targetDir: finalTargetDir,
	}

	fmt.Printf("Starting translation from %s to languages: %s\n", finalSourceLang, finalTargetLang)
	fmt.Printf("Source directory: %s\n", finalSourceDir)
	fmt.Printf("Target directory: %s\n", finalTargetDir)

	if err := runJobs(translator, jobs, finalSourceLang, finalConcurrency); err != nil {
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
