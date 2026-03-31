package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
)

type MensagemTCP struct {
	Tipo string `json:"tipo"`
	ID   string `json:"id"`
	Dado string `json:"dado"`
	Acao string `json:"acao"`
}

func enviar(conn net.Conn, id string, dado string, acao string) {

	msg := MensagemTCP{
		Tipo: "ATUADOR",
		ID:   id,
		Dado: dado,
		Acao: acao,
	}

	data, _ := json.Marshal(msg)

	conn.Write([]byte(string(data) + "\n"))
}

func main() {

	conn, err := net.Dial("tcp", "server:8000")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)

	msg := MensagemTCP{
		Tipo: "ATUADOR",
		ID:   "ATUADOR A",
		Dado: "3",
		Acao: "nil",
	}

	jsonData, _ := json.Marshal(msg)

	conn.Write([]byte(string(jsonData) + "\n"))

	for {

		var msgRec MensagemTCP
		fmt.Println("Conectado")
		buffer, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Servidor desconectado")
			return
		}

		err = json.Unmarshal([]byte(buffer), &msgRec)
		if err != nil {
			fmt.Println("Erro ao Ler Dados do Json")
			continue
		}

		fmt.Println("Recebido:", msgRec.Tipo)
		fmt.Println("Recebido:", msgRec.ID)
		fmt.Println("Recebido:", msgRec.Dado)
		fmt.Println("Recebido:", msgRec.Acao)
	}
}
