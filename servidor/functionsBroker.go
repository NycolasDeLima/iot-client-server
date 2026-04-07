package main

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

func removerAtuador(id string) {
	mutex.Lock()
	defer mutex.Unlock()

	delete(atuadoresConn, id)
	delete(atuadores, id)

	fmt.Println("Atuador desconectado: ", id)
}

func removerSensor() {
	for {
		time.Sleep(5 * time.Second)

		mutex.Lock()
		for id, sensor := range sensores {
			if time.Since(sensor.UltimoVisto) > 5*time.Second {
				fmt.Println("Sensor desconectado: ", id)
				delete(sensores, id)

				listaClientes := inscritos[sensor.ID]
				for _, conn := range listaClientes {
					err := enviar(conn, "nil", "SENSOR DESCONECTADO", "nil")
					if err != nil {
						fmt.Println("Erro ao enviar ao Cliente ", id, ":", err)
						continue
					}
				}

				delete(inscritos, sensor.ID)
			}
		}

		mutex.Unlock()
	}
}

func enviar(conn net.Conn, id string, dado string, acao string) error {

	msg := MensagemTCP{
		Tipo: "SERVIDOR",
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

func removerCliente(conn net.Conn, sensorID string) {
	mutex.Lock()
	defer mutex.Unlock()

	if sensorID == "nil" {

		for ID, lista := range inscritos {

			novaLista := []net.Conn{}

			for _, c := range lista {
				if c != conn {
					novaLista = append(novaLista, c)
				}
			}

			// atualiza a lista
			if len(novaLista) == 0 {
				delete(inscritos, ID)
			} else {
				inscritos[ID] = novaLista
			}
		}
	} else {
		lista := inscritos[sensorID]

		novaLista := []net.Conn{}

		for _, c := range lista {
			if c != conn {
				novaLista = append(novaLista, c)
			}
		}

		inscritos[sensorID] = novaLista
	}
}
