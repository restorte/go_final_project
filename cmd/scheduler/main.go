package main

import (
	"log"
	"scheduler/internal/server"
)

func main() {
	log.Println("Запуск сервера планировщика...")
	if err := server.StartServer(); err != nil {
		log.Fatal(err)
	}
}
