package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET")) // It's good practice to get this from an environment variable

func init() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// For development, you can have a default, but in production, enforce it.
		os.Setenv("JWT_SECRET", "933665fbcda138607c8cefcc9faed3e0f885be91c4040a3211207dffc37e7fdf")
		jwtSecret = []byte("933665fbcda138607c8cefcc9faed3e0f885be91c4040a3211207dffc37e7fdf")
		println("Warning: JWT_SECRET environment variable not set, using a default. This is not secure for production.")
	}
}

// GenerateJWT generates a new JWT for the given username.
func GenerateJWT(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Minute * 1).Unix(), // Token expires in 1 minute
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// VerifyJWT verifies the JWT and returns the username if valid.
func VerifyJWT(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username, ok := claims["username"].(string)
		if !ok {
			return "", jwt.ErrTokenInvalidClaims
		}
		return username, nil
	}

	return "", jwt.ErrInvalidKey
}
