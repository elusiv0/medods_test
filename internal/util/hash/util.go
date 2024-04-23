package hash

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func CryptToken(token string) (string, error) {
	tokenBytes := []byte(token)

	hashedToken, err := bcrypt.GenerateFromPassword(tokenBytes, bcrypt.MinCost)

	if err != nil {
		return "", fmt.Errorf("Hash Util - CryptToken: %w", err)
	}

	return string(hashedToken), nil
}

func Compare(hashedRefresh, refresh string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedRefresh), []byte(refresh))
}
