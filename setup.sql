-- Ativar a extensão para gerar UUIDs automaticamente
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Criar o schema schema_accounts
CREATE SCHEMA IF NOT EXISTS schema_accounts;

-- Criar a tabela accounts no schema schema_accounts
CREATE TABLE schema_accounts.accounts (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    mail VARCHAR(255) NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT TRUE
);

-- Criar o schema schema_auction
CREATE SCHEMA IF NOT EXISTS schema_auction;

-- Criar a tabela auctions no schema schema_auction
CREATE TABLE schema_auction.auctions (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY
);

-- Criação do schema leilao
CREATE SCHEMA IF NOT EXISTS leilao_schema;

-- Criação da tabela leilao
CREATE TABLE IF NOT EXISTS leilao_schema.leilao (
    id SERIAL PRIMARY KEY,
    produto VARCHAR(255) NOT NULL,
    lance_inicial DECIMAL(10, 2) NULL,
    data_finalizacao TIMESTAMP NOT NULL,
    criador VARCHAR(50) NOT NULL,
    vencedor VARCHAR(50) NULL,
    lance_final DECIMAL(10, 2) NULL,
    valor_maximo DECIMAL(10, 2) NULL
);

-- Criação da tabela participantes
CREATE TABLE IF NOT EXISTS leilao_schema.participantes (
    leilao_id BIGINT NOT NULL,
    usuario_email VARCHAR(50) NOT NULL,
    lance DECIMAL(10, 2) NOT NULL,
    PRIMARY KEY (leilao_id, usuario_email),
    FOREIGN KEY (leilao_id) REFERENCES leilao_schema.leilao (id) ON DELETE CASCADE
);