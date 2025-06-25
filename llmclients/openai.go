package llmclients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/matheusbuniotto/goagent/agent"
)

type LLMClient = agent.LLMClient

type Message = agent.Message

type Tool = agent.Tool

// OpenAI Client
type openAIClient struct {
	apiKey     string
	httpClient *http.Client
}

type openAIRequest struct {
	Model     string          `json:"model"`
	Messages  []agent.Message `json:"messages"`
	MaxTokens int             `json:"max_tokens,omitempty"`
}

type openAIResponse struct {
	Choices []struct {
		Message agent.Message `json:"message"`
	} `json:"choices"`
}

func NewOpenAIClient(apiKey string) LLMClient {
	return &openAIClient{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *openAIClient) GenerateResponse(ctx context.Context, history []Message, tools []Tool) (string, error) {
	systemPrompt := agent.BuildSystemPrompt(tools)
	messages := []agent.Message{{Role: "system", Content: systemPrompt}}
	messages = append(messages, history...)

	reqBody, err := json.Marshal(openAIRequest{
		Model:     "gpt-4.1-nano",
		Messages:  messages,
		MaxTokens: 9060,
	})
	if err != nil {
		return "", fmt.Errorf("erro ao codificar requisição para OpenAI: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("erro ao criar requisição para OpenAI: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("erro ao enviar requisição para OpenAI: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API da OpenAI retornou status não-OK: %s, Body: %s", resp.Status, string(bodyBytes))
	}

	var openAIResp openAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openAIResp); err != nil {
		return "", fmt.Errorf("erro ao decodificar resposta da OpenAI: %w", err)
	}

	if len(openAIResp.Choices) == 0 {
		return "", fmt.Errorf("resposta da OpenAI não contém escolhas")
	}

	return openAIResp.Choices[0].Message.Content, nil
}
