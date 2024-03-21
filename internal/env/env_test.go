package env

import (
	"os"
	"testing"
)

func Test_Load(t *testing.T) {
	Load("./mocks/.mockenv")
	expectedDictionaryApi := "https://api.dictionaryapi.dev/api/v2/entries/en/"
	dictionaryApi := os.Getenv("DICTIONARY_API")
	if dictionaryApi != expectedDictionaryApi {
		t.Fatalf("expected %s but got %s", expectedDictionaryApi, dictionaryApi)
	}
}
