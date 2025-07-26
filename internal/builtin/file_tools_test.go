package builtin

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

// TestCreateDirectory
// *testes orientados por tabela
func TestCreateDirectory(t *testing.T) {
	// Define casos de teste em uma tabela
	testCases := []struct {
		name         string // Nome do caso de teste.
		inputPath    string // O caminho
		expectErr    bool   // Indica se esperamos um erro
		setupTempDir bool   // Indica se op teste precisa de um temp dir
	}{
		{
			name:         "Sucesso - Cria um novo diretório",
			inputPath:    "meu-diretorio-de-teste",
			expectErr:    false,
			setupTempDir: true,
		},
		{
			name:         "Erro - Caminho vazio",
			inputPath:    "",
			expectErr:    true,
			setupTempDir: false,
		},
	}

	// Iterando sobre cada caso de teste.
	for _, tc := range testCases {
		// t.Run nos permite agrupar a lógica de cada teste e dar um nome a ele.
		t.Run(tc.name, func(t *testing.T) {
			// Prepara o caminho completo do diretório, se necessário.
			fullPath := tc.inputPath
			if tc.setupTempDir {
				tempDir := t.TempDir()
				fullPath = filepath.Join(tempDir, tc.inputPath)
			}

			// Monta o input JSON para a nossa função.
			inputData := CreateDirectoryInput{Path: fullPath}
			rawInput, _ := json.Marshal(inputData)

			// Executa a função que estamos testando.
			_, err := createDirectory(rawInput)

			// Verifica o resultado do erro.
			if (err != nil) != tc.expectErr {
				t.Fatalf("createDirectory() erro = %v, expectErr %v", err, tc.expectErr)
			}

			// Se não esperávamos um erro, fazemos uma verificação extra:
			// o diretório realmente existe no sistema de arquivos?
			if !tc.expectErr {
				if _, err := os.Stat(fullPath); os.IsNotExist(err) {
					t.Errorf("createDirectory() não criou o diretório esperado em %s", fullPath)
				}
			}
		})
	}
}

// TestListFiles
func TestListFiles(t *testing.T) {
	// Setup: Criar uma estrutura de diretórios e arquivos conhecida.
	tempDir := t.TempDir()

	// Cria uma sub pasta
	subDir := filepath.Join(tempDir, "subdir")
	_ = os.Mkdir(subDir, 0755)

	// Cria alguns arquivos.
	_ = os.WriteFile(filepath.Join(tempDir, "testando_som.txt"), []byte("a"), 0644)
	_ = os.WriteFile(filepath.Join(subDir, "testando_som_12.txt"), []byte("b"), 0644)

	// Os caminhos que esperamos que sejam listados.
	// filepath.Join para garantir que funcione em qualquer so
	expectedFiles := []string{
		filepath.Join(tempDir, "testando_som.txt"),
		filepath.Join(tempDir, "subdir"),
		filepath.Join(tempDir, "subdir", "testando_som_12.txt"),
	}

	// Execução: Chamar a função a ser testada.
	inputData := ListFilesInput{Path: tempDir}
	rawInput, _ := json.Marshal(inputData)

	resultJSON, err := listFiles(rawInput)

	// Verificação: Checar se o resultado está correto.
	if err != nil {
		t.Fatalf("listFiles() retornou um erro inesperado: %v", err)
	}

	var actualFiles []string
	if err := json.Unmarshal([]byte(resultJSON), &actualFiles); err != nil {
		t.Fatalf("Falha ao decodificar o resultado JSON: %v", err)
	}

	// Sort para ordernar os slices para comparar.
	sort.Strings(expectedFiles)
	sort.Strings(actualFiles)

	// reflect.DeepEqual é a forma idiomática para comparar slices, maps e structs.
	if !reflect.DeepEqual(expectedFiles, actualFiles) {
		t.Errorf("listFiles() resultado incorreto.\nEsperado: %v\nRecebido: %v", expectedFiles, actualFiles)
	}
}

// TestReadFile
func TestReadFile(t *testing.T) {
	testCases := []struct {
		name          string
		fileContent   string
		fileName      string
		expectErr     bool
		expectedError string
	}{
		{
			name:        "Sucesso - Lê arquivo existente",
			fileContent: "Conteúdo de teste\nSegunda linha",
			fileName:    "teste.txt",
			expectErr:   false,
		},
		{
			name:        "Sucesso - Arquivo vazio",
			fileContent: "",
			fileName:    "vazio.txt",
			expectErr:   false,
		},
		{
			name:          "Erro - Arquivo inexistente",
			fileName:      "inexistente.txt",
			expectErr:     true,
			expectedError: "no such file or directory",
		},
		{
			name:          "Erro - Caminho vazio",
			fileName:      "",
			expectErr:     true,
			expectedError: "argumento inválido",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tempDir := t.TempDir()
			var fullPath string

			if tc.fileName != "" {
				fullPath = filepath.Join(tempDir, tc.fileName)
				
				// Cria o arquivo se não for teste de arquivo inexistente
				if !tc.expectErr || tc.expectedError != "no such file or directory" {
					err := os.WriteFile(fullPath, []byte(tc.fileContent), 0644)
					if err != nil {
						t.Fatalf("Falha ao criar arquivo de teste: %v", err)
					}
				}
			}

			// Prepara input JSON
			inputData := ReadFileInput{Path: fullPath}
			rawInput, _ := json.Marshal(inputData)

			// Executa a função
			result, err := readFile(rawInput)

			// Verifica erro
			if (err != nil) != tc.expectErr {
				t.Fatalf("readFile() erro = %v, expectErr %v", err, tc.expectErr)
			}

			if tc.expectErr && tc.expectedError != "" {
				if err == nil || !contains(err.Error(), tc.expectedError) {
					t.Errorf("readFile() erro esperado contendo '%s', got '%v'", tc.expectedError, err)
				}
			}

			// Verifica conteúdo se sucesso
			if !tc.expectErr {
				if result != tc.fileContent {
					t.Errorf("readFile() conteúdo = %q, esperado %q", result, tc.fileContent)
				}
			}
		})
	}
}

