package token

import (
	"errors"
)

var (
	ErrRefreshTokenNotRegistered = errors.New("refresh token not found in registered tokens")
)
