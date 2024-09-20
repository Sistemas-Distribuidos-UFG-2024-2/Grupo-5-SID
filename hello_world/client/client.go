package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

// Função que gerencia a comunicação com o servidor TCP.
func startClient(done chan struct{}, loadBalancerAddress string) {
    // Ticker para enviar mensagens periodicamente.
    ticker := time.NewTicker(time.Second)
    defer ticker.Stop()

    reader := bufio.NewReader(os.Stdin)

    for {
        select {
        case <-done:
            log.Println("Received done signal, stopping client")
            return
        case <-ticker.C:
            conn, err := net.Dial("tcp", loadBalancerAddress)
            if err != nil {
                log.Println("dial:", err)
                continue
            }
            defer conn.Close()

            // input para receber os dados a serem calculados

            fmt.Print("Nome: ")
            nome, _ := reader.ReadString('\n')
            nome = strings.TrimSpace(nome)

            fmt.Print("Cargo: ")
            cargo, _ := reader.ReadString('\n')
            cargo = strings.TrimSpace(cargo)

            fmt.Print("Salário: ")
            salarioStr, _ := reader.ReadString('\n')
            salarioStr = strings.TrimSpace(salarioStr)
            salario, err := strconv.ParseFloat(salarioStr, 64)
            if err != nil {
                log.Println("Erro ao converter salário:", err)
                continue
            }

            // Envia dados para o servidor
            message := fmt.Sprintf("FUNCIONARIO,%s,%s,%.2f", nome, cargo, salario)
            fmt.Fprintf(conn, message+"\n")

            // Recebe resposta
            response, err := bufio.NewReader(conn).ReadString('\n')
            if err != nil {
                log.Println("read:", err)
                continue
            }
            log.Printf("Received: %s", strings.TrimSpace(response))

            time.Sleep(3000 * time.Millisecond)
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
	// ? #### Threads ####
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