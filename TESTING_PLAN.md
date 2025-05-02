# План тестирования проекта "Code Telescope"

## Критичные компоненты для тестирования

Основываясь на анализе архитектуры проекта, были выделены наиболее критичные компоненты, требующие приоритетного тестирования:

1. **Парсеры языков программирования** - от их точности зависит качество всей карты кода
2. **Конвертация структур данных** - ключевая функциональность для корректной генерации Markdown
3. **Взаимодействие с ЛЛМ** - основа для генерации описаний кода
4. **Оркестрация процесса** - координация всех компонентов

## Модульные тесты

### 1. Тестирование парсеров (internal/parser/languages)

#### 1.1. Тестирование парсера Go

```go
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
    // Парсинг файла
    // Проверка корректности извлечения публичной и приватной функций
    // Проверка параметров и возвращаемых значений
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

    // Парсинг файла
    // Проверка корректности извлечения публичного типа
    // Проверка корректности извлечения полей
    // Проверка корректности связи методов с типом
}

// Дополнительные тесты для импортов, экспортов, констант и переменных
```

#### 1.2. Тестирование парсера JavaScript

```go
// TestJavaScriptParserFunctions проверяет корректность извлечения функций из JS файла
func TestJavaScriptParserFunctions(t *testing.T) {
    // Пример содержимого JS файла с функциями
    content := `
function publicFunction(param1, param2 = "default") {
    // тело функции
    return param1;
}

const arrowFunction = (param) => {
    return param * 2;
};`

    // Парсинг файла
    // Проверка корректности извлечения функций
    // Проверка параметров и значений по умолчанию
}

// TestJavaScriptParserClasses проверяет корректность извлечения классов из JS файла
func TestJavaScriptParserClasses(t *testing.T) {
    // Пример содержимого JS файла с классами и методами
    content := `
class ExampleClass {
    constructor(param) {
        this.param = param;
    }
    
    publicMethod() {
        return this.param;
    }
}

export default ExampleClass;`

    // Парсинг файла
    // Проверка корректности извлечения класса
    // Проверка корректности извлечения методов
    // Проверка корректности извлечения экспорта
}

// Дополнительные тесты для импортов, экспортов и переменных
```

### 2. Тестирование преобразования структур (pkg/models)

```go
// TestCodeStructureConversion проверяет корректность преобразования между CodeStructure и FileStructure
func TestCodeStructureConversion(t *testing.T) {
    // Создание тестового экземпляра CodeStructure с импортами, экспортами, методами и типами
    codeStructure := &models.CodeStructure{
        // Инициализация полей
    }
    
    // Конвертация в FileStructure
    fileStructure := models.ConvertToFileStructure(codeStructure)
    
    // Проверка корректности преобразования всех полей
    // Особенное внимание на преобразование методов, принадлежащих типам
}
```

### 3. Тестирование взаимодействия с ЛЛМ (internal/llm)

```go
// TestLLMProviders проверяет корректность работы провайдеров ЛЛМ
func TestLLMProviders(t *testing.T) {
    // Тестирование с использованием моков HTTP-клиента для имитации API-ответов
    
    // Создание мок-клиента, возвращающего заранее подготовленные ответы
    mockClient := &MockHTTPClient{
        // Настройка мока
    }
    
    // Тестирование OpenAI провайдера
    // Тестирование Anthropic провайдера
    // Проверка обработки ошибок
    // Проверка корректности парсинга ответов
}

// TestPromptBuilder проверяет корректность построения промптов
func TestPromptBuilder(t *testing.T) {
    // Создание тестовых данных о методах и файлах
    methodInfo := models.MethodInfo{
        Name: "TestMethod",
        Signature: "func TestMethod(param string) error",
        // Другие поля
    }
    
    // Тестирование построения промпта для описания метода
    builder := llm.NewPromptBuilder(1000)
    prompt := builder.BuildMethodDescriptionPrompt(methodInfo, "Контекст файла")
    
    // Проверка содержит ли промпт все необходимые элементы
    // Проверка не превышает ли промпт максимальную длину
}
```

