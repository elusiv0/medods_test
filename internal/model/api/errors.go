package api

import (
	"errors"
)

var (
	ErrNoUUID             = errors.New("no uuid in query parameters")
	ErrNoAccessTokenFound = errors.New("no authorization token found in request headers")
	ErrInvalidAccessToken = errors.New("invalid token")
	ErrAccessTokenExpired = errors.New("token is expired")
	ErrTokenMismatch      = errors.New("tokens pair mismatch: invalid refresh token for access token")
	ErrBadRefreshRequest  = errors.New("refresh and access token are required")
)
