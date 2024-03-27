package dictionary

import (
	"bytes"
	"errors"
	"essentials/nerdle/internal/rest"
	"io/fs"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	testGoodResponseBody = []byte(
		`[{"word":"test","meanings":[{"definitions":[{"definition":"what you do!"}]}]}]`)
)

type MockRest struct {
	mock.Mock
}

func (m *MockRest) Get(url string) rest.RestClientIface {
	return m
}
func (m *MockRest) Do() (*rest.RestResponse, error) {
	args := m.Called()
	return args.Get(0).(*rest.RestResponse), args.Error(1)
}

type MockFS struct {
	mock.Mock
}

func (m *MockFS) Open(string) (fs.File, error) {
	args := m.Called()
	return args.Get(0).(fs.File), args.Error(1)
}

func Test_GetDefinition(t *testing.T) {
	testRc := new(MockRest)
	testFs := new(MockFS)
	testRc.On("Do").Return(&rest.RestResponse{Body: testGoodResponseBody}, nil)

	expected := DefinitionResponse{
		Word: "test",
		Meanings: []Meaning{
			{Definitions: []Definition{{DefinitionString: "what you do!"}}},
		},
	}
	testDict := Dictionary{FsClient: testFs, RestClient: testRc}

	actualDef, err := testDict.getDefinition(&DictionaryWord{Word: "test"})
	if err != nil {
		t.Errorf("got %s but expected nil error\n", err.Error())
	}
	if actualDef.Word != expected.Word {
		t.Errorf("expected %s but got %s\n", "test", actualDef.Word)
	}

	for i, meaning := range actualDef.Meanings {
		assert.Equal(t, meaning, expected.Meanings[i])
	}
}

func Test_scanSource(t *testing.T) {
	testRc := new(MockRest)
	testFs := new(MockFS)
	testReader := strings.NewReader("one\ntwo\n")
	testFs.On("Open").Return(testReader, nil)

	testDict := Dictionary{FsClient: testFs, RestClient: testRc}
	actualWord, err := testDict.scanSource(testReader, 2)
	if err != nil {
		t.Errorf("expected nil error but got %s\n", err)
	}
	if actualWord.Word != "two" {
		t.Errorf("expected %s but got %s\n", "two", actualWord.Word)
	}

	_, actualErr := testDict.scanSource(testReader, 4)
	if actualErr == nil {
		t.Error("expected error to be non-nil")
	}
	if actualErr.Error() != "word not found" {
		t.Errorf("expected %s but got %s\n", "word not found", actualErr.Error())
	}
}

type testReader struct {
	strings.Reader
}

func (tst *testReader) Open(s string) (fs.File, error) {
	if s == "good" {
		r := strings.NewReader("one\ntwo\n")
		return &testReader{*r}, nil
	}
	if s == "bad" {
		return nil, errors.New("bruh")
	}
	return nil, nil
}

func (tst *testReader) Close() error {
	return nil
}

func (tst *testReader) Stat() (fs.FileInfo, error) {
	return nil, nil
}

func Test_getWord(t *testing.T) {
	randomInt := 1
	testRc := new(MockRest)
	testFs := new(testReader)
	testDict := Dictionary{FsClient: testFs, RestClient: testRc, DictionarySource: "good"}

	actual, _ := testDict.getWord(randomInt)

	if actual.Word != "one" {
		t.Errorf("expected 'one' but got %s", actual.Word)
	}

	_, actualErr := testDict.getWord(4)
	if actualErr == nil {
		t.Errorf("expected non nil error")
	}

	testBadDict := Dictionary{FsClient: testFs, RestClient: testRc, DictionarySource: "bad"}

	_, fsClientErr := testBadDict.getWord(1)
	if fsClientErr == nil {
		t.Errorf("expected non nil error")
	}
}

func TestOrchestrate(t *testing.T) {
	randomInt := 1
	var testWriter bytes.Buffer
	testRc := new(MockRest)
	testFs := new(testReader)
	testRc.On("Do").Return(&rest.RestResponse{Body: testGoodResponseBody}, nil)
	testDict := Dictionary{FsClient: testFs, RestClient: testRc, DictionarySource: "good", Writer: &testWriter}

	solution := testDict.Orchestrate(randomInt)
	if solution != "test" {
		t.Errorf("got '%s' but expected 'test'", solution)
	}
}

func Test_showHints(t *testing.T) {
	testRc := new(MockRest)
	testFs := new(MockFS)
	var testWriter bytes.Buffer
	testDict := Dictionary{FsClient: testFs, RestClient: testRc, Writer: &testWriter}

	definitionResponse := DefinitionResponse{
		Word: "test",
		Meanings: []Meaning{
			{Definitions: []Definition{
				{DefinitionString: "what you do!"},
				{DefinitionString: "what else you going to do?"},
			}},
		},
	}
	expectedOutput := []string{
		"",
		"Hints:",
		"Your word is 4 letters long and beings with 't'.",
		"",
		"hint 1: what you do!",
		"hint 2: what else you going to do?",
		"_______________________",
		"",
	}

	testDict.showHints(&definitionResponse)

	values := testWriter.String()
	lines := strings.Split(values, "\n")
	for i, line := range lines {
		if line != expectedOutput[i] {
			t.Errorf("expected %s but got %s", expectedOutput[i], line)
		}
	}
}

func Test_parseToDefintionResponse(t *testing.T) {
	testBytes := []byte("[]")
	_, err := parseToDefintionResponse(testBytes)
	if err == nil {
		t.Error("expected non nil error")
	}

	testBytes = []byte("<>")
	_, parsingErr := parseToDefintionResponse(testBytes)
	if parsingErr == nil {
		t.Error("expected non nil error")
	}
}
