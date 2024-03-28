package types

import "essentials/nerdle/internal/dictionary"

type DictionaryIface interface {
	Orchestrate(int) string
	GetWordApi(int) *dictionary.DefinitionResponse
}
