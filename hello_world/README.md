# Como executar o projeto

Necessário ter go e java instalado no terminal

Executar em terminais diferentes:

```
go run verification/main.go
go run loadbalancer/loadbalancer.go
go run server.go {PORTA}
java Server.java {PORTA}
java Client.java
```

Exemplo: 

```
go run verification/main.go
go run loadbalancer/loadbalancer.go
go run server.go 5601
java Server.java 5602
java Client.java
```

Observação a porta do servidor de verificação por padrão é a 5610