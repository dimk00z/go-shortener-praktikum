package handlers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dimk00z/go-shortener-praktikum/internal/middleware/cookie"
	"github.com/dimk00z/go-shortener-praktikum/internal/shortenererrors"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/storageerrors"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/storageinterface"
	"github.com/dimk00z/go-shortener-praktikum/internal/util"
	"github.com/go-chi/chi"
)

type RootHandler struct {
	Storage storageinterface.Storage
	host    string
}

func NewRootHandler(host string, st storageinterface.Storage) *RootHandler {
	return &RootHandler{
		Storage: st,
		host:    host,
	}
}

func (h RootHandler) HandleGETRequest(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "shortURL")
	log.Println("Get " + shortURL + " shortURL")
	var err error
	shortURL, err = h.Storage.GetByShortURL(shortURL)
	statusCode := http.StatusTemporaryRedirect
	if err != nil {
		log.Println(shortURL, err)

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

func (h RootHandler) HandlePOSTRequest(w http.ResponseWriter, r *http.Request) {

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
		log.Fatal(err)
	}
}
