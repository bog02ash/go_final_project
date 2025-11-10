package api

import (
	"encoding/json"
	"my_project/pkd/db"
	"net/http"
)

func getHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		respErrJSON(w, "Неуказан иидентификатор", http.StatusBadRequest)
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		respErrJSON(w, "Задача не найдена", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}
