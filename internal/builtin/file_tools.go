package builtin

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/matheusbuniotto/goagent/pkg/toolkit"
)

// ::: Ferramenta: ListFiles :::

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
var ListFilesDef = toolkit.ToolDefinition{
	Name:        "list_files",
	Description: "Lista arquivos e diretórios em um caminho específico. Se nenhum caminho for fornecido, lista o conteúdo do diretório atual.",
	Function:    listFiles,
}

// ::: Ferramenta: WriteFile :::

type WriteFileInput struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

func writeFile(input json.RawMessage) (string, error) {
	var typedInput WriteFileInput
	if err := json.Unmarshal(input, &typedInput); err != nil {
		return "", fmt.Errorf("JSON inválido para argumentos: %w", err)
	}

	if typedInput.Path == "" {
		return "", fmt.Errorf("argumentos inválidos. 'path' e 'content' são obrigatórios")
	}

	err := os.WriteFile(typedInput.Path, []byte(typedInput.Content), 0644)
	if err != nil {
		return "", fmt.Errorf("erro ao escrever no arquivo '%s': %w", typedInput.Path, err)
	}
	return fmt.Sprintf("Arquivo '%s' escrito com sucesso.", typedInput.Path), nil
}

var WriteFileDef = toolkit.ToolDefinition{
	Name:        "write_file",
	Description: `Escreve o conteúdo fornecido em um arquivo. Requer um objeto JSON com as chaves "path" e "content". Exemplo: {"path": "caminho/arquivo.txt", "content": "Olá, mundo!"}`,
	Function:    writeFile,
}

// :::: Ferramenta: ReadFile :::

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

var ReadFileDef = toolkit.ToolDefinition{
	Name:        "read_file",
	Description: `Lê o conteúdo de um arquivo. Requer um objeto JSON com a chave "path". Exemplo: {"path": "caminho/arquivo.txt"}`,
	Function:    readFile,
}

// ::: Ferramenta: CreateDirectory :::
type CreateDirectoryInput struct {
	Path string `json:"path"`
}

func createDirectory(input json.RawMessage) (string, error) {
	var typedInput CreateDirectoryInput
	if err := json.Unmarshal(input, &typedInput); err != nil {
		return "", fmt.Errorf("JSON inválido para argumentos: %w", err)
	}

	if typedInput.Path == "" {
		return "", fmt.Errorf("argumento inválido. 'path' é obrigatório")
	}

	// 0755 são as permissões = leitura/execução para todos, escrita para o dono
	err := os.MkdirAll(typedInput.Path, 0755)
	if err != nil {
		return "", fmt.Errorf("erro ao criar o diretório '%s': %w", typedInput.Path, err)
	}

	return fmt.Sprintf("Diretório '%s' criado com sucesso.", typedInput.Path), nil
}

var CreateDirectoryDef = toolkit.ToolDefinition{
	Name:        "create_directory",
	Description: `Cria um novo diretório no caminho especificado, necessita de um nome. Exemplo: {"path": "meu/novo/nome_diretorio"}`,
	Function:    createDirectory,
}

// ::: Ferramenta: GrepSearch :::

// GrepSearchInput define os parâmetros para a função grepSearch.
type GrepSearchInput struct {
	Pattern string `json:"pattern"`
	Path    string `json:"path,omitempty"`
}

// grepMatch representa uma linha encontrada pelo grep.
type grepMatch struct {
	File    string `json:"file"`
	Line    int    `json:"line"`
	Content string `json:"content"`
}

func grepSearch(input json.RawMessage) (string, error) {
	var typedInput GrepSearchInput
	if err := json.Unmarshal(input, &typedInput); err != nil {
		return "", fmt.Errorf("JSON inválido para argumentos: %w", err)
	}

	if typedInput.Pattern == "" {
		return "", fmt.Errorf("argumento inválido. 'pattern' é obrigatório")
	}

	dir := "."
	if typedInput.Path != "" {
		dir = typedInput.Path
	}

	var matches []grepMatch
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		// Ignora diretórios ocultos e arquivos ocultos
		if strings.Contains(path, string(filepath.Separator)+".") || strings.HasPrefix(filepath.Base(path), ".") {
			return nil
		}

		content, readErr := os.ReadFile(path)
		if readErr != nil {
			// Ignora arquivos que não podem ser lidos (ex: binários, permissão negada)
			return nil
		}

		lines := strings.Split(string(content), "\n")
		for lineNum, line := range lines {
			if strings.Contains(line, typedInput.Pattern) {
				matches = append(matches, grepMatch{
					File:    path,
					Line:    lineNum + 1,
					Content: line,
				})
			}
		}
		return nil
	})

	if err != nil {
		return "", fmt.Errorf("erro ao buscar no diretório '%s': %w", dir, err)
	}

	if len(matches) == 0 {
		return "Nenhum resultado encontrado.", nil
	}

	result, err := json.Marshal(matches)
	if err != nil {
		return "", fmt.Errorf("erro ao serializar resultados: %w", err)
	}

	return string(result), nil
}

// GrepSearchDef é a definição pública da ferramenta de busca grep.
var GrepSearchDef = toolkit.ToolDefinition{
	Name:        "grep_search",
	Description: `Busca por um padrão de texto em arquivos dentro de um diretório. Requer um objeto JSON com a chave "pattern" e opcionalmente "path". Exemplo: {"pattern": "minha_funcao", "path": "src/"}`,
	Function:    grepSearch,
}
