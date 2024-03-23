package player

import (
	"essentials/nerdle/internal/guess"
	"fmt"
	"io"
)

type Player struct {
	Attempts []guess.Guess
	Id       string
	Name     string
	Reader   io.Reader
	Writer   io.Writer
}

func InitPlayer(reader io.Reader, writer io.Writer) (*Player, error) {
	name, err := getName(reader)
	if err != nil {
		return nil, err
	}
	return &Player{Attempts: make([]guess.Guess, 0), Name: name, Reader: reader, Writer: writer}, nil
}

func (p *Player) PrintAtempts() {
	fmt.Fprintf(p.Writer, "player %s guesses so far:\n", p.Name)
	for i, attempt := range p.Attempts {
		fmt.Fprintf(p.Writer, "%d) %s\n", i+1, attempt.GetColoredString())
	}
}

func (p *Player) HasAttempted(guess string) bool {
	for _, attempt := range p.Attempts {
		if guess == attempt.String() {
			return true
		}
	}
	return false
}

func getName(reader io.Reader) (string, error) {
	var name string
	fmt.Println("Enter your name to start:")
	fmt.Scanf("%s", &name)
	for name == "" {
		fmt.Println("you gotta enter a non-empty name")
		_, err := fmt.Fscanf(reader, "%s", &name)
		if err != nil {
			return "", err
		}
	}
	return name, nil
}
