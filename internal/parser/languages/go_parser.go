package languages

import (
	"strings"
	"unsafe"

	sitter "github.com/tree-sitter/go-tree-sitter"

	"code-telescope/internal/config"
	"code-telescope/internal/parser"
	"code-telescope/pkg/models"
)

// #cgo CFLAGS: -I${SRCDIR}/../../../vendor/github.com/tree-sitter/tree-sitter-go/src
// #include <tree_sitter/parser.h>
// extern TSLanguage *tree_sitter_go();
import "C"

// GetGoLanguage возвращает язык Go для tree-sitter
func GetGoLanguage() *sitter.Language {
	ptr := unsafe.Pointer(C.tree_sitter_go())
	return sitter.NewLanguage(ptr)
}

// GoParser реализует парсер для языка Go
type GoParser struct {
	parser.BaseTreeSitterParser
}

// NewGoParser создает новый экземпляр парсера Go
func NewGoParser(cfg *config.Config) *GoParser {
	language := GetGoLanguage()
	extensions := []string{".go"}

	return &GoParser{
		BaseTreeSitterParser: *parser.NewBaseTreeSitterParser(cfg, language, extensions, "Go"),
	}
}

// ParseTreeNode разбирает узлы дерева Go кода
func (p *GoParser) ParseTreeNode(node *sitter.Node, structure *models.CodeStructure, content []byte) error {
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
		case "package_clause":
			// Обработать объявление пакета
		case "import_declaration":
			p.parseImport(current, structure, content)
		case "function_declaration":
			p.parseFunction(current, structure, content)
		case "method_declaration":
			p.parseMethod(current, structure, content)
		case "type_declaration":
			p.parseType(current, structure, content)
		case "const_declaration":
			p.parseConstant(current, structure, content)
		case "var_declaration":
			p.parseVariable(current, structure, content)
		}

		if !cursor.GoToNextSibling() {
			break
		}
	}

	return nil
}

// parseImport извлекает импорты
func (p *GoParser) parseImport(node *sitter.Node, structure *models.CodeStructure, content []byte) {
	// Найти все import_spec внутри import_declaration
	cursor := sitter.NewTreeCursor(node)
	defer cursor.Close()

	if !cursor.GoToFirstChild() {
		return
	}

	for {
		current := cursor.CurrentNode()

		if current.Type() == "import_spec" {
			pathNode := current.ChildByFieldName("path")
			if pathNode != nil && pathNode.Type() == "interpreted_string_literal" {
				path := string(content[pathNode.StartByte():pathNode.EndByte()])
				// Удаляем кавычки
				path = strings.Trim(path, "\"")

				var alias string
				// Проверяем наличие псевдонима
				aliasNode := current.ChildByFieldName("name")
				if aliasNode != nil {
					alias = string(content[aliasNode.StartByte():aliasNode.EndByte()])
				}

				imp := &models.Import{
					Path:  path,
					Alias: alias,
					Position: models.Position{
						StartLine:   int(current.StartPoint().Row) + 1,
						StartColumn: int(current.StartPoint().Column) + 1,
						EndLine:     int(current.EndPoint().Row) + 1,
						EndColumn:   int(current.EndPoint().Column) + 1,
					},
				}

				structure.AddImport(imp)
			}
		}

		if !cursor.GoToNextSibling() {
			break
		}
	}
}

