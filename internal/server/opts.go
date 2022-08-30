package server

import (
	"github.com/dimk00z/go-shortener-praktikum/internal/worker"
	"github.com/dimk00z/go-shortener-praktikum/pkg/logger"
	"github.com/go-chi/chi"
)

type ShortenerServerOptions func(*ShortenerServer)

func SetSecretKey(key string) ShortenerServerOptions {
	return func(s *ShortenerServer) {
		s.secretKey = key
	}
}

func SetRouter(router *chi.Mux) ShortenerServerOptions {
	return func(s *ShortenerServer) {
		s.Router = router
	}
}

func SetPort(port string) ShortenerServerOptions {
	return func(s *ShortenerServer) {
		s.port = port
	}
}

func SetLoger(l *logger.Logger) ShortenerServerOptions {
	return func(s *ShortenerServer) {
		s.l = l
	}
}

func SetWorkersPool(wp worker.IWorkerPool) ShortenerServerOptions {
	return func(s *ShortenerServer) {
		s.wp = wp
	}
}

func SetTLSConfig(certFile string, keyFile string, port string) ShortenerServerOptions {
	return func(s *ShortenerServer) {
		s.tlsConfig = &serverTLSConfig{
			port:     port,
			certFile: certFile,
			keyFile:  keyFile,
		}
	}
}
