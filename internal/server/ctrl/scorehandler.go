package ctrl

import (
	"encoding/json"
	"essentials/nerdle/internal/scoreboard"
	"net/http"

	"github.com/rs/zerolog/log"
)

type ScoreHandler struct {
	ScoreBoard *scoreboard.ScoreboardV2
}

func (s *ScoreHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handleGetScores(w, r)
}

func (s *ScoreHandler) handleGetScores(w http.ResponseWriter, _ *http.Request) {
	playersScores := s.ScoreBoard.Players
	bts, err := json.Marshal(playersScores)
	if err != nil {
		log.Info().Msgf("error encountered %s", err)
		w.Write([]byte(err.Error()))
		return
	}
	log.Info().Msgf("bytes received %s", string(bts))
	w.Header().Add("Content-Type", "application/json")
	w.Write(bts)
}

func NewScoreHandler(sb *scoreboard.ScoreboardV2) *ScoreHandler {
	return &ScoreHandler{ScoreBoard: sb}
}
