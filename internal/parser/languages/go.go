package languages

import (
	"os"
	"strings"

	"code-telescope/internal/config"
	"code-telescope/internal/parser"
	"code-telescope/pkg/models"
)

// GoParser реализует парсер для языка Go
type GoParser struct {
	*parser.BaseParser
}

// NewGoParser создает новый экземпляр парсера Go
func NewGoParser(cfg *config.Config) *GoParser {
	return &GoParser{
		BaseParser: parser.NewBaseParser(cfg),
	}
}

// GetLanguageName возвращает название языка программирования
func (p *GoParser) GetLanguageName() string {
	return "Go"
}

// GetSupportedExtensions возвращает список поддерживаемых расширений файлов
func (p *GoParser) GetSupportedExtensions() []string {
	return []string{".go"}
}

// Parse разбирает файл Go и извлекает его структуру
func (p *GoParser) Parse(fileMetadata *models.FileMetadata) (*models.CodeStructure, error) {
	// Создаем пустую структуру кода
	codeStructure := models.NewCodeStructure(fileMetadata)

	// Читаем содержимое файла
	content, err := os.ReadFile(fileMetadata.AbsolutePath)
	if err != nil {
		return nil, err
	}

	// TODO: Использовать Tree-sitter для парсинга Go файла
	// В настоящее время используем простой подход с поиском маркеров

	// Разбиваем содержимое на строки
	lines := strings.Split(string(content), "\n")

	// Обрабатываем импорты
	p.parseImports(lines, codeStructure)

	// Обрабатываем типы и методы
	p.parseTypesAndMethods(lines, codeStructure)

	// Обрабатываем функции верхнего уровня
	p.parseFunctions(lines, codeStructure)

	// Обрабатываем переменные и константы
	p.parseVariablesAndConstants(lines, codeStructure)

	return codeStructure, nil
}

// parseImports извлекает импорты из Go файла
func (p *GoParser) parseImports(lines []string, cs *models.CodeStructure) {
	inImportBlock := false

	for lineNum, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// Находим начало блока импортов
		if strings.HasPrefix(trimmedLine, "import (") {
			inImportBlock = true
			continue
		}

		// Находим конец блока импортов
		if inImportBlock && trimmedLine == ")" {
			inImportBlock = false
			continue
		}

		// Обрабатываем однострочный импорт
		if strings.HasPrefix(trimmedLine, "import ") && !inImportBlock {
			importPath := strings.Trim(strings.TrimPrefix(trimmedLine, "import "), "\"")
			cs.AddImport(&models.Import{
				Path: importPath,
				Position: models.Position{
					StartLine: lineNum + 1,
					EndLine:   lineNum + 1,
				},
			})
			continue
		}

		// Обрабатываем импорты внутри блока
		if inImportBlock && trimmedLine != "" && !strings.HasPrefix(trimmedLine, "//") {
			// Проверяем наличие псевдонима
			importParts := strings.Fields(trimmedLine)

			var importPath string
			var alias string

			if len(importParts) > 1 {
				// Импорт с псевдонимом (например: alias "path/to/package")
				alias = importParts[0]
				importPath = strings.Trim(importParts[1], "\"")
			} else {
				// Обычный импорт (например: "path/to/package")
				importPath = strings.Trim(trimmedLine, "\"")
			}

			cs.AddImport(&models.Import{
				Path:  importPath,
				Alias: alias,
				Position: models.Position{
					StartLine: lineNum + 1,
					EndLine:   lineNum + 1,
				},
			})
		}
	}
}

