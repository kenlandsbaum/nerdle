package main

import (
	"io/fs"
	"os"
)

type Opener struct{}

func (o Opener) Open(s string) (fs.File, error) {
	return os.Open(s)
}
