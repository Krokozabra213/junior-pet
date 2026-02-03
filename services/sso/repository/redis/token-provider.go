package redis

import (
	"context"
	"log/slog"
	"time"
)

const tokenPrefix = "token:refresh:"

func (r *RedisRepository) SaveToken(ctx context.Context, userID int64, token string, expiresAt time.Time) error {
	const op = "repository.SaveToken"
	log := slog.With(
		slog.String("op", op),
		slog.Int64("user_id", userID),
	)

	expiration := time.Until(expiresAt)
	if expiration <= 0 {
		log.Warn("token already expired")
		return ErrTokenExpired
	}

	key := r.tokenKey(token)

	if err := r.client.SetEx(ctx, key, userID, expiration).Err(); err != nil {
		log.Error("failed save token", "error", err)
		return ErrInternal
	}

	return nil
}

func (r *RedisRepository) CheckToken(ctx context.Context, token string) (bool, error) {
	const op = "repository.CheckToken"
	log := slog.With(
		slog.String("op", op),
	)

	exists, err := r.client.Exists(ctx, token).Result()
	if err != nil {
		log.Error("failed check token", "error", err)
		return false, ErrInternal
	}

	return exists == 1, nil
}

// tokenKey генерирует ключ для токена
func (r *RedisRepository) tokenKey(token string) string {
	return tokenPrefix + token
}
