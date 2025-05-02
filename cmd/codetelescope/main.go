package main

import (
	"flag"
	"fmt"
	"os"

	"code-telescope/internal/config"
	"code-telescope/internal/orchestrator"
)

// Версия программы, устанавливается при сборке через ldflags
var Version = "dev"

func main() {
	// Парсим аргументы командной строки
	configPath := flag.String("config", "", "Путь к файлу конфигурации")
	outputPath := flag.String("output", "code_map.md", "Путь для сохранения карты кода")
	verbose := flag.Bool("verbose", false, "Подробный вывод")
	flag.Parse()

	// Проверяем наличие пути к проекту
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Необходимо указать путь к проекту для анализа")
		fmt.Println("Использование: codetelescope [опции] <путь_к_проекту>")
		flag.PrintDefaults()
		os.Exit(1)
	}
	projectPath := args[0]

	// Загружаем конфигурацию
	cfg, err := loadConfig(*configPath)
	if err != nil {
		fmt.Printf("Ошибка загрузки конфигурации: %s\n", err)
		os.Exit(1)
	}

	// Создаем оркестратор
	orch, err := orchestrator.New(cfg, *verbose)
	if err != nil {
		fmt.Printf("Ошибка инициализации оркестратора: %s\n", err)
		os.Exit(1)
	}

	// Генерируем карту кода
	if *verbose {
		fmt.Println("Начало генерации карты кода...")
	}

	codeMap, err := orch.GenerateCodeMap(projectPath)
	if err != nil {
		fmt.Printf("Ошибка генерации карты кода: %s\n", err)
		os.Exit(1)
	}

	// Сохраняем результат
	if err := orch.SaveCodeMap(codeMap, *outputPath); err != nil {
		fmt.Printf("Ошибка сохранения карты кода: %s\n", err)
		os.Exit(1)
	}

	if *verbose {
		fmt.Printf("Карта кода успешно сохранена в файл: %s\n", *outputPath)
	}
}

func loadConfig(configPath string) (*config.Config, error) {
	// Если путь к конфигурации не указан, используем конфигурацию по умолчанию
	if configPath == "" {
		return config.DefaultConfig(), nil
	}

	// Проверяем существование файла
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("файл конфигурации не найден: %s", configPath)
	}

	// Загружаем конфигурацию из файла
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
