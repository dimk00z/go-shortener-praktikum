package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/dimk00z/go-shortener-praktikum/internal/middleware/cookie"
	"github.com/dimk00z/go-shortener-praktikum/internal/models"
	"github.com/dimk00z/go-shortener-praktikum/internal/shortenererrors"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/storageerrors"
	"github.com/dimk00z/go-shortener-praktikum/internal/util"
	"github.com/go-chi/chi"
)

// GetByShortURL godoc
// @Summary      Get by shortURL
// @Description  get saved URL
// @Tags         API
// @Param        shortURL   path      string  true  "shortURL"
// @Success      307
// @Failure      404
// @Failure      410
// @Failure      500
// @Router       /{shortURL} [get]
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

// PostShortURL godoc
// @Summary      Save url
// @Description  post URL for saving
// @Tags         API
// @Param url body string true "URLforSave"
// @Success      307
// @Failure      400
// @Failure      409
// @Failure      500
// @Router       / [post]
func (h ShortenerHandler) PostShortURL(w http.ResponseWriter, r *http.Request) {

	if err := util.RequestBodyCheck(w, r); err != nil {
		return
	}

	body, err := io.ReadAll(r.Body)
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

type jsonResult struct {
	Result string `json:"result"`
}

// SaveJSON godoc
// @Summary      Save url by json
// @Description  post URL for saving
// @Accept json
// @Produce json
// @Tags         API
// @Param URL body models.URLRequest true "URL for shorting"
// @Success 201 {object} jsonResult
// @Failure 400
// @Failure 409
// @Failure 500
// @router /api/shortener [post]
func (h ShortenerHandler) SaveJSON(w http.ResponseWriter, r *http.Request) {
	if err := util.RequestBodyCheck(w, r); err != nil {
		return
	}

	var u models.URLRequest

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		util.JSONError(w, err.Error(), http.StatusBadRequest)
		h.l.Debug("ShortenerHandler - SaveJSON - json.NewDecoder: %w", err)
		return
	}
	if u.URL == "" {
		errMessage := "Request doesn't contain url field"
		util.JSONError(w, errMessage, http.StatusBadRequest)
		h.l.Debug("ShortenerHandler - SaveJSON: %s", errMessage)

		return
	}
	if !util.IsURL(u.URL) {
		errMessage := "Wrong URL given -" + u.URL

		util.JSONError(w, errMessage, http.StatusBadRequest)
		h.l.Debug("ShortenerHandler - SaveJSON: %s", errMessage)

		return
	}
	shortURL := util.ShortenLink(u.URL)
	userIDCtx := r.Context().Value(cookie.UserIDCtxName).(string)

	err := h.Storage.SaveURL(u.URL, shortURL, userIDCtx)
	resultStatus := http.StatusCreated
	if errors.Is(err, storageerrors.ErrURLAlreadySave) {
		resultStatus = http.StatusConflict
	}
	util.JSONResponse(w, jsonResult{
		Result: fmt.Sprintf("%s/%s", h.host, shortURL),
	}, resultStatus)

}

// SaveBatch godoc
// @Summary      Save url by json
// @Description  post URL for saving
// @Accept json
// @Produce json
// @Tags         API
// @Param batchURLs body models.BatchURLs true "URLs for shorting"
// @Success 201 {object} jsonResult
// @Failure 400
// @Failure 409
// @Failure 500
// @router /api/shortener/batch [post]
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
