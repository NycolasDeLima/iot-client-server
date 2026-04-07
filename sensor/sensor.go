package main

import (
	//"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"

	//"strings"
	"time"
)

type Mensagem struct {
	Tipo string `json:"tipo"`
	ID   string `json:"id"`
	Dado string `json:"dado"`
}

func main() {

	//input := bufio.NewReader(os.Stdin)

	serverAddrUDP, err := net.ResolveUDPAddr("udp", "server:8080")
	if err != nil {
		panic(err)
	}

	conn, err := net.DialUDP("udp", nil, serverAddrUDP)
	if err != nil {
		panic(err)
	}

	var (
		msg      Mensagem
		dado     int
		estado   string
		handlers = map[string]func(int, string) int{
			"bpm":  ajustarBPM,
			"spo2": ajustarSpO2,
		}
	)

	defer conn.Close()

	msg.Tipo = "SENSOR"

	if len(os.Args) < 3 {
		fmt.Println("Uso: go run main.go <tipoSensor> <id>")
		return
	}

	tipoSensor := os.Args[1]
	id := os.Args[2]

	if tipoSensor != "bpm" && tipoSensor != "spo2" {
		fmt.Println("Tipo inválido! Use bpm ou spo2")
		return
	}

	msg.ID = tipoSensor + "_" + id

	switch tipoSensor {

	case "bpm":
		dado = 75
		estado = "repouso"

	case "spo2":
		dado = 98
		estado = "normal"
	}

	for {

		if rand.Float64() < 0.01 {
			estado = mudarEstado(tipoSensor)
			fmt.Println("Mudou estado para:", estado)
		}

		if handler, ok := handlers[tipoSensor]; ok {
			dado = handler(dado, estado)
		}

		msg.Dado = strconv.Itoa(dado)

		jsondata, _ := json.Marshal(msg)

		conn.Write([]byte(jsondata))

		time.Sleep(100 * time.Millisecond)
	}

}
