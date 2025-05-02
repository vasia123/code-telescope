package languages

import (
	"os"
	"strings"

	"code-telescope/internal/config"
	"code-telescope/pkg/models"
	"code-telescope/pkg/utils"
)

// JavaScriptParser реализует парсер для языка JavaScript
type JavaScriptParser struct {
	Config *config.Config
}

// NewJavaScriptParser создает новый экземпляр парсера JavaScript
func NewJavaScriptParser(cfg *config.Config) *JavaScriptParser {
	return &JavaScriptParser{
		Config: cfg,
	}
}

// GetLanguageName возвращает название языка программирования
func (p *JavaScriptParser) GetLanguageName() string {
	return "JavaScript"
}

// GetSupportedExtensions возвращает список поддерживаемых расширений файлов
func (p *JavaScriptParser) GetSupportedExtensions() []string {
	return []string{".js", ".jsx", ".mjs"}
}

// Parse разбирает файл JavaScript и извлекает его структуру
func (p *JavaScriptParser) Parse(fileMetadata *models.FileMetadata) (*models.CodeStructure, error) {
	// Создаем пустую структуру кода
	codeStructure := models.NewCodeStructure(fileMetadata)

	// Читаем содержимое файла
	content, err := os.ReadFile(fileMetadata.AbsolutePath)
	if err != nil {
		return nil, err
	}

	// Разбиваем содержимое на строки
	lines := strings.Split(string(content), "\n")

	// Обрабатываем импорты
	p.parseImports(lines, codeStructure)

	// Обрабатываем экспорты
	p.parseExports(lines, codeStructure)

	// Обрабатываем функции
	p.parseFunctions(lines, codeStructure)

	// Обрабатываем классы и методы
	p.parseClassesAndMethods(lines, codeStructure)

	return codeStructure, nil
}

// parseImports извлекает импорты из JavaScript файла
func (p *JavaScriptParser) parseImports(lines []string, cs *models.CodeStructure) {
	for lineNum, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// import defaultExport from 'module-name';
		// import * as name from 'module-name';
		// import { export1, export2 } from 'module-name';
		if strings.HasPrefix(trimmedLine, "import ") {
			pathStartPos := strings.Index(trimmedLine, "from ")
			if pathStartPos != -1 {
				pathStartPos += 5 // длина "from "
				// Ищем путь в кавычках
				if pathStartPos < len(trimmedLine) {
					path := ""
					if strings.Contains(trimmedLine[pathStartPos:], "'") {
						// 'module-name'
						startQuote := strings.Index(trimmedLine[pathStartPos:], "'") + pathStartPos
						endQuote := strings.LastIndex(trimmedLine, "'")
						if startQuote != -1 && endQuote > startQuote {
							path = trimmedLine[startQuote+1 : endQuote]
						}
					} else if strings.Contains(trimmedLine[pathStartPos:], "\"") {
						// "module-name"
						startQuote := strings.Index(trimmedLine[pathStartPos:], "\"") + pathStartPos
						endQuote := strings.LastIndex(trimmedLine, "\"")
						if startQuote != -1 && endQuote > startQuote {
							path = trimmedLine[startQuote+1 : endQuote]
						}
					}

					if path != "" {
						// Извлекаем имя импорта
						alias := ""
						if strings.Contains(trimmedLine, "* as ") {
							// import * as name
							asPos := strings.Index(trimmedLine, "* as ") + 5
							spacePos := strings.Index(trimmedLine[asPos:], " ") + asPos
							if spacePos > asPos {
								alias = trimmedLine[asPos:spacePos]
							}
						} else if strings.HasPrefix(trimmedLine, "import ") && !strings.Contains(trimmedLine, "{") {
							// import defaultExport
							importKeywordLen := len("import ")
							spacePos := strings.Index(trimmedLine[importKeywordLen:], " ") + importKeywordLen
							if spacePos > importKeywordLen {
								alias = trimmedLine[importKeywordLen:spacePos]
							}
						}

						cs.AddImport(&models.Import{
							Path:  path,
							Alias: alias,
							Position: models.Position{
								StartLine: lineNum + 1,
								EndLine:   lineNum + 1,
							},
						})
					}
				}
			}
		}

		// require() синтаксис
		if strings.Contains(trimmedLine, "require(") {
			startPos := strings.Index(trimmedLine, "require(") + len("require(")
			endPos := strings.Index(trimmedLine[startPos:], ")") + startPos
			if startPos != -1 && endPos > startPos {
				path := trimmedLine[startPos:endPos]
				// Удаляем кавычки
				path = strings.Trim(path, "'\"")

				cs.AddImport(&models.Import{
					Path: path,
					Position: models.Position{
						StartLine: lineNum + 1,
						EndLine:   lineNum + 1,
					},
				})
			}
		}
	}
}

