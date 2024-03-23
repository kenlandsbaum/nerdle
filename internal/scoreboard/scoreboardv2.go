package scoreboard

import "essentials/nerdle/internal/player"

type ScoreboardV2 map[string]*player.ScoredPlayer

func (s ScoreboardV2) AddPlayerScore(p *player.Player) {
	if _, ok := s[p.Name]; !ok {
		s[p.Name] = &player.ScoredPlayer{Player: p}
	}
	s[p.Name].Score += 1
}
