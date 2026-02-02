package ssoconfig

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Default values
const (
	// Redis values
	defaultRedisPoolSize     = 20
	defaultRedisMinIdleConns = 5
	defaultRedisDialTimeout  = 3 * time.Second
	defaultRedisReadTimeout  = 2 * time.Second
	defaultRedisWriteTimeout = 2 * time.Second

	// Postgres values
	defaultPGSSLMode         = "disable"
	defaultPGConnectTimeout  = 5 * time.Second
	defaultPGMaxConns        = 10
	defaultPGMinConns        = 2
	defaultPGMaxConnLifeTime = 1 * time.Hour
	defaultPGMaxConnIdleTime = 15 * time.Minute

	// HTTP values
	defaultHTTPHost               = "localhost"
	defaultHTTPPort               = "8080"
	defaultHTTPWriteTimeout       = 10 * time.Second
	defaultHTTPReadTimeout        = 10 * time.Second
	defaultHTTPMaxHeaderMegabytes = 1

	// GRPC values
	defaultGRPCHost               = "localhost"
	defaultGRPCPort               = "44050"
	defaultGRPCReadTimeout        = 10 * time.Second
	defaultGRPCWriteTimeout       = 10 * time.Second
	defaultGRPCMaxHeaderMegabytes = 1

	// JWT values
	defaultAccessTokenTTL  = 15 * time.Minute
	defaultRefreshTokenTTL = 24 * time.Hour * 30
	defaultPrivateKeyPath  = "private.pem"
)

type (
	PostgresConfig struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string

		SSLMode         string        `mapstructure:"sslMode"`
		ConnectTimeout  time.Duration `mapstructure:"connectTimeout"`
		MaxConns        int           `mapstructure:"maxConns"`
		MinConns        int           `mapstructure:"minConns"`
		MaxConnLifeTime time.Duration `mapstructure:"maxConnLifeTime"`
		MaxConnIdleTime time.Duration `mapstructure:"maxConnIdleTime"`
	}

	GRPCConfig struct {
		Host               string        `mapstructure:"host"`
		Port               string        `mapstructure:"port"`
		ReadTimeout        time.Duration `mapstructure:"readTimeout"`
		WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
		MaxHeaderMegabytes int           `mapstructure:"maxHeaderBytes"`
	}

	HTTPConfig struct {
		Host               string        `mapstructure:"host"`
		Port               string        `mapstructure:"port"`
		ReadTimeout        time.Duration `mapstructure:"readTimeout"`
		WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
		MaxHeaderMegabytes int           `mapstructure:"maxHeaderBytes"`
	}

	AppConfig struct {
		AppSecretKey string
		Environment  string
	}

	JWTConfig struct {
		AccessTokenTTL  time.Duration `mapstructure:"accessTokenTTL"`
		RefreshTokenTTL time.Duration `mapstructure:"refreshTokenTTL"`
		PrivateKey      []byte
	}

	RedisConfig struct {
		Addr     string
		Password string
		Database int

		PoolSize     int           `mapstructure:"poolSize"`
		MinIdleConns int           `mapstructure:"minIdleConns"`
		DialTimeout  time.Duration `mapstructure:"dialTimeout"`
		ReadTimeout  time.Duration `mapstructure:"readTimeout"`
		WriteTimeout time.Duration `mapstructure:"writeTimeout"`
	}
)

type Config struct {
	JWT   JWTConfig
	App   AppConfig
	PG    PostgresConfig
	Redis RedisConfig
	HTTP  HTTPConfig
	GRPC  GRPCConfig
}

func newCfg() Config {
	cfg := Config{
		JWT:   JWTConfig{},
		App:   AppConfig{},
		PG:    PostgresConfig{},
		Redis: RedisConfig{},
		HTTP:  HTTPConfig{},
		GRPC:  GRPCConfig{},
	}
	return cfg
}

func populateDefault() {
	// HTTP defaults
	viper.SetDefault("http.host", defaultHTTPHost)
	viper.SetDefault("http.port", defaultHTTPPort)
	viper.SetDefault("http.maxHeaderMegabytes", defaultHTTPMaxHeaderMegabytes)
	viper.SetDefault("http.readTimeout", defaultHTTPReadTimeout)
	viper.SetDefault("http.writeTimeout", defaultHTTPWriteTimeout)

	// GRPC defaults
	viper.SetDefault("grpc.host", defaultGRPCHost)
	viper.SetDefault("grpc.port", defaultGRPCPort)
	viper.SetDefault("grpc.maxHeaderMegabytes", defaultGRPCMaxHeaderMegabytes)
	viper.SetDefault("grpc.readTimeout", defaultGRPCReadTimeout)
	viper.SetDefault("grpc.writeTimeout", defaultGRPCWriteTimeout)

	// JWT default
	viper.SetDefault("jwt.accessTokenTTL", defaultAccessTokenTTL)
	viper.SetDefault("jwt.refreshTokenTTL", defaultRefreshTokenTTL)
	viper.SetDefault("jwt.privateKeyPath", defaultPrivateKeyPath)

	// Redis defaults
	viper.SetDefault("redis.poolSize", defaultRedisPoolSize)
	viper.SetDefault("redis.minIdleConns", defaultRedisMinIdleConns)
	viper.SetDefault("redis.dialTimeout", defaultRedisDialTimeout)
	viper.SetDefault("redis.readTimeout", defaultRedisReadTimeout)
	viper.SetDefault("redis.writeTimeout", defaultRedisWriteTimeout)

	// Postgres defaults
	viper.SetDefault("postgres.sslMode", defaultPGSSLMode)
	viper.SetDefault("postgres.connectTimeout", defaultPGConnectTimeout)
	viper.SetDefault("postgres.maxConns", defaultPGMaxConns)
	viper.SetDefault("postgres.minConns", defaultPGMinConns)
	viper.SetDefault("postgres.maxConnLifeTime", defaultPGMaxConnLifeTime)
	viper.SetDefault("postgres.maxConnIdleTime", defaultPGMaxConnIdleTime)
}

