package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_handleHome(t *testing.T) {
	expectedResponseBody := `{"message":"welcome to the app"}`
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/", nil)

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
