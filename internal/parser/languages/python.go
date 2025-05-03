package languages

import (
	"os"
	"strings"
	"regexp"

	"code-telescope/internal/config"
	"code-telescope/pkg/models"
)

// PythonParser реализует парсер для языка Python
type PythonParser struct {
	Config *config.Config
}

// NewPythonParser создает новый экземпляр парсера Python
func NewPythonParser(cfg *config.Config) *PythonParser {
	return &PythonParser{
		Config: cfg,
	}
}

// GetLanguageName возвращает название языка программирования
func (p *PythonParser) GetLanguageName() string {
	return "Python"
}

// GetSupportedExtensions возвращает список поддерживаемых расширений файлов
func (p *PythonParser) GetSupportedExtensions() []string {
	return []string{".py"}
}

// Parse разбирает файл Python и извлекает его структуру
func (p *PythonParser) Parse(fileMetadata *models.FileMetadata) (*models.CodeStructure, error) {
	// Создаем пустую структуру кода
	codeStructure := models.NewCodeStructure(fileMetadata)

	// Читаем содержимое файла
	content, err := os.ReadFile(fileMetadata.AbsolutePath)
	if err != nil {
		return nil, err
	}

	// TODO: Использовать Tree-sitter для парсинга Python файла
	// В настоящее время используем простой подход с поиском маркеров

	// Разбиваем содержимое на строки
	lines := strings.Split(string(content), "\n")

	// Обрабатываем импорты
	p.parseImports(lines, codeStructure)

	// Обрабатываем классы и методы
	p.parseClassesAndMethods(lines, codeStructure)

	// Обрабатываем функции верхнего уровня
	p.parseFunctions(lines, codeStructure)

	// Обрабатываем переменные верхнего уровня
	p.parseVariables(lines, codeStructure)

	return codeStructure, nil
}

// parseImports извлекает импорты из Python файла
func (p *PythonParser) parseImports(lines []string, cs *models.CodeStructure) {
	for lineNum, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// Пропускаем пустые строки и комментарии
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "#") {
			continue
		}

		// Обрабатываем простые импорты (import module)
		if strings.HasPrefix(trimmedLine, "import ") {
			importNames := strings.Split(strings.TrimPrefix(trimmedLine, "import "), ",")
			for _, importName := range importNames {
				importName = strings.TrimSpace(importName)
				parts := strings.Fields(importName)
				
				// import module
				if len(parts) == 1 {
					cs.AddImport(&models.Import{
						Path: parts[0],
						Position: models.Position{
							StartLine: lineNum + 1,
							EndLine:   lineNum + 1,
						},
					})
				} else if len(parts) >= 3 && parts[1] == "as" {
					// import module as alias
					cs.AddImport(&models.Import{
						Path:  parts[0],
						Alias: parts[2],
						Position: models.Position{
							StartLine: lineNum + 1,
							EndLine:   lineNum + 1,
						},
					})
				}
			}
		}

		// Обрабатываем импорты from...import
		if strings.HasPrefix(trimmedLine, "from ") {
			fromImportPattern := regexp.MustCompile(`from\s+(\S+)\s+import\s+(.+)`)
			matches := fromImportPattern.FindStringSubmatch(trimmedLine)
			
			if len(matches) >= 3 {
				modulePath := matches[1]
				importItems := strings.Split(matches[2], ",")
				
				for _, item := range importItems {
					item = strings.TrimSpace(item)
					parts := strings.Fields(item)
					
					if len(parts) == 1 {
						// from module import item
						cs.AddImport(&models.Import{
							Path:     modulePath + "." + parts[0],
							Position: models.Position{
								StartLine: lineNum + 1,
								EndLine:   lineNum + 1,
							},
						})
					} else if len(parts) >= 3 && parts[1] == "as" {
						// from module import item as alias
						cs.AddImport(&models.Import{
							Path:     modulePath + "." + parts[0],
							Alias:    parts[2],
							Position: models.Position{
								StartLine: lineNum + 1,
								EndLine:   lineNum + 1,
							},
						})
					}
				}
			}
		}
	}
}

