package server

import (
	"net"

	"github.com/dimk00z/go-shortener-praktikum/config"
	"github.com/dimk00z/go-shortener-praktikum/internal/grpc/interceptors"
	pb "github.com/dimk00z/go-shortener-praktikum/internal/grpc/proto"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/storageinterface"
	"github.com/dimk00z/go-shortener-praktikum/internal/worker"
	"github.com/dimk00z/go-shortener-praktikum/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	pb.UnimplementedShortenerServer
	Service *Service
}

func newGRPCServer() *GRPCServer {
	return &GRPCServer{Service: &Service{Interceptor: &interceptors.Interceptor{}}}
}
func SetGRPC(
	st storageinterface.Storage,
	wp worker.IWorkerPool,
	l *logger.Logger,
	config *config.Config) *grpc.Server {
	server := newGRPCServer()
	opts := []ServiceOptions{
		SetLogger(l),
		SetWorkerPool(wp),
		SetStorage(st),
		SetSecretKey(config.Security.SecretKey),
		SetTrustedSubnet(config.Security.TrustedSubnet),
		SetShortenerHost(config.Server.Host),
	}
	for _, opt := range opts {
		opt(server.Service)
	}

	listen, err := net.Listen("tcp", config.GRPC.Port)
	if err != nil {
		l.Fatal(err)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(server.Service.Interceptor.AuthInterceptor))
	reflection.Register(s)
	pb.RegisterShortenerServer(s, server)
	go func() {
		l.Info("setGRPC - gRPC server started on " + config.GRPC.Port)
		if err := s.Serve(listen); err != nil {
			l.Fatal(err)
		}
	}()
	return s
}
