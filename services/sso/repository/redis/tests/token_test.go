//go:build integration

package tests

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	repository "github.com/Krokozabra213/schools_backend/services/sso/repository/redis"
)

func TestRevokeRefreshToken(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		cleanup(t)

		jti := uuid.New().String()
		expiresAt := time.Now().Add(15 * time.Minute)

		err := testRepo.RevokeRefreshToken(ctx, jti, expiresAt)

		require.NoError(t, err)

		// Проверяем что токен сохранён
		revoked, err := testRepo.IsTokenRevoked(ctx, jti)
		require.NoError(t, err)
		assert.True(t, revoked)
	})

	t.Run("expired token", func(t *testing.T) {
		cleanup(t)

		jti := uuid.New().String()
		expiresAt := time.Now().Add(-1 * time.Minute) // уже истёк

		err := testRepo.RevokeRefreshToken(ctx, jti, expiresAt)

		assert.ErrorIs(t, err, repository.ErrTokenExpired)
	})

	t.Run("token expires after ttl", func(t *testing.T) {
		cleanup(t)

		jti := uuid.New().String()
		expiresAt := time.Now().Add(2 * time.Second) // короткий TTL

		err := testRepo.RevokeRefreshToken(ctx, jti, expiresAt)
		require.NoError(t, err)

		// Сразу есть
		revoked, _ := testRepo.IsTokenRevoked(ctx, jti)
		assert.True(t, revoked)

		// Ждём истечения
		time.Sleep(3 * time.Second)

		// После TTL - нет
		revoked, err = testRepo.IsTokenRevoked(ctx, jti)
		require.NoError(t, err)
		assert.False(t, revoked)
	})
}

func TestIsTokenRevoked(t *testing.T) {
	ctx := context.Background()

	t.Run("token revoked", func(t *testing.T) {
		cleanup(t)

		jti := uuid.New().String()
		testRepo.RevokeRefreshToken(ctx, jti, time.Now().Add(10*time.Minute))

		revoked, err := testRepo.IsTokenRevoked(ctx, jti)

		require.NoError(t, err)
		assert.True(t, revoked)
	})

	t.Run("token not revoked", func(t *testing.T) {
		cleanup(t)

		jti := uuid.New().String() // не сохраняли

		revoked, err := testRepo.IsTokenRevoked(ctx, jti)

		require.NoError(t, err)
		assert.False(t, revoked)
	})
}
