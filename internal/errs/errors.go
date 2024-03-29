package errs

import (
	"fmt"
)

func PanicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

type NotFoundError struct {
	Word string
}

// var NotFoundError = errors.New("not found")

func (n NotFoundError) Error() string {
	if n.Word != "" {
		return fmt.Sprintf("word '%s' not found", n.Word)
	}
	return "word not found"
}
