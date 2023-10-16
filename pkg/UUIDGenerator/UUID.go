package UUIDGenerator

import (
	"math/rand"
	"time"
)

const pattern = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
const lengthUUID = 16

func NewUUID() string {
	random := newRandom()
	limit := len(pattern)

	UUID := make([]byte, lengthUUID)

	for i := 0; i < lengthUUID; i++ {
		UUID[i] = pattern[random.Intn(limit)]
	}
	return string(UUID)
}

func newRandom() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}
