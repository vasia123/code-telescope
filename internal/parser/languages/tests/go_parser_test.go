package tests

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"code-telescope/internal/config"
	"code-telescope/internal/parser/languages"
	"code-telescope/pkg/models"

	"github.com/stretchr/testify/assert"
)

// createTempFile создает временный файл с заданным содержимым и расширением
func createTempFile(t *testing.T, content, extension string) (*os.File, string) {
	tmpfile, err := ioutil.TempFile("", "test-*"+extension)
	if err != nil {
		t.Fatalf("Не удалось создать временный файл: %v", err)
	}

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
		t.Fatalf("Не удалось записать во временный файл: %v", err)
	}

	// Закрыть файл и вернуть, чтобы его можно было открыть для чтения позже
	if err := tmpfile.Close(); err != nil {
		os.Remove(tmpfile.Name())
		t.Fatalf("Не удалось закрыть временный файл: %v", err)
	}

	return tmpfile, filepath.Base(tmpfile.Name())
}

// createFileMetadata создает объект метаданных файла на основе пути
func createFileMetadata(t *testing.T, filePath string) *models.FileMetadata {
	metadata, err := models.NewFileMetadata(filePath, filepath.Dir(filePath))
	if err != nil {
		t.Fatalf("Не удалось создать метаданные файла: %v", err)
	}
	return metadata
}

// TestGoParserFunctions проверяет корректность извлечения функций из Go файла
func TestGoParserFunctions(t *testing.T) {
	// Пример содержимого Go файла с функциями
	content := `package example

func PublicFunction(param1 string, param2 int) (string, error) {
    // тело функции
    return "", nil
}

func privateFunction() {
    // тело функции
}`

	// Создание временного файла
	tmpfile, _ := createTempFile(t, content, ".go")
	defer os.Remove(tmpfile.Name())

	// Создание конфигурации и парсера
	cfg := config.DefaultConfig()
	parser := languages.NewGoParser(cfg)

	// Парсинг файла
	metadata := createFileMetadata(t, tmpfile.Name())
	structure, err := parser.Parse(metadata)

	// Проверки
	assert.NoError(t, err, "Парсинг должен выполняться без ошибок")
	assert.NotNil(t, structure, "Структура кода не должна быть nil")
	assert.Equal(t, 2, len(structure.Methods), "Должно быть извлечено 2 функции")

	// Проверка публичной функции
	var publicFunc, privateFunc *models.Method
	for _, m := range structure.Methods {
		if m.Name == "PublicFunction" {
			publicFunc = m
		} else if m.Name == "privateFunction" {
			privateFunc = m
		}
	}

	assert.NotNil(t, publicFunc, "Публичная функция должна быть извлечена")
	assert.True(t, publicFunc.IsPublic, "Функция PublicFunction должна быть публичной")
	assert.Equal(t, 2, len(publicFunc.Parameters), "Публичная функция должна иметь 2 параметра")
	assert.Equal(t, "string, error", publicFunc.ReturnType, "Публичная функция должна возвращать (string, error)")

	// Проверка приватной функции
	assert.NotNil(t, privateFunc, "Приватная функция должна быть извлечена")
	assert.False(t, privateFunc.IsPublic, "Функция privateFunction должна быть приватной")
	assert.Equal(t, 0, len(privateFunc.Parameters), "Приватная функция не должна иметь параметров")
	assert.Equal(t, "", privateFunc.ReturnType, "Приватная функция не должна иметь возвращаемых значений")
}

