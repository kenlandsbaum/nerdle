package id

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

func GetUlid() ulid.ULID {
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	ms := ulid.Timestamp(time.Now())
	u, err := ulid.New(ms, entropy)
	if err != nil {
		return ulid.Make()
	}
	return u
}
