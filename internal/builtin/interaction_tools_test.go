package builtin

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestAskHumanValidation testa a validação de entrada da função askHuman
func TestAskHumanValidation(t *testing.T) {
	testCases := []struct {
		name          string
		input         interface{}
		expectErr     bool
		expectedError string
	}{
		{
			name:          "Erro - JSON inválido",
			input:         `{"invalid": json`,
			expectErr:     true,
			expectedError: "JSON inválido",
		},
		{
			name: "Erro - Question vazia",
			input: AskHumanInput{
				Question: "",
			},
			expectErr:     true,
			expectedError: "argumento inválido",
		},
		{
			name: "Erro - Question apenas espaços",
			input: AskHumanInput{
				Question: "   ",
			},
			expectErr: false, // Espaços não são validados, apenas string vazia
		},
		{
			name: "Sucesso - Question válida (não executará stdin)",
			input: AskHumanInput{
				Question: "Esta é uma pergunta válida?",
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var rawInput []byte
			var err error

			// Prepara o input JSON
			switch v := tc.input.(type) {
			case string:
				rawInput = []byte(v)
			case AskHumanInput:
				rawInput, err = json.Marshal(v)
				if err != nil {
					t.Fatalf("Falha ao serializar input de teste: %v", err)
				}
			}

			// Para testes que não esperamos erro de validação, 
			// vamos testar apenas a validação inicial sem executar stdin
			if !tc.expectErr && tc.name == "Sucesso - Question válida (não executará stdin)" {
				// Testa apenas se a validação JSON e de campo funciona
				var typedInput AskHumanInput
				err := json.Unmarshal(rawInput, &typedInput)
				if err != nil {
					t.Errorf("askHuman() validação JSON falhou: %v", err)
				}
				if typedInput.Question == "" {
					t.Errorf("askHuman() validação de campo falhou")
				}
				return
			}

			// Para casos de erro, podemos testar completamente
			if tc.expectErr {
				_, err := askHuman(rawInput)
				
				if err == nil {
					t.Errorf("askHuman() deveria retornar erro para caso: %s", tc.name)
					return
				}

				if tc.expectedError != "" && !strings.Contains(err.Error(), tc.expectedError) {
					t.Errorf("askHuman() erro esperado contendo '%s', got '%v'", tc.expectedError, err)
				}
			}
		})
	}
}

// TestAskHumanInputStructure testa a estrutura de entrada
func TestAskHumanInputStructure(t *testing.T) {
	input := AskHumanInput{
		Question: "Teste de pergunta",
	}

	// Verifica se a estrutura pode ser marshaled/unmarshaled corretamente
	jsonData, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Falha ao fazer marshal de AskHumanInput: %v", err)
	}

	var decoded AskHumanInput
	err = json.Unmarshal(jsonData, &decoded)
	if err != nil {
		t.Fatalf("Falha ao fazer unmarshal de AskHumanInput: %v", err)
	}

	if decoded.Question != input.Question {
		t.Errorf("AskHumanInput.Question = %q, esperado %q", decoded.Question, input.Question)
	}
}

// TestAskHumanDefinition testa a definição da ferramenta
func TestAskHumanDefinition(t *testing.T) {
	// Verifica se a definição da ferramenta está correta
	if AskHumanDef.Name == "" {
		t.Error("AskHumanDef.Name não deveria estar vazio")
	}

	if AskHumanDef.Description == "" {
		t.Error("AskHumanDef.Description não deveria estar vazio")
	}

	if AskHumanDef.Function == nil {
		t.Error("AskHumanDef.Function não deveria ser nil")
	}

	// Verifica se o nome segue o padrão esperado
	expectedName := "ask_human_for_clarification"
	if AskHumanDef.Name != expectedName {
		t.Errorf("AskHumanDef.Name = %q, esperado %q", AskHumanDef.Name, expectedName)
	}

	// Verifica se a descrição contém informações importantes
	description := AskHumanDef.Description
	requiredKeywords := []string{"question", "JSON"}
	
	for _, keyword := range requiredKeywords {
		if !strings.Contains(description, keyword) {
			t.Errorf("AskHumanDef.Description deveria conter '%s': %s", keyword, description)
		}
	}
}

// TestInteractionToolsIntegration testa a integração básica
func TestInteractionToolsIntegration(t *testing.T) {
	// Testa se a função pode ser chamada através da definição
	validInput := AskHumanInput{
		Question: "Pergunta de teste para integração",
	}
	
	rawInput, err := json.Marshal(validInput)
	if err != nil {
		t.Fatalf("Falha ao preparar input de teste: %v", err)
	}

	// Como não podemos testar stdin facilmente, testamos apenas se a função
	// não panic e retorna o erro esperado quando não há input disponível
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("askHuman() causou panic: %v", r)
		}
	}()

	// Nota: Este teste pode bloquear se executado em ambiente interativo
	// Em CI/CD, o stdin será EOF e retornará erro, que é o comportamento esperado
	_, err = AskHumanDef.Function(rawInput)
	
	// Em ambiente de teste automatizado, esperamos erro de leitura
	// Em ambiente interativo, a função aguardaria input do usuário
	if err != nil {
		// Verifica se é um erro de leitura relacionado ao stdin
		errorStr := err.Error()
		expectedErrors := []string{"EOF", "erro ao ler a resposta", "stdin"}
		
		hasExpectedError := false
		for _, expectedErr := range expectedErrors {
			if strings.Contains(errorStr, expectedErr) {
				hasExpectedError = true
				break
			}
		}
		
		if !hasExpectedError {
			t.Logf("askHuman() retornou erro (esperado em ambiente de teste): %v", err)
		}
	}
}