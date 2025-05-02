package models

// ConvertToFileStructure преобразует CodeStructure в FileStructure
// для совместимости с модулем генерации Markdown
func ConvertToFileStructure(cs *CodeStructure) FileStructure {
	// Создаем базовую структуру
	fs := FileStructure{
		Path:     cs.Metadata.Path,
		Language: cs.Metadata.LanguageName(),
		Content:  "", // Содержимое файла в FileStructure не используется
	}

	// Преобразуем импорты
	imports := make([]string, 0, len(cs.Imports))
	for _, imp := range cs.Imports {
		if imp.Alias != "" {
			imports = append(imports, imp.Alias+" "+imp.Path)
		} else {
			imports = append(imports, imp.Path)
		}
	}
	fs.Imports = imports

	// Преобразуем экспорты
	exports := make([]string, 0, len(cs.Exports))
	for _, exp := range cs.Exports {
		exportStr := exp.Name
		if exp.Type != "" {
			exportStr += " (" + exp.Type + ")"
		}
		exports = append(exports, exportStr)
	}
	fs.Exports = exports

	// Преобразуем методы
	methods := make([]MethodInfo, 0, len(cs.Methods))
	for _, method := range cs.Methods {
		if !method.IsPublic {
			continue // Пропускаем непубличные методы
		}

		// Формируем параметры
		params := make([]string, 0, len(method.Parameters))
		for _, param := range method.Parameters {
			paramStr := param.Name
			if param.Type != "" {
				paramStr += ": " + param.Type
			}
			params = append(params, paramStr)
		}

		// Формируем сигнатуру
		signature := method.Name + "("
		if len(method.Parameters) > 0 {
			paramStrs := make([]string, 0, len(method.Parameters))
			for _, param := range method.Parameters {
				paramStr := param.Name
				if param.Type != "" {
					paramStr += " " + param.Type
				}
				paramStrs = append(paramStrs, paramStr)
			}
			signature += joinStrings(paramStrs, ", ")
		}
		signature += ")"
		if method.ReturnType != "" {
			signature += " " + method.ReturnType
		}

		// Создаем MethodInfo
		methodInfo := MethodInfo{
			Name:      method.Name,
			Signature: signature,
			Params:    params,
			Returns:   []string{method.ReturnType},
		}

		// Если метод имеет описание, добавляем его
		if method.Description != "" {
			methodInfo.Body = method.Description
		} else {
			methodInfo.Body = "Нет описания"
		}

		methods = append(methods, methodInfo)
	}
	fs.Methods = methods

	// Группируем классы/типы
	classes := make([]string, 0, len(cs.Types))
	for _, typ := range cs.Types {
		if typ.IsPublic {
			classes = append(classes, typ.Name)
		}
	}
	fs.Classes = classes

	return fs
}

// joinStrings объединяет строки с указанным разделителем
func joinStrings(strings []string, separator string) string {
	if len(strings) == 0 {
		return ""
	}

	result := strings[0]
	for i := 1; i < len(strings); i++ {
		result += separator + strings[i]
	}
	return result
}
