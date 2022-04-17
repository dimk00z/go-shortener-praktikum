package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/dimk00z/go-shortener-praktikum/internal/storage"
	"github.com/dimk00z/go-shortener-praktikum/internal/util"
)

type RootHandler struct {
	storage storage.URLStorage
	host    string
}

func NewRootHandler(host string) *RootHandler {
	return &RootHandler{
		storage: *storage.NewStorage(),
		host:    host,
	}
}

func (h RootHandler) HandleGETRequest(w http.ResponseWriter, r *http.Request) {
	shortURL := strings.Replace(r.URL.Path, "/", "", 1)
	log.Println("Get " + shortURL + " shortURL")
	var err error
	shortURL, err = h.storage.GetByShortURL(shortURL)
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
	if r.Body == http.NoBody {
		http.Error(w, "Request should have body", http.StatusBadRequest)
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
	shortURL := h.storage.SaveURL(URL)

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("%s/%s", h.host, shortURL)))
}
