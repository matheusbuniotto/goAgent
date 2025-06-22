package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/matheusbuniotto/goagent/agent"
	"github.com/matheusbuniotto/goagent/tools"
)

func main() {
	// Carrega chaves de API de variáveis de ambiente.
	openaiAPIKey := os.Getenv("OPENAI_API_KEY")
	geminiAPIKey := os.Getenv("GEMINI_API_KEY")

	if openaiAPIKey == "" && geminiAPIKey == "" {
		fmt.Println("\u001b[91mErro: Nenhuma chave de API encontrada. Por favor, defina OPENAI_API_KEY ou GEMINI_API_KEY.\u001b[0m")
		os.Exit(1)
	}

	// Instância de ferramentas que o agente poderá usar, vindas do pacte tools
	allTools := []agent.Tool{
		&tools.ToolAdapter{Definition: tools.ListFilesDef},
		&tools.ToolAdapter{Definition: tools.WriteFileDef},
		&tools.ToolAdapter{Definition: tools.ReadFileDef},
	}
	// --- Escolha qual LLM usar ---
	var llmClient agent.LLMClient

	// Dê prioridade ao Gemini se ambas as chaves estiverem definidas
	if geminiAPIKey != "" {
		fmt.Println("\u001b[92mUsando cliente Google Gemini.\u001b[0m")
		llmClient = agent.NewGeminiClient(geminiAPIKey)
	} else if openaiAPIKey != "" {
		fmt.Println("\u001b[92mUsando cliente OpenAI.\u001b[0m")
		llmClient = agent.NewOpenAIClient(openaiAPIKey)
	}

	// Inicializa o agente do pacote agent
	theAgent := agent.NewAgent(llmClient, allTools)

	// Prepara a função para ler a entrada do terminal
	scanner := bufio.NewScanner(os.Stdin)
	getUserInput := func() (string, bool) {
		if !scanner.Scan() {
			return "", false
		}
		return scanner.Text(), true
	}

	// Executa o agente no terminal
	fmt.Println("\u001b[92mChat com GoAgent ('ctrl-c' para sair)\u001b[0m")
	err := theAgent.Run(context.Background(), getUserInput)
	if err != nil {
		fmt.Printf("\u001b[91mErro fatal do agente: %s\u001b[0m\n", err.Error())
	}
}
