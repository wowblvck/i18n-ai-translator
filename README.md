# 🌍 i18n AI Translator

[![npm version](https://img.shields.io/npm/v/%40wowblvck%2Fi18n-ai-translator)](https://www.npmjs.com/package/@wowblvck/i18n-ai-translator)
[![License](https://img.shields.io/github/license/wowblvck/i18n-ai-translator)](https://github.com/wowblvck/i18n-ai-translator/blob/main/LICENSE.md)
[![npm downloads](https://img.shields.io/npm/dm/%40wowblvck%2Fi18n-ai-translator)](https://www.npmjs.com/package/@wowblvck/i18n-ai-translator)
[![GitHub stars](https://img.shields.io/github/stars/wowblvck/i18n-ai-translator?style=social)](https://github.com/wowblvck/i18n-ai-translator/stargazers)

Automatic translation tool for i18n JSON files using AI providers (OpenAI, Anthropic, DeepSeek, Mistral, Groq, Gemini, xAI, OpenRouter) and local models (Ollama, LM Studio).

## 📋 Description

`i18n-ai-translator` is a CLI tool written in Go that automates translation of localization (i18n) files. The full contents of each `.json` file are sent to the selected model as text. The CLI asks the model to preserve the JSON structure, placeholders, and formatting, extracts the translated text from the model's JSON response, and checks that the translated text is valid JSON before writing it.

### ✨ Key Features

- 🤖 **Multi-provider support**: OpenAI, Anthropic Claude, DeepSeek, Mistral, Groq, Gemini, xAI Grok, OpenRouter, Ollama, LM Studio
- 🏠 **Local inference**: Run translations offline via Ollama or LM Studio
- 📁 **Batch processing**: Translates multiple files and languages simultaneously
- ⚡ **Parallel processing**: Configurable number of concurrent workers
- 🔒 **Format-aware prompts**: Asks the model to preserve placeholders (`{{variable}}`), HTML tags, and i18n functions (`$t(key)`)
- 🌐 **Multi-language support**: Support for multiple target languages in a single run
- 📦 **Easy installation**: Pre-built binaries for Linux, macOS, and Windows

## 🚀 Installation

### Via npm/yarn

Requires Node.js 22 or newer. Prebuilt binaries are provided for macOS 13+ (x64 and ARM64), Linux (x64 and ARM64), and Windows (x64 and ARM64).

```bash
npm install -g @wowblvck/i18n-ai-translator
# or
yarn global add @wowblvck/i18n-ai-translator
```

### Via Go

Requires Go 1.26.8 or newer.

```bash
go install github.com/wowblvck/i18n-translator/cmd@latest
```

### From source
```bash
git clone https://github.com/wowblvck/i18n-ai-translator.git
cd i18n-ai-translator
npm run build
```

## 🔑 Setup

### ChatGPT (OpenAI)
1. Sign up at [OpenAI Platform](https://platform.openai.com/)
2. Go to [API Keys](https://platform.openai.com/api-keys) section
3. Create a new API key

### Groq
1. Sign up at [Groq Console](https://console.groq.com/)
2. Go to API Keys section
3. Create a new API key

### Gemini
1. Sign up at [Google AI Studio](https://aistudio.google.com/)
2. Go to [API keys](https://aistudio.google.com/app/apikey)
3. Create a new API key

### Anthropic Claude
1. Create an account at [Anthropic Console](https://console.anthropic.com/)
2. Create an API key in [Settings → API keys](https://console.anthropic.com/settings/keys)

### DeepSeek
1. Create an account at [DeepSeek Platform](https://platform.deepseek.com/)
2. Create an API key in [API Keys](https://platform.deepseek.com/api_keys)

### Mistral
1. Create an account at [Mistral La Plateforme](https://console.mistral.ai/)
2. Create an API key in [API keys](https://console.mistral.ai/api-keys/)

### xAI Grok
1. Create an account at [xAI Console](https://console.x.ai/)
2. Create an API key on the API Keys page

### OpenRouter
1. Create an account at [OpenRouter](https://openrouter.ai/)
2. Create an API key in [Settings → Keys](https://openrouter.ai/settings/keys)
3. Choose a model from the [OpenRouter model catalog](https://openrouter.ai/models)

### Ollama (local, no API key required)
1. Install [Ollama](https://ollama.com/)
2. Pull a model: `ollama pull llama3.2`
3. Start the server: `ollama serve`

### LM Studio (local, no API key required)
1. Install [LM Studio](https://lmstudio.ai/)
2. Download a model in the app
3. Start the local server in the "Local Server" tab

## 💻 Usage

### Basic examples

**ChatGPT:**
```bash
i18n-ai-translator --service=chatgpt --api-key=YOUR_OPENAI_KEY --from=en --to=ru,es,fr
```

**Groq:**
```bash
i18n-ai-translator --service=groq --api-key=YOUR_GROQ_KEY --from=en --to=ru,es,fr
```

**Gemini:**
```bash
i18n-ai-translator --service=gemini --api-key=YOUR_GEMINI_KEY --from=en --to=ru,es,fr
```

**Anthropic Claude:**
```bash
i18n-ai-translator --service=anthropic --api-key=YOUR_ANTHROPIC_KEY --from=en --to=ru,es,fr
```

**DeepSeek:**
```bash
i18n-ai-translator --service=deepseek --api-key=YOUR_DEEPSEEK_KEY --from=en --to=ru,es,fr
```

**Mistral:**
```bash
i18n-ai-translator --service=mistral --api-key=YOUR_MISTRAL_KEY --from=en --to=ru,es,fr
```

**xAI Grok:**
```bash
i18n-ai-translator --service=xai --api-key=YOUR_XAI_KEY --from=en --to=ru,es,fr
```

**OpenRouter** (choose a model from its catalog):
```bash
i18n-ai-translator --service=openrouter --api-key=YOUR_OPENROUTER_KEY --model=provider/model-id --from=en --to=ru,es,fr
```

**Ollama (local):**
```bash
i18n-ai-translator --service=ollama --model=llama3.2 --from=en --to=ru,es,fr
```

**LM Studio (local):**
```bash
i18n-ai-translator --service=lmstudio --model=your-loaded-model --from=en --to=ru,es,fr
```

### Advanced examples

**Translation with custom directories:**
```bash
i18n-ai-translator \
  --service=chatgpt \
  --source=./src/locales/en \
  --target=./src/locales \
  --api-key=YOUR_API_KEY \
  --from=en \
  --to=ru,uk,pl
```

**Using a specific model:**
```bash
i18n-ai-translator \
  --service=chatgpt \
  --api-key=YOUR_API_KEY \
  --model=gpt-4o \
  --from=en \
  --to=de,it,pt
```

**Configuring concurrency:**
```bash
i18n-ai-translator \
  --service=groq \
  --api-key=YOUR_GROQ_KEY \
  --from=en \
  --to=ja,ko,zh \
  --concurrency=8
```

**Translating a single file:**
```bash
i18n-ai-translator \
  --service=ollama \
  --model=llama3.2 \
  --source=./locales/en/common.json \
  --target=./locales \
  --from=en \
  --to=es
```

**Ollama with custom server URL:**
```bash
i18n-ai-translator \
  --service=ollama \
  --url=http://192.168.1.10:11434/v1 \
  --model=llama3.2 \
  --from=en \
  --to=ru
```

**Using config file:**
```bash
i18n-ai-translator --config=.i18n-translator.yaml
```

**Initialize config file:**
```bash
i18n-ai-translator init
```

**Dry run with existing-file skip:**
```bash
i18n-ai-translator --dry-run --skip-existing --to=ru,es
```

**List planned files only:**
```bash
i18n-ai-translator --list-files --to=ru,es
```

**Retries and timeout for unstable API responses:**
```bash
i18n-ai-translator --retries=3 --retry-delay=2s --timeout=60s --to=ru,es
```

**Fail fast (stop on first error):**
```bash
i18n-ai-translator --fail-fast --retries=1 --to=ru,es
```

Example file in repository: `.i18n-translator.example.yaml`.

### `init` command

Generate a starter config file in your project:

```bash
i18n-ai-translator init
```

Options:
- `--path` path to output config file (default: `.i18n-translator.yaml`)
- `--example` source template path (default: `.i18n-translator.example.yaml`)
- `--force` overwrite target file if it exists

The repository includes `.i18n-translator.example.yaml`. If the template file is unavailable (as in the npm package), `init` uses the equivalent built-in template.

## 🔧 Command-line options

| Parameter | Description | Default value |
|----------|----------|----------------------|
| `--service` | Translation service: `chatgpt`, `anthropic`, `deepseek`, `mistral`, `groq`, `gemini`, `xai`, `openrouter`, `ollama`, `lmstudio` | `chatgpt` |
| `--api-key` | API key (required for cloud services) | - |
| `--model` | Model name (see defaults per provider below) | - |
| `--url` | Base URL for `ollama` or `lmstudio` | see below |
| `--from` | Source language code | `en` |
| `--to` | Target language codes (comma-separated) | `es,fr,de` |
| `--source` | Source directory or file with translations | `./locales/en` |
| `--target` | Target directory for translations | `./locales` |
| `--concurrency` | Number of concurrent workers | `4` |
| `--config` | Path to YAML config file | auto-load `.i18n-translator.yaml/.yml` |
| `--dry-run` | Preview jobs without writing translated files | `false` |
| `--list-files` | Print planned source/target jobs and exit | `false` |
| `--skip-existing` | Skip jobs if target file already exists | `false` |
| `--fail-fast` | Stop scheduling new jobs after first translation error | `false` |
| `--retries` | Retries for temporary provider errors (429/5xx/timeout) | `0` |
| `--retry-delay` | Base delay between retries (`time.Duration`) | `1s` |
| `--timeout` | Per-request provider timeout (`time.Duration`), `0` disables | `0s` |
| `--help` | Show help message | - |
| `--version` | Show version | - |

### Config and environment fallbacks

Parameters are resolved in this order:
1. CLI flags
2. Environment variables
3. Config file values
4. Built-in defaults

Supported environment variables:
`I18N_SERVICE`, `I18N_API_KEY`, `I18N_MODEL`, `I18N_URL`, `I18N_FROM`, `I18N_TO`, `I18N_SOURCE`, `I18N_TARGET`, `I18N_CONCURRENCY`, `I18N_CONFIG`, `I18N_DRY_RUN`, `I18N_LIST_FILES`, `I18N_SKIP_EXISTING`, `I18N_FAIL_FAST`, `I18N_RETRIES`, `I18N_RETRY_DELAY`, `I18N_TIMEOUT`.

Example config file (`.i18n-translator.yaml`):
```yaml
service: gemini
api_key: your-api-key
from: en
to: ru,es,fr
source: ./locales/en
target: ./locales
concurrency: 6
dry_run: false
list_files: false
skip_existing: false
fail_fast: false
retries: 2
retry_delay: 1s
timeout: 60s
```

In a source checkout, you can start from `.i18n-translator.example.yaml` and copy it to `.i18n-translator.yaml`.

### Exit codes and summary

After a translation run, the CLI prints a summary with `Total`, `Succeeded`, `Failed`, `Skipped`, and `Retried` (jobs retried at least once). Preview and planning modes (`--dry-run` and `--list-files`) print their plans instead.
When `--fail-fast` is enabled, the summary also shows whether fail-fast was triggered and the first error that caused it.

- Exit code `0`: no job or file-write error was reported, and each translated result was valid JSON. The CLI does not compare the result's structure or keys with the source, so review generated files before using them.
- Exit code `1`: one or more jobs failed.

### Provider defaults

| Service | Default model | Default URL | API key |
|---------|--------------|-------------|---------|
| `chatgpt` | `gpt-4o-mini` | `https://api.openai.com/v1` | Required |
| `groq` | `openai/gpt-oss-120b` | `https://api.groq.com/openai/v1` | Required |
| `gemini` | `gemini-3.8-flash` | `https://generativelanguage.googleapis.com/v1beta/openai` | Required |
| `anthropic` | `claude-sonnet-5-5` | `https://api.anthropic.com` | Required |
| `deepseek` | `deepseek-flash` | `https://api.deepseek.com` | Required |
| `mistral` | `mistral-medium-3-5` | `https://api.mistral.ai/v1` | Required |
| `xai` | `grok-4.7` | `https://api.x.ai/v1` | Required |
| `openrouter` | *(specify via `--model`)* | `https://openrouter.ai/api/v1` | Required |
| `ollama` | `llama3.2` | `http://localhost:11434/v1` | Not required |
| `lmstudio` | *(specify via `--model`)* | `http://localhost:1234/v1` | Not required |

### 🌍 Supported languages

Language availability depends on the selected model. The CLI passes the `--from` and `--to` values to the model without validating language codes; the examples below use common ISO 639-1 codes:

| Code | Language | Code | Language |
|-----|------|-----|------|
| `en` | English | `ru` | Russian |
| `es` | Spanish | `uk` | Ukrainian |
| `fr` | French | `pl` | Polish |
| `de` | German | `ja` | Japanese |
| `it` | Italian | `ko` | Korean |
| `pt` | Portuguese | `zh` | Chinese |
| `nl` | Dutch | `ar` | Arabic |
| `tr` | Turkish | `hi` | Hindi |
| `sv` | Swedish | `cs` | Czech |
| `da` | Danish | `fi` | Finnish |
| `no` | Norwegian | `el` | Greek |
| `he` | Hebrew | `id` | Indonesian |
| `th` | Thai | `vi` | Vietnamese |

And many more...

## 📝 File format

The CLI accepts files whose names end in `.json`. It sends each complete file to the model as text, parses the model's response envelope to extract the translated text, and checks that text is valid JSON before writing the target file. If validation fails, the job reports an error and does not write the target file. The CLI does not compare the output structure or keys with the source, so review generated files before using them.

**Source file (en/common.json):**
```json
{
  "welcome": "Welcome to our app!",
  "greeting": "Hello, {{username}}!",
  "button": {
    "submit": "Submit",
    "cancel": "Cancel"
  },
  "message": "You have $t(notifications) new messages",
  "html_content": "Click the <button>Start</button> button to begin",
  "multiline": "First line{{NEWLINE}}Second line"
}
```

**Result (ru/common.json):**
```json
{
  "welcome": "Добро пожаловать в наше приложение!",
  "greeting": "Привет, {{username}}!",
  "button": {
    "submit": "Отправить",
    "cancel": "Отмена"
  },
  "message": "У вас $t(notifications) новых сообщений",
  "html_content": "Нажмите кнопку <button>Начать</button> чтобы начать",
  "multiline": "Первая строка{{NEWLINE}}Вторая строка"
}
```

## 🔒 Translation features

The translation prompt asks the model to handle special elements as follows. The model may still change them, so check the output, especially for production or user-facing text:

- **Placeholders**: Asked to keep `{{variable}}`, `{{username}}`, and `{{count}}` unchanged
- **i18n functions**: Asked to preserve `$t(key)` and `$t(notifications)`
- **HTML tags**: Asked to preserve tag names, attributes, and structure
- **Formatting**: Asked to preserve whitespace, case, and structure
- **Special characters**: Asked to keep `{{NEWLINE}}` instead of converting it to `\n`
- **Length**: Asked to keep roughly the same text length
- **Nested objects**: The complete JSON text is sent as one string; the CLI does not traverse object values itself

## 🛠️ Development

### Requirements

- Go 1.26.8+
- Node.js 22+

### Building the project

```bash
# Clone the repository
git clone https://github.com/wowblvck/i18n-ai-translator.git
cd i18n-ai-translator

# Install dependencies
go mod download

# Build for current platform
npm run build

# Build for all platforms
npm run build:all

# Build for specific platform
npm run build:linux    # Linux x64 and ARM64
npm run build:macos    # macOS x64 and ARM64
npm run build:windows  # Windows x64 and ARM64

# Clean build artifacts
npm run clean
```

The npm build scripts use POSIX-style shell syntax. On Windows, run them in WSL or configure npm to use a POSIX shell. To build natively in PowerShell, use `go build -o bin/i18n-ai-translator.exe ./cmd/main.go`.

### Linting (Go)

```bash
# Install golangci-lint (one-time)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run Go linters
npm run lint:go
```

## ⚠️ Limitations

- **API Key required**: Required for cloud services (`chatgpt`, `anthropic`, `deepseek`, `mistral`, `groq`, `gemini`, `xai`, and `openrouter`); not needed for `ollama` and `lmstudio`
- **Internet connection**: Required for cloud services; `ollama` and `lmstudio` work locally, though model downloads and setup may require internet access
- **Cost**: Cloud translation costs depend on:
  - Text volume
  - Number of target languages
  - Chosen model
  - Current provider pricing ([OpenAI](https://openai.com/api/pricing/), [Anthropic](https://platform.claude.com/docs/en/about-claude/pricing), [DeepSeek](https://api-docs.deepseek.com/quick_start/pricing), [Mistral](https://docs.mistral.ai/inference/pricing), [Groq](https://console.groq.com/docs/models), [Gemini](https://ai.google.dev/gemini-api/docs/pricing), [xAI](https://docs.x.ai/developers/pricing), [OpenRouter](https://openrouter.ai/models))
- **Rate limits**: Cloud providers have rate limits based on your account tier
- **Quality**: While AI translations are generally good, manual review is recommended for critical content
- **Local models**: Translation quality varies by model size — larger models generally produce better results

## ❓ FAQ

**Q: Do I need to pay to use this tool?**
A: Local inference through `ollama` or `lmstudio` does not incur this tool's cloud API charges, though it requires local hardware and a model whose license permits your use. Cloud providers require API access and may charge for usage.

**Q: Can I use this with free ChatGPT account?**
A: No, you need API access which is separate from the ChatGPT web interface.

**Q: How much does it cost to translate?**
A: It depends on the amount of input and output tokens and the selected model. Check the provider's current pricing; rates and model availability can change. Local inference has no per-request API charge from this tool.

**Q: Can I translate from multiple source languages?**
A: Currently, you specify one source language per run. You can run the tool multiple times for different sources.

**Q: Does it work offline?**
A: Yes, when using `ollama` or `lmstudio`. Cloud providers require internet access.

**Q: Will it overwrite existing translations?**
A: By default, yes. Use `--skip-existing` to leave existing target files untouched.

**Q: Can I customize the translation prompts?**
A: Not currently via CLI, but you can modify the code in `internal/providers/common.go`.

**Q: What if translation fails?**
A: The tool logs the error and continues scheduling other jobs by default. With `--fail-fast`, it stops scheduling new jobs after the first translation error; jobs already in progress may still finish. Check the console output for details.
