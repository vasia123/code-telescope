# Code Telescope

Инструмент для создания высокоуровневых карт кода проектов с использованием ЛЛМ.

## Описание

"Code Telescope" анализирует исходные файлы проекта, извлекает их структуру (импорты, экспорты, публичные методы) и генерирует Markdown-документацию, представляющую каждый файл как "черный ящик" с описанием его интерфейсов.

Инструмент использует Tree-sitter для точного парсинга кода и ЛЛМ (большие языковые модели) для генерации высокоуровневых описаний функций и методов.

## Установка

### Предварительные требования

- Go 1.18 или новее
- Tree-sitter и его грамматики для поддерживаемых языков

#### Установка Tree-sitter

1. **Windows**:

```bash
# Установка необходимых инструментов для сборки C-библиотек
# Требуется mingw-w64 или MSVC

# Если используете mingw-w64
go get -u github.com/tree-sitter/go-tree-sitter
go get -u github.com/tree-sitter/tree-sitter-go
go get -u github.com/tree-sitter/tree-sitter-javascript
go get -u github.com/tree-sitter/tree-sitter-python
```

2. **Linux/macOS**:

```bash
# Установка зависимостей для сборки
# Ubuntu/Debian
sudo apt-get install build-essential

# macOS (требуется Homebrew)
brew install cmake

# Установка Tree-sitter и его грамматик
go get -u github.com/tree-sitter/go-tree-sitter
go get -u github.com/tree-sitter/tree-sitter-go
go get -u github.com/tree-sitter/tree-sitter-javascript
go get -u github.com/tree-sitter/tree-sitter-python
```

### Установка Code Telescope

```bash
# Клонирование репозитория
git clone https://github.com/your-username/code-telescope.git
cd code-telescope

# Сборка
go build -o bin/code-telescope ./cmd/codetelescope
```

## Использование

```bash
# Анализ проекта
./bin/code-telescope -project /path/to/your/project -output map.md

# Использование с конкретной конфигурацией
./bin/code-telescope -project /path/to/your/project -config custom-config.yaml -output map.md

# Получение справки
./bin/code-telescope -help
```

## Поддерживаемые языки

В настоящее время поддерживаются следующие языки:
- Go
- JavaScript
- Python

## Конфигурация

Пример конфигурационного файла:

```yaml
parser:
  parsePrivateMethods: false
  includeDirectories:
    - src
    - lib
  excludeDirectories:
    - node_modules
    - vendor
    - .git
  excludePatterns:
    - "**/*_test.go"
    - "**/*.min.js"

llm:
  provider: "openai"
  model: "gpt-4-turbo"
  temperature: 0.2
  apiKey: "${OPENAI_API_KEY}"
  contextLimit: 4000

markdown:
  includePositionInfo: true
```

## Лицензия

MIT