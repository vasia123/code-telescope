package orchestrator

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"code-telescope/internal/config"
	"code-telescope/internal/filesystem"
	"code-telescope/internal/llm"
	"code-telescope/internal/logger"
	"code-telescope/internal/markdown"
	"code-telescope/internal/parser"
	"code-telescope/pkg/models"
)

// Orchestrator координирует весь процесс генерации карты кода
type Orchestrator struct {
	config        *config.Config
	verbose       bool
	scanner       *filesystem.Scanner
	parserFactory *parser.LanguageFactory
	llmProvider   llm.LLMProvider
	promptBuilder *llm.PromptBuilder
	mdGenerator   *markdown.Generator
}

// New создает новый экземпляр оркестратора
func New(cfg *config.Config, verbose bool) (*Orchestrator, error) {
	// Настраиваем уровень логирования в зависимости от verbose
	if verbose {
		logger.SetLevel(logger.DebugLevel)
	} else {
		logger.SetLevel(logger.InfoLevel)
	}

	logger.Info("Инициализация оркестратора")
	scanner := filesystem.New(cfg)
	parserFactory := parser.NewLanguageFactory(cfg)

	// Создаем конфигурацию для LLM провайдера
	llmConfig := map[string]interface{}{
		"api_key":         cfg.LLM.APIKey,
		"model":           cfg.LLM.Model,
		"timeout_seconds": 60,
	}

	// Инициализируем провайдера ЛЛМ
	logger.Infof("Инициализация провайдера ЛЛМ: %s", cfg.LLM.Provider)
	provider, err := llm.GetProvider(cfg.LLM.Provider, llmConfig)
	if err != nil {
		err = logger.OrchestratorError("не удалось инициализировать провайдера ЛЛМ", err)
		return nil, logger.LogError(err)
	}

	// Создаем конструктор промптов
	logger.Debug("Инициализация конструктора промптов")
	promptBuilder := llm.NewPromptBuilder(cfg.LLM.MaxTokens)

	// Создаем генератор Markdown
	logger.Debug("Инициализация генератора Markdown")
	mdGenerator := markdown.New(cfg)

	logger.Info("Оркестратор успешно инициализирован")
	return &Orchestrator{
		config:        cfg,
		verbose:       verbose,
		scanner:       scanner,
		parserFactory: parserFactory,
		llmProvider:   provider,
		promptBuilder: promptBuilder,
		mdGenerator:   mdGenerator,
	}, nil
}

