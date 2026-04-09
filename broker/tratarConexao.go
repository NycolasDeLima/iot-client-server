package main

import (
	"bufio"
	"cmp"
	"encoding/json"
	"log"
	"net"
	"slices"
	"strconv"
	"strings"
	"time"
)

const (
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
		log.Println("Sensor "+msg.ID+" Conectado: ", clientAddr)
	}

	sensores[msg.ID] = sensor

	listaClientes := inscritos[msg.ID]

	mutex.Unlock()

	for _, conn := range listaClientes {

		data, _ := json.Marshal(sensor)

		err := enviar(conn, msg.ID, string(data), "DADO SENSOR")
		if err != nil {
			log.Println("Erro ao enviar broadcast")
		}

	}

}

func tratarCliente(id string, conn net.Conn, reader *bufio.Reader) {
	defer conn.Close()

	log.Println("Dispositivo Conectado: ", id)

	for {

		buffer, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Aplicação Desconectada: ", id)
			removerCliente(conn, "nil")
			return
		}

		var msg MensagemTCP

		err = json.Unmarshal([]byte(buffer), &msg)
		if err != nil {
			log.Println("Erro JSON cliente ", id, ":", err)
			continue
		}

		// Ações do cliente
		if msg.Acao == AcaoAtuador {

			mutex.RLock()
			atuadorConn, existe := atuadoresConn[msg.ID]

			if !existe {
				enviar(conn, msg.ID, "ATUADOR NÃO ENCONTRADO", AcaoAtuador)
				continue
			}

			err := enviar(atuadorConn, msg.ID, msg.Dado, msg.Acao)
			if err != nil {
				log.Println("Erro ao enviar ao atuador ", id, ":", err)
				continue
			}

			mutex.RUnlock()

			err = enviar(conn, msg.ID, "ACAO ENVIADA COM SUCESSO", msg.Acao)
			if err != nil {
				log.Println("Erro ao enviar cliente ", id, ":", err)
				continue
			}

			log.Println("Dispositivo ", id, ": ", "Ação ", msg.Acao, " ", msg.ID)

		} else if msg.Acao == ListarSensores {

			var lista []Sensor

			mutex.RLock()
			for _, sensor := range sensores {
				lista = append(lista, sensor)
			}

			//Ordenação dos sensores é feita primeiro pelo prefixo (bpm, spo2) e depois pelo número (1, 2, 3)

			slices.SortFunc(lista, func(a, b Sensor) int {

				prefA, numAStr, _ := strings.Cut(a.ID, "_")
				prefB, numBStr, _ := strings.Cut(b.ID, "_")

				if prefA != prefB {
					return cmp.Compare(prefA, prefB)
				}

				numA, _ := strconv.Atoi(numAStr)
				numB, _ := strconv.Atoi(numBStr)

				return cmp.Compare(numA, numB)
			})

			data, _ := json.Marshal(lista)

			mutex.RUnlock()

			err := enviar(conn, "nil", string(data), msg.Acao)
			if err != nil {
				log.Println("Erro ao enviar ao Cliente ", id, ":", err)
				continue
			}

			log.Println("Dispositivo ", id, ": ", "Lista de Sensores")

		} else if msg.Acao == ListarAtuadores {

			var lista []Atuador

			mutex.RLock()

			for _, atuador := range atuadores {
				lista = append(lista, atuador)
			}

			slices.SortFunc(lista, func(a, b Atuador) int {

				prefA, numAStr, _ := strings.Cut(a.ID, "_")
				prefB, numBStr, _ := strings.Cut(b.ID, "_")

				if prefA != prefB {
					return cmp.Compare(prefA, prefB)
				}

				numA, _ := strconv.Atoi(numAStr)
				numB, _ := strconv.Atoi(numBStr)

				return cmp.Compare(numA, numB)
			})

			data, _ := json.Marshal(lista)

			mutex.RUnlock()

			err := enviar(conn, "nil", string(data), msg.Acao)

			if err != nil {
				log.Println("Erro ao enviar ao Cliente ", id, ":", err)
				continue
			}

			log.Println("Dispositivo ", id, ": ", "Lista Atuadores")

		} else if msg.Acao == VerDadoSensor {

			sensorID := msg.Dado

			mutex.RLock()
			_, existe := sensores[sensorID]
			mutex.RUnlock()

			if !existe {

				err := enviar(conn, "nil", "SENSOR NÃO ENCONTRADO", msg.Acao)
				if err != nil {
					log.Println("Erro ao enviar ao Cliente ", id, ":", err)
					continue
				}
				continue
			}

			mutex.Lock()
			inscritos[sensorID] = append(inscritos[sensorID], conn)
			mutex.Unlock()

			log.Println("Cliente ", id, "inscrito no sensor:", sensorID)

		} else if msg.Acao == RemoverInscrito {

			sensorID := msg.Dado

			removerCliente(conn, msg.Dado)

			log.Println("Cliente", id, "removido do sensor:", sensorID)
		}
	}
}

func tratarAtuador(id string, conn net.Conn, reader *bufio.Reader) {

	defer func() {
		conn.Close()
		removerAtuador(id)
	}()

	log.Println("Dispositivo Conectado: ", id)

	for {
		buffer, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Dispositivo Desconectado: ", id)
			return
		}

		var msg MensagemTCP

		err = json.Unmarshal([]byte(buffer), &msg)
		if err != nil {
			log.Println("Erro ao Ler Dados do Json")
			continue
		}

		log.Println("Dispositivo ", msg.ID, " Atualizado: ", msg.Dado)

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
		log.Println("JSON inválido:", err)
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
