package main

import (
	"app/route"
	"app/shared/database"
	"app/shared/server"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"runtime"
)

// struct de configuracao do sistema

type configuration struct {
	Database database.Database `json:"Database"`
	Server   server.Server     `json:"Server"`
}

//funcao para copiar os dados do arquivo config.json para a struct

func ParseJsonFile(configPath string) configuration {
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Println("File error %v", err)
		os.Exit(1)
	}
	var config configuration
	json.Unmarshal(file, &config)
	return config
}

func init() {
	log.SetFlags(log.Lshortfile)
	runtime.GOMAXPROCS(runtime.NumCPU())
}

// funcao principal que conecta ao banco de dados e roda o servidor

func main() {
	config := ParseJsonFile("../config/config.json")

	database.Connect(config.Database)

	server.Run(route.LoadHTTP(), config.Server)
}
