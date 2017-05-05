package server

import (
	"fmt"
	"log"
	"net/http"
)

// struct do server

type Server struct {
	Hostname string `json:"Hostname"`
	HTTPPort string `json:"HTTPPort"`
}

// funcao usada para rodar o servidor

func Run(httpHandlers http.Handler, s Server) {
	log.Fatal(http.ListenAndServe(s.HTTPPort, httpHandlers))
}

// funcao para retornar o acesso ao servidor

func httpAddress(s Server) string {
	return s.Hostname + ":" + fmt.Sprintf("%d", s.HTTPPort)
}
