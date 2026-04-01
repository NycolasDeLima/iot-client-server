package main

import (
	"math/rand"
)

func mudarEstado(tipoSensor string) string {

	var estados []string

	if tipoSensor == "bpm" {
		estados = []string{"repouso", "atividade", "taquicardia", "bradicardia"}
	}
	return estados[rand.Intn(len(estados))]
}

func ajustarBPM(atual int, estado string) int {
	var alvo int

	switch estado {
	case "repouso":
		alvo = 70
	case "atividade":
		alvo = 110
	case "taquicardia":
		alvo = 140
	case "bradicardia":
		alvo = 50
	}

	// aproxima suavemente do alvo
	if atual < alvo {
		atual += rand.Intn(3)
	} else if atual > alvo {
		atual -= rand.Intn(3)
	}

	// pequeno ruído
	atual += rand.Intn(3) - 1

	return limitar(atual, 40, 180)
}

func limitar(valor, min, max int) int {
	if valor < min {
		return min
	}
	if valor > max {
		return max
	}
	return valor
}
