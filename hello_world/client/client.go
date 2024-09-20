package main

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
    // Canal para capturar sinais de interrupção do sistema (como Ctrl+C).
    interrupt := make(chan os.Signal, 1)
    signal.Notify(interrupt, os.Interrupt)

    // URL do servidor WebSocket 
    u := url.URL{Scheme: "ws", Host: "localhost:5611",  Path: "/ws"}
    log.Printf("connecting to %s", u.String())

    // Conecta ao servidor WebSocket.
    c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
    if err != nil {
        log.Fatal("dial:", err)
    }
    defer c.Close()

    // Canal para sinalizar quando a leitura das mensagens estiver concluída.
    done := make(chan struct{})

    // Goroutine para ler mensagens do servidor WebSocket.
    // ? funciona de forma semelhante a threads, mas são mais leves e gerenciadas pelo Go
    go func() {
        defer close(done)
        for {
            _, message, err := c.ReadMessage()
            if err != nil {
                log.Println("read:", err)
                return
            }
            log.Printf("received: %s", message)
        }
    }()

    // Ticker para enviar mensagens periodicamente.
    ticker := time.NewTicker(time.Second)
    defer ticker.Stop()

    // Loop principal para gerenciar a comunicação com o servidor WebSocket.
    for {
        select {
        // Caso o canal 'done' receba um valor, a função retorna, encerrando a execução.
        case <-done:
            return

        // Caso o ticker envie um valor, escrevemos uma mensagem no WebSocket.
        case t := <-ticker.C:
            err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
            if err != nil {
                log.Println("write:", err)
                return
            }

        // Caso o canal 'interrupt' receba um valor (como um sinal de interrupção), fecha a conexão WebSocket.
        case <-interrupt:
            log.Println("interrupt")

            // Envia uma mensagem de fechamento para o WebSocket.
            err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
            if err != nil {
                log.Println("write close:", err)
                return
            }

            // Espera pelo fechamento da conexão ou um timeout de 1 segundo.
            select {
            case <-done:
            case <-time.After(time.Second):
            }
            return
        }
    }
}
