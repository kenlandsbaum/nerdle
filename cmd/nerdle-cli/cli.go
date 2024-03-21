package main

import (
	"essentials/nerdle/internal/dictionary"
	"essentials/nerdle/internal/env"
	"essentials/nerdle/internal/errs"
	"essentials/nerdle/internal/game"
	"essentials/nerdle/internal/player"
	"essentials/nerdle/internal/rest"
	"essentials/nerdle/internal/scoreboard"
	"math/rand/v2"
	"net/http"
	"os"
	"strconv"
)

func main() {
	env.Load(".env")
	restClient := rest.New(http.DefaultClient)

	dict := &dictionary.Dictionary{
		DictionaryApi:    os.Getenv("DICTIONARY_API"),
		DictionarySource: os.Getenv("DICTIONARY_SOURCE"),
		FsClient:         Opener{},
		RestClient:       restClient,
		Writer:           os.Stdout,
	}

	RunGame(dict)
}

func RunGame(dict dictionary.DictionaryIface) {
	reader := os.Stdin
	writer := os.Stdout
	scoreboard := scoreboard.Scoreboard{Board: make(map[string]int, 1), Writer: writer}

	dictionarySize, err := strconv.Atoi(os.Getenv("DICTIONARY_SIZE"))
	errs.PanicIfErr(err)
	for {
		exitFunc := checkDoesWantToPlay(reader)
		if exitFunc != nil {
			scoreboard.PrintScore()
			exitFunc(0)
		}

		gamePlayer, err := player.InitPlayer(reader, writer)
		errs.PanicIfErr(err)
		randomNumber := rand.IntN(dictionarySize)
		solutionText := dict.Orchestrate(randomNumber)

		gm := game.New(gamePlayer, reader, solutionText)
		isSuccessful, err := gm.Play()
		errs.PanicIfErr(err)
		scoreboard.UpdateScore(gamePlayer.Name, isSuccessful)
	}
}
