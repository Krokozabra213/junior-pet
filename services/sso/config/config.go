package ssoconfig

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	App   AppConfig      `yaml:"app"`
	HTTP  HTTPConfig     `yaml:"http"`
	GRPC  GRPCConfig     `yaml:"grpc"`
	PG    PostgresConfig `yaml:"postgres"`
	Redis RedisConfig    `yaml:"redis"`
	JWT   JWTConfig      `yaml:"jwt"`
}

type AppConfig struct {
	Environment  string `env:"SSO_ENV" env-default:"development"`
	AppSecretKey string `env:"SSO_APP_SECRET" env-required:"true"`
}

type PostgresConfig struct {
	Host     string `env:"SSO_POSTGRES_HOST" env-required:"true"`
	Port     string `env:"SSO_POSTGRES_PORT" env-default:"5432"`
	User     string `env:"SSO_POSTGRES_USER" env-required:"true"`
	Password string `env:"SSO_POSTGRES_PASSWORD" env-required:"true"`
	DBName   string `env:"SSO_POSTGRES_DB" env-required:"true"`

	MaxConns int `yaml:"maxConns" env:"SSO_PG_MAX_CONNS" env-default:"10"`
	MinConns int `yaml:"minConns" env:"SSO_PG_MIN_CONNS" env-default:"2"`
}

type RedisConfig struct {
	Addr     string `env:"SSO_REDIS_ADDR" env-required:"true"`
	Password string `env:"SSO_REDIS_PASSWORD" env-required:"true"`
	Database int    `env:"SSO_REDIS_DATABASE" env-default:"0"`

	PoolSize     int `yaml:"poolSize" env:"SSO_REDIS_POOL_SIZE" env-default:"10"`
	MinIdleConns int `yaml:"minIdleConns" env:"SSO_REDIS_MIN_IDLE_CONNS" env-default:"2"`
}

type HTTPConfig struct {
	Host               string        `yaml:"host" env:"SSO_HTTP_HOST" env-default:"0.0.0.0"`
	Port               string        `yaml:"port" env:"SSO_HTTP_PORT" env-default:"8080"`
	ReadTimeout        time.Duration `yaml:"readTimeout" env:"SSO_HTTP_READ_TIMEOUT" env-default:"10s"`
	WriteTimeout       time.Duration `yaml:"writeTimeout" env:"SSO_HTTP_WRITE_TIMEOUT" env-default:"10s"`
	MaxHeaderMegabytes int           `yaml:"maxHeaderBytes" env:"SSO_HTTP_MAX_HEADER_BYTES" env-default:"1"`
}

type GRPCConfig struct {
	Host               string        `yaml:"host" env:"SSO_GRPC_HOST" env-default:"0.0.0.0"`
	Port               string        `yaml:"port" env:"SSO_GRPC_PORT" env-default:"44050"`
	ReadTimeout        time.Duration `yaml:"readTimeout" env:"SSO_GRPC_READ_TIMEOUT" env-default:"10s"`
	WriteTimeout       time.Duration `yaml:"writeTimeout" env:"SSO_GRPC_WRITE_TIMEOUT" env-default:"10s"`
	MaxHeaderMegabytes int           `yaml:"maxHeaderBytes" env:"SSO_GRPC_MAX_HEADER_BYTES" env-default:"1"`
}

type JWTConfig struct {
	AccessTokenTTL  time.Duration `yaml:"accessTokenTTL" env:"SSO_JWT_ACCESS_TTL" env-default:"15m"`
	RefreshTokenTTL time.Duration `yaml:"refreshTokenTTL" env:"SSO_JWT_REFRESH_TTL" env-default:"720h"`
	PrivateKeyPath  string        `yaml:"privateKeyPath" env:"SSO_JWT_PRIVATE_KEY_PATH" env-default:"private.pem"`
	PrivateKey      []byte        `yaml:"-" env:"-"`
}

func MustInit(configFile string) *Config {
	cfg, err := Init(configFile)
	if err != nil {
		panic("config: " + err.Error())
	}
	return cfg
}

func Init(configFile string) (*Config, error) {

	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("load env file: %w", err)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configFile, &cfg); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	keyData, err := os.ReadFile(cfg.JWT.PrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("read private key %q: %w", cfg.JWT.PrivateKeyPath, err)
	}
	cfg.JWT.PrivateKey = keyData

	return &cfg, nil
}

func (c *Config) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("env", c.App.Environment),

		slog.Group("http",
			slog.String("address", c.HTTP.Host+":"+c.HTTP.Port),
			slog.Duration("read_timeout", c.HTTP.ReadTimeout),
			slog.Duration("write_timeout", c.HTTP.WriteTimeout),
			slog.Int("max_header_megabytes", c.HTTP.MaxHeaderMegabytes),
		),

		slog.Group("grpc",
			slog.String("address", c.GRPC.Host+":"+c.GRPC.Port),
			slog.Duration("read_timeout", c.GRPC.ReadTimeout),
			slog.Duration("write_timeout", c.GRPC.WriteTimeout),
			slog.Int("max_header_megabytes", c.GRPC.MaxHeaderMegabytes),
		),

		slog.Group("jwt",
			slog.Duration("access_token_ttl", c.JWT.AccessTokenTTL),
			slog.Duration("refresh_token_ttl", c.JWT.RefreshTokenTTL),
			slog.String("private_key_path", c.JWT.PrivateKeyPath),
		),

		slog.Group("postgres",
			slog.String("address", c.PG.Host+":"+c.PG.Port),
			slog.String("database", c.PG.DBName),
			slog.String("user", c.PG.User),
			slog.Int("max_conns", c.PG.MaxConns),
			slog.Int("min_conns", c.PG.MinConns),
		),

		slog.Group("redis",
			slog.String("address", c.Redis.Addr),
			slog.Int("database", c.Redis.Database),
			slog.Int("pool_size", c.Redis.PoolSize),
			slog.Int("min_idle_conns", c.Redis.MinIdleConns),
		),
	)
}
