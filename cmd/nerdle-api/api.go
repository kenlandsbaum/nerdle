package main

import (
	"essentials/nerdle/internal/env"
	"essentials/nerdle/internal/server"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/go-chi/chi/v5"
)

func main() {
	env.Load(".env")
	srv := server.New(chi.NewRouter())
	if err := srv.Run(); err != nil {
		log.Fatal().Msg(fmt.Sprintf("failed to run application %s", err))
	}
}
