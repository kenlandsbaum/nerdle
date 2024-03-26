package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func Test_authenticateUnauthorized(t *testing.T) {
	testRouter := chi.NewRouter()
	fn1 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	testRouter.Use(authenticate)
	testRouter.Get("/tst", fn1)

	r, _ := http.NewRequest(http.MethodGet, "/tst", nil)
	w := httptest.NewRecorder()

	testRouter.ServeHTTP(w, r)
	result := w.Result()

	if result.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 but got %d\n", result.StatusCode)
	}
}

func Test_authenticateOK(t *testing.T) {
	testRouter := chi.NewRouter()
	fn1 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	testRouter.Use(authenticate)
	testRouter.Get("/tst", fn1)

	r, _ := http.NewRequest(http.MethodGet, "/tst", nil)
	r.Header.Add("authorization", "somevalue")
	w := httptest.NewRecorder()

	testRouter.ServeHTTP(w, r)
	result := w.Result()

	if result.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 but got %d\n", result.StatusCode)
	}
	defer result.Body.Close()
	bts, err := io.ReadAll(result.Body)
	if err != nil {
		t.Fatalf("expected nil error but got %s\n", err)
	}
	if string(bts) != "success" {
		t.Fatalf("expected 'success' but got %s\n", string(bts))
	}
}
