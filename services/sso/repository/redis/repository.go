package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenProvider interface {
	RevokeRefreshToken(ctx context.Context, jti string, expiresAt time.Time) error
	IsTokenRevoked(ctx context.Context, jti string) (bool, error)
}

type RedisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Exists(ctx context.Context, keys ...string) *redis.IntCmd
	Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd
	SetEx(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
}

type RedisRepository struct {
	client RedisClient
}

func NewRepository(client RedisClient) *RedisRepository {
	return &RedisRepository{
		client: client,
	}
}

var _ TokenProvider = (*RedisRepository)(nil)
