package UUIDGenerator

import (
	"math/rand"
	"time"
)

func NewUUID() string {
	pattern := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	limit := len(pattern)
	var UUID string
	for i := 0; i < 16; i++ {
		UUID += string(pattern[random.Intn(limit)])
	}
	return UUID
}
