package parser

import (
	"fmt"
	"path/filepath"

	"code-telescope/internal/config"
	"code-telescope/internal/parser/languages"
	"code-telescope/pkg/models"
)

// LanguageFactory создает парсеры для различных языков программирования
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
	// Регистрируем парсер для Go
	lf.registerParser(languages.NewGoParser(lf.config))

	// Регистрируем парсер для JavaScript
	lf.registerParser(languages.NewJavaScriptParser(lf.config))

	// Регистрируем парсер для Python
	lf.registerParser(languages.NewPythonParser(lf.config))

	// TODO: Добавить парсеры для других языков в будущем
}

// registerParser регистрирует парсер в фабрике
func (lf *LanguageFactory) registerParser(parser models.Parser) {
	languageName := parser.GetLanguageName()
	lf.parsers[languageName] = parser

	// Связываем расширения файлов с этим парсером
	for _, ext := range parser.GetSupportedExtensions() {
		lf.extToParser[ext] = languageName
	}
}

// GetParserForFile возвращает подходящий парсер для указанного файла
func (lf *LanguageFactory) GetParserForFile(filePath string) (models.Parser, error) {
	ext := filepath.Ext(filePath)
	if ext == "" {
		return nil, fmt.Errorf("файл не имеет расширения: %s", filePath)
	}

	// Находим соответствующий парсер по расширению
	languageName, ok := lf.extToParser[ext]
	if !ok {
		return nil, fmt.Errorf("неподдерживаемое расширение файла: %s", ext)
	}

	parser, ok := lf.parsers[languageName]
	if !ok {
		return nil, fmt.Errorf("парсер для языка %s не найден", languageName)
	}

	return parser, nil
}

// GetSupportedLanguages возвращает список поддерживаемых языков
func (lf *LanguageFactory) GetSupportedLanguages() []string {
	languages := make([]string, 0, len(lf.parsers))

	for language := range lf.parsers {
		languages = append(languages, language)
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
