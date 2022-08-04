package main

import (
	"fmt"

	"github.com/dimk00z/go-shortener-praktikum/config"
	"github.com/dimk00z/go-shortener-praktikum/internal/app"
)

func main() {
	// Configuration
	cfg := config.LoadConfig()

	fmt.Println(cfg)
	// Run
	app.StartApp(cfg)
}
