# Grupo 5 - Sistema de Leilão Distribuído (SID)


<img src="assets/leilao.jpg" width="120" height="120" alt="Logo SID">

## Link Pitch

[Link Apresentação 08/11](https://docs.google.com/presentation/d/1eK21P_mmpFl_WHYtD-Mm9hy7vjfw-AyXXOsiqwvpQU8/edit#slide=id.g3045f246402_0_7)

### Integrantes

Roberta Assis de Carvalho - 202300425

Christian Alexandre - 201709627

Júlio César Vieira Cruz - 202306001

Cleverson Oliveira


<img src="assets/sid.jpg" width="80" height="80" alt="Logo SID">

# Resumo

O SID, nosso sistema de leilão distribuído possui duas aplicações para o back-end e um front-end.
Este sistema web possui a intenção de prover a ferramenta de um marketplace de leilão, onde as pessoas podem criar um leilão e/ou participar de leilões e no final é enviado e-mail para o vencedor.

A nossa aplicação Leilão é responsável por criar e finalizar leilão, enviar e-mail para o vencedor, a aplicação foi escrita em Java, salvando os dados em um banco de dados relacional(Postgres) e para o envio de e-mail é utilizando um message broker (RabbitMQ)

A aplicação Leilão Ativo é responsável por receber, tratar lances e responder os lances de um determinado leilão via websocket, e avisar à aplicação Leilão sobre quem foi o vencedor, a aplicação foi escrita em Golang salvando os dados em um banco não relacional em memória (Redis), realiza lock distribuído também com o Redis e utiliza message broker (RabbitMQ) para atualizar os clients do websocket.

A aplicação do Front-end é responsável por realizar as interações com os serviços do back-end, para criar, visualizar e finalizar leilões, a aplicação foi escrita em Flutter para web.

## Como executar o projeto ?

1. Executar: (antes instalar docker e abrir docker desktop)
```
docker-compose up -d
```
