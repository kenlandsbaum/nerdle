package guess

import (
	"errors"
	"fmt"
	"strings"
)

type Guess []Letter

func New(guessWord string) Guess {
	letters := make([]Letter, 0, len(guessWord))
	for _, l := range guessWord {
		letters = append(letters, Letter{Status: NotInWord, Character: CharacterType(l)})
	}
	return letters
}

func (g Guess) String() string {
	word := make([]CharacterType, 0, len(g))
	for _, l := range g {
		word = append(word, l.Character)
	}
	return string(word)
}

func (g Guess) Check(solution string) (bool, error) {
	if len(g) != len(solution) {
		return false, errors.New("wrong size guess")
	}

	if g.String() == solution {
		return true, nil
	}

	for i, l := range g {
		if rune(l.Character) == rune(solution[i]) {
			g[i].Status = CorrectPosition
		} else if strings.ContainsRune(solution, rune(l.Character)) {
			g[i].Status = IncorrectPosition
		} else {
			l.Status = NotInWord
		}
	}
	return false, nil
}

func (g Guess) GetColoredString() string {
	word := ""
	for _, l := range g {
		var highlighted string
		switch l.Status {
		case CorrectPosition:
			highlighted = l.Character.GetColorCode(colorGreen)
		case IncorrectPosition:
			highlighted = l.Character.GetColorCode(colorYellow)
		case NotInWord:
			highlighted = l.Character.GetColorCode(colorRed)
		default:
			highlighted = fmt.Sprintf("%c", l.Character)
		}
		word = word + highlighted
	}
	return word
}
