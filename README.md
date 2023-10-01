# Pismo

<img align="right" width="50px" src="https://raw.githubusercontent.com/swaggo/swag/master/assets/swaggo.png">

O serviço possui 3 API`s, que são: 
- Cadastro de uma conta
- Busca da conta por ID
- Transação da conta

## Pré-requistos

- Docker
- docker-compose
## Iniciando

Rode o build da aplicação no Docker

```sh
docker-compose build
```
Iniciar aplicação

```sh
docker-compose up
```

A API está apontando para a porta :8080

```
http:localhost:8080/api/v1/
```

## Documentação

Foi usado o Swagger UI para gerar a documentação das API's

Para atualizar uma nova alteração a documentação rode o seguinte comando: 

```sh
swag init -g ./server/server.go
```

Para acessar a documentação, com o projeto rodando acesse: 
```sh
http://localhost:8080/swagger/index.html
```

## Testes

Foi usado o mockgen para os mocks. Os arquivos estão na pasta /mocks. Para gerar os mocks use o comando: 

```sh
go generate ./...
```

Para visualizar .html a cobertura de testes do projeto use o comando: 

```sh
go test  ./... -coverprofile=coverage.out
```

será gerado um arquivo coverage.out, para visualizar rode o comando: 

```sh
go tool cover -html=coverage.out 
```

## Parar

Para parar o serviço, use o comando: 

```sh
docker-compose down
```
