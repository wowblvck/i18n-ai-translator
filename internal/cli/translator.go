package cli

import (
	"encoding/json"
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

	translatedJSON = strings.TrimSpace(translatedJSON)
	if !json.Valid([]byte(translatedJSON)) {
		return fmt.Errorf("translated content is not valid JSON")
	}

	return os.WriteFile(targetFile, []byte(translatedJSON), 0644)
}
