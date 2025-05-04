# Карта кода проекта "Code Telescope"(содержит внутри себя каждый файл проекта как "черный ящик" с его интерфейсами (импорты/экспорты) и публичными методами, включая их параметры и высокоуровневые описания)

## cmd/codetelescope/main.go

### Импорты/Экспорты
```
Импорты:
- flag из "flag"
- fmt из "fmt"
- os из "os"
- filepath из "path/filepath"
- config из "code-telescope/internal/config"
- orchestrator из "code-telescope/internal/orchestrator"

Экспорты:
- Функция main() (точка входа в приложение)
```

### Публичные методы

#### func main()
- **Описание**: Точка входа в приложение, обеспечивает парсинг аргументов командной строки, загрузку конфигурации и запуск процесса генерации карты кода.

#### func loadConfig(configPath string) (*config.Config, error)
- **Входные параметры**: 
  - configPath: string - путь к конфигурационному файлу
- **Выходные параметры**: 
  - *config.Config - объект конфигурации
  - error - ошибка при загрузке конфигурации
- **Описание**: Загружает конфигурацию из указанного файла или использует значения по умолчанию.

## internal/config/config.go

### Импорты/Экспорты
```
Импорты:
- fmt из "fmt"
- os из "os"
- yaml.v3 из "gopkg.in/yaml.v3"

Экспорты:
- Структура Config
- Структура FileSystemConfig
- Структура ParserConfig
- Структура LLMConfig
- Структура MarkdownConfig
- Функция LoadConfig
- Функция DefaultConfig
- Функция validateConfig
```

### Публичные методы

#### func LoadConfig(configPath string) (*Config, error)
- **Входные параметры**: 
  - configPath: string - путь к конфигурационному файлу
- **Выходные параметры**: 
  - *Config - объект конфигурации
  - error - ошибка при загрузке конфигурации
- **Описание**: Загружает конфигурацию из указанного YAML-файла.

#### func DefaultConfig() *Config
- **Выходные параметры**: 
  - *Config - объект конфигурации со значениями по умолчанию
- **Описание**: Создает конфигурацию с предустановленными значениями по умолчанию.

#### func validateConfig(cfg *Config) error
- **Входные параметры**: 
  - config: *Config - объект конфигурации для проверки
- **Выходные параметры**: 
  - error - ошибка валидации или nil при успешной проверке
- **Описание**: Проверяет корректность настроек конфигурации.

### Публичные типы и структуры

#### type Config struct
- **Поля**:
  - FileSystem: FileSystemConfig - настройки для модуля файловой системы
  - Parser: ParserConfig - настройки для модуля парсинга кода
  - LLM: LLMConfig - настройки для модуля взаимодействия с ЛЛМ
  - Markdown: MarkdownConfig - настройки для модуля генерации Markdown
- **Описание**: Представляет основную конфигурацию приложения.

#### type FileSystemConfig struct
- **Поля**:
  - IncludePatterns: []string - шаблоны для включения файлов
  - ExcludePatterns: []string - шаблоны для исключения файлов
  - MaxDepth: int - максимальная глубина сканирования
- **Описание**: Содержит настройки для модуля файловой системы.

#### type ParserConfig struct
- **Поля**:
  - ParsePrivateMethods: bool - флаг парсинга приватных методов
  - MaxFileSize: int64 - максимальный размер файла для парсинга
- **Описание**: Содержит настройки для модуля парсинга кода.

#### type LLMConfig struct
- **Поля**:
  - Provider: string - провайдер ЛЛМ
  - Model: string - модель ЛЛМ
  - Temperature: float64 - температура (креативность) генерации
  - MaxTokens: int - максимальное количество токенов
  - BatchSize: int - размер пакета запросов
  - BatchDelay: int - задержка между пакетами
- **Описание**: Содержит настройки для модуля взаимодействия с ЛЛМ.

#### type MarkdownConfig struct
- **Поля**:
  - IncludeTOC: bool - включать оглавление
  - IncludeFileInfo: bool - включать информацию о файле
  - MaxMethodDescriptionLen: int - максимальная длина описания метода
  - GroupMethodsByType: bool - группировать методы по типу
  - CodeStyle: string - стиль кода
- **Описание**: Содержит настройки для модуля генерации Markdown.

## internal/config/defaults.go

### Импорты/Экспорты
```
Импорты:
- Нет импортов

Экспорты:
- Константы по умолчанию
- Переменные по умолчанию
```

### Константы и переменные

#### Константы
- DefaultMaxDepth = 10
- DefaultMaxFileSize = 1048576 (1MB)
- DefaultParsePrivateItems = false
- DefaultLLMProvider = "openai"
- DefaultLLMModel = "gpt-4"
- DefaultTemperature = 0.3
- DefaultMaxTokens = 1000
- DefaultBatchSize = 5
- DefaultBatchDelay = 1
- DefaultIncludeTOC = true
- DefaultIncludeFileInfo = true
- DefaultMaxMethodDescriptionLen = 200
- DefaultGroupMethodsByType = true
- DefaultCodeStyle = "github"

#### Переменные
- DefaultIncludePatterns - шаблоны включения файлов по умолчанию
- DefaultExcludePatterns - шаблоны исключения файлов по умолчанию
- SupportedLLMProviders - поддерживаемые провайдеры ЛЛМ
- SupportedCodeStyles - поддерживаемые стили кода

