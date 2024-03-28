package player

import (
	"sync"

	"github.com/oklog/ulid/v2"
)

type ApiPlayer struct {
	Attempts []string
	Name     string
	Id       ulid.ULID
}

type ApiPlayers struct {
	mutex   *sync.RWMutex
	Players map[ulid.ULID]*ApiPlayer
}

func (a *ApiPlayers) Add(p *ApiPlayer) {
	a.mutex.Lock()
	a.Players[p.Id] = p
	a.mutex.Unlock()
}

func (a *ApiPlayers) Get() map[ulid.ULID]*ApiPlayer {
	return a.Players
}

func NewApiPlayers(mut *sync.RWMutex) *ApiPlayers {
	return &ApiPlayers{mutex: mut, Players: make(map[ulid.ULID]*ApiPlayer, 0)}
}
