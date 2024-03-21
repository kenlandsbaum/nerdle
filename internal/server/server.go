package server

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

type Server struct {
	Router chi.Router
	*http.Server
}

func New(router chi.Router) *Server {
	srv := http.Server{
		Addr:         os.Getenv("API_HOST"),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	return &Server{router, &srv}
}

func (s *Server) applyMiddleware() {
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.RealIP)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)
	s.Router.Use(authenticate)
	s.Router.Use(middleware.Timeout(3 * time.Second))
}

func (s *Server) routes() {
	s.Router.Get("/", s.handleHome)
}

func (s *Server) Run() error {
	s.applyMiddleware()
	s.routes()

	log.Info().Msg(fmt.Sprintf("starting server on %s", s.Addr))

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
