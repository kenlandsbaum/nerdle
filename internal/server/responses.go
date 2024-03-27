package server

import "net/http"

type PlayerCreatedResponse struct {
	PlayerID string `json:"player_id"`
}

type GameCreatedResponse struct {
	GameID string `json:"game_id"`
}

func respondOk(w http.ResponseWriter, body []byte) {
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func respondCreated(w http.ResponseWriter, body []byte) {
	w.WriteHeader(http.StatusCreated)
	w.Write(body)
}

func respondInternalErr(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}

func respondBadRequestErr(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(err.Error()))
}

func respondUnauthorizedErr(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(err.Error()))
}
