package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	AppEnv        string
	HTTPAddr      string
	DBDriver      string
	SQLitePath    string
	DatabaseURL   string
	JWTSecret     string
	AccessTTL     time.Duration
	AdminUsername string
	AdminPassword string
	StaticDir     string
}

func Load() (Config, error) {
	cfg := Config{
		AppEnv:        getEnv("APP_ENV", "development"),
		HTTPAddr:      getEnv("HTTP_ADDR", ":8080"),
		DBDriver:      strings.ToLower(getEnv("DB_DRIVER", "sqlite")),
		SQLitePath:    getEnv("DB_SQLITE_PATH", "./data/windnav.db"),
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		JWTSecret:     getEnv("JWT_SECRET", "change-me-in-production"),
		AccessTTL:     getDuration("ACCESS_TOKEN_TTL", 24*time.Hour),
		AdminUsername: getEnv("ADMIN_USERNAME", "admin"),
		AdminPassword: getEnv("ADMIN_PASSWORD", "admin123456"),
		StaticDir:     getEnv("STATIC_DIR", ""),
	}

	if cfg.DBDriver != "sqlite" && cfg.DBDriver != "postgres" {
		return Config{}, fmt.Errorf("unsupported DB_DRIVER %q", cfg.DBDriver)
	}
	if cfg.DBDriver == "postgres" && cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required when DB_DRIVER=postgres")
	}
	if cfg.AppEnv == "production" && cfg.JWTSecret == "change-me-in-production" {
		return Config{}, fmt.Errorf("JWT_SECRET must be changed in production")
	}
	return cfg, nil
}

func getEnv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func getDuration(key string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return duration
}
