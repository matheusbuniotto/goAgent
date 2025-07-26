package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/matheusbuniotto/goagent/pkg/agent"
)

type geminiClient struct {
	apiKey     string
	httpClient *http.Client
}

type geminiRequest struct {
	Contents         []geminiContent `json:"contents"`
	GenerationConfig geminiGenConfig `json:"generationConfig"`
}
type geminiContent struct {
	Role  string       `json:"role"`
	Parts []geminiPart `json:"parts"`
}
type geminiPart struct {
	Text string `json:"text"`
}
type geminiGenConfig struct {
	MaxOutputTokens int `json:"maxOutputTokens"`
}
type geminiResponse struct {
	Candidates []struct {
		Content geminiContent `json:"content"`
	} `json:"candidates"`
}

func NewGeminiClient(apiKey string) agent.LLMClient {
	return &geminiClient{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *geminiClient) GenerateResponse(ctx context.Context, history []agent.Message, tools []agent.Tool) (string, error) {
	systemPrompt := agent.BuildSystemPrompt(tools)
	fullPrompt := systemPrompt + "\n\nAqui está o histórico da conversa:\n"

	var geminiContents []geminiContent
	geminiContents = append(geminiContents, geminiContent{
		Role:  "user",
		Parts: []geminiPart{{Text: fullPrompt}},
	})
	geminiContents = append(geminiContents, geminiContent{
		Role:  "model",
		Parts: []geminiPart{{Text: "Entendido. Estou pronto para ajudar."}},
	})

	for _, msg := range history {
		role := "user"
		if msg.Role == "assistant" || msg.Role == "model" {
			role = "model"
		}
		geminiContents = append(geminiContents, geminiContent{
			Role:  role,
			Parts: []geminiPart{{Text: msg.Content}},
		})
	}

	reqBody, err := json.Marshal(geminiRequest{
		Contents:         geminiContents,
		GenerationConfig: geminiGenConfig{MaxOutputTokens: 10000},
	})
	if err != nil {
		return "", fmt.Errorf("erro ao codificar requisição para Gemini: %w", err)
	}

	apiURL := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash-lite:generateContent?key=%s", c.apiKey)
	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("erro ao criar requisição para Gemini: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("erro ao enviar requisição para Gemini: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API do Gemini retornou status não-OK: %s, Body: %s", resp.Status, string(bodyBytes))
	}

	var geminiResp geminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		return "", fmt.Errorf("erro ao decodificar resposta do Gemini: %w", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("resposta do Gemini está vazia ou em formato inesperado")
	}

	return geminiResp.Candidates[0].Content.Parts[0].Text, nil
}
