package handlers

import (
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/storageinterface"
	"github.com/dimk00z/go-shortener-praktikum/internal/worker"
)

type ShortenerHandler struct {
	Storage storageinterface.Storage
	host    string
	wp      worker.IWorkerPool
}

func NewShortenerHandler(host string, st storageinterface.Storage, wp worker.IWorkerPool) *ShortenerHandler {
	return &ShortenerHandler{
		Storage: st,
		host:    host,
		wp:      wp,
	}
}
