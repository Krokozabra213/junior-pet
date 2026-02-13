package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/Krokozabra213/schools_backend/pkg/logger"
	redisclient "github.com/Krokozabra213/schools_backend/pkg/redis-client"
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

	// Redis
	redisCfg := redisclient.NewConfig(
		cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.Database,
		cfg.Redis.DialTimeout, cfg.Redis.ReadTimeout, cfg.Redis.WriteTimeout,
		cfg.Redis.ConnMaxLifetime, cfg.Redis.ConnMaxIdletime,
		cfg.Redis.PoolSize, cfg.Redis.MinIdleConns,
	)

	redisClient, err := redisclient.New(context.Background(), redisCfg)
	if err != nil {
		log.Error("redis connect error: %v", "error", err)
	}
	defer redisClient.CloseConn()
	log.Info("connected to redis")

	return nil
}
