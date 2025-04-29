// Пакет api реализует и регистрирует API обработчики
package api

import (
	"net/http"
	"time"

	"go1f/pkg/db"
	"go1f/pkg/rules"
)

// doneTaskHandler обрабатывает POST-запрос на отметку задачи выполненной
func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "не указан идентификатор задачи"}, http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": "задача не найдена"}, http.StatusNotFound)
		return
	}

	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			writeJson(w, map[string]string{"error": "не удалось удалить задачу"}, http.StatusInternalServerError)
			return
		}
	} else {
		now := time.Now()
		nextDate, err := rules.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJson(w, map[string]string{"error": "ошибка расчета следующей даты"}, http.StatusInternalServerError)
			return
		}

		err = db.UpdateDate(nextDate, id)
		if err != nil {
			writeJson(w, map[string]string{"error": "не удалось обновить дату задачи"}, http.StatusInternalServerError)
			return
		}
	}

	writeJson(w, map[string]string{}, http.StatusOK)
}
