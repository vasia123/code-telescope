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

// TestGoParserInterface проверяет корректность извлечения интерфейсов из Go файла
func TestGoParserInterface(t *testing.T) {
	// Пример содержимого Go файла с интерфейсами
	content := `package example

// Repository определяет общий интерфейс для работы с хранилищем
type Repository interface {
	// Get возвращает сущность по ID
	Get(id string) (interface{}, error)
	
	// Save сохраняет сущность
	Save(entity interface{}) error
	
	// Delete удаляет сущность по ID
	Delete(id string) error
}

// Logger описывает интерфейс для логирования
type logger interface {
	Log(message string, level int)
	Error(err error)
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
	assert.Equal(t, 2, len(structure.Types), "Должно быть извлечено 2 интерфейса")

	// Проверка публичного интерфейса
	var publicInterface, privateInterface *models.Type
	for _, t := range structure.Types {
		if t.Name == "Repository" {
			publicInterface = t
		} else if t.Name == "logger" {
			privateInterface = t
		}
	}

	assert.NotNil(t, publicInterface, "Публичный интерфейс должен быть извлечен")
	assert.Equal(t, "interface", publicInterface.Kind, "Kind должен быть 'interface'")
	assert.True(t, publicInterface.IsPublic, "Интерфейс Repository должен быть публичным")
	assert.Equal(t, 3, len(publicInterface.Methods), "Публичный интерфейс должен иметь 3 метода")

	// Проверка методов интерфейса
	var getMethod, saveMethod, deleteMethod *models.Method
	for _, m := range publicInterface.Methods {
		if m.Name == "Get" {
			getMethod = m
		} else if m.Name == "Save" {
			saveMethod = m
		} else if m.Name == "Delete" {
			deleteMethod = m
		}
	}

	assert.NotNil(t, getMethod, "Метод Get должен быть извлечен")
	assert.Equal(t, 1, len(getMethod.Parameters), "Метод Get должен иметь 1 параметр")
	assert.Equal(t, "interface{}, error", getMethod.ReturnType, "Метод Get должен возвращать (interface{}, error)")

	assert.NotNil(t, saveMethod, "Метод Save должен быть извлечен")
	assert.NotNil(t, deleteMethod, "Метод Delete должен быть извлечен")

	// Проверка приватного интерфейса
	assert.NotNil(t, privateInterface, "Приватный интерфейс должен быть извлечен")
	assert.Equal(t, "interface", privateInterface.Kind, "Kind должен быть 'interface'")
	assert.False(t, privateInterface.IsPublic, "Интерфейс logger должен быть приватным")
	assert.Equal(t, 2, len(privateInterface.Methods), "Приватный интерфейс должен иметь 2 метода")
}

// TestGoParserGenerics проверяет корректность извлечения дженериков из Go файла
func TestGoParserGenerics(t *testing.T) {
	// Пример содержимого Go файла с дженериками
	content := `package example

// Stack представляет собой обобщенный стек
type Stack[T any] struct {
	items []T
}

// NewStack создает новый стек
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{items: make([]T, 0)}
}

// Push добавляет элемент в стек
func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

// Pop извлекает элемент из стека
func (s *Stack[T]) Pop() (T, bool) {
	var zero T
	if len(s.items) == 0 {
		return zero, false
	}
	
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item, true
}

