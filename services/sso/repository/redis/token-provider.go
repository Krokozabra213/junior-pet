package redis

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

const tokenPrefix = "token:revoked:"

// RevokeRefreshToken сохраняет refresh токен по JTI
func (r *RedisRepository) RevokeRefreshToken(ctx context.Context, jti string, expiresAt time.Time) error {
	const op = "repository.RevokeRefreshToken"
	log := slog.With(
		slog.String("op", op),
		slog.String("jti", jti),
	)

	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return ErrTokenExpired
	}

	key := r.refreshTokenKey(jti)

	if err := r.client.SetEx(ctx, key, "", ttl).Err(); err != nil {
		log.Error("failed save token", "error", err)
		return ErrInternal
	}

	return nil
}

// IsTokenRevoked проверяет существование токена
func (r *RedisRepository) IsTokenRevoked(ctx context.Context, jti string) (bool, error) {
	const op = "repository.IsTokenRevoked"
	log := slog.With(
		slog.String("op", op),
		slog.String("jti", jti),
	)

	key := r.refreshTokenKey(jti)

	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		log.Error("failed check token", "error", err)
		return false, ErrInternal
	}

	return exists == 1, nil
}

// refreshTokenKey генерирует ключ для refresh токена
func (r *RedisRepository) refreshTokenKey(jti string) string {
	return fmt.Sprintf("%s:%s", tokenPrefix, jti)
}
