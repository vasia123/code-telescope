package languages

import (
	"unsafe"

	sitter "github.com/tree-sitter/go-tree-sitter"

	"code-telescope/internal/config"
	"code-telescope/internal/parser"
	"code-telescope/pkg/models"
)

// #cgo CFLAGS: -I${SRCDIR}/../../../vendor/github.com/tree-sitter/tree-sitter-python/src
// #include <tree_sitter/parser.h>
// extern TSLanguage *tree_sitter_python();
import "C"

// GetPythonLanguage возвращает язык Python для tree-sitter
func GetPythonLanguage() *sitter.Language {
	ptr := unsafe.Pointer(C.tree_sitter_python())
	return sitter.NewLanguage(ptr)
}

// PythonParser реализует парсер для языка Python
type PythonParser struct {
	parser.BaseTreeSitterParser
}

// NewPythonParser создает новый экземпляр парсера Python
func NewPythonParser(cfg *config.Config) *PythonParser {
	language := GetPythonLanguage()
	extensions := []string{".py"}

	return &PythonParser{
		BaseTreeSitterParser: *parser.NewBaseTreeSitterParser(cfg, language, extensions, "Python"),
	}
}

// ParseTreeNode разбирает узлы дерева Python кода
func (p *PythonParser) ParseTreeNode(node *sitter.Node, structure *models.CodeStructure, content []byte) error {
	// Проходим по всем дочерним узлам корня
	cursor := sitter.NewTreeCursor(node)
	defer cursor.Close()

	if !cursor.GoToFirstChild() {
		return nil
	}

	// Обходим дерево и извлекаем структуру
	for {
		current := cursor.CurrentNode()
		nodeType := current.Type()

		switch nodeType {
		case "import_statement":
			p.parseImport(current, structure, content)
		case "import_from_statement":
			p.parseImportFrom(current, structure, content)
		case "function_definition":
			p.parseFunction(current, structure, content)
		case "class_definition":
			p.parseClass(current, structure, content)
		}

		if !cursor.GoToNextSibling() {
			break
		}
	}

	return nil
}

// parseImport извлекает импорты типа "import x"
func (p *PythonParser) parseImport(node *sitter.Node, structure *models.CodeStructure, content []byte) {
	// Находим все имена модулей в импорте
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child == nil {
			continue
		}

		if child.Type() == "dotted_name" {
			moduleName := string(content[child.StartByte():child.EndByte()])

			imp := &models.Import{
				Path: moduleName,
				Position: models.Position{
					StartLine:   int(node.StartPoint().Row) + 1,
					StartColumn: int(node.StartPoint().Column) + 1,
					EndLine:     int(node.EndPoint().Row) + 1,
					EndColumn:   int(node.EndPoint().Column) + 1,
				},
			}

			structure.AddImport(imp)
		}
	}
}

// parseImportFrom извлекает импорты типа "from x import y"
func (p *PythonParser) parseImportFrom(node *sitter.Node, structure *models.CodeStructure, content []byte) {
	// Извлекаем имя модуля (after "from")
	var moduleName string
	var importNames []string

	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child == nil {
			continue
		}

		if child.Type() == "dotted_name" {
			moduleName = string(content[child.StartByte():child.EndByte()])
		} else if child.Type() == "import_statement" {
			// Извлекаем имена импортируемых сущностей
			for j := 0; j < int(child.ChildCount()); j++ {
				importChild := child.Child(j)
				if importChild == nil {
					continue
				}

				if importChild.Type() == "dotted_name" {
					importNames = append(importNames, string(content[importChild.StartByte():importChild.EndByte()]))
				}
			}
		}
	}

	if moduleName != "" {
		for _, name := range importNames {
			imp := &models.Import{
				Path:  moduleName,
				Alias: name,
				Position: models.Position{
					StartLine:   int(node.StartPoint().Row) + 1,
					StartColumn: int(node.StartPoint().Column) + 1,
					EndLine:     int(node.EndPoint().Row) + 1,
					EndColumn:   int(node.EndPoint().Column) + 1,
				},
			}

			structure.AddImport(imp)
		}
	}
}

