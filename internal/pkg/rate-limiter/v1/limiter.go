package ratelimiterv1

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Limiter checks if a request is allowed.
type Limiter interface {
	Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, error)
}

// RedisLimiter implements sliding window rate limiting with Redis.
type RedisLimiter struct {
	client redis.Scripter
}

func NewRedisLimiter(client redis.Cmdable) *RedisLimiter {
	return &RedisLimiter{client: client}
}

// Lua script: sliding window counter.
// Атомарная операция — никаких race conditions.
var slidingWindowScript = redis.NewScript(`
	local key = KEYS[1]
	local window = tonumber(ARGV[1])
	local limit = tonumber(ARGV[2])
	local now = tonumber(ARGV[3])

	-- Генерируем уникальную часть
	local time_parts = redis.call('TIME')
	local micro_total = time_parts[1] * 1000000 + time_parts[2]
	math.randomseed(micro_total)
	local unique = math.random(1000000000)

	-- Удаляем записи старше окна
	redis.call('ZREMRANGEBYSCORE', key, '-inf', now - window)

	-- Считаем текущие запросы в окне
	local count = redis.call('ZCARD', key)

	if count < limit then
    -- Добавляем элемент с составным member
    redis.call('ZADD', key, now, now .. '-' .. micro_total .. '-' .. unique)
    redis.call('PEXPIRE', key, window)
    return 1
	end

	return 0
`)

// Allow returns true if request is within the rate limit.
func (r *RedisLimiter) Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	now := time.Now().UnixMilli()
	windowMs := window.Milliseconds()

	result, err := slidingWindowScript.Run(ctx, r.client,
		[]string{key},
		windowMs,
		limit,
		now,
	).Int()

	if err != nil {
		return false, fmt.Errorf("rate limiter script: %w", err)
	}

	return result == 1, nil
}
