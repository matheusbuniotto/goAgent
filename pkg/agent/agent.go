package agent

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/matheusbuniotto/goagent/internal/prompts"
)

// =================================================================================================
// Interfaces e Estruturas Públicas
// =================================================================================================

// Message define a estrutura de uma única mensagem na conversa.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Tool define a interface que todas as ferramentas devem implementar.
type Tool interface {
	Name() string
	Description() string
	Execute(args string) (string, error)
}

// LLMClient é a interface para comunicação com qualquer Large Language Model.
type LLMClient interface {
	GenerateResponse(ctx context.Context, history []Message, tools []Tool) (string, error)
}

// BuildSystemPrompt cria o prompt do sistema que instrui o LLM.
func BuildSystemPrompt(tools []Tool) string {
	prompt := prompts.SystemPrompt + "\n"
	for _, tool := range tools {
		prompt += fmt.Sprintf("- Ferramenta: %s\n  Descrição: %s\n", tool.Name(), tool.Description())
	}
	prompt += "\nDepois que uma ferramenta for chamada, eu fornecerei o resultado, e então você deve responder à pergunta original do usuário com base nesse resultado. Se você puder responder diretamente sem ferramentas, faça."
	return prompt
}

// BuildReasoningPrompt cria o prompt de raciocínio que instrui o LLM.
func BuildReasoningPrompt(tools []Tool) string {
	prompt := prompts.ReasoningPrompt + "\n"
	for _, tool := range tools {
		prompt += fmt.Sprintf("- Ferramenta: %s\n  Descrição: %s\n", tool.Name(), tool.Description())
	}
	return prompt
}

// ** ================================================================================================= **
// Implementação do Agente
// =================================================================================================

// Agent é a estrutura principal que orquestra todo o processo.
type Agent struct {
	llmClient     LLMClient
	tools         map[string]Tool
	history       []Message
	toolCallRegex *regexp.Regexp
}

// NewAgent cria uma nova instância do agente.
func NewAgent(client LLMClient, tools []Tool) *Agent {
	toolMap := make(map[string]Tool)
	for _, tool := range tools {
		toolMap[tool.Name()] = tool
	}

	return &Agent{
		llmClient: client,
		tools:     toolMap,
		history:   []Message{},
		// Regex atualizado para ser "ganancioso" e capturar objetos JSON complexos
		toolCallRegex: regexp.MustCompile(`TOOL_CALL:\s*(\w+)\((.*)\)`),
	}
}

// Run inicia o loop de interação principal do agente.
func (a *Agent) Run(ctx context.Context, getUserInput func() (string, bool)) error {
	for {
		fmt.Print("\u001b[94mHumano\u001b[0m: ")
		userInput, ok := getUserInput()
		if !ok {
			break
		}

		a.history = append(a.history, Message{Role: "user", Content: userInput})

		for {
			fmt.Println("\u001b[90mGoAgent está processando a mensagem...\u001b[0m")
			allTools := make([]Tool, 0, len(a.tools))
			for _, t := range a.tools {
				allTools = append(allTools, t)
			}

			llmResponse, err := a.llmClient.GenerateResponse(ctx, a.history, allTools)
			if err != nil {
				fmt.Printf("\u001b[91mErro ao chamar LLM: %v\u001b[0m\n", err)
				break
			}

			matches := a.toolCallRegex.FindStringSubmatch(llmResponse)
			if len(matches) == 3 {
				toolName := matches[1]
				toolArgs := matches[2]

				fmt.Printf("\u001b[95mGoAgent quer usar a ferramenta: %s(%s)\u001b[0m\n", toolName, toolArgs)
				a.history = append(a.history, Message{Role: "assistant", Content: llmResponse})

				tool, ok := a.tools[toolName]
				if !ok {
					fmt.Printf("\u001b[91mErro: Agente tentou usar uma ferramenta desconhecida: %s\u001b[0m\n", toolName)
					a.history = append(a.history, Message{Role: "user", Content: fmt.Sprintf("TOOL_ERROR: Ferramenta '%s' não encontrada.", toolName)})
					continue
				}

				toolResult, err := tool.Execute(toolArgs)
				if err != nil {
					fmt.Printf("\u001b[91mErro ao executar a ferramenta '%s': %v\u001b[0m\n", toolName, err)
					a.history = append(a.history, Message{Role: "user", Content: fmt.Sprintf("TOOL_ERROR: %v", err)})
					continue
				}

				fmt.Printf("\u001b[96mResultado da ferramenta: %s\u001b[0m\n", toolResult)
				a.history = append(a.history, Message{Role: "user", Content: fmt.Sprintf("TOOL_RESULT: %s", toolResult)})
				continue
			}

			fmt.Printf("\u001b[92mGoAgent\u001b[0m: %s\n", llmResponse)
			a.history = append(a.history, Message{Role: "assistant", Content: llmResponse})
			break
		}
	}
	return nil
}

