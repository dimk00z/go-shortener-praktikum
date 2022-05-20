package handlers

import (
	"log"
	"net/http"

	"github.com/dimk00z/go-shortener-praktikum/internal/middleware/cookie"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/storageinterface"
	"github.com/dimk00z/go-shortener-praktikum/internal/util"
)

type UserHandler struct {
	Storage storageinterface.Storage
	host    string
}

func NewUserHandler(host string, st storageinterface.Storage) *UserHandler {
	return &UserHandler{
		Storage: st,
		host:    host,
	}
}

func (h UserHandler) GetUserURLs(w http.ResponseWriter, r *http.Request) {

	type result struct {
		ShortURL string `json:"short_url"`
		URL      string `json:"original_url"`
	}
	resultStatus := http.StatusOK
	userIDCtx := r.Context().Value(cookie.UserIDCtxName).(string)
	userURLs, err := h.Storage.GetUserURLs(userIDCtx)
	if err != nil {
		log.Println(err)
		resultStatus = http.StatusNoContent
	}
	results := make([]result, len(userURLs))
	for index, userURL := range userURLs {
		results[index] = result{
			ShortURL: h.host + "/" + userURL.ShortURL,
			URL:      userURL.URL,
		}
	}
	log.Println(results)
	util.JSONResponse(w, results, resultStatus)

}
