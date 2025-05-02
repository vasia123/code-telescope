package markdown

// Шаблоны для генерации Markdown документации

// FileHeaderTemplate шаблон для заголовка файла
const FileHeaderTemplate = "## %s\n\n"

// ImportsExportsTemplate шаблон для секции импортов/экспортов
const ImportsExportsTemplate = "### Импорты/Экспорты\n```\n%s\n```\n\n"

// PublicMethodsHeaderTemplate шаблон для заголовка публичных методов
const PublicMethodsHeaderTemplate = "### Публичные методы\n\n"

// MethodTemplate шаблон для описания метода
const MethodTemplate = "#### %s\n%s\n%s\n- **Описание**: %s\n\n"

// CodeMapHeaderTemplate шаблон для заголовка карты кода
const CodeMapHeaderTemplate = "# Карта кода проекта %s\n\n"

// CodeMapIntroTemplate шаблон для введения карты кода
const CodeMapIntroTemplate = `## Общая информация

Эта карта кода представляет высокоуровневое описание проекта. Каждый файл представлен как "черный ящик" 
с его интерфейсами (импорты/экспорты) и публичными методами.

## Содержание

%s

`

// TableOfContentsItemTemplate шаблон для элемента оглавления
const TableOfContentsItemTemplate = "- [%s](#%s)\n"
