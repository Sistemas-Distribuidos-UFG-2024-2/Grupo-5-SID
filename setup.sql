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
