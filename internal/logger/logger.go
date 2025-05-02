package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

// Уровни логирования
const (
	DebugLevel = logrus.DebugLevel
	InfoLevel  = logrus.InfoLevel
	WarnLevel  = logrus.WarnLevel
	ErrorLevel = logrus.ErrorLevel
	FatalLevel = logrus.FatalLevel
	PanicLevel = logrus.PanicLevel
)

var (
	// стандартный логгер
	defaultLogger *logrus.Logger
)

// настраиваемые поля логгера
type Fields map[string]interface{}

// инициализация логгера
func init() {
	defaultLogger = logrus.New()
	defaultLogger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.Function, ".")
			funcName := s[len(s)-1]
			return funcName, fmt.Sprintf("%s:%d", filepath.Base(f.File), f.Line)
		},
	})
	defaultLogger.SetReportCaller(true)
	defaultLogger.SetOutput(os.Stdout)
	defaultLogger.SetLevel(logrus.InfoLevel)
}

// SetLevel устанавливает уровень логирования
func SetLevel(level logrus.Level) {
	defaultLogger.SetLevel(level)
}

// SetOutput устанавливает устройство вывода для логов
func SetOutput(output io.Writer) {
	defaultLogger.SetOutput(output)
}

// Debug логирует сообщение с уровнем Debug
func Debug(args ...interface{}) {
	defaultLogger.Debug(args...)
}

// Debugf логирует отформатированное сообщение с уровнем Debug
func Debugf(format string, args ...interface{}) {
	defaultLogger.Debugf(format, args...)
}

// Info логирует сообщение с уровнем Info
func Info(args ...interface{}) {
	defaultLogger.Info(args...)
}

// Infof логирует отформатированное сообщение с уровнем Info
func Infof(format string, args ...interface{}) {
	defaultLogger.Infof(format, args...)
}

// Warn логирует сообщение с уровнем Warn
func Warn(args ...interface{}) {
	defaultLogger.Warn(args...)
}

// Warnf логирует отформатированное сообщение с уровнем Warn
func Warnf(format string, args ...interface{}) {
	defaultLogger.Warnf(format, args...)
}

// Error логирует сообщение с уровнем Error
func Error(args ...interface{}) {
	defaultLogger.Error(args...)
}

// Errorf логирует отформатированное сообщение с уровнем Error
func Errorf(format string, args ...interface{}) {
	defaultLogger.Errorf(format, args...)
}

// Fatal логирует сообщение с уровнем Fatal и завершает приложение
func Fatal(args ...interface{}) {
	defaultLogger.Fatal(args...)
}

// Fatalf логирует отформатированное сообщение с уровнем Fatal и завершает приложение
func Fatalf(format string, args ...interface{}) {
	defaultLogger.Fatalf(format, args...)
}

// Panic логирует сообщение с уровнем Panic и вызывает панику
func Panic(args ...interface{}) {
	defaultLogger.Panic(args...)
}

// Panicf логирует отформатированное сообщение с уровнем Panic и вызывает панику
func Panicf(format string, args ...interface{}) {
	defaultLogger.Panicf(format, args...)
}

// WithFields добавляет структурированные поля к логу
func WithFields(fields Fields) *logrus.Entry {
	return defaultLogger.WithFields(logrus.Fields(fields))
}

// WithField добавляет одно поле к логу
func WithField(key string, value interface{}) *logrus.Entry {
	return defaultLogger.WithField(key, value)
}

// WithError добавляет ошибку к записи лога
func WithError(err error) *logrus.Entry {
	return defaultLogger.WithError(err)
}
