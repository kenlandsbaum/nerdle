package server

import (
	"essentials/nerdle/internal/player"
	"essentials/nerdle/internal/service/id"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

func (s Server) handleHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	respondOk(w, []byte(`{"message":"welcome to the app"}`))
}

func (s Server) handlePostPlayer(w http.ResponseWriter, r *http.Request) {
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
	newPlayer := player.ApiPlayer{Name: newPlayerRequest.Name, Id: id.GetUlid()}
	// do something with player
	log.Info().Msgf("new player created: %v\n", newPlayer)
	respondCreated(w, []byte(fmt.Sprintf("created player %s", newPlayer.Id.String())))
}
