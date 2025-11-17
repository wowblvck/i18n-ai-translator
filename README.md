# 🌍 i18n AI Translator

Automatic translation tool for i18n JSON files using ChatGPT (OpenAI).

## 📋 Description

`i18n-ai-translator` is a CLI tool written in Go that automates the translation process of localization (i18n) files for your projects. It uses OpenAI's API (ChatGPT) for high-quality translations while preserving JSON structure, placeholders, and formatting.

### ✨ Key Features

- 🤖 **AI-powered translation**: Uses ChatGPT for contextually accurate translations
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
Get your OpenAI API key:
1. Sign up at [OpenAI Platform](https://platform.openai.com/)
2. Go to [API Keys](https://platform.openai.com/api-keys) section
3. Create a new API key

## 💻 Usage

### Basic example

```bash
i18n-ai-translator --api-key=YOUR_OPENAI_API_KEY --from=en --to=ru,es,fr
```

### Advanced examples

**Translation with custom directories:**
```bash
i18n-ai-translator \
  --source=./src/locales/en \
  --target=./src/locales \
  --api-key=YOUR_API_KEY \
  --from=en \
  --to=ru,uk,pl
```

**Using a specific model:**
```bash
i18n-ai-translator \
  --api-key=YOUR_API_KEY \
  --model=gpt-4o \
  --from=en \
  --to=de,it,pt
```

**Configuring concurrency:**
```bash
i18n-ai-translator \
  --api-key=YOUR_API_KEY \
  --from=en \
  --to=ja,ko,zh \
  --concurrency=8
```

**Translating a single file:**
```bash
i18n-ai-translator \
  --source=./locales/en/common.json \
  --target=./locales \
  --api-key=YOUR_API_KEY \
  --from=en \
  --to=es
```

**Using environment variable for API key:**
```bash
export OPENAI_API_KEY="your-api-key-here"
i18n-ai-translator --from=en --to=ru,es
```

## 🔧 Command-line options

| Parameter | Description | Default value |
|----------|----------|----------------------|
| `--api-key` | OpenAI API key (required) | - |
| `--from` | Source language code | `en` |
| `--to` | Target language codes (comma-separated) | `es,fr,de` |
| `--source` | Source directory or file with translations | `./locales/en` |
| `--target` | Target directory for translations | `./locales` |
| `--model` | ChatGPT model | `gpt-4o-mini` |
| `--service` | Translation service | `chatgpt` |
| `--concurrency` | Number of concurrent workers | `4` |
| `--help` | Show help message | - |
| `--version` | Show version | - |

### 🌍 Supported languages

The tool supports all languages that ChatGPT works with. Use standard ISO 639-1 language codes:

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

## ⚠️ Limitations

- **API Key required**: Active OpenAI API key with available credits
- **Internet connection**: Required for API communication
- **Cost**: Translation costs depend on:
  - Text volume
  - Number of target languages
  - Chosen model
  - API pricing (check [OpenAI Pricing](https://openai.com/api/pricing/))
- **Rate limits**: OpenAI API has rate limits based on your account tier
- **Quality**: While AI translations are generally good, manual review is recommended for critical content
- **Service**: Currently only supports ChatGPT (OpenAI). Other providers may be added in future

## ❓ FAQ

**Q: Do I need to pay for OpenAI API?**  
A: Yes, you need an OpenAI account with available credits. Check [pricing](https://openai.com/api/pricing/).

**Q: Can I use this with free ChatGPT account?**  
A: No, you need API access which is separate from the ChatGPT web interface.

**Q: How much does it cost to translate?**  
A: Costs depend on text volume and model. For example, translating 100KB of text with `gpt-4o-mini` typically costs $0.01-0.05.

**Q: Can I translate from multiple source languages?**  
A: Currently, you specify one source language per run. You can run the tool multiple times for different sources.

**Q: Does it work offline?**  
A: No, internet connection is required to communicate with OpenAI API.

**Q: Will it overwrite existing translations?**  
A: Yes, existing target files will be overwritten. Use version control to track changes.

**Q: Can I customize the translation prompts?**  
A: Not currently via CLI, but you can modify the code in `internal/providers/chatgpt_translate.go`.

**Q: What if translation fails?**  
A: The tool logs errors and continues with other files. Check the console output for details.
