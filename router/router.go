package router

import (
	"net/http"
	"todo-app/handlers"

	"github.com/gorilla/mux"
)

// NewRouter создаёт новый роутер
func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/nextdate", handlers.NextDateHandler).Methods("GET")
	r.HandleFunc("/api/task", handlers.TaskHandler).Methods("POST")
	r.HandleFunc("/api/tasks", handlers.GetTasksHandler).Methods("GET")

	// Маршруты для статических файлов
	webDir := "./web"
	fs := http.FileServer(http.Dir(webDir))
	r.PathPrefix("/").Handler(fs)

	return r
}
