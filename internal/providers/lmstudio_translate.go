package providers

import (
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

type LMStudioService struct {
	client *openai.Client
	model  string
}

func LMStudioProvider(baseURL, model string) (*LMStudioService, error) {
	if strings.TrimSpace(baseURL) == "" {
		baseURL = "http://localhost:1234/v1"
	}
	config := openai.DefaultConfig("lm-studio")
	config.BaseURL = baseURL
	return &LMStudioService{
		client: openai.NewClientWithConfig(config),
		model:  model,
	}, nil
}

func (l *LMStudioService) Translate(text, sourceLang, targetLang string) (string, error) {
	return performTranslation(l.client, l.model, text, sourceLang, targetLang)
}
