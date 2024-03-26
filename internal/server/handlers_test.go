package server

import (
	"bytes"
	"essentials/nerdle/internal/player"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/oklog/ulid/v2"
)

func Test_handleHome(t *testing.T) {
	expectedResponseBody := `{"message":"welcome to the app"}`
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	s := Server{}
	s.handleHome(w, r)

	result := w.Result()
	if result.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 but got %d\n", result.StatusCode)
	}
	defer result.Body.Close()
	bts, _ := io.ReadAll(result.Body)
	if string(bts) != expectedResponseBody {
		t.Fatalf("expected %s but got %s\n", expectedResponseBody, string(bts))
	}
}

func Test_handlePostPlayerSuccess(t *testing.T) {
	expectedResponseBody := `created player`
	testRequestBody := []byte(`{"name":"ken"}`)
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/", bytes.NewReader(testRequestBody))

	s := Server{players: make(map[ulid.ULID]*player.ApiPlayer, 0)}
	s.handlePostPlayer(w, r)
	result := w.Result()

	if result.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 but got %d\n", result.StatusCode)
	}
	defer result.Body.Close()
	bts, _ := io.ReadAll(result.Body)

	if !strings.Contains(string(bts), expectedResponseBody) {
		t.Fatalf("expected %s but got %s\n", expectedResponseBody, string(bts))
	}
}

func Test_handlePostPlayerFailure(t *testing.T) {
	expectedResponseBody := `missing required property 'name'`
	testRequestBody := []byte(`{"notname":"ken"}`)
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/", bytes.NewReader(testRequestBody))

	s := Server{players: make(map[ulid.ULID]*player.ApiPlayer, 0)}
	s.handlePostPlayer(w, r)
	result := w.Result()

	if result.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 but got %d\n", result.StatusCode)
	}
	defer result.Body.Close()
	bts, _ := io.ReadAll(result.Body)

	if !strings.Contains(string(bts), expectedResponseBody) {
		t.Fatalf("expected %s but got %s\n", expectedResponseBody, string(bts))
	}
}
