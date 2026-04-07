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

		fmt.Println("ALARME CRÍTICO ATIVADO")

		idx = strings.Index(msg, ":")
		if idx != -1 {
			motivo = msg[idx+2:]
		}
		fmt.Println("Motivo: ", motivo)

		// simula sirene
		for i := 0; i < 3; i++ {
			fmt.Println("BEEP BEEP BEEP")
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

			fmt.Println("VMI LIGADA EM MODO CONTROLADO")
			ativo = true
			estado = modo

		case "MODO ASSISTO-CONTROLADO":

			fmt.Println("VMI LIGADA EM MODO ASSISTO-CONTROLADO")
			ativo = true
			estado = modo

		case "MODO ESPONTÂNEO":

			fmt.Println("VMI LIGADA EM MODO ESPONTÂNEO")
			ativo = false
			estado = modo

		default:
			fmt.Println("COMANDO NÃO IDENTIFICADO")
		}

	} else if strings.Contains(msg, "DESLIGAR VMI") {
		fmt.Println("VMI DESLIGADA")
		ativo = false
		estado = "DESLIGADA"
	}
}
