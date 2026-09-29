package providers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const xAIResponsesURL = "https://api.x.ai/v1/responses"

type XAIService struct {
	apiKey   string
	model    string
	client   *http.Client
	endpoint string
}

type xAIRequest struct {
	Model string       `json:"model"`
	Input []xAIMessage `json:"input"`
	Store bool         `json:"store"`
}

type xAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type xAIResponse struct {
	Status string `json:"status"`
	Error  *struct {
		Message string `json:"message"`
	} `json:"error"`
	IncompleteDetails *struct {
		Reason string `json:"reason"`
	} `json:"incomplete_details"`
	Output []struct {
		Type    string `json:"type"`
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	} `json:"output"`
}

func XAIProvider(apiKey, model string, timeout time.Duration) (*XAIService, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, errors.New("xai api key is required")
	}
	if strings.TrimSpace(model) == "" {
		model = "grok-4.7"
	}

	client := &http.Client{}
	if timeout > 0 {
		client.Timeout = timeout
	}
	return &XAIService{apiKey: apiKey, model: model, client: client, endpoint: xAIResponsesURL}, nil
}

func (s *XAIService) Translate(text, sourceLang, targetLang string) (string, error) {
	systemPrompt, userPrompt, err := translationPrompt(text, sourceLang, targetLang)
	if err != nil {
		return "", err
	}

	body, err := postJSON(s.client, s.endpoint, map[string]string{
		"Authorization": "Bearer " + s.apiKey,
	}, xAIRequest{
		Model: s.model,
		Input: []xAIMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Store: false,
	}, "xai")
	if err != nil {
		return "", err
	}

	var response xAIResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("decode xai response: %w", err)
	}
	if response.Status == "failed" && response.Error != nil {
		return "", fmt.Errorf("xai response failed: %s", response.Error.Message)
	}
	if response.Status != "" && response.Status != "completed" {
		if response.IncompleteDetails != nil && response.IncompleteDetails.Reason != "" {
			return "", fmt.Errorf("xai response %s: %s", response.Status, response.IncompleteDetails.Reason)
		}
		return "", fmt.Errorf("xai response status is %q", response.Status)
	}

	var content strings.Builder
	for _, item := range response.Output {
		if item.Type != "message" {
			continue
		}
		for _, block := range item.Content {
			if block.Type == "output_text" {
				content.WriteString(block.Text)
			}
		}
	}
	if strings.TrimSpace(content.String()) == "" {
		return "", errors.New("empty response from xai")
	}
	return parseTranslation(content.String())
}
