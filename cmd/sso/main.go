package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/Krokozabra213/schools_backend/pkg/logger"
	pgxclient "github.com/Krokozabra213/schools_backend/pkg/pgx-client"
	ssoconfig "github.com/Krokozabra213/schools_backend/services/sso/config"
)

const (
	// configFile      = "configs/main.yml"
	// envFile         = ".env"
	shutdownTimeout = 5 * time.Second
)

func main() {
	if err := run(); err != nil {
		slog.Error("application failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	// // Config
	cfg, err := ssoconfig.Init("configs/sso.yaml", ".env")
	if err != nil {
		return err
	}

	// Logger
	log := logger.Init(cfg.App.Environment)
	log.Info("initialized config", "config", cfg.LogValue())
	log.Info("starting application")

	// Database
	pgxConf := pgxclient.NewPGXConfig(cfg.PG.Host, cfg.PG.Port, cfg.PG.User, cfg.PG.Password,
		cfg.PG.DBName, cfg.PG.SSLMode, cfg.PG.ConnectTimeout, cfg.PG.MaxConnLifeTime,
		cfg.PG.MaxConnIdleTime, cfg.PG.MaxConns, cfg.PG.MinConns,
	)
	pgxClient, err := pgxclient.New(context.Background(), pgxConf)
	if err != nil {
		return err
	}
	defer func() {
		if err := pgxClient.Shutdown(shutdownTimeout); err != nil {
			log.Error("database shutdown error", "error", err)
		}
	}()
	log.Info("connected to postgres")

	return nil
}
