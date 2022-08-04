package handlers_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/dimk00z/go-shortener-praktikum/config"
	"github.com/dimk00z/go-shortener-praktikum/internal/server"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/storageinterface"
	"github.com/dimk00z/go-shortener-praktikum/internal/worker"
	"github.com/dimk00z/go-shortener-praktikum/pkg/logger"
)

const (
	shortenerPort = ":8080"
	host          = "http://localhost" + shortenerPort
)

func execRequest(req *http.Request, s *server.ShortenerServer) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.Router.ServeHTTP(rr, req)

	return rr
}
func getMockLogger() *logger.Logger {
	l := logger.New("debug")
	return l
}
func getMockWorkersPool() worker.IWorkerPool {
	return worker.GetWorkersPool(getMockLogger(), config.Workers{WorkersNumber: 2, PoolLength: 10})
}

func createMockServer(mockStorage storageinterface.Storage, wp worker.IWorkerPool) *server.ShortenerServer {
	server := server.NewServer(getMockLogger(),
		shortenerPort, wp, "MockSecret")
	server.MountHandlers(host, mockStorage)
	return server
}