### 4. Тестирование модуля конфигурации (internal/config)

```go
// TestConfigLoading проверяет корректность загрузки конфигурации из файла
func TestConfigLoading(t *testing.T) {
    // Создание временного файла конфигурации
    tempFile := createTempConfigFile(t, yamlContent)
    defer os.Remove(tempFile.Name())
    
    // Загрузка конфигурации
    cfg, err := config.LoadConfig(tempFile.Name())
    
    // Проверка отсутствия ошибок
    // Проверка корректности загруженных настроек
}

// TestConfigValidation проверяет валидацию конфигурации
func TestConfigValidation(t *testing.T) {
    // Тестирование различных невалидных конфигураций
    invalidConfigs := []struct{
        name string
        config *config.Config
        expectError bool
    }{
        // Различные тестовые случаи
    }
    
    // Проверка всех тестовых случаев
}
```

## Интеграционные тесты

### 1. Интеграция парсеров и файловой системы

```go
// TestFileSystemAndParser проверяет совместную работу модуля файловой системы и парсеров
func TestFileSystemAndParser(t *testing.T) {
    // Создание временной директории с тестовыми файлами разных типов
    testDir := createTestProjectStructure(t)
    defer os.RemoveAll(testDir)
    
    // Инициализация компонентов
    cfg := config.DefaultConfig()
    scanner := filesystem.New(cfg)
    factory := parser.NewLanguageFactory(cfg)
    
    // Сканирование директории
    files, err := scanner.ScanProject(testDir)
    
    // Проверка корректности сканирования
    
    // Парсинг каждого файла соответствующим парсером
    for _, file := range files {
        parser, err := factory.GetParserForFile(file.Path)
        // Проверка корректного выбора парсера
        
        structure, err := parser.Parse(file)
        // Проверка результатов парсинга
    }
}
```

### 2. Интеграция ЛЛМ-провайдеров и генератора Markdown

```go
// TestLLMAndMarkdownIntegration проверяет совместную работу ЛЛМ-провайдеров и генератора Markdown
func TestLLMAndMarkdownIntegration(t *testing.T) {
    // Создание моков и тестовых данных
    mockLLM := &MockLLMProvider{
        // Настройка мока для возврата предопределенных описаний
    }
    
    // Создание тестовой структуры кода
    codeStructure := createTestCodeStructure()
    
    // Создание генератора Markdown
    generator := markdown.NewGenerator(config.DefaultConfig().Markdown)
    
    // Генерация описаний методов через мок ЛЛМ
    
    // Генерация Markdown на основе структуры с описаниями
    markdownOutput := generator.GenerateFileSection(convertToFileStructure(codeStructure))
    
    // Проверка корректности сгенерированного Markdown
}
```

### 3. Тестирование оркестратора с моками

```go
// TestOrchestrator проверяет корректность работы оркестратора с моками компонентов
func TestOrchestrator(t *testing.T) {
    // Создание моков для всех зависимостей оркестратора
    mockScanner := &MockScanner{
        // Настройка мока
    }
    
    mockParserFactory := &MockParserFactory{
        // Настройка мока
    }
    
    mockLLM := &MockLLMProvider{
        // Настройка мока
    }
    
    mockMDGenerator := &MockMarkdownGenerator{
        // Настройка мока
    }
    
    // Создание оркестратора с моками
    orchestrator := orchestrator.NewWithDependencies(
        config.DefaultConfig(),
        mockScanner,
        mockParserFactory,
        mockLLM,
        mockMDGenerator,
    )
    
    // Запуск процесса генерации
    codeMap, err := orchestrator.GenerateCodeMap("test-project-path")
    
    // Проверка вызова всех зависимостей в правильном порядке
    // Проверка корректной обработки ошибок
    // Проверка результата
}
```

## E2E тесты

### 1. Тестирование на простом Go-проекте

