package main

import (
	"essentials/nerdle/internal/dictionary"
	"essentials/nerdle/internal/env"
	"essentials/nerdle/internal/images"
	"essentials/nerdle/internal/impl"
	"essentials/nerdle/internal/rest"
	"essentials/nerdle/internal/scoreboard"
	"essentials/nerdle/internal/server"
	"essentials/nerdle/internal/server/ctrl"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/go-chi/chi/v5"
)

func TestNetwork() {
	res, err := http.Get("https://randomuser.me/api")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	bts, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println("got response?", string(bts))
}

func main() {
	env.Load(".env")
	images.ProcessThumbnails(os.Getenv("IMAGES_FOLDER"))
	// TestNetwork()
	// App()
}

func App() {
	env.Load(".env")
	sb := scoreboard.New()
	scoreHandler := ctrl.NewScoreHandler(sb)
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
