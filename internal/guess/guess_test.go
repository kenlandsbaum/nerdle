package guess

import (
	"testing"
)

func Test_String(t *testing.T) {
	testGuess := New("test")
	result := testGuess.String()

	if result != "test" {
		t.Fatalf("expected %s but got %s", "test", result)
	}
}

func Test_Check(t *testing.T) {
	type tc struct {
		in       string
		expected bool
		statuses map[int]CharacterStatus
	}
	testCases := []tc{
		{in: "tst", expected: true},
		{
			in:       "tct",
			expected: false,
			statuses: map[int]CharacterStatus{0: CorrectPosition, 1: NotInWord, 2: CorrectPosition},
		},
		{
			in:       "abc",
			expected: false,
			statuses: map[int]CharacterStatus{0: NotInWord, 1: NotInWord, 2: NotInWord},
		},
		{
			in:       "act",
			expected: false,
			statuses: map[int]CharacterStatus{0: IncorrectPosition, 1: NotInWord, 2: CorrectPosition},
		},
		{
			in:       "acs",
			expected: false,
			statuses: map[int]CharacterStatus{0: NotInWord, 1: IncorrectPosition, 2: NotInWord},
		},
	}

	for _, ts := range testCases {
		testGuess := New("tst")
		actual, _ := testGuess.Check(ts.in)
		if actual != ts.expected {
			t.Fatalf("expected %v but got %v\n", ts.expected, actual)
		}
		if !actual {
			for k, enumVal := range ts.statuses {
				if testGuess[k].Status != enumVal {
					t.Fatalf(
						"test failed on guess character status assertion. wanted %d but got %d\n",
						enumVal,
						testGuess[k].Status,
					)
				}
			}
		}
	}
	testGuess := New("tst")
	_, err := testGuess.Check("lololol")
	if err == nil {
		t.Fatal("expected error to not be nil")
	}
}

func Test_GetColoredString(t *testing.T) {
	testGuess := New("test")

	testGuess[0].Status = CorrectPosition
	testGuess[1].Status = IncorrectPosition
	testGuess[2].Status = NotInWord
	testGuess[3].Status = CorrectPosition

	expected := "\033[32mt\033[0m\033[33me\033[0m\033[31ms\033[0m\033[32mt\033[0m"

	actual := testGuess.GetColoredString()
	if actual != expected {
		t.Fatalf("expected %s but got %s", expected, actual)
	}
}
