package providers

import (
	"errors"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

type GeminiService struct {
	client *openai.Client
	model  string
}

func GeminiProvider(apiKey, model string) (*GeminiService, error) {
	if apiKey == "" {
		return nil, errors.New("gemini api key is required")
	}
	if strings.TrimSpace(model) == "" {
		model = "gemini-2.0-flash"
	}
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://generativelanguage.googleapis.com/v1beta/openai"
	return &GeminiService{
		client: openai.NewClientWithConfig(config),
		model:  model,
	}, nil
}

func (g *GeminiService) Translate(text, sourceLang, targetLang string) (string, error) {
	return performTranslation(g.client, g.model, text, sourceLang, targetLang)
}
