package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"my_project/pkd/db"
)

type respErr struct {
	Error string `json:"error"`
}
type respTask struct {
	ID string `json:"id"`
}

func respErrJSON(w http.ResponseWriter, textErr string, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(respErr{Error: textErr}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func checkDate(task *db.Task) error {
	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format(formatDate)
	}
	t, err := time.Parse(formatDate, task.Date)
	if err != nil {
		return fmt.Errorf("invalid format date: %w", err)
	}
	if task.Repeat != "" {
		_, err = nextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("invalid next date: %w", err)
		}
	}
	if now.After(t) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format(formatDate)
		} else {
			task.Date = now.Format(formatDate)
		}

	}
	return nil
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respErrJSON(w, "Ошибка чтения тела", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	if err = json.Unmarshal(body, &task); err != nil {
		respErrJSON(w, "Ошибка десериализации JSON", http.StatusBadRequest)
		return
	}
	if task.Title == "" {
		respErrJSON(w, "Пустое поле заголовка", http.StatusBadRequest)
		return
	}
	if err = checkDate(&task); err != nil {
		respErrJSON(w, "Некорректное поле даты", http.StatusBadRequest)
		return
	}
	id, err := db.AddTask(&task)
	if err != nil {
		respErrJSON(w, "Ошибка добавление задачи в базу данных", http.StatusBadRequest)
		return
	}
	response := fmt.Sprintf("%v", id)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(respTask{ID: response}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