```go
// TestEndToEndGoProject проверяет полный процесс на простом Go-проекте
func TestEndToEndGoProject(t *testing.T) {
    // Пропуск теста при отсутствии API ключей
    if os.Getenv("OPENAI_API_KEY") == "" {
        t.Skip("Skipping E2E test: OPENAI_API_KEY not set")
    }
    
    // Создание временной директории с простым Go-проектом
    projectDir := createSimpleGoProject(t)
    defer os.RemoveAll(projectDir)
    
    // Создание конфигурации
    cfg := config.DefaultConfig()
    
    // Создание оркестратора
    orch, err := orchestrator.New(cfg, true)
    
    // Генерация карты кода
    codeMap, err := orch.GenerateCodeMap(projectDir)
    
    // Проверка наличия ожидаемых разделов в сгенерированной карте
    // Проверка корректности описаний
}
```

### 2. Тестирование на JavaScript-проекте

```go
// TestEndToEndJSProject проверяет полный процесс на JavaScript-проекте
func TestEndToEndJSProject(t *testing.T) {
    // Аналогично тесту на Go-проекте
}
```

## Бенчмарки

### 1. Бенчмарки парсеров

```go
// BenchmarkGoParser измеряет производительность парсера Go
func BenchmarkGoParser(b *testing.B) {
    // Подготовка тестовых данных
    content := loadTestFile("testdata/large_go_file.go")
    
    // Создание парсера
    parser := languages.NewGoParser(config.DefaultConfig())
    
    // Запуск бенчмарка
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        parser.Parse(createFileMetadata(content))
    }
}

// Аналогичные бенчмарки для других парсеров
```

### 2. Бенчмарки ЛЛМ-запросов

```go
// BenchmarkBatchLLMRequests измеряет производительность пакетной обработки запросов к ЛЛМ
func BenchmarkBatchLLMRequests(b *testing.B) {
    // Создание тестовых данных
    methods := createTestMethods(10)
    
    // Настройка провайдера ЛЛМ с мок-клиентом
    provider := setupMockLLMProvider()
    
    // Запуск бенчмарка
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        batchDescriptions(provider, methods)
    }
}
```

## Инструменты тестирования

### 1. Мок-объекты

Для эффективного модульного и интеграционного тестирования необходимо создать мок-объекты для ключевых интерфейсов:

```go
// MockParser мок парсера для тестирования
type MockParser struct {
    mock.Mock
}

func (m *MockParser) Parse(fileMetadata *models.FileMetadata) (*models.CodeStructure, error) {
    args := m.Called(fileMetadata)
    return args.Get(0).(*models.CodeStructure), args.Error(1)
}

func (m *MockParser) GetSupportedExtensions() []string {
    args := m.Called()
    return args.Get(0).([]string)
}

func (m *MockParser) GetLanguageName() string {
    args := m.Called()
    return args.String(0)
}

// Аналогичные моки для других интерфейсов
```

### 2. Временные файлы и директории

Для тестирования файловой системы и парсеров необходимы вспомогательные функции:

```go
// createTempFile создает временный файл с заданным содержимым
func createTempFile(t *testing.T, content, extension string) *os.File {
    // Создание временного файла
    // Запись содержимого
    return file
}

// createTestProjectStructure создает временную директорию с тестовой структурой проекта
func createTestProjectStructure(t *testing.T) string {
    // Создание временной директории
    // Создание файлов разных типов
    return dirPath
}
```

## График реализации тестов

1. **Неделя 1**: Модульные тесты для парсеров и преобразования структур
2. **Неделя 2**: Модульные тесты для ЛЛМ-интерфейсов и конфигурации
3. **Неделя 3**: Интеграционные тесты
4. **Неделя 4**: E2E тесты и бенчмарки

## Метрики и критерии успеха

- **Покрытие кода**: не менее 80% для критичных компонентов
- **Проходимость тестов**: 100% проходимость на CI
- **Производительность**: установленные бенчмарки для парсеров и ЛЛМ-запросов
- **Корректность карты кода**: генерация карты для тестовых проектов без ошибок