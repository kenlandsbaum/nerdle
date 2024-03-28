package game

import (
	"essentials/nerdle/internal/player"
	"essentials/nerdle/internal/service/id"
	"sync"
	"testing"

	"github.com/oklog/ulid/v2"
)

func getPlayersToAdd() map[ulid.ULID]*player.ApiPlayer {
	ps := make(map[ulid.ULID]*player.ApiPlayer, 0)
	for _, name := range []string{"ken", "jen", "sam", "jim", "jill", "jo", "pam", "dan"} {
		pId := id.GetUlid()
		ps[pId] = &player.ApiPlayer{Name: name, Id: pId}
	}
	return ps
}

func getGamesToAdd(players map[ulid.ULID]*player.ApiPlayer) map[ulid.ULID]*ApiGame {
	games := make(map[ulid.ULID]*ApiGame, 8)
	for _, p := range players {
		gId := id.GetUlid()
		game := NewApiGame(p, gId)
		games[game.Id] = game
	}
	return games
}
func TestAdd(t *testing.T) {
	testPlayers := getPlayersToAdd()
	testGamesToAdd := getGamesToAdd(testPlayers)

	testApiGames := NewApiGames(&sync.RWMutex{})

	wg := sync.WaitGroup{}
	for _, g := range testGamesToAdd {
		wg.Add(1)
		go func() {
			defer wg.Done()
			testApiGames.Add(g)
		}()
	}
	wg.Wait()

	if len(testApiGames.Games) != 8 {
		t.Errorf("got %d but expected %d\n", len(testApiGames.Games), 8)
	}

	for _, g := range testGamesToAdd {
		wg.Add(1)
		go func() {
			defer wg.Done()
			testApiGames.Delete(g.Id)
		}()
	}
	wg.Wait()
	if len(testApiGames.Games) != 0 {
		t.Errorf("got %d but expected %d\n", len(testApiGames.Games), 0)
	}
}
