package main

import (
	"log"
	"net/http"

	"github.com/dimk00z/go-shortener-praktikum/internal/handlers"
)

func main() {
	handler1 := handlers.MyHandler{
		Templ: []byte("Hola, Mundo"),
	}

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler1,
	}
	log.Fatal(server.ListenAndServe())
}
