package handlers

import (
	"net/http"

	"github.com/dimk00z/go-shortener-praktikum/internal/util"
)

const securityErrorMessage = "Host is not allowed to connect to endpoint"

// GetStats godoc
// Get db statistics
// @Summary Get db statistics
// @Tags         Internal
// @Accept json
// @Success 200 {object} models.Stat
// @Failure 403
// @Failure 500
// @router /api/internal/stats [get]
func (h *ShortenerHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	// TODO: check network logic
	if h.trustedSubnet != "" {
		util.JSONError(w, securityErrorMessage, http.StatusForbidden)
		return
	}
	status := http.StatusOK
	stat, err := h.Storage.GetStat()
	if err != nil {
		util.JSONError(w, err, http.StatusInternalServerError)
		return
	}
	util.JSONResponse(w, stat, status)
}
