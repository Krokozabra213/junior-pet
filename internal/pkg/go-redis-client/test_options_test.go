package gorediscli

import (
	"testing"
	"time"
)

func TestOptions(t *testing.T) {
	cfg := defaultConfig()

	opts := []Option{
		WithAddr("redis.example.com:6380"),
		WithPassword("supersecret"),
		WithDB(3),
		WithPoolSize(20),
		WithMinIdleConns(5),
		WithDialTimeout(10 * time.Second),
		WithReadTimeout(5 * time.Second),
		WithWriteTimeout(5 * time.Second),
		WithMaxConnLifetime(1 * time.Hour),
		WithMaxConnIdleTime(30 * time.Minute),
		WithPingTimeout(2 * time.Second),
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"addr", cfg.addr, "redis.example.com:6380"},
		{"password", *cfg.password, "supersecret"},
		{"db", cfg.db, 3},
		{"poolSize", cfg.poolSize, 20},
		{"minIdleConns", cfg.minIdleConns, 5},
		{"dialTimeout", cfg.dialTimeout, 10 * time.Second},
		{"readTimeout", cfg.readTimeout, 5 * time.Second},
		{"writeTimeout", cfg.writeTimeout, 5 * time.Second},
		{"maxConnLifetime", cfg.maxConnLifetime, 1 * time.Hour},
		{"maxConnIdleTime", cfg.maxConnIdleTime, 30 * time.Minute},
		{"pingTimeout", cfg.pingTimeout, 2 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestOptionOverridesDefault(t *testing.T) {
	cfg := defaultConfig()

	if cfg.addr != "localhost:6379" {
		t.Fatal("precondition failed: unexpected default addr")
	}

	WithAddr("other:6380")(&cfg)

	if cfg.addr != "other:6380" {
		t.Errorf("option did not override default, got %s", cfg.addr)
	}
}

func TestPasswordOptionSetsPointer(t *testing.T) {
	cfg := defaultConfig()

	if cfg.password != nil {
		t.Fatal("precondition failed: password should be nil by default")
	}

	WithPassword("test")(&cfg)

	if cfg.password == nil {
		t.Fatal("password should not be nil after WithPassword")
	}
	if *cfg.password != "test" {
		t.Errorf("got %s, want test", *cfg.password)
	}
}
