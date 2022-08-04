package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/dimk00z/go-shortener-praktikum/internal/middleware/cookie"
	"github.com/dimk00z/go-shortener-praktikum/internal/models"
	"github.com/dimk00z/go-shortener-praktikum/internal/util"
)

func (h ShortenerHandler) GetUserURLs(w http.ResponseWriter, r *http.Request) {

	resultStatus := http.StatusOK
	userIDCtx := r.Context().Value(cookie.UserIDCtxName).(string)
	userURLs, err := h.Storage.GetUserURLs(userIDCtx)
	if err != nil {
		h.l.Debug(err)
		resultStatus = http.StatusNoContent
	}
	results := make([]models.UserURL, len(userURLs))
	for index, userURL := range userURLs {
		results[index] = models.UserURL{
			ShortURL: h.host + "/" + userURL.ShortURL,
			URL:      userURL.URL,
		}
	}
	h.l.Debug(results)
	util.JSONResponse(w, results, resultStatus)

}

func (h ShortenerHandler) DeleteUserURLs(w http.ResponseWriter, r *http.Request) {

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
	h.l.Debug(shortURLs)

	deleteBatchTask := func(ctx context.Context) error {
		return h.Storage.DeleteBatch(ctx, shortURLs, userIDCtx)
	}
	h.wp.Push(deleteBatchTask)
	w.WriteHeader(resultStatus)

}
