package gorediscli

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	*redis.Client
}

func New(ctx context.Context, opts ...Option) (*Client, error) {
	cfg := defaultConfig()

	for _, opt := range opts {
		opt(&cfg)
	}

	if err := cfg.valid(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.addr,
		Password: *cfg.password,
		DB:       cfg.db,

		PoolSize:     cfg.poolSize,
		MinIdleConns: cfg.minIdleConns,

		DialTimeout:  cfg.dialTimeout,
		ReadTimeout:  cfg.readTimeout,
		WriteTimeout: cfg.writeTimeout,

		ConnMaxLifetime: cfg.maxConnLifetime,
		ConnMaxIdleTime: cfg.maxConnIdleTime,
	})

	pingCtx, cancel := context.WithTimeout(ctx, cfg.pingTimeout)
	defer cancel()

	if err := rdb.Ping(pingCtx).Err(); err != nil {
		_ = rdb.Close()
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return &Client{
		Client: rdb,
	}, nil
}
