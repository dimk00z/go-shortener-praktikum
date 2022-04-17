package main

import (
	"context"

	"github.com/dimk00z/go-shortener-praktikum/internal/handlers"
	"github.com/dimk00z/go-shortener-praktikum/internal/server"
	"github.com/dimk00z/go-shortener-praktikum/internal/settings"
)

func main() {
	config := settings.LoadConfig()
	shortenerPort := config.Server.Port
	host := "http://localhost:" + shortenerPort
	rootHandler := handlers.NewRootHandler(host)
	server := server.NewServer(":" + shortenerPort)
	server.MountHandlers(*rootHandler)
	ctx, cancel := context.WithCancel(context.Background())
	server.RunServer(ctx, cancel)
}
