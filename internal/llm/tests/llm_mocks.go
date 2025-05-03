package tests

import (
	"context"
	"net/http"

	"code-telescope/internal/llm"

	"github.com/stretchr/testify/mock"
)

// HTTPClientInterface определяет интерфейс для HTTP клиента
type HTTPClientInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

// MockHTTPClient - мок для HTTP клиента
type MockHTTPClient struct {
	mock.Mock
}

// Do - реализация метода Do для мока HTTP клиента
func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

// MockLLMProvider - мок для провайдера ЛЛМ
type MockLLMProvider struct {
	mock.Mock
}

// Name - реализация метода Name для мока провайдера ЛЛМ
func (m *MockLLMProvider) Name() string {
	args := m.Called()
	return args.String(0)
}

// GenerateText - реализация метода GenerateText для мока провайдера ЛЛМ
func (m *MockLLMProvider) GenerateText(ctx context.Context, request llm.LLMRequest) (llm.LLMResponse, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(llm.LLMResponse), args.Error(1)
}

// BatchGenerateText - реализация метода BatchGenerateText для мока провайдера ЛЛМ
func (m *MockLLMProvider) BatchGenerateText(ctx context.Context, requests []llm.LLMRequest) ([]llm.LLMResponse, error) {
	args := m.Called(ctx, requests)
	return args.Get(0).([]llm.LLMResponse), args.Error(1)
}

// CreateTestLLMResponse - вспомогательная функция для создания тестового ответа ЛЛМ
func CreateTestLLMResponse(text string) llm.LLMResponse {
	return llm.LLMResponse{
		Text:       text,
		TokensUsed: len(text) / 4, // примерная оценка токенов
		Truncated:  false,
	}
}

// SetupMockProvider - настраивает мок провайдера для возврата заданных ответов
func SetupMockProvider(methodDescriptions map[string]string) *MockLLMProvider {
	mockProvider := &MockLLMProvider{}

	// Настраиваем ответы для метода Name
	mockProvider.On("Name").Return("mock-llm")

	// Настраиваем ответы для GenerateText
	mockProvider.On("GenerateText", mock.Anything, mock.Anything).Return(
		func(ctx context.Context, request llm.LLMRequest) llm.LLMResponse {
			// Простая симуляция - возвращаем часть промпта как ответ
			return CreateTestLLMResponse("Это сгенерированное описание для метода: " + request.Prompt[:50] + "...")
		},
		func(ctx context.Context, request llm.LLMRequest) error {
			return nil
		},
	)

	// Настраиваем ответы для BatchGenerateText
	mockProvider.On("BatchGenerateText", mock.Anything, mock.Anything).Return(
		func(ctx context.Context, requests []llm.LLMRequest) []llm.LLMResponse {
			responses := make([]llm.LLMResponse, len(requests))
			for i, req := range requests {
				responses[i] = CreateTestLLMResponse("Пакетный ответ #" + string(i+'0') + ": " + req.Prompt[:30] + "...")
			}
			return responses
		},
		func(ctx context.Context, requests []llm.LLMRequest) error {
			return nil
		},
	)

	return mockProvider
}
