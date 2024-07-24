package router

import (
	"net/http"
	"todo-app/auth"
	"todo-app/handlers"

	"github.com/gorilla/mux"
)

// Создаем роутер
func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/signin", auth.SigninHandler).Methods("POST")
	r.HandleFunc("/api/nextdate", handlers.NextDateHandler).Methods("GET")
	r.Handle("/api/task", auth.AuthMiddleware(http.HandlerFunc(handlers.TaskHandler))).Methods("POST", "PUT", "GET", "DELETE")
	r.Handle("/api/task/done", auth.AuthMiddleware(http.HandlerFunc(handlers.HandleCompleteTask))).Methods("POST")
	r.Handle("/api/tasks", auth.AuthMiddleware(http.HandlerFunc(handlers.GetTasksHandler))).Methods("GET")

	// Маршрут для файлов фронтенда
	webDir := "./web"
	fs := http.FileServer(http.Dir(webDir))
	r.PathPrefix("/").Handler(fs)

	return r
}
