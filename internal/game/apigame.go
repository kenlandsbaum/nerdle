package game

import (
	"essentials/nerdle/internal/player"

	"github.com/oklog/ulid/v2"
)

type ApiGame struct {
	GamePlayer  *player.ApiPlayer
	MaxAttempts int
	Solution    string
	Id          ulid.ULID
}

func NewApiGame(p *player.ApiPlayer, id ulid.ULID) *ApiGame {
	return &ApiGame{GamePlayer: p, MaxAttempts: 5, Id: id}
}
