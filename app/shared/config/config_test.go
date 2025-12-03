package config

import (
	"os"
	"testing"
	"time"
)

// setMinimalDBConfig sets the minimum required database config for tests
func setMinimalDBConfig() {
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("DB_NAME", "testdb")
}

// clearDBConfig clears all database config
func clearDBConfig() {
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("DB_SSLMODE")
}

func TestLoad_WithAtomicParts(t *testing.T) {
	// Set up atomic database parts
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("DB_HOST", "testhost")
	os.Setenv("DB_PORT", "5555")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_SSLMODE", "require")
	defer clearDBConfig()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	expected := "postgres://testuser:testpass@testhost:5555/testdb?sslmode=require"
	if cfg.DatabaseURL != expected {
		t.Errorf("DatabaseURL = %v, want %v", cfg.DatabaseURL, expected)
	}
}

func TestLoad_WithDefaultHostAndPort(t *testing.T) {
	// Set only required parts, rely on defaults for host/port
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("DB_NAME", "testdb")
	defer clearDBConfig()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	expected := "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable"
	if cfg.DatabaseURL != expected {
		t.Errorf("DatabaseURL = %v, want %v", cfg.DatabaseURL, expected)
	}
}

func TestLoad_ValidationError(t *testing.T) {
	// Clear all database config
	clearDBConfig()

	_, err := Load()
	if err == nil {
		t.Error("Load() expected error when no database config provided, got nil")
	}
}

func TestLoad_ValidationError_MissingUser(t *testing.T) {
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("DB_NAME", "testdb")
	defer clearDBConfig()

	_, err := Load()
	if err == nil {
		t.Error("Load() expected error when DB_USER missing, got nil")
	}
}

func TestLoad_ValidationError_MissingName(t *testing.T) {
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	defer clearDBConfig()

	_, err := Load()
	if err == nil {
		t.Error("Load() expected error when DB_NAME missing, got nil")
	}
}

func TestLoad_ServerDefaults(t *testing.T) {
	setMinimalDBConfig()
	os.Unsetenv("PORT")
	os.Unsetenv("ENVIRONMENT")
	defer clearDBConfig()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Port != "5000" {
		t.Errorf("Port = %v, want 5000", cfg.Port)
	}
	if cfg.Environment != "development" {
		t.Errorf("Environment = %v, want development", cfg.Environment)
	}
}

