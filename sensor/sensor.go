package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
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
	var msg Mensagem

	defer conn.Close()

	msg.Tipo = "SENSOR"

	fmt.Println("Digite um número: ")

	msg.Dado, _ = input.ReadString('\n')
	msg.Dado = strings.TrimSpace(msg.Dado)

	fmt.Println("Digite um id: ")
	msg.Tipo, _ = input.ReadString('\n')
	msg.Tipo = strings.TrimSpace(msg.Tipo)

	for {
		time.Sleep(2 * time.Second)

		jsondata, _ := json.Marshal(msg)

		conn.Write([]byte(jsondata))
	}

}
