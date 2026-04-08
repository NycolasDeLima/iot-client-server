package main

import (
	"fmt"
	"strings"
)

func tratarAlarme(msg string) {

	var motivo string
	var idx int

	if strings.Contains(msg, "LIGAR ALARME: ") {

		ativo = true

		idx = strings.Index(msg, ":")
		if idx != -1 {
			motivo = msg[idx+2:]
		}

		estado = motivo

	} else if strings.Contains(msg, "DESLIGAR ALARME") {
		ativo = false
		estado = "DESLIGADO"

		fmt.Println("ALARME DESLIGADO")

	} else {
		fmt.Println("COMANDO NÃO IDENTIFICADO")
	}
}

func tratarVMI(msg string) {

	var idx int
	var modo string

	if strings.Contains(msg, "LIGAR VMI: ") {

		idx = strings.Index(msg, ":")
		if idx != -1 {
			modo = msg[idx+2:]
		}

		switch modo {
		case "MODO CONTROLADO":

			ativo = true
			estado = "VMI LIGADA EM MODO CONTROLADO"

		case "MODO ASSISTO-CONTROLADO":

			ativo = true
			estado = "VMI LIGADA EM MODO ASSISTO-CONTROLADO"

		case "MODO ESPONTÂNEO":

			ativo = true
			estado = "VMI LIGADA EM MODO ESPONTÂNEO"

		default:
			fmt.Println("COMANDO NÃO IDENTIFICADO")
		}

	} else if strings.Contains(msg, "DESLIGAR VMI") {
		fmt.Println("VMI DESLIGADA")
		ativo = false
		estado = "DESLIGADA"
	}
}

func limparTela() {
	fmt.Print("\033[H\033[2J")
}

func exibirAtuador(tipo, id, estado string, ativo bool, conectado bool) {
	limparTela()

	fmt.Println("====================================")
	fmt.Println("        PAINEL DO ATUADOR")
	fmt.Println("====================================")

	fmt.Printf("Tipo: %s\n", tipo)
	fmt.Printf("ID:   %s\n", id)

	if conectado {
		fmt.Printf("Status: DESCONECTADO\n")
	} else {
		fmt.Printf("Status: CONECTADO\n")
	}

	fmt.Println("------------------------------------")

	if ativo {
		fmt.Printf("Ativo: SIM\n")
	} else {
		fmt.Printf("Ativo: NÃO\n")
	}

	fmt.Println("------------------------------------")
	fmt.Printf("Mensagem: %s \n", estado)

	fmt.Println("====================================")
}