func Init(configFile, envFile string) (*Config, error) {
	populateDefault()

	if err := parseConfigFile(configFile); err != nil {
		return nil, err
	}

	cfg := newCfg()

	err := unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	err = setFromEnv(envFile, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func parseConfigFile(configPath string) error {
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("http", &cfg.HTTP); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("grpc", &cfg.GRPC); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("redis", &cfg.Redis); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("postgres", &cfg.PG); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("jwt", &cfg.JWT); err != nil {
		return err
	}

	keyPath := viper.GetString("jwt.privateKeyPath")
	data, err := PrivateKeyData(keyPath)
	if err != nil {
		return err
	}
	cfg.JWT.PrivateKey = data

	return nil
}

func PrivateKeyData(keyPath string) ([]byte, error) {
	if keyPath == "" {
		return []byte(nil), errors.New("absent private key path on env file")
	}
	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		return []byte(nil), fmt.Errorf("failed to read private key file: %v", err)
	}

	return keyData, nil
}

func setFromEnv(envpath string, cfg *Config) error {
	err := godotenv.Load(envpath)
	if err != nil {
		return err
	}

	// APP
	cfg.App.AppSecretKey = os.Getenv("SSO_APP_SECRET")
	cfg.App.Environment = os.Getenv("SSO_ENV")

	// POSTGRES
	cfg.PG.User = os.Getenv("SSO_POSTGRES_USER")
	cfg.PG.Host = os.Getenv("SSO_POSTGRES_HOST")
	cfg.PG.Port = os.Getenv("SSO_POSTGRES_PORT")
	cfg.PG.DBName = os.Getenv("SSO_POSTGRES_DB")
	cfg.PG.Password = os.Getenv("SSO_POSTGRES_PASSWORD")

	// REDIS
	cfg.Redis.Addr = os.Getenv("SSO_REDIS_ADDR")
	cfg.Redis.Password = os.Getenv("SSO_REDIS_PASSWORD")
	cfg.Redis.Database, err = strconv.Atoi(os.Getenv("SSO_REDIS_DATABASE"))
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("env", c.App.Environment),
		slog.Group("http",
			slog.String("address", c.HTTP.Host+":"+c.HTTP.Port),
			slog.Duration("read_timeout", c.HTTP.ReadTimeout),
			slog.Duration("write_timeout", c.HTTP.WriteTimeout),
			slog.Int("maxHeaderMegabytes", c.HTTP.MaxHeaderMegabytes),
		),
		slog.Group("grpc",
			slog.String("address", c.GRPC.Host+":"+c.HTTP.Port),
			slog.Duration("read_timeout", c.GRPC.ReadTimeout),
			slog.Duration("write_timeout", c.GRPC.WriteTimeout),
			slog.Int("maxHeaderMegabytes", c.GRPC.MaxHeaderMegabytes),
		),
		slog.Group("jwt",
			slog.Duration("access_token_ttl", c.JWT.AccessTokenTTL),
			slog.Duration("refresh_token_ttl", c.JWT.RefreshTokenTTL),
		),

		slog.Group("postgres",
			slog.String("address", c.PG.Host+":"+c.PG.Port),
			slog.String("ssl_mode", c.PG.SSLMode),
			slog.Duration("connect_timeout", c.PG.ConnectTimeout),
			slog.Int("max_conns", c.PG.MaxConns),
            slog.Int("min_conns", c.PG.MinConns),
            slog.Duration("max_conn_lifetime", c.PG.MaxConnLifeTime),
            slog.Duration("max_conn_idletime", c.PG.MaxConnIdleTime),
		),

		slog.Group("redis",
			slog.String("address", c.Redis.Addr),
			slog.Int("pool_size", c.Redis.PoolSize),
			slog.Int("min_edle_conns", c.Redis.MinIdleConns),
			slog.Duration("dial_timeout", c.Redis.DialTimeout),
			slog.Duration("read_timeout", c.Redis.ReadTimeout),
			slog.Duration("write_timeout", c.Redis.WriteTimeout),
		),
	)
}
