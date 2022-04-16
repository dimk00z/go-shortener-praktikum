package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dimk00z/go-shortener-praktikum/internal/handlers"
)

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
	case syscall.SIGQUIT:
		log.Print("Got SIGQUIT...")
	}
	log.Print("Done")
}
