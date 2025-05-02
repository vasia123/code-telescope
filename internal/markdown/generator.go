package markdown

import (
	"fmt"
	"regexp"
	"strings"

	"code-telescope/internal/config"
	"code-telescope/pkg/models"
)

// Generator представляет генератор Markdown документации
type Generator struct {
	config *config.Config
}

// New создает новый экземпляр генератора Markdown
func New(cfg *config.Config) *Generator {
	return &Generator{
		config: cfg,
	}
}

// GenerateCodeMap генерирует полную карту кода на основе структур файлов
func (g *Generator) GenerateCodeMap(fileStructures []models.FileStructure, projectName string) string {
	// Заголовок
	codeMap := fmt.Sprintf(CodeMapHeaderTemplate, projectName)

	// Создаем оглавление
	toc := g.generateTableOfContents(fileStructures)
	codeMap += fmt.Sprintf(CodeMapIntroTemplate, toc)

	// Добавляем разделы для каждого файла
	for _, fileStructure := range fileStructures {
		fileSection := g.GenerateFileSection(fileStructure)
		codeMap += fileSection
	}

	return codeMap
}

// GenerateFileSection генерирует Markdown секцию для одного файла
func (g *Generator) GenerateFileSection(fileStructure models.FileStructure) string {
	content := fmt.Sprintf(FileHeaderTemplate, fileStructure.Path)

	// Добавляем секцию импортов/экспортов
	importsExportsContent := g.generateImportsExportsSection(fileStructure.Imports, fileStructure.Exports)
	content += fmt.Sprintf(ImportsExportsTemplate, importsExportsContent)

	// Если есть методы, добавляем их
	if len(fileStructure.Methods) > 0 {
		content += PublicMethodsHeaderTemplate

		for _, method := range fileStructure.Methods {
			// Формируем параметры метода для отображения
			paramsStr := g.formatParameters(method.Params)

			// Формируем возвращаемые значения
			returnsStr := g.formatReturns(method.Returns)

			// Добавляем в карту кода
			content += fmt.Sprintf(MethodTemplate, method.Name, paramsStr, returnsStr, 
				method.Signature)
		}
	}

	return content + "\n"
}

// generateImportsExportsSection генерирует секцию импортов и экспортов
func (g *Generator) generateImportsExportsSection(imports []string, exports []string) string {
	var content string

	// Добавляем импорты
	if len(imports) > 0 {
		content += "Импорты:\n"
		for _, imp := range imports {
			content += fmt.Sprintf("- %s\n", imp)
		}
	}

	// Добавляем экспорты
	if len(exports) > 0 {
		if len(imports) > 0 {
			content += "\n"
		}
		content += "Экспорты:\n"
		for _, exp := range exports {
			content += fmt.Sprintf("- %s\n", exp)
		}
	}

	return content
}

// formatParameters форматирует параметры метода для отображения в Markdown
func (g *Generator) formatParameters(parameters []string) string {
	if len(parameters) == 0 {
		return ""
	}

	paramsStr := "- **Входные параметры**: \n"
	for _, param := range parameters {
		paramsStr += fmt.Sprintf("  - %s\n", param)
	}

	return paramsStr
}

// formatReturns форматирует возвращаемые значения метода для отображения в Markdown
func (g *Generator) formatReturns(returns []string) string {
	if len(returns) == 0 {
		return ""
	}

	returnsStr := "- **Выходные параметры**: \n"
	for _, ret := range returns {
		returnsStr += fmt.Sprintf("  - %s\n", ret)
	}

	return returnsStr
}

// generateTableOfContents генерирует оглавление для карты кода
func (g *Generator) generateTableOfContents(fileStructures []models.FileStructure) string {
	var toc string

	for _, fileStructure := range fileStructures {
		// Создаем якорь из пути файла
		anchor := g.createAnchor(fileStructure.Path)
		toc += fmt.Sprintf(TableOfContentsItemTemplate, fileStructure.Path, anchor)
	}

	return toc
}

// createAnchor создает якорь для оглавления из строки
func (g *Generator) createAnchor(text string) string {
	// Преобразуем в нижний регистр
	anchor := strings.ToLower(text)
	
	// Заменяем пробелы на дефисы
	anchor = strings.ReplaceAll(anchor, " ", "-")
	
	// Удаляем символы, не являющиеся буквами, цифрами, дефисами или подчеркиваниями
	re := regexp.MustCompile(`[^a-z0-9\-_]`)
	anchor = re.ReplaceAllString(anchor, "")
	
	return anchor
}
