package server

import (
	"essentials/nerdle/internal/game"
	"essentials/nerdle/internal/player"

	"github.com/oklog/ulid/v2"
)

type testApiGames struct {
	games map[ulid.ULID]*game.ApiGame
}

func (ta *testApiGames) Add(*game.ApiGame) {}
func (ta *testApiGames) Delete(id ulid.ULID) {
	delete(ta.games, id)
}
func (ta *testApiGames) GetById(id ulid.ULID) (*game.ApiGame, bool) {
	g, ok := ta.games[id]
	return g, ok
}

type testApiPlayers struct {
	players map[ulid.ULID]*player.ApiPlayer
}

func (tp *testApiPlayers) Add(p *player.ApiPlayer) {
	tp.players[p.Id] = p
}
func (tp *testApiPlayers) Get() map[ulid.ULID]*player.ApiPlayer {
	return tp.players
}
func (tp *testApiPlayers) GetById(id ulid.ULID) (*player.ApiPlayer, bool) {
	p, ok := tp.players[id]
	return p, ok
}
