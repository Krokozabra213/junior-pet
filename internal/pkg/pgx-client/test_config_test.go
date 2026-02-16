package pgxclient

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
		{"host", cfg.host, "localhost"},
		{"port", cfg.port, uint16(5432)},
		{"database", cfg.database, "postgres"},
		{"sslMode", cfg.sslMode, "disable"},
		{"connectTimeout", cfg.connectTimeout, 5 * time.Second},
		{"maxConns", cfg.maxConns, int32(10)},
		{"minConns", cfg.minConns, int32(2)},
		{"maxConnLifeTime", cfg.maxConnLifeTime, 2 * time.Hour},
		{"maxConnIdleTime", cfg.maxConnIdleTime, 15 * time.Minute},
		{"user is nil", cfg.user == nil, true},
		{"password is nil", cfg.password == nil, true},
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
	user := "testuser"
	password := "testpass"

	tests := []struct {
		name    string
		cfg     config
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: config{
				user:     &user,
				password: &password,
				port:     5432,
				maxConns: 10,
				minConns: 2,
			},
			wantErr: false,
		},
		{
			name: "nil user",
			cfg: config{
				user:     nil,
				password: &password,
			},
			wantErr: true,
		},
		{
			name: "nil password",
			cfg: config{
				user:     &user,
				password: nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.valid()
			if (err != nil) != tt.wantErr {
				t.Errorf("valid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
