package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dimk00z/go-shortener-praktikum/internal/handlers"
	"github.com/dimk00z/go-shortener-praktikum/internal/storage"
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
func (s *ShortenerServer) MountHandlers(host string, getStorage func() (*storage.URLStorage, error)) {
	// Mount all Middleware here
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)
	// Mount all handlers here
	// Sprint 1
	s.Router.Route("/", func(r chi.Router) {
		h := handlers.NewRootHandler(host,
			getStorage)

		r.Post("/", h.HandlePOSTRequest)
		r.Get("/{shortURL}", h.HandleGETRequest)
	})
	// Sprint 2
	shortenerRouter := chi.NewRouter()
	shortenerRouter.Route("/", func(r chi.Router) {
		h := handlers.NewShortenerAPIHandler(host,
			getStorage)
		r.Post("/", h.SaveJSON)
	})
	apiRouter := chi.NewRouter()
	apiRouter.Mount("/shorten", shortenerRouter)
	s.Router.Mount("/api", apiRouter)
}

func (s ShortenerServer) RunServer(ctx context.Context, cancel context.CancelFunc) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	go func() {
		log.Println("Server started at " + s.port)
		err := http.ListenAndServe(s.port, s.Router)
		if err != nil {
			log.Println(err)
		}
		cancel()
	}()
	select {
	case killSignal := <-interrupt:
		log.Print("Got ", killSignal)
	case <-ctx.Done():
	}
	log.Print("Done")
}
