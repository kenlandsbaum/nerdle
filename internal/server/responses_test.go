package server

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type SomeBody struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func Test_respondCreated(t *testing.T) {
	body := SomeBody{Field: "f1", Message: "m1"}
	expectedResponseBody := `{"field":"f1","message":"m1"}`

	w := httptest.NewRecorder()
	bodyBytes, _ := marshalToJson(body)

	respondCreated(w, bodyBytes)

	result := w.Result()
	if result.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 but got %d\n", result.StatusCode)
	}
	defer result.Body.Close()
	bts, err := io.ReadAll(result.Body)
	if err != nil {
		t.Fatalf("expected nil error but got %s", err)
	}
	if string(bts) != expectedResponseBody {
		t.Fatalf("expected %s but got %s", expectedResponseBody, string(bts))
	}
}

func Test_respondInternalErr(t *testing.T) {
	internalErr := errors.New("bruh, you suck")

	w := httptest.NewRecorder()

	respondInternalErr(w, internalErr)
	result := w.Result()
	if result.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500 but got %d", result.StatusCode)
	}
	defer result.Body.Close()
	bts, _ := io.ReadAll(result.Body)
	if string(bts) != "bruh, you suck" {
		t.Fatalf("unexpected body %s", string(bts))
	}
}

func Test_respondBadRequestErr(t *testing.T) {
	internalErr := errors.New("bruh, bad request")

	w := httptest.NewRecorder()

	respondBadRequestErr(w, internalErr)
	result := w.Result()
	if result.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 but got %d", result.StatusCode)
	}
	defer result.Body.Close()
	bts, _ := io.ReadAll(result.Body)
	if string(bts) != "bruh, bad request" {
		t.Fatalf("unexpected body %s", string(bts))
	}
}
