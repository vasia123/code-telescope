package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config представляет основную конфигурацию приложения
type Config struct {
	FileSystem FileSystemConfig `yaml:"filesystem"`
	Parser     ParserConfig     `yaml:"parser"`
	LLM        LLMConfig        `yaml:"llm"`
	Markdown   MarkdownConfig   `yaml:"markdown"`
}

// FileSystemConfig содержит настройки для модуля файловой системы
type FileSystemConfig struct {
	IncludePatterns []string `yaml:"include_patterns"`
	ExcludePatterns []string `yaml:"exclude_patterns"`
	MaxDepth        int      `yaml:"max_depth"`
}

// ParserConfig содержит настройки для модуля парсинга кода
type ParserConfig struct {
	ParsePrivateMethods bool  `yaml:"parse_private_methods"`
	MaxFileSize         int64 `yaml:"max_file_size"`
}

// LLMConfig содержит настройки для модуля взаимодействия с ЛЛМ
type LLMConfig struct {
	Provider    string  `yaml:"provider"`
	Model       string  `yaml:"model"`
	Temperature float64 `yaml:"temperature"`
	MaxTokens   int     `yaml:"max_tokens"`
	BatchSize   int     `yaml:"batch_size"`
	BatchDelay  int     `yaml:"batch_delay"`
}

// MarkdownConfig содержит настройки для модуля генерации Markdown
type MarkdownConfig struct {
	IncludeTOC              bool   `yaml:"include_toc"`
	IncludeFileInfo         bool   `yaml:"include_file_info"`
	MaxMethodDescriptionLen int    `yaml:"max_method_description_length"`
	GroupMethodsByType      bool   `yaml:"group_methods_by_type"`
	CodeStyle               string `yaml:"code_style"`
}

// LoadConfig загружает конфигурацию из файла YAML
func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла конфигурации: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("ошибка парсинга YAML: %w", err)
	}

	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("ошибка валидации конфигурации: %w", err)
	}

	return &config, nil
}

// DefaultConfig возвращает конфигурацию по умолчанию
func DefaultConfig() *Config {
	return &Config{
		FileSystem: FileSystemConfig{
			IncludePatterns: []string{
				"*.go", "*.js", "*.ts", "*.py", "*.java",
				"*.c", "*.cpp", "*.h", "*.hpp",
			},
			ExcludePatterns: []string{
				"*_test.go", "test_*.py", "**/test/**",
				"**/node_modules/**", "**/vendor/**",
				"**/dist/**", "**/build/**",
			},
			MaxDepth: 10,
		},
		Parser: ParserConfig{
			ParsePrivateMethods: false,
			MaxFileSize:         1048576, // 1MB
		},
		LLM: LLMConfig{
			Provider:    "openai",
			Model:       "gpt-4",
			Temperature: 0.3,
			MaxTokens:   1000,
			BatchSize:   5,
			BatchDelay:  1,
		},
		Markdown: MarkdownConfig{
			IncludeTOC:              true,
			IncludeFileInfo:         true,
			MaxMethodDescriptionLen: 200,
			GroupMethodsByType:      true,
			CodeStyle:               "github",
		},
	}
}

// validateConfig проверяет корректность настроек конфигурации
func validateConfig(cfg *Config) error {
	// Проверка настроек LLM
	if cfg.LLM.Provider != "openai" && cfg.LLM.Provider != "anthropic" {
		return fmt.Errorf("неподдерживаемый провайдер ЛЛМ: %s", cfg.LLM.Provider)
	}

	if cfg.LLM.Temperature < 0 || cfg.LLM.Temperature > 1 {
		return fmt.Errorf("температура должна быть в диапазоне [0, 1], получено: %f", cfg.LLM.Temperature)
	}

	if cfg.LLM.BatchSize < 1 {
		return fmt.Errorf("размер пакета должен быть положительным, получено: %d", cfg.LLM.BatchSize)
	}

	// Проверка настроек файловой системы
	if cfg.FileSystem.MaxDepth < 1 {
		return fmt.Errorf("максимальная глубина должна быть положительной, получено: %d", cfg.FileSystem.MaxDepth)
	}

	// Другие проверки...

	return nil
}
