package languages

import (
	"fmt"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/javascript"

	"code-telescope/internal/config"
	"code-telescope/internal/parser"
	"code-telescope/pkg/models"
)

// jsLanguage синглтон для языка JavaScript
var jsLanguage *sitter.Language

func init() {
	jsLanguage = javascript.GetLanguage()

	// Регистрация парсера
	extensions := []string{".js", ".jsx", ".mjs", ".cjs"}
	parser.RegisterParser("JavaScript", extensions, func(cfg *config.Config) parser.Parser {
		return NewJavaScriptParser(cfg)
	})
}

// GetJavaScriptLanguage возвращает инициализированный язык JavaScript для tree-sitter
func GetJavaScriptLanguage() *sitter.Language {
	return jsLanguage
}

// JavaScriptParser реализует интерфейс parser.Parser для языка JavaScript
type JavaScriptParser struct {
	baseParser *parser.TreeSitterParser
	config     *config.Config
}

// NewJavaScriptParser создает новый экземпляр парсера JavaScript
func NewJavaScriptParser(cfg *config.Config) parser.Parser { // Возвращаем интерфейс
	jsParser := &JavaScriptParser{
		config: cfg,
	}
	jsParser.baseParser = parser.NewTreeSitterParser(GetJavaScriptLanguage(), jsParser.ParseTreeNode)
	return jsParser
}

// Parse вызывает базовый парсер
func (p *JavaScriptParser) Parse(fileMetadata *models.FileMetadata) (*models.CodeStructure, error) {
	return p.baseParser.Parse(fileMetadata)
}

// GetLanguageName возвращает название языка программирования
func (p *JavaScriptParser) GetLanguageName() string {
	return "JavaScript"
}

// GetSupportedExtensions возвращает список поддерживаемых расширений файлов
func (p *JavaScriptParser) GetSupportedExtensions() []string {
	return []string{".js", ".jsx", ".mjs", ".cjs"}
}

// ParseTreeNode разбирает узлы дерева JavaScript кода
func (p *JavaScriptParser) ParseTreeNode(node *sitter.Node, structure *models.CodeStructure, content []byte) error {
	// Рекурсивно обходим дочерние узлы
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		switch child.Type() {
		case "import_statement":
			p.parseImport(child, structure, content)
		case "export_statement":
			p.parseExport(child, structure, content)
		case "function_declaration", "generator_function_declaration":
			p.parseFunction(child, structure, content, false) // isMethod = false
		case "class_declaration":
			p.parseClass(child, structure, content)
		case "variable_declaration": // var foo = ...
			p.parseVariableOrLexicalDeclaration(child, structure, content)
		case "lexical_declaration": // let foo = ..., const bar = ...
			p.parseVariableOrLexicalDeclaration(child, structure, content)
		default:
			// Рекурсивно обходим узлы, которые могут содержать нужные декларации
			if child.ChildCount() > 0 {
				if err := p.ParseTreeNode(child, structure, content); err != nil {
					fmt.Printf("Error parsing child node: %v\n", err) // Заменить на логгер
				}
			}
		}
	}
	return nil
}

// parseImport извлекает импорты
func (p *JavaScriptParser) parseImport(node *sitter.Node, structure *models.CodeStructure, content []byte) {
	// Находим source (строку с путем импорта)
	sourceNode := node.ChildByFieldName("source")
	if sourceNode == nil || sourceNode.Type() != "string" {
		// Может быть импорт без source, например, `import "./styles.css";`
		// Или импорт метаданных `import.meta.url`
		// Пока игнорируем такие случаи для простоты
		return
	}

	path := strings.Trim(sourceNode.Content(content), `"'`)
	isDynamic := false // TODO
	isNamespace := false
	isTypeImport := false // TODO

	// Проверяем тип импорта
	if node.ChildCount() > 0 {
		firstChild := node.Child(0)
		if firstChild.Type() == "import" {
			// Проверяем на namespace import
			if node.ChildCount() > 1 && node.Child(1).Type() == "*" {
				isNamespace = true
			}
		}
	}

	startLine := int(node.StartPoint().Row)
	startCol := int(node.StartPoint().Column)
	endLine := int(node.EndPoint().Row)
	endCol := int(node.EndPoint().Column)

	imp := &models.Import{
		Path:         path,
		IsDynamic:    isDynamic,
		IsNamespace:  isNamespace,
		IsTypeImport: isTypeImport,
		Position: models.Position{
			StartLine:   startLine + 1,
			StartColumn: startCol + 1,
			EndLine:     endLine + 1,
			EndColumn:   endCol + 1,
		},
	}

	structure.AddImport(imp)
}

