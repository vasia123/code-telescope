# План реализации проекта "Code Telescope"

## Общее описание
"Code Telescope" — инструмент на языке Go для генерации высокоуровневых карт кода проектов с использованием ЛЛМ. Инструмент анализирует исходные файлы проекта, извлекает их структуру (импорты, экспорты, публичные методы) и генерирует Markdown-документацию, представляющую каждый файл как "черный ящик" с описанием его интерфейсов.

## План реализации

### Этап 1: Настройка проекта
- [x] Создание структуры директорий проекта
- [x] Инициализация Go-модуля
- [x] Создание базового файла конфигурации
- [x] Настройка Makefile для сборки

### Этап 2: Модуль конфигурации
- [x] Реализация структуры конфигурации
- [x] Функции загрузки конфигурации из файла
- [x] Функции валидации конфигурации
- [x] Настройка параметров по умолчанию

### Этап 3: Модуль файловой системы
- [x] Реализация сканирования директорий
- [x] Фильтрация файлов на основе конфигурации
- [x] Сбор метаданных о файлах

### Этап 4: Модуль парсинга кода
- [x] Временная реализация без Tree-sitter
- [x] Исправление несоответствий между интерфейсами и реализациями
- [x] Доработка парсера Go для правильного анализа параметров и возвращаемых значений
- [x] Обработка методов, принадлежащих типам
- [x] Настройка Tree-sitter для более точного парсинга
- [x] Реализация парсеров для других языков (JavaScript)
- [x] Реализация парсера для Python

### Этап 5: Модуль взаимодействия с ЛЛМ
- [x] Абстрактный интерфейс ЛЛМ
- [x] Реализация для OpenAI API
- [x] Реализация для Anthropic API
- [x] Создание конструктора промптов
- [x] Доработка конструктора промптов для улучшения качества описаний

### Этап 6: Модуль генерации Markdown
- [x] Создание шаблонов для Markdown
- [x] Функции генерации разделов документации
- [x] Функции генерации полной карты кода
- [x] Исправление преобразования между типами CodeStructure и FileStructure

### Этап 7: Оркестратор
- [x] Координация всего процесса (базовый каркас)
- [x] Полная интеграция со всеми модулями
- [x] Обработка ошибок
- [x] Логирование
- [x] Исправление обработки методов в GenerateCodeMap
- [x] Правильное преобразование CodeStructure в формат FileStructure

### Этап 8: CLI и документация
- [x] Реализация CLI-интерфейса
- [ ] Создание пользовательской документации
- [ ] Создание примеров использования

### Этап 9: Тестирование
- [ ] Модульные тесты
  - [x] Тестирование модуля конфигурации
  - [ ] Тестирование модуля файловой системы
  - [ ] Тестирование парсеров языков программирования
    - [x] Начало тестирования парсера Go (базовая структура)
    - [x] Завершение тестирования парсера Go
    - [ ] Тестирование парсера JavaScript
    - [ ] Тестирование парсера Python (после реализации)
  - [x] Базовое тестирование преобразования между CodeStructure и FileStructure
  - [x] Тестирование ЛЛМ-интерфейсов
    - [x] Создание моков для тестирования
    - [x] Тестирование OpenAI провайдера
    - [x] Тестирование Anthropic провайдера
- [ ] Интеграционные тесты
  - [ ] Интеграция парсеров и модуля файловой системы
  - [ ] Интеграция ЛЛМ-провайдеров и генератора Markdown
  - [ ] Тестирование оркестратора с моками компонентов
- [ ] E2E тесты
  - [ ] Тестирование полного процесса на простом Go-проекте
  - [ ] Тестирование полного процесса на JavaScript-проекте
  - [ ] Тестирование на реальных проектах разного размера

## Текущий статус

На данный момент успешно завершена разработка тестов для парсера Go, который является ключевым компонентом системы. Тесты охватывают все основные аспекты синтаксиса Go, включая обработку интерфейсов, дженериков и встроенных типов.

Обнаружены ошибки совместимости в реализациях парсеров JavaScript и Python, которые необходимо исправить перед тем, как продолжить разработку тестов для них.

