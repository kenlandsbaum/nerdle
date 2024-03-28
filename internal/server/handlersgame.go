package server

import (
	"errors"
	"essentials/nerdle/internal/errs"
	"essentials/nerdle/internal/game"
	"essentials/nerdle/internal/player"
	"essentials/nerdle/internal/service/id"
	"fmt"
	"io"
	"log"
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
	s.games[gameId] = game.NewApiGame(p, gameId)
	respondCreated(w, mustMarshal(GameCreatedResponse{GameID: gameId.String()}))
}

func (s *Server) getPlayer(id ulid.ULID) (*player.ApiPlayer, error) {
	var p *player.ApiPlayer
	s.mutex.RLock()
	p, ok := s.players[id]
	s.mutex.RUnlock()
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
	game := s.games[startGameRequest.GameId]
	if game.GamePlayer.Id != startGameRequest.PlayerID {
		return errors.New("this player is not playing this game")
	}
	definitionResponse := s.dictionary.GetWordApi(s.getRandomInt())
	game.Solution = definitionResponse.Word
	definitionResponse.Word = mask(definitionResponse.Word)
	respondOk(w, mustMarshal(definitionResponse))
	return nil
}

func handleError(fn HandlerFuncErr) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			switch err.Error() {
			case "this player is not playing this game":
				respondBadRequestErr(w, err)
			default:
				respondInternalErr(w, err)
			}
		}
	}
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
	log.Println("body received:", string(bodyBytes))
	guessRequest, err := unmarshalToType[GuessRequest](bodyBytes)
	game, ok := s.games[guessRequest.GameId]
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

	respondOk(w, []byte(`{"status":"a winner is you!"}`))
	s.mutex.Lock()
	delete(s.games, guessRequest.GameId)
	s.mutex.Unlock()
}
