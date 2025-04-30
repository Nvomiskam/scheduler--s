// Пакет server предоставляет HTTP-сервер для обслуживания статических файлов фронтенда из директории "web".
// Сервер использует порт из переменной окружения TODO_PORT или 7540 по умолчанию.
package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

const defaultPort = 7540

// getPort определяет порт из переменной окружения TODO_PORT или возвращает defaultPort
func getPort() int {

	portStr := os.Getenv("TODO_PORT")

	if portStr == "" {
		return defaultPort

	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return defaultPort
	}

	return port

}

// Run запускает HTTP-сервер, используя файлы из папки "web"
func Run() error {
	port := getPort()
	http.Handle("/", http.FileServer(http.Dir("web")))
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
