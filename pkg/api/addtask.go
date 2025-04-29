// Пакет api реализует и регистрирует API обработчики
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go1f/pkg/db"
	"go1f/pkg/rules"
)

// addTaskHandler обрабатывает POST-запрос на добавление задачи в БД
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": "ошибка десериализации JSON"}, http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeJson(w, map[string]string{"error": "не указан заголовок задачи"}, http.StatusBadRequest)
		return
	}

	now := time.Now()
	var parsedDate time.Time

	if task.Date == "" {
		task.Date = now.Format(rules.DateFormat)
	}

	parsedDate, err = time.Parse(rules.DateFormat, task.Date)
	if err != nil {
		writeJson(w, map[string]string{"error": "дата представлена в неправильном формате"}, http.StatusBadRequest)
		return
	}

	var nextDate string
	if task.Repeat != "" {
		nextDate, err = rules.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJson(w, map[string]string{"error": "неверное правило повторения"}, http.StatusBadRequest)
			return
		}
	}

	// Проверяем - не в прошлом ли дата
	parsedDateStr := parsedDate.Format(rules.DateFormat)
	nowStr := now.Format(rules.DateFormat)

	if parsedDateStr < nowStr {
		if task.Repeat == "" {
			task.Date = nowStr
		} else {
			task.Date = nextDate
		}
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": "ошибка добавления задачи"}, http.StatusInternalServerError)
		return
	}

	writeJson(w, map[string]string{"id": fmt.Sprintf("%d", id)}, http.StatusOK)
}
