package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"todo-app/db"
)

// Обработчик для получения списка задач
func GetTasksHandler(w http.ResponseWriter, r *http.Request) {

	// В задании этого нет, но если фронтенд будет поддерживать пагинацию, то это пригодится

	limitParam := r.URL.Query().Get("limit")
	pageParam := r.URL.Query().Get("page")
	searchParam := r.URL.Query().Get("search")

	limit := 50 // Значение по умолчанию
	page := 1   // Значение по умолчанию

	if limitParam != "" {
		l, err := strconv.Atoi(limitParam)
		if err == nil && l >= 10 && l <= 50 {
			limit = l
		}
	}

	if pageParam != "" {
		p, err := strconv.Atoi(pageParam)
		if err == nil && p >= 1 {
			page = p
		}
	}

	offset := (page - 1) * limit

	var tasks []db.Task
	var err error

	if searchParam != "" {
		if isDate(searchParam) {
			tasks, err = db.GetTasksByDate(convertToDate(searchParam), limit, offset)
		} else {
			tasks, err = db.SearchTasks(searchParam, limit, offset)
		}
	} else {
		tasks, err = db.GetTasks(limit, offset)
	}

	if err != nil {
		http.Error(w, `{"error":"Ошибка при получении списка задач"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"tasks": tasks,
	}

	if tasks == nil {
		response["tasks"] = []db.Task{}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(response)
}

// Обработчик для получения задачи по идентификатору
func GetTaskByIDHandler(w http.ResponseWriter, r *http.Request) {
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

// Проверка на соответствие строки формату даты
func isDate(str string) bool {
	_, err := time.Parse("02.01.2006", str)
	return err == nil
}

// Преобразование строки формата 02.01.2006 в 20060102
func convertToDate(str string) string {
	date, _ := time.Parse("02.01.2006", str)
	return date.Format("20060102")
}
