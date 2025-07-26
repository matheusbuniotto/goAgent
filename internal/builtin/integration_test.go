package builtin

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matheusbuniotto/goagent/pkg/toolkit"
)

// TestToolDefinitionsIntegration testa todas as definições de ferramentas
func TestToolDefinitionsIntegration(t *testing.T) {
	definitions := []struct {
		name string
		def  toolkit.ToolDefinition
	}{
		{"ListFilesDef", ListFilesDef},
		{"WriteFileDef", WriteFileDef},
		{"ReadFileDef", ReadFileDef},
		{"CreateDirectoryDef", CreateDirectoryDef},
		{"AskHumanDef", AskHumanDef},
	}

	for _, td := range definitions {
		t.Run(td.name, func(t *testing.T) {
			// Verifica se todos os campos obrigatórios estão preenchidos
			if td.def.Name == "" {
				t.Errorf("%s.Name está vazio", td.name)
			}

			if td.def.Description == "" {
				t.Errorf("%s.Description está vazio", td.name)
			}

			if td.def.Function == nil {
				t.Errorf("%s.Function é nil", td.name)
			}

			// Verifica se a descrição é útil (deve ter pelo menos 20 caracteres)
			if len(td.def.Description) < 20 {
				t.Errorf("%s.Description muito curta: %s", td.name, td.def.Description)
			}
		})
	}
}

// TestToolAdapterIntegration testa a integração com o ToolAdapter
func TestToolAdapterIntegration(t *testing.T) {
	// Testa algumas ferramentas através do ToolAdapter
	testCases := []struct {
		name    string
		adapter *toolkit.ToolAdapter
		input   interface{}
		wantErr bool
	}{
		{
			name:    "CreateDirectory via ToolAdapter",
			adapter: &toolkit.ToolAdapter{Definition: CreateDirectoryDef},
			input:   CreateDirectoryInput{Path: filepath.Join(t.TempDir(), "test-dir")},
			wantErr: false,
		},
		{
			name:    "CreateDirectory erro via ToolAdapter",
			adapter: &toolkit.ToolAdapter{Definition: CreateDirectoryDef},
			input:   CreateDirectoryInput{Path: ""}, // Caminho vazio deve dar erro
			wantErr: true,
		},
		{
			name:    "WriteFile via ToolAdapter",
			adapter: &toolkit.ToolAdapter{Definition: WriteFileDef},
			input:   WriteFileInput{Path: filepath.Join(t.TempDir(), "test.txt"), Content: "conteúdo teste"},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Verifica métodos do ToolAdapter
			if tc.adapter.Name() == "" {
				t.Error("ToolAdapter.Name() retornou string vazia")
			}

			if tc.adapter.Description() == "" {
				t.Error("ToolAdapter.Description() retornou string vazia")
			}

			// Testa execução através do adapter
			inputJSON, err := json.Marshal(tc.input)
			if err != nil {
				t.Fatalf("Falha ao fazer marshal do input: %v", err)
			}

			result, err := tc.adapter.Execute(string(inputJSON))

			if (err != nil) != tc.wantErr {
				t.Errorf("ToolAdapter.Execute() erro = %v, wantErr %v", err, tc.wantErr)
			}

			if !tc.wantErr && result == "" {
				t.Error("ToolAdapter.Execute() retornou resultado vazio para caso de sucesso")
			}
		})
	}
}

