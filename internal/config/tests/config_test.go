package tests

import (
	"os"
	"testing"

	"code-telescope/internal/config"

	"github.com/stretchr/testify/assert"
)

// createTempConfigFile создает временный файл конфигурации с заданным содержимым
func createTempConfigFile(t *testing.T, content string) *os.File {
	tmpfile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatalf("Не удалось создать временный файл: %v", err)
	}

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
		t.Fatalf("Не удалось записать во временный файл: %v", err)
	}

	if err := tmpfile.Close(); err != nil {
		os.Remove(tmpfile.Name())
		t.Fatalf("Не удалось закрыть временный файл: %v", err)
	}

	return tmpfile
}

// TestDefaultConfig проверяет корректность значений конфигурации по умолчанию
func TestDefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig()

	// Проверка значений по умолчанию для FileSystem
	assert.Equal(t, 10, cfg.FileSystem.MaxDepth, "Максимальная глубина сканирования должна быть 10")
	assert.NotEmpty(t, cfg.FileSystem.IncludePatterns, "Шаблоны включения не должны быть пустыми")
	assert.NotEmpty(t, cfg.FileSystem.ExcludePatterns, "Шаблоны исключения не должны быть пустыми")

	// Проверка значений по умолчанию для Parser
	assert.False(t, cfg.Parser.ParsePrivateMethods, "По умолчанию приватные методы не должны парситься")
	assert.Equal(t, int64(1048576), cfg.Parser.MaxFileSize, "Максимальный размер файла должен быть 1MB")

	// Проверка значений по умолчанию для LLM
	assert.Equal(t, "openai", cfg.LLM.Provider, "Провайдер ЛЛМ по умолчанию должен быть OpenAI")
	assert.Equal(t, "gpt-4", cfg.LLM.Model, "Модель по умолчанию должна быть gpt-4")
	assert.Equal(t, 0.3, cfg.LLM.Temperature, "Температура по умолчанию должна быть 0.3")
	assert.Equal(t, 1000, cfg.LLM.MaxTokens, "Максимальное количество токенов по умолчанию должно быть 1000")
	assert.Equal(t, 5, cfg.LLM.BatchSize, "Размер пакета должен быть 5")
	assert.Equal(t, 1, cfg.LLM.BatchDelay, "Задержка между пакетами должна быть 1")

	// Проверка значений по умолчанию для Markdown
	assert.True(t, cfg.Markdown.IncludeTOC, "По умолчанию оглавление должно включаться")
	assert.True(t, cfg.Markdown.IncludeFileInfo, "По умолчанию информация о файле должна включаться")
	assert.Equal(t, 200, cfg.Markdown.MaxMethodDescriptionLen, "Максимальная длина описания метода должна быть 200")
	assert.Equal(t, "github", cfg.Markdown.CodeStyle, "Стиль кода по умолчанию должен быть github")
	assert.True(t, cfg.Markdown.GroupMethodsByType, "Группировка методов по типу должна быть включена")
}

// TestLoadConfig проверяет загрузку конфигурации из файла
func TestLoadConfig(t *testing.T) {
	// Создание YAML содержимого для тестовой конфигурации
	yamlContent := `
filesystem:
  include_patterns:
    - "*.go"
    - "*.js"
  exclude_patterns:
    - "*_test.go"
    - "node_modules/**"
  max_depth: 5

parser:
  parse_private_methods: true
  max_file_size: 524288

llm:
  provider: "anthropic"
  model: "claude-3"
  temperature: 0.5
  max_tokens: 2000
  batch_size: 10
  batch_delay: 2

markdown:
  include_toc: false
  include_file_info: true
  max_method_description_length: 300
  group_methods_by_type: true
  code_style: "monokai"
`

	// Создание временного файла конфигурации
	tmpfile := createTempConfigFile(t, yamlContent)
	defer os.Remove(tmpfile.Name())

	// Загрузка конфигурации
	cfg, err := config.LoadConfig(tmpfile.Name())

	// Проверки
	assert.NoError(t, err, "Загрузка конфигурации должна выполняться без ошибок")
	assert.NotNil(t, cfg, "Конфигурация не должна быть nil")

	// Проверка загруженных значений для FileSystem
	assert.Equal(t, 5, cfg.FileSystem.MaxDepth, "Максимальная глубина должна быть 5")
	assert.Equal(t, []string{"*.go", "*.js"}, cfg.FileSystem.IncludePatterns, "Шаблоны включения должны соответствовать")
	assert.Equal(t, []string{"*_test.go", "node_modules/**"}, cfg.FileSystem.ExcludePatterns, "Шаблоны исключения должны соответствовать")

	// Проверка загруженных значений для Parser
	assert.True(t, cfg.Parser.ParsePrivateMethods, "Парсинг приватных методов должен быть включен")
	assert.Equal(t, int64(524288), cfg.Parser.MaxFileSize, "Максимальный размер файла должен быть 512KB")

	// Проверка загруженных значений для LLM
	assert.Equal(t, "anthropic", cfg.LLM.Provider, "Провайдер ЛЛМ должен быть Anthropic")
	assert.Equal(t, "claude-3", cfg.LLM.Model, "Модель должна быть claude-3")
	assert.Equal(t, 0.5, cfg.LLM.Temperature, "Температура должна быть 0.5")
	assert.Equal(t, 2000, cfg.LLM.MaxTokens, "Максимальное количество токенов должно быть 2000")
	assert.Equal(t, 10, cfg.LLM.BatchSize, "Размер пакета должен быть 10")
	assert.Equal(t, 2, cfg.LLM.BatchDelay, "Задержка между пакетами должна быть 2")

	// Проверка загруженных значений для Markdown
	assert.False(t, cfg.Markdown.IncludeTOC, "Оглавление не должно включаться")
	assert.True(t, cfg.Markdown.IncludeFileInfo, "Информация о файле должна включаться")
	assert.Equal(t, 300, cfg.Markdown.MaxMethodDescriptionLen, "Максимальная длина описания метода должна быть 300")
	assert.Equal(t, "monokai", cfg.Markdown.CodeStyle, "Стиль кода должен быть monokai")
	assert.True(t, cfg.Markdown.GroupMethodsByType, "Группировка методов по типу должна быть включена")
}

// TestLoadInvalidConfig проверяет поведение при загрузке невалидной конфигурации
func TestLoadInvalidConfig(t *testing.T) {
	// Создаем временный файл с невалидным YAML
	yamlContent := `
filesystem:
  max_depth: "invalid" # должно быть числом
`
	tmpfile := createTempConfigFile(t, yamlContent)
	defer os.Remove(tmpfile.Name())

	// Пытаемся загрузить конфигурацию
	_, err := config.LoadConfig(tmpfile.Name())

	// Должна быть ошибка
	assert.Error(t, err, "Должна возникнуть ошибка при загрузке невалидной конфигурации")
}

// TestLoadNonExistentConfig проверяет поведение при загрузке несуществующего файла
func TestLoadNonExistentConfig(t *testing.T) {
	// Пытаемся загрузить конфигурацию из несуществующего файла
	cfg, err := config.LoadConfig("non-existent-file.yaml")

	// Должна быть ошибка, а конфигурация - nil
	assert.Error(t, err, "Должна возникнуть ошибка при загрузке несуществующего файла")
	assert.Nil(t, cfg, "Конфигурация должна быть nil")
}
