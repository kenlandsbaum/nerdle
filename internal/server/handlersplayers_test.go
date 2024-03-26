package server

import (
	"bytes"
	"essentials/nerdle/internal/player"
	"essentials/nerdle/internal/service/id"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/oklog/ulid/v2"
)

func Test_handlePostPlayerSuccess(t *testing.T) {
	expectedResponseBody := `created player`
	testRequestBody := []byte(`{"name":"ken"}`)
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/", bytes.NewReader(testRequestBody))

	s := Server{players: make(map[ulid.ULID]*player.ApiPlayer, 0), mutex: &sync.RWMutex{}}
	s.handlePostPlayer(w, r)
	result := w.Result()

	if result.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 but got %d\n", result.StatusCode)
	}
	defer result.Body.Close()
	bts, _ := io.ReadAll(result.Body)

	if !strings.Contains(string(bts), expectedResponseBody) {
		t.Fatalf("expected %s but got %s\n", expectedResponseBody, string(bts))
	}
}

func Test_handlePostPlayerFailure(t *testing.T) {
	expectedResponseBody := `missing required property 'name'`
	testRequestBody := []byte(`{"notname":"ken"}`)
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/", bytes.NewReader(testRequestBody))

	s := Server{players: make(map[ulid.ULID]*player.ApiPlayer, 0)}
	s.handlePostPlayer(w, r)
	result := w.Result()

	if result.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 but got %d\n", result.StatusCode)
	}
	defer result.Body.Close()
	bts, _ := io.ReadAll(result.Body)

	if !strings.Contains(string(bts), expectedResponseBody) {
		t.Fatalf("expected %s but got %s\n", expectedResponseBody, string(bts))
	}
}

func Test_handleGetPlayersSuccess(t *testing.T) {
	type testCase struct {
		in       int
		expected string
	}
	testCases := []testCase{{in: 0, expected: "player1"}, {in: 1, expected: "player2"}}

	s := Server{players: make(map[ulid.ULID]*player.ApiPlayer, 0)}
	u1 := id.GetUlid()
	u2 := id.GetUlid()

	s.players[u1] = &player.ApiPlayer{Name: "player1", Id: u1}
	s.players[u2] = &player.ApiPlayer{Name: "player2", Id: u2}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	s.handleGetPlayers(w, r)

	result := w.Result()

	if result.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 but got %d\n", result.StatusCode)
	}
	defer result.Body.Close()
	bts, _ := io.ReadAll(result.Body)

	players, err := unmarshalToType[[]*player.ApiPlayer](bts)
	if err != nil {
		t.Errorf("expected nil error but got %s\n", err)
	}
	if len(*players) != 2 {
		t.Errorf("expected 2 but got %d\n", len(*players))
	}

	for i, player := range *players {
		actual := player.Name
		expected := testCases[i].expected
		if expected != actual {
			t.Errorf("expected %s but got %s\n", expected, actual)
		}
	}
}
