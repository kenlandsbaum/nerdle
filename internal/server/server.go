package server

import (
	"essentials/nerdle/internal/fileserver"
	"essentials/nerdle/internal/game"
	"essentials/nerdle/internal/player"
	"essentials/nerdle/internal/types"
	"fmt"
	"math/rand/v2"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

type Server struct {
	Router chi.Router
	*http.Server
	players      ApiPlayersIface
	games        ApiGamesIface
	dictionary   types.DictionaryIface
	intFunc      func(int) int
	scoreChannel chan *player.ApiPlayer
}

func New(router chi.Router, dict types.DictionaryIface, scoreChannel chan *player.ApiPlayer) *Server {
	srv := http.Server{
		Addr:         os.Getenv("API_HOST"),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	players := player.NewApiPlayers(&sync.RWMutex{})
	games := game.NewApiGames(&sync.RWMutex{})
	intFunc := rand.IntN

	return &Server{router, &srv, players, games, dict, intFunc, scoreChannel}
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
	s.Router.Get("/", useTextContent(s.handleHome))
	s.Router.Get("/player", useJsonContent(s.handleGetPlayers))
	s.Router.Post("/player", useJsonContent(s.handlePostPlayer))
	s.Router.Post("/game", useJsonContent(s.handlePostGame))
	s.Router.Post("/start", useJsonContent(handleError(s.handleStartGame)))
	s.Router.Post("/guess", useJsonContent(s.handleGuess))
}

func (s *Server) ui() {
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "ui"))
	fileserver.FileServer(s.Router, "/ui", filesDir)
}

func (s *Server) Run() error {
	s.applyMiddleware()
	s.routes()
	s.ui()

	log.Info().Msg(fmt.Sprintf("starting server on %s", s.Addr))

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
