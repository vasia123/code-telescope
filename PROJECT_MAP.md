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
- fmt из "fmt"
- os из "os"
- config из "code-telescope/internal/config"

Экспорты:
- Структура Orchestrator
- Функция New
```

### Публичные типы и структуры

#### type Orchestrator struct
- **Поля**:
  - config: *config.Config - объект конфигурации
  - verbose: bool - флаг подробного вывода
- **Описание**: Оркестратор координирует весь процесс генерации карты кода.

### Публичные методы

#### func New(config *config.Config, verbose bool) *Orchestrator
- **Входные параметры**: 
  - config: *config.Config - объект конфигурации
  - verbose: bool - флаг подробного вывода
- **Выходные параметры**: 
  - *Orchestrator - экземпляр оркестратора
- **Описание**: Создает новый экземпляр оркестратора с указанной конфигурацией.

#### func (o *Orchestrator) GenerateCodeMap(projectPath string) (string, error)
- **Входные параметры**: 
  - projectPath: string - путь к директории проекта
- **Выходные параметры**: 
  - string - сгенерированная markdown-карта кода
  - error - ошибка при генерации
- **Описание**: Координирует весь процесс генерации карты кода, включая сканирование файловой системы, парсинг кода, взаимодействие с ЛЛМ и генерацию Markdown.

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

Экспорты:
- Интерфейс Parser
- Структура BaseParser
- Структура ParseOptions
- Функция NewBaseParser
- Функция NewParseOptions
```

### Публичные типы и интерфейсы

#### interface Parser
- **Метод**: Parse(fileMetadata *models.FileMetadata) (*models.CodeStructure, error)
- **Метод**: GetSupportedExtensions() []string
- **Метод**: GetLanguageName() string
- **Описание**: Интерфейс для парсеров различных языков программирования.

#### type BaseParser struct
- **Поля**:
  - Config: *config.Config - объект конфигурации
- **Описание**: Содержит общую функциональность для всех парсеров.

#### type ParseOptions struct
- **Поля**:
  - IncludePrivate: bool - включать приватные методы и свойства
  - Depth: int - глубина парсинга AST
  - LanguageSpecific: map[string]interface{} - опции специфичные для языка
- **Описание**: Содержит опции для процесса парсинга.

### Публичные методы

#### func NewBaseParser(cfg *config.Config) *BaseParser
- **Входные параметры**: 
  - cfg: *config.Config - объект конфигурации
- **Выходные параметры**: 
  - *BaseParser - экземпляр базового парсера
- **Описание**: Создает новый экземпляр базового парсера.

#### func NewParseOptions(cfg *config.Config) *ParseOptions
- **Входные параметры**: 
  - cfg: *config.Config - объект конфигурации
- **Выходные параметры**: 
  - *ParseOptions - опции парсинга
- **Описание**: Создает новые опции парсинга на основе конфигурации.

## internal/parser/language_factory.go

### Импорты/Экспорты
```
Импорты:
- fmt из "fmt"
- filepath из "path/filepath"
- config из "code-telescope/internal/config"

Экспорты:
- Структура LanguageFactory
- Функция NewLanguageFactory
```

### Публичные типы и структуры

#### type LanguageFactory struct
- **Поля**:
  - config: *config.Config - объект конфигурации
  - parsers: map[string]Parser - словарь парсеров по языкам
  - extToParser: map[string]string - словарь для сопоставления расширений с языками
- **Описание**: Фабрика для создания парсеров разных языков программирования.

### Публичные методы

#### func NewLanguageFactory(cfg *config.Config) *LanguageFactory
- **Входные параметры**: 
  - cfg: *config.Config - объект конфигурации
- **Выходные параметры**: 
  - *LanguageFactory - экземпляр фабрики парсеров
- **Описание**: Создает новый экземпляр фабрики парсеров с указанной конфигурацией.

#### func (lf *LanguageFactory) registerParsers()
- **Описание**: Регистрирует все поддерживаемые парсеры.

#### func (lf *LanguageFactory) registerParser(parser Parser)
- **Входные параметры**: 
  - parser: Parser - экземпляр парсера
- **Описание**: Регистрирует парсер в фабрике.

#### func (lf *LanguageFactory) GetParserForFile(filePath string) (Parser, error)
- **Входные параметры**: 
  - filePath: string - путь к файлу
- **Выходные параметры**: 
  - Parser - подходящий парсер
  - error - ошибка при получении парсера
- **Описание**: Возвращает подходящий парсер для указанного файла на основе его расширения.

#### func (lf *LanguageFactory) GetSupportedLanguages() []string
- **Выходные параметры**: 
  - []string - список поддерживаемых языков
- **Описание**: Возвращает список поддерживаемых языков программирования.

#### func (lf *LanguageFactory) GetSupportedExtensions() []string
- **Выходные параметры**: 
  - []string - список поддерживаемых расширений
- **Описание**: Возвращает список поддерживаемых расширений файлов.

## internal/parser/languages/go.go

### Импорты/Экспорты
```
Импорты:
- os из "os"
- strings из "strings"
- config из "code-telescope/internal/config"
- parser из "code-telescope/internal/parser"
- models из "code-telescope/pkg/models"

Экспорты:
- Структура GoParser
- Функция NewGoParser
```

### Публичные типы и структуры

#### type GoParser struct
- **Поля**:
  - *parser.BaseParser - встроенный базовый парсер
- **Описание**: Реализует парсер для языка Go.

### Публичные методы

#### func NewGoParser(cfg *config.Config) *GoParser
- **Входные параметры**: 
  - cfg: *config.Config - объект конфигурации
- **Выходные параметры**: 
  - *GoParser - экземпляр парсера Go
- **Описание**: Создает новый экземпляр парсера для языка Go.

#### func (p *GoParser) GetLanguageName() string
- **Выходные параметры**: 
  - string - название языка
- **Описание**: Возвращает название языка программирования (Go).

#### func (p *GoParser) GetSupportedExtensions() []string
- **Выходные параметры**: 
  - []string - поддерживаемые расширения
- **Описание**: Возвращает список поддерживаемых расширений файлов (.go).

#### func (p *GoParser) Parse(fileMetadata *models.FileMetadata) (*models.CodeStructure, error)
- **Входные параметры**: 
  - fileMetadata: *models.FileMetadata - метаданные файла
- **Выходные параметры**: 
  - *models.CodeStructure - структура кода
  - error - ошибка при парсинге
- **Описание**: Разбирает файл Go и извлекает его структуру: пакет, импорты, типы, функции, методы и переменные.

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

*Примечание: Этот документ отражает текущее состояние проекта и будет обновляться по мере его развития.* 