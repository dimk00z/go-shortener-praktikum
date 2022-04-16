package handlers

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/dimk00z/go-shortener-praktikum/internal/storage"
	"github.com/dimk00z/go-shortener-praktikum/internal/util"
)

type RootHandler struct {
	storage storage.URLsStorage
}

func NewRootHandler() *RootHandler {
	return &RootHandler{
		storage: *storage.NewStorage(),
	}
}

func (h RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		{
			shortURL := strings.Replace(r.URL.Path, "/", "", 1)
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
	case "POST":
		{
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
			w.Write([]byte(shortURL))
		}
	}
}
