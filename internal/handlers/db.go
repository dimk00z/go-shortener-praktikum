package handlers

import (
	"log"
	"net/http"

	"github.com/dimk00z/go-shortener-praktikum/internal/util"
)

func (h *ShortenerHandler) PingDB(w http.ResponseWriter, r *http.Request) {
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

}
