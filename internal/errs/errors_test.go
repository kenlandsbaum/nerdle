package errs

import (
	"errors"
	"testing"
)

var (
	badness = "badness"
)

func TestPanicIfErr(t *testing.T) {

	defer func() {
		err := recover()
		if err == nil {
			t.Error("expected non-nil error")
		}
		actual := err.(error).Error()
		if actual != badness {
			t.Errorf("got %s but expected %s\n", actual, badness)
		}
	}()

	PanicIfErr(errors.New(badness))
}
