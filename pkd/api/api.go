package api

import (
	"net/http"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getHandler(w, r)
	case http.MethodPut:
		updateHandler(w, r)
	case http.MethodDelete:
		deleteHandler(w, r)
	}
}

func Init() {
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/task", auth(taskHandler))
	http.HandleFunc("/api/tasks", auth(tasksHandler))
	http.HandleFunc("/api/task/done", auth(doneHandler))
	http.HandleFunc("/api/signin", signinHandler)
}