// parseExports извлекает экспорты из JavaScript файла
func (p *JavaScriptParser) parseExports(lines []string, cs *models.CodeStructure) {
	for lineNum, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// export default expression;
		if strings.HasPrefix(trimmedLine, "export default ") {
			name := "default"
			exportType := "default"

			// Пытаемся определить тип экспорта
			if strings.Contains(trimmedLine, "function ") {
				exportType = "function"
				// Извлекаем имя функции, если есть
				funcPos := strings.Index(trimmedLine, "function ") + len("function ")
				if funcPos < len(trimmedLine) {
					spacePos := strings.Index(trimmedLine[funcPos:], "(")
					if spacePos != -1 {
						name = strings.TrimSpace(trimmedLine[funcPos : funcPos+spacePos])
					}
				}
			} else if strings.Contains(trimmedLine, "class ") {
				exportType = "class"
				// Извлекаем имя класса
				classPos := strings.Index(trimmedLine, "class ") + len("class ")
				if classPos < len(trimmedLine) {
					spacePos := strings.Index(trimmedLine[classPos:], " ")
					if spacePos != -1 {
						name = strings.TrimSpace(trimmedLine[classPos : classPos+spacePos])
					}
				}
			}

			cs.AddExport(&models.Export{
				Name: name,
				Type: exportType,
				Position: models.Position{
					StartLine: lineNum + 1,
					EndLine:   lineNum + 1,
				},
			})
		}

		// export { name1, name2, …, nameN };
		// export { variable1 as name1, variable2 as name2, …, nameN };
		if strings.HasPrefix(trimmedLine, "export {") {
			startBrace := strings.Index(trimmedLine, "{")
			endBrace := strings.Index(trimmedLine, "}")
			if startBrace != -1 && endBrace > startBrace {
				exports := trimmedLine[startBrace+1 : endBrace]
				exportItems := strings.Split(exports, ",")
				for _, item := range exportItems {
					item = strings.TrimSpace(item)
					if item == "" {
						continue
					}

					name := item
					exportType := "variable"

					// Проверяем на 'as' синтаксис: variable as name
					if strings.Contains(item, " as ") {
						parts := strings.Split(item, " as ")
						name = strings.TrimSpace(parts[1])
					}

					cs.AddExport(&models.Export{
						Name: name,
						Type: exportType,
						Position: models.Position{
							StartLine: lineNum + 1,
							EndLine:   lineNum + 1,
						},
					})
				}
			}
		}

		// export const/let/var name = ...
		if strings.HasPrefix(trimmedLine, "export const ") ||
			strings.HasPrefix(trimmedLine, "export let ") ||
			strings.HasPrefix(trimmedLine, "export var ") {

			typeEndPos := len("export ")
			if strings.HasPrefix(trimmedLine, "export const ") {
				typeEndPos = len("export const ")
			} else if strings.HasPrefix(trimmedLine, "export let ") {
				typeEndPos = len("export let ")
			} else if strings.HasPrefix(trimmedLine, "export var ") {
				typeEndPos = len("export var ")
			}

			// Извлекаем имя переменной
			nameEndPos := strings.Index(trimmedLine[typeEndPos:], "=")
			if nameEndPos != -1 {
				name := strings.TrimSpace(trimmedLine[typeEndPos : typeEndPos+nameEndPos])
				cs.AddExport(&models.Export{
					Name: name,
					Type: "variable",
					Position: models.Position{
						StartLine: lineNum + 1,
						EndLine:   lineNum + 1,
					},
				})
			}
		}

		// export function name() { ... }
		if strings.HasPrefix(trimmedLine, "export function ") {
			funcPos := len("export function ")
			// Извлекаем имя функции
			parenPos := strings.Index(trimmedLine[funcPos:], "(")
			if parenPos != -1 {
				name := strings.TrimSpace(trimmedLine[funcPos : funcPos+parenPos])
				cs.AddExport(&models.Export{
					Name: name,
					Type: "function",
					Position: models.Position{
						StartLine: lineNum + 1,
						EndLine:   lineNum + 1,
					},
				})
			}
		}

		// export class Name { ... }
		if strings.HasPrefix(trimmedLine, "export class ") {
			classPos := len("export class ")
			// Извлекаем имя класса
			spacePos := strings.Index(trimmedLine[classPos:], " ")
			bracePos := strings.Index(trimmedLine[classPos:], "{")
			endPos := spacePos
			if (bracePos != -1 && bracePos < spacePos) || spacePos == -1 {
				endPos = bracePos
			}
			if endPos != -1 {
				name := strings.TrimSpace(trimmedLine[classPos : classPos+endPos])
				cs.AddExport(&models.Export{
					Name: name,
					Type: "class",
					Position: models.Position{
						StartLine: lineNum + 1,
						EndLine:   lineNum + 1,
					},
				})
			}
		}
	}
}