## internal/filesystem/filesystem.go

### Импорты/Экспорты
```
Импорты:
- fmt из "fmt"
- os из "os"
- filepath из "path/filepath"
- strings из "strings"
- config из "code-telescope/internal/config"
- models из "code-telescope/pkg/models"

Экспорты:
- Структура Scanner
- Функция New
```

### Публичные методы

#### func New(cfg *config.Config) *Scanner
- **Входные параметры**: 
  - config: *config.Config - объект конфигурации
- **Выходные параметры**: 
  - *Scanner - экземпляр модуля файловой системы
- **Описание**: Создает новый экземпляр Scanner с указанной конфигурацией.

#### func (s *Scanner) ScanProject(projectPath string) ([]*models.FileMetadata, error)
- **Входные параметры**: 
  - projectPath: string - путь к директории проекта
- **Выходные параметры**: 
  - []*models.FileMetadata - массив метаданных файлов
  - error - ошибка при сканировании
- **Описание**: Сканирует директорию проекта, выполняя фильтрацию файлов на основе конфигурации, и возвращает метаданные релевантных файлов кода.

#### func (s *Scanner) shouldInclude(relPath string) bool
- **Входные параметры**: 
  - relPath: string - относительный путь к файлу
- **Выходные параметры**: 
  - bool - результат проверки
- **Описание**: Проверяет, соответствует ли файл шаблонам включения из конфигурации.

#### func (s *Scanner) shouldExclude(relPath string, isDir bool) bool
- **Входные параметры**: 
  - relPath: string - относительный путь к файлу
  - isDir: bool - флаг директории
- **Выходные параметры**: 
  - bool - результат проверки
- **Описание**: Проверяет, соответствует ли файл или директория шаблонам исключения из конфигурации.

## internal/filesystem/types.go

### Импорты/Экспорты
```
Импорты:
- filepath из "path/filepath"
- models из "code-telescope/pkg/models"

Экспорты:
- Структура FileGroup
- Функция NewFileGroup
- Функция GroupFilesByDirectory
```

### Публичные типы и структуры

#### type FileGroup struct
- **Поля**:
  - Name: string - имя группы (обычно имя директории)
  - Path: string - путь к директории относительно корня проекта
  - Files: []*models.FileMetadata - файлы в этой группе
  - SubGroups: []*FileGroup - вложенные группы
- **Описание**: Представляет группу файлов, сгруппированных по директории.

### Публичные методы

#### func NewFileGroup(name, path string) *FileGroup
- **Входные параметры**: 
  - name: string - имя группы
  - path: string - путь к группе
- **Выходные параметры**: 
  - *FileGroup - новая группа файлов
- **Описание**: Создает новую группу файлов с указанным именем и путем.

#### func (fg *FileGroup) AddFile(file *models.FileMetadata)
- **Входные параметры**: 
  - file: *models.FileMetadata - файл для добавления
- **Описание**: Добавляет файл в группу.

#### func (fg *FileGroup) AddSubGroup(group *FileGroup)
- **Входные параметры**: 
  - group: *FileGroup - группа для добавления
- **Описание**: Добавляет вложенную группу.

#### func (fg *FileGroup) FindOrCreateSubGroup(path string) *FileGroup
- **Входные параметры**: 
  - path: string - путь к подгруппе
- **Выходные параметры**: 
  - *FileGroup - найденная или созданная группа
- **Описание**: Находит или создает вложенную группу по указанному пути.

#### func GroupFilesByDirectory(files []*models.FileMetadata) *FileGroup
- **Входные параметры**: 
  - files: []*models.FileMetadata - список файлов
- **Выходные параметры**: 
  - *FileGroup - корневая группа с группировкой по директориям
- **Описание**: Группирует файлы по директориям, создавая иерархическую структуру.

## internal/orchestrator/orchestrator.go

### Импорты/Экспорты
```
Импорты:
- context из "context"
- os из "os"
- filepath из "path/filepath"
- strings из "strings"
- time из "time"
- config из "code-telescope/internal/config"
- filesystem из "code-telescope/internal/filesystem"
- llm из "code-telescope/internal/llm"
- logger из "code-telescope/internal/logger"
- markdown из "code-telescope/internal/markdown"
- parser из "code-telescope/internal/parser"
- models из "code-telescope/pkg/models"

Экспорты:
- Структура Orchestrator
- Функция New
```

### Публичные типы и структуры

#### type Orchestrator struct
- **Поля**:
  - config: *config.Config - объект конфигурации
  - verbose: bool - флаг подробного вывода
  - scanner: *filesystem.Scanner - модуль файловой системы
  - parserFactory: *parser.LanguageFactory - фабрика парсеров
  - llmProvider: llm.LLMProvider - провайдер ЛЛМ
  - promptBuilder: *llm.PromptBuilder - конструктор промптов
  - mdGenerator: *markdown.Generator - генератор Markdown
- **Описание**: Оркестратор координирует весь процесс генерации карты кода.

### Публичные методы

#### func New(config *config.Config, verbose bool) (*Orchestrator, error)
- **Входные параметры**: 
  - config: *config.Config - объект конфигурации
  - verbose: bool - флаг подробного вывода
