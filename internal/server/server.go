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
func (s *ShortenerServer) MountHandlers(host string, st storageinterface.Storage) {
	// Mount all Middleware here
	hadlers := []func(http.Handler) http.Handler{
		middleware.RequestID,
		middleware.Logger,
		middleware.Recoverer,
		decompressor.DecompressHandler,
		cookie.CookieHandler,
		middleware.Compress(5),
	}
	for _, handler := range hadlers {
		s.Router.Use(handler)
	}
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
		// Sprint 3 Increment 12
		r.Post("/batch", h.SaveBatch)
	})

	// Sprint 3
	userRouter := chi.NewRouter()
	userRouter.Route("/", func(r chi.Router) {
		userHandler := handlers.NewUserHandler(
			host,
			st,
			s.wp,
		)
		r.Get("/urls", userHandler.GetUserURLs)
		r.Delete("/urls", userHandler.DeleteUserURLs)
	})
	dbRouter := chi.NewRouter()
	dbRouter.Route("/", func(r chi.Router) {
		dbHandler := handlers.NewDBHandler(host,
			st)
		r.Get("/", dbHandler.PingDB)
	})
	apiRouter := chi.NewRouter()
	apiRouter.Mount("/shorten", shortenerRouter)
	apiRouter.Mount("/user", userRouter)

	s.Router.Mount("/api", apiRouter)
	s.Router.Mount("/ping", dbRouter)
}

func (s ShortenerServer) RunServer(ctx context.Context, cancel context.CancelFunc, storage storageinterface.Storage) {
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
		cancel()
	case <-ctx.Done():
	}
	log.Print("Server closed")
}
