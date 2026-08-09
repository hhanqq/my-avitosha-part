package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"
)

type LogLevel string

const (
	AppEnvDev  = "dev"
	AppEnvTest = "test"
	AppEnvProd = "prod"

	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"

	defaultAccessTokenTTL  = 15 * time.Minute
	defaultRefreshTokenTTL = 30 * 24 * time.Hour
	defaultProxyAPIBaseURL = "https://api.proxyapi.ru/openrouter/v1"
	defaultProxyAPIModel   = "qwen/qwen-2.5-7b-instruct"
	defaultProxyAPITimeout = 4 * time.Second
)

type Config struct {
	AppEnv           string
	HTTPAddr         string
	DatabaseURL      string
	FrontendOrigin   string
	JWTSigningKey    string
	JWTIssuer        string
	JWTAudience      string
	LogLevel         LogLevel
	ShutdownTimeout  time.Duration
	HTTPReadTimeout  time.Duration
	HTTPWriteTimeout time.Duration
	HTTPIdleTimeout  time.Duration
	AccessTokenTTL   time.Duration
	RefreshTokenTTL  time.Duration
	ProxyAPIKey      string
	ProxyAPIBaseURL  string
	ProxyAPIModel    string
	ProxyAPITimeout  time.Duration
	GRPCAddr         string
	AuthGRPCAddr     string
	GameGRPCAddr     string
}

func Load() (Config, error) {
	return LoadFromEnv(os.Getenv)
}

func LoadGateway() (Config, error) {
	return loadForRole(os.Getenv, Config.ValidateGateway)
}

func LoadAuthService() (Config, error) {
	return loadForRole(os.Getenv, Config.ValidateAuthService)
}

func LoadGameService() (Config, error) {
	return loadForRole(os.Getenv, Config.ValidateGameService)
}

func LoadFromEnv(getenv func(string) string) (Config, error) {
	return loadForRole(getenv, Config.Validate)
}

