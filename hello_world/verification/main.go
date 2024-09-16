package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

const (
	PORT          = "5610"
	PrefixHealth  = "/health/"
	PrefixServers = "/servers/"
)

type ServerInfo struct {
	IP      string `json:"ip"`
	Port    int    `json:"port"`
	Enabled bool   `json:"enabled"`
}

var servers = map[int]ServerInfo{}

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

		switch {
		case strings.HasPrefix(message, PrefixHealth):
			err := HandleNewServer(conn, message)
			if err != nil {
				fmt.Fprintln(conn, "Erro ao processar a mensagem:", err)
				continue
			}
		case strings.HasPrefix(message, PrefixServers):
			err := HandleServerList(conn)
			if err != nil {
				fmt.Fprintln(conn, "Erro ao processar a lista de servidores:", err)
				continue
			}
		default:
			fmt.Fprintln(conn, "Comando desconhecido")
		}
	}
}

func HandleNewServer(conn net.Conn, message string) error {
	serverInfo, err := getServerInfoFromSocketMessage(message)
	if err != nil {
		fmt.Fprintln(conn, "Erro ao decodificar JSON")
		return err
	}

	if !serverInfo.Enabled {
		delete(servers, serverInfo.Port)
		return nil
	}

	servers[serverInfo.Port] = serverInfo

	fmt.Fprintln(conn, "Mensagem recebida")
	fmt.Printf("IP: %s, Port: %d\n", serverInfo.IP, serverInfo.Port)
	return nil
}

func getServerInfoFromSocketMessage(message string) (ServerInfo, error) {
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

func HandleServerList(conn net.Conn) error {
	serversSlice := make([]ServerInfo, 0, len(servers))
	for _, s := range servers {
		serversSlice = append(serversSlice, s)
	}

	serverList, err := json.Marshal(serversSlice)
	if err != nil {
		return err
	}

	// Envia a lista de servidores de volta ao cliente
	fmt.Fprintln(conn, string(serverList))
	return nil
}
