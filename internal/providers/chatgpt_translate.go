package providers

import (
	"errors"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

type ChatGPTService struct {
	client *openai.Client
	model  string
}

func ChatGPTProvider(apiKey string, model string) (*ChatGPTService, error) {
	if apiKey == "" {
		return nil, errors.New("chatgpt api key is required")
	}
	if strings.TrimSpace(model) == "" {
		model = "gpt-4o-mini"
	}
	return &ChatGPTService{
		client: openai.NewClient(apiKey),
		model:  model,
	}, nil
}

func (c *ChatGPTService) Translate(text, sourceLang, targetLang string) (string, error) {
	return performTranslation(c.client, c.model, text, sourceLang, targetLang)
}
