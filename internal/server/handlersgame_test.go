package server

import (
	"bytes"
	"essentials/nerdle/internal/game"
	"essentials/nerdle/internal/player"
	"essentials/nerdle/internal/service/id"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/oklog/ulid/v2"
)

func Test_handlePostGameSuccess(t *testing.T) {
	testId := id.GetUlid()
	testBody := fmt.Sprintf(`{"player_id":"%s"}`, testId)
	testRequestBody := []byte(testBody)
	testPlayer := player.ApiPlayer{Id: testId, Name: "ken"}
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/test", bytes.NewReader(testRequestBody))

	s := Server{
		mutex:   &sync.RWMutex{},
		players: map[ulid.ULID]*player.ApiPlayer{testId: &testPlayer},
		games:   make(map[ulid.ULID]*game.ApiGame, 0),
	}

	s.handlePostGame(w, r)

	result := w.Result()

	if result.StatusCode != http.StatusCreated {
		t.Errorf("expected %d but got %d\n", http.StatusCreated, result.StatusCode)
	}

	defer result.Body.Close()
	bts, _ := io.ReadAll(result.Body)

	if !strings.Contains(string(bts), `"game_id":"`) {
		t.Errorf("unexpected body %s\n", string(bts))
	}
}

func Test_handlePostGameError(t *testing.T) {
	testId := id.GetUlid()
	testBody := fmt.Sprintf(`{"player_id":"%s"}`, testId)
	testRequestBody := []byte(testBody)
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/test", bytes.NewReader(testRequestBody))

	s := Server{
		mutex:   &sync.RWMutex{},
		players: make(map[ulid.ULID]*player.ApiPlayer, 0),
		games:   make(map[ulid.ULID]*game.ApiGame, 0),
	}

	s.handlePostGame(w, r)

	result := w.Result()

	if result.StatusCode != http.StatusBadRequest {
		t.Errorf("got %d but expected %d\n", result.StatusCode, http.StatusBadRequest)
	}
}