- **Выходные параметры**: 
  - *Orchestrator - экземпляр оркестратора
  - error - ошибка при инициализации
- **Описание**: Создает новый экземпляр оркестратора с указанной конфигурацией и инициализирует все необходимые компоненты.

#### func (o *Orchestrator) GenerateCodeMap(projectPath string) (string, error)
- **Входные параметры**: 
  - projectPath: string - путь к директории проекта
- **Выходные параметры**: 
  - string - сгенерированная markdown-карта кода
  - error - ошибка при генерации
- **Описание**: Координирует весь процесс генерации карты кода, включая сканирование файловой системы, парсинг кода, взаимодействие с ЛЛМ и генерацию Markdown.
- **Требует доработки**: 
  1. Существует несоответствие между типами CodeStructure и FileStructure, что вызывает ошибки при обработке результатов парсинга.
  2. Неправильная обработка методов при отправке запросов к ЛЛМ.
  3. Необходима корректная группировка методов по принадлежности к типам.

#### func (o *Orchestrator) SaveCodeMap(codeMap string, outputPath string) error
- **Входные параметры**: 
  - codeMap: string - сгенерированная карта кода
  - outputPath: string - путь для сохранения выходного файла
- **Выходные параметры**: 
  - error - ошибка при сохранении
- **Описание**: Сохраняет сгенерированную карту кода в файл по указанному пути.

## internal/parser/parser.go

### Импорты/Экспорты
```
Импорты:
- config из "code-telescope/internal/config"
- models из "code-telescope/pkg/models"
- github.com/tree-sitter/go-tree-sitter

Экспорты:
- Интерфейс Parser
```

### Публичные интерфейсы

#### type Parser interface
- **Методы**:
    - `Parse(fileMetadata *models.FileMetadata) (*models.CodeStructure, error)`: Разбирает содержимое файла и возвращает его структуру.
    - `GetSupportedExtensions() []string`: Возвращает список расширений файлов, поддерживаемых этим парсером.
    - `GetLanguageName() string`: Возвращает имя языка программирования, поддерживаемого этим парсером.
    - `ParseTreeNode(node *sitter.Node, structure *models.CodeStructure, content []byte) error`: Разбирает узел дерева синтаксического анализа Tree-sitter (специфично для реализаций на Tree-sitter).
- **Описание**: Определяет общий интерфейс для всех парсеров кода, используемых в проекте.

## internal/parser/treesitter_parser.go

### Импорты/Экспорты
```
Импорты:
- os
- sync
- code-telescope/pkg/models
- github.com/tree-sitter/go-tree-sitter

Экспорты:
- Структура TreeSitterParser
- Функция NewTreeSitterParser
```

### Публичные методы и структуры

#### type TreeSitterParser struct
- **Описание**: Предоставляет базовую реализацию парсера на основе библиотеки Tree-sitter. Содержит общую логику инициализации парсера и разбора файла с использованием переданной функции `parseTreeNodeFunc`.

#### func NewTreeSitterParser(language *sitter.Language, parseTreeNodeFunc func(node *sitter.Node, structure *models.CodeStructure, content []byte) error) *TreeSitterParser
- **Входные параметры**:
    - language: *sitter.Language - Tree-sitter язык.
    - parseTreeNodeFunc: func(...) error - Функция для разбора узлов дерева, специфичная для конкретного языка.
- **Выходные параметры**: 
    - *TreeSitterParser - Новый экземпляр базового парсера.
- **Описание**: Создает новый базовый парсер Tree-sitter с указанным языком и функцией разбора узлов.

#### func (p *TreeSitterParser) Parse(fileMetadata *models.FileMetadata) (*models.CodeStructure, error)
- **Входные параметры**:
    - fileMetadata: *models.FileMetadata - Метаданные файла для парсинга.
- **Выходные параметры**: 
    - *models.CodeStructure - Извлеченная структура кода.
    - error - Ошибка при чтении файла или парсинге.
- **Описание**: Читает файл, выполняет парсинг с помощью Tree-sitter и вызывает специфичную для языка функцию `parseTreeNodeFunc` для обхода дерева и заполнения структуры кода.

## internal/parser/language_factory.go

### Импорты/Экспорты
```
Импорты:
- fmt
- path/filepath
- strings
- sync
- code-telescope/internal/config

Экспорты:
- Тип parserConstructor (не экспортируется, используется внутри пакета)
- Переменные parserConstructors, langToExtensions, extToLangName (не экспортируются)
- Функция RegisterParser
- Структура LanguageFactory
- Функция NewLanguageFactory
```

### Публичные методы и структуры

#### func RegisterParser(languageName string, extensions []string, constructor parserConstructor)
- **Входные параметры**:
    - languageName: string - Имя языка.
    - extensions: []string - Список расширений файлов для этого языка.
    - constructor: func(*config.Config) Parser - Функция-конструктор, создающая экземпляр парсера.
- **Описание**: Регистрирует конструктор парсера для указанного языка и его расширений. Должна вызываться из `init()` функций пакетов конкретных языковых парсеров для разрыва циклической зависимости.

#### type LanguageFactory struct
- **Поля**:
    - config: *config.Config - Конфигурация приложения.
