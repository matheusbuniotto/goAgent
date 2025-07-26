# goAgent
Implementação **bare bones** de agentes de IA utilizando a linguagem go.

![main-image-gemini-generated](https://github.com/user-attachments/assets/49df3432-b530-481c-bc67-20fadaa0d263)

## Sobre o projeto
O goAgent é uma implementação de um agente de IA desenvolvida em Go, sem auxilio de nenhum SDK externo. Ele possui diversas ferramentas que podem ser utilizadas para interagir com sistemas de arquivos, automatizar tarefas ou estender funcionalidades de acordo com as necessidades dos usuários. O objetivo principal é entender o funcionamento de um agente com acesso a ferramentas. Esse conteúdo é baseado no artigo [How to build an Agent](https://ampcode.com/how-to-build-an-agent) que implementa um agente usando o SDK do Claude. 

Aqui, focaremos em modelos da OpenAI e Gemini, por enquanto. Porém, a implementação é feita sem auxilio dos SDK dos provedores.

## Usando
Defina sua chave API OpenRouter, OpenAI ou Gemini.

```bash
export OPENROUTER_API_KEY=.....  # Recomendado - acesso a múltiplos modelos
export GEMINI_API_KEY=.....
export OPENAI_API_KEY=.....
```
ou fish
```fish
set -x GEMINI_API_KEY ......
```

```bash
go mod tidy
go run main.go
```
```bash
// Especifique o provedor, se preferir
go run ./cmd/goagent -model openrouter
go run ./cmd/goagent -model openai
go run ./cmd/goagent -model gemini

// NOVO: especifique se o modelo tem acesso ao reasoning 
go run ./cmd/goagent --agent reasoning 
```

Obs: OpenRouter usa gpt-4o-mini por padrão (modelo econômico), o modelo gemini é o flash e o modelo da OpenAI é o gpt-4.1-nano. É possível utilizar o Gemini de forma gratuíta gerando uma chave em https://aistudio.google.com/apikey e o OpenRouter oferece créditos iniciais gratuitos.

## Arquitetura (até o momento)
![diagram(1)](https://github.com/user-attachments/assets/b270a0ad-9665-4f94-a0d2-e57995b687f6)

## Ferramentas disponíveis

É possível verificar as ferramentas disponíveis perguntando ao agente. As ferramentas estão localizadas em /tools com uma arquitetura hexagonal de ports/adapters isolando a lógica de interação com o agente.


![image](https://github.com/user-attachments/assets/001025f1-716e-4659-94af-bd4d088dc44d)

**NOVO**: Modo reasoning (think), implementa lógica de racicionio para enriquecer o contexto.


## Como criar uma nova ferramenta

Para criar uma nova ferramenta, siga os passos abaixo:

1. **Implemente a lógica da sua ferramenta**: Crie uma função que siga o tipo `ToolFunction`. Essa função receberá argumentos como JSON cru (`json.RawMessage`) e retornará uma string ou erro.

2. **Crie a definição da ferramenta**: Instancie uma variável `ToolDefinition` com o nome, descrição e a função criada.

3. **Adapte para o sistema**: Crie um `ToolAdapter` usando a definição criada. Você pode registrar esse adaptador no sistema de ferramentas para que possa ser utilizado pelo agente.

4. **Registre a função no main**: Forneça acesso a ferramenta ao agente no bloco allTools no arquivo main. 

### Exemplo prático

```go
// ::: Ferramenta: CreateDirectory :::
type CreateDirectoryInput struct {
	Path string `json:"path"`
}

// Definindo a função
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

// cria a definição da nova ferramenta
var CreateDirectoryDef = ToolDefinition{
	Name:        "create_directory",
	Description: `Cria um novo diretório no caminho especificado, necessita de um nome. Exemplo: {"path": "meu/novo/nome_diretorio"}`, //muito importante para comunicar com o agente.
	Function:    createDirectory,
}

```

### Como usar
Quando fizer sentido, o agente usará a nova ferramenta `createDirectory{path}`, realizando a criação da pasta no local específicado:

```
Humano: crie uma pasta chamada pasta-nova-teste
GoAgent está processando a mensagem...
GoAgent quer usar a ferramenta: create_directory({"path": "pasta-nova-teste"})
Resultado da ferramenta: Diretório 'pasta-nova-teste' criado com sucesso.
GoAgent está processando a mensagem...
GoAgent: OK. A pasta "pasta-nova-teste" foi criada com sucesso.
```

## Observação
Esse README foi gerado quase integralmente através das ferramentas disponíveis no agente, self made README.md 🤣

## Roadmap
[ ] Makefile

[ ] Mais modelos ou routers

[ ] Expandir testes e abordagem TDD

[ ] Especificar modelo no args (gpt-4o, etc)

[ ] Adicionar mais ferramentas

[ ] Remover arquivos ocultos da leitura ou colocar um .agentigore

[ ] Melhorar a interação e leitura para edição em partes específicas dos arquivos/textos.

[ ] Paramêtro para ajustar o quanto o modelo vai pedir confirmações para ações (human in the loop)

[ ] Ferramenta otimizada para criar novas ferramentas para o modelo 🔁

[ ] Traduzir para inglês / repo bilingue 

[ ] ....
