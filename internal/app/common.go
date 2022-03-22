package app

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"github.com/slavkluev/gophermart/internal/app/middleware"
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

func LoginFromContext(ctx context.Context) (string, bool) {
	u, ok := ctx.Value(middleware.ContextLoginKey).(string)
	return u, ok
}
