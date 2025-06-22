package tools

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// ::: Ferramenta: Human in the loop :::

type AskHumanInput struct {
	Question string `json:"question"`
}

func askHuman(input json.RawMessage) (string, error) {
	var typedInput AskHumanInput
	if err := json.Unmarshal(input, &typedInput); err != nil {
		return "", fmt.Errorf("JSON inválido para argumentos: %w", err)
	}

	if typedInput.Question == "" {
		return "", fmt.Errorf("argumento inválido. 'question' é obrigatório")
	}

	// Traza pergunta do agente
	fmt.Printf("\n Chefe, preciso de ajuda: %s\n> ", typedInput.Question)

	// Cria um leitor para o input padrão, por enquanto cli.
	reader := bufio.NewReader(os.Stdin)
	// Lê a resposta humana até que pressione Enter
	response, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("erro ao ler a resposta do humano: %w", err)
	}

	// Remove espaços em branco e a quebra de linha da resposta
	return strings.TrimSpace(response), nil
}

var AskHumanDef = ToolDefinition{
	Name: "ask_human_for_clarification",
	Description: `Pausa a execução e faz uma pergunta ao usuário humano para obter mais informações ou esclarecimentos. Use quando os requisitos não forem claros ou quando julgar ser uma ação crítica. 
		Exemplo: {"question": "Qual nome você gostaria de dar ao  ou pasta?"}
		Exemplo 2: {question: "Posso prosseguir com essa alteração crítica?"}`,
	Function: askHuman,
}
