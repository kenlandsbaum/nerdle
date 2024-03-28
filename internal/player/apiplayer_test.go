package player

import (
	"essentials/nerdle/internal/service/id"
	"sync"
	"testing"

	"github.com/oklog/ulid/v2"
)

func getPlayersToAdd() map[ulid.ULID]*ApiPlayer {
	ps := make(map[ulid.ULID]*ApiPlayer, 8)
	for _, name := range []string{"ken", "jen", "sam", "jim", "jill", "jo", "pam", "dan"} {
		pId := id.GetUlid()
		ps[pId] = &ApiPlayer{Name: name, Id: pId}
	}
	return ps
}

func TestAdd(t *testing.T) {
	testPlayers := NewApiPlayers(&sync.RWMutex{})
	playersToAdd := getPlayersToAdd()

	wg := sync.WaitGroup{}
	for _, p := range playersToAdd {
		wg.Add(1)
		go func() {
			defer wg.Done()
			testPlayers.Add(p)
		}()
	}
	wg.Wait()

	actualPlayersAdded := testPlayers.Get()
	for id, player := range actualPlayersAdded {
		if player.Name != playersToAdd[id].Name {
			t.Errorf("got %s but expected %s\n", player.Name, playersToAdd[id].Name)
		}
	}
}
