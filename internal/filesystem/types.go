package filesystem

import (
	"path/filepath"

	"code-telescope/pkg/models"
)

// FileGroup представляет группу файлов, сгруппированных по директории
type FileGroup struct {
	// Имя группы (как правило, имя директории)
	Name string

	// Путь к директории относительно корня проекта
	Path string

	// Файлы в этой группе
	Files []*models.FileMetadata

	// Вложенные группы
	SubGroups []*FileGroup
}

// NewFileGroup создает новую группу файлов с указанным именем и путем
func NewFileGroup(name, path string) *FileGroup {
	return &FileGroup{
		Name:      name,
		Path:      path,
		Files:     make([]*models.FileMetadata, 0),
		SubGroups: make([]*FileGroup, 0),
	}
}

// AddFile добавляет файл в группу
func (fg *FileGroup) AddFile(file *models.FileMetadata) {
	fg.Files = append(fg.Files, file)
}

// AddSubGroup добавляет вложенную группу
func (fg *FileGroup) AddSubGroup(group *FileGroup) {
	fg.SubGroups = append(fg.SubGroups, group)
}

// FindOrCreateSubGroup находит или создает вложенную группу по пути
func (fg *FileGroup) FindOrCreateSubGroup(path string) *FileGroup {
	if path == "" || path == "." {
		return fg
	}

	parts := filepath.SplitList(path)
	if len(parts) == 0 {
		return fg
	}

	// Находим или создаем первую часть пути
	firstPart := parts[0]
	var subGroup *FileGroup

	for _, group := range fg.SubGroups {
		if group.Name == firstPart {
			subGroup = group
			break
		}
	}

	if subGroup == nil {
		// Создаем новую группу, если не нашли существующую
		subGroup = NewFileGroup(firstPart, filepath.Join(fg.Path, firstPart))
		fg.AddSubGroup(subGroup)
	}

	// Если путь состоит из нескольких частей, рекурсивно создаем группы
	if len(parts) > 1 {
		return subGroup.FindOrCreateSubGroup(filepath.Join(parts[1:]...))
	}

	return subGroup
}

// GroupFilesByDirectory группирует файлы по директориям
func GroupFilesByDirectory(files []*models.FileMetadata) *FileGroup {
	rootGroup := NewFileGroup("", "")

	for _, file := range files {
		dirPath := file.Directory
		group := rootGroup.FindOrCreateSubGroup(dirPath)
		group.AddFile(file)
	}

	return rootGroup
}
