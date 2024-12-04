# Etapa 1: Build da aplicação usando Maven
FROM maven:3.9.4-eclipse-temurin-17 AS build
WORKDIR /app

# Copiar o arquivo pom.xml e instalar dependências
COPY leilao/pom.xml . 
RUN mvn dependency:go-offline

# Copiar o código-fonte para o container e empacotar o projeto
COPY leilao/src ./src
RUN mvn clean package -DskipTests

# Etapa 2: Executar a aplicação usando uma imagem de runtime do Java
FROM eclipse-temurin:17-jre-alpine
WORKDIR /app

# Copiar o arquivo JAR gerado na etapa de build
COPY --from=build /app/target/*.jar app.jar

# Expor a porta que será usada pela aplicação
EXPOSE 8080

# Comando para iniciar a aplicação
ENTRYPOINT ["java", "-jar", "app.jar"]
