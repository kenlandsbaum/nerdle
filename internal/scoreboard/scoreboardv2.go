package scoreboard

import (
	"essentials/nerdle/internal/player"

	"github.com/rs/zerolog/log"
)

type ScoreboardV2 struct {
	Players      map[string]*player.ScoredPlayer
	ScoreChannel chan *player.ApiPlayer
}

func (s *ScoreboardV2) AddPlayerScore(p *player.ApiPlayer) {
	if _, ok := s.Players[p.Name]; !ok {
		s.Players[p.Name] = &player.ScoredPlayer{Player: p}
	}
	s.Players[p.Name].Score += 1
}

func (s *ScoreboardV2) ListenForPlayer() {
	for p := range s.ScoreChannel {
		log.Info().Msgf("received player %s", p.Id)
		s.AddPlayerScore(p)
	}
}

func New() *ScoreboardV2 {
	return &ScoreboardV2{Players: make(map[string]*player.ScoredPlayer, 0), ScoreChannel: make(chan *player.ApiPlayer)}
}
