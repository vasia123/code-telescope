package languages

import (
	"fmt"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/python"

	"code-telescope/internal/config"
	"code-telescope/internal/parser"
	"code-telescope/pkg/models"
)

// pyLanguage синглтон для языка Python
var pyLanguage *sitter.Language

func init() {
	pyLanguage = python.GetLanguage()

	// Регистрация парсера
	extensions := []string{".py", ".pyw"}
	parser.RegisterParser("Python", extensions, func(cfg *config.Config) parser.Parser {
		return NewPythonParser(cfg)
	})
}

// GetPythonLanguage возвращает инициализированный язык Python для tree-sitter
func GetPythonLanguage() *sitter.Language {
	return pyLanguage
}

// PythonParser реализует интерфейс parser.Parser для языка Python
type PythonParser struct {
	baseParser *parser.TreeSitterParser
	config     *config.Config
}

// NewPythonParser создает новый экземпляр парсера Python
func NewPythonParser(cfg *config.Config) parser.Parser { // Возвращаем интерфейс
	pyParser := &PythonParser{
		config: cfg,
	}
	pyParser.baseParser = parser.NewTreeSitterParser(GetPythonLanguage(), pyParser.parseTreeNode)
	return pyParser
}

// Parse вызывает базовый парсер
func (p *PythonParser) Parse(fileMetadata *models.FileMetadata) (*models.CodeStructure, error) {
	return p.baseParser.Parse(fileMetadata)
}

// GetLanguageName возвращает название языка программирования
func (p *PythonParser) GetLanguageName() string {
	return "Python"
}

// GetSupportedExtensions возвращает список поддерживаемых расширений файлов
func (p *PythonParser) GetSupportedExtensions() []string {
	return []string{".py", ".pyw"}
}

// parseTreeNode разбирает узлы дерева Python кода
func (p *PythonParser) parseTreeNode(node *sitter.Node, structure *models.CodeStructure, content []byte) error {
	// Рекурсивно обходим дочерние узлы
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		switch child.Type() {
		case "import_statement":
			p.parseImport(child, structure, content)
		case "import_from_statement":
			p.parseImportFrom(child, structure, content)
		case "function_definition":
			p.parseFunctionOrMethod(child, structure, content, nil)
		case "class_definition":
			p.parseClass(child, structure, content)
		case "expression_statement":
			// Некоторые глобальные переменные и модульные документации могут быть здесь
			if child.ChildCount() > 0 && child.Child(0).Type() == "assignment" {
				p.parseAssignment(child.Child(0), structure, content)
			}
		default:
			// Рекурсивно обходим узлы, которые могут содержать нужные декларации
			if child.ChildCount() > 0 {
				if err := p.parseTreeNode(child, structure, content); err != nil {
					fmt.Printf("Error parsing child node: %v\n", err) // Заменить на логгер
				}
			}
		}
	}
	return nil
}

// parseImport извлекает импорты
func (p *PythonParser) parseImport(node *sitter.Node, structure *models.CodeStructure, content []byte) {
	cursor := sitter.NewTreeCursor(node)
	defer cursor.Close()

	if !cursor.GoToFirstChild() {
		return
	}

	for {
		currentNode := cursor.CurrentNode()
		alias := ""
		path := ""

		if currentNode.Type() == "dotted_name" { // import module
			path = currentNode.Content(content)
			// Проверяем следующий узел на "as"
			if cursor.GoToNextSibling() && cursor.CurrentNode().Type() == "as" {
				if cursor.GoToNextSibling() && cursor.CurrentNode().Type() == "identifier" {
					alias = cursor.CurrentNode().Content(content)
				} else {
					continue // Ошибка в структуре `as`
				}
			} else {
				cursor.GoToParent()      // Возвращаемся, чтобы продолжить с этого узла
				cursor.GoToNextSibling() // Переходим к следующему элементу импорта (если есть)
			}
			structure.AddImport(&models.Import{
				Path:     path,
				Alias:    alias,
				Position: getNodePosition(currentNode), // Позиция имени модуля
			})
		} else if currentNode.Type() == "aliased_import" { // import module as alias
			nameNode := currentNode.ChildByFieldName("name")
			aliasNode := currentNode.ChildByFieldName("alias")
			if nameNode != nil && aliasNode != nil {
				structure.AddImport(&models.Import{
					Path:     nameNode.Content(content),
					Alias:    aliasNode.Content(content),
					Position: getNodePosition(currentNode),
				})
			}
		} else if currentNode.Type() == "," {
			continue // Пропускаем запятую
		}

		if !cursor.GoToNextSibling() {
			break
		}
	}
}

