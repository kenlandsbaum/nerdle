package main

import (
	"strings"
	"testing"
)

func Test_doYouWantToPlay(t *testing.T) {
	type testCase struct {
		arg      string
		expected bool
	}

	tests := []testCase{
		{arg: "y", expected: true},
		{arg: "yes", expected: true},
		{arg: "Y", expected: true},
		{arg: "Yes", expected: true},
		{arg: "n", expected: false},
		{arg: "no", expected: false},
		{arg: "N", expected: false},
		{arg: "No", expected: false},
	}

	for _, tc := range tests {
		testReader := strings.NewReader(tc.arg)

		actual := checkDoesWantToPlay(testReader)
		if tc.expected && actual != nil {
			t.Fatal("expected nil func returned")
		}
		if !tc.expected && actual == nil {
			t.Fatal("expected non-nil func returned")
		}
	}

	invalidTextTestReader := strings.NewReader("oopsy\nNo")

	actualRecovered := checkDoesWantToPlay(invalidTextTestReader)
	if actualRecovered == nil {
		t.Fatal("expected false for input 'oops\nNo'")
	}
}
