package api

import (
	"encoding/json"
	"net/http"
	"time"

	"my_project/pkd/db"
)

type tasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

const limit = 50

func writeJSON(w http.ResponseWriter, tasks []*db.Task) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(tasksResp{Tasks: tasks}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	var tasks []*db.Task
	var err error
	search := r.FormValue("search")
	if search != "" {
		if dateSearch, err := time.Parse("02.01.2006", search); err == nil {
			date := dateSearch.Format(formatDate)
			tasks, err = db.SearchData(date, limit)
			if err != nil {
				respErrJSON(w, "Ошибка получения задач по дате", http.StatusBadRequest)
				return
			}
		} else {
			tasks, err = db.SearchWord(search, limit)
			if err != nil {
				respErrJSON(w, "Ошибка получения задач по слову", http.StatusBadRequest)
				return
			}
		}
	} else {
		tasks, err = db.Tasks(limit)
		if err != nil {
			respErrJSON(w, "Ошибка при извлечении задач из базы данных", http.StatusBadRequest)
			return
		}
	}
	writeJSON(w, tasks)
}