// parseExport извлекает экспорты
func (p *JavaScriptParser) parseExport(node *sitter.Node, structure *models.CodeStructure, content []byte) {
	var exportedName string
	var exportType string = "unknown"
	var posNode *sitter.Node = node
	isDefault := false
	isTypeExport := false
	isNamespace := false

	valueNode := node.ChildByFieldName("value")
	declarationNode := node.ChildByFieldName("declaration")

	if declarationNode != nil {
		posNode = declarationNode
		switch declarationNode.Type() {
		case "function_declaration", "generator_function_declaration":
			nameNode := declarationNode.ChildByFieldName("name")
			if nameNode != nil {
				exportedName = nameNode.Content(content)
				exportType = "function"
				p.parseFunction(declarationNode, structure, content, false)
			}
		case "class_declaration":
			nameNode := declarationNode.ChildByFieldName("name")
			if nameNode != nil {
				exportedName = nameNode.Content(content)
				exportType = "class"
				p.parseClass(declarationNode, structure, content)
			}
		case "lexical_declaration", "variable_declaration":
			declarator := findFirstChildOfType(declarationNode, "variable_declarator")
			if declarator != nil {
				nameNode := declarator.ChildByFieldName("name")
				if nameNode != nil {
					exportedName = nameNode.Content(content)
					exportType = "variable"
					p.parseVariableOrLexicalDeclaration(declarationNode, structure, content)
				}
			}
		}
	} else if valueNode != nil && valueNode.Type() == "identifier" && node.ChildCount() > 1 && node.Child(1).Type() == "default" {
		exportedName = valueNode.Content(content)
		exportType = "default_identifier"
		isDefault = true
		posNode = valueNode
	} else if node.ChildCount() > 1 && node.Child(1).Type() == "default" {
		exportedName = "default"
		exportType = "default"
		isDefault = true
		posNode = node.Child(1)
	} else if node.ChildCount() > 0 && node.Child(0).Type() == "export_clause" {
		exportClause := node.Child(0)
		cursor := sitter.NewTreeCursor(exportClause)
		defer cursor.Close()

		if cursor.GoToFirstChild() {
			for {
				if cursor.CurrentNode().Type() == "export_specifier" {
					nameNode := cursor.CurrentNode().ChildByFieldName("name")
					if nameNode != nil {
						structure.AddExport(&models.Export{
							Name:         nameNode.Content(content),
							Type:         "named",
							IsDefault:    false,
							IsTypeExport: false,
							IsNamespace:  false,
							Position:     getNodePosition(cursor.CurrentNode()),
						})
					}
				}

				if !cursor.GoToNextSibling() {
					break
				}
			}
		}
		return
	}

	if exportedName != "" {
		startLine := int(posNode.StartPoint().Row)
		startCol := int(posNode.StartPoint().Column)
		endLine := int(posNode.EndPoint().Row)
		endCol := int(posNode.EndPoint().Column)

		structure.AddExport(&models.Export{
			Name:         exportedName,
			Type:         exportType,
			IsDefault:    isDefault,
			IsTypeExport: isTypeExport,
			IsNamespace:  isNamespace,
			Position: models.Position{
				StartLine:   startLine + 1,
				StartColumn: startCol + 1,
				EndLine:     endLine + 1,
				EndColumn:   endCol + 1,
			},
		})
	}
}

