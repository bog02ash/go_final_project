package api

import (
	"encoding/json"
	"net/http"

	"my_project/pkd/db"
)

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		respErrJSON(w, "Идентификатор не указан", http.StatusBadRequest)
		return
	}
	if err := db.DeleteTask(id); err != nil {
		respErrJSON(w, "Задача не удалена", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(struct{}{}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
