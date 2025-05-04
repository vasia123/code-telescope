# Code Telescope

Инструмент для генерации высокоуровневых карт кода проектов с использованием ЛЛМ. Инструмент анализирует исходные файлы проекта, извлекает их структуру (импорты, экспорты, публичные методы) и генерирует Markdown-документацию, представляющую каждый файл как "черный ящик" с описанием его интерфейсов.

## Описание

Основные возможности:
- Сканирование директорий проекта
- Парсинг кода на разных языках программирования
- Генерация высокоуровневых описаний с использованием ЛЛМ
- Создание структурированной Markdown-документации

## Установка

### Предварительные требования

- Go 1.21 или новее
- Установленные библиотеки Tree-sitter для поддерживаемых языков

### Установка Tree-sitter

Для точного парсинга кода используется библиотека Tree-sitter. Установите все необходимые зависимости:

```bash
# Установка Go-биндингов для Tree-sitter
go get github.com/tree-sitter/go-tree-sitter@latest

# Установка грамматик для поддерживаемых языков
go get github.com/tree-sitter/tree-sitter-go@latest
go get github.com/tree-sitter/tree-sitter-javascript@latest
go get github.com/tree-sitter/tree-sitter-python@latest
```

### Сборка проекта

```bash
# Клонирование репозитория
git clone https://github.com/your-username/code-telescope.git
cd code-telescope

# Загрузка зависимостей
go mod download

# Сборка проекта
go build -o code-telescope ./cmd/codetelescope
```

## Использование

### Базовое использование

```bash
# Анализ проекта и генерация карты кода
./code-telescope -project /путь/к/проекту -output карта-кода.md
```

### Дополнительные параметры

```bash
# Использование собственного конфигурационного файла
./code-telescope -project /путь/к/проекту -config my-config.yaml -output карта-кода.md

# Подробный вывод процесса анализа
./code-telescope -project /путь/к/проекту -output карта-кода.md -verbose
```

## Поддерживаемые языки

В настоящее время поддерживаются следующие языки программирования:

- Go
- JavaScript
- Python

## Конфигурация

Вы можете настроить работу инструмента с помощью конфигурационного файла YAML. Пример конфигурации:

```yaml
filesystem:
  include_patterns:
    - "*.go"
    - "*.js"
    - "*.py"
  exclude_patterns:
    - "*_test.go"
    - "vendor/**"
    - "node_modules/**"
  max_depth: 10

parser:
  parse_private_methods: false
  max_file_size: 1048576

llm:
  provider: "openai"
  model: "gpt-4"
  temperature: 0.3
  max_tokens: 1000
  batch_size: 5
  batch_delay: 1

markdown:
  include_toc: true
  include_file_info: true
  max_method_description_len: 200
  group_methods_by_type: true
  code_style: "github"
```

## Требования

- Доступ к API выбранного ЛЛМ-провайдера

## Архитектура

Проект имеет модульную архитектуру, состоящую из следующих компонентов:

1. **Модуль файловой системы** - сканирование директорий и сбор метаданных
2. **Модуль парсинга кода** - извлечение структурной информации из исходных файлов
3. **Модуль взаимодействия с ЛЛМ** - генерация высокоуровневых описаний
4. **Модуль генерации Markdown** - создание документации
5. **Модуль конфигурации** - управление настройками
6. **Оркестратор** - координация всего процесса

## Лицензия

MIT