// parseClassesAndMethods извлекает классы и их методы из Python файла
func (p *PythonParser) parseClassesAndMethods(lines []string, cs *models.CodeStructure) {
	classPattern := regexp.MustCompile(`^class\s+(\w+)(?:\s*\(([\w\s,]+)\))?:`)
	methodPattern := regexp.MustCompile(`^\s+def\s+(\w+)\s*\(([\w\s,=.*:]*)\)(?:\s*->\s*(\w+))?:`)
	
	currentClassName := ""
	currentClassIndentation := 0
	inClass := false
	
	for lineNum, line := range lines {
		if len(strings.TrimSpace(line)) == 0 || strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}
		
		// Определение уровня отступа
		indentation := len(line) - len(strings.TrimLeft(line, " \t"))
		
		// Если мы в классе и текущий отступ меньше отступа класса, выходим из класса
		if inClass && indentation <= currentClassIndentation {
			inClass = false
			currentClassName = ""
		}
		
		// Проверяем, является ли строка определением класса
		classMatches := classPattern.FindStringSubmatch(line)
		if len(classMatches) > 0 {
			className := classMatches[1]
			parentClasses := ""
			
			if len(classMatches) > 2 && classMatches[2] != "" {
				parentClasses = classMatches[2]
			}
			
			// Проверяем, публичный ли класс (в Python обычно считается, что класс публичный)
			isPublic := !strings.HasPrefix(className, "_")
			
			cs.AddType(&models.Type{
				Name:     className,
				Kind:     "class",
				Parent:   parentClasses,
				IsPublic: isPublic,
				Position: models.Position{
					StartLine: lineNum + 1,
					EndLine:   lineNum + 1,
				},
			})
			
			// Добавляем класс в экспорты, если он публичный
			if isPublic {
				cs.AddExport(&models.Export{
					Name: className,
					Type: "class",
					Position: models.Position{
						StartLine: lineNum + 1,
						EndLine:   lineNum + 1,
					},
				})
			}
			
			inClass = true
			currentClassName = className
			currentClassIndentation = indentation
			continue
		}
		
		// Проверяем, является ли строка определением метода
		methodMatches := methodPattern.FindStringSubmatch(line)
		if inClass && len(methodMatches) > 0 {
			methodName := methodMatches[1]
			paramString := ""
			returnType := ""
			
			if len(methodMatches) > 2 {
				paramString = methodMatches[2]
			}
			
			if len(methodMatches) > 3 && methodMatches[3] != "" {
				returnType = methodMatches[3]
			}
			
			// Проверяем, публичный ли метод
			isPublic := !strings.HasPrefix(methodName, "_") || methodName == "__init__"
			
			method := &models.Method{
				Name:      methodName,
				IsPublic:  isPublic,
				BelongsTo: currentClassName,
				ReturnType: returnType,
				Position: models.Position{
					StartLine: lineNum + 1,
					EndLine:   lineNum + 1,
				},
			}
			
			// Разбираем параметры
			method.Parameters = p.parseParameters(paramString, methodName == "__init__")
			
			cs.AddMethod(method)
			
			// Добавляем метод в экспорты, если он публичный
			if isPublic && currentClassName != "" {
				cs.AddExport(&models.Export{
					Name: currentClassName + "." + methodName,
					Type: "method",
					Position: models.Position{
						StartLine: lineNum + 1,
						EndLine:   lineNum + 1,
					},
				})
			}
		}
	}
}

// parseFunctions извлекает функции верхнего уровня из Python файла
func (p *PythonParser) parseFunctions(lines []string, cs *models.CodeStructure) {
	functionPattern := regexp.MustCompile(`^def\s+(\w+)\s*\(([\w\s,=.*:]*)\)(?:\s*->\s*(\w+))?:`)
	
	for lineNum, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		
		// Пропускаем пустые строки и комментарии
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "#") {
			continue
		}
		
		// Проверяем, является ли строка определением функции
		funcMatches := functionPattern.FindStringSubmatch(trimmedLine)
		if len(funcMatches) > 0 {
			funcName := funcMatches[1]
			paramString := ""
			returnType := ""
			
			if len(funcMatches) > 2 {
				paramString = funcMatches[2]
			}
			
			if len(funcMatches) > 3 && funcMatches[3] != "" {
				returnType = funcMatches[3]
			}
			
			// Проверяем, публичная ли функция
			isPublic := !strings.HasPrefix(funcName, "_")
			
			method := &models.Method{
				Name:       funcName,
				IsPublic:   isPublic,
				ReturnType: returnType,
				Position: models.Position{
					StartLine: lineNum + 1,
					EndLine:   lineNum + 1,
				},
			}
			
			// Разбираем параметры
			method.Parameters = p.parseParameters(paramString, false)
			
			cs.AddMethod(method)
			
			// Добавляем функцию в экспорты, если она публичная
			if isPublic {
				cs.AddExport(&models.Export{
					Name: funcName,
					Type: "function",
					Position: models.Position{
						StartLine: lineNum + 1,
						EndLine:   lineNum + 1,
					},
				})
			}
		}
	}
}

