package ordenate

import (
	"encoding/json"
	"log"
	"strings"
)

// struct de ordenacao

type Ordenate struct {
	Column string `json:"Column"`
	Order  string `json:"Order"`
}

// funcao usada para capturar o json da url e retornar um array de objetos que
// serao usados na ordenacao
func Order(order string) ([]Ordenate, error) {
	ordenate := []Ordenate{}
	dec := json.NewDecoder(strings.NewReader(order))
	_, err := dec.Token()
	if err != nil {
		log.Fatal(err)
		return ordenate, err
	} else {
		for dec.More() {
			var o Ordenate
			err := dec.Decode(&o)
			if err != nil {
				log.Fatal(err)
			}
			ordenate = append(ordenate, o)
		}
	}
	return ordenate, nil
}
