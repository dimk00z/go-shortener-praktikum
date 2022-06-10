package decompressor

import (
	"compress/gzip"
	"log"
	"net/http"
	"strings"
)

func DecompressHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return

		}
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer func() {
			if err := gz.Close(); err != nil {
				log.Println(err.Error())
			}
		}()
		r.Body = gz
		next.ServeHTTP(w, r)
	})
}
