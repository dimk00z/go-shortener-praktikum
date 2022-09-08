package handlers

import (
	"errors"
	"net"
	"net/http"
	"strings"

	"github.com/dimk00z/go-shortener-praktikum/internal/util"
	"golang.org/x/exp/slices"
)

func realIP(r *http.Request) string {
	// logic from https://github.com/go-chi/chi/blob/master/middleware/realip.go
	var ip string

	if tcip := r.Header.Get("True-Client-IP"); tcip != "" {
		ip = tcip
	} else if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		ip = xrip
	} else if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
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
func (h *ShortenerHandler) ProtectedByTrustedNetwork(allowedMethods []string, next http.Handler) http.Handler {
	var securityRealIPNotGivenError = errors.New("real ip was not given")
	var securityRealIPError = errors.New("host is not allowed to connect to endpoint")
	var securityTrustedNetworkError = errors.New("trusted network should be given")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !slices.Contains(allowedMethods, r.Method) {
			util.JSONError(w, "Wrong method", http.StatusForbidden)
			return
		}
		if h.trustedSubnet == "" {
			util.JSONError(w, securityTrustedNetworkError.Error(), http.StatusForbidden)
			return
		}
		ip := getRealIP(w, r)
		h.l.Debug("Client ip is " + ip)
		if ip == "" {
			util.JSONError(w, securityRealIPNotGivenError.Error(), http.StatusForbidden)
			return
		}
		_, ipnet, err := net.ParseCIDR(h.trustedSubnet)
		if err != nil {
			util.JSONError(w, err, http.StatusInternalServerError)
		}

		if !ipnet.Contains(net.ParseIP(ip)) {
			util.JSONError(w, securityRealIPError.Error(), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})

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
	status := http.StatusOK
	stat, err := h.Storage.GetStat()
	if err != nil {
		util.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	util.JSONResponse(w, stat, status)
}
