package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/matheusbuniotto/goagent/pkg/agent"
)

// Modelos pré-definidos do OpenRouter
type OpenRouterModel struct {
	ID          string
	Name        string
	Description string
	CostLevel   string // "Gratuito", "Baixo", "Médio", "Alto"
}

var PredefinedModels = []OpenRouterModel{
	{"openai/gpt-4o-mini", "GPT-4o Mini", "Modelo econômico da OpenAI", "Baixo"},
	{"openai/gpt-4o", "GPT-4o", "Modelo avançado da OpenAI", "Médio"},
	{"anthropic/claude-3.5-sonnet", "Claude 3.5 Sonnet", "Modelo avançado da Anthropic", "Médio"},
	{"google/gemini-pro-1.5", "Gemini Pro 1.5", "Modelo avançado do Google", "Médio"},
	{"meta-llama/llama-3.1-8b-instruct", "Llama 3.1 8B", "Modelo open source da Meta", "Baixo"},
	{"anthropic/claude-3-haiku", "Claude 3 Haiku", "Modelo rápido da Anthropic", "Baixo"},
	{"openai/gpt-3.5-turbo", "GPT-3.5 Turbo", "Modelo clássico da OpenAI", "Baixo"},
	{"mistralai/mixtral-8x7b-instruct", "Mixtral 8x7B", "Modelo da Mistral AI", "Baixo"},
}

// SelectOpenRouterModel permite ao usuário escolher um modelo interativamente
func SelectOpenRouterModel() string {
	fmt.Println("\n🤖 Selecione um modelo do OpenRouter:")
	fmt.Println("─────────────────────────────────────")
	
	for i, model := range PredefinedModels {
		fmt.Printf("%d. %s (%s)\n   %s [Custo: %s]\n", 
			i+1, model.Name, model.ID, model.Description, model.CostLevel)
	}
	
	fmt.Print("\nDigite o número do modelo desejado (1-8): ")
	
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if !scanner.Scan() {
			fmt.Println("Erro ao ler entrada. Usando modelo padrão.")
			return PredefinedModels[0].ID
		}
		
		input := strings.TrimSpace(scanner.Text())
		choice, err := strconv.Atoi(input)
		
		if err != nil || choice < 1 || choice > len(PredefinedModels) {
			fmt.Printf("Opção inválida. Digite um número de 1 a %d: ", len(PredefinedModels))
			continue
		}
		
		selectedModel := PredefinedModels[choice-1]
		fmt.Printf("✅ Modelo selecionado: %s\n", selectedModel.Name)
		return selectedModel.ID
	}
}

// OpenRouter Client
type openRouterClient struct {
	apiKey     string
	httpClient *http.Client
	model      string // Modelo padrão a ser usado
}

type openRouterRequest struct {
	Model     string          `json:"model"`
	Messages  []agent.Message `json:"messages"`
	MaxTokens int             `json:"max_tokens,omitempty"`
	Stream    bool            `json:"stream"`
}

type openRouterResponse struct {
	Choices []struct {
		Message agent.Message `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}

func NewOpenRouterClient(apiKey string) LLMClient {
	// Usa modelo padrão, mas pode ser alterado via seleção interativa
	model := "meta-llama/llama-3.1-8b-instruct" // Modelo barato para testes
	return &openRouterClient{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 60 * time.Second}, // Timeout maior para gateway
		model:      model,
	}
}

func NewOpenRouterClientWithModel(apiKey string, model string) LLMClient {
	return &openRouterClient{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 60 * time.Second}, // Timeout maior para gateway
		model:      model,
	}
}

func (c *openRouterClient) GenerateResponse(ctx context.Context, history []Message, tools []Tool) (string, error) {
	systemPrompt := agent.BuildSystemPrompt(tools)
	messages := []agent.Message{{Role: "system", Content: systemPrompt}}
	messages = append(messages, history...)

	reqBody := openRouterRequest{
		Model:     c.model,
		Messages:  messages,
		MaxTokens: 1000,
		Stream:    false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("erro ao fazer marshal do JSON: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("erro ao criar request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("HTTP-Referer", "https://github.com/matheusbuniotto/goagent") // Opcional mas recomendado
	req.Header.Set("X-Title", "goAgent") // Opcional mas recomendado

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("erro ao fazer request para OpenRouter: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("erro ao ler resposta: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("OpenRouter retornou erro %d: %s", resp.StatusCode, string(body))
	}

	var openRouterResp openRouterResponse
	if err := json.Unmarshal(body, &openRouterResp); err != nil {
		return "", fmt.Errorf("erro ao fazer unmarshal da resposta: %w", err)
	}

	if openRouterResp.Error != nil {
		return "", fmt.Errorf("erro da API OpenRouter: %s", openRouterResp.Error.Message)
	}

	if len(openRouterResp.Choices) == 0 {
		return "", fmt.Errorf("nenhuma resposta recebida do OpenRouter")
	}

	return openRouterResp.Choices[0].Message.Content, nil
}