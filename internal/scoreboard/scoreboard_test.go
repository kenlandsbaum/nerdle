package scoreboard

import (
	"bytes"
	"strings"
	"testing"
)

func Test_UpdateScore(t *testing.T) {
	testScoreboard := Scoreboard{Board: make(map[string]int, 0)}
	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			testScoreboard.UpdateScore("ken", true)
		} else {
			testScoreboard.UpdateScore("ken", false)
		}
	}
	actual := testScoreboard.Board["ken"]

	if actual != 5 {
		t.Fatalf("expected 5 but got %d", actual)
	}
}

func Test_PrintScore(t *testing.T) {
	var testWriter bytes.Buffer
	testScoreboard := Scoreboard{Board: make(map[string]int, 0), Writer: &testWriter}
	for i := 0; i < 3; i++ {
		testScoreboard.UpdateScore("ken", true)
	}
	for i := 0; i < 4; i++ {
		testScoreboard.UpdateScore("jen", true)
	}
	for i := 0; i < 1; i++ {
		testScoreboard.UpdateScore("jill", true)
	}

	testScoreboard.PrintScore()

	actual := testWriter.String()
	expected := []string{
		"",
		"Summary of games:",
		"",
		"ken: 3",
		"",
		"jen: 4",
		"",
		"jill: 1",
		"",
		"",
	}

	for _, line := range strings.Split(actual, "\n") {
		found := false
		for _, expectedLine := range expected {
			if expectedLine == line {
				found = true
			}
		}
		if !found {
			t.Fatal("expected line to be found")
		}
	}
}