// parseFunction извлекает функции
func (p *GoParser) parseFunction(node *sitter.Node, structure *models.CodeStructure, content []byte) {
	nameNode := node.ChildByFieldName("name")
	if nameNode == nil {
		return
	}

	name := string(content[nameNode.StartByte():nameNode.EndByte()])
	isPublic := isPublicName(name)

	method := &models.Method{
		Name:      name,
		IsPublic:  isPublic,
		BelongsTo: "", // Функция верхнего уровня
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

	// Извлекаем возвращаемые значения
	resultNode := node.ChildByFieldName("result")
	if resultNode != nil {
		method.ReturnType = p.parseResultType(resultNode, content)
	}

	structure.AddMethod(method)
}

// parseMethod извлекает методы
func (p *GoParser) parseMethod(node *sitter.Node, structure *models.CodeStructure, content []byte) {
	nameNode := node.ChildByFieldName("name")
	if nameNode == nil {
		return
	}

	name := string(content[nameNode.StartByte():nameNode.EndByte()])
	isPublic := isPublicName(name)

	// Извлекаем тип, к которому привязан метод
	receiverNode := node.ChildByFieldName("receiver")
	var belongsTo string

	if receiverNode != nil {
		// Извлекаем имя типа
		belongsTo = p.parseReceiverType(receiverNode, content)
	}

	method := &models.Method{
		Name:      name,
		IsPublic:  isPublic,
		BelongsTo: belongsTo,
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

	// Извлекаем возвращаемые значения
	resultNode := node.ChildByFieldName("result")
	if resultNode != nil {
		method.ReturnType = p.parseResultType(resultNode, content)
	}

	structure.AddMethod(method)
}

// Вспомогательные методы

// parseType извлекает типы/структуры
func (p *GoParser) parseType(node *sitter.Node, structure *models.CodeStructure, content []byte) {
	cursor := sitter.NewTreeCursor(node)
	defer cursor.Close()

	if !cursor.GoToFirstChild() {
		return
	}

	for {
		current := cursor.CurrentNode()

		if current.Type() == "type_spec" {
			nameNode := current.ChildByFieldName("name")
			if nameNode == nil {
				continue
			}

			name := string(content[nameNode.StartByte():nameNode.EndByte()])
			isPublic := isPublicName(name)

			// Определяем вид типа (структура, интерфейс и т.д.)
			typeNode := current.ChildByFieldName("type")
			var kind string
			var properties []*models.Property

			if typeNode != nil {
				kind = typeNode.Type()

				// Если это структура, извлекаем поля
				if kind == "struct_type" {
					properties = p.parseStructFields(typeNode, content)
				}
			}

			typ := &models.Type{
				Name:       name,
				Kind:       kind,
				IsPublic:   isPublic,
				Properties: properties,
				Position: models.Position{
					StartLine:   int(current.StartPoint().Row) + 1,
					StartColumn: int(current.StartPoint().Column) + 1,
					EndLine:     int(current.EndPoint().Row) + 1,
					EndColumn:   int(current.EndPoint().Column) + 1,
				},
			}

			structure.AddType(typ)
		}

		if !cursor.GoToNextSibling() {
			break
		}
	}
}

// parseConstant извлекает константы
func (p *GoParser) parseConstant(node *sitter.Node, structure *models.CodeStructure, content []byte) {
	// Реализация извлечения констант (заглушка)
}

// parseVariable извлекает переменные
func (p *GoParser) parseVariable(node *sitter.Node, structure *models.CodeStructure, content []byte) {
	// Реализация извлечения переменных (заглушка)
}

// parseParameters извлекает параметры функции/метода
func (p *GoParser) parseParameters(node *sitter.Node, content []byte) []*models.Parameter {
	var parameters []*models.Parameter

	cursor := sitter.NewTreeCursor(node)
	defer cursor.Close()

	if !cursor.GoToFirstChild() {
		return parameters
	}

	for {
		current := cursor.CurrentNode()

		if current.Type() == "parameter_declaration" {
			nameNode := current.ChildByFieldName("name")
			typeNode := current.ChildByFieldName("type")

			if nameNode != nil && typeNode != nil {
				name := string(content[nameNode.StartByte():nameNode.EndByte()])
				typeName := string(content[typeNode.StartByte():typeNode.EndByte()])

				param := &models.Parameter{
					Name:       name,
					Type:       typeName,
					IsRequired: true, // В Go все параметры обязательны
				}

				parameters = append(parameters, param)
			}
		}

		if !cursor.GoToNextSibling() {
			break
		}
	}

	return parameters
}

// parseResultType извлекает возвращаемые значения
func (p *GoParser) parseResultType(node *sitter.Node, content []byte) string {
	// Извлекаем всю строку с возвращаемыми значениями
	if node != nil {
		return string(content[node.StartByte():node.EndByte()])
	}
	return ""
}

// parseReceiverType извлекает тип получателя метода
func (p *GoParser) parseReceiverType(node *sitter.Node, content []byte) string {
	cursor := sitter.NewTreeCursor(node)
	defer cursor.Close()

	if !cursor.GoToFirstChild() {
		return ""
	}

	for {
		current := cursor.CurrentNode()

		if current.Type() == "parameter_declaration" {
			typeNode := current.ChildByFieldName("type")
			if typeNode != nil {
				typeName := string(content[typeNode.StartByte():typeNode.EndByte()])
				// Удаляем символы указателя, если есть
				typeName = strings.TrimPrefix(typeName, "*")
				return typeName
			}
		}

		if !cursor.GoToNextSibling() {
			break
		}
	}

	return ""
}

// parseStructFields извлекает поля структуры
func (p *GoParser) parseStructFields(node *sitter.Node, content []byte) []*models.Property {
	var properties []*models.Property

	cursor := sitter.NewTreeCursor(node)
	defer cursor.Close()

	if !cursor.GoToFirstChild() {
		return properties
	}

	for {
		current := cursor.CurrentNode()

		if current.Type() == "field_declaration" {
			nameNode := current.ChildByFieldName("name")
			typeNode := current.ChildByFieldName("type")

			if nameNode != nil && typeNode != nil {
				name := string(content[nameNode.StartByte():nameNode.EndByte()])
				typeName := string(content[typeNode.StartByte():typeNode.EndByte()])
				isPublic := isPublicName(name)

				property := &models.Property{
					Name:     name,
					Type:     typeName,
					IsPublic: isPublic,
					Position: models.Position{
						StartLine:   int(current.StartPoint().Row) + 1,
						StartColumn: int(current.StartPoint().Column) + 1,
						EndLine:     int(current.EndPoint().Row) + 1,
						EndColumn:   int(current.EndPoint().Column) + 1,
					},
				}

				properties = append(properties, property)
			}
		}

		if !cursor.GoToNextSibling() {
			break
		}
	}

	return properties
}

// isPublicName проверяет, является ли имя публичным по правилам Go
func isPublicName(name string) bool {
	if len(name) == 0 {
		return false
	}
	// В Go публичными являются идентификаторы, начинающиеся с заглавной буквы
	firstChar := name[0]
	return 'A' <= firstChar && firstChar <= 'Z'
}
