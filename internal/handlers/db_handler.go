package handlers

import (
	"log"
	"net/http"

	"github.com/dimk00z/go-shortener-praktikum/internal/storages/database"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/storageinterface"
	"github.com/dimk00z/go-shortener-praktikum/internal/util"
)

type DBHandler struct {
	host    string
	Storage storageinterface.Storage
}

func NewDBHandler(host string, st storageinterface.Storage) *DBHandler {
	return &DBHandler{
		Storage: st,
		host:    host,
	}
}

func (h *DBHandler) PingDB(w http.ResponseWriter, r *http.Request) {
	log.Printf("%T\n", h.Storage)
	switch h.Storage.(type) {
	case *database.DataBaseStorage:
		util.JSONResponse(w, struct {
			Message string `json:"message"`
		}{Message: "Correct storage"}, http.StatusOK)
	default:
		util.JSONError(w, "Wrong storage loaded", http.StatusInternalServerError)
	}
}
