package app

import (
	"context"
	"log"

	"github.com/dimk00z/go-shortener-praktikum/internal/server"
	"github.com/dimk00z/go-shortener-praktikum/internal/settings"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/storagedi"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/storageinterface"
	"github.com/dimk00z/go-shortener-praktikum/internal/worker"
)

func StartApp() {
	config := settings.LoadConfig()
	log.Printf("%+v\n", config)

	host := config.Server.Host
	wp := worker.GetWorkersPool(config.Workers)
	server := server.NewServer(config.Server.Port, wp)
	storage := storagedi.GetStorage(config.Storage)

	server.MountHandlers(host, storage)
	defer shutDown(wp, storage, server)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		wp.Run(ctx)
	}()
	server.RunServer(ctx, cancel, storage)
}

func shutDown(wp worker.IWorkerPool, st storageinterface.Storage, s *server.ShortenerServer) {
	wp.Close()
	st.Close()
	s.ShutDown()
}
