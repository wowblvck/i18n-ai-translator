package cli

import (
	"fmt"
	"time"

	"github.com/wowblvck/i18n-translator/internal/providers"
)

func createTranslationService(service, apiKey, model, baseURL string, timeout time.Duration) (TranslationService, error) {
	switch service {
	case "chatgpt":
		return providers.ChatGPTProvider(apiKey, model, timeout)
	case "groq":
		return providers.GroqProvider(apiKey, model, timeout)
	case "gemini":
		return providers.GeminiProvider(apiKey, model, timeout)
	case "ollama":
		return providers.OllamaProvider(baseURL, model, timeout)
	case "lmstudio":
		return providers.LMStudioProvider(baseURL, model, timeout)
	default:
		return nil, fmt.Errorf("unsupported translation service: %s (supported: chatgpt, groq, gemini, ollama, lmstudio)", service)
	}
}