- **Описание**: Фабрика, отвечающая за предоставление нужного парсера для файла на основе его расширения.

#### func NewLanguageFactory(cfg *config.Config) *LanguageFactory
- **Входные параметры**: 
    - config: *config.Config - Конфигурация приложения.
- **Выходные параметры**: 
    - *LanguageFactory - Новый экземпляр фабрики.
- **Описание**: Создает новую фабрику парсеров.

#### func (lf *LanguageFactory) GetParserForFile(filePath string) (Parser, error)
- **Входные параметры**: 
    - filePath: string - Путь к файлу.
- **Выходные параметры**: 
    - Parser - Экземпляр парсера, подходящий для файла.
    - error - Ошибка, если парсер для данного расширения не зарегистрирован.
- **Описание**: Определяет язык файла по расширению, находит зарегистрированный конструктор и создает (или возвращает из кэша, если реализовано) экземпляр парсера.

#### func (lf *LanguageFactory) GetSupportedLanguages() []string
- **Выходные параметры**: 
    - []string - Список имен зарегистрированных языков.
- **Описание**: Возвращает список всех языков, для которых зарегистрированы парсеры.

#### func (lf *LanguageFactory) GetSupportedExtensions() []string
- **Выходные параметры**: 
    - []string - Список всех расширений, для которых зарегистрированы парсеры.
- **Описание**: Возвращает список всех расширений файлов, поддерживаемых зарегистрированными парсерами.

## internal/parser/languages/go.go

### Импорты/Экспорты
```
Импорты:
- fmt
- strings
- unsafe
- github.com/tree-sitter/go-tree-sitter
- code-telescope/internal/config
- code-telescope/internal/parser
- code-telescope/pkg/models
- C (cgo)

Экспорты:
- Функция GetGoLanguage
- Структура GoParser (не экспортируется, но реализует интерфейс parser.Parser)
- Функция NewGoParser (возвращает parser.Parser)
```

### Публичные методы и структуры

#### func init()
- **Описание**: Регистрирует `GoParser` в `parser.LanguageFactory` при инициализации пакета.

#### func GetGoLanguage() *sitter.Language
- **Выходные параметры**: 
    - *sitter.Language - Tree-sitter язык для Go.
- **Описание**: Возвращает синглтон экземпляра языка Go для Tree-sitter.

#### func NewGoParser(cfg *config.Config) parser.Parser
- **Входные параметры**: 
    - cfg: *config.Config - Конфигурация.
- **Выходные параметры**: 
    - parser.Parser - Экземпляр парсера Go.
- **Описание**: Создает новый парсер для языка Go, инициализируя базовый `TreeSitterParser`.

#### Методы GoParser (реализация интерфейса parser.Parser)
- `Parse(fileMetadata *models.FileMetadata) (*models.CodeStructure, error)`
- `GetLanguageName() string`
- `GetSupportedExtensions() []string`
- `parseTreeNode(node *sitter.Node, structure *models.CodeStructure, content []byte) error` (не экспортируется)
- **Описание**: Реализует методы интерфейса `Parser` для специфики языка Go, используя Tree-sitter для разбора кода.

## internal/parser/languages/javascript.go

### Импорты/Экспорты
```
Импорты:
- fmt
- strings
- unsafe
- github.com/tree-sitter/go-tree-sitter
- code-telescope/internal/config
- code-telescope/internal/parser
- code-telescope/pkg/models
- C (cgo)

Экспорты:
- Функция GetJavaScriptLanguage
- Структура JavaScriptParser (не экспортируется, но реализует интерфейс parser.Parser)
- Функция NewJavaScriptParser (возвращает parser.Parser)
```

### Публичные методы и структуры

#### func init()
- **Описание**: Регистрирует `JavaScriptParser` в `parser.LanguageFactory` при инициализации пакета.

#### func GetJavaScriptLanguage() *sitter.Language
- **Выходные параметры**: 
    - *sitter.Language - Tree-sitter язык для JavaScript.
- **Описание**: Возвращает синглтон экземпляра языка JavaScript для Tree-sitter.

#### func NewJavaScriptParser(cfg *config.Config) parser.Parser
- **Входные параметры**: 
    - cfg: *config.Config - Конфигурация.
- **Выходные параметры**: 
    - parser.Parser - Экземпляр парсера JavaScript.
- **Описание**: Создает новый парсер для языка JavaScript, инициализируя базовый `TreeSitterParser`.

#### Методы JavaScriptParser (реализация интерфейса parser.Parser)
- `Parse(fileMetadata *models.FileMetadata) (*models.CodeStructure, error)`
- `GetLanguageName() string`
- `GetSupportedExtensions() []string`
- `parseTreeNode(node *sitter.Node, structure *models.CodeStructure, content []byte) error` (не экспортируется)
- **Описание**: Реализует методы интерфейса `Parser` для специфики языка JavaScript, используя Tree-sitter для разбора кода.

## internal/parser/languages/python.go

### Импорты/Экспорты
```
Импорты:
- fmt
- strings
- unsafe
- github.com/tree-sitter/go-tree-sitter
- code-telescope/internal/config
- code-telescope/internal/parser
- code-telescope/pkg/models
- C (cgo)

Экспорты:
- Функция GetPythonLanguage
- Структура PythonParser (не экспортируется, но реализует интерфейс parser.Parser)
- Функция NewPythonParser (возвращает parser.Parser)
```

