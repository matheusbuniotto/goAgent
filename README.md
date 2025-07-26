# goAgent
Implementação **bare bones** de agentes de IA utilizando Go com suporte a múltiplos provedores LLM.

![main-image-gemini-generated](https://github.com/user-attachments/assets/49df3432-b530-481c-bc67-20fadaa0d263)

## Sobre o projeto
O goAgent é uma implementação de um agente de IA desenvolvida em Go, sem auxílio de SDKs externos. Ele possui diversas ferramentas que podem ser utilizadas para interagir com sistemas de arquivos, automatizar tarefas ou estender funcionalidades de acordo com as necessidades dos usuários.

### 🎯 Características principais:
- **Múltiplos provedores**: OpenRouter, OpenAI e Gemini
- **Seleção interativa**: Escolha de provedor e modelo via interface
- **Modo reasoning**: Capacidade de raciocínio avançado
- **Arquitetura hexagonal**: Sistema de ferramentas modular
- **Zero dependências**: Implementação pura sem SDKs externos

### 🔗 Provedores suportados:
- **OpenRouter** (Recomendado): Acesso a GPT-4, Claude, Llama, Gemini e mais
- **Google Gemini**: Modelo Flash
- **OpenAI**: GPT-4o-mini e outros modelos

## 🚀 Instalação e Configuração

### 1. Configure sua chave API
```bash
# Escolha um ou mais provedores:
export OPENROUTER_API_KEY=your_key_here    # Recomendado - múltiplos modelos
export GEMINI_API_KEY=your_key_here        # Google Gemini
export OPENAI_API_KEY=your_key_here        # OpenAI GPT
```

### 2. Execute o projeto
```bash
go mod tidy
go run ./cmd/goagent
```

## Modos de Uso

### 🔄 Auto-detecção (Padrão)
```bash
go run ./cmd/goagent
# Detecta automaticamente: OpenRouter > Gemini > OpenAI
```

### 📋 Seleção Interativa
```bash
go run ./cmd/goagent -select
# Mostra menu para escolher provedor e modelo
```

### 🎯 Seleção Direta
```bash
go run ./cmd/goagent -model openrouter  # Pergunta qual modelo
go run ./cmd/goagent -model gemini      # Usa Gemini direto
go run ./cmd/goagent -model openai      # Usa OpenAI direto
```

### 🧠 Modo Reasoning
```bash
go run ./cmd/goagent --agent reasoning
# Ativa raciocínio avançado com tags <think>
```

## 🏗️ Arquitetura

O projeto segue o **layout padrão Go** com arquitetura hexagonal:

```
goAgent/
├── cmd/goagent/           # Aplicação principal
├── pkg/                   # Componentes reutilizáveis
│   ├── agent/            # Core do agente
│   └── toolkit/          # Sistema de ferramentas
├── internal/              # Código privado
│   ├── llm/              # Clientes LLM (OpenRouter, OpenAI, Gemini)
│   ├── builtin/          # Ferramentas built-in
│   └── prompts/          # Definições de prompts
└── examples/              # Exemplos de uso
```

![diagram(1)](https://github.com/user-attachments/assets/b270a0ad-9665-4f94-a0d2-e57995b687f6)

## 🛠️ Ferramentas Disponíveis

O agente possui ferramentas built-in para:

- **📁 Operações de arquivo**: Listar, ler e escrever arquivos
- **📂 Criação de diretórios**: Criar estruturas de pastas
- **🤔 Interação humana**: Perguntas diretas ao usuário
- **🔧 Sistema extensível**: Adicione suas próprias ferramentas facilmente

> 💡 **Dica**: Pergunte ao agente "quais ferramentas você tem disponível?" para ver a lista completa.


![image](https://github.com/user-attachments/assets/001025f1-716e-4659-94af-bd4d088dc44d)

**NOVO**: Modo reasoning (think), implementa lógica de racicionio para enriquecer o contexto.


## 🔧 Como criar uma nova ferramenta

Siga o padrão de arquitetura hexagonal:

### 1. Defina a estrutura de entrada
```go
// Em internal/builtin/minha_ferramenta.go
type MinhaFerramentaInput struct {
    Param1 string `json:"param1"`
    Param2 int    `json:"param2"`
}
```

### 2. Implemente a função
```go
func minhaFerramenta(input json.RawMessage) (string, error) {
    var typedInput MinhaFerramentaInput
    if err := json.Unmarshal(input, &typedInput); err != nil {
        return "", fmt.Errorf("JSON inválido: %w", err)
    }
    
    // Sua lógica aqui
    resultado := fmt.Sprintf("Processado: %s", typedInput.Param1)
    return resultado, nil
}
```

### 3. Crie a definição
```go
var MinhaFerramentaDef = toolkit.ToolDefinition{
    Name:        "minha_ferramenta",
    Description: "Descrição clara da ferramenta para o agente",
    Function:    minhaFerramenta,
}
```

### 4. Registre no main
```go
// Em cmd/goagent/main.go
allTools := []agent.Tool{
    // ... ferramentas existentes
    &toolkit.ToolAdapter{Definition: builtin.MinhaFerramentaDef},
}
``` 

## 🧪 Testes

```bash
# Executar todos os testes
go test ./...

# Testes com saída detalhada
go test -v ./internal/builtin

# Construir o projeto
go build ./cmd/goagent
```


## 🗺️ Roadmap

### ✅ Implementado
- [x] Suporte múltiplos provedores (OpenRouter, OpenAI, Gemini)
- [x] Seleção interativa de modelos
- [x] Arquitetura hexagonal
- [x] Modo reasoning
- [x] Layout padrão Go

### 🚧 Próximos passos
- [ ] Makefile para automação
- [ ] Expandir testes e abordagem TDD
- [ ] Adicionar mais ferramentas (web, APIs, etc.)
- [ ] Sistema de plugins
- [ ] Melhorar interação para edição de arquivos
- [ ] Configuração de confirmações (human-in-the-loop)
- [ ] Interface web opcional
- [ ] Suporte a diferentes formatos de saída
