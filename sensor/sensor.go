package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type Mensagem struct {
	Tipo string `json:"tipo"`
	ID   string `json:"id"`
	Dado string `json:"dado"`
}

func main() {

	input := bufio.NewReader(os.Stdin)

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
			"bpm": ajustarBPM,
			//"admin":  tratarAdmin,
			//"user":   tratarUser,
		}
	)

	defer conn.Close()

	msg.Tipo = "SENSOR"

	fmt.Println("Digite o tipo de Sensor: ")

	tipoSensor, _ := input.ReadString('\n')
	tipoSensor = strings.TrimSpace(tipoSensor)

	fmt.Println("\nDigite um id para o Sensor: ")
	id, _ := input.ReadString('\n')
	id = strings.TrimSpace(id)

	msg.ID = tipoSensor + "_" + id

	switch tipoSensor {

	case "bpm":
		dado = 75
		estado = "repouso"
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
