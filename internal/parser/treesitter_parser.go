package parser

import (
	"os"
	"sync"

	"code-telescope/pkg/models"

	sitter "github.com/smacker/go-tree-sitter"
)

// TreeSitterParser предоставляет базовую реализацию парсера на основе Tree-sitter.
// Он используется конкретными языковыми парсерами.
type TreeSitterParser struct {
	parser   *sitter.Parser
	language *sitter.Language
	initOnce sync.Once
	// Метод ParseTreeNode должен быть предоставлен конкретным парсером языка
	parseTreeNodeFunc func(node *sitter.Node, structure *models.CodeStructure, content []byte) error
}

// NewTreeSitterParser создает новый базовый Tree-sitter парсер.
// language: Tree-sitter язык.
// parseTreeNodeFunc: Функция для разбора узлов, специфичная для языка.
func NewTreeSitterParser(language *sitter.Language, parseTreeNodeFunc func(node *sitter.Node, structure *models.CodeStructure, content []byte) error) *TreeSitterParser {
	return &TreeSitterParser{
		language:          language,
		parseTreeNodeFunc: parseTreeNodeFunc,
	}
}

// initParser инициализирует внутренний Tree-sitter парсер
func (p *TreeSitterParser) initParser() {
	p.initOnce.Do(func() {
		p.parser = sitter.NewParser()
		p.parser.SetLanguage(p.language)
	})
}

// Parse разбирает файл и извлекает его структуру, используя предоставленную parseTreeNodeFunc.
func (p *TreeSitterParser) Parse(fileMetadata *models.FileMetadata) (*models.CodeStructure, error) {
	p.initParser()

	content, err := os.ReadFile(fileMetadata.AbsolutePath)
	if err != nil {
		return nil, err
	}

	tree := p.parser.Parse(nil, content)
	// Важно закрыть дерево после использования, чтобы избежать утечек памяти CGO
	// Однако, если root узел используется после этого вызова (например, в parseTreeNodeFunc),
	// закрытие здесь может вызвать проблемы. Передаем content для работы с исходным кодом.
	// defer tree.Close() // Закрытие будет управляться вызывающей стороной или при необходимости

	if tree == nil {
		// Обработка случая, когда парсинг не удался
		// Можно вернуть ошибку или пустую структуру
		return models.NewCodeStructure(fileMetadata), nil // Возвращаем пустую структуру
	}
	defer tree.Close() // Закрываем дерево здесь, после получения root

	root := tree.RootNode()
	if root == nil {
		return models.NewCodeStructure(fileMetadata), nil // Возвращаем пустую структуру, если нет корневого узла
	}

	// Создаем базовую структуру кода
	codeStructure := models.NewCodeStructure(fileMetadata)

	// Вызываем специфичную для языка функцию разбора узлов
	if p.parseTreeNodeFunc == nil {
		// Возвращаем ошибку или просто базовую структуру, если функция не задана
		return codeStructure, nil // Или вернуть ошибку fmt.Errorf("parseTreeNodeFunc is not set")
	}
	if err := p.parseTreeNodeFunc(root, codeStructure, content); err != nil {
		return nil, err
	}

	return codeStructure, nil
}

// Close освобождает ресурсы парсера. Пока не используется, но может понадобиться.
// func (p *TreeSitterParser) Close() {
// 	if p.parser != nil {
// 		// Согласно документации go-tree-sitter, Parser имеет метод Close(), но он не экспортирован.
// 		// Управление памятью CGO для Parser обрабатывается через runtime finalizer в библиотеке.
// 		// Закрывать нужно только Tree, TreeCursor, Query, QueryCursor, LookaheadIterator.
// 	}
// }
