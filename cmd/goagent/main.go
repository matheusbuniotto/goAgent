package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"flag"
	"log"

	"github.com/matheusbuniotto/goagent/pkg/agent"
	"github.com/matheusbuniotto/goagent/internal/llm"
	"github.com/matheusbuniotto/goagent/pkg/toolkit"
	"github.com/matheusbuniotto/goagent/internal/builtin"
)

// selectProvider permite ao usuário escolher um provedor interativamente
func selectProvider() string {
	fmt.Println("\n🤖 Selecione um provedor de LLM:")
	fmt.Println("─────────────────────────────────────")
	fmt.Println("1. OpenRouter (Acesso a múltiplos modelos)")
	fmt.Println("2. Gemini (Google)")
	fmt.Println("3. OpenAI")
	
	fmt.Print("\nDigite o número do provedor desejado (1-3): ")
	
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if !scanner.Scan() {
			fmt.Println("Erro ao ler entrada. Usando auto-detecção.")
			return "auto"
		}
		
		input := strings.TrimSpace(scanner.Text())
		switch input {
		case "1":
			fmt.Println("✅ OpenRouter selecionado")
			return "openrouter"
		case "2":
			fmt.Println("✅ Gemini selecionado")
			return "gemini"
		case "3":
			fmt.Println("✅ OpenAI selecionado")
			return "openai"
		default:
			fmt.Print("Opção inválida. Digite 1, 2 ou 3: ")
		}
	}
}

func main() {
	// Carrega chaves de API de variáveis de ambiente.
	openaiAPIKey := os.Getenv("OPENAI_API_KEY")
	geminiAPIKey := os.Getenv("GEMINI_API_KEY")
	openrouterAPIKey := os.Getenv("OPENROUTER_API_KEY")

	// Flag para escolher provedor ou usar menu interativo
	interactiveMode := flag.Bool("select", false, "Modo interativo para escolher provedor")
	model := flag.String("model", "", "O provedor a ser usado para o agente (gemini, openai ou openrouter). Sobrepõe a detecção automática.")
	// Novo: permite selecionar o ReasoningAgent
	agentType := flag.String("agent", "default", "Tipo de agente: default ou reasoning")
	// Configurações de reasoning
	reasoningDetail := flag.Int("reasoning-detail", 2, "Nível de detalhe do reasoning (1=básico, 2=médio, 3=detalhado)")
	reasoningTimestamp := flag.Bool("reasoning-timestamp", true, "Mostrar timestamp no reasoning")
	flag.Parse()

	var llmClient llm.LLMClient
	var selectedProvider string

	// Determina o provedor a ser usado
	if *interactiveMode {
		selectedProvider = selectProvider()
	} else if *model != "" {
		selectedProvider = *model
	} else {
		// Auto-detecção por chave de API (prioridade: OpenRouter > Gemini > OpenAI)
		fmt.Println("\u001b[92mNenhum provedor especificado, detectando automaticamente por chave de API...\u001b[0m")
		if openrouterAPIKey != "" {
			selectedProvider = "openrouter"
		} else if geminiAPIKey != "" {
			selectedProvider = "gemini"
		} else if openaiAPIKey != "" {
			selectedProvider = "openai"
		} else {
			log.Fatal("\u001b[91mErro: Nenhuma chave de API encontrada. Por favor, defina OPENROUTER_API_KEY, GEMINI_API_KEY ou OPENAI_API_KEY.\u001b[0m")
		}
	}

	// Cria o cliente baseado no provedor selecionado
	switch selectedProvider {
	case "gemini":
		if geminiAPIKey == "" {
			log.Fatal("\u001b[91mErro: Gemini selecionado, mas a chave GEMINI_API_KEY não foi encontrada.\u001b[0m")
		}
		fmt.Println("\u001b[92m✅ Usando cliente Google Gemini\u001b[0m")
		llmClient = llm.NewGeminiClient(geminiAPIKey)

	case "openai":
		if openaiAPIKey == "" {
			log.Fatal("\u001b[91mErro: OpenAI selecionado, mas a chave OPENAI_API_KEY não foi encontrada.\u001b[0m")
		}
		fmt.Println("\u001b[92m✅ Usando cliente OpenAI\u001b[0m")
		llmClient = llm.NewOpenAIClient(openaiAPIKey)

	case "openrouter":
		if openrouterAPIKey == "" {
			log.Fatal("\u001b[91mErro: OpenRouter selecionado, mas a chave OPENROUTER_API_KEY não foi encontrada.\u001b[0m")
		}
		fmt.Println("\u001b[92m✅ Usando cliente OpenRouter\u001b[0m")
		// Se OpenRouter for escolhido, sempre pergunta qual modelo usar
		selectedModel := llm.SelectOpenRouterModel()
		llmClient = llm.NewOpenRouterClientWithModel(openrouterAPIKey, selectedModel)

	case "auto":
		// Fallback para auto-detecção se seleção interativa falhou
		if openrouterAPIKey != "" {
			fmt.Println("\u001b[92m✅ Usando cliente OpenRouter (auto-detectado)\u001b[0m")
			// Quando auto-detectado, também permite escolher o modelo
			selectedModel := llm.SelectOpenRouterModel()
			llmClient = llm.NewOpenRouterClientWithModel(openrouterAPIKey, selectedModel)
		} else if geminiAPIKey != "" {
			fmt.Println("\u001b[92m✅ Usando cliente Google Gemini (auto-detectado)\u001b[0m")
			llmClient = llm.NewGeminiClient(geminiAPIKey)
		} else if openaiAPIKey != "" {
			fmt.Println("\u001b[92m✅ Usando cliente OpenAI (auto-detectado)\u001b[0m")
			llmClient = llm.NewOpenAIClient(openaiAPIKey)
		} else {
			log.Fatal("\u001b[91mErro: Nenhuma chave de API encontrada.\u001b[0m")
		}

	default:
		log.Fatalf("\u001b[91mErro: Provedor desconhecido '%s'.\u001b[0m", selectedProvider)
	}

	// Instância de ferramentas que o agente poderá usar, vindas do pacte builtin
	allTools := []agent.Tool{
		&toolkit.ToolAdapter{Definition: builtin.ListFilesDef},
		&toolkit.ToolAdapter{Definition: builtin.WriteFileDef},
		&toolkit.ToolAdapter{Definition: builtin.ReadFileDef},
		&toolkit.ToolAdapter{Definition: builtin.CreateDirectoryDef},
			&toolkit.ToolAdapter{Definition: builtin.GrepSearchDef},
		&toolkit.ToolAdapter{Definition: builtin.AskHumanDef},
		&toolkit.ToolAdapter{Definition: builtin.AnalyzeReasoningDef},
		&toolkit.ToolAdapter{Definition: builtin.ReviewDecisionDef},
	}

	// Inicializa o agente correto
	var theAgent interface {
		Run(context.Context, func() (string, bool)) error
	}
	if *agentType == "reasoning" || *agentType == "r" {
		theAgent = agent.WithRunWithReasoning(agent.NewAgent(llmClient, allTools))
		fmt.Printf("\u001b[92mModo Reasoning ativado (detalhe: %d, timestamp: %v).\u001b[0m\n", *reasoningDetail, *reasoningTimestamp)
	} else {
		theAgent = agent.NewAgent(llmClient, allTools)
		// Suprime warning das variáveis não usadas no modo default
		_ = *reasoningDetail
		_ = *reasoningTimestamp
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
