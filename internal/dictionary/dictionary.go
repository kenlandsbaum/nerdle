package dictionary

import (
	"bufio"
	"encoding/json"
	"errors"
	"essentials/nerdle/internal/rest"
	"fmt"
	"io"
	"io/fs"
	"log"
)

type DictionaryWord struct {
	Word       string
	Definition string
}

type DefinitionResponse struct {
	Word     string    `json:"word"`
	Meanings []Meaning `json:"meanings"`
}

type Meaning struct {
	PartOfSpeech string       `json:"partOfSpeech"`
	Definitions  []Definition `json:"definitions"`
}

type Definition struct {
	DefinitionString string `json:"definition"`
}

type DictionaryIface interface {
	Lookup(string) (*DictionaryWord, error)
	Orchestrate(randomInt int) string
}

type Dictionary struct {
	DictionaryApi    string
	DictionarySource string
	FsClient         fs.FS
	RestClient       rest.RestClientIface
	Writer           io.Writer
}

func (d Dictionary) Lookup(word string) (*DictionaryWord, error) {
	return nil, nil
}

func (d Dictionary) getWord(randomInt int) (*DictionaryWord, error) {
	file, err := d.FsClient.Open(d.DictionarySource)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	word, err := d.scanSource(file, randomInt)
	if err != nil {
		return nil, err
	}
	return word, nil
}

func (d Dictionary) scanSource(reader io.Reader, randomInt int) (*DictionaryWord, error) {
	scanner := bufio.NewScanner(reader)
	i := 0
	for scanner.Scan() {
		i += 1
		if i == randomInt {
			current := scanner.Text()
			return &DictionaryWord{Word: current}, nil
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("Error reading from dictionary file: %s", err)
		return nil, err
	}
	return nil, errors.New("word not found")
}

func (d Dictionary) getDefinition(word *DictionaryWord) (*DefinitionResponse, error) {
	wordUrl := fmt.Sprintf("%s%s", d.DictionaryApi, word.Word)
	apiResponse, err := d.RestClient.Get(wordUrl).Do()
	if err != nil {
		return nil, err
	}
	definitionResponse, err := parseToDefintionResponse(apiResponse.Body)
	if err != nil {
		return nil, err
	}
	return definitionResponse, nil
}

func parseToDefintionResponse(bts []byte) (*DefinitionResponse, error) {
	var response []DefinitionResponse
	if err := json.Unmarshal(bts, &response); err != nil {
		return nil, err
	}
	if len(response) == 1 {
		return &response[0], nil
	}
	return nil, errors.New("empty or pathological response")
}

func (d Dictionary) Orchestrate(randomInt int) string {
	word, _ := d.getWord(randomInt)
	def, err := d.getDefinition(word)
	for err != nil {
		randomInt += 1
		word, _ = d.getWord(randomInt)
		def, err = d.getDefinition(word)
	}
	if def == nil {
		log.Println("definition came back nil. You are on your own.")
		log.Printf("Your word is %d characters and begins with %c.", len(word.Word), word.Word[0])
		return word.Word
	}
	if def.Meanings != nil {
		d.showHints(def)
	} else {
		fmt.Printf(
			"\nNo definitions came back from the service. Your word is %d characters and begins with %c. Lots of luck.\n",
			len(def.Word),
			def.Word[0],
		)
	}
	return def.Word
}

func (d Dictionary) showHints(definitionResponse *DefinitionResponse) {
	meaning := definitionResponse.Meanings[0]
	fmt.Fprint(d.Writer, "\nHints:")
	fmt.Fprintf(
		d.Writer,
		"\nYour word is %d letters long and beings with '%c'.\n\n",
		len(definitionResponse.Word), definitionResponse.Word[0])
	for i, m := range meaning.Definitions {
		fmt.Fprintf(d.Writer, "hint %d: %v\n", i+1, m.DefinitionString)
	}
	fmt.Fprintln(d.Writer, "_______________________")
}
