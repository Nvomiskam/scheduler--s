// Пакет api реализует и регистрирует API обработчики
package api

import (
	"encoding/json"
	"net/http"
	"time"

	"go1f/pkg/db"
	"go1f/pkg/rules"
)

// updateTaskHandler обрабатывает PUT-запрос на обновление задачи
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": "ошибка десериализации JSON"}, http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		writeJson(w, map[string]string{"error": "не указан идентификатор задачи"}, http.StatusBadRequest)
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

	parsedDateStr := parsedDate.Format(rules.DateFormat)
	nowStr := now.Format(rules.DateFormat)

	if parsedDateStr < nowStr {
		if task.Repeat == "" {
			task.Date = nowStr
		} else {
			task.Date = nextDate
		}
	}

	dbInstance, err := db.NewDB(db.GetDBPath())
	if err != nil {
		writeJson(w, map[string]string{"error": "ошибка подключения к базе данных"}, http.StatusInternalServerError)
		return
	}

	err = dbInstance.UpdateTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": "ошибка обновления задачи"}, http.StatusInternalServerError)
		return
	}

	writeJson(w, map[string]string{}, http.StatusOK)
}
