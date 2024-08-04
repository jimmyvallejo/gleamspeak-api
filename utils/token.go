package utils

import (
	"net/http"
	"strings"
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
