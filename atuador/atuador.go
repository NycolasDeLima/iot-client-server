package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
)

var (
	ativo  bool
	estado string
	errCon bool
)

type MensagemTCP struct {
	Tipo string `json:"tipo"`
	ID   string `json:"id"`
	Dado string `json:"dado"`
	Acao string `json:"acao"`
}

func enviar(conn net.Conn, id string, dado string, acao string) error {

	msg := MensagemTCP{
		Tipo: "ATUADOR",
		ID:   id,
		Dado: dado,
		Acao: acao,
	}

	data, _ := json.Marshal(msg)

	_, err := conn.Write([]byte(string(data) + "\n"))
	if err != nil {
		return err
	}

	return nil
}

func conectar() net.Conn {
	for {
		conn, err := net.Dial("tcp", "server:8000")
		if err != nil {
			fmt.Println("Conectando...")
			time.Sleep(1 * time.Second)
			continue
		}

		fmt.Println("Conectado ao servidor.")
		return conn
	}
}

func main() {

	var msg MensagemTCP
	handlers := map[string]func(string){
		"alarme": tratarAlarme,
		"vmi":    tratarVMI,
	}

	msg.Tipo = "SENSOR"

	if len(os.Args) < 3 {
		fmt.Println("Uso: go run main.go <tipoAtuador> <id>")
		return
	}

	tipo := os.Args[1]
	id := os.Args[2]

	switch tipo {

	case "alarme":
		ativo = false
		estado = "DESLIGADO"

	case "vmi":
		ativo = false
		estado = "DESLIGADA"

	default:
		fmt.Println("Tipo inválido! Use alarme ou vmi")
		return
	}

	conn := conectar()
	defer conn.Close()

	fmt.Println("Conectado ao servidor: ", conn.LocalAddr().String())

	reader := bufio.NewReader(conn)

	id = tipo + "_" + id

	err := enviar(conn, id, estado, "nil")
	if err != nil {
		errCon = true
	}

	buffer, err := reader.ReadString('\n')
	if err != nil {
		errCon = true
	}

	err = json.Unmarshal([]byte(buffer), &msg)
	if err != nil {
		fmt.Println("Erro ao Ler Dados do Json")
		errCon = true
	}

	if msg.Dado != "ATUADOR CONECTADO" {
		fmt.Println("Conexão mal estabelecida")
		errCon = true
	}

	for {

		if errCon {
			fmt.Println("Servidor Desconectado. Tentando Reconexão")
			conn.Close()
			conn = conectar()
			reader = bufio.NewReader(conn)

			err := enviar(conn, id, estado, "nil")
			if err != nil {
				continue
			}

			buffer, err := reader.ReadString('\n')
			if err != nil {
				continue
			}

			err = json.Unmarshal([]byte(buffer), &msg)
			if err != nil {
				fmt.Println("Erro ao Ler Dados do Json")
				continue
			}

			if msg.Dado != "ATUADOR CONECTADO" {
				fmt.Println("Conexão mal estabelecida")
				continue
			}

			errCon = false
			continue
		}

		buffer, err := reader.ReadString('\n')
		if err != nil {
			errCon = true
			continue
		}

		err = json.Unmarshal([]byte(buffer), &msg)
		if err != nil {
			fmt.Println("Erro ao Ler Dados do Json")
			continue
		}

		fmt.Println("Recebido:", msg.Dado)

		if handler, ok := handlers[tipo]; ok {
			handler(msg.Dado)
		}

		err = enviar(conn, tipo, estado, "nil")
		if err != nil {
			errCon = true
			continue
		}

		if ativo {
			fmt.Println(estado)
		}
	}
}
