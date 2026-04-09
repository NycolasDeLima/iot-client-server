package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

// ============= variavel global ===============

var (
	atuadoresConn = make(map[string]net.Conn)
	clientes      = make(map[string]net.Conn)

	sensores  = make(map[string]Sensor)
	inscritos = make(map[string][]net.Conn)
	atuadores = make(map[string]Atuador)

	mutex sync.RWMutex
)

// =============     Protocolo de Comunicação      =============

type MensagemTCP struct {
	Tipo string `json:"tipo"`
	ID   string `json:"id"`
	Dado string `json:"dado"`
	Acao string `json:"acao"`
}

type MensagemUDP struct {
	Tipo string `json:"tipo"`
	ID   string `json:"id"`
	Dado string `json:"dado"`
}

// =============     Structs      =============

type Sensor struct {
	Tipo        string `json:"tipo"`
	ID          string `json:"id"`
	Dado        string `json:"dado"`
	UltimoVisto time.Time
}

type Atuador struct {
	Tipo   string `json:"tipo"`
	ID     string `json:"id"`
	Status string `json:"status"`
}

// =============     Broker      =============

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Uso: go run main.go <portUDP> <portTCP>")
		return
	}

	portUDP := ":" + os.Args[1]
	portTCP := ":" + os.Args[2]

	go servidorUDP(portUDP)
	go servidorTcp(portTCP)

	log.Println("Broker iniciado (UDP" + portUDP + " | " + "TCP" + portTCP + ")")
	select {}
}

// =============      Broker UDP     =============

func servidorUDP(portUDP string) {

	addrUDP, err := net.ResolveUDPAddr("udp", portUDP)
	if err != nil {
		panic(err)
	}

	connUDP, err := net.ListenUDP("udp", addrUDP)
	if err != nil {
		panic(err)
	}

	defer connUDP.Close()

	buffer := make([]byte, 1024)

	go removerSensor()

	for {

		var msg MensagemUDP
		n, clientAddr, err := connUDP.ReadFromUDP(buffer)
		if err != nil {
			log.Println("Erro:", err)
			continue
		}

		err = json.Unmarshal(buffer[:n], &msg)
		if err != nil {
			log.Println("JSON inválido:", err)
			continue
		}

		go tratarSensor(msg, clientAddr)

	}
}

// =============      Broker TCP     =================

func servidorTcp(portTCP string) {

	listenner, err := net.Listen("tcp", portTCP)
	if err != nil {
		panic(err)
	}

	log.Println("TCP Addr: " + listenner.Addr().String())

	for {
		conn, err := listenner.Accept()
		if err != nil {
			log.Println("Erro:", err)
			continue
		}

		go tratarConexaoTcp(conn)
	}
}
