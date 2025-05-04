package parser

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"code-telescope/internal/config"
	// Убираем прямой импорт languages "code-telescope/internal/parser/languages"
	// Импортируем models "code-telescope/pkg/models"
)

// parserConstructor тип функции для создания экземпляра парсера
type parserConstructor func(cfg *config.Config) Parser

var (
	// mu защищает доступ к глобальным мапам регистрации
	mu sync.RWMutex
	// parserConstructors хранит зарегистрированные конструкторы парсеров по имени языка
	parserConstructors = make(map[string]parserConstructor)
	// langToExtensions хранит расширения для каждого языка
	langToExtensions = make(map[string][]string)
	// extToLangName хранит имя языка для каждого расширения
	extToLangName = make(map[string]string)
)

// RegisterParser регистрирует конструктор парсера для указанного языка и расширений.
// Эта функция должна вызываться из init() пакетов парсеров языков.
func RegisterParser(languageName string, extensions []string, constructor parserConstructor) {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := parserConstructors[languageName]; exists {
		// Возможно, стоит логировать или возвращать ошибку, если язык уже зарегистрирован
		fmt.Printf("Warning: Parser for language '%s' already registered. Overwriting.\n", languageName)
	}
	parserConstructors[languageName] = constructor
	langToExtensions[languageName] = extensions

	for _, ext := range extensions {
		if existingLang, ok := extToLangName[ext]; ok {
			fmt.Printf("Warning: Extension '%s' already registered for language '%s'. Re-registering for '%s'.\n", ext, existingLang, languageName)
		}
		extToLangName[ext] = languageName
	}
}

// LanguageFactory создает и управляет парсерами для различных языков программирования
type LanguageFactory struct {
	config *config.Config
	// Кэш созданных парсеров (опционально, для производительности)
	// parsersCache map[string]Parser
	// muCache sync.Mutex
}

// NewLanguageFactory создает новый экземпляр фабрики парсеров
func NewLanguageFactory(cfg *config.Config) *LanguageFactory {
	factory := &LanguageFactory{
		config: cfg,
		// parsersCache: make(map[string]Parser),
	}
	// Регистрация происходит в init() функций пакетов языков
	// factory.registerParsers() // Убрано
	return factory
}

// GetParserForFile возвращает подходящий парсер для указанного файла на основе его расширения.
// Парсер создается при первом запросе для данного языка.
func (lf *LanguageFactory) GetParserForFile(filePath string) (Parser, error) {
	ext := strings.ToLower(filepath.Ext(filePath))

	mu.RLock()
	languageName, ok := extToLangName[ext]
	mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("no parser registered for file extension: %s", ext)
	}

	// TODO: Добавить кэширование созданных парсеров при необходимости
	// lf.muCache.Lock()
	// parser, cached := lf.parsersCache[languageName]
	// lf.muCache.Unlock()
	// if cached {
	// 	 return parser, nil
	// }

	mu.RLock()
	constructor, ok := parserConstructors[languageName]
	mu.RUnlock()

	if !ok {
		// Эта ситуация не должна возникать, если extToLangName и parserConstructors синхронизированы
		return nil, fmt.Errorf("internal error: constructor not found for registered language: %s", languageName)
	}

	parser := constructor(lf.config)

	// TODO: Добавить созданный парсер в кэш
	// lf.muCache.Lock()
	// lf.parsersCache[languageName] = parser
	// lf.muCache.Unlock()

	return parser, nil
}

// GetSupportedLanguages возвращает список поддерживаемых (зарегистрированных) языков программирования
func (lf *LanguageFactory) GetSupportedLanguages() []string {
	mu.RLock()
	defer mu.RUnlock()
	languages := make([]string, 0, len(parserConstructors))
	for lang := range parserConstructors {
		languages = append(languages, lang)
	}
	return languages
}

// GetSupportedExtensions возвращает список поддерживаемых (зарегистрированных) расширений файлов
func (lf *LanguageFactory) GetSupportedExtensions() []string {
	mu.RLock()
	defer mu.RUnlock()
	extensions := make([]string, 0, len(extToLangName))
	for ext := range extToLangName {
		extensions = append(extensions, ext)
	}
	return extensions
}

// --- Устаревшие методы ---
/*
// registerParsers регистрирует все поддерживаемые парсеры
func (lf *LanguageFactory) registerParsers() {
	// Регистрируем парсер Go
	lf.registerParser(languages.NewGoTreeSitterParser(lf.config))

	// Регистрируем парсер JavaScript
	lf.registerParser(languages.NewJavaScriptTreeSitterParser(lf.config))

	// Регистрируем парсер Python
	lf.registerParser(languages.NewPythonTreeSitterParser(lf.config))

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
*/
