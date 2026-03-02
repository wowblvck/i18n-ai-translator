package providers

import (
	"errors"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

type GroqService struct {
	client *openai.Client
	model  string
}

func GroqProvider(apiKey, model string) (*GroqService, error) {
	if apiKey == "" {
		return nil, errors.New("groq api key is required")
	}
	if strings.TrimSpace(model) == "" {
		model = "llama-3.3-70b-versatile"
	}
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://api.groq.com/openai/v1"
	return &GroqService{
		client: openai.NewClientWithConfig(config),
		model:  model,
	}, nil
}

func (g *GroqService) Translate(text, sourceLang, targetLang string) (string, error) {
	return performTranslation(g.client, g.model, text, sourceLang, targetLang)
}