func loadForRole(getenv func(string) string, validate func(Config) error) (Config, error) {
	cfg := Config{
		AppEnv:           envOrDefault(getenv, "APP_ENV", AppEnvDev),
		HTTPAddr:         strings.TrimSpace(getenv("HTTP_ADDR")),
		DatabaseURL:      strings.TrimSpace(getenv("DATABASE_URL")),
		FrontendOrigin:   strings.TrimSpace(getenv("FRONTEND_ORIGIN")),
		JWTSigningKey:    getenv("JWT_SIGNING_KEY"),
		JWTIssuer:        strings.TrimSpace(getenv("JWT_ISSUER")),
		JWTAudience:      strings.TrimSpace(getenv("JWT_AUDIENCE")),
		LogLevel:         LogLevel(envOrDefault(getenv, "LOG_LEVEL", string(LogLevelInfo))),
		ShutdownTimeout:  5 * time.Second,
		HTTPReadTimeout:  5 * time.Second,
		HTTPWriteTimeout: 10 * time.Second,
		HTTPIdleTimeout:  60 * time.Second,
		AccessTokenTTL:   defaultAccessTokenTTL,
		RefreshTokenTTL:  defaultRefreshTokenTTL,
		ProxyAPIKey:      strings.TrimSpace(getenv("PROXYAPI_API_KEY")),
		ProxyAPIBaseURL:  envOrDefault(getenv, "PROXYAPI_BASE_URL", defaultProxyAPIBaseURL),
		ProxyAPIModel:    envOrDefault(getenv, "PROXYAPI_MODEL", defaultProxyAPIModel),
		ProxyAPITimeout:  defaultProxyAPITimeout,
		GRPCAddr:         envOrDefault(getenv, "GRPC_ADDR", ":9090"),
		AuthGRPCAddr:     envOrDefault(getenv, "AUTH_GRPC_ADDR", "127.0.0.1:9091"),
		GameGRPCAddr:     envOrDefault(getenv, "GAME_GRPC_ADDR", "127.0.0.1:9092"),
	}

	if value := strings.TrimSpace(getenv("SHUTDOWN_TIMEOUT")); value != "" {
		timeout, err := time.ParseDuration(value)
		if err != nil {
			return Config{}, fmt.Errorf("SHUTDOWN_TIMEOUT must be a duration: %w", err)
		}
		cfg.ShutdownTimeout = timeout
	}
	if value := strings.TrimSpace(getenv("ACCESS_TOKEN_TTL")); value != "" {
		ttl, err := time.ParseDuration(value)
		if err != nil {
			return Config{}, fmt.Errorf("ACCESS_TOKEN_TTL must be a duration: %w", err)
		}
		cfg.AccessTokenTTL = ttl
	}
	if value := strings.TrimSpace(getenv("REFRESH_TOKEN_TTL")); value != "" {
		ttl, err := time.ParseDuration(value)
		if err != nil {
			return Config{}, fmt.Errorf("REFRESH_TOKEN_TTL must be a duration: %w", err)
		}
		cfg.RefreshTokenTTL = ttl
	}
	if value := strings.TrimSpace(getenv("PROXYAPI_TIMEOUT")); value != "" {
		timeout, err := time.ParseDuration(value)
		if err != nil {
			return Config{}, fmt.Errorf("PROXYAPI_TIMEOUT must be a duration: %w", err)
		}
		cfg.ProxyAPITimeout = timeout
	}

	if err := validate(cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (cfg Config) ValidateGateway() error {
	if err := cfg.validateRuntime(); err != nil {
		return err
	}
	if cfg.HTTPAddr == "" {
		return fmt.Errorf("HTTP_ADDR is required")
	}
	if cfg.FrontendOrigin == "" {
		return fmt.Errorf("FRONTEND_ORIGIN is required")
	}
	if cfg.RefreshTokenTTL <= 0 {
		return fmt.Errorf("REFRESH_TOKEN_TTL must be positive")
	}
	if cfg.AuthGRPCAddr == "" || cfg.GameGRPCAddr == "" {
		return fmt.Errorf("AUTH_GRPC_ADDR and GAME_GRPC_ADDR must not be empty")
	}
	return nil
}

func (cfg Config) ValidateAuthService() error {
	if err := cfg.validateRuntime(); err != nil {
		return err
	}
	if cfg.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.JWTSigningKey == "" {
		return fmt.Errorf("JWT_SIGNING_KEY is required")
	}
	if cfg.JWTIssuer == "" {
		return fmt.Errorf("JWT_ISSUER is required")
	}
	if cfg.JWTAudience == "" {
		return fmt.Errorf("JWT_AUDIENCE is required")
	}
	if cfg.AccessTokenTTL <= 0 || cfg.RefreshTokenTTL <= 0 {
		return fmt.Errorf("ACCESS_TOKEN_TTL and REFRESH_TOKEN_TTL must be positive")
	}
	if cfg.GRPCAddr == "" {
		return fmt.Errorf("GRPC_ADDR must not be empty")
	}
	return nil
}

func (cfg Config) ValidateGameService() error {
	if err := cfg.validateRuntime(); err != nil {
		return err
	}
	if cfg.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.GRPCAddr == "" {
		return fmt.Errorf("GRPC_ADDR must not be empty")
	}
	if cfg.ProxyAPIBaseURL == "" || cfg.ProxyAPIModel == "" || cfg.ProxyAPITimeout <= 0 {
		return fmt.Errorf("ProxyAPI URL/model must not be empty and timeout must be positive")
	}
	return nil
}

func (cfg Config) validateRuntime() error {
	switch cfg.AppEnv {
	case AppEnvDev, AppEnvTest, AppEnvProd:
	default:
		return fmt.Errorf("APP_ENV must be one of: dev, test, prod")
	}
	switch cfg.LogLevel {
	case LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError:
	default:
		return fmt.Errorf("LOG_LEVEL must be one of: debug, info, warn, error")
	}
	if cfg.ShutdownTimeout <= 0 {
		return fmt.Errorf("SHUTDOWN_TIMEOUT must be positive")
	}
	return nil
}

func (cfg Config) Validate() error {
	if cfg.HTTPAddr == "" {
		return fmt.Errorf("HTTP_ADDR is required")
	}
	if cfg.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.FrontendOrigin == "" {
		return fmt.Errorf("FRONTEND_ORIGIN is required")
	}
	if cfg.JWTSigningKey == "" {
		return fmt.Errorf("JWT_SIGNING_KEY is required")
	}
	if cfg.JWTIssuer == "" {
		return fmt.Errorf("JWT_ISSUER is required")
	}
	if cfg.JWTAudience == "" {
		return fmt.Errorf("JWT_AUDIENCE is required")
	}

	switch cfg.AppEnv {
	case AppEnvDev, AppEnvTest, AppEnvProd:
	default:
		return fmt.Errorf("APP_ENV must be one of: dev, test, prod")
	}

	switch cfg.LogLevel {
	case LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError:
	default:
		return fmt.Errorf("LOG_LEVEL must be one of: debug, info, warn, error")
	}

	if cfg.ShutdownTimeout <= 0 {
		return fmt.Errorf("SHUTDOWN_TIMEOUT must be positive")
	}
	if cfg.AccessTokenTTL <= 0 {
		return fmt.Errorf("ACCESS_TOKEN_TTL must be positive")
	}
	if cfg.RefreshTokenTTL <= 0 {
		return fmt.Errorf("REFRESH_TOKEN_TTL must be positive")
	}
	if cfg.ProxyAPIBaseURL == "" {
		return fmt.Errorf("PROXYAPI_BASE_URL must not be empty")
	}
	if cfg.ProxyAPIModel == "" {
		return fmt.Errorf("PROXYAPI_MODEL must not be empty")
	}
	if cfg.ProxyAPITimeout <= 0 {
		return fmt.Errorf("PROXYAPI_TIMEOUT must be positive")
	}
	if cfg.GRPCAddr == "" || cfg.AuthGRPCAddr == "" || cfg.GameGRPCAddr == "" {
		return fmt.Errorf("GRPC_ADDR, AUTH_GRPC_ADDR and GAME_GRPC_ADDR must not be empty")
	}

	return nil
}

func (level LogLevel) Level() slog.Level {
	switch level {
	case LogLevelDebug:
		return slog.LevelDebug
	case LogLevelWarn:
		return slog.LevelWarn
	case LogLevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func envOrDefault(getenv func(string) string, key, fallback string) string {
	value := strings.TrimSpace(getenv(key))
	if value == "" {
		return fallback
	}
	return value
}
