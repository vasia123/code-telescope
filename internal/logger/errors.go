package logger

import (
	"fmt"
	"runtime"
)

// ErrorType представляет тип ошибки
type ErrorType string

const (
	// ErrorTypeConfig ошибки конфигурации
	ErrorTypeConfig ErrorType = "CONFIG_ERROR"
	// ErrorTypeFileSystem ошибки файловой системы
	ErrorTypeFileSystem ErrorType = "FILESYSTEM_ERROR"
	// ErrorTypeParser ошибки парсинга
	ErrorTypeParser ErrorType = "PARSER_ERROR"
	// ErrorTypeLLM ошибки ЛЛМ
	ErrorTypeLLM ErrorType = "LLM_ERROR"
	// ErrorTypeMarkdown ошибки генерации Markdown
	ErrorTypeMarkdown ErrorType = "MARKDOWN_ERROR"
	// ErrorTypeOrchestrator ошибки оркестратора
	ErrorTypeOrchestrator ErrorType = "ORCHESTRATOR_ERROR"
	// ErrorTypeUnknown неизвестная ошибка
	ErrorTypeUnknown ErrorType = "UNKNOWN_ERROR"
)

// AppError представляет ошибку приложения со структурированной информацией
type AppError struct {
	Type    ErrorType // тип ошибки
	Message string    // сообщение об ошибке
	Err     error     // исходная ошибка
	File    string    // файл, в котором произошла ошибка
	Line    int       // строка, в которой произошла ошибка
	Func    string    // функция, в которой произошла ошибка
}

// NewAppError создает новую ошибку приложения
func NewAppError(errType ErrorType, message string, err error) *AppError {
	pc, file, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)

	return &AppError{
		Type:    errType,
		Message: message,
		Err:     err,
		File:    file,
		Line:    line,
		Func:    fn.Name(),
	}
}

// Error реализует интерфейс error
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v (at %s:%d)", e.Type, e.Message, e.Err, e.File, e.Line)
	}
	return fmt.Sprintf("[%s] %s (at %s:%d)", e.Type, e.Message, e.File, e.Line)
}

// Unwrap возвращает исходную ошибку
func (e *AppError) Unwrap() error {
	return e.Err
}

// LogError логирует ошибку и возвращает её
func LogError(err error) error {
	if err == nil {
		return nil
	}

	// Проверяем, является ли ошибка AppError
	if appErr, ok := err.(*AppError); ok {
		WithFields(Fields{
			"error_type": appErr.Type,
			"file":       appErr.File,
			"line":       appErr.Line,
			"function":   appErr.Func,
		}).Error(appErr.Message)
	} else {
		// Если это не AppError, создаем новый AppError
		pc, file, line, _ := runtime.Caller(1)
		fn := runtime.FuncForPC(pc)

		WithFields(Fields{
			"error_type": ErrorTypeUnknown,
			"file":       file,
			"line":       line,
			"function":   fn.Name(),
		}).Error(err.Error())
	}

	return err
}

// ConfigError создает ошибку конфигурации
func ConfigError(message string, err error) *AppError {
	return NewAppError(ErrorTypeConfig, message, err)
}

// FileSystemError создает ошибку файловой системы
func FileSystemError(message string, err error) *AppError {
	return NewAppError(ErrorTypeFileSystem, message, err)
}

// ParserError создает ошибку парсинга
func ParserError(message string, err error) *AppError {
	return NewAppError(ErrorTypeParser, message, err)
}

// LLMError создает ошибку ЛЛМ
func LLMError(message string, err error) *AppError {
	return NewAppError(ErrorTypeLLM, message, err)
}

// MarkdownError создает ошибку генерации Markdown
func MarkdownError(message string, err error) *AppError {
	return NewAppError(ErrorTypeMarkdown, message, err)
}

// OrchestratorError создает ошибку оркестратора
func OrchestratorError(message string, err error) *AppError {
	return NewAppError(ErrorTypeOrchestrator, message, err)
}
