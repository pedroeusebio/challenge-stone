# challenge-stone


## Indice

- [Instalação](https://github.com/pedroeusebio/challenge-stone#instalacao)
- [Bibliotecas](https://github.com/pedroeusebio/challenge-stone#bibliotecas)
- [Documentação da API](https://github.com/pedroeusebio/challenge-stone#documentacao-da-api)
- [Referências](https://github.com/pedroeusebio/challenge-stone#referencias) 

## Instalação 

O manual de instalação tem o intuito de preparar o ambiente do sistema, incluindo o banco de dados e as bibliotecas usadas.

Para ter acesso ao codigo basta um git clone.
```shell
	git clone https://github.com/pedroeusebio/challenge-stone.git
```
O proximo passo é configurar o docker para ter disponível o banco postgres configurado e com as tabelas preparadas para a aplicação:
```shell
    cd challenge-stone
    docker build -t <nome_desejado>
    docker volume create pgdata # criar um volume para armazenar os dados do postgres
    docker run -dt -p 5432:5432 --rm -v pgdata:/var/lib/postgresql/data <nome_desejado> 
```
Esse último comando serve para subir o banco sem travar o terminal (-dt) e fazendo um *forward* da porta 5432 local para a 5432 da maquina virtual do docker. Então caso a porta já esteja sendo utilizada, só alterar para a porta desejada.

A API foi desenvolvida utilizando o [glide](https://github.com/Masterminds/glide) como *Package Management*, dessa forma, basta utilizar os comando abaixo para obter as *libs* de acordo com a versão determinada no arquivo .yaml: 
```shell
    cd ./src && glide install
```

Tendo toda as libs e banco configurado, a ultima etapa é alterar o GOPATH para que o binaário encontre os arquivos das libs na pasta do projeto. Para isso é necessário colocar o *path* do repositório em algum arquivo que o seu shell leia (bash, zsh):
```shell
    pwd # descobrir o caminho da pasta do repositorio
```
Após copiar o caminho é só inserí-lo no .bashrc ou outro de preferência e finalizar com o carregamento dos dados : `source ~/.bashrc`.

Tendo o $GOPATH configurado para rodar a aplicação basta usar o comando: 
```shell
    go run src/challengestone.go
```

## Bibliotecas

Segue abaixo a lista de bibliotecas utilizadas e a justificativa para a sua utilização.
1. [Squirrel](https://github.com/Masterminds/squirrel)
    - Squirrel é uma lib que facilita a criação das Queries em SQL, dessa forma se fez extremamente necessario para trabalhar com o banco de dados de maneira mais rapida e facil.
2. [Sqlx](https://github.com/jmoiron/sqlx)
    - escrever alguma coisa
3. [Httprouter](https://github.com/julienschmidt/httprouter)
    - Usado para gerenciar as rotas da API com mais facildade que o net/http. 
4. [PQ](https://github.com/lib/pq)
    - A biblioteca basica para acesso aos bancos PostgreSQL.
5. [Validator](https://gopkg.in/go-playground/validator.v9)
    - Usado para realizar a validação na criação dos Usuários e das notas fiscais. 
6. [JWT-GO](https://github.com/dgrijalva/jwt-go)
    - Usado para realizar a autenticação e gerar o token do usuário logado.

## Documentação da API

A API possui 3 *endpoints* : 
- user : é a base de acesso para a criação e listagem dos usuários.
- invoice : é a base de acesso para a criação, listagem e remoção logica das notais fiscais.
- login : é a base de autenticação do usuário


### *User*

#### Listagem

O acesso da listagem é feito atraves da url `/user` com o metodo GET, e podem ser passados 4 parametros para filtrar os dados a serem exibidos:
- name : esse parametro é usado para solicitar um usuário que tenha a *substring* desejada. 
- page : parametro usado para solicitar a pagina deseja para visualização.
- length : parametro usado para determinar quantos dados serão exibidos na pagina.
- order : parametro usado para determinar a ordenação dos dados. É possivel utiliza-lo passando multiplas ordenações. O *order* solicita o envio de um objeto json com duas chaves : a coluna desejada e o sentido da ordenação (ASC ou DESC), caso não exista a coluna solicitada é retornado um erro. A ordenação da priopridade ao primeiro objeto do pedido.
    - exemplo de uso:
    `http://localhost:3000/user?order=[{"Column":"name","Order":"desc"},{"Column":"password","Order":"asc"}]`
    
#### Inserção

A inserção é feita através da url `/user/`, porém o método utilizado é o POST. Dessa forma, é necessário enviar um formulário para que o usuário possa ser devidamente inserido. O formado de envio do corpo da requisição é o `x-www-form-urlencoded`. A inserção respeita os requisitos do sistema, sendo assim, caso o *name* ou o *password* estejam fora das regras, é retornado um erro de validação com uma mensagem de exibição, podendo ser usada diretamente pela *client*. Caso tenho ocorrido um sucesso, é respondido com um json uma mensagem de sucesso e as informações do usuário inserido que pode ser usada para confirmar o pedido.


### *Invoice*

#### Listagem

Da mesma maneira que o usuário, foi utilizado o metodo GET e o pedido é feito atraves da url `/invoice`. Podem ser passados apenas 3 paramêtros que filtram os dados a serem exibidos:
- page
- length
- order

os 3 paramêtros acima se comportam exatamente da mesma maneira que a listagem dos usuários, retornando caso bem sucedido, um json contendo o metodo tipo de pedido solicitado e  um *payload* que é um *array* de objetos invoice, contendo todos os dados de cada invoice.

### Inserção

A inserção é deita através do url `/invoice` utilizando o metodo POST. Funciona da mesma forma do que o *user*, retornando um json com erro, caso não passe no teste de validação ou de algum erro na inserção no banco e no caso de sucesso é retornado o método com o objeto json do invoice inserido.

### Remoção

Como solicitado, a remoção é apenas lógica, ou seja, caso o usuário remova um *invoice*, a *flag* `is_active` passa a receber o valor de `false`. A operação é acessada através da url `/invoice/:id` usando o método DELETE. Para deletar é necessário inserir na url o id do *invoice* desejado. Caso a id exista, a *invoice* será deletada com sucesso e receberá um json com uma mensagem de confirmação e um *payload* do usuário que foi deletado. Já no caso de erro, onde o usuário não existe, receberá um json de com a mensagem de erro e um payload com um invoice vazio.

Obs: Os métodos de inserção e remoção são necessários a autenticação do usuário. Caso não seja enviado no cabeçalho do http *request*, o pedido será respondido com uma mensagem de erro informando a falta de autenticação.

### *Request Token*

#### *Login*

A operação de *login* é feita de através da url `/login` utilizando um metodo POST. É necessário o envio do *name* e *password* do usuário que deseja autenticar. Caso o usuário não exista ou a senha esteja incorreta, é retornado um json com a mensagem de erro e o token vazio. Caso o usuário consiga autenticar com sucesso, é gerado um JWT feita da composição do *name* e do *password* do usuário, com uma validade de 10 anos, dessa forma, é enviado como resposta um json contendo a mensagem de sucesso e token gerado. Para poder utilizar esse *token* para a validação do usuário é necessário adicionar a *header* `X-Token` no cabeçalho do protocolo HTTP que enviará o pedido.

