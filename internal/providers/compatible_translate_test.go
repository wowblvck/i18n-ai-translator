package providers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	openai "github.com/sashabaranov/go-openai"
)

func TestOpenAICompatibleProviderDefaults(t *testing.T) {
	tests := []struct {
		service     string
		model       string
		temperature float32
		apiKey      string
	}{
		{service: "chatgpt", model: "gpt-4o-mini", temperature: 0.2, apiKey: "key"},
		{service: "groq", model: "openai/gpt-oss-120b", temperature: 0.2, apiKey: "key"},
		{service: "gemini", model: "gemini-3.8-flash", apiKey: "key"},
		{service: "deepseek", model: "deepseek-flash", apiKey: "key"},
		{service: "mistral", model: "mistral-medium-3-5", apiKey: "key"},
		{service: "openrouter", model: "provider/model", apiKey: "key"},
		{service: "ollama", model: "llama3.2", temperature: 0.2},
		{service: "lmstudio", temperature: 0.2},
	}

	for _, test := range tests {
		t.Run(test.service, func(t *testing.T) {
			service, err := OpenAICompatibleProvider(test.service, test.apiKey, "", test.model, 0)
			if err != nil {
				t.Fatalf("OpenAICompatibleProvider() error = %v", err)
			}
			if service.model != test.model {
				t.Errorf("model = %q, want %q", service.model, test.model)
			}
			if service.temperature != test.temperature {
				t.Errorf("temperature = %v, want %v", service.temperature, test.temperature)
			}
		})
	}
}

func TestOpenAICompatibleProviderRequiresCredentialsAndModel(t *testing.T) {
	for _, service := range []string{"chatgpt", "groq", "gemini", "deepseek", "mistral", "openrouter"} {
		if _, err := OpenAICompatibleProvider(service, "", "", "model", 0); err == nil {
			t.Errorf("OpenAICompatibleProvider(%q) accepted a missing API key", service)
		}
	}
	if _, err := OpenAICompatibleProvider("openrouter", "key", "", "", 0); err == nil {
		t.Fatal("OpenRouter accepted a missing model")
	}
	if _, err := OpenAICompatibleProvider("unsupported", "key", "", "model", 0); err == nil {
		t.Fatal("unsupported service was accepted")
	}
}

func TestOpenAICompatibleServiceTranslate(t *testing.T) {
	modelResponse := translationEnvelope(t, `{"hello":"hola"}`)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat/completions" {
			t.Errorf("request path = %q, want /chat/completions", r.URL.Path)
		}

		var request struct {
			Model       string   `json:"model"`
			Temperature *float32 `json:"temperature"`
			Messages    []struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"messages"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Errorf("decode request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if request.Model != "test-model" {
			t.Errorf("model = %q, want test-model", request.Model)
		}
		if request.Temperature != nil {
			t.Errorf("temperature = %v, want omitted", *request.Temperature)
		}
		if len(request.Messages) != 2 || !strings.Contains(request.Messages[1].Content, "inputLanguage=en; outputLanguage=es;") {
			t.Errorf("unexpected messages: %+v", request.Messages)
		}

		_ = json.NewEncoder(w).Encode(map[string]any{
			"choices": []any{map[string]any{
				"message": map[string]string{"role": "assistant", "content": modelResponse},
			}},
		})
	}))
	defer server.Close()

	config := openai.DefaultConfig("test-key")
	config.BaseURL = server.URL
	service := &OpenAICompatibleService{
		client:      openai.NewClientWithConfig(config),
		model:       "test-model",
		temperature: 0,
	}
	got, err := service.Translate(`{"hello":"hello"}`, "en", "es")
	if err != nil {
		t.Fatalf("Translate() error = %v", err)
	}
	if got != `{"hello":"hola"}` {
		t.Errorf("Translate() = %q, want %q", got, `{"hello":"hola"}`)
	}
}

func TestAnthropicServiceTranslate(t *testing.T) {
	modelResponse := translationEnvelope(t, `{"hello":"hola"}`)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/messages" {
			t.Errorf("request path = %q, want /v1/messages", r.URL.Path)
		}
		if r.Header.Get("x-api-key") != "test-key" {
			t.Errorf("x-api-key = %q, want test-key", r.Header.Get("x-api-key"))
		}
		if r.Header.Get("anthropic-version") != "2023-06-01" {
			t.Errorf("anthropic-version = %q", r.Header.Get("anthropic-version"))
		}

		var request map[string]any
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Errorf("decode request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if request["model"] != "test-claude" || request["max_tokens"] != float64(64000) {
			t.Errorf("unexpected model/token settings: %+v", request)
		}
		if request["system"] == "" || request["temperature"] != nil {
			t.Errorf("unexpected Anthropic request fields: %+v", request)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"stop_reason": "end_turn",
			"content":     []any{map[string]string{"type": "text", "text": modelResponse}},
		})
	}))
	defer server.Close()

	service := &AnthropicService{
		apiKey:   "test-key",
		model:    "test-claude",
		client:   server.Client(),
		endpoint: server.URL + "/v1/messages",
	}
	got, err := service.Translate(`{"hello":"hello"}`, "en", "es")
	if err != nil {
		t.Fatalf("Translate() error = %v", err)
	}
	if got != `{"hello":"hola"}` {
		t.Errorf("Translate() = %q, want %q", got, `{"hello":"hola"}`)
	}
}

func TestXAIServiceTranslate(t *testing.T) {
	modelResponse := translationEnvelope(t, `{"hello":"hola"}`)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/responses" {
			t.Errorf("request path = %q, want /v1/responses", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Errorf("Authorization = %q", r.Header.Get("Authorization"))
		}

		var request map[string]any
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Errorf("decode request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if request["model"] != "test-grok" || request["store"] != false {
			t.Errorf("unexpected xAI request fields: %+v", request)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status": "completed",
			"output": []any{map[string]any{
				"type":    "message",
				"content": []any{map[string]string{"type": "output_text", "text": modelResponse}},
			}},
		})
	}))
	defer server.Close()

	service := &XAIService{
		apiKey:   "test-key",
		model:    "test-grok",
		client:   server.Client(),
		endpoint: server.URL + "/v1/responses",
	}
	got, err := service.Translate(`{"hello":"hello"}`, "en", "es")
	if err != nil {
		t.Fatalf("Translate() error = %v", err)
	}
	if got != `{"hello":"hola"}` {
		t.Errorf("Translate() = %q, want %q", got, `{"hello":"hola"}`)
	}
}

func translationEnvelope(t *testing.T, translated string) string {
	t.Helper()
	encoded, err := json.Marshal([]translateItem{{Original: "source", Translated: translated}})
	if err != nil {
		t.Fatalf("marshal translation envelope: %v", err)
	}
	return string(encoded)
}
