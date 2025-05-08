package languages

import (
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	golang_ts "github.com/smacker/go-tree-sitter/golang"

	"code-telescope/internal/config"
	"code-telescope/internal/parser"
	"code-telescope/pkg/models"
)

// GetGoLanguage возвращает язык Go для tree-sitter
func GetGoLanguage() *sitter.Language {
	return golang_ts.GetLanguage()
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
	isAsync := false     // Go не имеет асинхронных функций
	isGenerator := false // Go не имеет генераторов
	isArrow := false     // Go не имеет стрелочных функций
	isIIFE := false      // Go не имеет IIFE

	// Извлекаем параметры
	paramsNode := node.ChildByFieldName("parameters")
	params := p.parseParameters(paramsNode, content)

	// Извлекаем возвращаемые значения
	resultNode := node.ChildByFieldName("result")
	returnType := ""
	if resultNode != nil {
		returnType = p.parseResultType(resultNode, content)
	}

	startLine := int(node.StartPoint().Row)
	startCol := int(node.StartPoint().Column)
	endLine := int(node.EndPoint().Row)
	endCol := int(node.EndPoint().Column)

	fn := &models.Function{
		Name:        name,
		IsPublic:    isPublic,
		IsAsync:     isAsync,
		IsGenerator: isGenerator,
		IsArrow:     isArrow,
		IsIIFE:      isIIFE,
		Parameters:  params,
		ReturnType:  returnType,
		Position: models.Position{
			StartLine:   startLine + 1,
			StartColumn: startCol + 1,
			EndLine:     endLine + 1,
			EndColumn:   endCol + 1,
		},
	}
	structure.AddFunction(fn)
}

// parseMethod извлекает методы
func (p *GoParser) parseMethod(node *sitter.Node, structure *models.CodeStructure, content []byte) {
	nameNode := node.ChildByFieldName("name")
	if nameNode == nil {
		return
	}

	name := string(content[nameNode.StartByte():nameNode.EndByte()])
	isPublic := isPublicName(name)
	isStatic := false      // Go не имеет статических методов
	isAsync := false       // Go не имеет асинхронных методов
	isGenerator := false   // Go не имеет генераторов
	isDecorator := false   // Go не имеет декораторов
	isConstructor := false // Go не имеет конструкторов в том же смысле, что и другие языки
	kind := "method"

	// Извлекаем тип, к которому привязан метод
	receiverNode := node.ChildByFieldName("receiver")
	var belongsTo string

	if receiverNode != nil {
		// Извлекаем имя типа
		belongsTo = p.parseReceiverType(receiverNode, content)
	}

	// Извлекаем параметры
	paramsNode := node.ChildByFieldName("parameters")
	params := p.parseParameters(paramsNode, content)

	// Извлекаем возвращаемые значения
	resultNode := node.ChildByFieldName("result")
	returnType := ""
	if resultNode != nil {
		returnType = p.parseResultType(resultNode, content)
	}

	startLine := int(node.StartPoint().Row)
	startCol := int(node.StartPoint().Column)
	endLine := int(node.EndPoint().Row)
	endCol := int(node.EndPoint().Column)

	method := &models.Method{
		Name:          name,
		IsPublic:      isPublic,
		IsStatic:      isStatic,
		IsAsync:       isAsync,
		IsGenerator:   isGenerator,
		IsDecorator:   isDecorator,
		IsConstructor: isConstructor,
		Kind:          kind,
		BelongsTo:     belongsTo,
		Parameters:    params,
		ReturnType:    returnType,
		Position: models.Position{
			StartLine:   startLine + 1,
			StartColumn: startCol + 1,
			EndLine:     endLine + 1,
			EndColumn:   endCol + 1,
		},
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
			isAbstract := false // Go не имеет абстрактных классов
			isInterface := false
			isMixin := false   // Go не имеет миксинов
			isGeneric := false // TODO: Добавить поддержку дженериков
			isEnum := false    // Go не имеет перечислений в том же смысле, что и другие языки

			// Определяем вид типа (структура, интерфейс и т.д.)
			typeNode := current.ChildByFieldName("type")
			var kind string
			var properties []*models.Property
			var methods []*models.Method
			var parent string
			var implements []string
			var genericParameters []string

			if typeNode != nil {
				kind = typeNode.Type()

				// Если это структура, извлекаем поля
				if kind == "struct_type" {
					properties = p.parseStructFields(typeNode, content)
				} else if kind == "interface_type" {
					isInterface = true
					// TODO: Извлечь методы интерфейса
				}
			}

			typ := &models.Type{
				Name:              name,
				IsPublic:          isPublic,
				IsAbstract:        isAbstract,
				IsInterface:       isInterface,
				IsMixin:           isMixin,
				IsGeneric:         isGeneric,
				IsEnum:            isEnum,
				Kind:              kind,
				Parent:            parent,
				Implements:        implements,
				GenericParameters: genericParameters,
				Properties:        properties,
				Methods:           methods,
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
				isRequired := true // В Go все параметры обязательны
				isVariadic := false
				isDestructuredObject := false
				isDestructuredArray := false

				// Проверяем на вариативные параметры (...)
				if strings.HasPrefix(typeName, "...") {
					isVariadic = true
					typeName = strings.TrimPrefix(typeName, "...")
				}

				param := &models.Parameter{
					Name:                 name,
					Type:                 typeName,
					IsRequired:           isRequired,
					DefaultValue:         "", // Go не поддерживает значения по умолчанию
					IsVariadic:           isVariadic,
					IsDestructuredObject: isDestructuredObject,
					IsDestructuredArray:  isDestructuredArray,
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
				isStatic := false   // Go не имеет статических полей
				isComputed := false // Go не имеет вычисляемых свойств
				isPrivate := !isPublic
				isReadonly := false // TODO: Определить readonly поля

				property := &models.Property{
					Name:       name,
					Type:       typeName,
					IsPublic:   isPublic,
					IsStatic:   isStatic,
					IsComputed: isComputed,
					IsPrivate:  isPrivate,
					IsReadonly: isReadonly,
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
