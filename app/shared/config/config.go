// Package config provides centralized configuration management for the application.
// It loads configuration from environment variables and validates required settings.
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all application configuration loaded from environment variables.
// It serves as the single source of truth for configuration across the application.
type Config struct {
	// Server settings
	Port        string
	Environment string // "development" or "production"

	// Database settings (atomic parts for flexibility)
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
	DBSSLMode  string

	// Computed database URL (composed from atomic parts)
	DatabaseURL string

	// Test database settings
	TestDBUser      string
	TestDBPassword  string
	TestDBHost      string
	TestDBPort      string
	TestDBName      string
	TestDatabaseURL string

	// Database behavior
	DBAutoMigrate    bool
	DBMigrationsPath string

	// External services
	GeminiAPIKey string

	// Code execution settings
	MaxConcurrentExecutions int
	ExecutionQueueTimeout   time.Duration
	JSPoolSize              int
	PythonPoolSize          int
	RustPoolSize            int
}

// Load reads configuration from environment variables.
// Database configuration requires atomic parts: DB_USER, DB_PASSWORD, DB_NAME, etc.
func Load() (*Config, error) {
	cfg := &Config{
		// Server defaults
		Port:        getEnvOrDefault("PORT", "5000"),
		Environment: getEnvOrDefault("ENVIRONMENT", "development"),

		// Database atomic parts
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBHost:     getEnvOrDefault("DB_HOST", "localhost"),
		DBPort:     getEnvOrDefault("DB_PORT", "5442"),
		DBName:     os.Getenv("DB_NAME"),
		DBSSLMode:  getEnvOrDefault("DB_SSLMODE", "disable"),

		// Test database settings
		TestDBUser:     os.Getenv("TEST_DB_USER"),
		TestDBPassword: os.Getenv("TEST_DB_PASSWORD"),
		TestDBHost:     getEnvOrDefault("TEST_DB_HOST", "localhost"),
		TestDBPort:     getEnvOrDefault("TEST_DB_PORT", "5443"),
		TestDBName:     os.Getenv("TEST_DB_NAME"),

		// Database behavior
		DBAutoMigrate:    getEnvBool("DB_AUTO_MIGRATE", true),
		DBMigrationsPath: getEnvOrDefault("DB_MIGRATIONS_PATH", "file://migrations"),

		// Gemini
		GeminiAPIKey: os.Getenv("GEMINI_API_KEY"),

		// Code execution
		MaxConcurrentExecutions: getEnvInt("MAX_CONCURRENT_EXECUTIONS", 15),
		ExecutionQueueTimeout:   getEnvDuration("EXECUTION_QUEUE_TIMEOUT", 10*time.Second),
		JSPoolSize:              getEnvInt("JS_POOL_SIZE", 5),
		PythonPoolSize:          getEnvInt("PY_POOL_SIZE", 5),
		RustPoolSize:            getEnvInt("RUST_POOL_SIZE", 2),
	}

	// Compose DATABASE_URL from atomic parts if available
	if cfg.DBUser != "" && cfg.DBName != "" {
		cfg.DatabaseURL = buildDatabaseURL(
			cfg.DBUser, cfg.DBPassword, cfg.DBHost,
			cfg.DBPort, cfg.DBName, cfg.DBSSLMode,
		)
	}

	// Compose TEST_DATABASE_URL from atomic parts if available
	if cfg.TestDBUser != "" && cfg.TestDBName != "" {
		cfg.TestDatabaseURL = buildDatabaseURL(
			cfg.TestDBUser, cfg.TestDBPassword, cfg.TestDBHost,
			cfg.TestDBPort, cfg.TestDBName, cfg.DBSSLMode,
		)
	}

	return cfg, cfg.Validate()
}

// MustLoad loads configuration and panics on error.
// Use this in main() when configuration is required to start.
func MustLoad() *Config {
	cfg, err := Load()
	if err != nil {
		panic(fmt.Sprintf("configuration error: %v", err))
	}
	return cfg
}

// Validate checks that all required configuration is present.
// Returns an error describing all missing required fields.
func (c *Config) Validate() error {
	var missing []string

	if c.DatabaseURL == "" {
		missing = append(missing, "database config (set DB_USER, DB_PASSWORD, and DB_NAME)")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required configuration: %s", strings.Join(missing, ", "))
	}

	return nil
}

// IsDevelopment returns true if running in development mode.
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if running in production mode.
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// buildDatabaseURL constructs a PostgreSQL connection URL from parts.
func buildDatabaseURL(user, password, host, port, dbname, sslmode string) string {
	if user == "" || dbname == "" {
		return ""
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbname, sslmode)
}

// getEnvOrDefault returns the environment variable value or a default.
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt parses an integer from environment variable with a default.
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil && intVal > 0 {
			return intVal
		}
	}
	return defaultValue
}

// getEnvBool parses a boolean from environment variable with a default.
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true" || value == "1" || value == "yes"
	}
	return defaultValue
}

// getEnvDuration parses a duration from environment variable with a default.
func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
