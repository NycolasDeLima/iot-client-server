package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
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
				log.Println("Erro ao enviar ao atuador ", id, ":", err)
				continue
			}

			mutex.RUnlock()

			err = enviar(conn, msg.ID, "ACAO ENVIADA COM SUCESSO", msg.Acao)
			if err != nil {
				log.Println("Erro ao enviar cliente ", id, ":", err)
				continue
			}

		} else if msg.Acao == ListarSensores {

			var lista []string

			mutex.RLock()
			for idSensor := range sensores {
				lista = append(lista, idSensor)
			}

			log.Println(lista)

			data, _ := json.Marshal(sensores)

			mutex.RUnlock()

			err := enviar(conn, "nil", string(data), msg.Acao)
			if err != nil {
				log.Println("Erro ao enviar ao Cliente ", id, ":", err)
				continue
			}

		} else if msg.Acao == ListarAtuadores {

			var lista []string

			mutex.RLock()

			for id := range atuadores {
				lista = append(lista, id)
			}

			log.Println(lista)

			data, _ := json.Marshal(atuadores)

			mutex.RUnlock()

			err := enviar(conn, "nil", string(data), msg.Acao)

			if err != nil {
				log.Println("Erro ao enviar ao Cliente ", id, ":", err)
				continue
			} else {
				log.Println("Lista de atuadores enviada ao cliente ", id)
			}

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

			log.Println("Cliente inscrito no sensor:", sensorID)

		} else if msg.Acao == RemoverInscrito {

			sensorID := msg.Dado

			removerCliente(conn, msg.Dado)

			log.Println("Cliente removido do sensor:", sensorID)
		}
	}
}

func tratarAtuador(id string, conn net.Conn, reader *bufio.Reader) {

	defer func() {
		conn.Close()
		removerAtuador(id)
	}()

	log.Println("Atuador conectado: ", id)

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

		log.Println("Recebido:", msg.Acao)

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
