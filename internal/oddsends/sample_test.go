package oddsends

import (
	"strings"
	"testing"
)

type Doer interface {
	Do() int
}

type A struct{}
type B struct{}

func (d A) Do() int { return 42 }

type Tester interface {
	Pass(code string) bool
	Fail(code string) bool
}

type C struct{}
type D struct{}

func (c *C) Pass(code string) bool { return len(code) > 0 }
func (c *C) Fail(code string) bool { return strings.Contains(code, "fail") }

func (d *D) Pass(code string) bool { return len(code) > 0 }

func TestCheckDoesImplement(t *testing.T) {
	a := A{}
	b := B{}
	c := C{}
	d := D{}
	type tst struct {
		in       any
		expected bool
	}
	doerCases := map[string]tst{
		"should return true for a Doer": {
			in:       a,
			expected: true,
		},
		"should return false for a non-Doer": {
			in:       b,
			expected: false,
		},
		"should return true for pointer receiver": {
			in:       &a,
			expected: true,
		},
	}
	for name, testCase := range doerCases {
		t.Run(name, func(t *testing.T) {
			actual := CheckDoesImplement[Doer](testCase.in)
			if actual != testCase.expected {
				t.Errorf("got %v but wanted %v\n", actual, testCase.expected)
			}
		})
	}
	testerCases := map[string]tst{
		"should return true for a Tester": {
			in:       &c,
			expected: true,
		},
		"should return false for a non-Tester": {
			in:       &d,
			expected: false,
		},
		"should return false for non-pointer receiver": {
			in:       c,
			expected: false,
		},
	}
	for name, testCase := range testerCases {
		t.Run(name, func(t *testing.T) {
			actual := CheckDoesImplement[Tester](testCase.in)
			if actual != testCase.expected {
				t.Errorf("got %v but wanted %v\n", actual, testCase.expected)
			}
		})
	}
}
