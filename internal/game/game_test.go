package game

import (
	"bytes"
	"essentials/nerdle/internal/guess"
	"essentials/nerdle/internal/player"
	"strings"
	"testing"
)

func Test_getGuess(t *testing.T) {
	input := "testguess"
	testGame := Game{reader: strings.NewReader(input)}

	actual, _ := testGame.getGuess()

	if actual.String() != input {
		t.Fatalf("expected %s but got %s", input, actual)
	}

	testGame = Game{reader: strings.NewReader("\n")}
	_, err := testGame.getGuess()
	if err == nil {
		t.Fatal("expected an error but err was nil")
	}
}

func Test_handleAttempts(t *testing.T) {
	var testWriter bytes.Buffer
	testPlayer := player.Player{Name: "tst", Attempts: make([]guess.Guess, 0), Writer: &testWriter}
	testSuccessInput := "bad12 slice"
	testReader := strings.NewReader(testSuccessInput)
	testGame := Game{maxAttempts: 3, solution: "slice", reader: testReader, gamePlayer: &testPlayer}
	actual, _ := testGame.handleAttempts()
	if !actual {
		t.Fatalf("expected true for %s", testSuccessInput)
	}

	testPlayer = player.Player{Name: "tst", Attempts: make([]guess.Guess, 0), Writer: &testWriter}
	testFailureInput := "bad12 sclie slicey lices"
	testReader = strings.NewReader(testFailureInput)
	testGame = Game{maxAttempts: 3, solution: "slice", reader: testReader, gamePlayer: &testPlayer}
	actual, _ = testGame.handleAttempts()
	if actual {
		t.Fatalf("expected false for %s", testFailureInput)
	}
}

func Test_setExitMessage(t *testing.T) {
	testPlayer := player.Player{Name: "tst", Attempts: make([]guess.Guess, 0)}
	testGame := Game{maxAttempts: 3, solution: "slice", gamePlayer: &testPlayer}

	testGame.setExitMessage(true)
	if !strings.Contains(testGame.exitMessage, "winner") {
		t.Fatalf("unexpected string %s", testGame.exitMessage)
	}

	testGame.setExitMessage(false)
	if !strings.Contains(testGame.exitMessage, "lost") {
		t.Fatalf("unexpected string %s", testGame.exitMessage)
	}
}