// TestGoParserTypes проверяет корректность извлечения типов из Go файла
func TestGoParserTypes(t *testing.T) {
	// Пример содержимого Go файла с типами и методами
	content := `package example

type PublicType struct {
    PublicField string
    privateField int
}

func (p *PublicType) PublicMethod() string {
    return p.PublicField
}

func (p PublicType) privateMethod() {
    // тело метода
}`

	// Создание временного файла
	tmpfile, _ := createTempFile(t, content, ".go")
	defer os.Remove(tmpfile.Name())

	// Создание конфигурации и парсера
	cfg := config.DefaultConfig()
	parser := languages.NewGoParser(cfg)

	// Парсинг файла
	metadata := createFileMetadata(t, tmpfile.Name())
	structure, err := parser.Parse(metadata)

	// Проверки
	assert.NoError(t, err, "Парсинг должен выполняться без ошибок")
	assert.NotNil(t, structure, "Структура кода не должна быть nil")
	assert.Equal(t, 1, len(structure.Types), "Должен быть извлечен 1 тип")
	assert.Equal(t, 2, len(structure.Methods), "Должно быть извлечено 2 метода")

	// Проверка типа
	publicType := structure.Types[0]
	assert.Equal(t, "PublicType", publicType.Name, "Имя типа должно быть извлечено корректно")
	assert.True(t, publicType.IsPublic, "Тип PublicType должен быть публичным")
	assert.Equal(t, "struct", publicType.Kind, "Kind типа должен быть 'struct'")
	assert.Equal(t, 2, len(publicType.Properties), "Тип должен иметь 2 поля")

	// Проверка полей типа
	var publicField, privateField *models.Property
	for _, p := range publicType.Properties {
		if p.Name == "PublicField" {
			publicField = p
		} else if p.Name == "privateField" {
			privateField = p
		}
	}

	assert.NotNil(t, publicField, "Публичное поле должно быть извлечено")
	assert.True(t, publicField.IsPublic, "Поле PublicField должно быть публичным")
	assert.Equal(t, "string", publicField.Type, "Тип публичного поля должен быть string")

	assert.NotNil(t, privateField, "Приватное поле должно быть извлечено")
	assert.False(t, privateField.IsPublic, "Поле privateField должно быть приватным")
	assert.Equal(t, "int", privateField.Type, "Тип приватного поля должен быть int")

	// Проверка методов типа
	var publicMethod, privateMethod *models.Method
	for _, m := range structure.Methods {
		if m.Name == "PublicMethod" && m.BelongsTo == "PublicType" {
			publicMethod = m
		} else if m.Name == "privateMethod" && m.BelongsTo == "PublicType" {
			privateMethod = m
		}
	}

	assert.NotNil(t, publicMethod, "Публичный метод должен быть извлечен")
	assert.True(t, publicMethod.IsPublic, "Метод PublicMethod должен быть публичным")
	assert.Equal(t, "string", publicMethod.ReturnType, "Публичный метод должен возвращать string")

	assert.NotNil(t, privateMethod, "Приватный метод должен быть извлечен")
	assert.False(t, privateMethod.IsPublic, "Метод privateMethod должен быть приватным")
}

// TestGoParserImportsExports проверяет корректность извлечения импортов и экспортов
func TestGoParserImportsExports(t *testing.T) {
	// Пример содержимого Go файла с импортами
	content := `package example

import (
    "fmt"
    "io"
    custom "github.com/example/custom-pkg"
)

// Экспортируемые переменные и константы
const (
    PublicConst = "public"
    privateConst = "private"
)

var (
    PublicVar = 10
    privateVar = 20
)

// Экспортируемые типы и функции уже проверены в других тестах
`

	// Создание временного файла
	tmpfile, _ := createTempFile(t, content, ".go")
	defer os.Remove(tmpfile.Name())

	// Создание конфигурации и парсера
	cfg := config.DefaultConfig()
	parser := languages.NewGoParser(cfg)

	// Парсинг файла
	metadata := createFileMetadata(t, tmpfile.Name())
	structure, err := parser.Parse(metadata)

	// Проверки
	assert.NoError(t, err, "Парсинг должен выполняться без ошибок")
	assert.NotNil(t, structure, "Структура кода не должна быть nil")
	assert.Equal(t, 3, len(structure.Imports), "Должно быть извлечено 3 импорта")

	// Проверка импортов
	importPaths := map[string]string{}
	for _, imp := range structure.Imports {
		importPaths[imp.Path] = imp.Alias
	}

	assert.Contains(t, importPaths, "fmt")
	assert.Contains(t, importPaths, "io")
	assert.Contains(t, importPaths, "github.com/example/custom-pkg")
	assert.Equal(t, "custom", importPaths["github.com/example/custom-pkg"])

	// Проверка констант
	assert.Equal(t, 2, len(structure.Constants), "Должно быть извлечено 2 константы")
	var publicConst, privateConst *models.Constant
	for _, c := range structure.Constants {
		if c.Name == "PublicConst" {
			publicConst = c
		} else if c.Name == "privateConst" {
			privateConst = c
		}
	}

	assert.NotNil(t, publicConst, "Публичная константа должна быть извлечена")
	assert.Equal(t, "\"public\"", publicConst.Value, "Значение константы должно быть корректным")

	assert.NotNil(t, privateConst, "Приватная константа должна быть извлечена")
	assert.Equal(t, "\"private\"", privateConst.Value, "Значение константы должно быть корректным")

	// Проверка переменных
	assert.Equal(t, 2, len(structure.Variables), "Должно быть извлечено 2 переменные")
	var publicVar, privateVar *models.Variable
	for _, v := range structure.Variables {
		if v.Name == "PublicVar" {
			publicVar = v
		} else if v.Name == "privateVar" {
			privateVar = v
		}
	}

	assert.NotNil(t, publicVar, "Публичная переменная должна быть извлечена")
	assert.True(t, publicVar.IsPublic, "Переменная PublicVar должна быть публичной")

	assert.NotNil(t, privateVar, "Приватная переменная должна быть извлечена")
	assert.False(t, privateVar.IsPublic, "Переменная privateVar должна быть приватной")
}
