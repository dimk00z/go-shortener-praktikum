package handlers

import "net/http"

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Hello, World</h1>"))
}

type MyHandler struct {
	Templ []byte
}

func (h MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write(h.Templ)
}
