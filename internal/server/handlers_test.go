package server

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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

func Test_handlePostPlayer(t *testing.T) {
	expectedResponseBody := `created player`
	testRequestBody := []byte(`{"name":"ken"}`)
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/", bytes.NewReader(testRequestBody))

	s := Server{}
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
