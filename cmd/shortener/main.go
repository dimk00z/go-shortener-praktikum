package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dimk00z/go-shortener-praktikum/internal/handlers"
)

func oldMain() {
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
func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	shortenerAddress := ":8080"
	mux := http.NewServeMux()
	rootHandler := handlers.NewRootHandler()
	mux.Handle("/", rootHandler)
	go func() {
		log.Fatal(http.ListenAndServe(shortenerAddress, mux))
	}()
	killSignal := <-interrupt
	switch killSignal {
	case os.Interrupt:
		log.Print("Got SIGINT...")
	case syscall.SIGTERM:
		log.Print("Got SIGTERM...")
	}
	log.Print("Done")
}
