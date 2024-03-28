package player

import (
	"github.com/oklog/ulid/v2"
)

type ApiPlayer struct {
	Attempts []string
	Name     string
	Id       ulid.ULID
}
