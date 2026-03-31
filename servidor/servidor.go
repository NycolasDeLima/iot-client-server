package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"
)

// ============= variavel global ===============

var (
	atuadores = make(map[string]net.Conn)
	clientes  = make(map[string]net.Conn)
	sensores  = make(map[string]Sensor)
	mutex     sync.Mutex
)

// =============     struct      ====================

type MensagemTCP struct {
	Tipo string `json:"tipo"`
	ID   string `json:"id"`
	Dado string `json:"dado"`
	Acao string `json:"acao"`
}

type MensagemUDP struct {
	Tipo string `json:"tipo"`
	Dado string `json:"dado"`
}

type Sensor struct {
	Tipo        string `json:"tipo"`
	Dado        string `json:"dado"`
	UltimoVisto time.Time
	Online      bool
}

// ========= tratar ============

func tratarSensor(msg MensagemUDP, clientAddr *net.UDPAddr) {
	fmt.Println("Trantando Sensor: ", clientAddr)

	mutex.Lock()

	defer mutex.Unlock()

	sensores[msg.Tipo] = Sensor{
		Tipo:        msg.Tipo,
		Dado:        msg.Dado,
		UltimoVisto: time.Now(),
		Online:      true,
	}

}

func tratarCliente(id string, conn net.Conn, reader *bufio.Reader) {
	defer conn.Close()

	for {

		buffer, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Cliente", id, "desconectado")
			return
		}

		var msg MensagemTCP

		err = json.Unmarshal([]byte(buffer), &msg)
		if err != nil {
			fmt.Println("Erro JSON cliente ", id, ":", err)
			continue
		}

		if msg.Acao == "ACAO ATUADOR" {

			mutex.Lock()
			atuadorConn, existe := atuadores[msg.ID]
			mutex.Unlock()

			if !existe {
				conn.Write([]byte("Atuador não encontrado\n"))
				continue
			}

			data, _ := json.Marshal(msg)

			_, err := atuadorConn.Write([]byte(string(data) + "\n"))
			if err != nil {
				fmt.Println("Erro ao enviar cliente ", id, ":", err)
				continue
			}

			conn.Write([]byte("Ação enviada com sucesso\n"))

		} else if msg.Acao == "LISTAR SENSORES" {

			var lista []string

			for id := range sensores {
				lista = append(lista, id)
			}

			fmt.Println(lista)

			data, _ := json.Marshal(lista)
			msgEnv := MensagemTCP{
				Tipo: "Servidor",
				ID:   "nil",
				Dado: string(data),
				Acao: "LISTAR SENSORES",
			}

			data, _ = json.Marshal(msgEnv)

			conn.Write([]byte(string(data) + "\n"))

		} else if msg.Acao == "LISTAR ATUADORES" {

			var lista []string

			for id := range atuadores {
				lista = append(lista, id)
			}

			fmt.Println(lista)

			data, _ := json.Marshal(lista)
			msgEnv := MensagemTCP{
				Tipo: "Servidor",
				ID:   "nil",
				Dado: string(data),
				Acao: "LISTAR ATUADORES",
			}

			data, _ = json.Marshal(msgEnv)

			conn.Write([]byte(string(data) + "\n"))

		}
	}
}

func tratarAtuador(id string, conn net.Conn, reader *bufio.Reader) {

	defer func() {
		conn.Close()
		removerAtuador(id)
	}()

	fmt.Println("Tratando atuador: ", id)

	for {
		buffer, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Atuador desconectado: ", id)
			return
		}

		var msg MensagemTCP

		err = json.Unmarshal([]byte(buffer), &msg)
		if err != nil {
			fmt.Println("Erro ao Ler Dados do Json")
			continue
		}

		fmt.Println("Recebido:", msg.Acao)
	}
}

// =============== tratar conexoes ================

func tratarConexaoTcp(conn net.Conn) {

	var msg MensagemTCP

	addr := conn.RemoteAddr().String()
	fmt.Println("Cliente conectado: ", addr)

	reader := bufio.NewReader(conn)

	buffer, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(buffer), &msg)
	if err != nil {
		fmt.Println("JSON inválido:", err)
		return
	}

	if msg.Tipo == "ATUADOR" {
		id := msg.ID

		mutex.Lock()
		atuadores[id] = conn
		mutex.Unlock()

		go tratarAtuador(id, conn, reader)
	}

	if msg.Tipo == "CLIENTE" {
		id := msg.ID

		mutex.Lock()
		clientes[id] = conn
		mutex.Unlock()

		go tratarCliente(id, conn, reader)
	}

}

// ============ utilitario =============

func removerAtuador(id string) {
	mutex.Lock()
	defer mutex.Unlock()

	delete(atuadores, id)

	fmt.Println("Atuador: ", id, "removido")
}

func removerSensor() {
	for {
		time.Sleep(5 * time.Second)

		for id, sensor := range sensores {
			if time.Since(sensor.UltimoVisto) > 5*time.Second {
				fmt.Println("Sensor desconectado: ", id)
				delete(sensores, id)
			}
		}
	}
}

func main() {
	go servidorUDP()
	go servidorTcp()

	fmt.Println("Servidor iniciado (UDP:8080 | TCP:8000)")
	select {}
}

func servidorUDP() {

	addrUDP, err := net.ResolveUDPAddr("udp", ":8080")
	if err != nil {
		panic(err)
	}

	connUDP, err := net.ListenUDP("udp", addrUDP)
	if err != nil {
		panic(err)
	}

	defer connUDP.Close()

	fmt.Println("Servidor UDP rodando na porta 8080...")

	buffer := make([]byte, 1024)

	go removerSensor()

	for {

		var msg MensagemUDP
		n, clientAddr, err := connUDP.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Erro:", err)
			continue
		}

		err = json.Unmarshal(buffer[:n], &msg)
		if err != nil {
			fmt.Println("JSON inválido:", err)
			continue
		}

		go tratarSensor(msg, clientAddr)

	}
}

func servidorTcp() {

	listenner, err := net.Listen("tcp", ":8000")
	if err != nil {
		panic(err)
	}

	fmt.Println("Servidor rodando na por 5000")

	for {
		conn, err := listenner.Accept()
		if err != nil {
			fmt.Println("Erro:", err)
			continue
		}

		go tratarConexaoTcp(conn)
	}
}
