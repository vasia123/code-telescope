package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"code-telescope/internal/config"
	"code-telescope/internal/orchestrator"
)

// Версия программы, устанавливается при сборке через ldflags
var Version = "dev"

func main() {
	// Парсинг аргументов командной строки
	projectPath := flag.String("project-path", ".", "Путь к директории проекта")
	configPath := flag.String("config", "", "Путь к конфигурационному файлу")
	outputPath := flag.String("output", "code-map.md", "Путь к выходному Markdown-файлу")
	llmProvider := flag.String("llm-provider", "", "Провайдер ЛЛМ (openai, anthropic)")
	verbose := flag.Bool("verbose", false, "Подробный вывод")
	showVersion := flag.Bool("version", false, "Показать версию программы")

	flag.Parse()

	// Показать версию и выйти, если указан флаг
	if *showVersion {
		fmt.Printf("Code Telescope версия %s\n", Version)
		os.Exit(0)
	}

	// Загрузка конфигурации
	cfg, err := loadConfig(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка загрузки конфигурации: %v\n", err)
		os.Exit(1)
	}

	// Переопределение параметров из командной строки
	if *llmProvider != "" {
		cfg.LLM.Provider = *llmProvider
	}

	// Подготовка абсолютных путей
	absProjectPath, err := filepath.Abs(*projectPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка получения абсолютного пути проекта: %v\n", err)
		os.Exit(1)
	}

	absOutputPath, err := filepath.Abs(*outputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка получения абсолютного пути вывода: %v\n", err)
		os.Exit(1)
	}

	// Создание экземпляра оркестратора
	orchestr := orchestrator.New(cfg, *verbose)

	// Запуск процесса генерации карты кода
	fmt.Printf("Генерация карты кода для проекта: %s\n", absProjectPath)
	codeMap, err := orchestr.GenerateCodeMap(absProjectPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка генерации карты кода: %v\n", err)
		os.Exit(1)
	}

	// Сохранение результата
	if err := orchestr.SaveCodeMap(codeMap, absOutputPath); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка сохранения карты кода: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Карта кода успешно сгенерирована и сохранена в: %s\n", absOutputPath)
}

// Загрузка конфигурации из файла или использование значений по умолчанию
func loadConfig(configPath string) (*config.Config, error) {
	if configPath == "" {
		// Поиск конфигурации в стандартных местах
		defaultPaths := []string{
			"./configs/default.yaml",
			"./config.yaml",
			"./code-telescope.yaml",
		}

		for _, path := range defaultPaths {
			if _, err := os.Stat(path); err == nil {
				configPath = path
				break
			}
		}
	}

	var cfg *config.Config
	var err error

	if configPath != "" {
		cfg, err = config.LoadConfig(configPath)
		if err != nil {
			return nil, fmt.Errorf("не удалось загрузить конфигурацию из %s: %w", configPath, err)
		}
	} else {
		// Использование конфигурации по умолчанию
		cfg = config.DefaultConfig()
	}

	return cfg, nil
}
