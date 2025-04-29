// Пакет api реализует и регистрирует API обработчики
package api

import (
	"net/http"

	"go1f/pkg/db"
)

// getTaskHandler обрабатывает GET-запрос на получение задачи по id
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	writeJson(w, task, http.StatusOK)
}
