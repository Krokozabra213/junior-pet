package pgxclient

import "time"

type Option func(*config)

func WithHost(host string) Option {
	return func(c *config) {
		c.host = host
	}
}

func WithPort(port int) Option {
	return func(c *config) {
		c.port = uint16(port)
	}
}

func WithUser(user string) Option {
	return func(c *config) {
		c.user = &user
	}
}

func WithPassword(password string) Option {
	return func(c *config) {
		c.password = &password
	}
}

func WithDatabase(database string) Option {
	return func(c *config) {
		c.database = database
	}
}

func WithSSL(sslmode string) Option {
	return func(c *config) {
		c.sslMode = sslmode
	}
}

func WithConnectionTimeout(timeout time.Duration) Option {
	return func(c *config) {
		c.connectTimeout = timeout
	}
}

func WithMaxConns(maxConns int) Option {
	return func(c *config) {
		c.maxConns = int32(maxConns)
	}
}

func WithMinConns(minConns int) Option {
	return func(c *config) {
		c.minConns = int32(minConns)
	}
}

func WithMaxConnLifetime(lifetime time.Duration) Option {
	return func(c *config) {
		c.maxConnLifeTime = lifetime
	}
}

func WithMaxConnIdletime(idletime time.Duration) Option {
	return func(c *config) {
		c.maxConnIdleTime = idletime
	}
}
