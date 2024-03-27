package game

import (
	"essentials/nerdle/internal/player"

	"github.com/oklog/ulid/v2"
)

type ApiGame struct {
	GamePlayer  *player.ApiPlayer
	maxAttempts int
	// solution    string
	Id ulid.ULID
}

func NewApiGame(p *player.ApiPlayer, id ulid.ULID) *ApiGame {
	return &ApiGame{GamePlayer: p, maxAttempts: 5, Id: id}
}