// parseImportFrom извлекает импорты `from module import name`
func (p *PythonParser) parseImportFrom(node *sitter.Node, structure *models.CodeStructure, content []byte) {
	moduleNameNode := node.ChildByFieldName("module_name")
	if moduleNameNode == nil {
		return
	}
	modulePath := moduleNameNode.Content(content)

	// Ищем импортируемые имена (dotted_name или aliased_import)
	cursor := sitter.NewTreeCursor(node)
	defer cursor.Close()
	if !cursor.GoToFirstChild() {
		return
	}

	foundImportKeyword := false
	for {
		currentNode := cursor.CurrentNode()
		if currentNode.Type() == "import" {
			foundImportKeyword = true
			continue
		}
		if !foundImportKeyword {
			continue
		}

		// После 'import' идут импортируемые элементы
		if currentNode.Type() == "wildcard_import" { // from module import *
			structure.AddImport(&models.Import{
				Path:     modulePath + ".*",
				Alias:    "*",
				Position: getNodePosition(currentNode),
			})
		} else if currentNode.Type() == "dotted_name" { // from module import name
			name := currentNode.Content(content)
			alias := ""
			// Проверяем следующий узел на 'as'
			if cursor.GoToNextSibling() && cursor.CurrentNode().Type() == "as" {
				if cursor.GoToNextSibling() && cursor.CurrentNode().Type() == "identifier" {
					alias = cursor.CurrentNode().Content(content)
				}
			} else {
				cursor.GoToParent()
				cursor.GoToNextSibling()
			}
			structure.AddImport(&models.Import{
				Path:     modulePath + "." + name,
				Alias:    alias,
				Position: getNodePosition(currentNode),
			})
		} else if currentNode.Type() == "aliased_import" { // from module import name as alias
			nameNode := currentNode.ChildByFieldName("name")
			aliasNode := currentNode.ChildByFieldName("alias")
			if nameNode != nil && aliasNode != nil {
				structure.AddImport(&models.Import{
					Path:     modulePath + "." + nameNode.Content(content),
					Alias:    aliasNode.Content(content),
					Position: getNodePosition(currentNode),
				})
			}
		} else if currentNode.Type() == "(" || currentNode.Type() == ")" || currentNode.Type() == "," {
			continue // Пропускаем скобки и запятые
		}

		if !cursor.GoToNextSibling() {
			break
		}
	}
}

// parseFunctionOrMethod извлекает функции и методы
func (p *PythonParser) parseFunctionOrMethod(node *sitter.Node, structure *models.CodeStructure, content []byte, classModel *models.Type) {
	nameNode := node.ChildByFieldName("name")
	if nameNode == nil {
		return
	}
	funcName := nameNode.Content(content)

	// Определение публичности (начинается с _ или __?)
	isPublic := !strings.HasPrefix(funcName, "_")
	isMethod := classModel != nil
	// В Python self/cls обычно первый параметр метода
	paramsNode := node.ChildByFieldName("parameters")
	params := p.parseParameters(paramsNode, content, isMethod)

	// Тип возвращаемого значения
	returnType := ""
	returnTypeNode := node.ChildByFieldName("return_type")
	if returnTypeNode != nil {
		returnType = returnTypeNode.Content(content)
	}

	startLine, startCol := node.StartPosition().Line, node.StartPosition().Column
	endLine, endCol := node.EndPosition().Line, node.EndPosition().Column

	if isMethod {
		// Это метод класса
		method := &models.Method{
			Name:       funcName,
			IsPublic:   isPublic,
			IsStatic:   false, // TODO: Определить staticmethod/classmethod декораторы
			Kind:       "method",
			BelongsTo:  classModel.Name,
			Parameters: params,
			ReturnType: returnType,
			Position: models.Position{
				StartLine: int(startLine) + 1,
				StartCol:  int(startCol) + 1,
				EndLine:   int(endLine) + 1,
				EndCol:    int(endCol) + 1,
			},
		}
		classModel.Methods = append(classModel.Methods, method)
	} else {
		// Это функция верхнего уровня
		fn := &models.Function{
			Name:       funcName,
			IsPublic:   isPublic,
			Parameters: params,
			ReturnType: returnType,
			Position: models.Position{
				StartLine: int(startLine) + 1,
				StartCol:  int(startCol) + 1,
				EndLine:   int(endLine) + 1,
				EndCol:    int(endCol) + 1,
			},
		}
		structure.AddFunction(fn)
		// Добавляем публичные функции в экспорты модуля
		if isPublic {
			structure.AddExport(&models.Export{
				Name:     fn.Name,
				Type:     "function",
				Position: fn.Position,
			})
		}
	}
}

