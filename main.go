package main

import (
	"go1f/pkg/api"
	"go1f/pkg/db"
	"go1f/pkg/server"
)

func main() {
	// Инициализация подключения к базе данных
	scheduler, err := db.NewDB(db.GetDBPath())
	if err != nil {
		panic(err)
	}
	// Регистрация API обработчиков
	api.Init(scheduler)

	// Запуск HTTP-сервера
	if err := server.Run(); err != nil {
		panic(err)
	}
}
