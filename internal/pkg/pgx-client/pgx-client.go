package pgxclient

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ Database = (*Client)(nil)

type Database interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Begin(ctx context.Context) (pgx.Tx, error)
	Ping(ctx context.Context) error
	Pool() *pgxpool.Pool
	Close()
}

type Client struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, opts ...Option) (*Client, error) {
	cfg := defaultConfig()

	for _, opt := range opts {
		opt(&cfg)
	}

	if err := cfg.valid(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	poolCfg, err := createPGXConfig(&cfg)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	// ленивое соединение
	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, cfg.connectTimeout)
	defer cancel()

	// Проверяем соединение
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return &Client{
		pool: pool,
	}, nil
}

func createPGXConfig(cfg *config) (*pgxpool.Config, error) {
	pgxConf, err := pgxpool.ParseConfig("")
	if err != nil {
		return nil, err
	}

	pgxConf.ConnConfig.Host = cfg.host
	pgxConf.ConnConfig.Port = cfg.port
	pgxConf.ConnConfig.User = *cfg.user
	pgxConf.ConnConfig.Password = *cfg.password
	pgxConf.ConnConfig.Database = cfg.database

	if err := configureTLS(pgxConf.ConnConfig, cfg.sslMode); err != nil {
		return nil, fmt.Errorf("configure TLS: %w", err)
	}

	pgxConf.MaxConns = cfg.maxConns
	pgxConf.MinConns = cfg.minConns
	pgxConf.MaxConnLifetime = cfg.maxConnLifeTime
	pgxConf.MaxConnIdleTime = cfg.maxConnIdleTime

	return pgxConf, nil
}

func configureTLS(connConfig *pgx.ConnConfig, sslMode string) error {
	switch sslMode {
	case "disable":
		connConfig.TLSConfig = nil

	case "require":
		connConfig.TLSConfig = &tls.Config{
			InsecureSkipVerify: true, // шифрование есть, проверки сертификата нет
		}

	default:
		return fmt.Errorf("unknown sslmode: %s", sslMode)
	}

	return nil
}

func (c *Client) Close() {
	c.pool.Close()
}

func (c *Client) Pool() *pgxpool.Pool {
	return c.pool
}

func (c *Client) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return c.pool.Query(ctx, sql, args...)
}

func (c *Client) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return c.pool.Exec(ctx, sql, args...)
}

func (c *Client) Begin(ctx context.Context) (pgx.Tx, error) {
	return c.pool.Begin(ctx)
}

func (c *Client) Ping(ctx context.Context) error {
	return c.pool.Ping(ctx)
}

func (c *Client) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return c.pool.QueryRow(ctx, sql, args...)
}
