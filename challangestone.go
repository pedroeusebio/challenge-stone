package main

import (
	"encoding/json"
	"log"
	"os"
	"io/ioutil"
	"runtime"
	"app/shared/database"
	"app/shared/server"
	"app/route"
)

type configuration struct {
	Database database.Database `json:"Database"`
	Server   server.Server `json:"Server"`
}

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

func main() {
	config := ParseJsonFile("config/config.json")

	database.Connect(config.Database)

	server.Run(route.LoadHTTP(), config.Server)
}