### Публичные методы и структуры

#### func init()
- **Описание**: Регистрирует `PythonParser` в `parser.LanguageFactory` при инициализации пакета.

#### func GetPythonLanguage() *sitter.Language
- **Выходные параметры**: 
    - *sitter.Language - Tree-sitter язык для Python.
- **Описание**: Возвращает синглтон экземпляра языка Python для Tree-sitter.

#### func NewPythonParser(cfg *config.Config) parser.Parser
- **Входные параметры**: 
    - cfg: *config.Config - Конфигурация.
- **Выходные параметры**: 
    - parser.Parser - Экземпляр парсера Python.
- **Описание**: Создает новый парсер для языка Python, инициализируя базовый `TreeSitterParser`.

#### Методы PythonParser (реализация интерфейса parser.Parser)
- `Parse(fileMetadata *models.FileMetadata) (*models.CodeStructure, error)`
- `GetLanguageName() string`
- `GetSupportedExtensions() []string`
- `parseTreeNode(node *sitter.Node, structure *models.CodeStructure, content []byte) error` (не экспортируется)
- **Описание**: Реализует методы интерфейса `Parser` для специфики языка Python, используя Tree-sitter для разбора кода.

## internal/llm/llm.go

### Импорты/Экспорты
```
Импорты:
- context из "context"
- errors из "errors"

Экспорты:
- Структура LLMRequest
- Структура LLMResponse
- Интерфейс LLMProvider
- Тип ProviderFactory
- Функция RegisterProvider
- Функция GetProvider
- Константы ошибок: ErrLLMRequestFailed, ErrInvalidPrompt, ErrInvalidResponse
```

### Публичные типы и интерфейсы

#### type LLMRequest struct
- **Поля**:
  - Prompt: string - текст промпта
  - MaxTokens: int - максимальное количество токенов в ответе
  - Temperature: float64 - температура (креативность) генерации
  - Metadata: map[string]string - дополнительные метаданные
- **Описание**: Представляет запрос к ЛЛМ.

#### type LLMResponse struct
- **Поля**:
  - Text: string - сгенерированный текст
  - TokensUsed: int - количество использованных токенов
  - Truncated: bool - флаг, указывающий, был ли ответ обрезан
- **Описание**: Представляет ответ от ЛЛМ.

#### interface LLMProvider
- **Метод**: GenerateText(ctx context.Context, request LLMRequest) (LLMResponse, error)
- **Метод**: Name() string
- **Метод**: BatchGenerateText(ctx context.Context, requests []LLMRequest) ([]LLMResponse, error)
- **Описание**: Интерфейс для взаимодействия с различными провайдерами ЛЛМ.

#### type ProviderFactory func(config map[string]interface{}) (LLMProvider, error)
- **Описание**: Тип функции, создающей экземпляр провайдера ЛЛМ на основе конфигурации.

### Публичные методы и функции

#### func RegisterProvider(name string, factory ProviderFactory)
- **Входные параметры**: 
  - name: string - имя провайдера
  - factory: ProviderFactory - фабрика для создания провайдера
- **Описание**: Регистрирует новый провайдер ЛЛМ в глобальном реестре.

#### func GetProvider(name string, config map[string]interface{}) (LLMProvider, error)
- **Входные параметры**: 
  - name: string - имя провайдера
  - config: map[string]interface{} - конфигурация провайдера
- **Выходные параметры**: 
  - LLMProvider - экземпляр провайдера
  - error - ошибка при получении провайдера
- **Описание**: Возвращает провайдер ЛЛМ по имени из реестра.

## internal/llm/openai.go

### Импорты/Экспорты
```
Импорты:
- bytes из "bytes"
- context из "context"
- encoding/json из "encoding/json"
- fmt из "fmt"
- io из "io"
- http из "net/http"
- time из "time"

Экспорты:
- Структура OpenAIProvider
- Структура OpenAIConfig
- Функция NewOpenAIProvider
```

### Публичные типы и структуры

#### type OpenAIProvider struct
- **Поля**:
  - apiKey: string - ключ API OpenAI
  - model: string - модель OpenAI
  - baseURL: string - базовый URL API
  - httpClient: *http.Client - HTTP-клиент
- **Описание**: Реализует интерфейс LLMProvider для взаимодействия с OpenAI API.

#### type OpenAIConfig struct
- **Поля**:
  - APIKey: string - ключ API OpenAI
  - Model: string - модель OpenAI
  - BaseURL: string - базовый URL API
  - Timeout: int - тайм-аут запросов в секундах
- **Описание**: Содержит параметры конфигурации для OpenAI.

### Публичные методы и функции

#### func NewOpenAIProvider(config map[string]interface{}) (LLMProvider, error)
- **Входные параметры**: 
  - config: map[string]interface{} - конфигурация провайдера
- **Выходные параметры**: 
  - LLMProvider - экземпляр провайдера OpenAI
  - error - ошибка при создании провайдера
- **Описание**: Создает новый экземпляр OpenAIProvider с указанной конфигурацией.

#### func (p *OpenAIProvider) Name() string
- **Выходные параметры**: 
  - string - имя провайдера