// parseFunction извлекает функции (включая методы внутри классов)
func (p *JavaScriptParser) parseFunction(node *sitter.Node, structure *models.CodeStructure, content []byte, isMethod bool) {
	nameNode := node.ChildByFieldName("name")
	if nameNode == nil {
		// Может быть анонимная функция, присвоенная переменной, или метод класса без явного имени
		return
	}

	funcName := nameNode.Content(content)
	isPublic := true // В JS функции и методы всегда публичные, если они не скрыты замыканием
	isAsync := false
	isGenerator := false
	isArrow := node.Type() == "arrow_function"
	isIIFE := false // TODO: Определить IIFE

	// Проверяем модификаторы
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child.Type() == "async" {
			isAsync = true
		} else if child.Type() == "generator" {
			isGenerator = true
		}
	}

	paramsNode := node.ChildByFieldName("parameters")
	params := p.parseParameters(paramsNode, content)

	startLine := int(node.StartPoint().Row)
	startCol := int(node.StartPoint().Column)
	endLine := int(node.EndPoint().Row)
	endCol := int(node.EndPoint().Column)

	if !isMethod {
		// Это функция верхнего уровня
		fn := &models.Function{
			Name:        funcName,
			IsPublic:    isPublic,
			IsAsync:     isAsync,
			IsGenerator: isGenerator,
			IsArrow:     isArrow,
			IsIIFE:      isIIFE,
			Parameters:  params,
			ReturnType:  "", // В JS тип возвращаемого значения обычно не указывается статически
			Position: models.Position{
				StartLine:   startLine + 1,
				StartColumn: startCol + 1,
				EndLine:     endLine + 1,
				EndColumn:   endCol + 1,
			},
		}
		structure.AddFunction(fn)
	}
	// Если это метод, он будет добавлен при парсинге тела класса
}

// parseClass извлекает классы
func (p *JavaScriptParser) parseClass(node *sitter.Node, structure *models.CodeStructure, content []byte) {
	nameNode := node.ChildByFieldName("name")
	if nameNode == nil {
		return
	}

	className := nameNode.Content(content)
	isPublic := true    // Считаем публичным по умолчанию
	isAbstract := false // TODO: Определить абстрактный класс
	isInterface := false
	isMixin := false
	isGeneric := false
	isEnum := false

	// TODO: Извлечь родительские классы и интерфейсы
	var parent string
	var implements []string

	startLine := int(node.StartPoint().Row)
	startCol := int(node.StartPoint().Column)
	endLine := int(node.EndPoint().Row)
	endCol := int(node.EndPoint().Column)

	classModel := &models.Type{
		Name:              className,
		IsPublic:          isPublic,
		IsAbstract:        isAbstract,
		IsInterface:       isInterface,
		IsMixin:           isMixin,
		IsGeneric:         isGeneric,
		IsEnum:            isEnum,
		Kind:              "class",
		Parent:            parent,
		Implements:        implements,
		GenericParameters: make([]string, 0),
		Position: models.Position{
			StartLine:   startLine + 1,
			StartColumn: startCol + 1,
			EndLine:     endLine + 1,
			EndColumn:   endCol + 1,
		},
		Methods:    make([]*models.Method, 0),
		Properties: make([]*models.Property, 0),
	}

	// Парсим тело класса для извлечения методов и свойств
	bodyNode := node.ChildByFieldName("body")
	if bodyNode != nil {
		p.parseClassBody(bodyNode, classModel, content)
	}

	structure.AddType(classModel)
}

// parseClassBody извлекает методы и свойства из тела класса
func (p *JavaScriptParser) parseClassBody(bodyNode *sitter.Node, classModel *models.Type, content []byte) {
	cursor := sitter.NewTreeCursor(bodyNode)
	defer cursor.Close()
	if !cursor.GoToFirstChild() { // Пропускаем { и }
		return
	}

	for {
		currentNode := cursor.CurrentNode()
		switch currentNode.Type() {
		case "method_definition":
			p.parseMethod(currentNode, classModel, content)
		case "field_definition", "public_field_definition": // field_definition - старая грамматика?
			p.parseField(currentNode, classModel, content)
		}

		if !cursor.GoToNextSibling() {
			break
		}
	}
}

