package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunInitCommandCreatesConfigFromExample(t *testing.T) {
	tempDir := t.TempDir()
	examplePath := filepath.Join(tempDir, "example.yaml")
	targetPath := filepath.Join(tempDir, "generated.yaml")
	example := "service: ollama\nfrom: en\nto: ru\n"

	if err := os.WriteFile(examplePath, []byte(example), 0644); err != nil {
		t.Fatalf("failed to write example: %v", err)
	}

	code := runInitCommand("i18n-ai-translator", []string{"--path", targetPath, "--example", examplePath})
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}

	got, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatalf("failed to read generated file: %v", err)
	}
	if string(got) != example {
		t.Fatalf("generated file mismatch\nexpected:\n%s\ngot:\n%s", example, string(got))
	}
}

func TestRunInitCommandRejectsOverwriteWithoutForce(t *testing.T) {
	tempDir := t.TempDir()
	targetPath := filepath.Join(tempDir, "config.yaml")

	if err := os.WriteFile(targetPath, []byte("existing"), 0644); err != nil {
		t.Fatalf("failed to seed existing file: %v", err)
	}

	code := runInitCommand("i18n-ai-translator", []string{"--path", targetPath})
	if code == 0 {
		t.Fatalf("expected non-zero exit code when file exists without --force")
	}
}

func TestRunInitCommandForceOverwrite(t *testing.T) {
	tempDir := t.TempDir()
	examplePath := filepath.Join(tempDir, "example.yaml")
	targetPath := filepath.Join(tempDir, "config.yaml")

	if err := os.WriteFile(examplePath, []byte("service: chatgpt\n"), 0644); err != nil {
		t.Fatalf("failed to write example: %v", err)
	}
	if err := os.WriteFile(targetPath, []byte("old"), 0644); err != nil {
		t.Fatalf("failed to seed existing file: %v", err)
	}

	code := runInitCommand("i18n-ai-translator", []string{"--path", targetPath, "--example", examplePath, "--force"})
	if code != 0 {
		t.Fatalf("expected exit code 0 with --force, got %d", code)
	}

	got, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatalf("failed to read overwritten file: %v", err)
	}
	if string(got) != "service: chatgpt\n" {
		t.Fatalf("unexpected file content after force overwrite: %q", string(got))
	}
}
