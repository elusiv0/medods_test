package token

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/elusiv0/medods_test/internal/model/api"
	"github.com/elusiv0/medods_test/internal/util/hash"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TokenManager struct {
	lifeTime time.Duration
	secret   string
}

type TokenInfo struct {
	UUID         string             `json:"uuid"`
	RefreshToken string             `json:"refresh_token"`
	RefreshId    primitive.ObjectID `json:"refresh_id`
}
type Claims struct {
	TokenInfo
	jwt.RegisteredClaims
}

func New(time time.Duration, secret string) *TokenManager {
	return &TokenManager{
		lifeTime: time,
		secret:   secret,
	}
}

func (tokenManager *TokenManager) CheckTokenMatch(tokenInfo TokenInfo, refreshToken string) (bool, error) {
	if err := hash.Compare(tokenInfo.RefreshToken, refreshToken); err != nil {
		return false, fmt.Errorf("TokenManager - CheckTokenMatch: %w", api.ErrTokenMismatch)
	}

	return true, nil
}

func (tokenManager *TokenManager) NewJWTToken(refreshToken string, uuid string, id primitive.ObjectID) (string, error) {
	claims := &Claims{
		TokenInfo{
			UUID:         uuid,
			RefreshToken: refreshToken,
			RefreshId:    id,
		},
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenManager.lifeTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tok, err := token.SignedString([]byte(tokenManager.secret))
	if err != nil {
		return "", fmt.Errorf("TokenManager - NewJWTToken - SignedString: %w", err)
	}

	return tok, nil
}

func (tokenManager *TokenManager) ValidateJWT(accessToken string) (TokenInfo, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(accessToken, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(tokenManager.secret), nil
	})

	switch {
	case token.Valid:
		return token.Claims.(*Claims).TokenInfo, nil
	case errors.Is(err, jwt.ErrTokenMalformed) || errors.Is(err, jwt.ErrTokenSignatureInvalid):
		return TokenInfo{}, fmt.Errorf("TokenManager - ValidateJwt: %w", api.ErrInvalidAccessToken)
	case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
		return token.Claims.(*Claims).TokenInfo, fmt.Errorf("TokenManager - ValidateJwt: %w", api.ErrAccessTokenExpired)
	default:
		return TokenInfo{}, fmt.Errorf("TokenManager - ValidateJwt: %w", err)
	}
}

func (tokenManager *TokenManager) NewRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("TokenManager - NewRefreshToken: %w", err)
	}
	return hex.EncodeToString(b), nil
}
