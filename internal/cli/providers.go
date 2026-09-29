package cli

import (
	"fmt"
	"time"

	"github.com/wowblvck/i18n-translator/internal/providers"
)

func createTranslationService(service, apiKey, model, baseURL string, timeout time.Duration) (TranslationService, error) {
	switch service {
	case "anthropic":
		return providers.AnthropicProvider(apiKey, model, timeout)
	case "xai":
		return providers.XAIProvider(apiKey, model, timeout)
	case "chatgpt", "groq", "gemini", "deepseek", "mistral", "openrouter", "ollama", "lmstudio":
		return providers.OpenAICompatibleProvider(service, apiKey, baseURL, model, timeout)
	default:
		return nil, fmt.Errorf("unsupported translation service: %s (supported: chatgpt, groq, gemini, anthropic, deepseek, mistral, xai, openrouter, ollama, lmstudio)", service)
	}
}
