package jwtv1

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func generateTokenID() (string, error) {
	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generate random bytes: %w", err)
	}

	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(bytes), nil
}

func parseToken[T jwt.Claims](tokenString string, claims T, keyFunc jwt.Keyfunc) (T, error) {
	token, err := jwt.ParseWithClaims(tokenString, claims, keyFunc)
	if err != nil {
		// Отличаем "просрочен" от "сломан"
		if errors.Is(err, jwt.ErrTokenExpired) {
			return claims, ErrTokenExpired
		}
		return claims, fmt.Errorf("%w: %w", ErrTokenParse, err)
	}

	if !token.Valid {
		return claims, ErrTokenInvalid
	}

	return claims, nil
}
