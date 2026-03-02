package providers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

type translateItem struct {
	Original   string `json:"original"`
	Translated string `json:"translated"`
	Context    string `json:"context,omitempty"`
	Failure    string `json:"failure,omitempty"`
}

func performTranslation(client *openai.Client, model, text, sourceLang, targetLang string) (string, error) {
	systemPrompt := strings.TrimSpace(fmt.Sprintf(`
Translate from %s to %s.

- Translate each object in the array.
- 'original' is the text to be translated.
- 'translated' must not be empty.
- 'context' is additional info if needed.
- 'failure' explains why the previous translation failed.
- Preserve text formatting, case sensitivity, whitespace, and keep roughly the same length.

Special Instructions:
- Preserve text formatting, case sensitivity, whitespace, and keep roughly the same length.
- Do NOT translate or modify placeholders like {{variableName}}; keep them exactly as-is.
- Do NOT add new placeholders or variables; keep the same count and names.
- Do NOT convert {{NEWLINE}} to \n.
- Do NOT translate or modify i18n function calls in the form $t(key); return them verbatim (e.g., $t(ago) stays $t(ago)).
- Do NOT translate or modify HTML/XML tags (e.g., <button>...</button>, <icon/>, <actionButton/>, <bracketsButton/>); preserve tag names, attributes, and structure.

Return the translation as JSON.
`, sourceLang, targetLang))

	items := []translateItem{
		{Original: text},
	}
	payload, err := json.Marshal(items)
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
			{Role: openai.ChatMessageRoleUser, Content: "inputLanguage=" + sourceLang + "; outputLanguage=" + targetLang + ";\n" + string(payload)},
		},
		Temperature: 0.2,
	})
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", errors.New("empty response from provider")
	}

	content := strings.TrimSpace(resp.Choices[0].Message.Content)
	jsonText := extractJSON(content)

	var out []translateItem
	if err := json.Unmarshal([]byte(jsonText), &out); err != nil {
		return "", err
	}
	if len(out) == 0 {
		return "", errors.New("no translations returned")
	}
	if strings.TrimSpace(out[0].Translated) == "" {
		return "", errors.New("translated field is empty in response")
	}
	return out[0].Translated, nil
}

func extractJSON(s string) string {
	if strings.HasPrefix(s, "```") {
		s = strings.TrimPrefix(s, "```json")
		s = strings.TrimPrefix(s, "```")
		s = strings.TrimSuffix(s, "```")
		s = strings.TrimSpace(s)
	}
	startArray := strings.Index(s, "[")
	startObj := strings.Index(s, "{")
	start := startObj
	if startArray >= 0 && (startArray < startObj || startObj < 0) {
		start = startArray
	}
	if start < 0 {
		return s
	}
	end := strings.LastIndex(s, "]")
	if end < 0 || end < start {
		end = strings.LastIndex(s, "}")
	}
	if end > start {
		return s[start : end+1]
	}
	return s[start:]
}
