package cli

import (
	"fmt"
	"os"
	"strings"
)

func (t *translator) TranslateFile(sourceFile, targetFile, sourceLang, targetLang string) error {
	data, err := os.ReadFile(sourceFile)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}

	translatedJSON, err := t.service.Translate(string(data), sourceLang, targetLang)
	if err != nil {
		return fmt.Errorf("failed to translate content: %w", err)
	}

	return os.WriteFile(targetFile, []byte(strings.TrimSpace(translatedJSON)), 0644)
}
