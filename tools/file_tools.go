// tools/file_tools.go
package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// --- Ferramenta: ListFiles ---

// ListFilesInput define os parâmetros para a função ListFiles.
type ListFilesInput struct {
	Path string `json:"path,omitempty"`
}

// listFiles é a função lógica que varre um diretório.
func listFiles(input json.RawMessage) (string, error) {
	var typedInput ListFilesInput
	if len(input) > 0 && string(input) != "null" {
		if err := json.Unmarshal(input, &typedInput); err != nil {
			return "", fmt.Errorf("JSON inválido para argumentos: %w", err)
		}
	}

	dir := "." // Valor padrão: diretório atual
	if typedInput.Path != "" {
		dir = typedInput.Path
	}

	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path != dir {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	result, err := json.Marshal(files)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

// ListFilesDef é a definição pública da nossa ferramenta.
var ListFilesDef = ToolDefinition{
	Name:        "list_files",
	Description: "Lista arquivos e diretórios em um caminho específico. Se nenhum caminho for fornecido, lista o conteúdo do diretório atual.",
	Function:    listFiles,
}

// --- Ferramenta: WriteFile ---

type WriteFileInput struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

func writeFile(input json.RawMessage) (string, error) {
	var typedInput WriteFileInput
	if err := json.Unmarshal(input, &typedInput); err != nil {
		return "", fmt.Errorf("JSON inválido para argumentos: %w", err)
	}

	if typedInput.Path == "" || typedInput.Content == "" {
		return "", fmt.Errorf("argumentos inválidos. 'path' e 'content' são obrigatórios")
	}

	err := os.WriteFile(typedInput.Path, []byte(typedInput.Content), 0644)
	if err != nil {
		return "", fmt.Errorf("erro ao escrever no arquivo '%s': %w", typedInput.Path, err)
	}
	return fmt.Sprintf("Arquivo '%s' escrito com sucesso.", typedInput.Path), nil
}

var WriteFileDef = ToolDefinition{
	Name:        "write_file",
	Description: `Escreve o conteúdo fornecido em um arquivo. Requer um objeto JSON com as chaves "path" e "content". Exemplo: {"path": "caminho/arquivo.txt", "content": "Olá, mundo!"}`,
	Function:    writeFile,
}

// --- Ferramenta: ReadFile ---

type ReadFileInput struct {
	Path string `json:"path"`
}

func readFile(input json.RawMessage) (string, error) {
	var typedInput ReadFileInput
	if err := json.Unmarshal(input, &typedInput); err != nil {
		return "", fmt.Errorf("JSON inválido para argumentos: %w", err)
	}

	if typedInput.Path == "" {
		return "", fmt.Errorf("argumento inválido. 'path' é obrigatório")
	}

	content, err := os.ReadFile(typedInput.Path)
	if err != nil {
		return "", fmt.Errorf("erro ao ler o arquivo '%s': %w", typedInput.Path, err)
	}
	return string(content), nil
}

var ReadFileDef = ToolDefinition{
	Name:        "read_file",
	Description: `Lê o conteúdo de um arquivo. Requer um objeto JSON com a chave "path". Exemplo: {"path": "caminho/arquivo.txt"}`,
	Function:    readFile,
}
