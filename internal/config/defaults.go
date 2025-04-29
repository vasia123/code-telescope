package config

// Константы значений по умолчанию для конфигурации
const (
	// FileSystem
	DefaultMaxDepth = 10

	// Parser
	DefaultMaxFileSize       = 1048576 // 1MB
	DefaultParsePrivateItems = false

	// LLM
	DefaultLLMProvider = "openai"
	DefaultLLMModel    = "gpt-4"
	DefaultTemperature = 0.3
	DefaultMaxTokens   = 1000
	DefaultBatchSize   = 5
	DefaultBatchDelay  = 1

	// Markdown
	DefaultIncludeTOC              = true
	DefaultIncludeFileInfo         = true
	DefaultMaxMethodDescriptionLen = 200
	DefaultGroupMethodsByType      = true
	DefaultCodeStyle               = "github"
)

// Константы для шаблонов включения/исключения файлов
var (
	// Шаблоны по умолчанию для включения файлов
	DefaultIncludePatterns = []string{
		"*.go", "*.js", "*.ts", "*.py", "*.java",
		"*.c", "*.cpp", "*.h", "*.hpp",
	}

	// Шаблоны по умолчанию для исключения файлов
	DefaultExcludePatterns = []string{
		"*_test.go", "test_*.py", "**/test/**",
		"**/node_modules/**", "**/vendor/**",
		"**/dist/**", "**/build/**",
	}

	// Поддерживаемые провайдеры ЛЛМ
	SupportedLLMProviders = []string{
		"openai",
		"anthropic",
	}

	// Поддерживаемые стили кода в Markdown
	SupportedCodeStyles = []string{
		"github",
		"default",
		"monokai",
		"solarized-dark",
		"solarized-light",
	}
)
