package guess

import "fmt"

const (
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
	colorReset  = "\033[0m"
)

type CharacterType rune
type CharacterStatus int

const (
	NotInWord CharacterStatus = iota
	CorrectPosition
	IncorrectPosition
)

type Letter struct {
	Character CharacterType
	Status    CharacterStatus
}

func (c CharacterType) GetColorCode(color string) string {
	return fmt.Sprintf("%s%c%s", color, c, colorReset)
}
