package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"main/pkg"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		number := rand.Int()

		cliente := pkg.Customer{
			Nome:  fmt.Sprintf("Client %d", number),
			Email: fmt.Sprintf("%d@mail.com", number),
		}

		fmt.Printf("Nome: %s, Email: %s\n", cliente.Nome, cliente.Email)

		// Converte o cliente para JSON
		clienteJSON, err := json.Marshal(cliente)
		if err != nil {
			fmt.Println("Erro ao converter para JSON:", err)
			return
		}

		// Envia a requisição POST para o servidor
		resp, err := http.Post("http://localhost:8080/cliente", "application/json", bytes.NewBuffer(clienteJSON))
		if err != nil {
			fmt.Println("Erro ao enviar requisição:", err)
			return
		}
		defer resp.Body.Close()

		// Lê a resposta do servidor
		fmt.Println("Resposta do servidor:", resp.Status)
	}
}
