package llm

import (
	"context"
	"errors"
)

// Возможные ошибки
var (
	ErrLLMRequestFailed = errors.New("запрос к ЛЛМ не удался")
	ErrInvalidPrompt    = errors.New("некорректный промпт")
	ErrInvalidResponse  = errors.New("некорректный ответ от ЛЛМ")
)

// LLMRequest представляет запрос к ЛЛМ
type LLMRequest struct {
	Prompt      string            // Текст промпта
	MaxTokens   int               // Максимальное количество токенов в ответе
	Temperature float64           // Температура (креативность) генерации
	Metadata    map[string]string // Дополнительные метаданные
}

// LLMResponse представляет ответ от ЛЛМ
type LLMResponse struct {
	Text       string // Сгенерированный текст
	TokensUsed int    // Количество использованных токенов
	Truncated  bool   // Флаг, указывающий, был ли ответ обрезан
}

// LLMProvider интерфейс для взаимодействия с различными провайдерами ЛЛМ
type LLMProvider interface {
	// GenerateText отправляет запрос к ЛЛМ и возвращает сгенерированный текст
	GenerateText(ctx context.Context, request LLMRequest) (LLMResponse, error)

	// Name возвращает имя провайдера
	Name() string

	// BatchGenerateText отправляет несколько запросов к ЛЛМ пакетом (если поддерживается)
	BatchGenerateText(ctx context.Context, requests []LLMRequest) ([]LLMResponse, error)
}

// ProviderFactory создает экземпляр провайдера ЛЛМ на основе конфигурации
type ProviderFactory func(config map[string]interface{}) (LLMProvider, error)

// Реестр провайдеров ЛЛМ
var providerRegistry = make(map[string]ProviderFactory)

// RegisterProvider регистрирует новый провайдер ЛЛМ
func RegisterProvider(name string, factory ProviderFactory) {
	providerRegistry[name] = factory
}

// GetProvider возвращает провайдер ЛЛМ по имени
func GetProvider(name string, config map[string]interface{}) (LLMProvider, error) {
	factory, exists := providerRegistry[name]
	if !exists {
		return nil, errors.New("провайдер ЛЛМ не найден: " + name)
	}

	return factory(config)
}
