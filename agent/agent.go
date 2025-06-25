// Package agent contém a lógica principal para o agente de IA,
// incluindo as interfaces, o loop de execução e os clientes LLM.
package agent

import (
	"context"
	"fmt"
	"regexp"
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
	prompt := `Você é GoAgent, um assistente que pode usar ferramentas para interagir com o sistema do usuário.
		Para usar uma ferramenta, você **DEVE responder EXATAMENTE** no seguinte formato: TOOL_CALL: ToolName({"arg_name": "value", "another_arg": "value"})
		### IMPORTANTE: Os argumentos da ferramenta **DEVEM ser um objeto JSON válido**.
		Se uma ferramenta não requer argumentos, use um objeto JSON vazio: TOOL_CALL: ToolName({})
		As ferramentas disponíveis estão listadas abaixo com sua descrição:
`

	for _, tool := range tools {
		prompt += fmt.Sprintf("- Ferramenta: %s\n  Descrição: %s\n", tool.Name(), tool.Description())
	}

	prompt += "\nDepois que uma ferramenta for chamada, eu fornecerei o resultado, e então você deve responder à pergunta original do usuário com base nesse resultado. Se você puder responder diretamente sem ferramentas, faça."
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