// parseFunctions извлекает функции из JavaScript файла
func (p *JavaScriptParser) parseFunctions(lines []string, cs *models.CodeStructure) {
	for lineNum, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// Функции: function name() { ... }
		if (strings.HasPrefix(trimmedLine, "function ") || strings.HasPrefix(trimmedLine, "const ") ||
			strings.HasPrefix(trimmedLine, "let ") || strings.HasPrefix(trimmedLine, "var ")) &&
			strings.Contains(trimmedLine, "(") {

			var funcName string
			var isPublic bool
			var funcStart int

			if strings.HasPrefix(trimmedLine, "function ") {
				// Классическое определение функции
				funcStart = len("function ")
				funcEnd := strings.Index(trimmedLine[funcStart:], "(") + funcStart
				if funcEnd > funcStart {
					funcName = strings.TrimSpace(trimmedLine[funcStart:funcEnd])
					isPublic = true // В JavaScript все функции публичные, если они экспортированы
				}
			} else if strings.Contains(trimmedLine, " = function(") ||
				strings.Contains(trimmedLine, " = (") ||
				strings.Contains(trimmedLine, " => ") {
				// Функциональные выражения или стрелочные функции
				// const/let/var name = function() { ... } или const/let/var name = () => { ... }
				if strings.HasPrefix(trimmedLine, "const ") {
					funcStart = len("const ")
				} else if strings.HasPrefix(trimmedLine, "let ") {
					funcStart = len("let ")
				} else if strings.HasPrefix(trimmedLine, "var ") {
					funcStart = len("var ")
				}

				funcEnd := strings.Index(trimmedLine[funcStart:], " =") + funcStart
				if funcEnd > funcStart {
					funcName = strings.TrimSpace(trimmedLine[funcStart:funcEnd])
					isPublic = true
				}
			}

			if funcName != "" {
				method := &models.Method{
					Name:     funcName,
					IsPublic: isPublic,
					Position: models.Position{
						StartLine: lineNum + 1,
						EndLine:   lineNum + 1,
					},
				}

				// Извлекаем параметры
				paramStartPos := strings.Index(trimmedLine, "(")
				if paramStartPos != -1 {
					paramEndPos := utils.FindMatchingCloseBracket(trimmedLine, paramStartPos)
					if paramEndPos > paramStartPos {
						paramString := trimmedLine[paramStartPos+1 : paramEndPos]
						parameters := p.parseJSParameters(paramString)
						method.Parameters = parameters
					}
				}

				cs.AddMethod(method)
			}
		}
	}
}