- **Описание**: Возвращает имя провайдера (openai).

#### func (p *OpenAIProvider) GenerateText(ctx context.Context, request LLMRequest) (LLMResponse, error)
- **Входные параметры**: 
  - ctx: context.Context - контекст запроса
  - request: LLMRequest - запрос к ЛЛМ
- **Выходные параметры**: 
  - LLMResponse - ответ от ЛЛМ
  - error - ошибка при генерации
- **Описание**: Отправляет запрос к OpenAI API и возвращает сгенерированный текст.

#### func (p *OpenAIProvider) BatchGenerateText(ctx context.Context, requests []LLMRequest) ([]LLMResponse, error)
- **Входные параметры**: 
  - ctx: context.Context - контекст запроса
  - requests: []LLMRequest - запросы к ЛЛМ
- **Выходные параметры**: 
  - []LLMResponse - ответы от ЛЛМ
  - error - ошибка при генерации
- **Описание**: Отправляет несколько запросов к OpenAI API последовательно и возвращает результаты.

## internal/llm/anthropic.go

### Импорты/Экспорты
```
Импорты:
- bytes из "bytes"
- context из "context"
- encoding/json из "encoding/json"
- fmt из "fmt"
- io из "io"
- http из "net/http"
- time из "time"

Экспорты:
- Структура AnthropicProvider
- Структура AnthropicConfig
- Функция NewAnthropicProvider
```

### Публичные типы и структуры

#### type AnthropicProvider struct
- **Поля**:
  - apiKey: string - ключ API Anthropic
  - model: string - модель Anthropic
  - baseURL: string - базовый URL API
  - httpClient: *http.Client - HTTP-клиент
- **Описание**: Реализует интерфейс LLMProvider для взаимодействия с Anthropic API.

#### type AnthropicConfig struct
- **Поля**:
  - APIKey: string - ключ API Anthropic
  - Model: string - модель Anthropic
  - BaseURL: string - базовый URL API
  - Timeout: int - тайм-аут запросов в секундах
- **Описание**: Содержит параметры конфигурации для Anthropic.

### Публичные методы и функции

#### func NewAnthropicProvider(config map[string]interface{}) (LLMProvider, error)
- **Входные параметры**: 
  - config: map[string]interface{} - конфигурация провайдера
- **Выходные параметры**: 
  - LLMProvider - экземпляр провайдера Anthropic
  - error - ошибка при создании провайдера
- **Описание**: Создает новый экземпляр AnthropicProvider с указанной конфигурацией.

#### func (p *AnthropicProvider) Name() string
- **Выходные параметры**: 
  - string - имя провайдера
- **Описание**: Возвращает имя провайдера (anthropic).

#### func (p *AnthropicProvider) GenerateText(ctx context.Context, request LLMRequest) (LLMResponse, error)
- **Входные параметры**: 
  - ctx: context.Context - контекст запроса
  - request: LLMRequest - запрос к ЛЛМ
- **Выходные параметры**: 
  - LLMResponse - ответ от ЛЛМ
  - error - ошибка при генерации
- **Описание**: Отправляет запрос к Anthropic API и возвращает сгенерированный текст.

#### func (p *AnthropicProvider) BatchGenerateText(ctx context.Context, requests []LLMRequest) ([]LLMResponse, error)
- **Входные параметры**: 
  - ctx: context.Context - контекст запроса
  - requests: []LLMRequest - запросы к ЛЛМ
- **Выходные параметры**: 
  - []LLMResponse - ответы от ЛЛМ
  - error - ошибка при генерации
- **Описание**: Отправляет несколько запросов к Anthropic API последовательно и возвращает результаты.

## internal/llm/prompt_builder.go

### Импорты/Экспорты
```
Импорты:
- fmt из "fmt"
- strings из "strings"
- models из "code-telescope/internal/models"

Экспорты:
- Структура PromptBuilder
- Функция NewPromptBuilder
```

### Публичные типы и структуры

#### type PromptBuilder struct
- **Поля**:
  - maxContextLength: int - максимальная длина контекста
- **Описание**: Предоставляет методы для создания промптов для различных задач.

### Публичные методы и функции

#### func NewPromptBuilder(maxContextLength int) *PromptBuilder
- **Входные параметры**: 
  - maxContextLength: int - максимальная длина контекста
- **Выходные параметры**: 
  - *PromptBuilder - экземпляр PromptBuilder
- **Описание**: Создает новый экземпляр PromptBuilder.

#### func (pb *PromptBuilder) BuildMethodDescriptionPrompt(methodInfo models.MethodInfo, fileContext string) string
- **Входные параметры**: 
  - methodInfo: models.MethodInfo - информация о методе
  - fileContext: string - контекст файла
- **Выходные параметры**: 
  - string - подготовленный prompt для ЛЛМ
- **Описание**: Формирует prompt для генерации описания метода.

#### func (pb *PromptBuilder) BuildFileSummaryPrompt(fileInfo models.FileStructure) string
- **Входные параметры**: 
  - fileInfo: models.FileStructure - информация о файле
- **Выходные параметры**: 
  - string - подготовленный prompt для ЛЛМ
- **Описание**: Формирует prompt для генерации общего описания файла.

