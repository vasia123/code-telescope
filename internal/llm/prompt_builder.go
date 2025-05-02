package llm

import (
	"fmt"
	"strings"

	"code-telescope/pkg/models"
)

// PromptBuilder предоставляет методы для создания промптов для различных задач
type PromptBuilder struct {
	maxContextLength int
}

// NewPromptBuilder создает новый экземпляр PromptBuilder
func NewPromptBuilder(maxContextLength int) *PromptBuilder {
	if maxContextLength <= 0 {
		maxContextLength = 8000
	}
	return &PromptBuilder{
		maxContextLength: maxContextLength,
	}
}

// BuildMethodDescriptionPrompt создает промпт для генерации описания метода
func (pb *PromptBuilder) BuildMethodDescriptionPrompt(methodInfo models.MethodInfo, fileContext string) string {
	// Базовый шаблон промпта
	templateStr := `Проанализируй следующий код метода и предоставь краткое, точное описание 
его функциональности в одном абзаце (3-4 предложения максимум). 
Фокусируйся на том, что метод делает, его входных и выходных данных, и основных побочных эффектах.

Метод: %s
Сигнатура: %s
Контекст файла: %s

Предоставь только описание метода без дополнительного форматирования, пояснений или вступлений.`

	// Обрезаем контекст файла, если он слишком длинный
	truncatedContext := fileContext
	if len(fileContext) > pb.maxContextLength {
		truncatedContext = fileContext[:pb.maxContextLength] + "...[контекст обрезан из-за длины]"
	}

	return fmt.Sprintf(templateStr,
		methodInfo.Name,
		methodInfo.Signature,
		truncatedContext)
}

// BuildFileSummaryPrompt создает промпт для генерации общего описания файла
func (pb *PromptBuilder) BuildFileSummaryPrompt(fileInfo models.FileStructure) string {
	// Собираем импорты в строку
	imports := strings.Join(fileInfo.Imports, "\n")

	// Собираем экспорты в строку
	var exports strings.Builder
	for _, exp := range fileInfo.Exports {
		exports.WriteString(exp)
		exports.WriteString("\n")
	}

	// Собираем методы в строку
	var methods strings.Builder
	for _, method := range fileInfo.Methods {
		methods.WriteString(fmt.Sprintf("Метод: %s\nСигнатура: %s\n\n",
			method.Name, method.Signature))
	}

	// Базовый шаблон промпта
	templateStr := `Проанализируй структуру файла и предоставь краткое описание его назначения 
и функциональности в одном абзаце (максимум 3-4 предложения).
Фокусируйся на том, что файл реализует, его основном назначении и взаимодействии с другими компонентами.

Информация о файле:
Имя файла: %s
Язык: %s
Импорты:
%s

Экспорты:
%s

Публичные методы:
%s

Предоставь только описание файла без дополнительного форматирования, пояснений или вступлений.`

	return fmt.Sprintf(templateStr,
		fileInfo.Path,
		fileInfo.Language,
		imports,
		exports.String(),
		methods.String())
}

// BuildBatchMethodPrompt создает промпт для пакетной обработки методов
func (pb *PromptBuilder) BuildBatchMethodPrompt(methods []models.MethodInfo, fileContext string) string {
	var methodsStr strings.Builder

	for i, method := range methods {
		methodsStr.WriteString(fmt.Sprintf("Метод %d: %s\nСигнатура: %s\n\n",
			i+1, method.Name, method.Signature))
	}

	// Обрезаем контекст файла, если он слишком длинный
	truncatedContext := fileContext
	if len(fileContext) > pb.maxContextLength {
		truncatedContext = fileContext[:pb.maxContextLength] + "...[контекст обрезан из-за длины]"
	}

	templateStr := `Проанализируй следующие методы из одного файла и предоставь краткое, точное описание 
для каждого метода. Для каждого метода напиши один абзац (3-4 предложения максимум).
Фокусируйся на том, что метод делает, его входных и выходных данных, и основных побочных эффектах.

Методы:
%s

Контекст файла:
%s

Формат вывода:
Метод 1: [Описание метода 1]
Метод 2: [Описание метода 2]
...и так далее

Предоставь только описания методов в указанном формате без дополнительных пояснений или вступлений.`

	return fmt.Sprintf(templateStr, methodsStr.String(), truncatedContext)
}

// ParseBatchResponse разбирает ответ от ЛЛМ, содержащий описания нескольких методов
func (pb *PromptBuilder) ParseBatchResponse(response string, methods []models.MethodInfo) map[string]string {
	result := make(map[string]string)

	lines := strings.Split(response, "\n")
	var currentMethod string
	var currentDescription strings.Builder

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Пропускаем пустые строки
		if trimmed == "" {
			continue
		}

		// Проверяем, является ли строка заголовком метода
		isMethodHeader := false
		for i, method := range methods {
			prefix := fmt.Sprintf("Метод %d:", i+1)
			if strings.HasPrefix(trimmed, prefix) {
				// Если у нас уже есть текущий метод, сохраняем его описание
				if currentMethod != "" {
					result[currentMethod] = strings.TrimSpace(currentDescription.String())
					currentDescription.Reset()
				}

				// Устанавливаем новый текущий метод
				currentMethod = method.Name
				descriptionPart := strings.TrimPrefix(trimmed, prefix)
				currentDescription.WriteString(strings.TrimSpace(descriptionPart))
				isMethodHeader = true
				break
			}
		}

		// Если это не заголовок метода, добавляем строку к текущему описанию
		if !isMethodHeader && currentMethod != "" {
			currentDescription.WriteString(" ")
			currentDescription.WriteString(trimmed)
		}
	}

	// Сохраняем последнее описание, если оно есть
	if currentMethod != "" {
		result[currentMethod] = strings.TrimSpace(currentDescription.String())
	}

	return result
}
