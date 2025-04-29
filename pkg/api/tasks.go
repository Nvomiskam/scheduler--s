// Пакет api реализует и регистрирует API обработчики
package api

import (
	"net/http"

	"go1f/pkg/db"
)

// TasksResp — структура для JSON-ответа со списком задач
type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// tasksHandler обрабатывает запрос для возвращения списка задач
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	var tasks []*db.Task
	var err error

	if search != "" {
		tasks, err = db.GetFilteredTasks(search, 50)
	} else {
		tasks, err = db.Tasks(50)
	}

	if err != nil {
		writeJson(w, map[string]string{"error": "ошибка при возвращении задачи"}, http.StatusInternalServerError)
		return
	}

	writeJson(w, TasksResp{Tasks: tasks}, http.StatusOK)
}
