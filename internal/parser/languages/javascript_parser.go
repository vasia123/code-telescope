package languages

import (
	"unsafe"

	sitter "github.com/tree-sitter/go-tree-sitter"

	"code-telescope/internal/config"
	"code-telescope/internal/parser"
	"code-telescope/pkg/models"
)

// #cgo CFLAGS: -I${SRCDIR}/../../../vendor/github.com/tree-sitter/tree-sitter-javascript/src
// #include <tree_sitter/parser.h>
// extern TSLanguage *tree_sitter_javascript();
import "C"

// GetJavaScriptLanguage возвращает язык JavaScript для tree-sitter
func GetJavaScriptLanguage() *sitter.Language {
	ptr := unsafe.Pointer(C.tree_sitter_javascript())
	return sitter.NewLanguage(ptr)
}

// JavaScriptParser реализует парсер для языка JavaScript
type JavaScriptParser struct {
	parser.BaseTreeSitterParser
}

// NewJavaScriptParser создает новый экземпляр парсера JavaScript
func NewJavaScriptParser(cfg *config.Config) *JavaScriptParser {
	language := GetJavaScriptLanguage()
	extensions := []string{".js", ".jsx", ".ts", ".tsx"}

	return &JavaScriptParser{
		BaseTreeSitterParser: *parser.NewBaseTreeSitterParser(cfg, language, extensions, "JavaScript"),
	}
}

// ParseTreeNode разбирает узлы дерева JavaScript кода
func (p *JavaScriptParser) ParseTreeNode(node *sitter.Node, structure *models.CodeStructure, content []byte) error {
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
		case "export_statement":
			p.parseExport(current, structure, content)
		case "function_declaration":
			p.parseFunction(current, structure, content)
		case "method_definition":
			p.parseMethod(current, structure, content)
		case "class_declaration":
			p.parseClass(current, structure, content)
		case "variable_declaration":
			p.parseVariableDeclaration(current, structure, content)
		}

		if !cursor.GoToNextSibling() {
			break
		}
	}

	return nil
}

