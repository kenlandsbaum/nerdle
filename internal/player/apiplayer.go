package player

import (
	"essentials/nerdle/internal/guess"

	"github.com/oklog/ulid/v2"
)

type ApiPlayer struct {
	Attempts []guess.Guess
	Name     string
	Id       ulid.ULID
}
