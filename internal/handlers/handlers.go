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
	storage storage.UrlsStorage
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
			shortUrl := strings.Replace(r.URL.Path, "/", "", 1)
			var err error
			shortUrl, err = h.storage.GetByShortUrl(shortUrl)
			if err != nil {
				log.Println(shortUrl, err)
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			} else {
				w.Header().Set("Location", shortUrl)
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
			url := string(body)
			if !util.IsUrl(url) {
				http.Error(w, "Wrong url given", http.StatusBadRequest)
				return
			}
			shortURL := h.storage.SaveUrl(url)

			w.Header().Add("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(shortURL))
		}
	}
}
