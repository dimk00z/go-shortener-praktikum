package main

import (
	"github.com/dimk00z/go-shortener-praktikum/config"
	"github.com/dimk00z/go-shortener-praktikum/internal/app"
)

func main() {
	// Configuration
	cfg := config.LoadConfig()
	// Run
	app.StartApp(cfg)
}
