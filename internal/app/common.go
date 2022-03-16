package app

import (
	"crypto/sha512"
	"encoding/hex"
)

func Hash(password string) string {
	s := sha512.New()
	s.Write([]byte(password))
	return hex.EncodeToString(s.Sum(nil))
}
