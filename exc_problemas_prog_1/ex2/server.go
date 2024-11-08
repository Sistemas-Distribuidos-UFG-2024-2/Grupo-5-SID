package main

import (
	"fmt"
	"net"
	"net/rpc"
)

type Pessoa struct {
	Nome  string
	Sexo  string
	Idade int
}

type Resultado struct {
	Mensagem string
}

type ServicoMaioridade struct{}

func (s *ServicoMaioridade) VerificaMaioridade(pessoa *Pessoa, resultado *Resultado) error {
	if (pessoa.Sexo == "M" && pessoa.Idade >= 18) || (pessoa.Sexo == "F" && pessoa.Idade >= 21) {
		resultado.Mensagem = fmt.Sprintf("%s já atingiu a maioridade.", pessoa.Nome)
	} else {
		resultado.Mensagem = fmt.Sprintf("%s não atingiu a maioridade.", pessoa.Nome)
	}
	return nil
}

func main() {
	servico := new(ServicoMaioridade)
	rpc.Register(servico)
	
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Servidor RPC iniciado na porta 1234")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Erro ao aceitar conexão:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
