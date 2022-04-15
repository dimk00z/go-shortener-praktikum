package main

import (
	"log"
	"net/http"
	"time"

	"github.com/dimk00z/go-shortener-praktikum/internal/handlers"
)

func main() {
	handler1 := handlers.MyHandler{
		// Templ: []byte("Hola, Mundo"),
	}

	mux := http.NewServeMux()
	mux.Handle("/", handler1)

	nameHandler := handlers.HameHandler{}
	mux.Handle("/hello/", nameHandler)

	th := handlers.TimeHandler{Format: time.RFC1123}
	mux.Handle("/time", th)
	mux.Handle("/ya", http.RedirectHandler("https://ya.ru/", 301))
	log.Fatal(http.ListenAndServe(":8080", mux))
}