// parseVariables извлекает переменные верхнего уровня из Python файла
func (p *PythonParser) parseVariables(lines []string, cs *models.CodeStructure) {
	variablePattern := regexp.MustCompile(`^(\w+)\s*=\s*(.+)$`)
	
	for lineNum, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		
		// Пропускаем пустые строки, комментарии и строки с импортами/классами/функциями
		if trimmedLine == "" || 
		   strings.HasPrefix(trimmedLine, "#") || 
		   strings.HasPrefix(trimmedLine, "import ") || 
		   strings.HasPrefix(trimmedLine, "from ") || 
		   strings.HasPrefix(trimmedLine, "class ") || 
		   strings.HasPrefix(trimmedLine, "def ") {
			continue
		}
		
		// Проверяем, является ли строка определением переменной
		varMatches := variablePattern.FindStringSubmatch(trimmedLine)
		if len(varMatches) > 0 {
			varName := varMatches[1]
			varValue := varMatches[2]
			
			// Проверяем, публичная ли переменная
			isPublic := !strings.HasPrefix(varName, "_")
			
			// Определяем, константа это или переменная (в Python нет встроенного понятия константы,
			// но по соглашению константы пишутся в UPPER_CASE)
			isConstant := strings.ToUpper(varName) == varName
			
			if isConstant {
				cs.AddConstant(&models.Constant{
					Name:  varName,
					Value: varValue,
					Position: models.Position{
						StartLine: lineNum + 1,
						EndLine:   lineNum + 1,
					},
				})
				
				if isPublic {
					cs.AddExport(&models.Export{
						Name: varName,
						Type: "constant",
						Position: models.Position{
							StartLine: lineNum + 1,
							EndLine:   lineNum + 1,
						},
					})
				}
			} else {
				cs.AddVariable(&models.Variable{
					Name:     varName,
					IsPublic: isPublic,
					Position: models.Position{
						StartLine: lineNum + 1,
						EndLine:   lineNum + 1,
					},
				})
				
				if isPublic {
					cs.AddExport(&models.Export{
						Name: varName,
						Type: "variable",
						Position: models.Position{
							StartLine: lineNum + 1,
							EndLine:   lineNum + 1,
						},
					})
				}
			}
		}
	}
}

// parseParameters разбирает строку параметров Python и создает массив Parameter
func (p *PythonParser) parseParameters(paramString string, isInitMethod bool) []*models.Parameter {
	if strings.TrimSpace(paramString) == "" {
		return []*models.Parameter{}
	}
	
	params := make([]*models.Parameter, 0)
	
	// Обрабатываем специальный случай с self/cls для методов
	if isInitMethod && strings.HasPrefix(strings.TrimSpace(paramString), "self") {
		paramString = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(paramString), "self"))
		if strings.HasPrefix(paramString, ",") {
			paramString = strings.TrimSpace(paramString[1:])
		}
	}
	
	// Если параметров нет - возвращаем пустой массив
	if strings.TrimSpace(paramString) == "" {
		return params
	}
	
	// Разбиваем параметры по запятой
	paramParts := strings.Split(paramString, ",")
	
	for _, part := range paramParts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		
		// Проверяем наличие значения по умолчанию
		isRequired := true
		defaultValue := ""
		
		if strings.Contains(part, "=") {
			parts := strings.SplitN(part, "=", 2)
			part = strings.TrimSpace(parts[0])
			defaultValue = strings.TrimSpace(parts[1])
			isRequired = false
		}
		
		// Проверяем наличие типа
		paramType := ""
		paramName := part
		
		if strings.Contains(part, ":") {
			typeParts := strings.SplitN(part, ":", 2)
			paramName = strings.TrimSpace(typeParts[0])
			paramType = strings.TrimSpace(typeParts[1])
		}
		
		// Проверяем специальные случаи (*args, **kwargs)
		if strings.HasPrefix(paramName, "*") {
			isRequired = false
		}
		
		param := &models.Parameter{
			Name:         paramName,
			Type:         paramType,
			IsRequired:   isRequired,
			DefaultValue: defaultValue,
		}
		
		params = append(params, param)
	}
	
	return params
}
