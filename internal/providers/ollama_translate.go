package providers

import (
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

type OllamaService struct {
	client *openai.Client
	model  string
}

func OllamaProvider(baseURL, model string) (*OllamaService, error) {
	if strings.TrimSpace(baseURL) == "" {
		baseURL = "http://localhost:11434/v1"
	}
	if strings.TrimSpace(model) == "" {
		model = "llama3.2"
	}
	config := openai.DefaultConfig("ollama")
	config.BaseURL = baseURL
	return &OllamaService{
		client: openai.NewClientWithConfig(config),
		model:  model,
	}, nil
}

func (o *OllamaService) Translate(text, sourceLang, targetLang string) (string, error) {
	return performTranslation(o.client, o.model, text, sourceLang, targetLang)
}
