package main

import (
	"github.com/dimk00z/go-shortener-praktikum/config"
	"github.com/dimk00z/go-shortener-praktikum/internal/app"
)

// @title Shortener
// @version 1.0.0
// @description Go Shortener project
// @contact.name Kuznetsov Dmitriy
// @contact.url https://github.com/dimk00z
// @contact.email dimk0z@yandex.ru
// @host localhost:8080
// @BasePath /
func main() {
	// Configuration
	cfg := config.LoadConfig()
	// Run
	app.StartApp(cfg)
}
