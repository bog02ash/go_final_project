package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const formatDate = "20060102"

func nextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("repetition is empty")
	}
	date, err := time.Parse(formatDate, dstart)
	if err != nil {
		return "", fmt.Errorf("invalid date format %w", err)
	}
	if repeat == "y" {
		for {
			date = date.AddDate(1, 0, 0)
			if date.After(now) {
				return date.Format(formatDate), nil
			}
		}
	}
	repeatSplit := strings.Split(repeat, " ")
	if len(repeatSplit) < 2 && repeatSplit[0] == "d" {
		return "", errors.New("interval missing for daily repeat")
	}
	if repeatSplit[0] != "d" && repeatSplit[0] != "y" {
		return "", errors.New("unsupported formatting")
	}
	interval, err := strconv.Atoi(repeatSplit[1])
	if err != nil {
		return "", fmt.Errorf("invalid interval %w", err)
	}
	if interval < 0 || interval > 400 {
		return "", errors.New("invalid interval")
	}
	if repeatSplit[0] == "d" {
		for {
			date = date.AddDate(0, 0, interval)
			if date.After(now) {
				break
			}
		}
	}
	return date.Format(formatDate), nil

}
func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
		return
	}
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")
	nowStr := r.FormValue("now")
	var nowTime time.Time
	var err error
	if date == "" {
		http.Error(w, "Неуказана дата", http.StatusBadRequest)
		return
	}
	if repeat == "" {
		http.Error(w, "Не указано повторение", http.StatusBadRequest)
		return
	}
	if nowStr == "" {
		nowTime = time.Now()
	} else {
		nowTime, err = time.Parse(formatDate, nowStr)
		if err != nil {
			http.Error(w, "Неверные входные параметры", http.StatusBadRequest)
			return
		}
	}
	writer, err := nextDate(nowTime, date, repeat)
	if err != nil {
		http.Error(w, "Неверные входные параметры", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(writer))
}
