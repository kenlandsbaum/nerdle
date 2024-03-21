package player

import (
	"bytes"
	"essentials/nerdle/internal/guess"
	"os"
	"strings"
	"testing"
)

func Test_getName(t *testing.T) {
	input := "ken"
	actual, _ := getName(strings.NewReader(input))

	if actual != input {
		t.Fatalf("expected %s but got %s", input, actual)
	}

	name, err := getName(strings.NewReader("\n"))
	if name != "" {
		t.Fatalf("expected name to be empty but got %s", name)
	}
	if err == nil {
		t.Fatal("expected an error but err was nil")
	}
}

func Test_Init(t *testing.T) {
	input := "ken"
	testPlayer, _ := InitPlayer(strings.NewReader(input), os.Stdout)
	if len(testPlayer.Attempts) != 0 {
		t.Fatal("expected empty initialization")
	}
	if testPlayer.Name != "ken" {
		t.Fatal("expected name initialized")
	}
}

func Test_HasAttempted(t *testing.T) {
	testPlayer, _ := InitPlayer(strings.NewReader("ken"), os.Stdout)

	testPlayer.Attempts = append(testPlayer.Attempts, guess.New("one"))
	testPlayer.Attempts = append(testPlayer.Attempts, guess.New("two"))
	testPlayer.Attempts = append(testPlayer.Attempts, guess.New("three"))

	if !testPlayer.HasAttempted("one") || !testPlayer.HasAttempted("two") || !testPlayer.HasAttempted("three") {
		t.Fatal("expected all true")
	}
	if testPlayer.HasAttempted("four") {
		t.Fatal("expected false")
	}
}

func Test_Print(t *testing.T) {
	reader := strings.NewReader("")
	var writer bytes.Buffer
	testPlayer := Player{Attempts: make([]guess.Guess, 0), Name: "ken", Reader: reader, Writer: &writer}

	testPlayer.Attempts = append(testPlayer.Attempts, guess.New("one"))
	testPlayer.Attempts = append(testPlayer.Attempts, guess.New("two"))

	testPlayer.PrintAtempts()

	actual := writer.String()
	expected := []string{
		"player ken guesses so far:",
		"1) \033[31mo\033[0m\033[31mn\033[0m\033[31me\033[0m",
		"2) \033[31mt\033[0m\033[31mw\033[0m\033[31mo\033[0m",
		"",
	}
	for i, line := range strings.Split(actual, "\n") {
		if line != expected[i] {
			t.Fatalf("expected %s but got %s", expected[i], actual)
		}
	}
}