// Generic функция
func Map[T, U any](slice []T, f func(T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = f(v)
	}
	return result
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

	// Проверка типа с дженериками
	var stackType *models.Type
	for _, t := range structure.Types {
		if t.Name == "Stack" {
			stackType = t
			break
		}
	}

	assert.NotNil(t, stackType, "Тип Stack должен быть извлечен")
	assert.True(t, stackType.IsPublic, "Тип Stack должен быть публичным")
	assert.Equal(t, "struct", stackType.Kind, "Kind типа должен быть 'struct'")
	assert.Contains(t, stackType.Name, "Stack", "Имя типа должно содержать Stack")

	// Проверка методов дженерик-типа
	var pushMethod, popMethod *models.Method
	for _, m := range structure.Methods {
		if m.Name == "Push" && m.BelongsTo == "Stack" {
			pushMethod = m
		} else if m.Name == "Pop" && m.BelongsTo == "Stack" {
			popMethod = m
		}
	}

	assert.NotNil(t, pushMethod, "Метод Push должен быть извлечен")
	assert.Equal(t, 1, len(pushMethod.Parameters), "Метод Push должен иметь 1 параметр")

	assert.NotNil(t, popMethod, "Метод Pop должен быть извлечен")
	assert.Equal(t, "T, bool", popMethod.ReturnType, "Метод Pop должен возвращать (T, bool)")

	// Проверка функции с дженериками
	var mapFunction *models.Method
	for _, m := range structure.Methods {
		if m.Name == "Map" && m.BelongsTo == "" {
			mapFunction = m
			break
		}
	}

	assert.NotNil(t, mapFunction, "Функция Map должна быть извлечена")
	assert.Equal(t, "Map", mapFunction.Name, "Имя функции должно быть Map")
	assert.Equal(t, 2, len(mapFunction.Parameters), "Функция Map должна иметь 2 параметра")
	assert.Equal(t, "[]U", mapFunction.ReturnType, "Функция Map должна возвращать []U")
}

// TestGoParserEmbeddedTypes проверяет корректность извлечения встроенных типов
func TestGoParserEmbeddedTypes(t *testing.T) {
	// Пример содержимого Go файла со встроенными типами
	content := `package example

import "io"

// BaseReader предоставляет базовую функциональность чтения
type BaseReader struct {
	source io.Reader
}

// Read реализует интерфейс io.Reader
func (br *BaseReader) Read(p []byte) (n int, err error) {
	return br.source.Read(p)
}

// EnhancedReader расширяет BaseReader дополнительной функциональностью
type EnhancedReader struct {
	BaseReader
	bufferSize int
}

// NewEnhancedReader создает новый EnhancedReader
func NewEnhancedReader(r io.Reader, size int) *EnhancedReader {
	return &EnhancedReader{
		BaseReader: BaseReader{source: r},
		bufferSize: size,
	}
}

// ReadAll читает все данные
func (er *EnhancedReader) ReadAll() ([]byte, error) {
	// реализация метода
	return nil, nil
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

	// Проверка базового типа
	var baseReader *models.Type
	for _, t := range structure.Types {
		if t.Name == "BaseReader" {
			baseReader = t
			break
		}
	}

	assert.NotNil(t, baseReader, "Тип BaseReader должен быть извлечен")
	assert.True(t, baseReader.IsPublic, "Тип BaseReader должен быть публичным")

	// Проверка расширенного типа
	var enhancedReader *models.Type
	for _, t := range structure.Types {
		if t.Name == "EnhancedReader" {
			enhancedReader = t
			break
		}
	}

	assert.NotNil(t, enhancedReader, "Тип EnhancedReader должен быть извлечен")
	assert.True(t, enhancedReader.IsPublic, "Тип EnhancedReader должен быть публичным")

	// Проверка встроенного типа
	foundEmbedded := false
	for _, prop := range enhancedReader.Properties {
		if prop.Name == "BaseReader" && prop.Type == "" {
			// В Go встроенные типы обычно не имеют явно указанного типа
			foundEmbedded = true
			break
		}
	}

	assert.True(t, foundEmbedded, "Встроенный тип BaseReader должен быть обнаружен")

	// Проверка методов расширенного типа
	var readAllMethod *models.Method
	for _, m := range structure.Methods {
		if m.Name == "ReadAll" && m.BelongsTo == "EnhancedReader" {
			readAllMethod = m
			break
		}
	}

	assert.NotNil(t, readAllMethod, "Метод ReadAll должен быть извлечен")
	assert.Equal(t, "[]byte, error", readAllMethod.ReturnType, "Метод ReadAll должен возвращать ([]byte, error)")
}
