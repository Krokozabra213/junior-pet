package pgxclient

import (
	"errors"
	"time"
)

const (
	defaultHost            = "localhost"
	defaultPort            = 5432
	defaultDatabase        = "postgres"
	defaultSSLMode         = "disable"
	defaultConnectTimeout  = 5 * time.Second
	defaultMaxConns        = 10
	defaultMinConns        = 2
	defaultMaxConnLifeTime = 2 * time.Hour
	defaultMaxConnIdleTime = 15 * time.Minute
)

type config struct {
	host            string
	port            uint16
	user            *string
	password        *string
	database        string
	sslMode         string
	connectTimeout  time.Duration
	maxConns        int32
	minConns        int32
	maxConnLifeTime time.Duration
	maxConnIdleTime time.Duration
}

func defaultConfig() config {
	return config{
		host:            defaultHost,
		port:            defaultPort,
		database:        defaultDatabase,
		sslMode:         defaultSSLMode,
		connectTimeout:  defaultConnectTimeout,
		maxConnLifeTime: defaultMaxConnLifeTime,
		maxConnIdleTime: defaultMaxConnIdleTime,
		maxConns:        defaultMaxConns,
		minConns:        defaultMinConns,
	}
}

func (c config) valid() error {
	if c.user == nil {
		return errors.New("user is required")
	}
	if c.password == nil {
		return errors.New("password is required")
	}
	if c.port == 0 {
		return errors.New("invalid port")
	}
	if c.maxConns < 1 {
		return errors.New("maxConns must be >= 1")
	}
	if c.minConns < 0 {
		return errors.New("minConns must be >= 0")
	}
	if c.minConns > c.maxConns {
		return errors.New("minConns must be <= maxConns")
	}
	return nil
}