// RunWithReasoning executa o agente "padrão", mas antes insere um raciocínio gerado no histórico.
func (a *Agent) RunWithReasoning(ctx context.Context, getUserInput func() (string, bool)) error {
	for {
		fmt.Print("\u001b[94mHumano\u001b[0m: ") // Garante o mesmo prompt do modo regular
		userInput, ok := getUserInput()
		if !ok {
			break
		}

		// 1. Gera raciocínio usando o helper do ReasoningAgent
		allTools := make([]Tool, 0, len(a.tools))
		for _, t := range a.tools {
			allTools = append(allTools, t)
		}
		reasoning, err := GenerateReasoningTrace(ctx, a.llmClient, userInput, a.history, allTools)
		if err != nil {
			fmt.Printf("\u001b[91mErro ao gerar raciocínio: %v\u001b[0m\n", err)
			continue
		}
		if reasoning != "" {
			fmt.Println("\u001b[96mRaciocínio do agente:\u001b[0m")
			fmt.Println(reasoning)
			// Adiciona o raciocínio ao histórico como mensagem de sistema
			a.history = append(a.history, Message{Role: "system", Content: "Raciocínio para solução:\n" + reasoning})
		}

		// 2. Adiciona a pergunta do usuário
		a.history = append(a.history, Message{Role: "user", Content: userInput})

		// 3. Executa o loop normal do agente
		for {
			fmt.Println("\u001b[90mGoAgent está processando a mensagem...\u001b[0m")
			llmResponse, err := a.llmClient.GenerateResponse(ctx, a.history, allTools)
			if err != nil {
				fmt.Printf("\u001b[91mErro ao chamar LLM: %v\u001b[0m\n", err)
				break
			}

			matches := a.toolCallRegex.FindStringSubmatch(llmResponse)
			if len(matches) == 3 {
				toolName := matches[1]
				toolArgs := matches[2]

				fmt.Printf("\u001b[95mGoAgent quer usar a ferramenta: %s(%s)\u001b[0m\n", toolName, toolArgs)
				a.history = append(a.history, Message{Role: "assistant", Content: llmResponse})

				tool, ok := a.tools[toolName]
				if !ok {
					fmt.Printf("\u001b[91mErro: Agente tentou usar uma ferramenta desconhecida: %s\u001b[0m\n", toolName)
					a.history = append(a.history, Message{Role: "user", Content: fmt.Sprintf("TOOL_ERROR: Ferramenta '%s' não encontrada.", toolName)})
					continue
				}

				toolResult, err := tool.Execute(toolArgs)
				if err != nil {
					fmt.Printf("\u001b[91mErro ao executar a ferramenta '%s': %v\u001b[0m\n", toolName, err)
					a.history = append(a.history, Message{Role: "user", Content: fmt.Sprintf("TOOL_ERROR: %v", err)})
					continue
				}

				fmt.Printf("\u001b[96mResultado da ferramenta: %s\u001b[0m\n", toolResult)
				a.history = append(a.history, Message{Role: "user", Content: fmt.Sprintf("TOOL_RESULT: %s", toolResult)})
				continue
			}

			fmt.Printf("\u001b[92mGoAgent\u001b[0m: %s\n", llmResponse)
			a.history = append(a.history, Message{Role: "assistant", Content: llmResponse})
			break
		}
	}
	return nil
}

// ReasoningConfig configura parâmetros do reasoning
type ReasoningConfig struct {
	MaxTokens     int  // Tokens máximos para reasoning
	ShowTimestamp bool // Mostrar timestamp no reasoning
	DetailLevel   int  // 1=básico, 2=médio, 3=detalhado
}

