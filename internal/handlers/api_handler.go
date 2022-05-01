package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dimk00z/go-shortener-praktikum/internal/storages/storage_interface"
	"github.com/dimk00z/go-shortener-praktikum/internal/util"
)

type URLRequest struct {
	URL string `json:"url"`
}

type ShortenerAPIHandler struct {
	Storage storage_interface.Storage
	host    string
}

func NewShortenerAPIHandler(host string, st storage_interface.Storage) *ShortenerAPIHandler {
	return &ShortenerAPIHandler{
		Storage: st,
		host:    host,
	}
}

func (h ShortenerAPIHandler) SaveJSON(w http.ResponseWriter, r *http.Request) {
	if err := util.RequestBodyCheck(w, r); err != nil {
		return
	}

	var u URLRequest // целевой объект

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		util.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if u.URL == "" {
		util.JSONError(w, "Request doesn't contain url field", http.StatusBadRequest)
		return
	}
	if !util.IsURL(u.URL) {
		util.JSONError(w, "Wrong URL given", http.StatusBadRequest)
		return
	}
	shortURL := h.Storage.SaveURL(u.URL)
	util.JSONResponse(w, struct {
		Result string `json:"result"`
	}{
		Result: fmt.Sprintf("%s/%s", h.host, shortURL),
	}, http.StatusCreated)

}
