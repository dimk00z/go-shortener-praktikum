package app

import (
	"context"
	"log"

	"github.com/dimk00z/go-shortener-praktikum/internal/server"
	"github.com/dimk00z/go-shortener-praktikum/internal/settings"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/storagedi"
)

func StartApp() {
	config := settings.LoadConfig()
	log.Printf("%+v\n", config)

	host := config.Server.Host
	srv := server.NewServer(config.Server.Port)
	storage := storagedi.GetStorage(config.Storage)

	srv.MountHandlers(host, storage)
	ctx, cancel := context.WithCancel(context.Background())
	srv.RunServer(ctx, cancel, storage)
}
