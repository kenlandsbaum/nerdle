package server

import (
	"errors"
	"essentials/nerdle/internal/game"
	"essentials/nerdle/internal/player"
	"essentials/nerdle/internal/service/id"
	"fmt"
	"io"
	"net/http"

	"github.com/oklog/ulid/v2"
)

func (s *Server) handlePostGame(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	defer body.Close()
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		respondBadRequestErr(w, err)
		return
	}
	newGameRequest, err := unmarshalToType[NewGameRequest](bodyBytes)
	if err != nil {
		respondBadRequestErr(w, err)
		return
	}
	p, err := s.getPlayer(newGameRequest.PlayerID)
	if err != nil {
		respondBadRequestErr(w, err)
		return
	}
	gameId := id.GetUlid()
	s.games[gameId] = game.NewApiGame(p, gameId)
	respondCreated(w, []byte(fmt.Sprintf("game created with id %v", gameId)))
}

func (s *Server) getPlayer(id ulid.ULID) (*player.ApiPlayer, error) {
	var p *player.ApiPlayer
	s.mutex.RLock()
	p, ok := s.players[id]
	s.mutex.RUnlock()
	if !ok {
		return nil, errors.New("player not found")
	}
	return p, nil
}
