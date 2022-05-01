package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dimk00z/go-shortener-praktikum/internal/storages/storage_interface"
	"github.com/dimk00z/go-shortener-praktikum/internal/util"
	"github.com/go-chi/chi"
)

type RootHandler struct {
	Storage storage_interface.Storage
	host    string
}

func NewRootHandler(host string, st storage_interface.Storage) *RootHandler {
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
	if err != nil {
		log.Println(shortURL, err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else {
		w.Header().Set("Location", shortURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
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
		http.Error(w, "Wrong URL given", http.StatusBadRequest)
		return
	}
	shortURL := h.Storage.SaveURL(URL)

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(fmt.Sprintf("%s/%s", h.host, shortURL)))
	if err != nil {
		log.Fatal(err)
	}
}
