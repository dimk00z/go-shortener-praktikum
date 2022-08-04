package app

import (
	"context"

	"github.com/dimk00z/go-shortener-praktikum/config"
	"github.com/dimk00z/go-shortener-praktikum/internal/server"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/storagedi"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/storageinterface"
	"github.com/dimk00z/go-shortener-praktikum/internal/worker"
	"github.com/dimk00z/go-shortener-praktikum/pkg/logger"
)

func StartApp(config *config.Config) {

	l := logger.New(config.Log.Level)

	l.Debug("%+v\n", config)

	host := config.Server.Host

	wp := worker.GetWorkersPool(l, config.Workers)
	server := server.NewServer(l, config.Server.Port, wp, config.Security.SecretKey)

	if config.Storage.DataSourceName != "" {
		doMigrations(config.Storage.DataSourceName, l)
	}

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
	st.Close()
	s.ShutDown()
}
