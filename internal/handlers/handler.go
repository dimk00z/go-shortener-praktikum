package handlers

import (
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/storageinterface"
	"github.com/dimk00z/go-shortener-praktikum/internal/worker"
	"github.com/dimk00z/go-shortener-praktikum/pkg/logger"
)

type ShortenerHandler struct {
	Storage       storageinterface.Storage
	host          string
	wp            worker.IWorkerPool
	l             *logger.Logger
	trustedSubnet string
}

func NewShortenerHandler() *ShortenerHandler {
	return &ShortenerHandler{}
}

type ShortenerOptions func(*ShortenerHandler)

func SetStorage(st storageinterface.Storage) ShortenerOptions {
	return func(s *ShortenerHandler) {
		s.Storage = st
	}
}

func SetTrustedSubnet(trustedSubnet string) ShortenerOptions {
	return func(s *ShortenerHandler) {
		s.trustedSubnet = trustedSubnet
	}
}
func SetHost(host string) ShortenerOptions {
	return func(s *ShortenerHandler) {
		s.host = host
	}
}
func SetWorkerPool(wp worker.IWorkerPool) ShortenerOptions {
	return func(s *ShortenerHandler) {
		s.wp = wp
	}
}
func SetLoger(l *logger.Logger) ShortenerOptions {
	return func(s *ShortenerHandler) {
		s.l = l
	}
}