// parseClass извлекает классы
func (p *PythonParser) parseClass(node *sitter.Node, structure *models.CodeStructure, content []byte) {
	nameNode := node.ChildByFieldName("name")
	if nameNode == nil {
		return
	}
	className := nameNode.Content(content)
	isPublic := !strings.HasPrefix(className, "_")

	// TODO: Извлечь родительские классы (superclasses)

	startLine, startCol := node.StartPosition().Line, node.StartPosition().Column
	endLine, endCol := node.EndPosition().Line, node.EndPosition().Column

	classModel := &models.Type{
		Name:     className,
		IsPublic: isPublic,
		Kind:     "class",
		Position: models.Position{
			StartLine: int(startLine) + 1,
			StartCol:  int(startCol) + 1,
			EndLine:   int(endLine) + 1,
			EndCol:    int(endCol) + 1,
		},
		Methods:    make([]*models.Method, 0),
		Properties: make([]*models.Property, 0),
	}

	// Парсим тело класса (блок)
	bodyNode := node.ChildByFieldName("body")
	if bodyNode != nil && bodyNode.Type() == "block" {
		// Обходим узлы внутри блока класса
		for i := 0; i < int(bodyNode.ChildCount()); i++ {
			child := bodyNode.Child(i)
			if child.Type() == "function_definition" {
				p.parseFunctionOrMethod(child, structure, content, classModel)
			} else if child.Type() == "assignment" || child.Type() == "typed_assignment" {
				// TODO: Обработать присваивания внутри класса (атрибуты)
				p.parseClassAssignment(child, classModel, content)
			}
		}
	}

	structure.AddType(classModel)
	// Добавляем публичные классы в экспорты модуля
	if isPublic {
		structure.AddExport(&models.Export{
			Name:     classModel.Name,
			Type:     "class",
			Position: classModel.Position,
		})
	}
}

// parseAssignment извлекает переменные верхнего уровня
func (p *PythonParser) parseAssignment(node *sitter.Node, structure *models.CodeStructure, content []byte) {
	// Интересуют только присваивания на верхнем уровне (не внутри функций или классов)
	// TODO: Добавить проверку на уровень вложенности, если необходимо

	leftNode := node.ChildByFieldName("left")
	// rightNode := node.ChildByFieldName("right") // Правая часть пока не нужна
	// typeNode := node.ChildByFieldName("type") // Для typed_assignment

	if leftNode != nil && leftNode.Type() == "identifier" {
		varName := leftNode.Content(content)
		isPublic := !strings.HasPrefix(varName, "_")

		startLine, startCol := leftNode.StartPosition().Line, leftNode.StartPosition().Column
		endLine, endCol := leftNode.EndPosition().Line, leftNode.EndPosition().Column

		variable := &models.Variable{
			Name:     varName,
			IsPublic: isPublic,
			Type:     "", // TODO: Попытаться извлечь тип из typed_assignment или аннотаций
			Position: models.Position{
				StartLine: int(startLine) + 1,
				StartCol:  int(startCol) + 1,
				EndLine:   int(endLine) + 1,
				EndCol:    int(endCol) + 1,
			},
		}
		structure.AddVariable(variable)

		// Добавляем публичные переменные в экспорты
		if isPublic {
			structure.AddExport(&models.Export{
				Name:     variable.Name,
				Type:     "variable",
				Position: variable.Position,
			})
		}
	}
	// TODO: Обработать множественные присваивания (a, b = 1, 2) и присваивания атрибутам/индексам
}

