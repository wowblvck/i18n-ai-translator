package providers

import (
	"errors"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

type OpenAICompatibleService struct {
	client      *openai.Client
	model       string
	temperature float32
}

type openAICompatibleProviderConfig struct {
	baseURL        string
	defaultModel   string
	temperature    float32
	apiKeyRequired bool
	staticAPIKey   string
	allowCustomURL bool
	modelRequired  bool
}

var openAICompatibleProviders = map[string]openAICompatibleProviderConfig{
	"chatgpt": {
		baseURL:        "https://api.openai.com/v1",
		defaultModel:   "gpt-4o-mini",
		temperature:    0.2,
		apiKeyRequired: true,
	},
	"groq": {
		baseURL:        "https://api.groq.com/openai/v1",
		defaultModel:   "openai/gpt-oss-120b",
		temperature:    0.2,
		apiKeyRequired: true,
	},
	"gemini": {
		baseURL:        "https://generativelanguage.googleapis.com/v1beta/openai",
		defaultModel:   "gemini-3.8-flash",
		apiKeyRequired: true,
	},
	"deepseek": {
		baseURL:        "https://api.deepseek.com",
		defaultModel:   "deepseek-flash",
		apiKeyRequired: true,
	},
	"mistral": {
		baseURL:        "https://api.mistral.ai/v1",
		defaultModel:   "mistral-medium-3-5",
		apiKeyRequired: true,
	},
	"openrouter": {
		baseURL:        "https://openrouter.ai/api/v1",
		apiKeyRequired: true,
		modelRequired:  true,
	},
	"ollama": {
		baseURL:        "http://localhost:11434/v1",
		defaultModel:   "llama3.2",
		temperature:    0.2,
		staticAPIKey:   "ollama",
		allowCustomURL: true,
	},
	"lmstudio": {
		baseURL:        "http://localhost:1234/v1",
		temperature:    0.2,
		staticAPIKey:   "lm-studio",
		allowCustomURL: true,
	},
}

func OpenAICompatibleProvider(service, apiKey, baseURL, model string, timeout time.Duration) (*OpenAICompatibleService, error) {
	settings, ok := openAICompatibleProviders[service]
	if !ok {
		return nil, errors.New("unsupported OpenAI-compatible service: " + service)
	}

	if settings.apiKeyRequired && strings.TrimSpace(apiKey) == "" {
		return nil, errors.New(service + " api key is required")
	}
	if settings.staticAPIKey != "" {
		apiKey = settings.staticAPIKey
	}
	if strings.TrimSpace(model) == "" {
		model = settings.defaultModel
	}
	if settings.modelRequired && strings.TrimSpace(model) == "" {
		return nil, errors.New(service + " model is required; set it with --model")
	}
	if settings.allowCustomURL && strings.TrimSpace(baseURL) != "" {
		settings.baseURL = strings.TrimSpace(baseURL)
	}

	clientConfig := openai.DefaultConfig(apiKey)
	clientConfig.BaseURL = settings.baseURL
	applyHTTPTimeout(&clientConfig, timeout)

	return &OpenAICompatibleService{
		client:      openai.NewClientWithConfig(clientConfig),
		model:       model,
		temperature: settings.temperature,
	}, nil
}

func (s *OpenAICompatibleService) Translate(text, sourceLang, targetLang string) (string, error) {
	return performTranslation(s.client, s.model, text, sourceLang, targetLang, s.temperature)
}