// parseClassesAndMethods извлекает классы и методы из JavaScript файла
func (p *JavaScriptParser) parseClassesAndMethods(lines []string, cs *models.CodeStructure) {
	inClass := false
	currentClass := ""

	for lineNum, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// Поиск определений класса
		if strings.HasPrefix(trimmedLine, "class ") {
			classNameStart := len("class ")
			classNameEnd := strings.Index(trimmedLine[classNameStart:], " ")
			bracePos := strings.Index(trimmedLine[classNameStart:], "{")

			if classNameEnd == -1 || (bracePos != -1 && bracePos < classNameEnd) {
				classNameEnd = bracePos
			}

			if classNameEnd != -1 {
				className := trimmedLine[classNameStart : classNameStart+classNameEnd]
				currentClass = className
				inClass = true

				// Создаем тип для класса
				typ := &models.Type{
					Name:     className,
					Kind:     "class",
					IsPublic: true, // Предполагаем, что класс публичный
					Position: models.Position{
						StartLine: lineNum + 1,
						EndLine:   lineNum + 1,
					},
				}

				cs.AddType(typ)
			}
		}

		// Закрывающая скобка класса
		if inClass && trimmedLine == "}" {
			inClass = false
			currentClass = ""
		}

		// Методы класса
		if inClass && (strings.Contains(trimmedLine, "(") || strings.Contains(trimmedLine, "=")) {
			// constructor() { ... }
			// methodName() { ... }
			// get property() { ... }
			// set property(value) { ... }
			// static methodName() { ... }

			isStatic := strings.HasPrefix(trimmedLine, "static ")
			isGetter := strings.HasPrefix(trimmedLine, "get ")
			isSetter := strings.HasPrefix(trimmedLine, "set ")

			// Определяем начальную позицию имени метода
			methodStart := 0
			if isStatic {
				methodStart = len("static ")
			} else if isGetter {
				methodStart = len("get ")
			} else if isSetter {
				methodStart = len("set ")
			}

			// Ищем конец имени метода (скобка параметров)
			methodEnd := strings.Index(trimmedLine[methodStart:], "(") + methodStart
			if methodEnd > methodStart {
				methodName := strings.TrimSpace(trimmedLine[methodStart:methodEnd])

				// Создаем метод
				method := &models.Method{
					Name:      methodName,
					IsPublic:  !strings.HasPrefix(methodName, "_"), // В JavaScript приватные методы часто начинаются с _
					IsStatic:  isStatic,
					BelongsTo: currentClass,
					Position: models.Position{
						StartLine: lineNum + 1,
						EndLine:   lineNum + 1,
					},
				}

				// Извлекаем параметры
				paramStartPos := methodEnd
				paramEndPos := utils.FindMatchingCloseBracket(trimmedLine, paramStartPos)
				if paramEndPos > paramStartPos {
					paramString := trimmedLine[paramStartPos+1 : paramEndPos]
					parameters := p.parseJSParameters(paramString)
					method.Parameters = parameters
				}

				cs.AddMethod(method)
			}
		}
	}
}

// parseJSParameters разбирает строку параметров JS и создает массив Parameter
func (p *JavaScriptParser) parseJSParameters(paramString string) []*models.Parameter {
	if strings.TrimSpace(paramString) == "" {
		return []*models.Parameter{}
	}

	params := make([]*models.Parameter, 0)
	paramParts := strings.Split(paramString, ",")

	for _, part := range paramParts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Проверяем на параметры с значениями по умолчанию: name = value
		defaultValue := ""
		isRequired := true
		name := part

		if strings.Contains(part, "=") {
			eqPos := strings.Index(part, "=")
			name = strings.TrimSpace(part[:eqPos])
			defaultValue = strings.TrimSpace(part[eqPos+1:])
			isRequired = false
		}

		// Проверяем на деструктуризацию и рест-параметры
		if strings.HasPrefix(name, "{") || strings.HasPrefix(name, "[") || strings.HasPrefix(name, "...") {
			// Упрощенная обработка для сложных случаев
			param := &models.Parameter{
				Name:         name,
				Type:         "any", // JavaScript не имеет статической типизации
				DefaultValue: defaultValue,
				IsRequired:   isRequired,
			}
			params = append(params, param)
		} else {
			// Обычный параметр
			param := &models.Parameter{
				Name:         name,
				Type:         "any", // JavaScript не имеет статической типизации
				DefaultValue: defaultValue,
				IsRequired:   isRequired,
			}
			params = append(params, param)
		}
	}

	return params
}
