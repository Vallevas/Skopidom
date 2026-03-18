// Package config loads and validates application configuration from the environment.
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

const (
	// minJWTSecretLen is the minimum acceptable JWT secret length for HS256.
	minJWTSecretLen = 32
)

// Config holds all runtime configuration for the inventory server.
type Config struct {
	// Debug enables verbose error responses when true.
	// Should be false in production.
	Debug    bool
	Server   ServerConfig
	Postgres PostgresConfig
	JWT      JWTConfig
	Storage  StorageConfig
}

// IsDevelopment reports whether debug mode is enabled.
func (c *Config) IsDevelopment() bool {
	return c.Debug
}

// ServerConfig contains HTTP server settings.
type ServerConfig struct {
	Port           string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	AllowedOrigins []string
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

	// Default true — safer for new developers who forget to set the variable.
	debug, err := parseBoolEnv("DEBUG", true)
	if err != nil {
		return nil, fmt.Errorf("config: DEBUG: %w", err)
	}

	jwtTTL, err := parseDurationEnv("JWT_TTL", 24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("config: JWT_TTL: %w", err)
	}

	jwtSecret := requireEnv("JWT_SECRET")
	if len(jwtSecret) < minJWTSecretLen {
		return nil, fmt.Errorf(
			"config: JWT_SECRET must be at least %d characters, got %d",
			minJWTSecretLen, len(jwtSecret),
		)
	}

	allowedOrigins := splitTrimmed(getEnv("ALLOWED_ORIGINS", "*"), ",")

	cfg := &Config{
		Debug: debug,
		Server: ServerConfig{
			Port:           getEnv("SERVER_PORT", "8080"),
			ReadTimeout:    15 * time.Second,
			WriteTimeout:   15 * time.Second,
			AllowedOrigins: allowedOrigins,
		},
		Postgres: PostgresConfig{
			DSN:            requireEnv("DATABASE_URL"),
			MigrationsPath: getEnv("MIGRATIONS_PATH", "internal/infrastructure/postgres/migrations"),
		},
		JWT: JWTConfig{
			Secret: jwtSecret,
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

// parseBoolEnv parses a boolean environment variable.
// Accepts "1", "t", "T", "TRUE", "true", "True",
// "0", "f", "F", "FALSE", "false", "False"
// Returns defaultVal if the variable is not set.
func parseBoolEnv(key string, defaultVal bool) (bool, error) {
	raw := os.Getenv(key)
	if raw == "" {
		return defaultVal, nil
	}
	val, err := strconv.ParseBool(raw)
	if err != nil {
		return false, fmt.Errorf(
			"%q is not a valid boolean (use True/False or 1/0), got %q", key, raw,
		)
	}
	return val, nil
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

// splitTrimmed splits s by sep and trims whitespace from each element.
func splitTrimmed(s, sep string) []string {
	parts := strings.Split(s, sep)
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