// parseFunction извлекает функции
func (p *PythonParser) parseFunction(node *sitter.Node, structure *models.CodeStructure, content []byte) {
	nameNode := node.ChildByFieldName("name")
	if nameNode == nil {
		return
	}

	name := string(content[nameNode.StartByte():nameNode.EndByte()])
	isPublic := !startsWith_(name) // В Python "приватные" методы начинаются с _

	method := &models.Method{
		Name:      name,
		IsPublic:  isPublic,
		BelongsTo: "", // Для методов класса будет заполнено при разборе класса
		Position: models.Position{
			StartLine:   int(node.StartPoint().Row) + 1,
			StartColumn: int(node.StartPoint().Column) + 1,
			EndLine:     int(node.EndPoint().Row) + 1,
			EndColumn:   int(node.EndPoint().Column) + 1,
		},
	}

	// Извлекаем параметры
	paramsNode := node.ChildByFieldName("parameters")
	if paramsNode != nil {
		method.Parameters = p.parseParameters(paramsNode, content)
	}

	structure.AddMethod(method)
}

// parseClass извлекает классы
func (p *PythonParser) parseClass(node *sitter.Node, structure *models.CodeStructure, content []byte) {
	nameNode := node.ChildByFieldName("name")
	if nameNode == nil {
		return
	}

	name := string(content[nameNode.StartByte():nameNode.EndByte()])
	isPublic := !startsWith_(name)

	typ := &models.Type{
		Name:     name,
		Kind:     "class",
		IsPublic: isPublic,
		Position: models.Position{
			StartLine:   int(node.StartPoint().Row) + 1,
			StartColumn: int(node.StartPoint().Column) + 1,
			EndLine:     int(node.EndPoint().Row) + 1,
			EndColumn:   int(node.EndPoint().Column) + 1,
		},
	}

	// Извлекаем методы класса из его тела
	bodyNode := node.ChildByFieldName("body")
	if bodyNode != nil {
		p.parseClassBody(bodyNode, structure, content, name)
	}

	structure.AddType(typ)
}

// parseClassBody извлекает содержимое класса
func (p *PythonParser) parseClassBody(node *sitter.Node, structure *models.CodeStructure, content []byte, className string) {
	cursor := sitter.NewTreeCursor(node)
	defer cursor.Close()

	if !cursor.GoToFirstChild() {
		return
	}

	for {
		current := cursor.CurrentNode()
		nodeType := current.Type()

		if nodeType == "function_definition" {
			p.parseClassMethod(current, structure, content, className)
		} else if nodeType == "expression_statement" {
			// Можно извлекать атрибуты класса, определенные на уровне класса
		}

		if !cursor.GoToNextSibling() {
			break
		}
	}
}

// parseClassMethod извлекает метод класса
func (p *PythonParser) parseClassMethod(node *sitter.Node, structure *models.CodeStructure, content []byte, className string) {
	nameNode := node.ChildByFieldName("name")
	if nameNode == nil {
		return
	}

	name := string(content[nameNode.StartByte():nameNode.EndByte()])
	isPublic := !startsWith_(name)

	method := &models.Method{
		Name:      name,
		IsPublic:  isPublic,
		BelongsTo: className,
		Position: models.Position{
			StartLine:   int(node.StartPoint().Row) + 1,
			StartColumn: int(node.StartPoint().Column) + 1,
			EndLine:     int(node.EndPoint().Row) + 1,
			EndColumn:   int(node.EndPoint().Column) + 1,
		},
	}

	// Извлекаем параметры
	paramsNode := node.ChildByFieldName("parameters")
	if paramsNode != nil {
		params := p.parseParameters(paramsNode, content)

		// В Python первым параметром метода класса является self
		if len(params) > 0 && (params[0].Name == "self" || params[0].Name == "cls") {
			params = params[1:] // Пропускаем self/cls в списке параметров
		}

		method.Parameters = params
	}

	structure.AddMethod(method)
}

// parseParameters извлекает параметры функции/метода
func (p *PythonParser) parseParameters(node *sitter.Node, content []byte) []*models.Parameter {
	var parameters []*models.Parameter

	// Ищем список параметров (parameter_list)
	var paramListNode *sitter.Node
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child != nil && child.Type() == "parameters" {
			paramListNode = child
			break
		}
	}

	if paramListNode == nil {
		return parameters
	}

	// Извлекаем каждый параметр
	for i := 0; i < int(paramListNode.ChildCount()); i++ {
		child := paramListNode.Child(i)
		if child == nil {
			continue
		}

		if child.Type() == "identifier" {
			name := string(content[child.StartByte():child.EndByte()])

			param := &models.Parameter{
				Name:       name,
				Type:       "", // Python часто не указывает типы
				IsRequired: true,
			}

			parameters = append(parameters, param)
		}
	}

	return parameters
}

// startsWith_ проверяет начинается ли строка с подчеркивания
func startsWith_(s string) bool {
	return len(s) > 0 && s[0] == '_'
}
