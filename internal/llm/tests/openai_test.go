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

// Тест для конструктора OpenAI провайдера
func TestNewOpenAIProvider(t *testing.T) {
	// Создаем конфигурацию
	config := map[string]interface{}{
		"APIKey": "test-api-key",
		"Model":  "gpt-3.5-turbo",
	}

	// Получаем провайдера
	provider, err := llm.GetProvider("openai", config)

	// Проверяем, что провайдер создан без ошибок
	assert.NoError(t, err)
	assert.NotNil(t, provider)
	assert.Equal(t, "openai", provider.Name())
}

// Тест для метода GenerateText
func TestOpenAIProviderGenerateText(t *testing.T) {
	// Создаем мок HTTP клиента
	mockHTTPClient := new(MockHTTPClient)

	// Создаем ожидаемый ответ
	expectedResponse := `{
		"choices": [
			{
				"message": {
					"content": "Это тестовое описание функции."
				},
				"finish_reason": "stop"
			}
		],
		"usage": {
			"total_tokens": 42
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
		"Model":  "gpt-3.5-turbo",
	}

	// Получаем провайдера через фабрику и внедряем мок HTTP клиента
	providerInterface, err := llm.GetProvider("openai", config)
	assert.NoError(t, err)

	// Приводим к конкретному типу и устанавливаем мок HTTP клиента
	openaiProvider, ok := providerInterface.(*llm.OpenAIProvider)
	if assert.True(t, ok, "Провайдер должен быть типа *llm.OpenAIProvider") {
		openaiProvider.SetHTTPClient(mockHTTPClient)
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
	assert.Equal(t, "Это тестовое описание функции.", llmResponse.Text)
	assert.Equal(t, 42, llmResponse.TokensUsed)
	assert.False(t, llmResponse.Truncated)

	// Проверяем, что HTTP клиент был вызван
	mockHTTPClient.AssertExpectations(t)
}

// Тест для метода BatchGenerateText
func TestOpenAIProviderBatchGenerateText(t *testing.T) {
	// Создаем мок HTTP клиента
	mockHTTPClient := new(MockHTTPClient)

	// Создаем ожидаемые ответы для двух запросов
	expectedResponses := []string{
		`{
			"choices": [
				{
					"message": {
						"content": "Описание функции 1"
					},
					"finish_reason": "stop"
				}
			],
			"usage": {
				"total_tokens": 30
			}
		}`,
		`{
			"choices": [
				{
					"message": {
						"content": "Описание функции 2"
					},
					"finish_reason": "stop"
				}
			],
			"usage": {
				"total_tokens": 32
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
		"Model":  "gpt-3.5-turbo",
	}

	// Получаем провайдера
	providerInterface, err := llm.GetProvider("openai", config)
	assert.NoError(t, err)

	// Устанавливаем мок HTTP клиента
	openaiProvider, ok := providerInterface.(*llm.OpenAIProvider)
	if assert.True(t, ok) {
		openaiProvider.SetHTTPClient(mockHTTPClient)
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
	assert.Equal(t, "Описание функции 1", responses[0].Text)
	assert.Equal(t, 30, responses[0].TokensUsed)
	assert.Equal(t, "Описание функции 2", responses[1].Text)
	assert.Equal(t, 32, responses[1].TokensUsed)

	// Проверяем, что HTTP клиент был вызван дважды
	mockHTTPClient.AssertExpectations(t)
}
