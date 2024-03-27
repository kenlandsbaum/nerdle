package server

import (
	"bytes"
	"essentials/nerdle/internal/dictionary"
	"essentials/nerdle/internal/game"
	"essentials/nerdle/internal/player"
	"essentials/nerdle/internal/service/id"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

type MockDict struct {
	mock.Mock
}

func (md *MockDict) GetWordApi(int) *dictionary.DefinitionResponse {
	args := md.Called()
	return args.Get(0).(*dictionary.DefinitionResponse)
}
func (md *MockDict) Orchestrate(int) string {
	return ""
}

func Test_handleStartGame(t *testing.T) {
	testIntFunc := func(i int) int { return i - 1 }
	os.Setenv("DICTIONARY_SIZE", "10")
	expected := dictionary.DefinitionResponse{
		Word: "test",
		Meanings: []dictionary.Meaning{{
			Definitions: []dictionary.Definition{
				{DefinitionString: "what you do"},
				{DefinitionString: "another sense"}},
		}},
	}

	mockDict := new(MockDict)
	mockDict.On("GetWordApi").Return(&expected)

	testPlayerId := id.GetUlid()
	testGameId := id.GetUlid()
	testRequestbody := []byte(
		fmt.Sprintf(`{"game_id":"%s","player_id":"%s"}`, testGameId, testPlayerId))

	testPlayer1 := player.ApiPlayer{Id: testPlayerId}
	testGame := game.ApiGame{GamePlayer: &testPlayer1}

	s := Server{games: map[ulid.ULID]*game.ApiGame{testGameId: &testGame}, intFunc: testIntFunc, dictionary: mockDict}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/test", bytes.NewReader(testRequestbody))

	err := s.handleStartGame(w, r)

	if err != nil {
		t.Errorf("expected nil error but got %s\n", err)
	}

	result := w.Result()

	if result.StatusCode != http.StatusOK {
		t.Errorf("got %d but expected %d\n", result.StatusCode, http.StatusOK)
	}
	body := result.Body
	defer body.Close()

	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		t.Errorf("expected nil error but got %s\n", err)
	}
	parsedResponse, err := unmarshalToType[dictionary.DefinitionResponse](bodyBytes)
	if err != nil {
		t.Errorf("expected nil error but got %s\n", err)
	}
	assert.EqualValues(t, &expected, parsedResponse)

	testBadRequestbody := []byte(
		fmt.Sprintf(`{"game_id":"%s","player_id":"%s"}`, testGameId, id.GetUlid()))

	w2 := httptest.NewRecorder()
	r2, _ := http.NewRequest(http.MethodPost, "/test", bytes.NewReader(testBadRequestbody))
	expectedError := s.handleStartGame(w2, r2)

	assert.EqualValues(t, "this player is not playing this game", expectedError.Error())
}
