package toolkit

import "encoding/json"

// ToolFunction define a estrutura da função principal de uma ferramenta
// Ela recebe os argumentos como um JSON "cru" e retorna o resultado ou um erro
type ToolFunction func(input json.RawMessage) (string, error)

// ToolDefinition é a estruturada de definir uma ferramenta.
type ToolDefinition struct {
	Name        string       // Nome da ferramenta
	Description string       // Descrição da ferramenta, usada para informar o agente
	Function    ToolFunction // A função que implementa a lógica da ferramenta
}
