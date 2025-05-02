package models

import (
	"code-telescope/internal/config"
)

// Parser определяет интерфейс для парсеров различных языков программирования
type Parser interface {
	// Parse разбирает файл и извлекает его структуру
	Parse(fileMetadata *FileMetadata) (*FileStructure, error)

	// GetSupportedExtensions возвращает список поддерживаемых расширений файлов
	GetSupportedExtensions() []string

	// GetLanguageName возвращает название языка программирования
	GetLanguageName() string
}

// BaseParser содержит общую функциональность для всех парсеров
type BaseParser struct {
	Config *config.Config
}

// NewBaseParser создает новый экземпляр BaseParser
func NewBaseParser(cfg *config.Config) *BaseParser {
	return &BaseParser{
		Config: cfg,
	}
}

// ParseOptions содержит опции для процесса парсинга
type ParseOptions struct {
	// Включать приватные методы и свойства
	IncludePrivate bool

	// Глубина парсинга AST
	Depth int

	// Дополнительные опции, специфичные для конкретного языка
	LanguageSpecific map[string]interface{}
}

// NewParseOptions создает новые опции парсинга на основе конфигурации
func NewParseOptions(cfg *config.Config) *ParseOptions {
	return &ParseOptions{
		IncludePrivate:   cfg.Parser.ParsePrivateMethods,
		Depth:            5, // Значение по умолчанию
		LanguageSpecific: make(map[string]interface{}),
	}
}
