package game

import (
	"essentials/nerdle/internal/guess"
	"essentials/nerdle/internal/player"
	"fmt"
	"io"
	"log"
)

const (
	winningMessage = "\033[32mA winner is you!\033[0m"
	losingMessage  = "\033[31mYou lost the game.\033[0m"
)

type Game struct {
	exitMessage string
	gamePlayer  *player.Player
	greeting    string
	maxAttempts int
	reader      io.Reader
	solution    string
}

func New(gamePlayer *player.Player, reader io.Reader, solution string) *Game {
	return &Game{
		gamePlayer:  gamePlayer,
		greeting:    "Welcome to Nerdle",
		maxAttempts: 5,
		reader:      reader,
		solution:    solution,
	}
}

func (g *Game) Play() (bool, error) {
	defer g.finish()
	g.greet()
	isSuccessful, err := g.handleAttempts()
	if err != nil {
		return false, err
	}
	g.setExitMessage(isSuccessful)
	return isSuccessful, nil
}

func (g *Game) greet() {
	fmt.Printf("%s, %s!\n", g.greeting, g.gamePlayer.Name)
}

func (g *Game) handleAttempts() (bool, error) {
	for len(g.gamePlayer.Attempts) < g.maxAttempts {
		fmt.Printf("You have %d more tries to guess the word\n", g.maxAttempts-len(g.gamePlayer.Attempts))
		guess, writeErr := g.getGuess()

		if writeErr != nil {
			log.Print("error writing?", writeErr)
			return false, writeErr
		}
		if g.gamePlayer.HasAttempted(guess.String()) {
			fmt.Printf("You have already guessed %s. Try again\n", guess)
			continue
		}
		isCorrect, err := guess.Check(g.solution)
		if err != nil {
			fmt.Printf("error in solution: %s\n", err.Error())
			continue
		}
		g.gamePlayer.Attempts = append(g.gamePlayer.Attempts, guess)
		if isCorrect {
			fmt.Printf("You guessed '%s' correctly!\n", guess)
			return true, nil
		} else {
			fmt.Printf("Your guess of %s, is incorrect\n", guess)
			g.gamePlayer.PrintAtempts()
			continue
		}
	}
	return false, nil
}

func (g *Game) getGuess() (guess.Guess, error) {
	var guessWord string
	if _, err := fmt.Fscanf(g.reader, "%s", &guessWord); err != nil {
		return guess.New(""), err
	}
	return guess.New(guessWord), nil
}

func (g *Game) setExitMessage(didSucceed bool) {
	if didSucceed {
		g.exitMessage = winningMessage
		return
	}
	g.exitMessage = losingMessage + fmt.Sprintf("\nYour word was '%s'.", g.solution)
}

func (g *Game) finish() {
	fmt.Println(g.exitMessage, "\nThanks for playing. Goodbye!")
}
