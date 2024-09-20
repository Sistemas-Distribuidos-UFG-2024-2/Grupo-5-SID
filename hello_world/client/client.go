package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

// Função que gerencia a comunicação com o servidor TCP.
func startClient(done chan struct{}, loadBalancerAddress string) {
    // Ticker para enviar mensagens periodicamente.
    ticker := time.NewTicker(time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-done:
            log.Println("Received done signal, stopping client")
            return
        case <-ticker.C:
            // log.Println("Attempting to connect to server")
            conn, err := net.Dial("tcp", loadBalancerAddress)
            if err != nil {
                log.Println("dial:", err)
                continue
            }
            defer conn.Close()

            // log.Println("Connected to server, sending message")
            _, err = fmt.Fprintf(conn, "hello\n")
            if err != nil {
                log.Println("write:", err)
                continue
            }
            // log.Println("Message sent successfully")

            // Recebe a resposta do servidor
            reader := bufio.NewReader(conn)
            message, err := reader.ReadString('\n')
            if err != nil {
                log.Println("read:", err)
                continue
            }
            log.Printf("Received: %s", message)
        }
    }
}

func main() {
    // Canal para capturar sinais de interrupção do sistema (como Ctrl+C).
    interrupt := make(chan os.Signal, 1)
    signal.Notify(interrupt, os.Interrupt)

    // Endereço do load balancer
    loadBalancerAddress := "localhost:5611"
    log.Printf("connecting to %s", loadBalancerAddress)

    // Canal para sinalizar quando a leitura das mensagens estiver concluída.
    done := make(chan struct{})

    // Inicia a goroutine para gerenciar a comunicação com o servidor TCP.
    go func() {
        log.Println("starting client")
        startClient(done, loadBalancerAddress)
    }()

    // Aguarda um sinal de interrupção (como Ctrl+C).
    <-interrupt

    // Fecha o canal 'done' para sinalizar a goroutine para parar.
    close(done)

    log.Println("client stopped")
}