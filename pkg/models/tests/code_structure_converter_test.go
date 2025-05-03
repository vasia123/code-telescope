package tests

import (
	"testing"

	"code-telescope/pkg/models"

	"github.com/stretchr/testify/assert"
)

// TestConvertToFileStructure проверяет корректность преобразования между CodeStructure и FileStructure
func TestConvertToFileStructure(t *testing.T) {
	// Создание тестового экземпляра FileMetadata
	metadata := &models.FileMetadata{
		Path:         "test/path/file.go",
		AbsolutePath: "/absolute/test/path/file.go",
		Name:         "file.go",
		Extension:    ".go",
		Directory:    "test/path",
	}

	// Создание тестового экземпляра CodeStructure
	codeStructure := models.NewCodeStructure(metadata)

	// Добавление импортов
	codeStructure.AddImport(&models.Import{
		Path:  "fmt",
		Alias: "",
	})
	codeStructure.AddImport(&models.Import{
		Path:  "github.com/example/custom-pkg",
		Alias: "custom",
	})

	// Добавление экспортов
	codeStructure.AddExport(&models.Export{
		Name: "PublicFunction",
		Type: "function",
	})
	codeStructure.AddExport(&models.Export{
		Name: "PublicType",
		Type: "type",
	})

	// Добавление типа
	publicType := &models.Type{
		Name:     "PublicType",
		Kind:     "struct",
		IsPublic: true,
		Properties: []*models.Property{
			{
				Name:     "PublicField",
				Type:     "string",
				IsPublic: true,
			},
			{
				Name:     "privateField",
				Type:     "int",
				IsPublic: false,
			},
		},
		Methods: []*models.Method{},
	}
	codeStructure.AddType(publicType)

	// Добавление методов
	publicMethod := &models.Method{
		Name:       "PublicMethod",
		IsPublic:   true,
		BelongsTo:  "PublicType",
		ReturnType: "string",
		Parameters: []*models.Parameter{
			{
				Name: "param",
				Type: "string",
			},
		},
		Description: "Публичный метод типа PublicType",
	}
	codeStructure.AddMethod(publicMethod)

	privateMethod := &models.Method{
		Name:        "privateMethod",
		IsPublic:    false,
		BelongsTo:   "PublicType",
		ReturnType:  "",
		Parameters:  []*models.Parameter{},
		Description: "Приватный метод типа PublicType",
	}
	codeStructure.AddMethod(privateMethod)

	// Добавление функции верхнего уровня
	publicFunction := &models.Method{
		Name:       "PublicFunction",
		IsPublic:   true,
		BelongsTo:  "",
		ReturnType: "string, error",
		Parameters: []*models.Parameter{
			{
				Name: "param1",
				Type: "string",
			},
			{
				Name: "param2",
				Type: "int",
			},
		},
		Description: "Публичная функция верхнего уровня",
	}
	codeStructure.AddMethod(publicFunction)

	// Конвертация в FileStructure
	fileStructure := models.ConvertToFileStructure(codeStructure)

	// Проверки
	assert.Equal(t, "test/path/file.go", fileStructure.Path, "Путь должен быть сохранен")
	assert.Equal(t, "Go", fileStructure.Language, "Язык должен быть определен корректно")

	// Проверка импортов
	assert.Equal(t, 2, len(fileStructure.Imports), "Должно быть 2 импорта")
	assert.Contains(t, fileStructure.Imports, "fmt")
	assert.Contains(t, fileStructure.Imports, "custom \"github.com/example/custom-pkg\"")

	// Проверка экспортов
	assert.Equal(t, 2, len(fileStructure.Exports), "Должно быть 2 экспорта")
	assert.Contains(t, fileStructure.Exports, "PublicFunction (function)")
	assert.Contains(t, fileStructure.Exports, "PublicType (type)")

	// Проверка методов
	assert.Equal(t, 3, len(fileStructure.Methods), "Должно быть 3 метода")

	// Находим метод типа и функцию по имени
	var foundPublicFunction, foundPublicMethod, foundPrivateMethod bool

	for _, method := range fileStructure.Methods {
		if method.Name == "PublicFunction" {
			foundPublicFunction = true
			assert.Equal(t, "func PublicFunction(param1 string, param2 int) (string, error)", method.Signature, "Сигнатура публичной функции должна быть корректной")
			assert.Equal(t, "Публичная функция верхнего уровня", method.Description, "Описание публичной функции должно быть сохранено")
			assert.Equal(t, []string{"param1 string", "param2 int"}, method.Params, "Параметры публичной функции должны быть корректными")
			assert.Equal(t, []string{"string", "error"}, method.Returns, "Возвращаемые значения публичной функции должны быть корректными")
		} else if method.Name == "PublicMethod" {
			foundPublicMethod = true
			assert.Equal(t, "func (PublicType) PublicMethod(param string) string", method.Signature, "Сигнатура публичного метода должна быть корректной")
			assert.Equal(t, "Публичный метод типа PublicType", method.Description, "Описание публичного метода должно быть сохранено")
		} else if method.Name == "privateMethod" {
			foundPrivateMethod = true
			assert.Equal(t, "func (PublicType) privateMethod()", method.Signature, "Сигнатура приватного метода должна быть корректной")
			assert.Equal(t, "Приватный метод типа PublicType", method.Description, "Описание приватного метода должно быть сохранено")
		}
	}

	assert.True(t, foundPublicFunction, "Публичная функция должна быть в FileStructure")
	assert.True(t, foundPublicMethod, "Публичный метод должен быть в FileStructure")
	assert.True(t, foundPrivateMethod, "Приватный метод должен быть в FileStructure")

	// Проверка классов
	assert.Equal(t, 1, len(fileStructure.Classes), "Должен быть 1 класс")
	assert.Contains(t, fileStructure.Classes, "PublicType")
}
