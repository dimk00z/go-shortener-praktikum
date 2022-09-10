package server

import (
	"github.com/dimk00z/go-shortener-praktikum/internal/grpc/interceptors"
	pb "github.com/dimk00z/go-shortener-praktikum/internal/grpc/proto"
)

type GRPCServer struct {
	pb.UnimplementedShortenerServer
	Service *Service
}

func NewGRPCServer() *GRPCServer {
	return &GRPCServer{Service: &Service{Interceptor: &interceptors.Interceptor{}}}
}
