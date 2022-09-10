package server

import (
	"github.com/dimk00z/go-shortener-praktikum/internal/grpc/interceptors"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/storageinterface"
	"github.com/dimk00z/go-shortener-praktikum/internal/worker"
	"github.com/dimk00z/go-shortener-praktikum/pkg/logger"
)

type Service struct {
	Interceptor   *interceptors.Interceptor
	st            storageinterface.Storage
	wp            worker.IWorkerPool
	l             *logger.Logger
	secretKey     string
	trustedSubnet string
	host          string
}

type ServiceOptions func(*Service)

func SetStorage(st storageinterface.Storage) ServiceOptions {
	return func(s *Service) {
		s.st = st
	}
}

func SetLogger(l *logger.Logger) ServiceOptions {
	return func(s *Service) {
		s.l = l
		interceptors.SetInterceptorLogger(l)(s.Interceptor)
	}
}

func SetWorkerPool(wp worker.IWorkerPool) ServiceOptions {
	return func(s *Service) {
		s.wp = wp
	}
}
func SetSecretKey(secretKey string) ServiceOptions {
	return func(s *Service) {
		s.secretKey = secretKey
		interceptors.SetSecretKey(secretKey)(s.Interceptor)
	}
}
func SetTrustedSubnet(trustedSubnet string) ServiceOptions {
	return func(s *Service) {
		s.trustedSubnet = trustedSubnet
	}
}

func SetShortenerHost(host string) ServiceOptions {
	return func(s *Service) {
		s.host = host
	}
}
