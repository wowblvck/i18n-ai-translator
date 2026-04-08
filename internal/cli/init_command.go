package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func runInitCommand(binaryName string, args []string) int {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	path := fs.String("path", defaultConfigPath, "Path for generated config file")
	example := fs.String("example", defaultExamplePath, "Path to example YAML template")
	force := fs.Bool("force", false, "Overwrite target config file if it already exists")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s init [options]\n\n", binaryName)
		fmt.Fprintln(os.Stderr, "Create a starter config file for i18n-ai-translator.")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Options:")
		fs.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nExamples:")
		fmt.Fprintf(os.Stderr, "  %s init\n", binaryName)
		fmt.Fprintf(os.Stderr, "  %s init --path=./configs/i18n.yaml\n", binaryName)
		fmt.Fprintf(os.Stderr, "  %s init --force\n", binaryName)
	}

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}
	if fs.NArg() > 0 {
		fmt.Fprintf(os.Stderr, "Error: unexpected arguments: %s\n", strings.Join(fs.Args(), ", "))
		return 1
	}

	targetPath := strings.TrimSpace(*path)
	if targetPath == "" {
		fmt.Fprintln(os.Stderr, "Error: path cannot be empty")
		return 1
	}
	if dir := filepath.Dir(targetPath); dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to create directory %s: %v\n", dir, err)
			return 1
		}
	}

	templateData, templateSource, err := loadInitTemplate(strings.TrimSpace(*example))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}
	if err := writeConfigFile(targetPath, templateData, *force); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	fmt.Printf("Created config: %s\n", targetPath)
	fmt.Printf("Template source: %s\n", templateSource)
	fmt.Println("Next step: set your api_key and run the translator.")
	return 0
}

func loadInitTemplate(examplePath string) ([]byte, string, error) {
	data, err := os.ReadFile(examplePath)
	if err == nil {
		return data, examplePath, nil
	}
	if os.IsNotExist(err) {
		return []byte(defaultConfigBody), "built-in template", nil
	}
	return nil, "", fmt.Errorf("failed to read example file %s: %w", examplePath, err)
}

func writeConfigFile(path string, content []byte, force bool) error {
	flags := os.O_WRONLY | os.O_CREATE
	if force {
		flags |= os.O_TRUNC
	} else {
		flags |= os.O_EXCL
	}

	file, err := os.OpenFile(path, flags, 0644)
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("file already exists: %s (use --force to overwrite)", path)
		}
		return fmt.Errorf("failed to open config file %s: %w", path, err)
	}

	if _, err := file.Write(content); err != nil {
		_ = file.Close()
		return fmt.Errorf("failed to write config file %s: %w", path, err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("failed to close config file %s: %w", path, err)
	}

	return nil
}
