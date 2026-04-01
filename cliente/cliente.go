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

// ================= structs ====================

type MensagemTCP struct {
	Tipo string `json:"tipo"`
	ID   string `json:"id"`
	Dado string `json:"dado"`
	Acao string `json:"acao"`
}

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

// ======================== functions ============

func enviar(conn net.Conn, id string, dado string, acao string) {

	msg := MensagemTCP{
		Tipo: "CLIENTE",
		ID:   id,
		Dado: dado,
		Acao: acao,
	}

	data, _ := json.Marshal(msg)

	conn.Write([]byte(string(data) + "\n"))
}

func ler(reader *bufio.Reader, msg MensagemTCP) error {
	buffer, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Erro no recebimento da Resposta")
		return err
	}

	err = json.Unmarshal([]byte(buffer), &msg)
	if err != nil {
		fmt.Println("Erro ao Ler Dados do Json")
		return err
	}

	return nil
}

// =============== main

func main() {

	var id string

	fmt.Print("Digite um id: ")
	fmt.Scanln(&id)

	fmt.Println("Conectando ao Servidor...")

	conn, err := net.Dial("tcp", "server:8000")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)

	msg := MensagemTCP{
		Tipo: "CLIENTE",
		ID:   id,
		Dado: "2",
		Acao: "nil",
	}

	data, _ := json.Marshal(msg)

	input := bufio.NewReader(os.Stdin)

	conn.Write(append(data, '\n'))

	for {
		fmt.Println("\n===== MENU =====")
		fmt.Println("1 - Visualizar sensores")
		fmt.Println("2 - Visualizar dados do sensor")
		fmt.Println("3 - Visualizar atuadores")
		fmt.Println("4 - Enviar comando")
		fmt.Println("5 - Sair")
		fmt.Print("Escolha: ")

		opcao, _ := input.ReadString('\n')
		opcao = strings.TrimSpace(opcao)

		switch opcao {

		case "1":
			enviar(conn, "nil", "nil", "LISTAR SENSORES")

			fmt.Println("Aguardando mensagem do servidor...")

			buffer, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Erro no recebimento da Resposta")
				continue
			}

			err = json.Unmarshal([]byte(buffer), &msg)
			if err != nil {
				fmt.Println("Erro ao Ler Dados do Json")
				continue
			}

			// var lista []string

			var lista map[string]Sensor

			err = json.Unmarshal([]byte(msg.Dado), &lista)
			if err != nil {
				fmt.Println("Erro ao Ler Dados da Lista")
				continue
			}

			//fmt.Printf("%-5s | %-15s | %-15s | %-10s\n", "ID", "Nome", "Tipo", "Status")
			fmt.Printf("%-5s | %-15s | %-15s | %-10s\n", "ID", "Tipo", "Dado", "UltimoVisto")
			fmt.Println("----------------------------------------------------------")

			for _, sensor := range lista {
				fmt.Printf("%-5s | %-15s | %-15s | %-10s\n",
					sensor.ID, sensor.Tipo, sensor.Dado, sensor.UltimoVisto.Format("15:04:05"),
				)
			}

		case "2":

			fmt.Println("\nDigite o id do Sensor: ")
			id, _ := input.ReadString('\n')
			id = strings.TrimSpace(id)

			enviar(conn, id, "nil", "VER DADO SENSOR")

			fmt.Println("\nLendo dados do sensor... Aperte ENTER para sair")

			stopChan := make(chan bool)

			// 🔹 Goroutine que recebe dados
			go func() {
				for {
					select {
					case <-stopChan:
						return
					default:
						buffer, err := reader.ReadString('\n')
						if err != nil {
							fmt.Println("Erro ao receber")
							continue
						}

						var msg MensagemTCP
						err = json.Unmarshal([]byte(buffer), &msg)
						if err != nil {
							fmt.Println("Erro ao Ler Mensagem do Servidor")
							continue
						}

						if msg.Dado == "SENSOR NÃO ENCONTRADO" {
							fmt.Println("Sensor não encontrado")
							continue
						}

						var sensor Sensor
						err = json.Unmarshal([]byte(msg.Dado), &sensor)
						if err != nil {
							fmt.Println("Erro ao Ler Dados do sensor")
							continue
						}

						fmt.Printf("Sensor %s | Dado: %s | Hora: %s\n",
							sensor.ID,
							sensor.Dado,
							sensor.UltimoVisto.In(time.Local).Format("15:04:05"),
						)
					}
				}
			}()

			input.ReadString('\n')

			close(stopChan)

			enviar(conn, id, "nil", "REMOVER INSCRITO")

			fmt.Println("Leitura de sensor finalizada")

		case "3":

			enviar(conn, "nil", "nil", "LISTAR ATUADORES")

			fmt.Println("Aguardando mensagem do servidor...")

			buffer, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Erro no recebimento da Resposta")
				continue
			}

			err = json.Unmarshal([]byte(buffer), &msg)
			if err != nil {
				fmt.Println("Erro ao Ler Dados do Json")
				continue
			}

			// var lista []string

			var lista map[string]Atuador

			err = json.Unmarshal([]byte(msg.Dado), &lista)
			if err != nil {
				fmt.Println("Erro ao Ler Dados da Lista")
				continue
			}

			//fmt.Printf("%-5s | %-15s | %-15s | %-10s\n", "ID", "Nome", "Tipo", "Status")
			fmt.Printf("%-5s | %-15s | %-15s\n", "ID", "Tipo", "Status")
			fmt.Println("----------------------------------------------------------")

			for _, atuador := range lista {
				fmt.Printf("%-5s | %-15s | %-15s\n",
					atuador.ID, atuador.Tipo, atuador.Status,
				)
			}

		case "4":

			fmt.Println("\nDigite o id do Atuador: ")
			id, _ := input.ReadString('\n')
			id = strings.TrimSpace(id)

			fmt.Println("\nDigite a ação: ")
			dado, _ := input.ReadString('\n')
			dado = strings.TrimSpace(dado)

			enviar(conn, id, dado, "ACAO ATUADOR")

		}

	}

}
