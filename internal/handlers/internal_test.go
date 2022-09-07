package handlers

import (
	"net/http"
	"testing"

	"github.com/dimk00z/go-shortener-praktikum/internal/storages/storageinterface"
	"github.com/dimk00z/go-shortener-praktikum/internal/worker"
	"github.com/dimk00z/go-shortener-praktikum/pkg/logger"
)

func TestShortenerHandler_GetStats(t *testing.T) {
	type fields struct {
		Storage storageinterface.Storage
		host    string
		wp      worker.IWorkerPool
		l       *logger.Logger
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &ShortenerHandler{
				Storage: tt.fields.Storage,
				host:    tt.fields.host,
				wp:      tt.fields.wp,
				l:       tt.fields.l,
			}
			h.GetStats(tt.args.w, tt.args.r)
		})
	}
}
