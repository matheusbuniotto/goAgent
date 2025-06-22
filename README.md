
# goAgent
IA Agent implementation in Go lang

## Sobre o projeto
O goAgent é uma implementação de um agente de IA desenvolvida em Go. Ele possui diversas ferramentas que podem ser utilizadas para interagir com sistemas de arquivos, automatizar tarefas ou estender funcionalidades de acordo com as necessidades dos usuários.

## Como contribuir
Quer contribuir com novas ferramentas ou melhorar o projeto? Sinta-se à vontade para criar novas funcionalidades ou oferecer sugestões. Para isso, adicione suas ferramentas na pasta *tools* e envie suas contribuições através de um pull request. Sempre que possível, documente bem suas implementações para facilitar a integração com o projeto principal.

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