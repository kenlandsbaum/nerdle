package server

import "github.com/oklog/ulid/v2"

type NewPlayerRequest struct {
	Name string `json:"name"`
}

type NewGameRequest struct {
	PlayerID ulid.ULID `json:"player_id"`
}

type StartGameRequest struct {
	PlayerID ulid.ULID `json:"player_id"`
	GameId   ulid.ULID `json:"game_id"`
}

type GuessRequest struct {
	GameId ulid.ULID `json:"game_id"`
	Guess  string    `json:"guess"`
}
