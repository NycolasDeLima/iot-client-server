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

// ============= Constantes ===============

const (
	ListarSensores  = "LISTAR SENSORES"
	ListarAtuadores = "LISTAR ATUADORES"
	AcaoAtuador     = "ACAO ATUADOR"
	VerDadoSensor   = "VER DADO SENSOR"
	RemoverInscrito = "REMOVER INSCRITO"
)

// ================= Protocolo de Comunicação ====================

type MensagemTCP struct {
	Tipo string `json:"tipo"`
	ID   string `json:"id"`
	Dado string `json:"dado"`
	Acao string `json:"acao"`
}

func enviar(conn net.Conn, id string, dado string, acao string) error {

	msg := MensagemTCP{
		Tipo: "CLIENTE",
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

// ============= Structs ===============

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

// ============= Conecta com o Servidor ===============

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

func limparTela() {
	fmt.Print("\033[H\033[2J")
}

// ============= Cliente ===============

func main() {

	loc, _ := time.LoadLocation("America/Sao_Paulo")

	var (
		idCliente   string
		tipoAtuador string
		tipoSensor  string
		msg         MensagemTCP
		errCon      bool = false
	)

	if len(os.Args) < 3 {
		fmt.Println("Uso: go run main.go <id> <serverIP>")
		return
	}

	idCliente = os.Args[1]
	serverIP := os.Args[2]

	idCliente = "CLIENTE_" + idCliente

	fmt.Println("Conectando ao Servidor...")

	conn := conectar(serverIP)
	defer conn.Close()

	reader := bufio.NewReader(conn)

	err := enviar(conn, idCliente, "nil", "nil")
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

	if msg.Dado != "CLIENTE CONECTADO" {
		fmt.Println("Conexão mal estabelecida")
		errCon = true
	}

	input := bufio.NewReader(os.Stdin)

	for {


		if errCon {
			limparTela()

			fmt.Println("----------------------------------------------------------")
			fmt.Println("        Servidor Desconectado. Tentando Reconexão...      ")
			fmt.Println("----------------------------------------------------------")
			conn.Close()
			conn = conectar(serverIP)
			reader = bufio.NewReader(conn)

			err := enviar(conn, idCliente, "nil", "nil")
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

			if msg.Dado != "CLIENTE CONECTADO" {
				fmt.Println("Conexão mal estabelecida")
				continue
			}

			fmt.Println("----------------------------------------------------------")
			fmt.Println("            Conectado ao Servidor com Sucesso!            ")
			fmt.Println("----------------------------------------------------------")
			input.ReadString('\n')

			errCon = false
			continue

		}
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
			err = enviar(conn, "nil", "nil", ListarSensores)
			if err != nil {
				errCon = true
				continue
			}

			fmt.Println("\nAguardando mensagem do servidor...")

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

			var lista []Sensor

			err = json.Unmarshal([]byte(msg.Dado), &lista)
			if err != nil {
				fmt.Println("Erro ao Ler Dados da Lista")
				continue
			}

			var sensorID string
			var sensorTipo string

			fmt.Printf("\n%-10s | %-15s | %-15s | %-10s\n", "Tipo", "ID", "Dado", "UltimoVisto")
			fmt.Println("----------------------------------------------------------")

			for _, sensor := range lista {

				idx := strings.Index(sensor.ID, "_")
				if idx != -1 {
					sensorTipo = sensor.ID[:idx]
					sensorID = sensor.ID[idx+1:]
				}

				fmt.Printf("%-10s | %-15s | %-15s | %-10s\n",
					sensorTipo, sensorID, sensor.Dado, sensor.UltimoVisto.Format("15:04:05"),
				)
			}

		case "2":

			fmt.Println("\n===== TIPO DO Sensor =====")
			fmt.Println("1 - bpm")
			fmt.Println("2 - SpO2")
			fmt.Print("Escolha o tipo do Sensor: ")
			tipoSensor, _ = input.ReadString('\n')
			tipoSensor = strings.TrimSpace(tipoSensor)

			switch tipoSensor {

			case "1":
				tipoSensor = "bpm"

			case "2":
				tipoSensor = "spo2"

			default:
				fmt.Println("\nOpção inválida!")
				continue
			}

			fmt.Println("\nDigite o id do Sensor: ")
			dado, _ := input.ReadString('\n')
			dado = tipoSensor + "_" + strings.TrimSpace(dado)

			err = enviar(conn, idCliente, dado, VerDadoSensor)
			if err != nil {
				errCon = true
				continue
			}

			fmt.Println("\nAguardando mensagem do servidor...")

			stopChan := make(chan bool)
			done := make(chan bool)

			// 🔹 Goroutine que recebe dados
			go func() {

				defer func() { done <- true }()
				for {
					select {
					case <-stopChan:
						return
					default:
						buffer, err := reader.ReadString('\n')
						if err != nil {
							fmt.Println("Servidor Desconectado... Aperte ENTER para sair")
							errCon = true
							return
						}

						var msg MensagemTCP
						err = json.Unmarshal([]byte(buffer), &msg)
						if err != nil {
							fmt.Println("Erro ao Ler Mensagem do Servidor")
							continue
						}

						switch msg.Dado {
						case "SENSOR NÃO ENCONTRADO":
							fmt.Println("Sensor não encontrado... Aperte ENTER para sair")
							return
						case "SENSOR DESCONECTADO":
							fmt.Println("Sensor Desconectado... Aperte ENTER para sair")
							return

						}

						var sensor Sensor
						err = json.Unmarshal([]byte(msg.Dado), &sensor)
						if err != nil {
							fmt.Println("Erro ao Ler Dados do sensor")
							continue
						}

						limparTela()

						fmt.Println("\nLendo dados do sensor... Aperte ENTER para sair")

						fmt.Printf("Sensor: %s | Dado: %s | Hora: %s\n",
							sensor.ID,
							sensor.Dado,
							sensor.UltimoVisto.In(loc).Format("15:04:05"),
						)
					}
				}
			}()

			input.ReadString('\n')

			close(stopChan)

			<-done

			err = enviar(conn, idCliente, dado, RemoverInscrito)
			if err != nil {
				errCon = true
				continue
			}

			fmt.Println("Leitura de sensor finalizada")

		case "3":

			err = enviar(conn, "nil", "nil", ListarAtuadores)
			if err != nil {
				errCon = true
				continue
			}

			fmt.Println("\nAguardando mensagem do servidor...")

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

			var lista []Atuador

			err = json.Unmarshal([]byte(msg.Dado), &lista)
			if err != nil {
				fmt.Println("Erro ao Ler Dados da Lista")
				continue
			}

			var atuadID string
			var atuadTipo string

			fmt.Printf("\n%-10s | %-15s | %-15s\n", "Tipo", "ID", "Status")
			fmt.Println("----------------------------------------------------------")

			for _, atuador := range lista {

				idx := strings.Index(atuador.ID, "_")
				if idx != -1 {
					atuadTipo = atuador.ID[:idx]
					atuadID = atuador.ID[idx+1:]
				}

				fmt.Printf("%-10s | %-15s | %-15s\n",
					atuadTipo, atuadID, atuador.Status,
				)
			}

		case "4":

			var dado string

			fmt.Println("\n===== TIPO DO ATUADOR =====")
			fmt.Println("1 - Alarme")
			fmt.Println("2 - VMI")
			fmt.Print("Escolha o tipo do Atuador: ")
			tipoAtuador, _ = input.ReadString('\n')
			tipoAtuador = strings.TrimSpace(tipoAtuador)

			switch tipoAtuador {

			case "1":
				tipoAtuador = "alarme"

			case "2":
				tipoAtuador = "vmi"

			default:
				fmt.Println("\nOpção inválida!")
				continue
			}

			fmt.Println("\nDigite o id do Atuador: ")
			id, _ := input.ReadString('\n')
			id = strings.TrimSpace(id)
			id = tipoAtuador + "_" + id

			switch tipoAtuador {

			case "alarme":
				fmt.Println("\n===== COMANDOS =====")
				fmt.Println("1 - Ligar Alarme")
				fmt.Println("2 - Desligar Alarme")
				fmt.Print("Escolha um comando: ")
				dado, _ = input.ReadString('\n')
				dado = strings.TrimSpace(dado)

				switch dado {
				case "1":
					fmt.Print("Digite a mensagem do alarme: ")
					info, _ := input.ReadString('\n')
					info = strings.TrimSpace(info)

					dado = "LIGAR ALARME: " + info

				case "2":
					dado = "DESLIGAR ALARME"

				default:
					fmt.Println("\nOpção inválida!")
					continue
				}

			case "vmi":
				fmt.Println("\n===== COMANDOS =====")
				fmt.Println("1 - Ligar VMI em modo controlado")
				fmt.Println("2 - Ligar VMI em modo assisto-controlado")
				fmt.Println("3 - Ligar VMI em modo espontâneo")
				fmt.Println("4 - Desligar VMI")
				fmt.Print("Escolha um comando: ")
				dado, _ = input.ReadString('\n')
				dado = strings.TrimSpace(dado)

				switch dado {
				case "1":
					dado = "LIGAR VMI: MODO CONTROLADO"

				case "2":
					dado = "LIGAR VMI: MODO ASSISTO-CONTROLADO"

				case "3":
					dado = "LIGAR VMI: MODO ESPONTÂNEO"

				case "4":
					dado = "DESLIGAR VMI"

				default:
					fmt.Println("\nOpção inválida!")
					continue
				}

			}

			err = enviar(conn, id, dado, AcaoAtuador)
			if err != nil {
				errCon = true
				continue
			}

			fmt.Println("\nEnviando mensagem ao servidor...")

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

			switch msg.Dado {
			case "ACAO ENVIADA COM SUCESSO":
				fmt.Println("Comando Enviado com Sucesso")
			case "ATUADOR NÃO ENCONTRADO":
				fmt.Println("Erro ao enviar comando: Atuador " + id + "não encontrado")
			}

		case "5":
			fmt.Println("Encerrando cliente...")
			return

		default:
			fmt.Println("\nOpção inválida!")

		}

	}

}
