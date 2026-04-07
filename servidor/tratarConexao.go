package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

const (
	//ENVIAR
	ListarSensores  = "LISTAR SENSORES"
	ListarAtuadores = "LISTAR ATUADORES"
	AcaoAtuador     = "ACAO ATUADOR"
	VerDadoSensor   = "VER DADO SENSOR"
	RemoverInscrito = "REMOVER INSCRITO"
)

func tratarSensor(msg MensagemUDP, clientAddr *net.UDPAddr) {

	sensor := Sensor{
		Tipo:        msg.Tipo,
		ID:          msg.ID,
		Dado:        msg.Dado,
		UltimoVisto: time.Now(),
	}

	mutex.Lock()

	_, ok := sensores[msg.ID]

	if !ok {
		fmt.Println("Sensor "+msg.ID+" Conectado: ", clientAddr)
	}

	sensores[msg.ID] = sensor

	listaClientes := inscritos[msg.ID]

	mutex.Unlock()

	for _, conn := range listaClientes {

		data, _ := json.Marshal(sensor)

		err := enviar(conn, msg.ID, string(data), "DADO SENSOR")
		if err != nil {
			fmt.Println("Erro ao enviar broadcast")
		}

	}

}

func tratarCliente(id string, conn net.Conn, reader *bufio.Reader) {
	defer conn.Close()

	for {

		buffer, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Cliente", id, "desconectado")
			removerCliente(conn, "nil")
			return
		}

		var msg MensagemTCP

		err = json.Unmarshal([]byte(buffer), &msg)
		if err != nil {
			fmt.Println("Erro JSON cliente ", id, ":", err)
			continue
		}

		if msg.Acao == AcaoAtuador {

			mutex.RLock()
			atuadorConn, existe := atuadoresConn[msg.ID]
			//mutex.RUnlock()

			if !existe {
				enviar(conn, msg.ID, "ATUADOR NÃO ENCONTRADO", AcaoAtuador)
				continue
			}

			err := enviar(atuadorConn, msg.ID, msg.Dado, msg.Acao)
			if err != nil {
				fmt.Println("Erro ao enviar ao atuador ", id, ":", err)
				continue
			}

			mutex.RUnlock()

			err = enviar(conn, msg.ID, "ACAO ENVIADA COM SUCESSO", msg.Acao)
			if err != nil {
				fmt.Println("Erro ao enviar cliente ", id, ":", err)
				continue
			}

		} else if msg.Acao == ListarSensores {

			var lista []string

			mutex.RLock()
			for id := range sensores {
				lista = append(lista, id)
			}

			fmt.Println(lista)

			data, _ := json.Marshal(sensores)

			mutex.RUnlock()

			err := enviar(conn, "nil", string(data), msg.Acao)
			if err != nil {
				fmt.Println("Erro ao enviar ao Cliente ", id, ":", err)
				continue
			}

		} else if msg.Acao == ListarAtuadores {

			var lista []string

			mutex.RLock()

			for id := range atuadores {
				lista = append(lista, id)
			}

			fmt.Println(lista)

			data, _ := json.Marshal(atuadores)

			mutex.RUnlock()

			err := enviar(conn, "nil", string(data), msg.Acao)
			if err != nil {
				fmt.Println("Erro ao enviar ao Cliente ", id, ":", err)
				continue
			}

		} else if msg.Acao == VerDadoSensor {

			sensorID := msg.Dado

			mutex.RLock()
			_, existe := sensores[sensorID]
			mutex.RUnlock()

			if !existe {

				err := enviar(conn, "nil", "SENSOR NÃO ENCONTRADO", msg.Acao)
				if err != nil {
					fmt.Println("Erro ao enviar ao Cliente ", id, ":", err)
					continue
				}
				continue
			}

			mutex.Lock()
			inscritos[sensorID] = append(inscritos[sensorID], conn)
			mutex.Unlock()

			fmt.Println("Cliente inscrito no sensor:", sensorID)

		} else if msg.Acao == RemoverInscrito {

			sensorID := msg.Dado

			removerCliente(conn, msg.Dado)

			fmt.Println("Cliente removido do sensor:", sensorID)
		}
	}
}

func tratarAtuador(id string, conn net.Conn, reader *bufio.Reader) {

	defer func() {
		conn.Close()
		removerAtuador(id)
	}()

	fmt.Println("Atuador conectado: ", id)

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

		mutex.Lock()

		atuadores[id] = Atuador{
			Tipo:   "ATUADOR",
			ID:     id,
			Status: msg.Dado,
		}

		mutex.Unlock()
	}
}

// =============== tratar conexoes ================

func tratarConexaoTcp(conn net.Conn) {

	var msg MensagemTCP


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
		atuadoresConn[id] = conn
		atuadores[id] = Atuador{
			Tipo:   "ATUADOR",
			ID:     id,
			Status: msg.Dado,
		}
		mutex.Unlock()

		enviar(conn, "nil", "ATUADOR CONECTADO", "HANDSHAKE")

		go tratarAtuador(id, conn, reader)
	}

	if msg.Tipo == "CLIENTE" {
		id := msg.ID

		mutex.Lock()
		clientes[id] = conn
		mutex.Unlock()

		enviar(conn, "nil", "CLIENTE CONECTADO", "HANDSHAKE")

		go tratarCliente(id, conn, reader)
	}

}
