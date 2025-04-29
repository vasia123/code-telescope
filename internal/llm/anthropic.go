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
	RegisterProvider("anthropic", NewAnthropicProvider)
}

// AnthropicProvider реализует интерфейс LLMProvider для взаимодействия с Anthropic API
type AnthropicProvider struct {
	apiKey     string
	model      string
	baseURL    string
	httpClient *http.Client
}

// AnthropicConfig содержит параметры конфигурации для Anthropic
type AnthropicConfig struct {
	APIKey  string `json:"api_key"`
	Model   string `json:"model"`
	BaseURL string `json:"base_url"`
	Timeout int    `json:"timeout_seconds"`
}

// NewAnthropicProvider создает новый экземпляр AnthropicProvider
func NewAnthropicProvider(config map[string]interface{}) (LLMProvider, error) {
	// Преобразование map в структуру конфигурации
	jsonConfig, err := json.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("ошибка при маршалинге конфигурации: %w", err)
	}

	var cfg AnthropicConfig
	if err := json.Unmarshal(jsonConfig, &cfg); err != nil {
		return nil, fmt.Errorf("ошибка при анмаршалинге конфигурации: %w", err)
	}

	// Установка значений по умолчанию
	if cfg.Model == "" {
		cfg.Model = "claude-3-opus-20240229"
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.anthropic.com/v1/messages"
	}

	if cfg.Timeout == 0 {
		cfg.Timeout = 30
	}

	if cfg.APIKey == "" {
		return nil, fmt.Errorf("не указан API ключ для Anthropic")
	}

	return &AnthropicProvider{
		apiKey:  cfg.APIKey,
		model:   cfg.Model,
		baseURL: cfg.BaseURL,
		httpClient: &http.Client{
			Timeout: time.Duration(cfg.Timeout) * time.Second,
		},
	}, nil
}

// Name возвращает имя провайдера
func (p *AnthropicProvider) Name() string {
	return "anthropic"
}

// anthropicRequest представляет запрос к Anthropic API
type anthropicRequest struct {
	Model       string             `json:"model"`
	Messages    []anthropicMessage `json:"messages"`
	MaxTokens   int                `json:"max_tokens,omitempty"`
	Temperature float64            `json:"temperature"`
}

// anthropicMessage представляет сообщение в запросе к Anthropic API
type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// anthropicResponse представляет ответ от Anthropic API
type anthropicResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
	Usage struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
	StopReason string `json:"stop_reason"`
}

// GenerateText отправляет запрос к Anthropic API и возвращает сгенерированный текст
func (p *AnthropicProvider) GenerateText(ctx context.Context, request LLMRequest) (LLMResponse, error) {
	apiRequest := anthropicRequest{
		Model: p.model,
		Messages: []anthropicMessage{
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
	httpReq.Header.Set("X-Api-Key", p.apiKey)
	httpReq.Header.Set("Anthropic-Version", "2023-06-01")

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return LLMResponse{}, fmt.Errorf("ошибка при выполнении HTTP-запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return LLMResponse{}, fmt.Errorf("ошибка API: статус %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var apiResponse anthropicResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return LLMResponse{}, fmt.Errorf("ошибка при декодировании ответа: %w", err)
	}

	if len(apiResponse.Content) == 0 {
		return LLMResponse{}, ErrInvalidResponse
	}

	truncated := false
	if apiResponse.StopReason == "max_tokens" {
		truncated = true
	}

	totalTokens := apiResponse.Usage.InputTokens + apiResponse.Usage.OutputTokens

	return LLMResponse{
		Text:       apiResponse.Content[0].Text,
		TokensUsed: totalTokens,
		Truncated:  truncated,
	}, nil
}

// BatchGenerateText отправляет несколько запросов к Anthropic API пакетом
func (p *AnthropicProvider) BatchGenerateText(ctx context.Context, requests []LLMRequest) ([]LLMResponse, error) {
	responses := make([]LLMResponse, len(requests))

	// Anthropic также не поддерживает нативный батчинг, обрабатываем запросы последовательно
	for i, req := range requests {
		resp, err := p.GenerateText(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("ошибка при обработке запроса %d: %w", i, err)
		}
		responses[i] = resp
	}

	return responses, nil
}
