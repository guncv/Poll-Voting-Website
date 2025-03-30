package util

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	AccessTokenSecret  = []byte("your-access-token-secret")
	RefreshTokenSecret = []byte("your-refresh-token-secret")
)

// GenerateAccessToken generates a JWT access token with a short expiry.
func GenerateAccessToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(15 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(AccessTokenSecret)
}

// GenerateRefreshToken generates a JWT refresh token with a longer expiry.
func GenerateRefreshToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(RefreshTokenSecret)
}

// ValidateAccessToken validates the access token.
func ValidateAccessToken(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return AccessTokenSecret, nil
	})
}

// ValidateRefreshToken validates the refresh token.
func ValidateRefreshToken(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return RefreshTokenSecret, nil
	})
}
