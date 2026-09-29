package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type staticTranslationService struct {
	output string
}

func (s staticTranslationService) Translate(string, string, string) (string, error) {
	return s.output, nil
}

func TestTranslateFileRejectsInvalidJSONWithoutOverwritingTarget(t *testing.T) {
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "source.json")
	targetPath := filepath.Join(dir, "target.json")
	if err := os.WriteFile(sourcePath, []byte(`{"hello":"hello"}`), 0644); err != nil {
		t.Fatalf("write source file: %v", err)
	}
	if err := os.WriteFile(targetPath, []byte(`{"old":"translation"}`), 0644); err != nil {
		t.Fatalf("write existing target file: %v", err)
	}

	translator := &translator{service: staticTranslationService{output: `{"hello": }`}}
	err := translator.TranslateFile(sourcePath, targetPath, "en", "es")
	if err == nil || !strings.Contains(err.Error(), "translated content is not valid JSON") {
		t.Fatalf("TranslateFile() error = %v, want invalid JSON error", err)
	}

	got, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatalf("read existing target file: %v", err)
	}
	if string(got) != `{"old":"translation"}` {
		t.Errorf("target file changed after invalid response: %s", got)
	}
}

func TestTranslateFileWritesValidJSON(t *testing.T) {
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "source.json")
	targetPath := filepath.Join(dir, "target.json")
	if err := os.WriteFile(sourcePath, []byte(`{"hello":"hello"}`), 0644); err != nil {
		t.Fatalf("write source file: %v", err)
	}

	translator := &translator{service: staticTranslationService{output: "  {\n  \"hello\": \"hola\"\n}  "}}
	if err := translator.TranslateFile(sourcePath, targetPath, "en", "es"); err != nil {
		t.Fatalf("TranslateFile() error = %v", err)
	}
	got, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatalf("read target file: %v", err)
	}
	if string(got) != "{\n  \"hello\": \"hola\"\n}" {
		t.Errorf("target file = %q, want trimmed valid JSON", got)
	}
}
