// Пакет api реализует и регистрирует API обработчики
package api

import (
	"net/http"
	"time"

	"go1f/pkg/rules"
)

// nextDateHandler обрабатывает запрос на расчёт следующей даты
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")
	nowStr := r.FormValue("now")

	var now time.Time
	var err error

	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(rules.DateFormat, nowStr)
		if err != nil {
			writeJson(w, map[string]string{"error": "некорректный формат даты"}, http.StatusBadRequest)
			return
		}
	}

	result, err := rules.NextDate(now, dateStr, repeat)
	if err != nil {
		writeJson(w, map[string]string{"error": "ошибка вычисления даты"}, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(result))
}
