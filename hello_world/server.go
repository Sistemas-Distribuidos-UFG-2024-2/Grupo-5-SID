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
	"time"
)

const (
	VerificationPort   = 5610
	VerificationPrefix = "/health/"
	IP                 = "localhost"

	HealthCheckTime = 15 * time.Second
)

type ServerInfo struct {
	IP      string `json:"ip"`
	Port    int    `json:"port"`
	Enabled bool   `json:"enabled"`
}

var (
	serverInfo ServerInfo
)

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

	// Inicializa as informações do servidor
	serverInfo = ServerInfo{IP: IP, Port: port, Enabled: true}

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
			fmt.Println("Encerrando o servidor:", sig)
			os.Exit(1)
		}
	}()

	// Envia informações do servidor para o servidor de verificação
	go func() {
		sendServerInfo(true)
		ticker := time.NewTicker(HealthCheckTime)
		defer ticker.Stop()
		for range ticker.C {
			sendServerInfo(true)
		}
	}()

	// Aceita conexões de clientes
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Erro de conexão:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func sendServerInfo(enabled bool) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", IP, VerificationPort))
	if err != nil {
		fmt.Println("Erro ao se conectar ao servidor de verificação:", err)
		return
	}
	defer conn.Close()

	serverInfo.Enabled = enabled

	// Serializa os dados do servidor para JSON
	serverInfoJSON, err := json.Marshal(serverInfo)
	if err != nil {
		fmt.Println("Erro ao converter ServerInfo para JSON:", err)
		return
	}

	// Envia a solicitação de verificação
	_, err = fmt.Fprintf(conn, "%s%s\n", VerificationPrefix, serverInfoJSON)
	if err != nil {
		fmt.Println("Erro ao enviar a solicitação de verificação:", err)
		return
	}

	// Lê a resposta do servidor de verificação
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Erro ao ler a resposta do servidor de verificação:", err)
	} else {
		fmt.Println("Resposta do servidor de verificação:", response)
	}
}
func handleConnection(conn net.Conn) {
    defer conn.Close()

    reader := bufio.NewReader(conn)
    for {
        message, _ := reader.ReadString('\n')
        message = strings.TrimSpace(message)

        if strings.HasPrefix(message, "FUNCIONARIO") {
            parts := strings.Split(message, ",")
            nome := parts[1]
            cargo := parts[2]
            salario, _ := strconv.ParseFloat(parts[3], 64)

            salarioReajustado := calculaReajuste(cargo, salario)
            response := fmt.Sprintf("Nome: %s, Salário Reajustado: %.2f", nome, salarioReajustado)
            fmt.Fprintln(conn, response)

            fmt.Println("Dados recebidos e processados:", message)
        }
    }
}

func calculaReajuste(cargo string, salario float64) float64 {
    if strings.ToLower(cargo) == "operador" {
        return salario * 1.20
    } else if strings.ToLower(cargo) == "programador" {
        return salario * 1.18
    }
    return salario
}
