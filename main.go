package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"flag"
	"log"

	"github.com/matheusbuniotto/goagent/agent"
	"github.com/matheusbuniotto/goagent/llmclients"
	"github.com/matheusbuniotto/goagent/tools"
)

func main() {
	// Carrega chaves de API de variáveis de ambiente.
	openaiAPIKey := os.Getenv("OPENAI_API_KEY")
	geminiAPIKey := os.Getenv("GEMINI_API_KEY")

	// O valor padrão "" significa que o flag não foi usado.
	model := flag.String("model", "", "O provedor a ser usado para o agente (gemini ou openai). Sobrepõe a detecção automática.")
	// Novo: permite selecionar o ReasoningAgent
	agentType := flag.String("agent", "default", "Tipo de agente: default ou reasoning")
	flag.Parse()

	var llmClient llmclients.LLMClient

	// A prioridade é o flag da linha de comando -model
	switch *model {
	case "gemini":
		if geminiAPIKey == "" {
			log.Fatal("\u001b[91mErro: Modelo 'gemini' especificado, mas a chave GEMINI_API_KEY não foi encontrada.\u001b[0m")
		}
		fmt.Println("\u001b[92mUsando cliente Google Gemini (especificado via flag).\u001b[0m")
		llmClient = llmclients.NewGeminiClient(geminiAPIKey)

	case "openai":
		if openaiAPIKey == "" {
			log.Fatal("\u001b[91mErro: Modelo 'openai' especificado, mas a chave OPENAI_API_KEY não foi encontrada.\u001b[0m")
		}
		fmt.Println("\u001b[92mUsando cliente OpenAI (especificado via flag).\u001b[0m")
		llmClient = llmclients.NewOpenAIClient(openaiAPIKey)

	case "":
		// Se NENHUM flag for passado, usa gemini por padrão a não ser que não esteja definido.
		fmt.Println("\u001b[92mNenhum modelo especificado, detectando automaticamente por chave de API...\u001b[0m")
		if geminiAPIKey != "" {
			fmt.Println("\u001b[92mUsando cliente Google Gemini.\u001b[0m")
			llmClient = llmclients.NewGeminiClient(geminiAPIKey)
		} else if openaiAPIKey != "" {
			fmt.Println("\u001b[92mUsando cliente OpenAI.\u001b[0m")
			llmClient = llmclients.NewOpenAIClient(openaiAPIKey)
		} else {
			log.Fatal("\u001b[91mErro: Nenhuma chave de API encontrada. Por favor, defina OPENAI_API_KEY ou GEMINI_API_KEY.\u001b[0m")
		}

	default:
		log.Fatalf("\u001b[91mErro: Modelo desconhecido '%s' especificado. Use 'gemini' ou 'openai'.\u001b[0m", *model)
	}

	// Instância de ferramentas que o agente poderá usar, vindas do pacte tools
	allTools := []agent.Tool{
		&tools.ToolAdapter{Definition: tools.ListFilesDef},
		&tools.ToolAdapter{Definition: tools.WriteFileDef},
		&tools.ToolAdapter{Definition: tools.ReadFileDef},
		&tools.ToolAdapter{Definition: tools.CreateDirectoryDef},
		&tools.ToolAdapter{Definition: tools.AskHumanDef},
	}

	// Inicializa o agente correto
	var theAgent interface {
		Run(context.Context, func() (string, bool)) error
	}
	if *agentType == "reasoning" || *agentType == "r" {
		theAgent = agent.WithRunWithReasoning(agent.NewAgent(llmClient, allTools))
		fmt.Println("\u001b[92mModo Reasoning ativado.\u001b[0m")
	} else {
		theAgent = agent.NewAgent(llmClient, allTools)
	}

	// Prepara a função para ler o input
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
