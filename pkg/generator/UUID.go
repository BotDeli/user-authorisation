package generator

import (
	"math/rand"
	"time"
)

const patternDigitsLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
const lengthUUIDDigitsLetters = 16

func NewUUIDDigitsLetters() string {
	return generateUUID(patternDigitsLetters, lengthUUIDDigitsLetters)
}

func generateUUID(pattern string, lengthUUID int) string {
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

const patternDigits = "0123456789"
const lengthUUIDDigits = 7

func NewUUIDDigits() string {
	return generateUUID(patternDigits, lengthUUIDDigits)
}
