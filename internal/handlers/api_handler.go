package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/dimk00z/go-shortener-praktikum/internal/middleware/cookie"
	"github.com/dimk00z/go-shortener-praktikum/internal/models"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/storageinterface"
	"github.com/dimk00z/go-shortener-praktikum/internal/util"
)

type URLRequest struct {
	URL string `json:"url"`
}

type ShortenerAPIHandler struct {
	Storage storageinterface.Storage
	host    string
}

func NewShortenerAPIHandler(host string, st storageinterface.Storage) *ShortenerAPIHandler {
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
		util.JSONError(w, "Wrong URL given -"+u.URL, http.StatusBadRequest)
		return
	}
	shortURL := util.ShortenLink(u.URL)
	userIDCtx := r.Context().Value(cookie.UserIDCtxName).(string)

	h.Storage.SaveURL(u.URL, shortURL, userIDCtx)

	util.JSONResponse(w, struct {
		Result string `json:"result"`
	}{
		Result: fmt.Sprintf("%s/%s", h.host, shortURL),
	}, http.StatusCreated)

}

// Sprint 3 Increment 12
func (h ShortenerAPIHandler) SaveBatch(w http.ResponseWriter, r *http.Request) {
	if err := util.RequestBodyCheck(w, r); err != nil {
		return
	}
	var requestData models.BatchURLs
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		util.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	for index, field := range requestData {
		requestData[index].ShortURL = util.ShortenLink(field.OriginalURL)
	}
	log.Println(requestData)
	userIDCtx := r.Context().Value(cookie.UserIDCtxName).(string)
	resultURLs, err := h.Storage.SaveBatch(requestData, userIDCtx)
	if err != nil {
		util.JSONError(w, err, http.StatusBadRequest)
		return
	}
	util.JSONResponse(w, resultURLs, http.StatusCreated)
}
