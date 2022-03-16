package app

import (
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"strconv"

	"github.com/theplant/luhn"
)

func Hash(password string) string {
	s := sha512.New()
	s.Write([]byte(password))
	return hex.EncodeToString(s.Sum(nil))
}

func CheckOrderNumber(number string) error {
	orderInt, err := strconv.Atoi(number)
	if err != nil {
		return err
	}

	if !luhn.Valid(orderInt) {
		return errors.New("invalid order number")
	}

	return nil
}
