// Package tools contém implementações de ferramentas que o agente pode usar.
package tools

import (
	"fmt"
	"os"
	"strings"
)

// WriteFileTool é uma ferramenta para escrever texto em um arquivo.
type WriteFileTool struct{}

// Name retorna o nome da ferramenta.
func (t *WriteFileTool) Name() string {
	return "WriteFile"
}

// Description descreve o que a ferramenta faz.
func (t *WriteFileTool) Description() string {
	return `Escreve o conteúdo fornecido em um arquivo. Aceita dois argumentos separados por vírgula: o caminho do arquivo e o conteúdo a ser escrito. Exemplo de uso: WriteFile(caminho/do/arquivo.txt, "Olá, mundo!")`
}

// Execute executa a lógica da ferramenta.
func (t *WriteFileTool) Execute(args string) (string, error) {
	parts := strings.SplitN(args, ",", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("argumentos inválidos. Esperado: 'caminho,conteúdo'. Recebido: '%s'", args)
	}
	filePath := strings.TrimSpace(parts[0])
	content := strings.TrimSpace(parts[1])

	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return "", fmt.Errorf("erro ao escrever no arquivo '%s': %w", filePath, err)
	}
	return fmt.Sprintf("Arquivo '%s' escrito com sucesso.", filePath), nil
}

// ReadFileTool é uma ferramenta para ler o conteúdo de um arquivo.
type ReadFileTool struct{}

// Name retorna o nome da ferramenta.
func (t *ReadFileTool) Name() string {
	return "ReadFile"
}

// Description descreve o que a ferramenta faz.
func (t *ReadFileTool) Description() string {
	return "Lê e retorna o conteúdo de um arquivo. Aceita um único argumento: o caminho do arquivo. Exemplo de uso: ReadFile(caminho/do/arquivo.txt)"
}

// Execute executa a lógica da ferramenta.
func (t *ReadFileTool) Execute(args string) (string, error) {
	filePath := strings.TrimSpace(args)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("erro ao ler o arquivo '%s': %w", err)
	}
	return string(content), nil
}
