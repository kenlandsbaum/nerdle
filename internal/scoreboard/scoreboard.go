package scoreboard

import (
	"fmt"
	"io"
)

type Scoreboard struct {
	Board  map[string]int
	Writer io.Writer
}

func (s Scoreboard) UpdateScore(name string, isSuccessful bool) {
	scoreForName, ok := s.Board[name]
	if ok {
		s.Board[name] = scoreForName + zeroOrOne(isSuccessful)
		return
	}
	s.Board[name] = zeroOrOne(isSuccessful)
}

func zeroOrOne(isSuccessful bool) int {
	if isSuccessful {
		return 1
	}
	return 0
}

func (s Scoreboard) PrintScore() {
	msg := "\nSummary of games:\n"
	for name, score := range s.Board {
		msg += fmt.Sprintf("\n%s: %d\n", name, score)
	}
	fmt.Fprintln(s.Writer, msg)
}
