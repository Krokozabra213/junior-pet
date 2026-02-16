package jwtv1

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	defaultAccessTTL  = 15 * time.Minute
	defaultRefreshTTL = 15 * 24 * time.Hour
)

type Manager struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
	accessTTL  time.Duration
	refreshTTL time.Duration
}

type Option func(*Manager)

func WithAccessTTL(ttl time.Duration) Option {
	return func(m *Manager) {
		m.accessTTL = ttl
	}
}

func WithRefreshTTL(ttl time.Duration) Option {
	return func(m *Manager) {
		m.refreshTTL = ttl
	}
}

func New(private *rsa.PrivateKey, public *rsa.PublicKey, opts ...Option) (*Manager, error) {
	if private == nil {
		return nil, errors.New("private key is required")
	}
	if public == nil {
		return nil, errors.New("public key is required")
	}

	m := &Manager{
		publicKey:  public,
		privateKey: private,
		accessTTL:  defaultAccessTTL,
		refreshTTL: defaultRefreshTTL,
	}

	for _, opt := range opts {
		opt(m)
	}

	return m, nil
}

func (m *Manager) GenerateTokens(data TokenData) (access, refresh string, err error) {
	if err := data.valid(); err != nil {
		return "", "", fmt.Errorf("%w: %w", ErrInvalidData, err)
	}

	access, err = m.GenerateAccess(data)
	if err != nil {
		return "", "", fmt.Errorf("generate access: %w", err)
	}

	refresh, err = m.GenerateRefresh(data)
	if err != nil {
		return "", "", fmt.Errorf("generate refresh: %w", err)
	}

	return access, refresh, nil
}

// GenerateAccess creates a signed access token.
func (m *Manager) GenerateAccess(data TokenData) (string, error) {
	now := time.Now()

	claims := accessJWTClaims{
		UserID:   data.UserID,
		Username: data.Username,
		Email:    data.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(m.accessTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	signed, err := token.SignedString(m.privateKey)
	if err != nil {
		return "", fmt.Errorf("sign access token: %w", err)
	}

	return signed, nil
}

// GenerateRefresh creates a signed refresh token with unique JTI.
func (m *Manager) GenerateRefresh(data TokenData) (string, error) {
	now := time.Now()

	jwtID, err := generateTokenID()
	if err != nil {
		return "", fmt.Errorf("generate jti: %w", err)
	}

	claims := refreshJWTClaims{
		UserID: data.UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jwtID,
			ExpiresAt: jwt.NewNumericDate(now.Add(m.refreshTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	signed, err := token.SignedString(m.privateKey)
	if err != nil {
		return "", fmt.Errorf("sign refresh token: %w", err)
	}

	return signed, nil
}

func (m *Manager) ParseAccess(tokenString string) (*AccessClaims, error) {
	claims := &accessJWTClaims{}

	claims, err := parseToken(tokenString, claims, m.keyFunc)
	if err != nil {
		return nil, err // ErrTokenExpired, ErrTokenParse, или ErrTokenInvalid
	}

	return &AccessClaims{
		UserID:   claims.UserID,
		Username: claims.Username,
		Email:    claims.Email,
		Exp:      claims.ExpiresAt.Time,
	}, nil
}

func (m *Manager) ParseRefresh(tokenString string) (*RefreshClaims, error) {
	claims := &refreshJWTClaims{}

	claims, err := parseToken(tokenString, claims, m.keyFunc)
	if err != nil {
		return nil, err
	}

	return &RefreshClaims{
		JWTID:  claims.ID,
		UserID: claims.UserID,
		Exp:    claims.ExpiresAt.Time,
	}, nil
}

func (m *Manager) keyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, fmt.Errorf("%w: %v", ErrSigningMethod, token.Header["alg"])
	}
	return m.publicKey, nil
}
