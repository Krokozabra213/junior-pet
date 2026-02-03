package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/Krokozabra213/schools_backend/services/sso/domain"
	"github.com/redis/go-redis/v9"
)

func (r *RedisRepository) CacheUserProfile(ctx context.Context, profile *domain.UserCacheProfile, ttl time.Duration) error {
	const op = "repository.CacheUserProfile"
	log := slog.With(
		slog.String("op", op),
		slog.Int64("user_id", profile.ID),
	)

	key := r.userProfileKey(profile.ID)

	data, err := json.Marshal(profile)
	if err != nil {
		log.Error("failed marshal profile", "error", err)
		return ErrInternal
	}

	if err := r.client.Set(ctx, key, data, ttl).Err(); err != nil {
		log.Error("failed set profile", "error", err)
		return ErrInternal
	}

	return nil
}

// GetUserProfile получает профиль из кеша
func (r *RedisRepository) GetUserProfile(ctx context.Context, userID int64) (*domain.UserCacheProfile, error) {
	const op = "repository.GetUserProfile"
	log := slog.With(
		slog.String("op", op),
		slog.Int64("user_id", userID),
	)

	key := r.userProfileKey(userID)

	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrNotFound
		}
		log.Error("failed get profile", "error", err)
		return nil, ErrInternal
	}

	var profile domain.UserCacheProfile
	if err := json.Unmarshal(data, &profile); err != nil {
		log.Error("failed unmarshal profile", "error", err)
		return nil, ErrInternal
	}

	return &profile, nil
}

// UpdateUserProfileTTL обновляет TTL профиля
func (r *RedisRepository) UpdateUserProfileTTL(ctx context.Context, userID int64, ttl time.Duration) error {
	const op = "repository.UpdateUserProfileTTL"
	log := slog.With(
		slog.String("op", op),
		slog.Int64("user_id", userID),
	)

	key := r.userProfileKey(userID)

	if err := r.client.Expire(ctx, key, ttl).Err(); err != nil {
		log.Error("failed update profile", "error", err)
		return ErrInternal
	}

	return nil
}

func (r *RedisRepository) DeleteUserProfile(ctx context.Context, userID int64) error {
	const op = "repository.DeleteUserProfile"
	log := slog.With(
		slog.String("op", op),
		slog.Int64("user_id", userID),
	)

	key := r.userProfileKey(userID)

	if err := r.client.Del(ctx, key).Err(); err != nil {
		log.Error("failed to delete profile", "error", err)
		return ErrInternal
	}

	return nil
}

func (r *RedisRepository) userProfileKey(userID int64) string {
	return fmt.Sprintf("user:profile:%d", userID)
}
