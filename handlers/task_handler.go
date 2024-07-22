package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"todo-app/db"
	"todo-app/utils"
)

type Task struct {
	ID      string `db:"id" json:"id"`
	Date    string `db:"date" json:"date"`
	Title   string `db:"title" json:"title"`
	Comment string `db:"comment" json:"comment"`
	Repeat  string `db:"repeat" json:"repeat"`
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handleCreateTask(w, r)
	case http.MethodPut:
		handleUpdateTask(w, r)
	case http.MethodGet:
		handleGetTask(w, r)
	default:
		http.Error(w, `{"error": "Метод не поддерживается"}`, http.StatusMethodNotAllowed)
	}
}

func handleCreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		response := map[string]string{"error": "Ошибка десериализации JSON"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if task.Title == "" {
		response := map[string]string{"error": "Не указан заголовок задачи"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	const layout = "20060102"
	now := time.Now()
	nowStr := now.Format(layout)

	// Если поле Date пустое, устанавливаем сегодняшнюю дату
	if task.Date == "" {
		task.Date = nowStr
	} else {
		// Парсим дату из строки
		parsedDate, err := time.Parse(layout, task.Date)
		if err != nil {
			response := map[string]string{"error": "Дата указана в неверном формате"}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		if task.Date == nowStr {
			// Если дата задачи совпадает с текущей или больше, ничего не меняем
		} else if parsedDate.Before(now) {
			// Если дата задачи раньше текущей
			if task.Repeat == "" {
				task.Date = nowStr
			} else {
				nextDate, err := utils.NextDate(now, task.Date, task.Repeat)
				if err != nil {
					response := map[string]string{"error": err.Error()}
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(response)
					return
				}
				task.Date = nextDate
			}
		}

	}

	id, err := db.AddTask(task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		response := map[string]string{"error": err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := map[string]string{"id": fmt.Sprintf("%d", id)}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		response := map[string]string{"error": "Ошибка десериализации JSON"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if task.ID == "" {
		response := map[string]string{"error": "Не указан идентификатор задачи"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if task.Title == "" {
		response := map[string]string{"error": "Не указан заголовок задачи"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	const layout = "20060102"
	now := time.Now()
	nowStr := now.Format(layout)

	// Если поле Date пустое, устанавливаем сегодняшнюю дату
	if task.Date == "" {
		task.Date = nowStr
	} else {
		// Парсим дату из строки
		parsedDate, err := time.Parse(layout, task.Date)
		if err != nil {
			response := map[string]string{"error": "Дата указана в неверном формате"}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		if task.Date == nowStr {
			// Если дата задачи совпадает с текущей или больше, ничего не меняем
		} else if parsedDate.Before(now) {
			// Если дата задачи раньше текущей
			if task.Repeat == "" {
				task.Date = nowStr
			} else {
				nextDate, err := utils.NextDate(now, task.Date, task.Repeat)
				if err != nil {
					response := map[string]string{"error": err.Error()}
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(response)
					return
				}
				task.Date = nextDate
			}
		}

	}

	err := db.UpdateTask(task.ID, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		response := map[string]string{"error": err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{})
}

func handleGetTask(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		http.Error(w, `{"error":"Не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		http.Error(w, `{"error":"Некорректный идентификатор"}`, http.StatusBadRequest)
		return
	}

	task, err := db.GetTaskByID(id)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(task)
}
