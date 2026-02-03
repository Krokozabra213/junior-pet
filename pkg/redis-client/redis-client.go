package redisclient

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Addr     string
	Password string
	DB       int

	// Pool
	PoolSize     int
	MinIdleConns int

	// Timeouts
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	// Connection lifecycle
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

func NewConfig(addr, password string, db int, dialTimeout, readTimeout, writeTimeout, maxConnLifetime, maxConnIdleTime time.Duration,
	poolSize, minIdleConns int,
) Config {
	return Config{
		Addr:            addr,
		Password:        password,
		DB:              db,
		DialTimeout:     dialTimeout,
		ReadTimeout:     readTimeout,
		WriteTimeout:    writeTimeout,
		MaxConnLifetime: maxConnLifetime,
		MaxConnIdleTime: maxConnIdleTime,
		PoolSize:        poolSize,
		MinIdleConns:    minIdleConns,
	}
}

type RedisClient struct {
	*redis.Client
}

func New(ctx context.Context, cfg Config) (*RedisClient, error) {
	options := &redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,

		// Pool
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,

		// Timeouts
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,

		// Connection lifecycle
		ConnMaxLifetime: cfg.MaxConnLifetime,
		ConnMaxIdleTime: cfg.MaxConnIdleTime,
	}

	client := redis.NewClient(options)

	// Таймаут на подключение
	pingCtx, cancel := context.WithTimeout(ctx, cfg.DialTimeout)
	defer cancel()

	// Проверяем соединение
	if err := client.Ping(pingCtx).Err(); err != nil {
		client.Close()
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return &RedisClient{
		Client: client,
	}, nil
}

// Close — немедленное закрытие (для defer)
func (c *RedisClient) Close() {
	c.Close()
}
