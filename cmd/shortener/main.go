package main

import (
	"github.com/dimk00z/go-shortener-praktikum/config"
	"github.com/dimk00z/go-shortener-praktikum/internal/app"
	"github.com/dimk00z/go-shortener-praktikum/internal/util"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
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
	// Print build info
	util.PrintBulidInfo(buildVersion, buildDate, buildCommit)
	// Configuration
	cfg := config.LoadConfig()
	// Run
	app.StartApp(cfg)
}
