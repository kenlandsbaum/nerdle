package server

import "github.com/oklog/ulid/v2"

type NewPlayerRequest struct {
	Name string `json:"name"`
}

type NewGameRequest struct {
	PlayerID ulid.ULID `json:"player_id"`
}
