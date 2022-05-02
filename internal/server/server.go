package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dimk00z/go-shortener-praktikum/internal/handlers"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/storageinterface"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// type Handler interface {
// 	HandlePOSTRequest(w http.ResponseWriter, r *http.Request)
// 	HandleGETRequest(w http.ResponseWriter, r *http.Request)
// }

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
func (s *ShortenerServer) MountHandlers(host string, st storageinterface.Storage) {
	// Mount all Middleware here
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)

	//Инкремент 8: вот тут не понял, обязательно ли писать компессор ибо он в коробке Chi есть?
	compressedTypes := []string{
		"text/html",
		"text/css",
		"text/plain",
		"text/xml",
		"application/javascript",
		"application/json",
	}
	s.Router.Use(middleware.Compress(5, compressedTypes...))

	// Mount all handlers here
	// Sprint 1
	s.Router.Route("/", func(r chi.Router) {
		h := handlers.NewRootHandler(host,
			st)

		r.Post("/", h.HandlePOSTRequest)
		r.Get("/{shortURL}", h.HandleGETRequest)
	})
	// Sprint 2
	shortenerRouter := chi.NewRouter()
	shortenerRouter.Route("/", func(r chi.Router) {
		h := handlers.NewShortenerAPIHandler(host,
			st)
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
