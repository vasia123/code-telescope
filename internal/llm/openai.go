package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func init() {
	RegisterProvider("openai", NewOpenAIProvider)
}

// OpenAIProvider реализует интерфейс LLMProvider для взаимодействия с OpenAI API
type OpenAIProvider struct {
	apiKey     string
	model      string
	baseURL    string
	httpClient *http.Client
}

// OpenAIConfig содержит параметры конфигурации для OpenAI
type OpenAIConfig struct {
	APIKey  string `json:"api_key"`
	Model   string `json:"model"`
	BaseURL string `json:"base_url"`
	Timeout int    `json:"timeout_seconds"`
}

// SetHTTPClient устанавливает HTTP клиент для провайдера (используется для тестирования)
func (p *OpenAIProvider) SetHTTPClient(client interface{}) {
	if httpClient, ok := client.(interface {
		Do(*http.Request) (*http.Response, error)
	}); ok {
		p.httpClient = &http.Client{
			Transport: &openaiTransport{client: httpClient},
		}
	}
}

// openaiTransport реализует http.RoundTripper, используя клиент с методом Do
type openaiTransport struct {
	client interface {
		Do(*http.Request) (*http.Response, error)
	}
}

// RoundTrip выполняет HTTP запрос
func (t *openaiTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.client.Do(req)
}

// NewOpenAIProvider создает новый экземпляр OpenAIProvider
func NewOpenAIProvider(config map[string]interface{}) (LLMProvider, error) {
	// Преобразование map в структуру конфигурации
	jsonConfig, err := json.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("ошибка при маршалинге конфигурации: %w", err)
	}

	var cfg OpenAIConfig
	if err := json.Unmarshal(jsonConfig, &cfg); err != nil {
		return nil, fmt.Errorf("ошибка при анмаршалинге конфигурации: %w", err)
	}

	// Установка значений по умолчанию
	if cfg.Model == "" {
		cfg.Model = "gpt-4o"
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.openai.com/v1/chat/completions"
	}

	if cfg.Timeout == 0 {
		cfg.Timeout = 30
	}

	if cfg.APIKey == "" {
		return nil, fmt.Errorf("не указан API ключ для OpenAI")
	}

	return &OpenAIProvider{
		apiKey:  cfg.APIKey,
		model:   cfg.Model,
		baseURL: cfg.BaseURL,
		httpClient: &http.Client{
			Timeout: time.Duration(cfg.Timeout) * time.Second,
		},
	}, nil
}

// Name возвращает имя провайдера
func (p *OpenAIProvider) Name() string {
	return "openai"
}

// openAIRequestMessage представляет сообщение в запросе к OpenAI API
type openAIRequestMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// openAIRequest представляет запрос к OpenAI API
type openAIRequest struct {
	Model       string                 `json:"model"`
	Messages    []openAIRequestMessage `json:"messages"`
	MaxTokens   int                    `json:"max_tokens,omitempty"`
	Temperature float64                `json:"temperature"`
}

// openAIResponseChoice представляет выбор в ответе от OpenAI API
type openAIResponseChoice struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
	FinishReason string `json:"finish_reason"`
}

// openAIResponse представляет ответ от OpenAI API
type openAIResponse struct {
	Choices []openAIResponseChoice `json:"choices"`
	Usage   struct {
		TotalTokens int `json:"total_tokens"`
	} `json:"usage"`
}

// GenerateText отправляет запрос к OpenAI API и возвращает сгенерированный текст
func (p *OpenAIProvider) GenerateText(ctx context.Context, request LLMRequest) (LLMResponse, error) {
	apiRequest := openAIRequest{
		Model: p.model,
		Messages: []openAIRequestMessage{
			{
				Role:    "user",
				Content: request.Prompt,
			},
		},
		MaxTokens:   request.MaxTokens,
		Temperature: request.Temperature,
	}

	jsonData, err := json.Marshal(apiRequest)
	if err != nil {
		return LLMResponse{}, fmt.Errorf("ошибка при маршалинге запроса: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return LLMResponse{}, fmt.Errorf("ошибка при создании HTTP-запроса: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return LLMResponse{}, fmt.Errorf("ошибка при выполнении HTTP-запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return LLMResponse{}, fmt.Errorf("ошибка API: статус %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var apiResponse openAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return LLMResponse{}, fmt.Errorf("ошибка при декодировании ответа: %w", err)
	}

	if len(apiResponse.Choices) == 0 {
		return LLMResponse{}, ErrInvalidResponse
	}

	truncated := false
	if apiResponse.Choices[0].FinishReason == "length" {
		truncated = true
	}

	return LLMResponse{
		Text:       apiResponse.Choices[0].Message.Content,
		TokensUsed: apiResponse.Usage.TotalTokens,
		Truncated:  truncated,
	}, nil
}

// BatchGenerateText отправляет несколько запросов к OpenAI API пакетом
func (p *OpenAIProvider) BatchGenerateText(ctx context.Context, requests []LLMRequest) ([]LLMResponse, error) {
	responses := make([]LLMResponse, len(requests))

	// OpenAI не поддерживает батчинг, поэтому обрабатываем запросы последовательно
	for i, req := range requests {
		resp, err := p.GenerateText(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("ошибка при обработке запроса %d: %w", i, err)
		}
		responses[i] = resp
	}

	return responses, nil
}
