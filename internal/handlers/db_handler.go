package handlers

import (
	"log"
	"net/http"

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
	message := "connection established"
	status := http.StatusOK
	if err := h.Storage.CheckConnection(r.Context()); err != nil {
		message = err.Error()
		status = http.StatusInternalServerError
	}
	util.JSONResponse(w, struct {
		Message string `json:"message"`
	}{Message: message}, status)

	// switch h.Storage.(type) {
	// case *database.DataBaseStorage:
	// 	message := "connection established"
	// 	status := http.StatusOK
	// 	if err := h.Storage.CheckConnection(r.Context()); err != nil {
	// 		message = err.Error()
	// 		status = http.StatusInternalServerError
	// 	}
	// 	util.JSONResponse(w, struct {
	// 		Message string `json:"message"`
	// 	}{Message: message}, status)
	// default:
	// 	util.JSONError(w, "Wrong storage loaded", http.StatusInternalServerError)
	// }
}
