package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dimk00z/go-shortener-praktikum/internal/handlers"
	"github.com/dimk00z/go-shortener-praktikum/internal/middleware/cookie"
	"github.com/dimk00z/go-shortener-praktikum/internal/middleware/decompressor"
	"github.com/dimk00z/go-shortener-praktikum/internal/storages/storageinterface"
	"github.com/dimk00z/go-shortener-praktikum/internal/util"
	"github.com/dimk00z/go-shortener-praktikum/internal/worker"
	"github.com/dimk00z/go-shortener-praktikum/pkg/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	// _ "github.com/swaggo/http-swagger/example/go-chi/docs"
)

type serverTLSConfig struct {
	port     string
	certFile string
	keyFile  string
}
type ShortenerServer struct {
	port      string
	Router    *chi.Mux
	wp        worker.IWorkerPool
	secretKey string
	l         *logger.Logger
	tlsConfig *serverTLSConfig
}

func NewServer(l *logger.Logger, port string, wp worker.IWorkerPool, secretKey string) *ShortenerServer {
	return &ShortenerServer{
		port:      port,
		Router:    chi.NewRouter(),
		wp:        wp,
		secretKey: secretKey,
		l:         l,
	}
}
func (s *ShortenerServer) mountMiddleware() {
	// Mount all Middleware here
	cookieHandler := cookie.CookieHandler{
		SecretKey: s.secretKey, L: s.l,
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
		handlers.SetLoger(s.l),
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

	apiRouter.Mount("/internal", chi.NewRouter().Route("/", func(r chi.Router) {
		r.Get("/stats", h.GetStats)
	}))

	s.Router.Mount("/api", apiRouter)

	s.Router.Mount("/ping",
		chi.NewRouter().Route("/", func(r chi.Router) {
			r.Get("/", h.PingDB)
		}))
	fileServer := http.FileServer(http.Dir("./docs/"))

	s.Router.HandleFunc("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
	})
	s.Router.Handle("/swagger/", http.StripPrefix("/swagger", util.ContentType(fileServer)))
	s.Router.Handle("/swagger/*", http.StripPrefix("/swagger", util.ContentType(fileServer)))
}

func (s *ShortenerServer) listenAndServe() (err error) {
	// HTTP case
	if s.tlsConfig == nil {
		s.l.Debug("Server started at " + s.port)
		err = http.ListenAndServe(s.port, s.Router)
		return err
	}

	// HTTPS case
	s.l.Debug("Server with TLS started at " + s.tlsConfig.port)
	if s.tlsConfig.port != ":443" {
		s.l.Warn("Default port for https is 443")
	}
	err = http.ListenAndServeTLS(
		s.tlsConfig.port,
		s.tlsConfig.certFile,
		s.tlsConfig.keyFile,
		s.Router)
	return err
}

func (s *ShortenerServer) RunServer(ctx context.Context, cancel context.CancelFunc, storage storageinterface.Storage) {
	interrupt := make(chan os.Signal, 1)
	defer s.wp.Close()
	shutdownSignals := []os.Signal{
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
	}
	signal.Notify(interrupt, shutdownSignals...)
	go func() {
		err := s.listenAndServe()
		if err != nil {
			s.l.Debug(err)
		}
		cancel()
	}()
	select {
	case killSignal := <-interrupt:
		s.l.Debug("Got ", killSignal)
		cancel()
	case <-ctx.Done():
	}

}

func (s ShortenerServer) ShutDown() {
	s.l.Debug("Server closed")
}
