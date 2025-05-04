package parser

import (
	"code-telescope/pkg/models"

	sitter "github.com/smacker/go-tree-sitter"
)

// Parser определяет интерфейс для парсеров кода
type Parser interface {
	// Parse разбирает содержимое файла и возвращает его структуру.
	Parse(fileMetadata *models.FileMetadata) (*models.CodeStructure, error)

	// GetSupportedExtensions возвращает список расширений файлов, поддерживаемых этим парсером.
	GetSupportedExtensions() []string

	// GetLanguageName возвращает имя языка программирования, поддерживаемого этим парсером.
	GetLanguageName() string

	// ParseTreeNode разбирает узел дерева синтаксического анализа Tree-sitter.
	// Этот метод специфичен для Tree-sitter и может потребовать рефакторинга
	// или быть частью внутренней реализации конкретных парсеров.
	// Пока оставляем его здесь для совместимости с существующей структурой,
	// но он должен быть реализован парсерами языков.
	ParseTreeNode(node *sitter.Node, structure *models.CodeStructure, content []byte) error
}
