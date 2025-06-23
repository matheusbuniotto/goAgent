
# goAgent
Implementação **bare bones** de agentes de IA utilizando a linguagem go.

## Sobre o projeto
O goAgent é uma implementação de um agente de IA desenvolvida em Go, sem auxilio de nenhum SDK externo. Ele possui diversas ferramentas que podem ser utilizadas para interagir com sistemas de arquivos, automatizar tarefas ou estender funcionalidades de acordo com as necessidades dos usuários. O objetivo principal é entender o funcionamento de um agente com acesso a ferramentas. Esse conteúdo é baseado no artigo https://ampcode.com/how-to-build-an-agent que implementa um agente usando o SDK do Claude. 

Aqui, focaremos em modelos da OpenAI e Gemini, por enquanto. Porém, a implementação é feita sem auxilio dos SDK dos provedores.

## Como contribuir
Quer contribuir com novas ferramentas ou melhorar o projeto? Sinta-se à vontade para criar novas funcionalidades ou oferecer sugestões. Para isso, adicione suas ferramentas na pasta *tools* e envie suas contribuições através de um pull request. Sempre que possível, documente bem suas implementações para facilitar a integração com o projeto principal.

## Usando
Defina sua chave API OpenAI ou Gemini.

```bash
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
go run main.go -model openai
go run main.go -model gemini
```

Obs: o modelo gemini é o flash e o modelo da OpenAI é o gpt-4.1-nano, ambos modelos bem economicos. É possível utilizar o Gemini de forma gratuíta gerando uma chave em https://aistudio.google.com/apikey

## Como criar uma nova ferramenta

Para criar uma nova ferramenta, siga os passos abaixo:

1. **Implemente a lógica da sua ferramenta**: Crie uma função que siga o tipo `ToolFunction`. Essa função receberá argumentos como JSON cru (`json.RawMessage`) e retornará uma string ou erro.

2. **Crie a definição da ferramenta**: Instancie uma variável `ToolDefinition` com o nome, descrição e a função criada.

3. **Adapte para o sistema**: Crie um `ToolAdapter` usando a definição criada. Você pode registrar esse adaptador no sistema de ferramentas para que possa ser utilizado pelo agente.

### Exemplo prático

```go
// Função que implementa a lógica da ferramenta
func GreetTool(input json.RawMessage) (string, error) {
    var args struct {
        Name string `json:"name"`
    }
    if err := json.Unmarshal(input, &args); err != nil {
        return "", err
    }
    return fmt.Sprintf("Olá, %s!", args.Name), nil
}

// Criando a definição da ferramenta
var GreetDefinition = tools.ToolDefinition{
    Name:        "greet",
    Description: "Retorna uma saudação personalizada.",
    Function:    GreetTool,
}

// Criando o adaptador
var GreetToolAdapter = tools.ToolAdapter{
    Definition: GreetDefinition,
}
```

### Como usar
Quando fizer sentido, o agente usará a nova função `GreetToolAdapter.Execute('{"name": "Carlos"}')`, retornando:

```
"Olá, Carlos!"
```

## Observação
Esse README foi gerado quase integralmente através das ferramentas disponíveis no agente.
![image](https://github.com/user-attachments/assets/3d0df82b-51ba-4b02-9992-e676ac099dd0)

## Roadmap
[ ] Mais modelos ou routers
[ ] Especificar modelo no args (gpt-4o, etc)
[ ] Adicionar mais ferramentas
[ ] Remover arquivos ocultos da leitura ou colocar um .agentigore
[ ] Melhorar a interação e leitura para edição em partes específicas dos arquivos/textos.
[ ] Paramêtro para ajustar o quanto o modelo vai pedir confirmações para ações (human in the loop)
[ ] ....
