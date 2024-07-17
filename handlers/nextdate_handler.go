package handlers

import (
	"net/http"
	"time"

	"todo-app/utils"
)

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.URL.Query().Get("now")
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	const layout = "20060102"
	now, err := time.Parse(layout, nowStr)
	if err != nil {
		http.Error(w, "время не может быть преобразовано в корректную дату", http.StatusBadRequest)
		return
	}

	nextDate, err := utils.NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(nextDate))
}
