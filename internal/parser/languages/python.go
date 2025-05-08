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
	pyParser.baseParser = parser.NewTreeSitterParser(GetPythonLanguage(), pyParser.ParseTreeNode)
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

// ParseTreeNode разбирает узлы дерева Python кода
func (p *PythonParser) ParseTreeNode(node *sitter.Node, structure *models.CodeStructure, content []byte) error {
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
				if err := p.ParseTreeNode(child, structure, content); err != nil {
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
		isDynamic := false
		isNamespace := false
		isTypeImport := false

		if currentNode.Type() == "dotted_name" {
			path = currentNode.Content(content)
			if cursor.GoToNextSibling() && cursor.CurrentNode().Type() == "as" {
				if cursor.GoToNextSibling() && cursor.CurrentNode().Type() == "identifier" {
					alias = cursor.CurrentNode().Content(content)
				}
			} else {
				cursor.GoToParent()
				cursor.GoToNextSibling()
			}
			structure.AddImport(&models.Import{
				Path:         path,
				Alias:        alias,
				IsDynamic:    isDynamic,
				IsNamespace:  isNamespace,
				IsTypeImport: isTypeImport,
				Position:     getNodePosition(currentNode),
			})
		} else if currentNode.Type() == "aliased_import" {
			nameNode := currentNode.ChildByFieldName("name")
			aliasNode := currentNode.ChildByFieldName("alias")
			if nameNode != nil && aliasNode != nil {
				structure.AddImport(&models.Import{
					Path:         nameNode.Content(content),
					Alias:        aliasNode.Content(content),
					IsDynamic:    isDynamic,
					IsNamespace:  isNamespace,
					IsTypeImport: isTypeImport,
					Position:     getNodePosition(currentNode),
				})
			}
		} else if currentNode.Type() == "," {
			continue
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

		if currentNode.Type() == "wildcard_import" {
			structure.AddImport(&models.Import{
				Path:         modulePath + ".*",
				Alias:        "*",
				IsDynamic:    false,
				IsNamespace:  true,
				IsTypeImport: false,
				Position:     getNodePosition(currentNode),
			})
		} else if currentNode.Type() == "dotted_name" {
			name := currentNode.Content(content)
			alias := ""
			if cursor.GoToNextSibling() && cursor.CurrentNode().Type() == "as" {
				if cursor.GoToNextSibling() && cursor.CurrentNode().Type() == "identifier" {
					alias = cursor.CurrentNode().Content(content)
				}
			} else {
				cursor.GoToParent()
				cursor.GoToNextSibling()
			}
			structure.AddImport(&models.Import{
				Path:         modulePath + "." + name,
				Alias:        alias,
				IsDynamic:    false,
				IsNamespace:  false,
				IsTypeImport: false,
				Position:     getNodePosition(currentNode),
			})
		} else if currentNode.Type() == "aliased_import" {
			nameNode := currentNode.ChildByFieldName("name")
			aliasNode := currentNode.ChildByFieldName("alias")
			if nameNode != nil && aliasNode != nil {
				structure.AddImport(&models.Import{
					Path:         modulePath + "." + nameNode.Content(content),
					Alias:        aliasNode.Content(content),
					IsDynamic:    false,
					IsNamespace:  false,
					IsTypeImport: false,
					Position:     getNodePosition(currentNode),
				})
			}
		} else if currentNode.Type() == "(" || currentNode.Type() == ")" || currentNode.Type() == "," {
			continue
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
	isPublic := !strings.HasPrefix(funcName, "_")
	isAsync := false
	isGenerator := false
	isArrow := false // Python не имеет стрелочных функций
	isIIFE := false

	// Проверяем декораторы на async
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child.Type() == "decorator" {
			decoratorName := child.ChildByFieldName("name")
			if decoratorName != nil && decoratorName.Content(content) == "async" {
				isAsync = true
			}
		}
	}

	// Извлекаем параметры
	paramsNode := node.ChildByFieldName("parameters")
	params := p.parseParameters(paramsNode, content)

	// Извлекаем возвращаемый тип (если есть аннотация)
	var returnType string
	returnTypeNode := node.ChildByFieldName("return_type")
	if returnTypeNode != nil {
		returnType = returnTypeNode.Content(content)
	}

	startLine := int(node.StartPoint().Row)
	startCol := int(node.StartPoint().Column)
	endLine := int(node.EndPoint().Row)
	endCol := int(node.EndPoint().Column)

	if classModel != nil {
		// Это метод класса
		method := &models.Method{
			Name:          funcName,
			IsPublic:      isPublic,
			IsAsync:       isAsync,
			IsGenerator:   isGenerator,
			IsDecorator:   false, // TODO: Определить декораторы
			IsConstructor: funcName == "__init__",
			Kind:          "method",
			BelongsTo:     classModel.Name,
			Parameters:    params,
			ReturnType:    returnType,
			Position: models.Position{
				StartLine:   startLine + 1,
				StartColumn: startCol + 1,
				EndLine:     endLine + 1,
				EndColumn:   endCol + 1,
			},
		}
		classModel.Methods = append(classModel.Methods, method)
	} else {
		// Это функция верхнего уровня
		fn := &models.Function{
			Name:        funcName,
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

		// Добавляем публичные функции в экспорты
		if isPublic {
			structure.AddExport(&models.Export{
				Name:         fn.Name,
				Type:         "function",
				IsDefault:    false,
				IsTypeExport: false,
				IsNamespace:  false,
				Position:     fn.Position,
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
	isAbstract := false
	isInterface := false
	isMixin := false
	isGeneric := false
	isEnum := false

	// Извлекаем родительские классы
	var parent string
	var implements []string
	argumentsNode := node.ChildByFieldName("arguments")
	if argumentsNode != nil {
		cursor := sitter.NewTreeCursor(argumentsNode)
		defer cursor.Close()
		if cursor.GoToFirstChild() {
			for {
				currentNode := cursor.CurrentNode()
				if currentNode.Type() == "identifier" {
					if parent == "" {
						parent = currentNode.Content(content)
					} else {
						implements = append(implements, currentNode.Content(content))
					}
				}
				if !cursor.GoToNextSibling() {
					break
				}
			}
		}
	}

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

	// Парсим тело класса
	bodyNode := node.ChildByFieldName("body")
	if bodyNode != nil && bodyNode.Type() == "block" {
		for i := 0; i < int(bodyNode.ChildCount()); i++ {
			child := bodyNode.Child(i)
			if child.Type() == "function_definition" {
				p.parseFunctionOrMethod(child, structure, content, classModel)
			} else if child.Type() == "assignment" || child.Type() == "typed_assignment" {
				p.parseClassAssignment(child, classModel, content)
			}
		}
	}

	structure.AddType(classModel)
	// Добавляем публичные классы в экспорты модуля
	if isPublic {
		structure.AddExport(&models.Export{
			Name:         classModel.Name,
			Type:         "class",
			IsDefault:    false,
			IsTypeExport: false,
			IsNamespace:  false,
			Position:     classModel.Position,
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

		startLine := int(leftNode.StartPoint().Row)
		startCol := int(leftNode.StartPoint().Column)
		endLine := int(leftNode.EndPoint().Row)
		endCol := int(leftNode.EndPoint().Column)

		variable := &models.Variable{
			Name:     varName,
			IsPublic: isPublic,
			Type:     "", // TODO: Попытаться извлечь тип из typed_assignment или аннотаций
			Position: models.Position{
				StartLine:   startLine + 1,
				StartColumn: startCol + 1,
				EndLine:     endLine + 1,
				EndColumn:   endCol + 1,
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
		return
	}

	propName := leftNode.Content(content)
	isPublic := !strings.HasPrefix(propName, "_")
	isStatic := false
	isComputed := false
	isPrivate := strings.HasPrefix(propName, "__")
	isReadonly := false

	// TODO: Определить, является ли свойство статическим или readonly
	// В Python это обычно определяется через декораторы или соглашения

	startLine := int(leftNode.StartPoint().Row)
	startCol := int(leftNode.StartPoint().Column)
	endLine := int(node.EndPoint().Row)
	endCol := int(node.EndPoint().Column)

	prop := &models.Property{
		Name:       propName,
		IsPublic:   isPublic,
		IsStatic:   isStatic,
		IsComputed: isComputed,
		IsPrivate:  isPrivate,
		IsReadonly: isReadonly,
		Type:       "", // TODO: Извлечь тип из аннотации
		Position: models.Position{
			StartLine:   startLine + 1,
			StartColumn: startCol + 1,
			EndLine:     endLine + 1,
			EndColumn:   endCol + 1,
		},
	}
	classModel.Properties = append(classModel.Properties, prop)
}

// parseParameters извлекает параметры функции/метода
func (p *PythonParser) parseParameters(node *sitter.Node, content []byte) []*models.Parameter {
	params := make([]*models.Parameter, 0)
	if node == nil || node.Type() != "parameters" {
		return params
	}

	cursor := sitter.NewTreeCursor(node)
	defer cursor.Close()
	if !cursor.GoToFirstChild() {
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

		nodeType := currentNode.Type()

		if nodeType == "identifier" {
			paramName = currentNode.Content(content)
		} else if nodeType == "typed_parameter" {
			nameNode := currentNode.ChildByFieldName("name")
			typeNode := currentNode.ChildByFieldName("type")
			if nameNode != nil {
				paramName = nameNode.Content(content)
			}
			if typeNode != nil {
				paramType = typeNode.Content(content)
			}
		} else if nodeType == "default_parameter" {
			isRequired = false
			nameNode := currentNode.ChildByFieldName("name")
			typeNode := currentNode.ChildByFieldName("type")
			valueNode := currentNode.ChildByFieldName("value")
			if nameNode != nil {
				paramName = nameNode.Content(content)
			}
			if typeNode != nil {
				paramType = typeNode.Content(content)
			}
			if valueNode != nil {
				defaultValue = valueNode.Content(content)
			}
		} else if nodeType == "list_splat_pattern" || nodeType == "dictionary_splat_pattern" {
			nameNode := findFirstChildOfType(currentNode, "identifier")
			if nameNode != nil {
				paramName = nameNode.Content(content)
			}
			isRequired = false
		} else if nodeType == "," || nodeType == "(" || nodeType == ")" {
			continue
		}

		if firstParam && (paramName == "self" || paramName == "cls") {
			isSelfOrCls = true
		}
		firstParam = false

		if paramName != "" && !isSelfOrCls {
			params = append(params, &models.Parameter{
				Name:                 paramName,
				Type:                 paramType,
				IsRequired:           isRequired,
				DefaultValue:         defaultValue,
				IsVariadic:           nodeType == "list_splat_pattern" || nodeType == "dictionary_splat_pattern",
				IsDestructuredObject: false,
				IsDestructuredArray:  false,
			})
		}

		if !cursor.GoToNextSibling() {
			break
		}
	}

	return params
}
