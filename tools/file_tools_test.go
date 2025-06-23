package tools

import (
	"encoding/json"
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
