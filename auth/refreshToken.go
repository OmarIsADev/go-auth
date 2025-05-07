package auth

import (
	"errors"
	"math/rand"

	"github.com/omarisadev/go-auth/database"
)

// returns a refresh token and saves it to memory DB
func GenerateRefreshToken(username string) (string, error) {
	if username == "" {
		return "", errors.New("username cannot be empty")
	}

	bytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 32)
	for i := range b {
		b[i] = bytes[rand.Intn(len(bytes))]
	}

	database.InMemoryDB().SaveRefreshToken(username, string(b))
	return string(b), nil
}
