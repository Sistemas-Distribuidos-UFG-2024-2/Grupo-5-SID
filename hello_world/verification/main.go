package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	PORT                = "5610"
	PrefixHealth        = "/health/"
	PrefixServers       = "/servers/"
	HealthCheckInterval = 10 * time.Second
	maxPoolSize         = 5 // Quantidade máxima de conexões no pool
)

type ServerInfo struct {
	IP      string `json:"ip"`
	Port    int    `json:"port"`
	Enabled bool   `json:"enabled"`
}

var (
	servers      = map[int]ServerInfo{}      // Mapa de servidores
	connPool     = make(map[string]net.Conn) // Pool de conexões TCP reutilizáveis
	connPoolMux  sync.Mutex                  // Mutex para proteger o acesso ao pool
	serversMutex sync.Mutex                  // Mutex para proteger a lista de servidores
)

func main() {
	ln, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
		return
	}
	defer ln.Close()
	fmt.Println("Servidor de verificação na porta " + PORT)

	// Goroutine para verificar a saúde dos servidores periodicamente
	go healthCheckServers()

	// Aceita conexões de clientes
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
		message = strings.TrimSpace(message)

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

	serversMutex.Lock()
	defer serversMutex.Unlock()

	if !serverInfo.Enabled {
		fmt.Printf("Servidor avisou que está indisponível: %d\n", serverInfo.Port)
		delete(servers, serverInfo.Port)
		return nil
	}

	servers[serverInfo.Port] = serverInfo
	fmt.Printf("Novo Servidor, IP: %s, Port: %d\n", serverInfo.IP, serverInfo.Port)
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
	serversMutex.Lock()
	defer serversMutex.Unlock()

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

// Função para verificar a saúde dos servidores periodicamente com reutilização de conexões
func healthCheckServers() {
	for {
		time.Sleep(HealthCheckInterval)

		serversMutex.Lock()
		for key, server := range servers {
			address := fmt.Sprintf("%s:%d", server.IP, server.Port)

			conn, err := getConnection(address)
			if err != nil {
				fmt.Printf("Servidor %s:%d não está acessível\n", server.IP, server.Port)
				server.Enabled = false
				delete(servers, key)
				continue
			}

			// Verifica a saúde do servidor
			fmt.Fprintf(conn, "%s\n", PrefixHealth)
			_, err = bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				fmt.Printf("Erro ao verificar o servidor %s:%d: %v\n", server.IP, server.Port, err)
				server.Enabled = false
				delete(servers, key)
			} else {
				server.Enabled = true
			}

			releaseConnection(address, conn)
		}
		serversMutex.Unlock()
	}
}

// Obtém uma conexão reutilizada ou cria uma nova se necessário
func getConnection(address string) (net.Conn, error) {
	connPoolMux.Lock()
	defer connPoolMux.Unlock()

	// Verifica se já existe uma conexão disponível no pool
	if conn, exists := connPool[address]; exists && conn != nil {
		return conn, nil
	}

	// Se não houver uma conexão no pool, cria uma nova
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	connPool[address] = conn
	return conn, nil
}

// Libera a conexão de volta para o pool
func releaseConnection(address string, conn net.Conn) {
	connPoolMux.Lock()
	defer connPoolMux.Unlock()

	// Adiciona a conexão de volta ao pool, para reutilização futura
	connPool[address] = conn
}
