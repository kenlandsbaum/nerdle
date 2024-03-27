package server

import (
	"context"
	"essentials/nerdle/internal/dictionary"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
)

type testDict struct{}

func (ts testDict) Orchestrate(_ int) string {
	return ""
}
func (ts testDict) GetWordApi(int) *dictionary.DefinitionResponse {
	return nil
}

func Test_serverRun(t *testing.T) {
	os.Setenv("API_HOST", "localhost:8888")
	r := chi.NewRouter()

	srv := New(r, testDict{})

	go func() {
		time.Sleep(time.Second * 1)
		srv.Shutdown(context.Background())
	}()

	err := srv.Run()
	if err != nil && err != http.ErrServerClosed {
		t.Fatalf("expected nil error but got %s\n", err.Error())
	}
}
