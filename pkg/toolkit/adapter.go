// pkg/toolkit/adapter.go
package toolkit

import "encoding/json"

// ToolAdapter faz a "ponte" entre o ToolDefinition e a interface
type ToolAdapter struct {
	Definition ToolDefinition
}

// Name retorna o nome da ferramenta. Parte da interface agent.Tool.
func (a *ToolAdapter) Name() string {
	return a.Definition.Name
}

// Description retorna a descrição da ferramenta. Parte da interface  do agent.Tool.
func (a *ToolAdapter) Description() string {
	return a.Definition.Description
}

// Execute é a chave do adaptador. Ele pega a string de argumentos do agente,
// a trata como JSON e a passa para a função real da nossa ToolDefinition.
func (a *ToolAdapter) Execute(args string) (string, error) {
	// O LLM é instruído a fornecer argumentos como um JSON.
	// Ex: `{"path": "meu_dir", "content": "olá"}` em vez de "meu_dir,olá"
	rawJSON := json.RawMessage(args)
	return a.Definition.Function(rawJSON)
}
