package server

import (
	"essentials/nerdle/internal/game"
	"essentials/nerdle/internal/player"

	"github.com/oklog/ulid/v2"
)

type ApiPlayersIface interface {
	Add(*player.ApiPlayer)
	Get() map[ulid.ULID]*player.ApiPlayer
	GetById(ulid.ULID) (*player.ApiPlayer, bool)
}

type ApiGamesIface interface {
	Add(*game.ApiGame)
	Delete(ulid.ULID)
	GetById(ulid.ULID) (*game.ApiGame, bool)
}
