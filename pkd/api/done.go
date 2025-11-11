package api

import (
	"encoding/json"
	"net/http"
	"time"

	"my_project/pkd/db"
)

func doneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		respErrJSON(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		respErrJSON(w, "Идентификатор не указан", http.StatusBadRequest)
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		respErrJSON(w, "Задача не найдена", http.StatusBadRequest)
		return
	}
	if task.Repeat == "" {
		if err = db.DeleteTask(id); err != nil {
			respErrJSON(w, "Задача не удалена", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if err := json.NewEncoder(w).Encode(struct{}{}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	now := time.Now()
	nextDate, err := nextDate(now, task.Date, task.Repeat)
	if err != nil {
		respErrJSON(w, "Не удалось расчитать дату повторения", http.StatusBadRequest)
		return
	}
	err = db.UpdateDate(nextDate, id)
	if err != nil {
		respErrJSON(w, "Не удалось обновить дату", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(struct{}{}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