#### func (pb *PromptBuilder) BuildBatchMethodPrompt(methods []models.MethodInfo, fileContext string) string
- **Входные параметры**: 
  - methods: []models.MethodInfo - список методов
  - fileContext: string - контекст файла
- **Выходные параметры**: 
  - string - подготовленный prompt для ЛЛМ
- **Описание**: Формирует prompt для пакетной обработки методов.

#### func (pb *PromptBuilder) ParseBatchResponse(response string, methods []models.MethodInfo) map[string]string
- **Входные параметры**: 
  - response: string - ответ от ЛЛМ
  - methods: []models.MethodInfo - список методов
- **Выходные параметры**: 
  - map[string]string - словарь сопоставляющий имена методов с их описаниями
- **Описание**: Разбирает ответ от ЛЛМ, содержащий описания нескольких методов.

## pkg/models/file_metadata.go

### Импорты/Экспорты
```
Импорты:
- os из "os"
- filepath из "path/filepath"
- time из "time"

Экспорты:
- Структура FileMetadata
- Функция NewFileMetadata
```

### Публичные типы и структуры

#### type FileMetadata struct
- **Поля**: 
  - Path: string - путь к файлу (относительный от корня проекта)
  - AbsolutePath: string - абсолютный путь к файлу
  - Name: string - имя файла
  - Extension: string - расширение файла
  - Size: int64 - размер файла в байтах
  - ModTime: time.Time - дата последнего изменения
  - Directory: string - родительская директория
- **Описание**: Содержит метаданные о файле исходного кода.

### Публичные методы

#### func NewFileMetadata(filePath, projectRoot string) (*FileMetadata, error)
- **Входные параметры**: 
  - filePath: string - путь к файлу
  - projectRoot: string - корень проекта
- **Выходные параметры**: 
  - *FileMetadata - метаданные файла
  - error - ошибка при создании метаданных
- **Описание**: Создает новый экземпляр FileMetadata из пути к файлу и корня проекта.

#### func (fm *FileMetadata) IsSupported() bool
- **Выходные параметры**: 
  - bool - флаг поддержки
- **Описание**: Проверяет, поддерживается ли этот тип файла.

#### func (fm *FileMetadata) IsTest() bool
- **Выходные параметры**: 
  - bool - флаг тестового файла
- **Описание**: Проверяет, является ли файл тестовым.

#### func (fm *FileMetadata) Description() string
- **Выходные параметры**: 
  - string - описание файла
- **Описание**: Возвращает строковое описание файла для вывода.

#### func (fm *FileMetadata) LanguageName() string
- **Выходные параметры**: 
  - string - название языка
- **Описание**: Возвращает название языка программирования по расширению файла.

## pkg/models/code_structure.go

### Импорты/Экспорты
```
Импорты:
- Нет явных импортов

Экспорты:
- Структура CodeStructure
- Структура Import
- Структура Export
- Структура Method
- Структура Parameter
- Структура Type
- Структура Property
- Структура Variable
- Структура Constant
- Структура Position
- Функция NewCodeStructure
```

### Публичные типы и структуры

#### type CodeStructure struct
- **Поля**: 
  - Metadata: *FileMetadata - метаданные файла
  - Imports: []*Import - импорты файла
  - Exports: []*Export - экспорты файла
  - Methods: []*Method - методы/функции файла
  - Types: []*Type - типы/классы файла
  - Variables: []*Variable - переменные верхнего уровня
  - Constants: []*Constant - константы
- **Описание**: Представляет структуру файла кода.

#### type Import struct
- **Поля**: 
  - Path: string - путь импорта
  - Alias: string - псевдоним импорта
  - Position: Position - позиция в файле
- **Описание**: Представляет импорт в файле.

#### type Export struct
- **Поля**: 
  - Name: string - имя экспортируемого элемента
  - Type: string - тип экспортируемого элемента
  - Position: Position - позиция в файле
- **Описание**: Представляет экспортируемый элемент.

#### type Method struct
- **Поля**: 
  - Name: string - имя метода
  - Parameters: []*Parameter - параметры метода
  - ReturnType: string - тип возвращаемого значения
  - IsPublic: bool - является ли метод публичным
  - IsStatic: bool - является ли метод статическим
  - Position: Position - позиция в файле
  - Description: string - описание метода
  - BelongsTo: string - принадлежность к классу/типу
- **Описание**: Представляет метод или функцию.

#### type Parameter struct
- **Поля**: 
  - Name: string - имя параметра
  - Type: string - тип параметра
  - DefaultValue: string - значение по умолчанию
  - IsRequired: bool - является ли параметр обязательным
- **Описание**: Представляет параметр метода или функции.

#### type Type struct
- **Поля**: 
  - Name: string - имя типа
  - Kind: string - тип сущности
  - IsPublic: bool - является ли публичным
  - Position: Position - позиция в файле
  - Properties: []*Property - свойства типа
  - Methods: []*Method - методы типа
- **Описание**: Представляет тип или класс.

#### type Property struct
- **Поля**: 
  - Name: string - имя свойства
  - Type: string - тип свойства
  - IsPublic: bool - является ли публичным
  - Position: Position - позиция в файле
- **Описание**: Представляет свойство класса или поле структуры.

