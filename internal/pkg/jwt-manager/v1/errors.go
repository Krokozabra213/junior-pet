package jwtv1

import "errors"

// var (
// 	ErrValidExp      = errors.New("err token expired")
// 	ErrParseJWT      = errors.New("failed parse token")
// 	ErrValidToken    = errors.New("err valid token")
// 	ErrUserID        = errors.New("user_id not found in token")
// 	ErrExp           = errors.New("exp not found in token")
// 	ErrInvalidClaims = errors.New("invalid claims")
// )

var (
	ErrTokenExpired  = errors.New("token expired")
	ErrTokenParse    = errors.New("failed to parse token")
	ErrTokenInvalid  = errors.New("invalid token")
	ErrInvalidData   = errors.New("invalid token data")
	ErrSigningMethod = errors.New("unexpected signing method")
)