### Основные достижения
1. Исправлено несоответствие между интерфейсами Parser в различных модулях
2. Доработан парсер Go для правильного извлечения параметров и возвращаемых значений
3. Создан функциональный парсер JavaScript
4. Реализовано преобразование между CodeStructure и FileStructure
5. Обновлен оркестратор для корректной обработки типов данных
6. Выделены общие утилиты парсинга в отдельный пакет
7. Начата разработка тестов для критических компонентов
8. Исправлены и улучшены тесты модуля конфигурации
9. Создан и реализован парсер для Python
10. Разработаны тесты для ЛЛМ-провайдеров (OpenAI и Anthropic)
11. Интегрирован Tree-sitter для точного парсинга кода
    - Создан базовый класс TreeSitterParser
    - Реализован парсер Go с использованием Tree-sitter
    - Реализован парсер JavaScript с использованием Tree-sitter
    - Реализован парсер Python с использованием Tree-sitter
    - Обновлена LanguageFactory для работы с Tree-sitter парсерами
12. Завершена разработка тестов парсера Go
    - Тесты базовой структуры файлов
    - Тесты для функций и методов
    - Тесты для интерфейсов
    - Тесты для дженериков
    - Тесты для встроенных типов

### Проблемы, требующие решения
1. ✓ Отсутствует парсер для Python
2. ✓ Не реализована интеграция с Tree-sitter для более точного парсинга
3. Требуется создание пользовательской документации
4. Нужны примеры использования для разных языков программирования
5. ✓ Необходимо доработать тесты для парсера Go
6. Необходимо разработать тесты для парсеров JavaScript и Python
7. Существует несоответствие в реализациях парсеров языков, вызывающее ошибки компиляции
   - В JavaScript-парсере используются несуществующие типы (models.Function)
   - В Python-парсере метод parseTreeNode не соответствует требуемому интерфейсу Parser (строчная буква в названии)
   - Используются поля, которых нет в структурах моделей (IsStatic, IsVariadic и т.д.)

## План на ближайший этап

1. Исправить несоответствия в реализациях парсеров JavaScript и Python:
   - Заменить использование несуществующего типа models.Function
   - Исправить метод parseTreeNode на ParseTreeNode в Python-парсере
   - Заменить несуществующие поля (IsStatic, IsVariadic и т.д.) в моделях

2. После исправления ошибок компиляции, реализовать тесты для парсеров JavaScript и Python по аналогии с тестами для парсера Go

3. Приступить к разработке пользовательской документации и примеров использования

## Следующие шаги
1. ✓ Исправить несоответствия в тестах модуля конфигурации
2. ✓ Завершить разработку тестов для парсера Go
3. Исправить несоответствия и ошибки в реализациях парсеров JavaScript и Python
4. Разработать тесты для парсеров JavaScript и Python
5. ✓ Реализовать тесты для ЛЛМ-провайдеров
6. ✓ Реализовать парсер для Python
7. ✓ Интегрировать Tree-sitter для более точного парсинга
8. Создать пользовательскую документацию
9. Подготовить примеры использования
10. Разработать интеграционные тесты
11. Провести тестирование на реальных проектах разного размера

## План рефакторинга интеграции с Tree-sitter

### Этап 1: Анализ текущего состояния (1-2 дня)
- [x] Оценка качества существующей интеграции с Tree-sitter
- [x] Выявление проблем с текущими парсерами языков
- [x] Проверка и обновление зависимостей Tree-sitter

### Этап 2: Реорганизация базовых классов (2-3 дня)
- [x] Устранение дублирования между TreeSitterParser и BaseTreeSitterParser
- [x] Создание единого интерфейса для Tree-sitter парсеров
- [x] Упрощение механизма инициализации Tree-sitter

### Этап 3: Улучшение парсеров языков (3-4 дня)
- [x] Переработка парсера Go для повышения точности и производительности
- [x] Переработка парсера JavaScript
- [x] Переработка парсера Python
- [ ] Унификация подхода к извлечению структурных элементов

### Этап 4: Тестирование (2-3 дня)
- [ ] Создание модульных тестов для парсеров Go, JavaScript и Python
- [ ] Тестирование на различных образцах кода
- [ ] Проверка корректности извлечения импортов, экспортов и методов

