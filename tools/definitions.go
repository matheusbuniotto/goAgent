// tools/definitions.go
package tools

import "encoding/json"

// ToolFunction define a assinatura da função principal de uma ferramenta.
// Ela recebe os argumentos como um JSON "cru" e retorna o resultado ou um erro.
type ToolFunction func(input json.RawMessage) (string, error)

// ToolDefinition é a nova forma estruturada de definir uma ferramenta.
// É mais descritiva e robusta do que a implementação anterior baseada em interface.
type ToolDefinition struct {
	Name        string
	Description string
	Function    ToolFunction
}