// DefaultReasoningConfig retorna configuração padrão
func DefaultReasoningConfig() ReasoningConfig {
	return ReasoningConfig{
		MaxTokens:     800,
		ShowTimestamp: true,
		DetailLevel:   2,
	}
}

// GenerateReasoningTrace gera um trace de raciocínio avançado com extração estruturada
func GenerateReasoningTrace(ctx context.Context, llmClient LLMClient, userInput string, history []Message, tools []Tool) (string, error) {
	return GenerateReasoningTraceWithConfig(ctx, llmClient, userInput, history, tools, DefaultReasoningConfig())
}

// GenerateReasoningTraceWithConfig gera trace com configuração customizada
func GenerateReasoningTraceWithConfig(ctx context.Context, llmClient LLMClient, userInput string, history []Message, tools []Tool, config ReasoningConfig) (string, error) {
	reasoningPrompt := BuildReasoningPrompt(tools)
	messages := append([]Message{{Role: "system", Content: reasoningPrompt}}, history...)
	messages = append(messages, Message{Role: "user", Content: userInput})
	
	llmResponse, err := llmClient.GenerateResponse(ctx, messages, tools)
	if err != nil {
		return "", err
	}
	
	// Extrai seções estruturadas do reasoning
	return extractStructuredReasoning(llmResponse, config), nil
}

// extractStructuredReasoning extrai e formata o conteúdo do reasoning
func extractStructuredReasoning(llmResponse string, config ReasoningConfig) string {
	// Regex para extrair conteúdo <think>
	re := regexp.MustCompile(`(?s)<think>(.*?)</think>`)
	matches := re.FindAllStringSubmatch(llmResponse, -1)
	
	if len(matches) == 0 {
		return "❌ Nenhum trace de raciocínio encontrado"
	}
	
	var result strings.Builder
	
	if config.ShowTimestamp {
		result.WriteString(fmt.Sprintf("⏰ Reasoning gerado em: %s\n", time.Now().Format("15:04:05")))
		result.WriteString("─────────────────────────────────────\n")
	}
	
	for i, match := range matches {
		if len(match) > 1 {
			reasoning := strings.TrimSpace(match[1])
			
			// Destaca seções importantes
			reasoning = highlightReasoningSections(reasoning)
			
			if len(matches) > 1 {
				result.WriteString(fmt.Sprintf("🧠 Trace %d:\n", i+1))
			}
			result.WriteString(reasoning)
			result.WriteString("\n")
		}
	}
	
	return result.String()
}

// highlightReasoningSections destaca seções importantes do reasoning
func highlightReasoningSections(reasoning string) string {
	// Destaca emojis e seções estruturadas
	patterns := map[string]string{
		`🎯 OBJETIVO:`:        "🎯 \033[1;33mOBJETIVO:\033[0m",
		`📊 ANÁLISE DO CONTEXTO:`: "📊 \033[1;34mANÁLISE DO CONTEXTO:\033[0m",
		`🛠️ ESTRATÉGIA:`:      "🛠️ \033[1;32mESTRATÉGIA:\033[0m",
		`⚡ MOMENTO AHA!:`:    "⚡ \033[1;31mMOMENTO AHA!:\033[0m",
		`🔍 VALIDAÇÃO:`:       "🔍 \033[1;35mVALIDAÇÃO:\033[0m",
		`🎯 PRÓXIMA AÇÃO:`:    "🎯 \033[1;36mPRÓXIMA AÇÃO:\033[0m",
	}
	
	result := reasoning
	for pattern, replacement := range patterns {
		result = strings.ReplaceAll(result, pattern, replacement)
	}
	
	return result
}

// WithRunWithReasoning retorna um wrapper que implementa Run chamando RunWithReasoning.
func WithRunWithReasoning(a *Agent) interface {
	Run(context.Context, func() (string, bool)) error
} {
	return &runWithReasoningAdapter{a}
}

type runWithReasoningAdapter struct {
	a *Agent
}

func (r *runWithReasoningAdapter) Run(ctx context.Context, getUserInput func() (string, bool)) error {
	return r.a.RunWithReasoning(ctx, getUserInput)
}
