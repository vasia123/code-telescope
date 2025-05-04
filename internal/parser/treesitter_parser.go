package parser

import (
	"io/ioutil"
	"sync"

	"code-telescope/internal/config"
	"code-telescope/pkg/models"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

// TreeSitterParser предоставляет базовую реализацию парсера на основе Tree-sitter
type TreeSitterParser struct {
	*BaseParser
	parser     *tree_sitter.Parser
	language   *tree_sitter.Language
	initOnce   sync.Once
	extensions []string
	langName   string
}

// NewTreeSitterParser создает новый Tree-sitter парсер с указанной конфигурацией
func NewTreeSitterParser(cfg *config.Config, language *tree_sitter.Language, extensions []string, langName string) *TreeSitterParser {
	return &TreeSitterParser{
		BaseParser: NewBaseParser(cfg),
		language:   language,
		extensions: extensions,
		langName:   langName,
	}
}

// initParser инициализирует Tree-sitter парсер
func (p *TreeSitterParser) initParser() {
	p.initOnce.Do(func() {
		p.parser = tree_sitter.NewParser()
		p.parser.SetLanguage(*p.language)
	})
}

// Parse разбирает файл и извлекает его структуру
func (p *TreeSitterParser) Parse(fileMetadata *models.FileMetadata) (*models.CodeStructure, error) {
	p.initParser()

	content, err := ioutil.ReadFile(fileMetadata.AbsolutePath)
	if err != nil {
		return nil, err
	}

	tree := p.parser.Parse(content, nil)
	defer tree.Close()

	root := tree.RootNode()

	// Создаем базовую структуру кода
	codeStructure := models.NewCodeStructure(fileMetadata)

	// Метод ParseTreeNode должен быть реализован в конкретных парсерах
	if err := p.ParseTreeNode(root, codeStructure, content); err != nil {
		return nil, err
	}

	return codeStructure, nil
}

// ParseTreeNode абстрактный метод для разбора дерева Tree-sitter
// Должен быть переопределен в конкретных парсерах
func (p *TreeSitterParser) ParseTreeNode(node *tree_sitter.Node, structure *models.CodeStructure, content []byte) error {
	// Реализуется в конкретных парсерах
	return nil
}

// GetSupportedExtensions возвращает поддерживаемые расширения файлов
func (p *TreeSitterParser) GetSupportedExtensions() []string {
	return p.extensions
}

// GetLanguageName возвращает название языка программирования
func (p *TreeSitterParser) GetLanguageName() string {
	return p.langName
}
