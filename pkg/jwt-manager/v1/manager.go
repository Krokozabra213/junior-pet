package jwtmanager

import (
	"crypto/rsa"
	"time"
)

const (
	ClaimJTI       = "jti"
	ClaimUserID    = "user_id"
	ClaimUsername  = "username"
	ClaimEmail     = "email"
	ClaimExpiredAt = "exp"
)

type JWTManager struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
	accessTTL  time.Duration
	refreshTTL time.Duration
}

type Data struct {
	UserID   uint64
	Username string
	Email    string
}

func New(
	public *rsa.PublicKey, private *rsa.PrivateKey, accessTTL, refreshTTL time.Duration,
) (*JWTManager, error) {
	if public == nil {
		return nil, ErrEmptyPublicKey
	}

	if private == nil {
		return nil, ErrEmptyPrivateKey
	}

	return &JWTManager{
		publicKey: public, privateKey: private, accessTTL: accessTTL, refreshTTL: refreshTTL,
	}, nil
}

func (m *JWTManager) GenerateTokens(data *Data) (string, string, error) {
	accessToken, err := m.GenerateAccess(
		data.UserID, data.Username, data.AppID,
	)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := m.GenerateRefresh(
		data.UserID, data.Username, data.AppID,
	)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
