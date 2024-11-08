package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

const (
	PORT                = "5610"
	PrefixHealth        = "/health/"
	PrefixServers       = "/servers/"
	HealthCheckInterval = 3 * time.Second
)

type ServerInfo struct {
	IP      string `json:"ip"`
	Port    int    `json:"port"`
	Enabled bool   `json:"enabled"`
}

var (
	servers     = map[int]ServerInfo{}
	connPool    = make(map[string]net.Conn) // Pool de conexões
	connPoolMux = make(chan struct{}, 1)    // Mutex simples com canal
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

	if !serverInfo.Enabled {
		fmt.Printf("Servidor avisou que está indisponível: %d", serverInfo.Port)
		delete(servers, serverInfo.Port)
		removeConnFromPool(serverInfo.IP, serverInfo.Port) // Remove a conexão do pool
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

	fmt.Fprintln(conn, string(serverList))
	return nil
}

func healthCheckServers() {
	ticker := time.NewTicker(HealthCheckInterval) // Cria o ticker com o intervalo de tempo definido
	defer ticker.Stop()                           // Para o ticker quando a função for encerrada

	for range ticker.C { // Executa a cada intervalo definido por HealthCheckInterval
		if len(servers) == 0 {
			fmt.Println("Nenhum servidor registrado para verificação.")
			continue
		}

		for key, server := range servers {
			address := fmt.Sprintf("%s:%d", server.IP, server.Port)

			// Conexão direta com o servidor para health check
			conn, err := net.Dial("tcp", address) // Cria uma nova conexão para o servidor
			if err != nil {
				fmt.Printf("Servidor %s:%d não está acessível. Removendo do registro.\n", server.IP, server.Port)
				delete(servers, key) // Remove o servidor do mapa diretamente
				continue
			}

			serverReader := bufio.NewReader(conn)
			serverWriter := bufio.NewWriter(conn)
			fmt.Fprintf(serverWriter, "%s\n", PrefixHealth)
			serverWriter.Flush() // Envia o comando de health check

			_, err = serverReader.ReadString('\n')
			if err != nil {
				fmt.Printf("Erro ao verificar o servidor %s:%d: %v. Removendo do registro.\n", server.IP, server.Port, err)
				delete(servers, key) // Remove o servidor do mapa diretamente
			} else {
				server.Enabled = true
				servers[key] = server
				fmt.Printf("Servidor %s:%d está saudável.\n", server.IP, server.Port)
			}
			conn.Close() // Fecha a conexão após a verificação
		}
	}
}

// Função para criar ou reutilizar uma conexão
func getOrCreateConnection(address string) (net.Conn, error) {
	// Bloqueia o pool de conexões
	connPoolMux <- struct{}{}
	defer func() { <-connPoolMux }()

	// Verifica se a conexão já existe no pool
	if conn, ok := connPool[address]; ok {
		// Verifica se a conexão ainda está ativa
		if isConnAlive(conn) {
			return conn, nil
		}
		// Se a conexão não estiver ativa, remove-a do pool
		conn.Close()
		delete(connPool, address)
	}

	// Cria uma nova conexão se não houver conexão reutilizável
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	// Adiciona a nova conexão ao pool
	connPool[address] = conn
	return conn, nil
}

// Função para verificar se uma conexão está viva
func isConnAlive(conn net.Conn) bool {
	conn.SetReadDeadline(time.Now().Add(15 * time.Second))
	one := []byte{}
	// Tentativa de leitura não bloqueante
	if _, err := conn.Read(one); err == nil || err == io.EOF {
		// Se a conexão está ativa ou EOF, ainda podemos usar a conexão
		conn.SetReadDeadline(time.Time{}) // Remova o deadline após o teste
		return true
	}
	return false
}

// Função para remover uma conexão do pool
func removeConnFromPool(ip string, port int) {
	address := fmt.Sprintf("%s:%d", ip, port)

	connPoolMux <- struct{}{}
	defer func() { <-connPoolMux }()

	if conn, ok := connPool[address]; ok {
		conn.Close()
		delete(connPool, address)
	}
}
