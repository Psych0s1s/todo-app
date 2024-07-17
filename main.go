package main

import (
	"log"
	"net/http"
	"os"
	"todo-app/db"
	"todo-app/router"
)

func main() {
	// Инициализация БД
	db.InitDB()
	defer db.DB.Close()

	// Определение порта
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540" // Порт по умолчанию
	}

	r := router.NewRouter()

	log.Printf("Starting server on :%s\n", port)
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal(err)
	}
}
