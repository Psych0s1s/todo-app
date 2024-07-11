package main

import (
	"log"
	"net/http"
	"os"
	"todo-app/db"
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

	// Директория с файлами веб-интерфейса
	webDir := "./web"
	fs := http.FileServer(http.Dir(webDir))
	http.Handle("/", fs)

	log.Printf("Starting server on :%s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