// parseTypesAndMethods извлекает типы и методы из Go файла
func (p *GoParser) parseTypesAndMethods(lines []string, cs *models.CodeStructure) {
	// Простая реализация для демонстрации
	for lineNum, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// Поиск определений типов (struct, interface)
		if strings.HasPrefix(trimmedLine, "type ") && (strings.Contains(trimmedLine, " struct ") || strings.Contains(trimmedLine, " interface ")) {
			parts := strings.Fields(trimmedLine)
			if len(parts) >= 3 {
				typeName := parts[1]
				typeKind := parts[2]

				// Проверяем, публичный ли тип (начинается с заглавной буквы)
				isPublic := strings.ToUpper(typeName[:1]) == typeName[:1]

				cs.AddType(&models.Type{
					Name:     typeName,
					Kind:     typeKind,
					IsPublic: isPublic,
					Position: models.Position{
						StartLine: lineNum + 1,
						EndLine:   lineNum + 1,
					},
				})

				// Добавляем тип в экспорты, если он публичный
				if isPublic {
					cs.AddExport(&models.Export{
						Name: typeName,
						Type: "type",
						Position: models.Position{
							StartLine: lineNum + 1,
							EndLine:   lineNum + 1,
						},
					})
				}
			}
		}

		// Поиск методов
		if strings.HasPrefix(trimmedLine, "func (") {
			// Пример: func (p *Parser) Parse(...)
			openBracketPos := strings.Index(trimmedLine, "(")
			closeBracketPos := strings.Index(trimmedLine, ")")

			if openBracketPos != -1 && closeBracketPos != -1 && closeBracketPos > openBracketPos {
				receiverStr := trimmedLine[openBracketPos+1 : closeBracketPos]
				receiverParts := strings.Fields(receiverStr)

				if len(receiverParts) >= 2 {
					receiverType := strings.TrimPrefix(receiverParts[1], "*")

					// Получаем имя метода
					methodNameStart := strings.Index(trimmedLine[closeBracketPos:], " ") + closeBracketPos + 1
					methodNameEnd := strings.Index(trimmedLine[methodNameStart:], "(") + methodNameStart

					if methodNameEnd > methodNameStart {
						methodName := trimmedLine[methodNameStart:methodNameEnd]

						// Проверяем, публичный ли метод (начинается с заглавной буквы)
						isPublic := strings.ToUpper(methodName[:1]) == methodName[:1]

						method := &models.Method{
							Name:      methodName,
							IsPublic:  isPublic,
							BelongsTo: receiverType,
							Position: models.Position{
								StartLine: lineNum + 1,
								EndLine:   lineNum + 1,
							},
						}

						cs.AddMethod(method)

						// Добавляем метод в экспорты, если он публичный
						if isPublic {
							cs.AddExport(&models.Export{
								Name: methodName,
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
		}
	}
}

// parseFunctions извлекает функции верхнего уровня из Go файла
func (p *GoParser) parseFunctions(lines []string, cs *models.CodeStructure) {
	for lineNum, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// Поиск функций верхнего уровня (не методов)
		if strings.HasPrefix(trimmedLine, "func ") && !strings.HasPrefix(trimmedLine, "func (") {
			// Пример: func ParseFile(...) или func main()
			funcNameStart := len("func ")
			funcNameEnd := strings.Index(trimmedLine[funcNameStart:], "(") + funcNameStart

			if funcNameEnd > funcNameStart {
				funcName := trimmedLine[funcNameStart:funcNameEnd]

				// Проверяем, публичная ли функция (начинается с заглавной буквы)
				isPublic := strings.ToUpper(funcName[:1]) == funcName[:1]

				method := &models.Method{
					Name:     funcName,
					IsPublic: isPublic,
					Position: models.Position{
						StartLine: lineNum + 1,
						EndLine:   lineNum + 1,
					},
				}

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
}

// parseVariablesAndConstants извлекает переменные и константы из Go файла
func (p *GoParser) parseVariablesAndConstants(lines []string, cs *models.CodeStructure) {
	inVarBlock := false
	inConstBlock := false

	for lineNum, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// Обработка блоков переменных
		if strings.HasPrefix(trimmedLine, "var (") {
			inVarBlock = true
			continue
		}

		if inVarBlock && trimmedLine == ")" {
			inVarBlock = false
			continue
		}

		// Обработка блоков констант
		if strings.HasPrefix(trimmedLine, "const (") {
			inConstBlock = true
			continue
		}

		if inConstBlock && trimmedLine == ")" {
			inConstBlock = false
			continue
		}

		// Обработка однострочных объявлений переменных
		if strings.HasPrefix(trimmedLine, "var ") && !inVarBlock {
			parts := strings.Fields(trimmedLine[4:])
			if len(parts) >= 1 {
				varName := parts[0]
				varType := ""

				if len(parts) >= 3 && parts[1] != "=" {
					varType = parts[1]
				}

				// Проверяем, публичная ли переменная
				isPublic := strings.ToUpper(varName[:1]) == varName[:1]

				cs.AddVariable(&models.Variable{
					Name:     varName,
					Type:     varType,
					IsPublic: isPublic,
					Position: models.Position{
						StartLine: lineNum + 1,
						EndLine:   lineNum + 1,
					},
				})

				// Добавляем переменную в экспорты, если она публичная
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

		// Обработка однострочных объявлений констант
		if strings.HasPrefix(trimmedLine, "const ") && !inConstBlock {
			parts := strings.Fields(trimmedLine[6:])
			if len(parts) >= 1 {
				constName := parts[0]
				constType := ""
				constValue := ""

				// Попытка извлечь тип и значение
				if len(parts) >= 3 {
					if parts[1] == "=" {
						constValue = strings.Join(parts[2:], " ")
					} else {
						constType = parts[1]
						if len(parts) >= 4 && parts[2] == "=" {
							constValue = strings.Join(parts[3:], " ")
						}
					}
				}

				// Проверяем, публичная ли константа
				isPublic := strings.ToUpper(constName[:1]) == constName[:1]

				cs.AddConstant(&models.Constant{
					Name:  constName,
					Type:  constType,
					Value: constValue,
					Position: models.Position{
						StartLine: lineNum + 1,
						EndLine:   lineNum + 1,
					},
				})

				// Добавляем константу в экспорты, если она публичная
				if isPublic {
					cs.AddExport(&models.Export{
						Name: constName,
						Type: "constant",
						Position: models.Position{
							StartLine: lineNum + 1,
							EndLine:   lineNum + 1,
						},
					})
				}
			}
		}

		// Обработка элементов внутри блоков
		if (inVarBlock || inConstBlock) && trimmedLine != "" && !strings.HasPrefix(trimmedLine, "//") {
			parts := strings.Fields(trimmedLine)
			if len(parts) >= 1 {
				name := parts[0]
				typ := ""
				value := ""

				// Попытка извлечь тип и значение
				if len(parts) >= 3 {
					if parts[1] == "=" {
						value = strings.Join(parts[2:], " ")
					} else {
						typ = parts[1]
						if len(parts) >= 4 && parts[2] == "=" {
							value = strings.Join(parts[3:], " ")
						}
					}
				}

				// Проверяем, публичный ли элемент
				isPublic := strings.ToUpper(name[:1]) == name[:1]

				if inVarBlock {
					cs.AddVariable(&models.Variable{
						Name:     name,
						Type:     typ,
						IsPublic: isPublic,
						Position: models.Position{
							StartLine: lineNum + 1,
							EndLine:   lineNum + 1,
						},
					})

					// Добавляем переменную в экспорты, если она публичная
					if isPublic {
						cs.AddExport(&models.Export{
							Name: name,
							Type: "variable",
							Position: models.Position{
								StartLine: lineNum + 1,
								EndLine:   lineNum + 1,
							},
						})
					}
				} else if inConstBlock {
					cs.AddConstant(&models.Constant{
						Name:  name,
						Type:  typ,
						Value: value,
						Position: models.Position{
							StartLine: lineNum + 1,
							EndLine:   lineNum + 1,
						},
					})

					// Добавляем константу в экспорты, если она публичная
					if isPublic {
						cs.AddExport(&models.Export{
							Name: name,
							Type: "constant",
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
