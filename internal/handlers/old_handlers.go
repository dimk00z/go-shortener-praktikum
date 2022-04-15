package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type MyHandler struct {
	Templ string
}

type Subj struct {
	Product string `json:"name"`
	Price   int    `json:"price"`
}

func (h MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Templ = "<h1>Hello, World</h1>"
	fmt.Println(r.Header)
	switch r.Method {
	case "GET":
		q := r.URL.Query().Get("query")
		h.Templ = "GET:" + h.Templ
		if q != "" {
			h.Templ = h.Templ + q

		}
		// w.Write(h.Templ)
		fmt.Fprintln(w, string(h.Templ))
	case "POST":

		w.Header().Set("content-type", "application/json")
		// устанавливаем статус-код 200
		w.WriteHeader(http.StatusOK)
		// собираем данные
		subj := Subj{"Milk", 50}
		// кодируем JSON
		resp, err := json.Marshal(subj)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		// пишем тело ответа
		w.Write(resp)

	}
}

type TimeHandler struct {
	Format string
}

func (th TimeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tm := time.Now().Format(th.Format)
	w.Write([]byte("The time is: " + tm))
}

type HameHandler struct {
}

func (nm HameHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	name := strings.Replace(r.URL.Path, "/hello/", "", 1)

	fmt.Fprintln(w, fmt.Sprintf("Hello %s", name))

}