#### type Variable struct
- **Поля**: 
  - Name: string - имя переменной
  - Type: string - тип переменной
  - IsPublic: bool - является ли публичной
  - Position: Position - позиция в файле
- **Описание**: Представляет переменную.

#### type Constant struct
- **Поля**: 
  - Name: string - имя константы
  - Type: string - тип константы
  - Value: string - значение константы
  - Position: Position - позиция в файле
- **Описание**: Представляет константу.

#### type Position struct
- **Поля**: 
  - StartLine: int - начальная строка
  - StartColumn: int - начальная колонка
  - EndLine: int - конечная строка
  - EndColumn: int - конечная колонка
- **Описание**: Представляет позицию в файле.

### Публичные методы

#### func NewCodeStructure(metadata *FileMetadata) *CodeStructure
- **Входные параметры**: 
  - metadata: *FileMetadata - метаданные файла
- **Выходные параметры**: 
  - *CodeStructure - структура кода
- **Описание**: Создает новую структуру кода для файла.

#### func (cs *CodeStructure) AddImport(imp *Import)
- **Входные параметры**: 
  - imp: *Import - импорт для добавления
- **Описание**: Добавляет импорт в структуру кода.

#### func (cs *CodeStructure) AddExport(exp *Export)
- **Входные параметры**: 
  - exp: *Export - экспорт для добавления
- **Описание**: Добавляет экспорт в структуру кода.

#### func (cs *CodeStructure) AddMethod(method *Method)
- **Входные параметры**: 
  - method: *Method - метод для добавления
- **Описание**: Добавляет метод в структуру кода.

#### func (cs *CodeStructure) AddType(typ *Type)
- **Входные параметры**: 
  - typ: *Type - тип для добавления
- **Описание**: Добавляет тип в структуру кода.

#### func (cs *CodeStructure) AddVariable(variable *Variable)
- **Входные параметры**: 
  - variable: *Variable - переменная для добавления
- **Описание**: Добавляет переменную в структуру кода.

#### func (cs *CodeStructure) AddConstant(constant *Constant)
- **Входные параметры**: 
  - constant: *Constant - константа для добавления
- **Описание**: Добавляет константу в структуру кода.

#### func (cs *CodeStructure) GetPublicMethods() []*Method
- **Выходные параметры**: 
  - []*Method - публичные методы
- **Описание**: Возвращает список публичных методов.

#### func (cs *CodeStructure) GetPublicTypes() []*Type
- **Выходные параметры**: 
  - []*Type - публичные типы
- **Описание**: Возвращает список публичных типов.

## pkg/models/method_info.go

### Импорты/Экспорты
```
Импорты:
- Нет импортов

Экспорты:
- Структура MethodInfo
```

### Публичные типы и поля

#### type MethodInfo struct
- **Поля**: 
  - Name: string - имя метода
  - Signature: string - полная сигнатура метода
  - Body: string - тело метода
  - Params: []string - параметры метода
  - Returns: []string - возвращаемые значения
- **Описание**: Содержит информацию о методе или функции.

## pkg/models/file_structure.go

### Импорты/Экспорты
```
Импорты:
- Нет импортов

Экспорты:
- Структура FileStructure
```

### Публичные типы и поля

#### type FileStructure struct
- **Поля**: 
  - Path: string - путь к файлу
  - Language: string - язык программирования
  - Imports: []string - импорты файла
  - Exports: []string - экспорты файла
  - Methods: []MethodInfo - методы файла
  - Classes: []string - классы в файле
  - Content: string - содержимое файла
  - Description: string - описание файла
- **Описание**: Представляет структурную информацию о файле кода.

## pkg/models/code_structure_converter.go

### Импорты/Экспорты
```
Импорты:
- Нет импортов

Экспорты:
- Функция ConvertToFileStructure
- Функция joinStrings
```

### Публичные методы

#### func ConvertToFileStructure(cs *CodeStructure) FileStructure
- **Входные параметры**: 
  - cs: *CodeStructure - структура кода для преобразования
- **Выходные параметры**: 
  - FileStructure - преобразованная структура файла
- **Описание**: Преобразует CodeStructure в FileStructure для совместимости с модулем генерации Markdown, обеспечивая правильное преобразование импортов, экспортов, методов и типов.

#### func joinStrings(strings []string, separator string) string
- **Входные параметры**: 
  - strings: []string - массив строк для объединения
  - separator: string - разделитель между строками
- **Выходные параметры**: 
  - string - объединенная строка
- **Описание**: Вспомогательная функция для объединения строк с указанным разделителем.

## pkg/utils/parser_utils.go

### Импорты/Экспорты
```
Импорты:
- Нет импортов

Экспорты:
- Функция FindMatchingCloseBracket
```

### Публичные методы

#### func FindMatchingCloseBracket(str string, openPos int) int
- **Входные параметры**: 
  - str: string - строка для анализа
  - openPos: int - позиция открывающей скобки
- **Выходные параметры**: 
  - int - позиция закрывающей скобки или -1, если не найдена
- **Описание**: Находит позицию закрывающей скобки, соответствующей открывающей скобке в указанной позиции. Учитывает вложенность скобок. Эта функция используется различными парсерами для обработки синтаксиса.

*Примечание: Этот документ отражает текущее состояние проекта и будет обновляться по мере его развития.* 