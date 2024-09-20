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
	VerificationServer  = "localhost:5610"
	PrefixServers       = "/servers/"
	PORT                = "5611"
	HealthCheckInterval = 10 * time.Second
)

type ServerInfo struct {
	IP      string `json:"ip"`
	Port    int    `json:"port"`
	Enabled bool   `json:"enabled"`
}

var (
	servers        []ServerInfo        // Lista de servidores ativos
	serverIndex    int                 // Índice do servidor para distribuição
	serverIndexMux sync.Mutex          // Mutex para evitar condições de corrida no acesso ao índice do servidor
	connPool       map[string]net.Conn // Pool de conexões ativas
	connPoolMux    sync.Mutex          // Mutex para proteger o acesso ao pool de conexões
)

func main() {
	connPool = make(map[string]net.Conn) // Inicializa o pool de conexões

	// Atualiza periodicamente a lista de servidores
	go func() {
		ticker := time.NewTicker(HealthCheckInterval)
		defer ticker.Stop() // Para garantir que o ticker seja parado quando a função for encerrada

		for {
			select {
			case <-ticker.C:
				err := updateServerList()
				if err != nil {
					fmt.Printf("Erro ao atualizar a lista de servidores: %v\n", err)
				}
			}
		}
	}()

	// Inicia o load balancer para aceitar requisições de clientes
	startLoadBalancer()
}

func startLoadBalancer() {
	listener, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		fmt.Printf("Erro ao iniciar o Load Balancer: %v\n", err)
		return
	}
	defer listener.Close()
	fmt.Printf("Load Balancer rodando na porta %s\n", PORT)

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Erro ao aceitar conexão do cliente: %v\n", err)
			continue
		}
		go handleClientConnection(clientConn)
	}
}

// Consulta o servidor de verificação para obter a lista de servidores
func updateServerList() error {
	conn, err := net.Dial("tcp", VerificationServer)
	if err != nil {
		return fmt.Errorf("Erro ao conectar ao servidor de verificação: %v", err)
	}
	defer conn.Close()

	_, err = fmt.Fprintln(conn, PrefixServers)
	if err != nil {
		return fmt.Errorf("Erro ao enviar requisição ao servidor de verificação: %v", err)
	}

	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("Erro ao ler resposta do servidor de verificação: %v", err)
	}

	// Remove o caractere de nova linha do final da resposta
	response = strings.TrimSpace(response)

	var serverList []ServerInfo
	err = json.Unmarshal([]byte(response), &serverList)
	if err != nil {
		return fmt.Errorf("Erro ao decodificar lista de servidores: %v", err)
	}

	// Atualiza a lista de servidores globalmente
	connPoolMux.Lock()
	defer connPoolMux.Unlock()
	servers = serverList

	fmt.Println("Lista de servidores atualizada:", servers)
	fmt.Println("--------------------")

	return nil
}

func getNextServer() ServerInfo {
	serverIndexMux.Lock()
	defer serverIndexMux.Unlock()

	if len(servers) == 0 {
		return ServerInfo{}
	}

	// Seleciona o próximo servidor na lista
	selectedServer := servers[serverIndex]
	serverIndex = (serverIndex + 1) % len(servers)

	return selectedServer
}

// Redireciona a requisição do cliente para o servidor selecionado e retorna a resposta ao cliente
func handleClientConnection(clientConn net.Conn) {
	defer clientConn.Close()

	server := getNextServer()

	// Verifica se há servidores disponíveis
	if server.IP == "" {
		fmt.Fprintf(clientConn, "Nenhum servidor disponível\n")
		return
	}

	serverAddress := fmt.Sprintf("%s:%d", server.IP, server.Port)

	// Reutiliza uma conexão do pool, se disponível
	serverConn, err := getServerConnection(serverAddress)
	if err != nil {
		fmt.Fprintf(clientConn, "Erro ao conectar ao servidor: %v\n", err)
		return
	}
	defer releaseServerConnection(serverAddress, serverConn)

	// Cria leitores e escritores para as conexões
	clientReader := bufio.NewReader(clientConn)
	serverWriter := bufio.NewWriter(serverConn)
	serverReader := bufio.NewReader(serverConn)
	clientWriter := bufio.NewWriter(clientConn)

	// Envia a requisição do cliente para o servidor
	message, err := clientReader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(clientConn, "Erro ao ler mensagem do cliente: %v\n", err)
		return
	}

	_, err = serverWriter.WriteString(message)
	if err != nil {
		fmt.Fprintf(clientConn, "Erro ao enviar mensagem para o servidor: %v\n", err)
		return
	}
	serverWriter.Flush()

	// Recebe a resposta do servidor e envia de volta ao cliente
	response, err := serverReader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(clientConn, "Erro ao ler resposta do servidor: %v\n", err)
		return
	}
	_, err = clientWriter.WriteString(response)
	if err != nil {
		fmt.Fprintf(clientConn, "Erro ao enviar resposta ao cliente: %v\n", err)
		return
	}
	clientWriter.Flush()

	fmt.Printf("Mensagem redirecionada para o servidor %s e resposta enviada ao cliente\n", serverAddress)
}

// getServerConnection reutiliza ou cria uma nova conexão para o servidor
func getServerConnection(serverAddress string) (net.Conn, error) {
	connPoolMux.Lock()
	defer connPoolMux.Unlock()

	if conn, exists := connPool[serverAddress]; exists {
		return conn, nil
	}

	// Se não houver conexão no pool, cria uma nova
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		return nil, fmt.Errorf("Erro ao conectar ao servidor: %v", err)
	}

	// Adiciona ao pool para reutilização futura
	connPool[serverAddress] = conn
	return conn, nil
}

// releaseServerConnection devolve a conexão ao pool
func releaseServerConnection(serverAddress string, conn net.Conn) {
	connPoolMux.Lock()
	defer connPoolMux.Unlock()

	// Se a conexão ainda estiver ativa, devolve-a ao pool
	if _, exists := connPool[serverAddress]; !exists {
		connPool[serverAddress] = conn
	}
}
