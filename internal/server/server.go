package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dimk00z/go-shortener-praktikum/internal/handlers"
	"github.com/dimk00z/go-shortener-praktikum/internal/middleware/cookie"
	"github.com/dimk00z/go-shortener-praktikum/internal/middleware/decompressor"
	"github.com/dimk00z/go-shortener-praktikum/internal/settings"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/storageinterface"
	"github.com/dimk00z/go-shortener-praktikum/internal/worker"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type ShortenerServer struct {
	port   string
	Router *chi.Mux
	wp     worker.IWorkerPool
}

func NewServer(port string, wp worker.IWorkerPool) *ShortenerServer {
	return &ShortenerServer{
		port:   port,
		Router: chi.NewRouter(),
		wp:     wp,
	}
}
func (s *ShortenerServer) mountMiddleware() {
	// Mount all Middleware here
	cookieHandler := cookie.CookieHandler{
		SecretKey: settings.LoadConfig().Security.SecretKey,
	}
	decompressHandler := decompressor.DecompressHandler{}

	middlewareHadlers := []func(http.Handler) http.Handler{
		middleware.RequestID,
		middleware.Logger,
		middleware.Recoverer,
		decompressHandler.Handle,
		cookieHandler.Handle,
		middleware.Compress(5),
	}

	for _, handler := range middlewareHadlers {
		s.Router.Use(handler)
	}
}
func (s *ShortenerServer) MountHandlers(host string, st storageinterface.Storage) {

	s.mountMiddleware()
	h := handlers.NewShortenerHandler()
	handlerOptions := []handlers.ShortenerOptions{
		handlers.SetHost(host),
		handlers.SetStorage(st),
		handlers.SetWorkerPool(s.wp),
	}

	for _, opt := range handlerOptions {
		opt(h)
	}

	s.Router.Route("/", func(r chi.Router) {
		r.Post("/", h.PostShortURL)
		r.Get("/{shortURL}", h.GetByShortURL)
	})

	apiRouter := chi.NewRouter()
	apiRouter.Mount("/shorten", chi.NewRouter().Route("/", func(r chi.Router) {
		r.Post("/", h.SaveJSON)
		r.Post("/batch", h.SaveBatch)
	}))

	apiRouter.Mount("/user", chi.NewRouter().Route("/", func(r chi.Router) {
		r.Get("/urls", h.GetUserURLs)
		r.Delete("/urls", h.DeleteUserURLs)
	}))

	s.Router.Mount("/api", apiRouter)

	s.Router.Mount("/ping",
		chi.NewRouter().Route("/", func(r chi.Router) {
			r.Get("/", h.PingDB)
		}))
}

func (s ShortenerServer) RunServer(ctx context.Context, cancel context.CancelFunc, storage storageinterface.Storage) {
	interrupt := make(chan os.Signal, 1)
	defer s.wp.Close()
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
		cancel()
	case <-ctx.Done():
	}

}

func (s ShortenerServer) ShutDown() {
	log.Print("Server closed")
}
