// tools/adapter.go
package tools

import "encoding/json"

// ToolAdapter faz a "ponte" entre o novo padrão (ToolDefinition) e a interface
// agent.Tool que seu `main.go` espera. Isso permite a transição sem quebrar o código existente.
type ToolAdapter struct {
	Definition ToolDefinition
}

// Name retorna o nome da ferramenta. Parte da interface agent.Tool.
func (a *ToolAdapter) Name() string {
	return a.Definition.Name
}

// Description retorna a descrição da ferramenta. Parte da interface agent.Tool.
func (a *ToolAdapter) Description() string {
	return a.Definition.Description
}

// Execute é a chave do adaptador. Ele pega a string de argumentos do agente,
// a trata como JSON e a passa para a função real da nossa ToolDefinition.
func (a *ToolAdapter) Execute(args string) (string, error) {
	// O LLM agora deve ser instruído a fornecer argumentos como um JSON.
	// Ex: `{"path": "meu_dir", "content": "olá"}` em vez de "meu_dir,olá"
	rawJSON := json.RawMessage(args)
	return a.Definition.Function(rawJSON)
}
