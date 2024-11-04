package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"main/pkg"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

// Função para inicializar o banco de dados SQLite
func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./clientes.db")
	if err != nil {
		return nil, err
	}

	// Cria a tabela de clientes se não existir
	query := `
	CREATE TABLE IF NOT EXISTS cliente (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nome TEXT,
		leilao TEXT
	);`
	_, err = db.Exec(query)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Handler para a rota POST /cliente
func clienteHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		var cliente pkg.Customer
		err := json.NewDecoder(r.Body).Decode(&cliente)
		if err != nil {
			http.Error(w, "Erro ao processar JSON", http.StatusBadRequest)
			return
		}

		query := `INSERT INTO cliente (nome, email) VALUES (?, ?)`
		_, err = db.Exec(query, cliente.Nome, cliente.Email)
		if err != nil {
			http.Error(w, "Erro ao salvar no banco de dados", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "Customer %s salvo com sucesso", cliente.Nome)
	}
}

func main() {
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/cliente", clienteHandler(db))
	fmt.Println("Servidor iniciado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
