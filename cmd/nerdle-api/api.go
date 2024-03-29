package main

import (
	"essentials/nerdle/internal/dictionary"
	"essentials/nerdle/internal/env"
	"essentials/nerdle/internal/impl"
	"essentials/nerdle/internal/rest"
	"essentials/nerdle/internal/scoreboard"
	"essentials/nerdle/internal/server"
	"essentials/nerdle/internal/server/ctrl"
	"fmt"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/go-chi/chi/v5"
)

func main() {
	env.Load(".env")
	sb := scoreboard.New()
	scoreHandler := ctrl.New(sb)
	go sb.ListenForPlayer()

	restClient := rest.New(http.DefaultClient)
	dict := &dictionary.Dictionary{
		DictionaryApi:    os.Getenv("DICTIONARY_API"),
		DictionarySource: os.Getenv("DICTIONARY_SOURCE"),
		RestClient:       restClient,
		FsClient:         impl.Opener{},
	}
	srv := server.New(chi.NewRouter(), dict, sb.ScoreChannel, scoreHandler)
	if err := srv.Run(); err != nil {
		log.Fatal().Msg(fmt.Sprintf("failed to run application %s", err))
	}
}
