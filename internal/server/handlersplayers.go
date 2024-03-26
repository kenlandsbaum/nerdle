package server

import (
	"errors"
	"essentials/nerdle/internal/player"
	"essentials/nerdle/internal/service/id"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

func (s *Server) handlePostPlayer(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	defer body.Close()
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		respondBadRequestErr(w, err)
		return
	}
	newPlayerRequest, err := unmarshalToType[NewPlayerRequest](bodyBytes)
	if err != nil {
		respondBadRequestErr(w, err)
		return
	}
	if newPlayerRequest.Name == "" {
		respondBadRequestErr(w, errors.New("player missing required property 'name'"))
		return
	}
	newPlayer := player.ApiPlayer{Name: newPlayerRequest.Name, Id: id.GetUlid()}
	s.players[newPlayer.Id] = &newPlayer
	log.Info().Msgf("new player created: %v\n", newPlayer)
	respondCreated(w, []byte(fmt.Sprintf("created player %s", newPlayer.Id.String())))
}

func (s *Server) handleGetPlayers(w http.ResponseWriter, _ *http.Request) {
	playersSlice := make([]*player.ApiPlayer, 0)
	for _, p := range s.players {
		playersSlice = append(playersSlice, p)
	}
	playersBytes, err := marshalToJson(playersSlice)
	if err != nil {
		respondInternalErr(w, err)
		return
	}
	respondOk(w, playersBytes)
}
