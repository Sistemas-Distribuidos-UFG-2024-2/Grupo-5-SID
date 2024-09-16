package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
)

const (
	VerificationPort   = 5610
	VerificationPrefix = "/health"
	IP                 = "localhost"
)

type ServerInfo struct {
	IP      string `json:"ip"`
	Port    int    `json:"port"`
	Enabled bool   `json:"enabled"`
}

var serverInfo ServerInfo

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run main.go <porta>")
		return
	}

	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Porta inválida. Use um número inteiro.")
		return
	}
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
		return
	}
	defer ln.Close()
	fmt.Printf("Servidor na porta %d\n", port)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			fmt.Println("Sinal recebido:", sig)
			fmt.Println("Encerrando o servidor...")
			os.Exit(1)
		}
	}()

	serverInfo = ServerInfo{IP: IP, Port: port, Enabled: true}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", IP, VerificationPort))
	if err != nil {
		fmt.Println("Erro ao se conectar ao servidor de verificação:", err)
	} else {
		defer conn.Close()

		// Serializa os dados do servidor para JSON
		serverInfoJSON, err := json.Marshal(serverInfo)
		if err != nil {
			fmt.Println("Erro ao converter ServerInfo para JSON:", err)
			return
		}

		// Envia a solicitação de verificação
		fmt.Fprintf(conn, "%s/%s\n", VerificationPrefix, serverInfoJSON)

		// Lê a resposta do servidor de verificação
		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Erro ao ler a resposta do servidor de verificação:", err)
		} else {
			fmt.Println("Resposta do servidor de verificação:", response)
		}
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Erro de conexão:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		message, _ := reader.ReadString('\n')

		message = strings.TrimSpace(message)
		if strings.EqualFold(message, "hello") {
			fmt.Fprintln(conn, "world")
			fmt.Println("Mensagem recebida:", message)
		}
	}
}
