package handlers

import (
	"fmt"
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
			fmt.Println(shortUrl, err)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
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
			url := r.FormValue("url")
			if url == "" {
				http.Error(w, "Url field is required", http.StatusBadRequest)
				return
			} else if !util.IsUrl(url) {
				http.Error(w, "Wrong url given", http.StatusBadRequest)
				return
			}
			shortUrl := h.storage.SaveUrl(url)
			w.Header().Add("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(shortUrl))
		}
	}
}
