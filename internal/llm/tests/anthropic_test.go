package tests

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"code-telescope/internal/llm"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Тест для конструктора Anthropic провайдера
func TestNewAnthropicProvider(t *testing.T) {
	// Создаем конфигурацию
	config := map[string]interface{}{
		"APIKey": "test-api-key",
		"Model":  "claude-1",
	}

	// Получаем провайдера
	provider, err := llm.GetProvider("anthropic", config)

	// Проверяем, что провайдер создан без ошибок
	assert.NoError(t, err)
	assert.NotNil(t, provider)
	assert.Equal(t, "anthropic", provider.Name())
}

// Тест для метода GenerateText
func TestAnthropicProviderGenerateText(t *testing.T) {
	// Создаем мок HTTP клиента
	mockHTTPClient := new(MockHTTPClient)

	// Создаем ожидаемый ответ
	expectedResponse := `{
		"completion": "Это тестовое описание функции от Claude.",
		"stop_reason": "end_turn",
		"usage": {
			"input_tokens": 20,
			"output_tokens": 22
		}
	}`

	// Создаем Response
	response := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(expectedResponse)),
	}

	// Настраиваем ожидание вызова
	mockHTTPClient.On("Do", mock.Anything).Return(response, nil)

	// Создаем конфигурацию
	config := map[string]interface{}{
		"APIKey": "test-api-key",
		"Model":  "claude-1",
	}

	// Получаем провайдера через фабрику
	providerInterface, err := llm.GetProvider("anthropic", config)
	assert.NoError(t, err)

	// Приводим к конкретному типу и устанавливаем мок HTTP клиента
	anthropicProvider, ok := providerInterface.(*llm.AnthropicProvider)
	if assert.True(t, ok, "Провайдер должен быть типа *llm.AnthropicProvider") {
		anthropicProvider.SetHTTPClient(mockHTTPClient)
	}

	// Создаем запрос к ЛЛМ
	request := llm.LLMRequest{
		Prompt:      "Опиши функцию generateText",
		MaxTokens:   100,
		Temperature: 0.3,
	}

	// Вызываем метод генерации текста
	llmResponse, err := providerInterface.GenerateText(context.Background(), request)

	// Проверяем результаты
	assert.NoError(t, err)
	assert.Equal(t, "Это тестовое описание функции от Claude.", llmResponse.Text)
	assert.Equal(t, 42, llmResponse.TokensUsed) // Это значение может отличаться в зависимости от реализации
	assert.False(t, llmResponse.Truncated)

	// Проверяем, что HTTP клиент был вызван
	mockHTTPClient.AssertExpectations(t)
}

// Тест для метода BatchGenerateText
func TestAnthropicProviderBatchGenerateText(t *testing.T) {
	// Создаем мок HTTP клиента
	mockHTTPClient := new(MockHTTPClient)

	// Создаем ожидаемые ответы для двух запросов
	expectedResponses := []string{
		`{
			"completion": "Описание функции 1 от Claude",
			"stop_reason": "end_turn",
			"usage": {
				"input_tokens": 15,
				"output_tokens": 25
			}
		}`,
		`{
			"completion": "Описание функции 2 от Claude",
			"stop_reason": "end_turn",
			"usage": {
				"input_tokens": 15,
				"output_tokens": 25
			}
		}`,
	}

	// Настраиваем поведение мока для последовательных вызовов
	for _, respText := range expectedResponses {
		response := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(respText)),
		}
		mockHTTPClient.On("Do", mock.Anything).Return(response, nil).Once()
	}

	// Создаем конфигурацию
	config := map[string]interface{}{
		"APIKey": "test-api-key",
		"Model":  "claude-1",
	}

	// Получаем провайдера
	providerInterface, err := llm.GetProvider("anthropic", config)
	assert.NoError(t, err)

	// Устанавливаем мок HTTP клиента
	anthropicProvider, ok := providerInterface.(*llm.AnthropicProvider)
	if assert.True(t, ok) {
		anthropicProvider.SetHTTPClient(mockHTTPClient)
	}

	// Создаем запросы
	requests := []llm.LLMRequest{
		{
			Prompt:      "Опиши функцию 1",
			MaxTokens:   100,
			Temperature: 0.3,
		},
		{
			Prompt:      "Опиши функцию 2",
			MaxTokens:   100,
			Temperature: 0.3,
		},
	}

	// Вызываем метод пакетной генерации текста
	responses, err := providerInterface.BatchGenerateText(context.Background(), requests)

	// Проверяем результаты
	assert.NoError(t, err)
	assert.Len(t, responses, 2)
	assert.Equal(t, "Описание функции 1 от Claude", responses[0].Text)
	assert.Equal(t, 40, responses[0].TokensUsed) // 15 + 25 = 40
	assert.Equal(t, "Описание функции 2 от Claude", responses[1].Text)
	assert.Equal(t, 40, responses[1].TokensUsed)

	// Проверяем, что HTTP клиент был вызван дважды
	mockHTTPClient.AssertExpectations(t)
}
