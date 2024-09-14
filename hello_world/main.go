package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
)

func main() {
	ln, err := net.Listen("tcp", ":5602")
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
		return
	}
	defer ln.Close()
	fmt.Println("Servidor 2 na porta 5602")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			println(sig)
			println("MORRI")
			os.Exit(1)
		}
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Erro de conexão no Servidor 2:", err)
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
		if strings.EqualFold(message, "hello") {
			fmt.Fprintln(conn, "world")
			fmt.Println("Mensagem recebida:", message)
		}
	}
}
