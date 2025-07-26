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

O goAgent oferece **múltiplas formas** de configurar chaves de API para máxima flexibilidade:

### 🔑 Configuração de API Keys

#### **1. 🌍 Variáveis de Ambiente (Tradicional)**
```bash
export OPENROUTER_API_KEY=your_key_here    # Recomendado - múltiplos modelos
export GEMINI_API_KEY=your_key_here        # Google Gemini
export OPENAI_API_KEY=your_key_here        # OpenAI GPT
```

#### **2. 🚩 Flags de Linha de Comando**
```bash
./goagent --openrouter-key "your_key" --agent reasoning
./goagent --gemini-key "your_key" --agent reasoning
./goagent --openai-key "your_key" --agent reasoning
```

#### **3. 📝 Arquivo .env**
```bash
# No diretório atual
echo "OPENROUTER_API_KEY=your_key" > .env

# Ou no home directory (global)
echo "OPENROUTER_API_KEY=your_key" > ~/.goagent.env
```

#### **4. 💬 Prompt Interativo (Automático)**
```bash
# Se nenhuma chave for encontrada, o sistema pergunta automaticamente:
./goagent --agent reasoning

# Resultado:
# ⚠️ Nenhuma chave de API encontrada. Vamos configurar uma:
# 🤖 Selecione um provedor de LLM:
# 1. OpenRouter  2. Gemini  3. OpenAI
# 🔑 Por favor, insira sua chave de API: [INPUT]
```

**Prioridade de verificação**: Flags → Env vars → .env local → .env home → Prompt interativo

### 🏃 Execução Rápida
```bash
# Com Go instalado
go mod tidy && go run ./cmd/goagent --agent reasoning

# Com binário compilado
./goagent --agent reasoning
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

O projeto segue o **padrão Hexagonal (Ports & Adapters)** com layout Go padrão:

### 📁 Estrutura de Diretórios
```
goAgent/
├── cmd/goagent/           # 🎯 Entry Point - CLI, flags, inicialização
├── pkg/                   # 🧠 Core Domain - Lógica de negócio
│   ├── agent/            #    • Orquestração, conversação, tool calling
│   └── toolkit/          #    • Sistema de ferramentas (ports/adapters)
├── internal/              # 🔧 Adapters - Implementações específicas
│   ├── llm/              #    • Clientes LLM (OpenRouter, OpenAI, Gemini)
│   ├── builtin/          #    • Ferramentas built-in (arquivos, interação)
│   └── prompts/          #    • Templates de prompts (sistema, reasoning)
└── dist/                  # 📦 Binários compilados
```

### 🏛️ Diagrama da Arquitetura Hexagonal

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              🚀 goAgent Architecture                        │
│                         (Hexagonal/Ports & Adapters)                       │
└─────────────────────────────────────────────────────────────────────────────┘

                               ┌─────────────────┐
                               │   cmd/goagent   │
                               │  (Entry Point)  │
                               │   🎯 CLI Flags  │
                               │   🔑 API Keys   │
                               │   📋 Selection  │
                               └─────────┬───────┘
                                         │
                                         ▼
┌────────────────────────────────────────────────────────────────────────────┐
│                           🧠 CORE DOMAIN (pkg/agent)                      │
│                                                                            │
│  ┌─────────────────┐              ┌─────────────────┐                     │
│  │      Agent      │◄────────────►│  LLMClient      │ (Port)              │
│  │   (Orchestrator)│              │   Interface     │                     │
│  │                 │              └─────────────────┘                     │
│  │ • Run()         │                        ▲                             │
│  │ • RunReasoning()│                        │                             │
│  │ • Tool Calling  │                        │                             │
│  │ • Conversation  │              ┌─────────────────┐                     │
│  │   Management    │              │      Tool       │ (Port)              │
│  │                 │◄────────────►│   Interface     │                     │
│  └─────────────────┘              │                 │                     │
│                                   │ • Name()        │                     │
│                                   │ • Description() │                     │
│                                   │ • Execute()     │                     │
│                                   └─────────────────┘                     │
└────────────────────────────────────────┬───────────────────────────────────┘
                                         │
                                         ▼
┌────────────────────────────────────────────────────────────────────────────┐
│                        🔧 TOOLKIT LAYER (pkg/toolkit)                     │
│                                                                            │
│  ┌─────────────────┐              ┌─────────────────┐                     │
│  │  ToolDefinition │              │   ToolAdapter   │ (Adapter)           │
│  │                 │◄────────────►│                 │                     │
│  │ • Name          │              │ Implements Tool │                     │
│  │ • Description   │              │   Interface     │                     │
│  │ • Function      │              │                 │                     │
│  │   (Logic)       │              │ • Name()        │                     │
│  └─────────────────┘              │ • Description() │                     │
│                                   │ • Execute()     │                     │
│                                   └─────────────────┘                     │
└────────────────────────────────────────┬───────────────────────────────────┘
                                         │
                                         ▼
┌────────────────────────────────────────────────────────────────────────────┐
│                     🛠️ ADAPTERS LAYER (internal/)                         │
│                                                                            │
│ ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐             │
│ │  LLM Adapters   │  │ Builtin Tools   │  │   Reasoning     │             │
│ │  (internal/llm) │  │(internal/builtin)│  │    Tools        │             │
│ │                 │  │                 │  │                 │             │
│ │ • OpenRouter    │  │ • file_tools    │  │ • analyze_      │             │
│ │ • Gemini        │  │ • interaction_  │  │   reasoning     │             │
│ │ • OpenAI        │  │   tools         │  │ • review_       │             │
│ │                 │  │                 │  │   decision      │             │
│ └─────────────────┘  └─────────────────┘  └─────────────────┘             │
└────────────────────────────────────────────────────────────────────────────┘

                         ┌─────────────────┐
                         │ 🌐 EXTERNAL WORLD │
                         │                 │
                         │ • APIs (LLM)    │
                         │ • File System   │
                         │ • Human (stdin) │
                         └─────────────────┘
```

