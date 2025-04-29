package models

import (
	"os"
	"path/filepath"
	"time"
)

// FileMetadata содержит метаданные о файле исходного кода
type FileMetadata struct {
	// Путь к файлу (относительный от корня проекта)
	Path string

	// Абсолютный путь к файлу
	AbsolutePath string

	// Имя файла
	Name string

	// Расширение файла
	Extension string

	// Размер файла в байтах
	Size int64

	// Дата последнего изменения
	ModTime time.Time

	// Родительская директория
	Directory string
}

// NewFileMetadata создает новый экземпляр FileMetadata из пути к файлу и корня проекта
func NewFileMetadata(filePath, projectRoot string) (*FileMetadata, error) {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, err
	}

	fileInfo, err := os.Stat(absPath)
	if err != nil {
		return nil, err
	}

	// Создаем относительный путь от корня проекта
	relPath, err := filepath.Rel(projectRoot, absPath)
	if err != nil {
		return nil, err
	}

	return &FileMetadata{
		Path:         relPath,
		AbsolutePath: absPath,
		Name:         fileInfo.Name(),
		Extension:    filepath.Ext(fileInfo.Name()),
		Size:         fileInfo.Size(),
		ModTime:      fileInfo.ModTime(),
		Directory:    filepath.Dir(relPath),
	}, nil
}

// IsSupported проверяет, поддерживается ли этот тип файла
func (fm *FileMetadata) IsSupported() bool {
	// Список поддерживаемых расширений файлов
	supportedExtensions := map[string]bool{
		".go":   true,
		".js":   true,
		".ts":   true,
		".py":   true,
		".java": true,
		".c":    true,
		".cpp":  true,
		".h":    true,
		".hpp":  true,
	}

	return supportedExtensions[fm.Extension]
}

// IsTest проверяет, является ли файл тестовым
func (fm *FileMetadata) IsTest() bool {
	// Простая проверка на тестовый файл по имени
	switch fm.Extension {
	case ".go":
		return filepath.Base(fm.Path)[len(filepath.Base(fm.Path))-8:] == "_test.go"
	case ".py":
		name := filepath.Base(fm.Path)
		return len(name) >= 5 && name[:5] == "test_"
	default:
		return false
	}
}

// Description возвращает строковое описание файла для вывода
func (fm *FileMetadata) Description() string {
	return filepath.Join(fm.Directory, fm.Name)
}

// LanguageName возвращает название языка программирования
func (fm *FileMetadata) LanguageName() string {
	switch fm.Extension {
	case ".go":
		return "Go"
	case ".js":
		return "JavaScript"
	case ".ts":
		return "TypeScript"
	case ".py":
		return "Python"
	case ".java":
		return "Java"
	case ".c":
		return "C"
	case ".cpp":
		return "C++"
	case ".h", ".hpp":
		return "C/C++ Header"
	default:
		return "Unknown"
	}
}