func TestLoad_ExecutionConfig(t *testing.T) {
	setMinimalDBConfig()
	os.Setenv("MAX_CONCURRENT_EXECUTIONS", "25")
	os.Setenv("EXECUTION_QUEUE_TIMEOUT", "30s")
	os.Setenv("JS_POOL_SIZE", "10")
	os.Setenv("PY_POOL_SIZE", "8")
	os.Setenv("RUST_POOL_SIZE", "4")
	defer func() {
		clearDBConfig()
		os.Unsetenv("MAX_CONCURRENT_EXECUTIONS")
		os.Unsetenv("EXECUTION_QUEUE_TIMEOUT")
		os.Unsetenv("JS_POOL_SIZE")
		os.Unsetenv("PY_POOL_SIZE")
		os.Unsetenv("RUST_POOL_SIZE")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.MaxConcurrentExecutions != 25 {
		t.Errorf("MaxConcurrentExecutions = %v, want 25", cfg.MaxConcurrentExecutions)
	}
	if cfg.ExecutionQueueTimeout != 30*time.Second {
		t.Errorf("ExecutionQueueTimeout = %v, want 30s", cfg.ExecutionQueueTimeout)
	}
	if cfg.JSPoolSize != 10 {
		t.Errorf("JSPoolSize = %v, want 10", cfg.JSPoolSize)
	}
	if cfg.PythonPoolSize != 8 {
		t.Errorf("PythonPoolSize = %v, want 8", cfg.PythonPoolSize)
	}
	if cfg.RustPoolSize != 4 {
		t.Errorf("RustPoolSize = %v, want 4", cfg.RustPoolSize)
	}
}

func TestLoad_ExecutionConfigDefaults(t *testing.T) {
	setMinimalDBConfig()
	os.Unsetenv("MAX_CONCURRENT_EXECUTIONS")
	os.Unsetenv("EXECUTION_QUEUE_TIMEOUT")
	os.Unsetenv("JS_POOL_SIZE")
	os.Unsetenv("PY_POOL_SIZE")
	os.Unsetenv("RUST_POOL_SIZE")
	defer clearDBConfig()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.MaxConcurrentExecutions != 15 {
		t.Errorf("MaxConcurrentExecutions = %v, want 15", cfg.MaxConcurrentExecutions)
	}
	if cfg.ExecutionQueueTimeout != 10*time.Second {
		t.Errorf("ExecutionQueueTimeout = %v, want 10s", cfg.ExecutionQueueTimeout)
	}
	if cfg.JSPoolSize != 5 {
		t.Errorf("JSPoolSize = %v, want 5", cfg.JSPoolSize)
	}
	if cfg.PythonPoolSize != 5 {
		t.Errorf("PythonPoolSize = %v, want 5", cfg.PythonPoolSize)
	}
	if cfg.RustPoolSize != 2 {
		t.Errorf("RustPoolSize = %v, want 2", cfg.RustPoolSize)
	}
}

func TestConfig_IsDevelopment(t *testing.T) {
	tests := []struct {
		env      string
		expected bool
	}{
		{"development", true},
		{"production", false},
		{"staging", false},
	}

	for _, tt := range tests {
		cfg := &Config{Environment: tt.env}
		if got := cfg.IsDevelopment(); got != tt.expected {
			t.Errorf("IsDevelopment() with env=%q = %v, want %v", tt.env, got, tt.expected)
		}
	}
}

func TestConfig_IsProduction(t *testing.T) {
	tests := []struct {
		env      string
		expected bool
	}{
		{"development", false},
		{"production", true},
		{"staging", false},
	}

	for _, tt := range tests {
		cfg := &Config{Environment: tt.env}
		if got := cfg.IsProduction(); got != tt.expected {
			t.Errorf("IsProduction() with env=%q = %v, want %v", tt.env, got, tt.expected)
		}
	}
}

func TestLoad_TestDatabaseConfig(t *testing.T) {
	setMinimalDBConfig()
	os.Setenv("TEST_DB_USER", "testdbuser")
	os.Setenv("TEST_DB_PASSWORD", "testdbpass")
	os.Setenv("TEST_DB_HOST", "testdbhost")
	os.Setenv("TEST_DB_PORT", "5555")
	os.Setenv("TEST_DB_NAME", "testdbname")
	defer func() {
		clearDBConfig()
		os.Unsetenv("TEST_DB_USER")
		os.Unsetenv("TEST_DB_PASSWORD")
		os.Unsetenv("TEST_DB_HOST")
		os.Unsetenv("TEST_DB_PORT")
		os.Unsetenv("TEST_DB_NAME")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	expected := "postgres://testdbuser:testdbpass@testdbhost:5555/testdbname?sslmode=disable"
	if cfg.TestDatabaseURL != expected {
		t.Errorf("TestDatabaseURL = %v, want %v", cfg.TestDatabaseURL, expected)
	}
}

func TestLoad_DBAutoMigrate(t *testing.T) {
	setMinimalDBConfig()
	defer clearDBConfig()

	tests := []struct {
		envValue string
		expected bool
	}{
		{"", true},         // default
		{"true", true},     // explicit true
		{"false", false},   // explicit false
		{"1", true},        // truthy
		{"yes", true},      // truthy
		{"no", false},      // falsy
		{"invalid", false}, // invalid treated as false
	}

	for _, tt := range tests {
		if tt.envValue == "" {
			os.Unsetenv("DB_AUTO_MIGRATE")
		} else {
			os.Setenv("DB_AUTO_MIGRATE", tt.envValue)
		}

		cfg, err := Load()
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}

		if cfg.DBAutoMigrate != tt.expected {
			t.Errorf("DBAutoMigrate with env=%q = %v, want %v", tt.envValue, cfg.DBAutoMigrate, tt.expected)
		}
	}
	os.Unsetenv("DB_AUTO_MIGRATE")
}
