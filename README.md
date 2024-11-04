# Grupo 5 - Sistema de Leilão Distribuído (SID)


<img src="assets/leilao.jpg" width="120" height="120" alt="Logo SID">

## Link Pitch

[Link do Pitch da Apresentação](https://docs.google.com/presentation/d/1eK21P_mmpFl_WHYtD-Mm9hy7vjfw-AyXXOsiqwvpQU8/edit#slide=id.g3045f246402_0_7)

### Integrantes

Roberta Assis de Carvalho - 202300425

Christian Alexandre - 201709627

Júlio César Vieira Cruz - 202306001

Cleverson Oliveira


<img src="assets/sid.jpg" width="80" height="80" alt="Logo SID">


## Como executar o projeto ?

1. Executar: (antes instalar docker e abrir docker desktop)
```
docker-compose up -d
```

2. Acessar o Grafana e o Prometheus
   O Grafana estará disponível em http://localhost:3000 (credenciais padrão: admin / admin123).
   O Prometheus estará disponível em http://localhost:9090.
3. Configurar o Grafana
   Após acessar o Grafana, você pode adicionar o Prometheus como uma fonte de dados:

- Faça login no Grafana.
- Vá para Configuration (Configuração) > Data Sources (Fontes de Dados).
- Clique em Add data source (Adicionar fonte de dados).
- Selecione Prometheus e configure a URL como http://prometheus_container:9090.
- Salve as configurações.
