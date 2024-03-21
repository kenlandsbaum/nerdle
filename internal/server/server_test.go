package server

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
)

func Test_serverRun(t *testing.T) {
	os.Setenv("API_HOST", "localhost:8888")
	r := chi.NewRouter()

	srv := New(r)

	go func() {
		time.Sleep(time.Second * 1)
		srv.Shutdown(context.Background())
	}()

	err := srv.Run()
	if err != nil && err != http.ErrServerClosed {
		t.Fatalf("expected nil error but got %s\n", err.Error())
	}
}
