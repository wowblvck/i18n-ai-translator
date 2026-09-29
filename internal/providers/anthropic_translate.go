package providers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const anthropicMessagesURL = "https://api.anthropic.com/v1/messages"

type AnthropicService struct {
	apiKey   string
	model    string
	client   *http.Client
	endpoint string
}

type anthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	System    string             `json:"system"`
	Messages  []anthropicMessage `json:"messages"`
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicResponse struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	StopReason string `json:"stop_reason"`
}

func AnthropicProvider(apiKey, model string, timeout time.Duration) (*AnthropicService, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, errors.New("anthropic api key is required")
	}
	if strings.TrimSpace(model) == "" {
		model = "claude-sonnet-5-5"
	}

	client := &http.Client{}
	if timeout > 0 {
		client.Timeout = timeout
	}
	return &AnthropicService{apiKey: apiKey, model: model, client: client, endpoint: anthropicMessagesURL}, nil
}

func (s *AnthropicService) Translate(text, sourceLang, targetLang string) (string, error) {
	systemPrompt, userPrompt, err := translationPrompt(text, sourceLang, targetLang)
	if err != nil {
		return "", err
	}

	body, err := postJSON(s.client, s.endpoint, map[string]string{
		"x-api-key":         s.apiKey,
		"anthropic-version": "2023-06-01",
	}, anthropicRequest{
		Model:     s.model,
		MaxTokens: 64000,
		System:    systemPrompt,
		Messages:  []anthropicMessage{{Role: "user", Content: userPrompt}},
	}, "anthropic")
	if err != nil {
		return "", err
	}

	var response anthropicResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("decode anthropic response: %w", err)
	}
	if response.StopReason == "max_tokens" {
		return "", errors.New("anthropic response reached the max_tokens limit")
	}

	var content strings.Builder
	for _, block := range response.Content {
		if block.Type == "text" {
			content.WriteString(block.Text)
		}
	}
	if strings.TrimSpace(content.String()) == "" {
		return "", errors.New("empty response from anthropic")
	}
	return parseTranslation(content.String())
}
