package server

import (
	"net/http"
	"fmt"
	"log"
)

type Server struct {
	Hostname string `json:"Hostname"`
	HTTPPort string `json:"HTTPPort"`
}

func Run(httpHandlers http.Handler, s Server) {
	log.Fatal(http.ListenAndServe(s.HTTPPort, httpHandlers))
}

func httpAddress(s Server) string {
	return s.Hostname + ":" + fmt.Sprintf("%d", s.HTTPPort)
}
