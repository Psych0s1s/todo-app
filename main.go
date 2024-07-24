package main

import (
	"log"
	"net/http"
	"os"
	"todo-app/db"
	"todo-app/router"

	"github.com/joho/godotenv"
)

func main() {
	// Загружаем переменные окружения из файла .env
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file, using default values")
	}

	// Инициализация БД
	db.InitDB()
	defer db.DB.Close()

	// Определение порта
	port := os.Getenv("PORT")
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
