package parser

import (
	"io/ioutil"
	"sync"

	"code-telescope/internal/config"
	"code-telescope/pkg/models"

	sitter "github.com/tree-sitter/go-tree-sitter"
)

// BaseTreeSitterParser предоставляет базовую реализацию парсера на основе Tree-sitter
type BaseTreeSitterParser struct {
	parser   *sitter.Parser
	Language *sitter.Language
	initOnce sync.Once
	Name     string
	Exts     []string
	Config   *config.Config
}

// NewBaseTreeSitterParser создает новый базовый Tree-sitter парсер
func NewBaseTreeSitterParser(cfg *config.Config, language *sitter.Language, extensions []string, langName string) *BaseTreeSitterParser {
	return &BaseTreeSitterParser{
		Language: language,
		Exts:     extensions,
		Name:     langName,
		Config:   cfg,
	}
}

// initParser инициализирует Tree-sitter парсер
func (p *BaseTreeSitterParser) initParser() {
	p.initOnce.Do(func() {
		p.parser = sitter.NewParser()
		p.parser.SetLanguage(p.Language)
	})
}

// GetLanguageName возвращает название языка
func (p *BaseTreeSitterParser) GetLanguageName() string {
	return p.Name
}

// GetSupportedExtensions возвращает поддерживаемые расширения
func (p *BaseTreeSitterParser) GetSupportedExtensions() []string {
	return p.Exts
}

// Parse разбирает файл и извлекает его структуру
func (p *BaseTreeSitterParser) Parse(fileMetadata *models.FileMetadata) (*models.CodeStructure, error) {
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
func (p *BaseTreeSitterParser) ParseTreeNode(node *sitter.Node, structure *models.CodeStructure, content []byte) error {
	// Этот метод должен быть переопределен в конкретных реализациях
	return nil
}
