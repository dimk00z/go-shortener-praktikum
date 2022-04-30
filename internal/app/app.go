package app

import (
	"context"

	"github.com/dimk00z/go-shortener-praktikum/internal/server"
	"github.com/dimk00z/go-shortener-praktikum/internal/settings"
	"github.com/dimk00z/go-shortener-praktikum/internal/storage"
)

func StartApp() {
	config := settings.LoadConfig()
	// host := "http://" + config.Server.Host + ":" + config.Server.Port
	host := config.Server.Host
	// server := server.NewServer(":" + config.Server.Port)
	server := server.NewServer(config.Server.Port)
	server.MountHandlers(host, storage.GetStorage)
	ctx, cancel := context.WithCancel(context.Background())
	server.RunServer(ctx, cancel)
}
