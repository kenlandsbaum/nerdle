package scoreboard

import (
	"essentials/nerdle/internal/player"
	"testing"
)

func TestAddPlayerScore(t *testing.T) {
	p1 := player.ApiPlayer{Name: "ken"}
	p2 := player.ApiPlayer{Name: "sam"}

	s := ScoreboardV2{Players: make(map[string]*player.ScoredPlayer, 2)}

	for i := 0; i < 10; i++ {
		s.AddPlayerScore(&p1)
	}
	for j := 0; j < 7; j++ {
		s.AddPlayerScore(&p2)
	}

	actual1, ok1 := s.Players["ken"]
	if !ok1 {
		t.Errorf("expected ok true")
	}
	if actual1.Score != 10 {
		t.Errorf("got %d expected 10\n", actual1.Score)
	}
	actual2, ok2 := s.Players["sam"]
	if !ok2 {
		t.Errorf("expected ok true")
	}
	if actual2.Score != 7 {
		t.Errorf("got %d expected 10\n", actual2.Score)
	}
}
