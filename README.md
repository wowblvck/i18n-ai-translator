# 🌍 i18n AI Translator

Automatic translation tool for i18n JSON files using AI providers (ChatGPT, Groq, Gemini, Ollama, LM Studio).

## 📋 Description

`i18n-ai-translator` is a CLI tool written in Go that automates the translation process of localization (i18n) files for your projects. It supports multiple AI providers — including cloud-based (OpenAI, Groq, Gemini) and local (Ollama, LM Studio) — while preserving JSON structure, placeholders, and formatting.

### ✨ Key Features

- 🤖 **Multi-provider support**: ChatGPT, Groq, Gemini, Ollama, LM Studio
- 🏠 **Local inference**: Run translations offline via Ollama or LM Studio
- 📁 **Batch processing**: Translates multiple files and languages simultaneously
- ⚡ **Parallel processing**: Configurable number of concurrent workers
- 🔒 **Format preservation**: Keeps placeholders (`{{variable}}`), HTML tags, and i18n functions (`$t(key)`) intact
- 🌐 **Multi-language support**: Support for multiple target languages in a single run
- 📦 **Easy installation**: Pre-built binaries for Linux, macOS, and Windows

## 🚀 Installation

### Via npm/yarn

```bash
npm install -g @wowblvck/i18n-ai-translator
# or
yarn global add @wowblvck/i18n-ai-translator
```

### Via Go
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

## 🔧 Command-line options

| Parameter | Description | Default value |
|----------|----------|----------------------|
| `--service` | Translation service: `chatgpt`, `groq`, `gemini`, `ollama`, `lmstudio` | `chatgpt` |
| `--api-key` | API key (required for `chatgpt`, `groq`, and `gemini`) | - |
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

You can start from `.i18n-translator.example.yaml` and copy it to `.i18n-translator.yaml`.

### Exit codes and summary

At the end of each run, the CLI prints a summary with `Total`, `Succeeded`, `Failed`, `Skipped`, and `Retried`.
When `--fail-fast` is enabled, summary also shows whether fail-fast was triggered and the first error that caused stop.

- Exit code `0`: all scheduled jobs finished successfully.
- Exit code `1`: one or more jobs failed.

### Provider defaults

| Service | Default model | Default URL | API key |
|---------|--------------|-------------|---------|
| `chatgpt` | `gpt-4o-mini` | `https://api.openai.com/v1` | Required |
| `groq` | `llama-3.3-70b-versatile` | `https://api.groq.com/openai/v1` | Required |
| `gemini` | `gemini-2.0-flash` | `https://generativelanguage.googleapis.com/v1beta/openai` | Required |
| `ollama` | `llama3.2` | `http://localhost:11434/v1` | Not required |
| `lmstudio` | *(specify via `--model`)* | `http://localhost:1234/v1` | Not required |

### 🌍 Supported languages

The tool supports all languages available in the chosen AI provider. Use standard ISO 639-1 language codes:

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

The tool works with any i18n JSON files:

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

The translator intelligently handles special elements:

- ✅ **Placeholders**: `{{variable}}`, `{{username}}`, `{{count}}` remain unchanged
- ✅ **i18n functions**: `$t(key)`, `$t(notifications)` are preserved as-is
- ✅ **HTML tags**: `<button>`, `<icon/>`, `<div>` are not modified
- ✅ **Formatting**: Whitespace, case sensitivity, and structure are preserved
- ✅ **Special characters**: `{{NEWLINE}}` is not converted to `\n`
- ✅ **Length preservation**: Keeps approximately the same text length
- ✅ **Nested objects**: Handles complex JSON structures with any nesting level

## 🛠️ Development

### Requirements

- Go 1.25.1+
- Node.js 16+

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
npm run build:linux    # Linux x64
npm run build:macos    # macOS x64 and ARM64
npm run build:windows  # Windows x64

# Clean build artifacts
npm run clean
```

### Linting (Go)

```bash
# Install golangci-lint (one-time)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run Go linters
npm run lint:go
```

## ⚠️ Limitations

- **API Key required**: Required for `chatgpt`, `groq`, and `gemini`; not needed for `ollama` and `lmstudio`
- **Internet connection**: Required for `chatgpt`, `groq`, and `gemini`; `ollama` and `lmstudio` work fully offline
- **Cost**: Cloud translation costs depend on:
  - Text volume
  - Number of target languages
  - Chosen model
  - API pricing ([OpenAI Pricing](https://openai.com/api/pricing/), [Groq Pricing](https://groq.com/pricing/))
- **Rate limits**: Cloud providers have rate limits based on your account tier
- **Quality**: While AI translations are generally good, manual review is recommended for critical content
- **Local models**: Translation quality varies by model size — larger models generally produce better results

## ❓ FAQ

**Q: Do I need to pay to use this tool?**
A: Not necessarily. `ollama` and `lmstudio` are completely free and work offline. For `chatgpt`, `groq`, and `gemini`, you need an account with API access.

**Q: Can I use this with free ChatGPT account?**
A: No, you need API access which is separate from the ChatGPT web interface.

**Q: How much does it cost to translate?**
A: Costs depend on text volume and model. For example, translating 100KB of text with `gpt-4o-mini` typically costs $0.01-0.05. Groq is significantly cheaper. Ollama and LM Studio are free.

**Q: Can I translate from multiple source languages?**
A: Currently, you specify one source language per run. You can run the tool multiple times for different sources.

**Q: Does it work offline?**
A: Yes, when using `ollama` or `lmstudio`. Cloud providers (`chatgpt`, `groq`, `gemini`) require internet access.

**Q: Will it overwrite existing translations?**
A: Yes, existing target files will be overwritten. Use version control to track changes.

**Q: Can I customize the translation prompts?**
A: Not currently via CLI, but you can modify the code in `internal/providers/common.go`.

**Q: What if translation fails?**
A: The tool logs errors and continues with other files. Check the console output for details.
