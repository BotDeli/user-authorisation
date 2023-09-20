package hasher

import (
	"crypto/md5"
	"encoding/hex"
)

func Hashing(str string) string {
	hash := md5.Sum([]byte(str))
	hashedStr := hex.EncodeToString(hash[:])
	return hashedStr
}
