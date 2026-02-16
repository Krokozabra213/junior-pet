package gorediscli

import (
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := defaultConfig()

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"addr", cfg.addr, "localhost:6379"},
		{"db", cfg.db, 0},
		{"password is nil", cfg.password == nil, true},
		{"poolSize", cfg.poolSize, 10},
		{"minIdleConns", cfg.minIdleConns, 2},
		{"dialTimeout", cfg.dialTimeout, 5 * time.Second},
		{"readTimeout", cfg.readTimeout, 3 * time.Second},
		{"writeTimeout", cfg.writeTimeout, 3 * time.Second},
		{"maxConnLifetime", cfg.maxConnLifetime, 2 * time.Hour},
		{"maxConnIdleTime", cfg.maxConnIdleTime, 15 * time.Minute},
		{"pingTimeout", cfg.pingTimeout, 3 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestConfigValid(t *testing.T) {
	password := "secret"

	validBase := func() config {
		cfg := defaultConfig()
		cfg.password = &password
		return cfg
	}

	tests := []struct {
		name    string
		modify  func(*config)
		wantErr bool
	}{
		{
			name:    "valid config",
			modify:  func(c *config) {},
			wantErr: false,
		},
		{
			name:    "nil password",
			modify:  func(c *config) { c.password = nil },
			wantErr: true,
		},
		{
			name:    "empty addr",
			modify:  func(c *config) { c.addr = "" },
			wantErr: true,
		},
		{
			name:    "negative db",
			modify:  func(c *config) { c.db = -1 },
			wantErr: true,
		},
		{
			name:    "db too high",
			modify:  func(c *config) { c.db = 16 },
			wantErr: true,
		},
		{
			name:    "db max valid",
			modify:  func(c *config) { c.db = 15 },
			wantErr: false,
		},
		{
			name:    "zero pool size",
			modify:  func(c *config) { c.poolSize = 0 },
			wantErr: true,
		},
		{
			name:    "negative minIdleConns",
			modify:  func(c *config) { c.minIdleConns = -1 },
			wantErr: true,
		},
		{
			name: "minIdleConns exceeds poolSize",
			modify: func(c *config) {
				c.poolSize = 5
				c.minIdleConns = 10
			},
			wantErr: true,
		},
		{
			name: "minIdleConns equals poolSize",
			modify: func(c *config) {
				c.poolSize = 5
				c.minIdleConns = 5
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := validBase()
			tt.modify(&cfg)

			err := cfg.valid()
			if (err != nil) != tt.wantErr {
				t.Errorf("valid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
