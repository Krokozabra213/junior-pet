package gorediscli

import (
	"errors"
	"time"
)

const (
	defaultAddr            = "localhost:6379"
	defaultDB              = 0
	defaultPoolSize        = 10
	defaultMinIdleConns    = 2
	defaultDialTimeout     = 5 * time.Second
	defaultReadTimeout     = 3 * time.Second
	defaultWriteTimeout    = 3 * time.Second
	defaultMaxConnLifetime = 2 * time.Hour
	defaultMaxConnIdleTime = 15 * time.Minute
	defaultPingTimeout     = 3 * time.Second
)

type config struct {
	addr            string
	password        *string
	db              int
	poolSize        int
	minIdleConns    int
	dialTimeout     time.Duration
	readTimeout     time.Duration
	writeTimeout    time.Duration
	maxConnLifetime time.Duration
	maxConnIdleTime time.Duration
	pingTimeout     time.Duration
}

func defaultConfig() config {
	return config{
		addr:            defaultAddr,
		db:              defaultDB,
		poolSize:        defaultPoolSize,
		minIdleConns:    defaultMinIdleConns,
		dialTimeout:     defaultDialTimeout,
		readTimeout:     defaultReadTimeout,
		writeTimeout:    defaultWriteTimeout,
		maxConnLifetime: defaultMaxConnLifetime,
		maxConnIdleTime: defaultMaxConnIdleTime,
		pingTimeout:     defaultPingTimeout,
	}
}

func (c config) valid() error {
	if c.addr == "" {
		return errors.New("addr is required")
	}
	if c.password == nil {
		return errors.New("password is required")
	}
	if c.db < 0 || c.db > 15 {
		return errors.New("db must be between 0 and 15")
	}
	if c.poolSize < 1 {
		return errors.New("poolSize must be >= 1")
	}
	if c.minIdleConns < 0 {
		return errors.New("minIdleConns must be >= 0")
	}
	if c.minIdleConns > c.poolSize {
		return errors.New("minIdleConns must be <= poolSize")
	}

	return nil
}
