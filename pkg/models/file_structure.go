package models

// FileStructure представляет структурную информацию о файле кода
type FileStructure struct {
	Path        string       // Путь к файлу
	Language    string       // Язык программирования
	Imports     []string     // Импорты файла
	Exports     []string     // Экспорты файла
	Methods     []MethodInfo // Методы файла
	Classes     []string     // Классы в файле (для объектно-ориентированных языков)
	Content     string       // Содержимое файла
	Description string       // Описание файла (может быть заполнено с помощью ЛЛМ)
}
