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

// Константа для лимита количества задач
const taskLimit = 50

// tasksHandler обрабатывает запрос для возвращения списка задач
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	var tasks []*db.Task
	var err error

	dbInstance, err := db.NewDB(db.GetDBPath())
	if err != nil {
		writeJson(w, map[string]string{"error": "ошибка подключения к базе данных"}, http.StatusInternalServerError)
		return
	}
	if search != "" {
		tasks, err = dbInstance.GetFilteredTasks(search, taskLimit)
	} else {
		tasks, err = dbInstance.Tasks(taskLimit)
	}

	if err != nil {
		writeJson(w, map[string]string{"error": "ошибка при возвращении задачи"}, http.StatusInternalServerError)
		return
	}

	writeJson(w, TasksResp{Tasks: tasks}, http.StatusOK)
}
