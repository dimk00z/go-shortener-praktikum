package handlers

import (
	"net"
	"net/http"
	"strings"

	"github.com/dimk00z/go-shortener-praktikum/internal/util"
)

const securityErrorMessage = "Host is not allowed to connect to endpoint"

func realIP(r *http.Request) string {
	var ip string

	if tcip := r.Header.Get(http.CanonicalHeaderKey("True-Client-IP")); tcip != "" {
		ip = tcip
	} else if xrip := r.Header.Get(http.CanonicalHeaderKey("X-Real-IP")); xrip != "" {
		ip = xrip
	} else if xff := r.Header.Get(http.CanonicalHeaderKey("X-Forwarded-For")); xff != "" {
		i := strings.Index(xff, ",")
		if i == -1 {
			i = len(xff)
		}
		ip = xff[:i]
	}
	if ip == "" || net.ParseIP(ip) == nil {
		return ""
	}
	return ip
}

func getRealIP(w http.ResponseWriter, r *http.Request) string {
	if rip := realIP(r); rip != "" {
		r.RemoteAddr = rip
		return rip
	}
	return ""

}

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

	if h.trustedSubnet == "" {
		util.JSONError(w, securityErrorMessage, http.StatusForbidden)
		return
	}
	ip := getRealIP(w, r)
	h.l.Debug("Client ip is " + ip)
	if ip == "" {
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
