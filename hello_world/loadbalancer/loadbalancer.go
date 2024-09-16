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
	servers        []ServerInfo // Lista de servidores ativos
	serverIndex    int          // Índice do servidor para distribuição
	serverIndexMux sync.Mutex   // Mutex para evitar condições de corrida no acesso ao índice do servidor
)

func main() {
	// Atualiza periodicamente a lista de servidores
	go func() {
		for {
			err := updateServerList()
			if err != nil {
				fmt.Printf("Erro ao atualizar a lista de servidores: %v\n", err)
			}
			// Aguarda antes da próxima atualização
			time.Sleep(HealthCheckInterval)
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

	fmt.Println("Lista de servidores:")
	fmt.Println(serverList)
	fmt.Println("--------------------")

	servers = serverList
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

	for i := 0; i < len(servers); i++ {
		server := getNextServer()

		// Verifica se há servidores disponíveis
		if server.IP == "" {
			return
		}

		serverAddress := fmt.Sprintf("%s:%d", server.IP, server.Port)
		serverConn, err := net.Dial("tcp", serverAddress)
		if err != nil {
			fmt.Fprintf(clientConn, "Erro ao conectar ao servidor: %v\n", err)
			return
		}
		defer serverConn.Close()

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

		fmt.Printf("Mensagem retornada, enviada ao servidor: %v\n", server.Port)
	}
}
