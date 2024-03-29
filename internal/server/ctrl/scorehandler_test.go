package ctrl

import (
	"encoding/json"
	"essentials/nerdle/internal/player"
	"essentials/nerdle/internal/scoreboard"
	"essentials/nerdle/internal/service/id"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestScoreHandler(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/test", nil)

	testId := id.GetUlid()
	testPlayer := player.ApiPlayer{Name: "ken", Id: testId, Attempts: []string{"one"}}
	testScoreboard := scoreboard.New()

	for i := 0; i < 5; i++ {
		testScoreboard.AddPlayerScore(&testPlayer)
	}

	scoreHandler := New(testScoreboard)
	scoreHandler.ServeHTTP(w, r)

	result := w.Result()
	defer result.Body.Close()
	bts, err := io.ReadAll(io.NopCloser(result.Body))
	if err != nil {
		t.Errorf("expected non-nil error but got %s\n", err.Error())
	}
	if bts == nil {
		t.Error("expected non-nil body")
	}
	var actual map[string]*player.ScoredPlayer
	parseErr := json.Unmarshal(bts, &actual)
	if parseErr != nil {
		t.Errorf("expected nil error but got %s\n", parseErr)
	}
	actualName := actual["ken"].Player.Name
	if actualName != "ken" {
		t.Errorf("expected 'ken' but got %s\n", actualName)
	}
}
