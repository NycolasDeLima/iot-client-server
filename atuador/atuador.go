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

func conectar(serverIP string) net.Conn {
	for {
		conn, err := net.Dial("tcp", serverIP)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		return conn
	}
}

func main() {

	var msg MensagemTCP
	handlers := map[string]func(string){
		"alarme": tratarAlarme,
		"vmi":    tratarVMI,
	}

	if len(os.Args) < 4 {
		fmt.Println("Uso: go run main.go <tipoSensor> <id> <serverIP>")
		return
	}

	tipo := os.Args[1]
	id1 := os.Args[2]
	serverIP := os.Args[3]

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

	conn := conectar(serverIP)
	defer conn.Close()

	fmt.Println("Conectado ao servidor: ", conn.LocalAddr().String())

	reader := bufio.NewReader(conn)

	id := tipo + "_" + id1

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

		conn.SetReadDeadline(time.Now().Add(1 * time.Minute))

		if errCon {
			exibirAtuador(tipo, id1, estado, ativo, errCon)
			conn.Close()
			conn = conectar(serverIP)
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

		exibirAtuador(tipo, id1, estado, ativo, errCon)

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

		if handler, ok := handlers[tipo]; ok {
			handler(msg.Dado)
		}

		err = enviar(conn, tipo, estado, "nil")
		if err != nil {
			errCon = true
			continue
		}

	}
}
