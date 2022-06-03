package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/dimk00z/go-shortener-praktikum/internal/middleware/cookie"
	"github.com/dimk00z/go-shortener-praktikum/internal/models"
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

	resultStatus := http.StatusOK
	userIDCtx := r.Context().Value(cookie.UserIDCtxName).(string)
	userURLs, err := h.Storage.GetUserURLs(userIDCtx)
	if err != nil {
		log.Println(err)
		resultStatus = http.StatusNoContent
	}
	results := make([]models.UserURL, len(userURLs))
	for index, userURL := range userURLs {
		results[index] = models.UserURL{
			ShortURL: h.host + "/" + userURL.ShortURL,
			URL:      userURL.URL,
		}
	}
	log.Println(results)
	util.JSONResponse(w, results, resultStatus)

}

func (h UserHandler) DeleteUserURLs(w http.ResponseWriter, r *http.Request) {

	resultStatus := http.StatusAccepted
	if err := util.RequestBodyCheck(w, r); err != nil {
		return
	}
	userIDCtx := r.Context().Value(cookie.UserIDCtxName).(string)
	var shortURLs models.BatchForDelete
	if err := json.NewDecoder(r.Body).Decode(&shortURLs); err != nil {
		util.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println(shortURLs)
	w.WriteHeader(resultStatus)

	// TODO добавить удаление сюда!
	ctx := context.Background()
	err := h.Storage.DeleteBatch(ctx, shortURLs, userIDCtx)
	if err != nil {
		log.Println(err)
	}
}
