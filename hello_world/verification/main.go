package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
)

const (
	PORT         = "5610"
	PrefixHealth = "/health/"
)

type ServerInfo struct {
	IP      string `json:"ip"`
	Port    int    `json:"port"`
	Enabled bool   `json:"enabled"`
}

var servers = map[string]ServerInfo{}

func main() {
	ln, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
		return
	}
	defer ln.Close()
	fmt.Println("Servidor de verificação na porta " + PORT)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Erro de conexão no servidor de verificação:", err)
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

		serverInfo, err := getServerInfoFromSocketMessage(message)
		if err != nil {
			fmt.Fprintln(conn, "Erro ao decodificar JSON")
			continue
		}

		if !serverInfo.Enabled {
			delete(servers, serverInfo.IP)
			continue
		}

		servers[serverInfo.IP] = serverInfo

		fmt.Fprintln(conn, "Mensagem recebida")
		fmt.Printf("IP: %s, Port: %d\n", serverInfo.IP, serverInfo.Port)
	}
}

func getServerInfoFromSocketMessage(message string) (ServerInfo, error) {
	if !strings.HasPrefix(message, PrefixHealth) {
		return ServerInfo{}, errors.New("O prefix /health/ nao foi enviado")
	}
	message = strings.TrimPrefix(message, PrefixHealth)
	message = strings.TrimSuffix(message, "\n")
	var serverInfo ServerInfo

	// Convertendo a string JSON para a struct
	err := json.Unmarshal([]byte(message), &serverInfo)
	if err != nil {
		return ServerInfo{}, err
	}

	return serverInfo, nil
}
