package server

import (
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/storageinterface"
	"github.com/dimk00z/go-shortener-praktikum/internal/worker"
	"github.com/dimk00z/go-shortener-praktikum/pkg/logger"
)

type Service struct {
	st        storageinterface.Storage
	wp        worker.IWorkerPool
	l         *logger.Logger
	secretKey string
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
	}
}
