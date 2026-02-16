package jwtv1

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenData struct {
	UserID   int64
	Username string
	Email    string
}

func (d TokenData) valid() error {
	if d.UserID == 0 {
		return errors.New("user_id is required")
	}
	if d.Username == "" {
		return errors.New("username is required")
	}
	if d.Email == "" {
		return errors.New("email is required")
	}
	return nil
}

// при парсе возвращаются
type AccessClaims struct {
	UserID   int64     `json:"user_id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Exp      time.Time `json:"exp"`
}

type RefreshClaims struct {
	JWTID  string    `json:"jti"`
	UserID int64     `json:"user_id"`
	Exp    time.Time `json:"exp"`
}

// при генерации передаются в claims
type accessJWTClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

type refreshJWTClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}
