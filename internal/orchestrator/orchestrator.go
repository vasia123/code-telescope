package orchestrator

import (
	"fmt"
	"os"

	"code-telescope/internal/config"
)

// Orchestrator координирует весь процесс генерации карты кода
type Orchestrator struct {
	config  *config.Config
	verbose bool
}

// New создает новый экземпляр оркестратора
func New(cfg *config.Config, verbose bool) *Orchestrator {
	return &Orchestrator{
		config:  cfg,
		verbose: verbose,
	}
}

// GenerateCodeMap генерирует карту кода для указанного проекта
func (o *Orchestrator) GenerateCodeMap(projectPath string) (string, error) {
	// Заглушка для интерфейса
	if o.verbose {
		fmt.Printf("Запуск генерации карты кода для проекта: %s\n", projectPath)
	}

	// TODO: Реализовать полный процесс генерации карты кода
	// 1. Сканирование файловой системы
	// 2. Парсинг кода
	// 3. Взаимодействие с ЛЛМ
	// 4. Генерация Markdown

	// Временная заглушка
	return "# Заглушка для карты кода\n\nЭта карта кода ещё не реализована.", nil
}

// SaveCodeMap сохраняет сгенерированную карту кода в файл
func (o *Orchestrator) SaveCodeMap(codeMap, outputPath string) error {
	if o.verbose {
		fmt.Printf("Сохранение карты кода в файл: %s\n", outputPath)
	}

	return os.WriteFile(outputPath, []byte(codeMap), 0644)
}