// GenerateCodeMap генерирует карту кода для указанного проекта
func (o *Orchestrator) GenerateCodeMap(projectPath string) (string, error) {
	startTime := time.Now()
	logger.WithFields(logger.Fields{
		"project_path": projectPath,
	}).Info("Запуск генерации карты кода")

	// Шаг 1: Сканирование файловой системы
	logger.Infof("Сканирование файловой системы в директории: %s", projectPath)
	files, err := o.scanner.ScanProject(projectPath)
	if err != nil {
		err = logger.OrchestratorError("ошибка при сканировании проекта", err)
		return "", logger.LogError(err)
	}

	logger.Infof("Найдено %d файлов для анализа", len(files))

	// Шаг 2: Парсинг кода и генерация описаний
	ctx := context.Background()

	// Подготовка коллекции файловых структур для генератора Markdown
	fileStructures := make([]models.FileStructure, 0, len(files))

	logger.Info("Парсинг файлов и генерация описаний")
	for _, file := range files {
		logger.WithField("file", file.Path).Debug("Обработка файла")

		// Получаем парсер для текущего файла
		parser, err := o.parserFactory.GetParserForFile(file.Path)
		if err != nil {
			logger.WithFields(logger.Fields{
				"file":  file.Path,
				"error": err.Error(),
			}).Warn("Пропуск файла (нет подходящего парсера)")
			continue
		}

		// Парсим файл
		logger.WithField("file", file.Path).Debug("Парсинг файла")
		codeStructure, err := parser.Parse(file)
		if err != nil {
			logger.WithFields(logger.Fields{
				"file":  file.Path,
				"error": err.Error(),
			}).Warn("Ошибка при парсинге файла")
			continue
		}

		// Получаем публичные методы для обработки через ЛЛМ
		logger.WithField("file", file.Path).Debug("Извлечение публичных методов")
		publicMethods := codeStructure.GetPublicMethods()
		if len(publicMethods) > 0 {
			logger.Debugf("Найдено %d публичных методов в файле %s", len(publicMethods), file.Path)

			// Конвертируем методы в формат для ЛЛМ
			methodInfos := make([]models.MethodInfo, 0, len(publicMethods))
			for _, method := range publicMethods {
				paramStrings := make([]string, 0, len(method.Parameters))
				for _, param := range method.Parameters {
					paramStrings = append(paramStrings, param.Name+": "+param.Type)
				}

				methodInfo := models.MethodInfo{
					Name:      method.Name,
					Signature: method.Name + "(" + strings.Join(paramStrings, ", ") + ") " + method.ReturnType,
					// Тело метода и другие поля могут быть добавлены при необходимости
				}
				methodInfos = append(methodInfos, methodInfo)
			}

			// Получаем описания методов через ЛЛМ (пакетная обработка)
			if len(methodInfos) > 0 {
				logger.WithField("file", file.Path).Debug("Генерация описаний методов через ЛЛМ")

				// Если методов много, обрабатываем их пакетами
				batchSize := o.config.LLM.BatchSize
				if batchSize <= 0 {
					batchSize = 5 // Значение по умолчанию
				}

				logger.Debugf("Обработка методов пакетами по %d", batchSize)
				for i := 0; i < len(methodInfos); i += batchSize {
					end := i + batchSize
					if end > len(methodInfos) {
						end = len(methodInfos)
					}

					batchMethods := methodInfos[i:end]
					prompt := o.promptBuilder.BuildBatchMethodPrompt(batchMethods, "")

					llmRequest := llm.LLMRequest{
						Prompt:      prompt,
						MaxTokens:   o.config.LLM.MaxTokens,
						Temperature: o.config.LLM.Temperature,
					}

					logger.Debugf("Отправка запроса к ЛЛМ для пакета из %d методов", len(batchMethods))
					response, err := o.llmProvider.GenerateText(ctx, llmRequest)
					if err != nil {
						logger.WithError(err).Warn("Ошибка при получении описаний методов от ЛЛМ")
						continue
					}

					logger.Debug("Парсинг ответа от ЛЛМ")
					methodDescriptions := o.promptBuilder.ParseBatchResponse(response.Text, batchMethods)

					// Добавляем описания к методам
					logger.Debug("Применение описаний к методам")
					for j := i; j < end && j-i < len(batchMethods); j++ {
						if j < len(publicMethods) {
							methodInfo := batchMethods[j-i]
							description, ok := methodDescriptions[methodInfo.Name]
							if ok {
								logger.Debugf("Добавлено описание для метода %s", methodInfo.Name)
								publicMethods[j].Description = description
							} else {
								logger.Warnf("Не удалось получить описание для метода %s", methodInfo.Name)
							}
						}
					}
				}
			}
		}

		// Добавляем структуру файла в коллекцию
		logger.WithField("file", file.Path).Debug("Добавление структуры файла в коллекцию")
		fileStructures = append(fileStructures, codeStructure)
	}

	// Шаг 3: Генерация Markdown с использованием нового модуля
	logger.Info("Генерация Markdown-документации")
	projectName := filepath.Base(projectPath)
	codeMapContent := o.mdGenerator.GenerateCodeMap(fileStructures, projectName)

	// Расчет времени выполнения
	elapsedTime := time.Since(startTime)
	logger.WithFields(logger.Fields{
		"elapsed_time": elapsedTime.String(),
		"files_count":  len(fileStructures),
	}).Info("Генерация карты кода завершена успешно")

	return codeMapContent, nil
}

// SaveCodeMap сохраняет сгенерированную карту кода в файл
func (o *Orchestrator) SaveCodeMap(codeMap, outputPath string) error {
	logger.WithField("output_path", outputPath).Info("Сохранение карты кода в файл")

	// Создаем директории, если они не существуют
	dir := filepath.Dir(outputPath)
	logger.Debugf("Проверка/создание директории: %s", dir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		err = logger.FileSystemError("ошибка при создании директории", err)
		return logger.LogError(err)
	}

	logger.Debugf("Запись %d байт в файл", len(codeMap))
	err := os.WriteFile(outputPath, []byte(codeMap), 0644)
	if err != nil {
		err = logger.FileSystemError("ошибка при записи в файл", err)
		return logger.LogError(err)
	}

	logger.Info("Карта кода успешно сохранена")
	return nil
}