### Этап 5: Документация и примеры (1-2 дня)
- [ ] Документирование подхода к интеграции с Tree-sitter
- [ ] Создание руководства по добавлению новых языков
- [ ] Обновление основной документации

## План разработки

- [ ] **Этап 1: Базовая структура и настройка**
    - [x] Инициализация проекта Go (`go mod init`)
    - [x] Определение основной структуры директорий (cmd, internal, pkg, etc.)
    - [x] Настройка базовой конфигурации (`internal/config`)
    - [x] Настройка логирования (`internal/utils/logger`)
    - [x] Создание точки входа (`cmd/codetelescope/main.go`)
    - [x] Определение базовых моделей данных (`pkg/models`)
- [ ] **Этап 2: Сканирование файловой системы**
    - [ ] Реализация `FileSystemModule` (`internal/filesystem`)
    - [ ] Метод `scanProject` для обхода директорий
    - [ ] Фильтрация файлов по расширениям и исключениям из конфигурации
    - [ ] Сбор метаданных (`FileMetadata`)
- [ ] **Этап 3: Парсинг кода (Рефакторинг)**
    - [ ] Диагностика и устранение циклической зависимости `internal/parser <-> internal/parser/languages`
    - [ ] Рефакторинг `internal/parser` для четкого разделения интерфейса и реализаций
        - [ ] Выделение интерфейса `Parser`
        - [ ] Рефакторинг `TreeSitterParser` как базовой структуры/реализации
        - [ ] Обновление фабрики парсеров (`language_factory.go`)
        - [ ] Адаптация парсеров языков (`internal/parser/languages/*`)
    - [ ] Интеграция Tree-sitter (`internal/parser/treesitter_parser.go`) - *базовая интеграция есть, требует доработки и очистки*
    - [ ] Реализация парсера для Go (`internal/parser/languages/go.go`) - *требует ревизии и обновления*
    - [ ] Добавление парсеров для других языков (JavaScript, Python) - *требует ревизии и обновления*
    - [ ] Обновление тестов для модуля парсера
- [ ] **Этап 4: Взаимодействие с ЛЛМ**
    - [ ] Реализация `LLMInterfaceModule` (`internal/llm`)
    - [ ] Абстрактный интерфейс для ЛЛМ
    - [ ] Реализация для OpenAI API (`openai.go`)
    - [ ] Реализация для Anthropic API (`anthropic.go`)
    - [ ] Конструктор промптов (`prompt_builder.go`)
    - [ ] Метод `generateMethodDescription`
    - [ ] Метод `batchGenerateDescriptions` (оптимизация)
- [ ] **Этап 5: Генерация Markdown**
    - [ ] Реализация `MarkdownGeneratorModule` (`internal/markdown`)
    - [ ] Метод `generateFileSection`
    - [ ] Метод `generateCodeMap`
    - [ ] Использование шаблонов (`templates.go`)
- [ ] **Этап 6: Оркестрация и сборка**
    - [ ] Реализация `OrchestratorModule` (`internal/orchestrator`)
    - [ ] Метод `generateCodeMap` для координации процесса
    - [ ] Обработка ошибок (`error_handler.go`)
    - [ ] Метод `saveCodeMap`
- [ ] **Этап 7: Завершение и документация**
    - [ ] Написание `README.md`
    - [ ] Написание руководства пользователя (`docs/usage.md`)
    - [ ] Написание руководства разработчика (`docs/development.md`)
    - [ ] Добавление примеров (`/examples`)
    - [ ] Настройка Makefile или скриптов сборки (`/scripts`)
    - [ ] Финальное тестирование и отладка
    - [ ] Обновление `PROJECT_MAP.md`

## Отслеживание прогресса

*   **YYYY-MM-DD:** Начало проекта. Инициализация структуры.
*   **YYYY-MM-DD:** Добавлен FileSystemModule (частично).
*   **YYYY-MM-DD:** Начата интеграция TreeSitter (требует рефакторинга).
*   **< сегодняшняя дата >:** Запланирован рефакторинг модуля парсера (`internal/parser`) для устранения циклической зависимости и улучшения структуры.