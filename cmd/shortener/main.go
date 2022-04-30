package main

import (
	"context"

	"github.com/dimk00z/go-shortener-praktikum/internal/handlers"
	"github.com/dimk00z/go-shortener-praktikum/internal/server"
	"github.com/dimk00z/go-shortener-praktikum/internal/settings"
	"github.com/dimk00z/go-shortener-praktikum/internal/storage"
)

func main() {
	config := settings.LoadConfig()
	rootHandler := handlers.NewRootHandler(
		"http://"+config.Server.Host+":"+config.Server.Port,
		storage.GetStorage)
	server := server.NewServer(":" + config.Server.Port)
	server.MountHandlers(*rootHandler)
	ctx, cancel := context.WithCancel(context.Background())
	server.RunServer(ctx, cancel)
}
