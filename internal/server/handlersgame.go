package server

import (
	"errors"
	"essentials/nerdle/internal/errs"
	"essentials/nerdle/internal/game"
	"essentials/nerdle/internal/player"
	"essentials/nerdle/internal/service/id"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/oklog/ulid/v2"
)

func (s *Server) handlePostGame(w http.ResponseWriter, r *http.Request) {
	newGameRequest, err := decodeRequestBody[NewGameRequest](r.Body)
	if err != nil {
		respondBadRequestErr(w, err)
		return
	}
	p, err := s.getPlayer(newGameRequest.PlayerID)
	if err != nil {
		respondBadRequestErr(w, err)
		return
	}
	gameId := id.GetUlid()
	s.games.Add(game.NewApiGame(p, gameId))
	respondCreated(w, mustMarshal(GameCreatedResponse{GameID: gameId.String()}))
}

func (s *Server) getPlayer(id ulid.ULID) (*player.ApiPlayer, error) {
	var p *player.ApiPlayer
	p, ok := s.players.GetById(id)
	if !ok {
		return nil, errors.New("player not found")
	}
	return p, nil
}

type HandlerFuncErr func(w http.ResponseWriter, r *http.Request) error

func (s *Server) handleStartGame(w http.ResponseWriter, r *http.Request) error {
	startGameRequest, err := decodeRequestBody[StartGameRequest](r.Body)
	if err != nil {
		return err
	}
	game, ok := s.games.GetById(startGameRequest.GameId) //[startGameRequest.GameId]
	if !ok {
		return errors.New("you must create a new game")
	}
	if game.GamePlayer.Id != startGameRequest.PlayerID {
		return errors.New("this player is not playing this game")
	}
	definitionResponse := s.dictionary.GetWordApi(s.getRandomInt())
	game.Solution = definitionResponse.Word
	definitionResponse.Word = mask(definitionResponse.Word)
	respondOk(w, mustMarshal(definitionResponse))
	return nil
}

func (s *Server) getRandomInt() int {
	dictionarySize, err := strconv.Atoi(os.Getenv("DICTIONARY_SIZE"))
	errs.PanicIfErr(err)
	return s.intFunc(dictionarySize)
}

func mask(s string) string {
	var masked string
	for i, c := range s {
		if i == 0 {
			masked += string(c)
		} else {
			masked += "*"
		}
	}
	return masked
}

func (s *Server) handleGuess(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	defer body.Close()
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		respondInternalErr(w, err)
		return
	}
	guessRequest, err := unmarshalToType[GuessRequest](bodyBytes)
	game, ok := s.games.GetById(guessRequest.GameId)
	if !ok {
		respondBadRequestErr(w, err)
		return
	}
	game.GamePlayer.Attempts = append(game.GamePlayer.Attempts, guessRequest.Guess)
	if game.Solution != guessRequest.Guess {
		remaining := game.MaxAttempts - len(game.GamePlayer.Attempts)
		respondOk(w, []byte(fmt.Sprintf(`{"status":"you suck","remainingAttempts":%d}`, remaining)))
		return
	}
	s.scoreChannel <- game.GamePlayer
	respondOk(w, []byte(`{"status":"a winner is you!"}`))
	s.games.Delete(guessRequest.GameId)
}
