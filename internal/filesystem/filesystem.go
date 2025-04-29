package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"code-telescope/internal/config"
	"code-telescope/pkg/models"
)

// Scanner отвечает за сканирование файловой системы и сбор метаданных о файлах
type Scanner struct {
	config *config.Config
}

// New создает новый экземпляр Scanner
func New(cfg *config.Config) *Scanner {
	return &Scanner{
		config: cfg,
	}
}

// ScanProject сканирует директорию проекта и возвращает метаданные всех релевантных файлов
func (s *Scanner) ScanProject(projectPath string) ([]*models.FileMetadata, error) {
	// Получаем абсолютный путь к проекту
	absProjectPath, err := filepath.Abs(projectPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения абсолютного пути: %w", err)
	}

	// Проверяем, существует ли директория
	info, err := os.Stat(absProjectPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка доступа к директории проекта: %w", err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("указанный путь не является директорией: %s", absProjectPath)
	}

	// Список файлов для результата
	var files []*models.FileMetadata

	// Рекурсивно обходим директорию
	err = filepath.Walk(absProjectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Пропускаем недоступные файлы и директории
			return nil
		}

		// Пропускаем директории
		if info.IsDir() {
			// Проверяем глубину рекурсии
			relPath, err := filepath.Rel(absProjectPath, path)
			if err != nil {
				return nil
			}

			// Если это не корневая директория, проверяем глубину
			if relPath != "." {
				depth := len(strings.Split(relPath, string(os.PathSeparator)))
				if depth > s.config.FileSystem.MaxDepth {
					return filepath.SkipDir
				}
			}

			// Проверяем, не нужно ли пропустить директорию на основе шаблонов исключения
			if s.shouldExclude(relPath, true) {
				return filepath.SkipDir
			}

			return nil
		}

		// Работаем только с файлами
		relPath, err := filepath.Rel(absProjectPath, path)
		if err != nil {
			return nil
		}

		// Проверяем, подходит ли файл по шаблонам включения/исключения
		if !s.shouldInclude(relPath) || s.shouldExclude(relPath, false) {
			return nil
		}

		// Создаем метаданные файла
		fileMetadata, err := models.NewFileMetadata(path, absProjectPath)
		if err != nil {
			return nil
		}

		// Проверяем размер файла
		if fileMetadata.Size > s.config.Parser.MaxFileSize {
			return nil
		}

		// Добавляем файл в результат
		files = append(files, fileMetadata)

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("ошибка при сканировании директории: %w", err)
	}

	return files, nil
}

// shouldInclude проверяет, соответствует ли файл шаблонам включения
func (s *Scanner) shouldInclude(relPath string) bool {
	// Если шаблоны включения не указаны, включаем все файлы
	if len(s.config.FileSystem.IncludePatterns) == 0 {
		return true
	}

	for _, pattern := range s.config.FileSystem.IncludePatterns {
		matched, err := filepath.Match(pattern, filepath.Base(relPath))
		if err == nil && matched {
			return true
		}
	}

	return false
}

// shouldExclude проверяет, соответствует ли файл или директория шаблонам исключения
func (s *Scanner) shouldExclude(relPath string, isDir bool) bool {
	for _, pattern := range s.config.FileSystem.ExcludePatterns {
		// Проверяем полный путь
		matched, err := filepath.Match(pattern, relPath)
		if err == nil && matched {
			return true
		}

		// Проверяем только имя файла или директории
		matched, err = filepath.Match(pattern, filepath.Base(relPath))
		if err == nil && matched {
			return true
		}
	}

	return false
}
