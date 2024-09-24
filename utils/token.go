package utils

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenError struct {
	Message string
	Code    int
}

func (e *TokenError) Error() string {
	return e.Message
}

var (
	ErrEmptyHeader   = &TokenError{Message: "Authorization header is missing", Code: http.StatusUnauthorized}
	ErrInvalidPrefix = &TokenError{Message: "Invalid authorization prefix", Code: http.StatusUnauthorized}
	ErrEmptyToken    = &TokenError{Message: "Token is empty", Code: http.StatusUnauthorized}
)

func CreateToken(id uuid.UUID, jwtSecret string, expiresInSeconds int) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "gleamspeak",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Duration(expiresInSeconds) * time.Second)),
		Subject:   id.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to create token: %w", err)
	}
	return signedToken, nil
}


func ExtractToken(r *http.Request, prefix string) (string, error) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		return "", ErrEmptyHeader
	}

	if !strings.HasPrefix(authHeader, prefix) {
		return "", ErrInvalidPrefix
	}

	token := strings.TrimPrefix(authHeader, prefix)

	if token == "" {
		return "", ErrEmptyToken
	}
	return token, nil
}

func SetTokenCookie(w http.ResponseWriter, name, value string, maxAgeSeconds int) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  time.Now().UTC().Add(time.Duration(maxAgeSeconds) * time.Second),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
}

func ValidateToken(tokenString, jwtSecret string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
