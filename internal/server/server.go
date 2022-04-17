package server

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Handler interface {
	HandlePOSTRequest(w http.ResponseWriter, r *http.Request)
	HandleGETRequest(w http.ResponseWriter, r *http.Request)
}

type ShortenerServer struct {
	port   string
	Router *chi.Mux
}

func NewServer(port string) *ShortenerServer {
	return &ShortenerServer{
		port:   port,
		Router: chi.NewRouter(),
	}
}
func (s *ShortenerServer) MountHandlers(r Handler) {
	// Mount all Middleware here
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)

	// Mount all handlers here

	s.Router.Post("/", r.HandlePOSTRequest)
	s.Router.Get("/{shortURL}", r.HandleGETRequest)

}

func (s ShortenerServer) RunServer() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	go func() {
		log.Println("Server started at " + s.port)
		http.ListenAndServe(s.port, s.Router)
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