### 🔄 Fluxo de Dados

1. **CLI** → Agent (comandos, chaves API)
2. **Agent** → LLMClient (conversação, lista de tools)
3. **LLMClient** → APIs Externas (requisições HTTP)
4. **Agent** → ToolAdapter → ToolDefinition → Lógica da Tool
5. **Tool** → Mundo Externo (sistema de arquivos, input humano)
6. Resultados retornam pelo mesmo caminho

### 🎯 Benefícios da Arquitetura

- **🧪 Testável**: Lógica central isolada de dependências externas
- **🔄 Flexível**: Fácil trocar provedores LLM ou adicionar novas tools
- **🛠️ Manutenível**: Separação clara de responsabilidades
- **🚀 Extensível**: Novos adapters sem alterar o core

## 🛠️ Ferramentas Disponíveis

O agente possui um conjunto robusto de ferramentas built-in organizadas por categoria:

### 📁 **Operações de Sistema**
- **`list_files`**: Lista arquivos e diretórios recursivamente
- **`read_file`**: Lê conteúdo de arquivos  
- **`write_file`**: Escreve/sobrescreve arquivos
- **`create_directory`**: Cria estruturas de diretórios

### 🤔 **Interação Humana**
- **`ask_human_for_clarification`**: Solicita esclarecimentos críticos do usuário

### 🧠 **Ferramentas de Reasoning (NOVO!)**
- **`analyze_reasoning`**: Analisa e valida qualidade do próprio raciocínio
- **`review_decision`**: Revisa criticamente decisões tomadas com scoring

### 🔧 **Sistema Extensível**
> 💡 **Arquitetura modular**: Adicione facilmente novas ferramentas usando o padrão ToolDefinition + ToolAdapter

**Exemplo de uso**:
```bash
# O agente pode usar ferramentas automaticamente:
"Analise meu raciocínio sobre usar React vs Vue para este projeto"
# → Usa analyze_reasoning automaticamente

# Ou você pode perguntar:
"Quais ferramentas você tem disponível?"
# → Lista todas as ferramentas com descrições
```


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
