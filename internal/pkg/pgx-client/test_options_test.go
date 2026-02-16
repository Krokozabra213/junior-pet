package pgxclient

import (
	"testing"
	"time"
)

func TestOptions(t *testing.T) {
	cfg := defaultConfig()

	opts := []Option{
		WithHost("db.example.com"),
		WithPort(5433),
		WithUser("admin"),
		WithPassword("secret"),
		WithDatabase("mydb"),
		WithSSL("require"),
		WithConnectionTimeout(10 * time.Second),
		WithMaxConns(20),
		WithMinConns(5),
		WithMaxConnLifetime(1 * time.Hour),
		WithMaxConnIdletime(30 * time.Minute),
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"host", cfg.host, "db.example.com"},
		{"port", cfg.port, uint16(5433)},
		{"user", *cfg.user, "admin"},
		{"password", *cfg.password, "secret"},
		{"database", cfg.database, "mydb"},
		{"sslMode", cfg.sslMode, "require"},
		{"connectTimeout", cfg.connectTimeout, 10 * time.Second},
		{"maxConns", cfg.maxConns, int32(20)},
		{"minConns", cfg.minConns, int32(5)},
		{"maxConnLifeTime", cfg.maxConnLifeTime, 1 * time.Hour},
		{"maxConnIdleTime", cfg.maxConnIdleTime, 30 * time.Minute},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}
