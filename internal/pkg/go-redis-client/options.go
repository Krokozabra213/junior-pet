package gorediscli

import "time"

type Option func(*config)

func WithAddr(addr string) Option {
	return func(c *config) {
		c.addr = addr
	}
}

func WithPassword(password string) Option {
	return func(c *config) {
		c.password = &password
	}
}

func WithDB(db int) Option {
	return func(c *config) {
		c.db = db
	}
}

func WithPoolSize(size int) Option {
	return func(c *config) {
		c.poolSize = size
	}
}

func WithMinIdleConns(minIdle int) Option {
	return func(c *config) {
		c.minIdleConns = minIdle
	}
}

func WithDialTimeout(timeout time.Duration) Option {
	return func(c *config) {
		c.dialTimeout = timeout
	}
}

func WithReadTimeout(timeout time.Duration) Option {
	return func(c *config) {
		c.readTimeout = timeout
	}
}

func WithWriteTimeout(timeout time.Duration) Option {
	return func(c *config) {
		c.writeTimeout = timeout
	}
}

func WithMaxConnLifetime(lifetime time.Duration) Option {
	return func(c *config) {
		c.maxConnLifetime = lifetime
	}
}

func WithMaxConnIdleTime(idleTime time.Duration) Option {
	return func(c *config) {
		c.maxConnIdleTime = idleTime
	}
}

func WithPingTimeout(timeout time.Duration) Option {
	return func(c *config) {
		c.pingTimeout = timeout
	}
}
