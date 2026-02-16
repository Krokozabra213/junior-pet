package jwtv1

import (
	"crypto/rsa"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// Validator can only parse/validate tokens (no generation).
// Use in services that receive tokens but don't issue them.
type Validator struct {
	publicKey *rsa.PublicKey
}

func NewValidator(publicKey *rsa.PublicKey) (*Validator, error) {
	if publicKey == nil {
		return nil, errors.New("public key is required")
	}
	return &Validator{publicKey: publicKey}, nil
}

// ValidateAccess parses and validates an access token.
func (v *Validator) ValidateAccess(tokenString string) (*AccessClaims, error) {
	claims := &accessJWTClaims{}

	claims, err := parseToken(tokenString, claims, v.keyFunc)
	if err != nil {
		return nil, err
	}

	return &AccessClaims{
		UserID:   claims.UserID,
		Username: claims.Username,
		Email:    claims.Email,
		Exp:      claims.ExpiresAt.Time,
	}, nil
}

// ValidateRefresh parses and validates a refresh token.
func (v *Validator) ValidateRefresh(tokenString string) (*RefreshClaims, error) {
	claims := &refreshJWTClaims{}

	claims, err := parseToken(tokenString, claims, v.keyFunc)
	if err != nil {
		return nil, err
	}

	return &RefreshClaims{
		JWTID:  claims.ID,
		UserID: claims.UserID,
		Exp:    claims.ExpiresAt.Time,
	}, nil
}

func (v *Validator) keyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, fmt.Errorf("%w: %v", ErrSigningMethod, token.Header["alg"])
	}
	return v.publicKey, nil
}