// TestFileToolsWorkflow testa um fluxo completo de operações de arquivo
func TestFileToolsWorkflow(t *testing.T) {
	tempDir := t.TempDir()
	
	// Workflow: Criar diretório -> Escrever arquivo -> Ler arquivo -> Listar arquivos
	
	// 1. Criar diretório
	subDir := filepath.Join(tempDir, "workflow-test")
	createInput := CreateDirectoryInput{Path: subDir}
	createJSON, _ := json.Marshal(createInput)
	
	result, err := createDirectory(createJSON)
	if err != nil {
		t.Fatalf("Falha ao criar diretório: %v", err)
	}
	
	if !strings.Contains(result, "criado com sucesso") {
		t.Errorf("Resultado de criação de diretório inesperado: %s", result)
	}

	// 2. Escrever arquivo
	testFile := filepath.Join(subDir, "workflow.txt")
	testContent := "Este é um teste de workflow\ncom múltiplas linhas"
	writeInput := WriteFileInput{Path: testFile, Content: testContent}
	writeJSON, _ := json.Marshal(writeInput)
	
	result, err = writeFile(writeJSON)
	if err != nil {
		t.Fatalf("Falha ao escrever arquivo: %v", err)
	}
	
	if !strings.Contains(result, "escrito com sucesso") {
		t.Errorf("Resultado de escrita de arquivo inesperado: %s", result)
	}

	// 3. Ler arquivo
	readInput := ReadFileInput{Path: testFile}
	readJSON, _ := json.Marshal(readInput)
	
	result, err = readFile(readJSON)
	if err != nil {
		t.Fatalf("Falha ao ler arquivo: %v", err)
	}
	
	if result != testContent {
		t.Errorf("Conteúdo lido = %q, esperado %q", result, testContent)
	}

	// 4. Listar arquivos
	listInput := ListFilesInput{Path: subDir}
	listJSON, _ := json.Marshal(listInput)
	
	result, err = listFiles(listJSON)
	if err != nil {
		t.Fatalf("Falha ao listar arquivos: %v", err)
	}

	var files []string
	err = json.Unmarshal([]byte(result), &files)
	if err != nil {
		t.Fatalf("Falha ao decodificar lista de arquivos: %v", err)
	}

	// Verifica se o arquivo criado está na lista
	found := false
	for _, file := range files {
		if file == testFile {
			found = true
			break
		}
	}
	
	if !found {
		t.Errorf("Arquivo criado não encontrado na listagem: %v", files)
	}
}

// TestErrorHandlingConsistency testa consistência no tratamento de erros
func TestErrorHandlingConsistency(t *testing.T) {
	testCases := []struct {
		name     string
		function func(json.RawMessage) (string, error)
		input    string
	}{
		{"createDirectory", createDirectory, `{}`},
		{"listFiles", listFiles, `{"path": "/diretorio/inexistente"}`},
		{"readFile", readFile, `{}`},
		{"writeFile", writeFile, `{}`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.function(json.RawMessage(tc.input))
			
			// Todas as funções devem retornar erro para inputs inválidos
			if err == nil {
				t.Errorf("%s() deveria retornar erro para input inválido", tc.name)
			}

			// Verifica se a mensagem de erro é informativa
			errorMsg := err.Error()
			if len(errorMsg) < 10 {
				t.Errorf("%s() mensagem de erro muito curta: %s", tc.name, errorMsg)
			}

			// Verifica se não há panic
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("%s() causou panic: %v", tc.name, r)
				}
			}()
		})
	}
}

// TestConcurrentAccess testa acesso concorrente às ferramentas
func TestConcurrentAccess(t *testing.T) {
	tempDir := t.TempDir()
	
	// Testa criação concorrente de diretórios
	const numGoroutines = 10
	results := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			dirPath := filepath.Join(tempDir, "concurrent", "dir"+string(rune('0'+id)))
			input := CreateDirectoryInput{Path: dirPath}
			inputJSON, _ := json.Marshal(input)
			
			_, err := createDirectory(inputJSON)
			results <- err
		}(i)
	}

	// Coleta resultados
	for i := 0; i < numGoroutines; i++ {
		err := <-results
		if err != nil {
			t.Errorf("Erro na criação concorrente de diretório %d: %v", i, err)
		}
	}

	// Verifica se todos os diretórios foram criados
	for i := 0; i < numGoroutines; i++ {
		dirPath := filepath.Join(tempDir, "concurrent", "dir"+string(rune('0'+i)))
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			t.Errorf("Diretório %d não foi criado: %s", i, dirPath)
		}
	}
}

// BenchmarkFileOperations benchmark das operações de arquivo
func BenchmarkFileOperations(b *testing.B) {
	tempDir := b.TempDir()
	
	b.Run("CreateDirectory", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			dirPath := filepath.Join(tempDir, "bench", "dir", "subdir"+string(rune('0'+i%10)))
			input := CreateDirectoryInput{Path: dirPath}
			inputJSON, _ := json.Marshal(input)
			
			_, err := createDirectory(inputJSON)
			if err != nil {
				b.Fatalf("Erro no benchmark: %v", err)
			}
		}
	})

	b.Run("WriteFile", func(b *testing.B) {
		content := "Conteúdo de teste para benchmark"
		for i := 0; i < b.N; i++ {
			filePath := filepath.Join(tempDir, "bench", "file"+string(rune('0'+i%10))+".txt")
			input := WriteFileInput{Path: filePath, Content: content}
			inputJSON, _ := json.Marshal(input)
			
			_, err := writeFile(inputJSON)
			if err != nil {
				b.Fatalf("Erro no benchmark: %v", err)
			}
		}
	})
}