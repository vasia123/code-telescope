package parser

import (
	"context"
	"fmt"
	"os"

	"code-telescope/internal/config"
	"code-telescope/pkg/models"

	sitter "github.com/smacker/go-tree-sitter"
)

// BaseTreeSitterParser предоставляет базовую реализацию для парсеров, использующих Tree-sitter.
type BaseTreeSitterParser struct {
	Cfg        *config.Config
	Language   *sitter.Language
	Extensions []string
	Name       string
}

// NewBaseTreeSitterParser создает новый экземпляр BaseTreeSitterParser.
func NewBaseTreeSitterParser(cfg *config.Config, language *sitter.Language, extensions []string, name string) *BaseTreeSitterParser {
	return &BaseTreeSitterParser{
		Cfg:        cfg,
		Language:   language,
		Extensions: extensions,
		Name:       name,
	}
}

// GetSupportedExtensions возвращает список поддерживаемых расширений файлов.
func (p *BaseTreeSitterParser) GetSupportedExtensions() []string {
	return p.Extensions
}

// GetLanguageName возвращает имя языка.
func (p *BaseTreeSitterParser) GetLanguageName() string {
	return p.Name
}

// Parse разбирает файл с использованием Tree-sitter.
// Это общая реализация, которая читает файл и вызывает ParseTreeNode.
func (p *BaseTreeSitterParser) Parse(fileMetadata *models.FileMetadata) (*models.CodeStructure, error) {
	content, err := os.ReadFile(fileMetadata.Path)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла %s: %w", fileMetadata.Path, err)
	}

	parser := sitter.NewParser()
	parser.SetLanguage(p.Language)

	tree, err := parser.ParseCtx(context.Background(), nil, content)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга файла %s: %w", fileMetadata.Path, err)
	}
	defer tree.Close()

	structure := models.NewCodeStructure(fileMetadata)

	// Получаем ссылку на метод ParseTreeNode конкретного парсера
	// Так как Go не поддерживает прямые вызовы методов переопределенных в дочерних структурах
	// из базовой структуры, нам нужен способ получить конкретный парсер.
	// Один из способов - передать его в Parse или использовать type assertion/switch,
	// но это усложняет интерфейс. Пока оставляем заглушку.
	// TODO: Реализовать вызов специфичного ParseTreeNode для конкретного языка.
	// if specificParser, ok := p.(Parser); ok { // Это не сработает напрямую
	//  err = specificParser.ParseTreeNode(tree.RootNode(), structure, content)
	// }

	// Временное решение: предполагаем, что p - это уже нужный тип парсера
	// и вызываем ParseTreeNode через интерфейс Parser, что потребует
	// реализации этого метода в BaseTreeSitterParser (что нелогично)
	// или использования рефлексии/другого подхода.

	// Пока вызываем ParseTreeNode напрямую, предполагая, что будет найдена
	// реализация в конкретном парсере языка (что не так для Go)
	// Необходимо будет переопределить Parse в каждом парсере языка
	// или изменить архитектуру.

	// Пока оставим так, но это нужно будет исправить
	if parserWithTreeNode, ok := interface{}(p).(interface {
		ParseTreeNode(node *sitter.Node, structure *models.CodeStructure, content []byte) error
	}); ok {
		err = parserWithTreeNode.ParseTreeNode(tree.RootNode(), structure, content)
		if err != nil {
			return nil, fmt.Errorf("ошибка разбора узлов файла %s: %w", fileMetadata.Path, err)
		}
	} else {
		// Этого не должно происходить, если все парсеры реализуют ParseTreeNode
		return nil, fmt.Errorf("парсер для языка %s не реализует ParseTreeNode", p.Name)
	}

	return structure, nil
}

// ParseTreeNode является заглушкой. Конкретные парсеры должны его переопределить.
func (p *BaseTreeSitterParser) ParseTreeNode(node *sitter.Node, structure *models.CodeStructure, content []byte) error {
	// Этот метод должен быть реализован в каждом конкретном парсере языка
	// (GoParser, PythonParser и т.д.), так как логика обхода дерева
	// зависит от грамматики языка.
	return fmt.Errorf("ParseTreeNode не реализован для базового парсера %s", p.Name)
}
