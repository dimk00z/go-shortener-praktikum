package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dimk00z/go-shortener-praktikum/internal/middleware/cookie"
	"github.com/dimk00z/go-shortener-praktikum/internal/models"
	"github.com/dimk00z/go-shortener-praktikum/internal/shortenererrors"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/storageerrors"
	"github.com/dimk00z/go-shortener-praktikum/internal/util"
	"github.com/go-chi/chi"
)

func (h ShortenerHandler) GetByShortURL(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "shortURL")
	h.l.Debug("Get " + shortURL + " shortURL")
	var err error
	shortURL, err = h.Storage.GetByShortURL(shortURL)
	statusCode := http.StatusTemporaryRedirect
	if err != nil {
		h.l.Debug(shortURL, err)

		if errors.Is(err, shortenererrors.ErrURLNotFound) {
			statusCode = http.StatusNotFound
		}
		if errors.Is(err, shortenererrors.ErrURLDeleted) {
			statusCode = http.StatusGone
		}
		http.Error(w, err.Error(), statusCode)
		return
	}
	w.Header().Set("Location", shortURL)
	w.WriteHeader(statusCode)

}

func (h ShortenerHandler) PostShortURL(w http.ResponseWriter, r *http.Request) {

	if err := util.RequestBodyCheck(w, r); err != nil {
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	URL := string(body)
	if !util.IsURL(URL) {
		http.Error(w, "Wrong URL given -"+URL, http.StatusBadRequest)
		return
	}
	shortURL := util.ShortenLink(URL)
	userIDCtx := r.Context().Value(cookie.UserIDCtxName).(string)
	err = h.Storage.SaveURL(URL, shortURL, userIDCtx)
	resultStatus := http.StatusCreated
	if errors.Is(err, storageerrors.ErrURLAlreadySave) {
		resultStatus = http.StatusConflict
	}

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(resultStatus)
	_, err = w.Write([]byte(fmt.Sprintf("%s/%s", h.host, shortURL)))
	if err != nil {
		h.l.Debug(err)
	}
}

func (h ShortenerHandler) SaveJSON(w http.ResponseWriter, r *http.Request) {
	if err := util.RequestBodyCheck(w, r); err != nil {
		return
	}

	var u models.URLRequest

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

	err := h.Storage.SaveURL(u.URL, shortURL, userIDCtx)
	resultStatus := http.StatusCreated
	if errors.Is(err, storageerrors.ErrURLAlreadySave) {
		resultStatus = http.StatusConflict
	}
	util.JSONResponse(w, struct {
		Result string `json:"result"`
	}{
		Result: fmt.Sprintf("%s/%s", h.host, shortURL),
	}, resultStatus)

}

// Sprint 3 Increment 12
func (h ShortenerHandler) SaveBatch(w http.ResponseWriter, r *http.Request) {
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
	h.l.Debug(requestData)
	userIDCtx := r.Context().Value(cookie.UserIDCtxName).(string)
	resultURLs, err := h.Storage.SaveBatch(requestData, userIDCtx)
	if err != nil {
		util.JSONError(w, err, http.StatusBadRequest)
		return
	}
	for index := range resultURLs {
		resultURLs[index].ShortURL = fmt.Sprintf("%s/%s", h.host, resultURLs[index].ShortURL)
	}
	util.JSONResponse(w, resultURLs, http.StatusCreated)
}
