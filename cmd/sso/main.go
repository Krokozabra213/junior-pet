package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/Krokozabra213/schools_backend/internal/pkg/logger"
	ssoconfig "github.com/Krokozabra213/schools_backend/services/sso/config"
)

func main() {
	if err := run(); err != nil {
		slog.Error("application failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	var configPath string
	flag.StringVar(&configPath, "config", "configs/sso.yaml", "path to configuration file")
	flag.Parse()

	// // Config
	cfg, err := ssoconfig.Init(configPath)
	if err != nil {
		return err
	}

	// Logger
	log := logger.Init(cfg.App.Environment).WithService("sso", "v0.1.0")
	log.Info("initialized config", "config", cfg.LogValue())
	log.Info("starting application")

	log.Info(configPath)

	// Database
	log.Info("connected to postgres")

	// Cache
	log.Info("connected to redis")

	return nil
}