// parseImport извлекает импорты
func (p *JavaScriptParser) parseImport(node *sitter.Node, structure *models.CodeStructure, content []byte) {
	// Находим источник импорта (путь)
	var sourceNode *sitter.Node
	var specifiersText string

	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child == nil {
			continue
		}

		// Ищем строку источника
		if child.Type() == "string" {
			sourceNode = child
		}

		// Собираем информацию о спецификаторах импорта
		if child.Type() == "import_clause" {
			specifiersText = string(content[child.StartByte():child.EndByte()])
		}
	}

	if sourceNode != nil {
		path := string(content[sourceNode.StartByte():sourceNode.EndByte()])
		// Удаляем кавычки
		path = path[1 : len(path)-1]

		imp := &models.Import{
			Path:  path,
			Alias: "", // JavaScript не использует алиасы как Go
			Info:  specifiersText,
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

// parseExport извлекает экспорты
func (p *JavaScriptParser) parseExport(node *sitter.Node, structure *models.CodeStructure, content []byte) {
	// Определяем тип экспорта
	var name string
	var exportedValue string

	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child == nil {
			continue
		}

		// Если экспортируется функция
		if child.Type() == "function_declaration" {
			nameNode := child.ChildByFieldName("name")
			if nameNode != nil {
				name = string(content[nameNode.StartByte():nameNode.EndByte()])
				exportedValue = "function"
			}
		}

		// Если экспортируется класс
		if child.Type() == "class_declaration" {
			nameNode := child.ChildByFieldName("name")
			if nameNode != nil {
				name = string(content[nameNode.StartByte():nameNode.EndByte()])
				exportedValue = "class"
			}
		}

		// Если используется default export
		if child.Type() == "default" {
			name = "default"
		}
	}

	if name != "" {
		exp := &models.Export{
			Name:  name,
			Value: exportedValue,
			Position: models.Position{
				StartLine:   int(node.StartPoint().Row) + 1,
				StartColumn: int(node.StartPoint().Column) + 1,
				EndLine:     int(node.EndPoint().Row) + 1,
				EndColumn:   int(node.EndPoint().Column) + 1,
			},
		}

		structure.AddExport(exp)
	}
}

// parseFunction извлекает функции
func (p *JavaScriptParser) parseFunction(node *sitter.Node, structure *models.CodeStructure, content []byte) {
	nameNode := node.ChildByFieldName("name")
	if nameNode == nil {
		return
	}

	name := string(content[nameNode.StartByte():nameNode.EndByte()])

	method := &models.Method{
		Name:      name,
		IsPublic:  true, // В JavaScript все функции публичные по умолчанию
		BelongsTo: "",   // Функция верхнего уровня
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

// parseMethod извлекает методы
func (p *JavaScriptParser) parseMethod(node *sitter.Node, structure *models.CodeStructure, content []byte) {
	nameNode := node.ChildByFieldName("name")
	if nameNode == nil {
		return
	}

	name := string(content[nameNode.StartByte():nameNode.EndByte()])

	method := &models.Method{
		Name:      name,
		IsPublic:  true, // В JavaScript методы классов публичные по умолчанию
		BelongsTo: "",   // Заполнится при разборе класса
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
func (p *JavaScriptParser) parseClass(node *sitter.Node, structure *models.CodeStructure, content []byte) {
	nameNode := node.ChildByFieldName("name")
	if nameNode == nil {
		return
	}

	name := string(content[nameNode.StartByte():nameNode.EndByte()])

	typ := &models.Type{
		Name:     name,
		Kind:     "class",
		IsPublic: true, // В JavaScript классы публичные по умолчанию
		Position: models.Position{
			StartLine:   int(node.StartPoint().Row) + 1,
			StartColumn: int(node.StartPoint().Column) + 1,
			EndLine:     int(node.EndPoint().Row) + 1,
			EndColumn:   int(node.EndPoint().Column) + 1,
		},
	}

	// Извлекаем методы и свойства класса
	bodyNode := node.ChildByFieldName("body")
	if bodyNode != nil {
		p.parseClassBody(bodyNode, structure, content, name)
	}

	structure.AddType(typ)
}

// parseClassBody извлекает содержимое класса
func (p *JavaScriptParser) parseClassBody(node *sitter.Node, structure *models.CodeStructure, content []byte, className string) {
	cursor := sitter.NewTreeCursor(node)
	defer cursor.Close()

	if !cursor.GoToFirstChild() {
		return
	}

	for {
		current := cursor.CurrentNode()
		nodeType := current.Type()

		if nodeType == "method_definition" {
			p.parseClassMethod(current, structure, content, className)
		} else if nodeType == "class_property" {
			p.parseClassProperty(current, structure, content, className)
		}

		if !cursor.GoToNextSibling() {
			break
		}
	}
}

// parseClassMethod извлекает метод класса
func (p *JavaScriptParser) parseClassMethod(node *sitter.Node, structure *models.CodeStructure, content []byte, className string) {
	nameNode := node.ChildByFieldName("name")
	if nameNode == nil {
		return
	}

	name := string(content[nameNode.StartByte():nameNode.EndByte()])

	method := &models.Method{
		Name:      name,
		IsPublic:  true, // В JavaScript методы публичные по умолчанию
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
		method.Parameters = p.parseParameters(paramsNode, content)
	}

	structure.AddMethod(method)
}

// parseClassProperty извлекает свойство класса
func (p *JavaScriptParser) parseClassProperty(node *sitter.Node, structure *models.CodeStructure, content []byte, className string) {
	// Извлечение свойств класса
}

// parseVariableDeclaration извлекает объявления переменных
func (p *JavaScriptParser) parseVariableDeclaration(node *sitter.Node, structure *models.CodeStructure, content []byte) {
	// Извлечение объявлений переменных
}

// parseParameters извлекает параметры функции/метода
func (p *JavaScriptParser) parseParameters(node *sitter.Node, content []byte) []*models.Parameter {
	var parameters []*models.Parameter

	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child == nil {
			continue
		}

		if child.Type() == "identifier" {
			name := string(content[child.StartByte():child.EndByte()])

			param := &models.Parameter{
				Name:       name,
				Type:       "", // JavaScript не имеет явных типов
				IsRequired: true,
			}

			parameters = append(parameters, param)
		}
	}

	return parameters
}
