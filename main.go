package main

import (
	"log"
	"os"
	"scheduler/internal/db"
	"scheduler/internal/server"
)

func main() {
	dbfile := os.Getenv("TODO_DBFILE")
	if dbfile == "" {
		dbfile = "scheduler.db"
	}

	if err := db.Init(dbfile); err != nil {
		log.Fatal("Ошибка инициализации БД:", err)
	}
	defer db.Close()

	log.Println("Запуск сервера планировщика...")
	if err := server.StartServer(); err != nil {
		log.Fatal(err)
	}
}
