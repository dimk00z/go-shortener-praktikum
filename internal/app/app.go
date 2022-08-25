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
	s := server.NewServer(l, config.Server.Port, wp, config.Security.SecretKey)
	if config.Security.EnableHTTPS {
		server.SetTLSConfig(
			config.Security.CertFile,
			config.Security.KeyFile,
			config.Server.Port,
		)(s)
	}

	if config.Storage.DataSourceName != "" {
		doMigrations(l, config.Storage.DataSourceName)
	}

	storage := storagedi.GetStorage(l, config.Storage)

	s.MountHandlers(host, storage)
	defer shutDown(wp, storage, s)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		wp.Run(ctx)
	}()
	s.RunServer(ctx, cancel, storage)
}

func shutDown(wp worker.IWorkerPool, st storageinterface.Storage, s *server.ShortenerServer) {
	st.Close()
	s.ShutDown()
}
