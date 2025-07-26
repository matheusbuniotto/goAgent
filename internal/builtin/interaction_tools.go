package builtin

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/matheusbuniotto/goagent/pkg/toolkit"
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
	fmt.Printf("\u001b[95mMe responda\u001b[0m: %s\n ", typedInput.Question)

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

var AskHumanDef = toolkit.ToolDefinition{
	Name:        "ask_human_for_clarification",
	Description: `Quando necessário,pede ajuda para tirar dúvidas **CRÍTICAS**. Requer um objeto JSON com a chaves "question", contendo sua pergunta. Exemplo: {"question": "Você pode me responder...?"}`,
	Function:    askHuman,
}
