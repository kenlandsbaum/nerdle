package game

import (
	"essentials/nerdle/internal/player"

	"github.com/oklog/ulid/v2"
)

type ApiGame struct {
	gamePlayer  *player.ApiPlayer
	maxAttempts int
	// solution    string
	Id ulid.ULID
}

func NewApiGame(p *player.ApiPlayer, id ulid.ULID) *ApiGame {
	return &ApiGame{gamePlayer: p, maxAttempts: 5, Id: id}
}
