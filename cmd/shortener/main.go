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
	shortenerPort := ":8080"
	host := "http://localhost" + shortenerPort
	mux := http.NewServeMux()
	rootHandler := handlers.NewRootHandler(host)
	mux.Handle("/", rootHandler)
	go func() {
		log.Fatal(http.ListenAndServe(shortenerPort, mux))
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
