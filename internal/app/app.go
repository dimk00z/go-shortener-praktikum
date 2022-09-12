package app

import (
	"context"

	"github.com/dimk00z/go-shortener-praktikum/config"
	grpcServer "github.com/dimk00z/go-shortener-praktikum/internal/grpc/server"
	"github.com/dimk00z/go-shortener-praktikum/internal/server"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/storagedi"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/storageinterface"
	"github.com/dimk00z/go-shortener-praktikum/internal/worker"
	"github.com/dimk00z/go-shortener-praktikum/pkg/logger"
	"google.golang.org/grpc"
)

func StartApp(config *config.Config) {

	l := logger.New(config.Log.Level)
	ctx, cancel := context.WithCancel(context.Background())

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
	server.SetTrustedSubnet(config.Security.TrustedSubnet)(s)

	if config.Storage.DataSourceName != "" {
		doMigrations(l, config.Storage.DataSourceName)
	}

	storage := storagedi.GetStorage(l, config.Storage)

	s.MountHandlers(host, storage)
	var grpcShortenerServer *grpc.Server
	if config.GRPC.EnableGRPC {
		grpcShortenerServer = grpcServer.SetGRPC(storage, wp, l, config)
	}
	defer shutDown(l, wp, storage, s, grpcShortenerServer)

	go func() {
		wp.Run(ctx)
	}()

	s.RunServer(ctx, cancel, storage)
}

func shutDown(
	l *logger.Logger,
	wp worker.IWorkerPool,
	st storageinterface.Storage,
	s *server.ShortenerServer,
	grpcShortenerServer *grpc.Server) {
	st.Close()
	s.ShutDown()
	if grpcShortenerServer != nil {
		grpcShortenerServer.GracefulStop()
	}
	l.Debug("Services stopped grafefully")
}
