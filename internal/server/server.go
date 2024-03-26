package server

import (
	"essentials/nerdle/internal/game"
	"essentials/nerdle/internal/player"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type Server struct {
	Router chi.Router
	*http.Server
	mutex   *sync.RWMutex
	players map[ulid.ULID]*player.ApiPlayer
	games   map[ulid.ULID]*game.ApiGame
}

func New(router chi.Router) *Server {
	srv := http.Server{
		Addr:         os.Getenv("API_HOST"),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	mutex := sync.RWMutex{}
	players := make(map[ulid.ULID]*player.ApiPlayer, 0)
	games := make(map[ulid.ULID]*game.ApiGame, 0)

	return &Server{router, &srv, &mutex, players, games}
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
	s.Router.Get("/", useJsonContent(s.handleHome))
	s.Router.Get("/player", useJsonContent(s.handleGetPlayers))
	s.Router.Post("/player", useJsonContent(s.handlePostPlayer))
	s.Router.Post("/game", useJsonContent(s.handlePostGame))
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
