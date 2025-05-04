package parser

import (
	"fmt"
	"path/filepath"

	"code-telescope/internal/config"
	"code-telescope/internal/parser/languages"
	"code-telescope/pkg/models"
)

// LanguageFactory создает и управляет парсерами для различных языков программирования
type LanguageFactory struct {
	config      *config.Config
	parsers     map[string]models.Parser
	extToParser map[string]string
}

// NewLanguageFactory создает новый экземпляр фабрики парсеров
func NewLanguageFactory(cfg *config.Config) *LanguageFactory {
	factory := &LanguageFactory{
		config:      cfg,
		parsers:     make(map[string]models.Parser),
		extToParser: make(map[string]string),
	}
	factory.registerParsers()
	return factory
}

// registerParsers регистрирует все поддерживаемые парсеры
func (lf *LanguageFactory) registerParsers() {
	// Регистрируем парсер Go
	lf.registerParser(languages.NewGoParser(lf.config))

	// Регистрируем парсер JavaScript
	lf.registerParser(languages.NewJavaScriptParser(lf.config))

	// Регистрируем парсер Python
	lf.registerParser(languages.NewPythonParser(lf.config))

	// Здесь можно добавить регистрацию других парсеров
}

// registerParser регистрирует парсер и связывает его с поддерживаемыми расширениями файлов
func (lf *LanguageFactory) registerParser(parser models.Parser) {
	languageName := parser.GetLanguageName()
	lf.parsers[languageName] = parser

	// Связываем расширения файлов с этим парсером
	for _, ext := range parser.GetSupportedExtensions() {
		lf.extToParser[ext] = languageName
	}
}

// GetParserForFile возвращает подходящий парсер для указанного файла на основе его расширения
func (lf *LanguageFactory) GetParserForFile(filePath string) (models.Parser, error) {
	ext := filepath.Ext(filePath)
	if languageName, ok := lf.extToParser[ext]; ok {
		if parser, ok := lf.parsers[languageName]; ok {
			return parser, nil
		}
	}
	return nil, fmt.Errorf("no parser available for file extension: %s", ext)
}

// GetSupportedLanguages возвращает список поддерживаемых языков программирования
func (lf *LanguageFactory) GetSupportedLanguages() []string {
	languages := make([]string, 0, len(lf.parsers))
	for lang := range lf.parsers {
		languages = append(languages, lang)
	}
	return languages
}

// GetSupportedExtensions возвращает список поддерживаемых расширений файлов
func (lf *LanguageFactory) GetSupportedExtensions() []string {
	extensions := make([]string, 0, len(lf.extToParser))
	for ext := range lf.extToParser {
		extensions = append(extensions, ext)
	}
	return extensions
}
