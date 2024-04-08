package server

import (
	"essentials/nerdle/internal/fileserver"
	"essentials/nerdle/internal/game"
	"essentials/nerdle/internal/player"
	"essentials/nerdle/internal/pro"
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
	"google.golang.org/protobuf/proto"
)

type Server struct {
	Router chi.Router
	*http.Server
	players      ApiPlayersIface
	games        ApiGamesIface
	dictionary   types.DictionaryIface
	intFunc      func(int) int
	scoreChannel chan *player.ApiPlayer
	scoreHandler http.Handler
}

func New(
	router chi.Router,
	dict types.DictionaryIface,
	scoreChannel chan *player.ApiPlayer,
	scoreHandler http.Handler) *Server {
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

	return &Server{router, &srv, players, games, dict, intFunc, scoreChannel, scoreHandler}
}

func (s *Server) applyMiddleware() {
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.RealIP)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)
	// s.Router.Use(authenticate)
	s.Router.Use(middleware.Timeout(3 * time.Second))
}

func (s *Server) handlePersonTest(w http.ResponseWriter, _ *http.Request) {
	address := pro.Address{Street: "123 elm st", State: "TX"}
	person := &pro.Person{FirstName: "jen", LastName: "lee", Email: "jen@mail.com", Address: &address}
	bts, err := proto.Marshal(person)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("protobuf stuff fell on its face"))
		return
	}
	w.WriteHeader(200)
	w.Write(bts)
}

func (s *Server) routes() {
	s.Router.Get("/", useTextContent(s.handleHome))
	s.Router.Get("/player", useJsonContent(s.handleGetPlayers))
	s.Router.Post("/player", useJsonContent(s.handlePostPlayer))
	s.Router.Post("/game", useJsonContent(s.handlePostGame))
	s.Router.Post("/start", useJsonContent(handleError(s.handleStartGame)))
	s.Router.Post("/guess", useJsonContent(s.handleGuess))
	s.Router.Mount("/score", s.scoreHandler)
	s.Router.Get("/pro", s.handlePersonTest)
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
