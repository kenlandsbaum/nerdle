package player

import (
	"essentials/nerdle/internal/player"
	"essentials/nerdle/internal/service/id"
)

func CreatePlayer(name string) *player.Player {
	return &player.Player{Id: id.GetUlid(), Name: name}
}
