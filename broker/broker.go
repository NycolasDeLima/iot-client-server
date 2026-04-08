package main

import (
	"encoding/json"
	"log"
	"net"
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

	go servidorUDP()
	go servidorTcp()

	log.Println("Broker iniciado (UDP:8080 | TCP:8000)")
	select {}
}

// =============      Broker UDP     =============

func servidorUDP() {

	addrUDP, err := net.ResolveUDPAddr("udp", ":5000")
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

func servidorTcp() {

	listenner, err := net.Listen("tcp", ":8000")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listenner.Accept()
		if err != nil {
			log.Println("Erro:", err)
			continue
		}

		go tratarConexaoTcp(conn)
	}
}
