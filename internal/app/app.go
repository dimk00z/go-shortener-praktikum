package app

import (
	"context"
	"log"

	"github.com/dimk00z/go-shortener-praktikum/internal/server"
	"github.com/dimk00z/go-shortener-praktikum/internal/settings"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/storage_di"
)

func StartApp() {
	config := settings.LoadConfig()
	log.Printf("%+v\n", config)

	host := config.Server.Host
	server := server.NewServer(config.Server.Port)
	storage := storage_di.GetStorage(config.Storage)
	defer storage.Close()

	server.MountHandlers(host, storage)
	ctx, cancel := context.WithCancel(context.Background())
	server.RunServer(ctx, cancel)
}