// TestWriteFile
func TestWriteFile(t *testing.T) {
	testCases := []struct {
		name          string
		fileName      string
		content       string
		expectErr     bool
		expectedError string
	}{
		{
			name:      "Sucesso - Escreve arquivo novo",
			fileName:  "novo.txt",
			content:   "Conteúdo novo",
			expectErr: false,
		},
		{
			name:      "Sucesso - Sobrescreve arquivo existente",
			fileName:  "existente.txt",
			content:   "Conteúdo atualizado",
			expectErr: false,
		},
		{
			name:      "Sucesso - Conteúdo vazio",
			fileName:  "vazio.txt",
			content:   "",
			expectErr: false,
		},
		{
			name:          "Erro - Nome de arquivo vazio",
			fileName:      "",
			content:       "conteúdo",
			expectErr:     true,
			expectedError: "argumentos inválidos",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tempDir := t.TempDir()
			var fullPath string

			if tc.fileName != "" {
				fullPath = filepath.Join(tempDir, tc.fileName)
				
				// Se for teste de sobrescrita, cria arquivo existente
				if tc.name == "Sucesso - Sobrescreve arquivo existente" {
					err := os.WriteFile(fullPath, []byte("conteúdo antigo"), 0644)
					if err != nil {
						t.Fatalf("Falha ao criar arquivo existente: %v", err)
					}
				}
			}

			// Prepara input JSON
			inputData := WriteFileInput{Path: fullPath, Content: tc.content}
			rawInput, _ := json.Marshal(inputData)

			// Executa a função
			result, err := writeFile(rawInput)

			// Verifica erro
			if (err != nil) != tc.expectErr {
				t.Fatalf("writeFile() erro = %v, expectErr %v", err, tc.expectErr)
			}

			if tc.expectErr && tc.expectedError != "" {
				if err == nil || !contains(err.Error(), tc.expectedError) {
					t.Errorf("writeFile() erro esperado contendo '%s', got '%v'", tc.expectedError, err)
				}
				return
			}

			// Verifica se foi criado com sucesso
			if !tc.expectErr {
				// Verifica se o resultado da função está correto
				expectedResult := fmt.Sprintf("Arquivo '%s' escrito com sucesso.", fullPath)
				if result != expectedResult {
					t.Errorf("writeFile() resultado = %q, esperado %q", result, expectedResult)
				}

				// Verifica se o arquivo foi realmente criado/atualizado
				actualContent, err := os.ReadFile(fullPath)
				if err != nil {
					t.Errorf("writeFile() não criou o arquivo: %v", err)
				}

				if string(actualContent) != tc.content {
					t.Errorf("writeFile() conteúdo do arquivo = %q, esperado %q", string(actualContent), tc.content)
				}
			}
		})
	}
}

// TestListFilesEdgeCases - Testa casos extremos do listFiles
func TestListFilesEdgeCases(t *testing.T) {
	testCases := []struct {
		name          string
		inputPath     string
		expectErr     bool
		expectedError string
	}{
		{
			name:      "Sucesso - Diretório vazio",
			inputPath: "",
			expectErr: false, // Deve usar diretório atual
		},
		{
			name:          "Erro - Diretório inexistente",
			inputPath:     "/diretorio/que/nao/existe",
			expectErr:     true,
			expectedError: "no such file or directory",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var inputData ListFilesInput
			var rawInput []byte

			if tc.inputPath == "" {
				// Teste com input vazio - deve usar diretório atual
				rawInput = []byte("{}")
			} else {
				inputData = ListFilesInput{Path: tc.inputPath}
				rawInput, _ = json.Marshal(inputData)
			}

			_, err := listFiles(rawInput)

			if (err != nil) != tc.expectErr {
				t.Fatalf("listFiles() erro = %v, expectErr %v", err, tc.expectErr)
			}

			if tc.expectErr && tc.expectedError != "" {
				if err == nil || !contains(err.Error(), tc.expectedError) {
					t.Errorf("listFiles() erro esperado contendo '%s', got '%v'", tc.expectedError, err)
				}
			}
		})
	}
}

// TestInvalidJSON - Testa JSON inválido para todas as funções
func TestInvalidJSON(t *testing.T) {
	testCases := []struct {
		name     string
		function func(json.RawMessage) (string, error)
	}{
		{"createDirectory", createDirectory},
		{"listFiles", listFiles},
		{"readFile", readFile},
		{"writeFile", writeFile},
	}

	invalidJSON := []byte(`{"invalid": json`)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.function(invalidJSON)
			if err == nil {
				t.Errorf("%s() deveria retornar erro para JSON inválido", tc.name)
			}
			if !contains(err.Error(), "JSON inválido") {
				t.Errorf("%s() erro deveria mencionar JSON inválido, got: %v", tc.name, err)
			}
		})
	}
}

// Helper function para verificar se uma string contém outra
func contains(s, substr string) bool {
	return len(s) >= len(substr) && 
		   (s == substr || 
		    (len(s) > len(substr) && 
		     (s[:len(substr)] == substr || 
		      s[len(s)-len(substr):] == substr || 
		      containsMiddle(s, substr))))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
