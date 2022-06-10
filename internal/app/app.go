package app

import (
	"context"
	"log"

	"github.com/dimk00z/go-shortener-praktikum/internal/server"
	"github.com/dimk00z/go-shortener-praktikum/internal/settings"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/storagedi"
	"github.com/dimk00z/go-shortener-praktikum/internal/worker"
)

func StartApp() {
	config := settings.LoadConfig()
	log.Printf("%+v\n", config)

	host := config.Server.Host
	wp := worker.GetWorkersPool(config.Workers)
	defer wp.Close()
	srv := server.NewServer(config.Server.Port, wp)
	storage := storagedi.GetStorage(config.Storage)
	defer func() {
		if err := storage.Close(); err != nil {
			log.Println(err.Error())
		}
	}()
	srv.MountHandlers(host, storage)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		wp.Run(ctx)
	}()
	srv.RunServer(ctx, cancel)
}
