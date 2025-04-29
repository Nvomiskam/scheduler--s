package api

import (
	"net/http"

	"go1f/pkg/db"
)

// deleteTaskHandler обрабатывает DELETE-запрос на удаление задачи
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "не указан идентификатор задачи"}, http.StatusBadRequest)
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": "не удалось удалить задачу"}, http.StatusInternalServerError)
		return
	}

	writeJson(w, map[string]string{}, http.StatusOK)
}
