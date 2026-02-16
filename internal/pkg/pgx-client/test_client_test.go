package pgxclient

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
)

func TestConfigureTLS(t *testing.T) {
	tests := []struct {
		name       string
		sslMode    string
		wantNilTLS bool
		wantErr    bool
	}{
		{
			name:       "disable",
			sslMode:    "disable",
			wantNilTLS: true,
			wantErr:    false,
		},
		{
			name:       "require",
			sslMode:    "require",
			wantNilTLS: false,
			wantErr:    false,
		},
		{
			name:    "unknown mode",
			sslMode: "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			connConfig := &pgx.ConnConfig{}

			err := configureTLS(connConfig, tt.sslMode)

			if (err != nil) != tt.wantErr {
				t.Errorf("configureTLS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if tt.wantNilTLS && connConfig.TLSConfig != nil {
				t.Error("expected nil TLSConfig for disable mode")
			}

			if !tt.wantNilTLS {
				if connConfig.TLSConfig == nil {
					t.Fatal("expected non-nil TLSConfig")
				}
				if !connConfig.TLSConfig.InsecureSkipVerify {
					t.Error("require mode should set InsecureSkipVerify=true")
				}
			}
		})
	}
}

func TestCreatePGXConfig(t *testing.T) {
	user := "testuser"
	password := "testpass"

	cfg := &config{
		host:            "db.example.com",
		port:            5433,
		user:            &user,
		password:        &password,
		database:        "testdb",
		sslMode:         "disable",
		maxConns:        20,
		minConns:        5,
		maxConnLifeTime: 1 * time.Hour,
		maxConnIdleTime: 30 * time.Minute,
	}

	pgxCfg, err := createPGXConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"host", pgxCfg.ConnConfig.Host, "db.example.com"},
		{"port", pgxCfg.ConnConfig.Port, uint16(5433)},
		{"user", pgxCfg.ConnConfig.User, "testuser"},
		{"password", pgxCfg.ConnConfig.Password, "testpass"},
		{"database", pgxCfg.ConnConfig.Database, "testdb"},
		{"maxConns", pgxCfg.MaxConns, int32(20)},
		{"minConns", pgxCfg.MinConns, int32(5)},
		{"maxConnLifetime", pgxCfg.MaxConnLifetime, 1 * time.Hour},
		{"maxConnIdleTime", pgxCfg.MaxConnIdleTime, 30 * time.Minute},
		{"tls is nil", pgxCfg.ConnConfig.TLSConfig == nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestCreatePGXConfigInvalidSSL(t *testing.T) {
	user := "testuser"
	password := "testpass"

	cfg := &config{
		user:     &user,
		password: &password,
		sslMode:  "invalid",
	}

	_, err := createPGXConfig(cfg)
	if err == nil {
		t.Error("expected error for invalid sslMode")
	}
}
