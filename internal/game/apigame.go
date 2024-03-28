package game

import (
	"essentials/nerdle/internal/player"
	"sync"

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

type ApiGames struct {
	mutex *sync.RWMutex
	Games map[ulid.ULID]*ApiGame
}

func (a *ApiGames) Add(g *ApiGame) {
	a.mutex.Lock()
	a.Games[g.Id] = g
	a.mutex.Unlock()
}

func (a *ApiGames) GetById(id ulid.ULID) (*ApiGame, bool) {
	a.mutex.Lock()
	game, ok := a.Games[id]
	a.mutex.Unlock()
	return game, ok
}

func (a *ApiGames) Delete(id ulid.ULID) {
	a.mutex.Lock()
	delete(a.Games, id)
	a.mutex.Unlock()
}

func NewApiGames(mut *sync.RWMutex) *ApiGames {
	return &ApiGames{mut, make(map[ulid.ULID]*ApiGame, 0)}
}
