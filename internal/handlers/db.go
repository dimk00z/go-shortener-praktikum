package handlers

import (
	"net/http"

	"github.com/dimk00z/go-shortener-praktikum/internal/util"
)

// PingDB godoc
// Check DB connection
// @Summary Check DB connection
// @Tags         DB
// @Produce json
// @Description Simple pinger
// @router /ping [get]
func (h *ShortenerHandler) PingDB(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("%T\n", h.Storage)
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