// parseMethod извлекает метод класса
func (p *JavaScriptParser) parseMethod(node *sitter.Node, classModel *models.Type, content []byte) {
	nameNode := node.ChildByFieldName("name")
	if nameNode == nil {
		// Может быть конструктор
		if node.ChildCount() > 0 && node.Child(0).Content(content) == "constructor" {
			nameNode = node.Child(0)
		} else {
			return
		}
	}

	methodName := nameNode.Content(content)
	isPublic := true // Методы класса по умолчанию публичные (до private #)
	isStatic := false
	isAsync := false
	isGenerator := false
	isDecorator := false
	isConstructor := methodName == "constructor"
	kind := "method"

	// Проверяем модификаторы
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child.Type() == "static" {
			isStatic = true
		} else if child.Type() == "async" {
			isAsync = true
		} else if child.Type() == "generator" {
			isGenerator = true
		} else if child.Type() == "get" {
			kind = "getter"
		} else if child.Type() == "set" {
			kind = "setter"
		}
	}

	if strings.HasPrefix(methodName, "#") {
		isPublic = false
		methodName = strings.TrimPrefix(methodName, "#")
	}

	paramsNode := node.ChildByFieldName("parameters")
	params := p.parseParameters(paramsNode, content)

	startLine := int(node.StartPoint().Row)
	startCol := int(node.StartPoint().Column)
	endLine := int(node.EndPoint().Row)
	endCol := int(node.EndPoint().Column)

	method := &models.Method{
		Name:          methodName,
		IsPublic:      isPublic,
		IsStatic:      isStatic,
		IsAsync:       isAsync,
		IsGenerator:   isGenerator,
		IsDecorator:   isDecorator,
		IsConstructor: isConstructor,
		Kind:          kind,
		BelongsTo:     classModel.Name,
		Parameters:    params,
		ReturnType:    "", // Не извлекаем для JS
		Position: models.Position{
			StartLine:   startLine + 1,
			StartColumn: startCol + 1,
			EndLine:     endLine + 1,
			EndColumn:   endCol + 1,
		},
	}
	classModel.Methods = append(classModel.Methods, method)
}

// parseField извлекает свойство класса
func (p *JavaScriptParser) parseField(node *sitter.Node, classModel *models.Type, content []byte) {
	nameNode := node.ChildByFieldName("name")
	if nameNode == nil {
		return
	}

	fieldName := nameNode.Content(content)
	isPublic := true
	isStatic := false
	isComputed := false
	isPrivate := false
	isReadonly := false

	// Проверяем модификаторы
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child.Type() == "static" {
			isStatic = true
		}
	}

	if strings.HasPrefix(fieldName, "#") {
		isPublic = false
		isPrivate = true
		fieldName = strings.TrimPrefix(fieldName, "#")
	}

	startLine := int(node.StartPoint().Row)
	startCol := int(node.StartPoint().Column)
	endLine := int(node.EndPoint().Row)
	endCol := int(node.EndPoint().Column)

	prop := &models.Property{
		Name:       fieldName,
		IsPublic:   isPublic,
		IsStatic:   isStatic,
		IsComputed: isComputed,
		IsPrivate:  isPrivate,
		IsReadonly: isReadonly,
		Type:       "", // Тип поля редко указывается статически в JS
		Position: models.Position{
			StartLine:   startLine + 1,
			StartColumn: startCol + 1,
			EndLine:     endLine + 1,
			EndColumn:   endCol + 1,
		},
	}
	classModel.Properties = append(classModel.Properties, prop)
}

// parseVariableOrLexicalDeclaration извлекает переменные (var, let, const)
func (p *JavaScriptParser) parseVariableOrLexicalDeclaration(node *sitter.Node, structure *models.CodeStructure, content []byte) {
	// Ищем все variable_declarator внутри
	cursor := sitter.NewTreeCursor(node)
	defer cursor.Close()

	if !cursor.GoToFirstChild() {
		return
	}

	for {
		if cursor.CurrentNode().Type() == "variable_declarator" {
			declarator := cursor.CurrentNode()

			switch {
			case declarator.ChildCount() >= 2 &&
				declarator.Child(0).Type() == "identifier" &&
				declarator.Child(1).Type() == "=":
				// var/let/const name = value
				name := declarator.Child(0).Content(content)

				// Проверяем, не функция ли это (function expression или arrow function)
				if declarator.ChildCount() >= 3 {
					valueNode := declarator.Child(2)
					if valueNode.Type() == "function" || valueNode.Type() == "arrow_function" {
						// Это функция, добавляем как функцию
						structure.AddFunction(&models.Function{
							Name:       name,
							IsPublic:   true, // В JS всё публичное по умолчанию
							Parameters: p.parseParameters(findFirstChildOfType(valueNode, "formal_parameters"), content),
							Position:   getNodePosition(declarator),
						})

						// Если не начинается с _, то считаем публичным API
						if !strings.HasPrefix(name, "_") {
							structure.AddExport(&models.Export{
								Name:     name,
								Type:     "function",
								Position: getNodePosition(declarator),
							})
						}
					}
				}
			}
		}

		if !cursor.GoToNextSibling() {
			break
		}
	}
}

