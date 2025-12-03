package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Database holds the connection pool and configuration
type Database struct {
	Pool *pgxpool.Pool
	URL  string
}

// Config holds database configuration options
type Config struct {
	URL                string
	MaxConns           int32
	MinConns           int32
	MaxConnLifetime    time.Duration
	MaxConnIdleTime    time.Duration
	HealthCheckPeriod  time.Duration
	AutoMigrate        bool // Run migrations on startup
	MigrationsPath     string
}

// DefaultConfig returns sensible defaults for database configuration
func DefaultConfig() *Config {
	return &Config{
		MaxConns:          25,
		MinConns:          5,
		MaxConnLifetime:   time.Hour,
		MaxConnIdleTime:   time.Minute * 30,
		HealthCheckPeriod: time.Minute,
		AutoMigrate:       true,
		MigrationsPath:    "file://migrations",
	}
}

// New creates a new database connection with the given configuration
// It establishes a connection pool, runs migrations (if enabled), and verifies connectivity
func New(ctx context.Context, cfg *Config) (*Database, error) {
	if cfg.URL == "" {
		return nil, fmt.Errorf("database URL is required")
	}

	// Parse connection config
	poolConfig, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	// Apply connection pool settings
	poolConfig.MaxConns = cfg.MaxConns
	poolConfig.MinConns = cfg.MinConns
	poolConfig.MaxConnLifetime = cfg.MaxConnLifetime
	poolConfig.MaxConnIdleTime = cfg.MaxConnIdleTime
	poolConfig.HealthCheckPeriod = cfg.HealthCheckPeriod

	// Create connection pool
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db := &Database{
		Pool: pool,
		URL:  cfg.URL,
	}

	// Run migrations if enabled
	if cfg.AutoMigrate {
		if err := db.Migrate(cfg.MigrationsPath); err != nil {
			pool.Close()
			return nil, fmt.Errorf("failed to run migrations: %w", err)
		}
	}

	return db, nil
}

// Migrate runs database migrations using golang-migrate
func (db *Database) Migrate(migrationsPath string) error {
	m, err := migrate.New(migrationsPath, db.URL)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			log.Printf("Warning: error closing migration source: %v", srcErr)
		}
		if dbErr != nil {
			log.Printf("Warning: error closing migration database: %v", dbErr)
		}
	}()

	// Check current version
	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	// Handle dirty state
	if dirty {
		log.Printf("Warning: database is in dirty state at version %d, forcing version", version)
		if err := m.Force(int(version)); err != nil {
			return fmt.Errorf("failed to force migration version: %w", err)
		}
	}

	// Run migrations
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("Database schema is up to date")
			return nil
		}
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	// Get new version
	newVersion, _, err := m.Version()
	if err != nil {
		log.Println("Migrations applied successfully")
	} else {
		log.Printf("Migrations applied successfully (version: %d)", newVersion)
	}

	return nil
}

// Close closes the database connection pool
func (db *Database) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}

// Ping verifies the database connection is alive
func (db *Database) Ping(ctx context.Context) error {
	return db.Pool.Ping(ctx)
}

// Stats returns connection pool statistics
func (db *Database) Stats() *PoolStats {
	stats := db.Pool.Stat()
	return &PoolStats{
		AcquireCount:            stats.AcquireCount(),
		AcquireDuration:         stats.AcquireDuration(),
		AcquiredConns:           stats.AcquiredConns(),
		CanceledAcquireCount:    stats.CanceledAcquireCount(),
		ConstructingConns:       stats.ConstructingConns(),
		EmptyAcquireCount:       stats.EmptyAcquireCount(),
		IdleConns:               stats.IdleConns(),
		MaxConns:                stats.MaxConns(),
		TotalConns:              stats.TotalConns(),
		NewConnsCount:           stats.NewConnsCount(),
		MaxLifetimeDestroyCount: stats.MaxLifetimeDestroyCount(),
		MaxIdleDestroyCount:     stats.MaxIdleDestroyCount(),
	}
}

// WithTransaction executes a function within a database transaction
// If the function returns an error, the transaction is rolled back
// If the function succeeds, the transaction is committed
func (db *Database) WithTransaction(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			// Rollback on panic and re-throw
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	err = fn(tx)
	if err != nil {
		// Rollback on error
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("transaction failed: %w (rollback error: %v)", err, rbErr)
		}
		return err
	}

	// Commit on success
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// PoolStats holds connection pool statistics
type PoolStats struct {
	AcquireCount            int64
	AcquireDuration         time.Duration
	AcquiredConns           int32
	CanceledAcquireCount    int64
	ConstructingConns       int32
	EmptyAcquireCount       int64
	IdleConns               int32
	MaxConns                int32
	TotalConns              int32
	NewConnsCount           int64
	MaxLifetimeDestroyCount int64
	MaxIdleDestroyCount     int64
}

// NewFromEnv creates a database connection using environment variables.
// Deprecated: Use New() with config from shared/config package instead.
// This function is kept for backward compatibility.
func NewFromEnv(ctx context.Context) (*Database, error) {
	cfg := DefaultConfig()
	cfg.URL = os.Getenv("DATABASE_URL")

	// Allow disabling auto-migration via env var
	if os.Getenv("DB_AUTO_MIGRATE") == "false" {
		cfg.AutoMigrate = false
	}

	// Allow custom migrations path
	if path := os.Getenv("DB_MIGRATIONS_PATH"); path != "" {
		cfg.MigrationsPath = path
	}

	return New(ctx, cfg)
}