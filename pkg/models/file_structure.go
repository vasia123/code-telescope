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

// GetPublicMethods возвращает все публичные методы файла
// В случае с Go, это методы, имена которых начинаются с заглавной буквы
func (fs *FileStructure) GetPublicMethods() []MethodInfo {
	// В текущей реализации просто вернём все методы, так как фильтрация по публичности 
	// уже должна была произойти в парсере
	return fs.Methods
}
