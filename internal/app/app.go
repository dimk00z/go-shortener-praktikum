package app

import (
	"context"
	"net"

	"github.com/dimk00z/go-shortener-praktikum/config"
	pb "github.com/dimk00z/go-shortener-praktikum/internal/grpc/proto"
	grpcServer "github.com/dimk00z/go-shortener-praktikum/internal/grpc/server"
	"github.com/dimk00z/go-shortener-praktikum/internal/server"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/storagedi"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/storageinterface"
	"github.com/dimk00z/go-shortener-praktikum/internal/worker"
	"github.com/dimk00z/go-shortener-praktikum/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
	server.SetTrustedSubnet(config.Security.TrustedSubnet)(s)

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

	if config.GRPC.EnableGRPC {
		setGRPC(storage, wp, l, config.Security.SecretKey, config.GRPC.Port)
	}
	s.RunServer(ctx, cancel, storage)
}

func setGRPC(
	st storageinterface.Storage,
	wp worker.IWorkerPool,
	l *logger.Logger,
	secretKey string, grpcPort string) {
	server := grpcServer.NewGRPCServer()
	opts := []grpcServer.ServiceOptions{
		grpcServer.SetLogger(l),
		grpcServer.SetWorkerPool(wp),
		grpcServer.SetStorage(st),
		grpcServer.SetSecretKey(secretKey),
	}
	for _, opt := range opts {
		opt(server.Service)
	}

	listen, err := net.Listen("tcp", grpcPort)
	if err != nil {
		l.Fatal(err)
	}
	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterShortenerServer(s, server)
	go func() {
		l.Info("setGRPC - gRPC server started on " + grpcPort)
		if err := s.Serve(listen); err != nil {
			l.Fatal(err)
			// TODO: add graceful shutdown
		}
	}()
}

func shutDown(wp worker.IWorkerPool, st storageinterface.Storage, s *server.ShortenerServer) {
	st.Close()
	s.ShutDown()
}
