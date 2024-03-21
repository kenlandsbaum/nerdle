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
	Router chi.Router //http.Handler
}

func New(router chi.Router) *Server {
	return &Server{Router: router}
}

func (s *Server) applyMiddleware() {
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.RealIP)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)
	s.Router.Use(middleware.Timeout(2 * time.Second))
}

func (s *Server) routes() {
	s.Router.Get("/", s.handleHome)
}

func (s Server) Run() error {
	host := os.Getenv("API_HOST")
	s.applyMiddleware()
	s.routes()

	srv := &http.Server{Addr: host, Handler: s.Router}

	log.Info().Msg(fmt.Sprintf("starting server on %s", host))
	if err := srv.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
