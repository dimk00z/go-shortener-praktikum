package main

import (
	"net/http"

	"github.com/dimk00z/go-shortener-praktikum/internal/handlers"
)

func main() {
	// маршрутизация запросов обработчику
	http.HandleFunc("/", handlers.HelloWorld)
	// запуск сервера с адресом localhost, порт 8080
	http.ListenAndServe(":8080", nil)
}
