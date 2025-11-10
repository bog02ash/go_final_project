package api

import (
	"encoding/json"
	"io"
	"my_project/pkd/db"
	"net/http"
)

func updateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respErrJSON(w, "Ошибка чтения тела запроса", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	var task db.Task
	if err = json.Unmarshal(body, &task); err != nil {
		respErrJSON(w, "Ошибка десериализации JSON", http.StatusBadRequest)
		return
	}
	if task.Title == "" {
		respErrJSON(w, "Пустой заголовок", http.StatusBadRequest)
		return
	}
	if err = checkDate(&task); err != nil {
		respErrJSON(w, "Некорректная дата", http.StatusBadRequest)
		return
	}
	err = db.UpdateTask(&task)
	if err != nil {
		respErrJSON(w, "Ошибка редактирования задачи", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(struct{}{}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
