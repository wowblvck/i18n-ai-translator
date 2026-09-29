package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

func applyHTTPTimeout(config *openai.ClientConfig, timeout time.Duration) {
	if timeout <= 0 {
		return
	}
	config.HTTPClient = &http.Client{Timeout: timeout}
}

type translateItem struct {
	Original   string `json:"original"`
	Translated string `json:"translated"`
	Context    string `json:"context,omitempty"`
	Failure    string `json:"failure,omitempty"`
}

func performTranslation(client *openai.Client, model, text, sourceLang, targetLang string, temperature float32) (string, error) {
	systemPrompt, userPrompt, err := translationPrompt(text, sourceLang, targetLang)
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
			{Role: openai.ChatMessageRoleUser, Content: userPrompt},
		},
		Temperature: temperature,
	})
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", errors.New("empty response from provider")
	}

	return parseTranslation(resp.Choices[0].Message.Content)
}

func translationPrompt(text, sourceLang, targetLang string) (string, string, error) {
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
- The 'translated' value must be a complete, valid JSON document with the same object/array structure as the original text.
- Return only the JSON response; do not wrap it in Markdown code fences.

Return the translation as JSON.
`, sourceLang, targetLang))

	items := []translateItem{
		{Original: text},
	}
	payload, err := json.Marshal(items)
	if err != nil {
		return "", "", err
	}

	return systemPrompt, "inputLanguage=" + sourceLang + "; outputLanguage=" + targetLang + ";\n" + string(payload), nil
}

func parseTranslation(content string) (string, error) {
	content = strings.TrimSpace(content)
	jsonText := extractJSON(content)

	var out []translateItem
	if err := json.Unmarshal([]byte(jsonText), &out); err != nil {
		return "", fmt.Errorf("parse translation response: %w", err)
	}
	if len(out) == 0 {
		return "", errors.New("no translations returned")
	}
	if strings.TrimSpace(out[0].Translated) == "" {
		return "", errors.New("translated field is empty in response")
	}
	return out[0].Translated, nil
}

func postJSON(client *http.Client, endpoint string, headers map[string]string, payload any, provider string) ([]byte, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("encode %s request: %w", provider, err)
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("create %s request: %w", provider, err)
	}
	req.Header.Set("Content-Type", "application/json")
	for name, value := range headers {
		req.Header.Set(name, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send %s request: %w", provider, err)
	}
	defer func() { _ = resp.Body.Close() }()

	const maxResponseSize = 16 << 20
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseSize+1))
	if err != nil {
		return nil, fmt.Errorf("read %s response: %w", provider, err)
	}
	if len(body) > maxResponseSize {
		return nil, fmt.Errorf("%s response exceeds %d bytes", provider, maxResponseSize)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("%s API returned %s: %s", provider, resp.Status, strings.TrimSpace(string(body)))
	}
	return body, nil
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
