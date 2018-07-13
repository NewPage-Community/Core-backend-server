package main

import (
	//core
	"log"

	"./core"
)

func main() {
	defer core.CloseDB()
	if ok := core.InitSetting(); !ok {
		log.Println("Failed to load setting")
		return
	}

	if ok := core.ConnectDB(); !ok {
		log.Println("Failed to connect to mysql")
		return
	}

	core.StartTCPServer()
}