// parseParameters извлекает параметры функции/метода
func (p *JavaScriptParser) parseParameters(node *sitter.Node, content []byte) []*models.Parameter {
	params := make([]*models.Parameter, 0)
	if node == nil || node.Type() != "formal_parameters" {
		return params
	}

	cursor := sitter.NewTreeCursor(node)
	defer cursor.Close()
	if !cursor.GoToFirstChild() { // Пропускаем ( и )
		return params
	}

	for {
		currentNode := cursor.CurrentNode()
		var paramName string
		var defaultValue string = ""
		isRequired := true
		isRest := false
		isDestructuredObject := false
		isDestructuredArray := false

		switch currentNode.Type() {
		case "identifier": // Простой параметр: func(a)
			paramName = currentNode.Content(content)
		case "required_parameter": // func(a)
			patternNode := currentNode.ChildByFieldName("pattern")
			if patternNode != nil && patternNode.Type() == "identifier" {
				paramName = patternNode.Content(content)
			}
		case "optional_parameter": // func(a = 1)
			isRequired = false
			patternNode := currentNode.ChildByFieldName("pattern")
			valueNode := currentNode.ChildByFieldName("value")
			if patternNode != nil && patternNode.Type() == "identifier" {
				paramName = patternNode.Content(content)
			}
			if valueNode != nil {
				defaultValue = valueNode.Content(content)
			}
		case "rest_parameter": // func(...args)
			isRest = true
			nameNode := findFirstChildOfType(currentNode, "identifier")
			if nameNode != nil {
				paramName = nameNode.Content(content)
			}
		case "assignment_pattern": // func({a = 1})
			isRequired = false
			leftNode := currentNode.ChildByFieldName("left")
			rightNode := currentNode.ChildByFieldName("right")
			if leftNode != nil && leftNode.Type() == "identifier" {
				paramName = leftNode.Content(content)
			}
			if rightNode != nil {
				defaultValue = rightNode.Content(content)
			}
		case "object_pattern":
			isDestructuredObject = true
			paramName = currentNode.Content(content)
		case "array_pattern":
			isDestructuredArray = true
			paramName = currentNode.Content(content)
		}

		if paramName != "" {
			params = append(params, &models.Parameter{
				Name:                 paramName,
				Type:                 "", // Тип не извлекаем для JS
				IsRequired:           isRequired,
				DefaultValue:         defaultValue,
				IsVariadic:           isRest,
				IsDestructuredObject: isDestructuredObject,
				IsDestructuredArray:  isDestructuredArray,
			})
		}

		if !cursor.GoToNextSibling() {
			break
		}
	}

	return params
}

// findFirstChildOfType вспомогательная функция для поиска первого дочернего узла заданного типа
func findFirstChildOfType(node *sitter.Node, childType string) *sitter.Node {
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child != nil && child.Type() == childType {
			return child
		}
	}
	return nil
}

// getNodePosition вспомогательная функция для получения позиции узла
func getNodePosition(node *sitter.Node) models.Position {
	startLine := int(node.StartPoint().Row)
	startCol := int(node.StartPoint().Column)
	endLine := int(node.EndPoint().Row)
	endCol := int(node.EndPoint().Column)

	return models.Position{
		StartLine:   startLine + 1,
		StartColumn: startCol + 1,
		EndLine:     endLine + 1,
		EndColumn:   endCol + 1,
	}
}
