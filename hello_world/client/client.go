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

func main() {
	// Canal para capturar sinais de interrupção do sistema (como Ctrl+C).
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Endereço do load balancer
	loadBalancerAddress := "localhost:5611"
	log.Printf("connecting to %s", loadBalancerAddress)

	// Conecta ao load balancer via TCP

	// Canal para sinalizar quando a leitura das mensagens estiver concluída.
	done := make(chan struct{})

	// Ticker para enviar mensagens periodicamente.
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	// Loop principal para gerenciar a comunicação com o servidor TCP.
	for {
		select {
		// Caso o canal 'done' receba um valor, a função retorna, encerrando a execução.
		case <-done:
			return

		// Caso o ticker envie um valor, escrevemos uma mensagem no TCP.
		case <-ticker.C:
			conn, err := net.Dial("tcp", loadBalancerAddress)
			if err != nil {
				log.Fatal("dial:", err)
			}
			defer conn.Close()

			_, err = fmt.Fprintf(conn, "hello\n")
			if err != nil {
				log.Println("write:", err)
				return
			}

			reader := bufio.NewReader(conn)
			message, err := reader.ReadString('\n')
			if err != nil {
				continue
			}
			log.Printf("received: %s", message)
		}
	}
}
