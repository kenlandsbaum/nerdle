package server

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	good        = "good"
	bad         = "bad"
	expectedMsg = "this player is not playing this game"
	otherMsg    = "other"
)

func testHandlerFunc(w http.ResponseWriter, r *http.Request) error {
	u := r.URL.Path
	if u == good {
		w.Write([]byte(good))
		return nil
	}
	if u == bad {
		return errors.New(expectedMsg)
	}
	return errors.New(otherMsg)
}

func Test_handleError(t *testing.T) {
	t.Run("should have no error to handle", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, good, nil)

		handleError(testHandlerFunc)(w, r)
		result := w.Result()
		defer result.Body.Close()
		bts, _ := io.ReadAll(result.Body)
		assert.Equal(t, good, string(bts))
	})

	t.Run("should write expected message", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, bad, nil)

		handleError(testHandlerFunc)(w, r)
		result := w.Result()
		defer result.Body.Close()
		bts, _ := io.ReadAll(result.Body)
		assert.Equal(t, expectedMsg, string(bts))
	})

	t.Run("should write other message", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "lol", nil)

		handleError(testHandlerFunc)(w, r)
		result := w.Result()
		defer result.Body.Close()
		bts, _ := io.ReadAll(result.Body)
		assert.Equal(t, otherMsg, string(bts))
	})
}