// parseClassAssignment обрабатывает присваивания внутри класса (для атрибутов)
func (p *PythonParser) parseClassAssignment(node *sitter.Node, classModel *models.Type, content []byte) {
	leftNode := node.ChildByFieldName("left")
	if leftNode == nil || leftNode.Type() != "identifier" {
		return // Интересуют только простые атрибуты вида name = ...
	}

	propName := leftNode.Content(content)
	isPublic := !strings.HasPrefix(propName, "_")

	startLine, startCol := leftNode.StartPosition().Line, leftNode.StartPosition().Column
	endLine, endCol := leftNode.EndPosition().Line, node.EndPosition().Column

	prop := &models.Property{
		Name:     propName,
		IsPublic: isPublic,
		IsStatic: false, // Атрибуты экземпляра по умолчанию
		Type:     "",    // TODO: Извлечь тип из аннотации (typed_assignment)
		Position: models.Position{
			StartLine: int(startLine) + 1,
			StartCol:  int(startCol) + 1,
			EndLine:   int(endLine) + 1,
			EndCol:    int(endCol) + 1,
		},
	}
	classModel.Properties = append(classModel.Properties, prop)
}

// parseParameters извлекает параметры функции/метода
func (p *PythonParser) parseParameters(node *sitter.Node, content []byte, isMethod bool) []*models.Parameter {
	params := make([]*models.Parameter, 0)
	if node == nil || node.Type() != "parameters" {
		return params
	}

	cursor := sitter.NewTreeCursor(node)
	defer cursor.Close()
	if !cursor.GoToFirstChild() { // Пропускаем ( и )
		return params
	}

	firstParam := true
	for {
		currentNode := cursor.CurrentNode()
		var paramName string
		var paramType string = ""
		var defaultValue string = ""
		isRequired := true
		isSelfOrCls := false

		// Обрабатываем разные типы параметров: identifier, typed_parameter, default_parameter, etc.
		nodeType := currentNode.Type()

		if nodeType == "identifier" { // Обычный параметр без типа и значения по умолчанию
			paramName = currentNode.Content(content)
		} else if nodeType == "typed_parameter" { // Параметр с типом: name: type
			nameNode := currentNode.ChildByFieldName("name")
			typeNode := currentNode.ChildByFieldName("type")
			if nameNode != nil {
				paramName = nameNode.Content(content)
			}
			if typeNode != nil {
				paramType = typeNode.Content(content)
			}
		} else if nodeType == "default_parameter" { // Параметр со значением по умолчанию: name=value или name:type=value
			isRequired = false
			nameNode := currentNode.ChildByFieldName("name")
			typeNode := currentNode.ChildByFieldName("type")   // Опционально
			valueNode := currentNode.ChildByFieldName("value") // Обязательно
			if nameNode != nil {
				paramName = nameNode.Content(content)
			}
			if typeNode != nil {
				paramType = typeNode.Content(content)
			}
			if valueNode != nil {
				defaultValue = valueNode.Content(content)
			}
		} else if nodeType == "list_splat_pattern" || nodeType == "dictionary_splat_pattern" { // *args, **kwargs
			// Ищем identifier внутри
			nameNode := findFirstChildOfType(currentNode, "identifier")
			if nameNode != nil {
				paramName = nameNode.Content(content)
			}
			isRequired = false // Обычно не считаются обязательными в сигнатуре
		} else if nodeType == "," || nodeType == "(" || nodeType == ")" {
			continue // Пропускаем разделители и скобки
		}

		// Проверяем, является ли первый параметр метода self или cls
		if isMethod && firstParam && (paramName == "self" || paramName == "cls") {
			isSelfOrCls = true
		}
		firstParam = false // Сбрасываем флаг после первого параметра

		// Добавляем параметр, если это не self/cls
		if paramName != "" && !isSelfOrCls {
			params = append(params, &models.Parameter{
				Name:         paramName,
				Type:         paramType,
				IsRequired:   isRequired,
				DefaultValue: defaultValue,
				// IsVariadic определить сложнее, зависит от * или ** перед именем
			})
		}

		if !cursor.GoToNextSibling() {
			break
		}
	}

	return params
}
