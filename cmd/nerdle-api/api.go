package main

import (
	"essentials/nerdle/internal/env"
	"essentials/nerdle/internal/server"
	"log"

	"github.com/go-chi/chi/v5"
)

func main() {
	env.Load(".env")
	srv := server.New(chi.NewRouter())
	if err := srv.Run(); err != nil {
		log.Fatalf("failed to initialize application %s", err)
	}
}
