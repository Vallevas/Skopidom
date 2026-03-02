// Package config loads and validates application configuration from the environment.
package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"time"
)

// Config holds all runtime configuration for the inventory server.
type Config struct {
	Server   ServerConfig
	Postgres PostgresConfig
	JWT      JWTConfig
	Storage  StorageConfig
}

// ServerConfig contains HTTP server settings.
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// PostgresConfig contains database connection settings.
type PostgresConfig struct {
	DSN            string
	MigrationsPath string
}

// JWTConfig contains token signing settings.
type JWTConfig struct {
	Secret string
	TTL    time.Duration
}

// StorageConfig contains local file storage settings.
type StorageConfig struct {
	// Dir is the filesystem directory where uploaded photos are stored.
	Dir string
	// BaseURL is the URL prefix used to construct public photo URLs.
	BaseURL string
}

// Load reads configuration from environment variables with sensible defaults.
func Load() (*Config, error) {
	_ = godotenv.Load()

	jwtTTL, err := parseDurationEnv("JWT_TTL", 24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("config: JWT_TTL: %w", err)
	}

	cfg := &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
		},
		Postgres: PostgresConfig{
			DSN:            requireEnv("DATABASE_URL"),
			MigrationsPath: getEnv("MIGRATIONS_PATH", "internal/infrastructure/postgres/migrations"),
		},
		JWT: JWTConfig{
			Secret: requireEnv("JWT_SECRET"),
			TTL:    jwtTTL,
		},
		Storage: StorageConfig{
			Dir:     getEnv("STORAGE_DIR", "./uploads"),
			BaseURL: getEnv("STORAGE_BASE_URL", "http://localhost:8080/static"),
		},
	}

	return cfg, nil
}

// getEnv returns the value of the environment variable or the given default.
func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// requireEnv returns the value of the environment variable or panics.
func requireEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("required environment variable %q is not set", key))
	}
	return val
}

// parseDurationEnv parses a duration string env var or returns the default.
func parseDurationEnv(key string, defaultVal time.Duration) (time.Duration, error) {
	raw := os.Getenv(key)
	if raw == "" {
		return defaultVal, nil
	}

	// Accept plain integer as hours for convenience (e.g. "48" → 48h).
	if hours, err := strconv.Atoi(raw); err == nil {
		return time.Duration(hours) * time.Hour, nil
	}

	return time.ParseDuration(raw)
}
