package database

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TestDatabase holds a connection pool for testing
type TestDatabase struct {
	Pool *pgxpool.Pool
	URL  string
}

// SetupTestDB creates a fresh test database with migrations
// Call this in your test setup to get a clean database state
//
// Example:
//
//	func TestMyHandler(t *testing.T) {
//	    db := database.SetupTestDB(t)
//	    // Use db.Pool for queries
//	}
func SetupTestDB(t *testing.T) *TestDatabase {
	t.Helper()

	// Get test database URL from environment
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://code_grader:test_password@localhost:5433/code_auto_grader_test?sslmode=disable"
	}

	// Run migrations
	m, err := migrate.New(
		"file://../migrations",
		dbURL,
	)
	if err != nil {
		t.Fatalf("Failed to create migrate instance: %v", err)
	}

	// Reset to clean state: first down to version 0, then up
	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		t.Fatalf("Failed to get migration version: %v", err)
	}

	// If migrations exist and database is dirty, force to current version
	if dirty {
		t.Logf("Database is in dirty state, forcing version %d", version)
		if err := m.Force(int(version)); err != nil {
			t.Fatalf("Failed to force migration version: %v", err)
		}
	}

	// Migrate down to clean state (if migrations exist)
	if err == nil || err == migrate.ErrNilVersion {
		if downErr := m.Down(); downErr != nil && downErr != migrate.ErrNoChange {
			t.Logf("Warning: Could not migrate down: %v", downErr)
		}
	}

	// Run all migrations (now only schema, no seed data)
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Close migration connection
	srcErr, dbErr := m.Close()
	if srcErr != nil {
		t.Logf("Warning: Error closing migration source: %v", srcErr)
	}
	if dbErr != nil {
		t.Logf("Warning: Error closing migration database: %v", dbErr)
	}

	// Create connection pool
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Verify connection
	if err := pool.Ping(context.Background()); err != nil {
		t.Fatalf("Failed to ping test database: %v", err)
	}

	// Cleanup when test finishes
	t.Cleanup(func() {
		pool.Close()
	})

	return &TestDatabase{
		Pool: pool,
		URL:  dbURL,
	}
}

// TruncateTables removes all data but keeps schema
// Useful for resetting state between subtests
//
// Example:
//
//	t.Run("TestCase1", func(t *testing.T) {
//	    db.TruncateTables(t) // Fresh state
//	    // Test logic
//	})
func (db *TestDatabase) TruncateTables(t *testing.T) {
	t.Helper()

	// Order matters: delete child tables first (due to foreign keys)
	tables := []string{
		"user_progress",
		"submissions",
		"node_prerequisites",
		"exercises",
		"nodes",
		"tech_trees",
		"users",
	}

	ctx := context.Background()
	for _, table := range tables {
		_, err := db.Pool.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			t.Fatalf("Failed to truncate table %s: %v", table, err)
		}
	}
}

// SeedUser creates a test user and returns the UUID
// Convenience method for tests that need a user
func (db *TestDatabase) SeedUser(t *testing.T, username, email string) string {
	t.Helper()

	var userID string
	err := db.Pool.QueryRow(
		context.Background(),
		`INSERT INTO users (username, email, display_name)
		 VALUES ($1, $2, $3)
		 RETURNING id`,
		username,
		email,
		username,
	).Scan(&userID)

	if err != nil {
		t.Fatalf("Failed to seed user: %v", err)
	}

	return userID
}

// SeedTechTree creates a test tech tree and returns the UUID
func (db *TestDatabase) SeedTechTree(t *testing.T, name, slug, language string) string {
	t.Helper()

	var treeID string
	err := db.Pool.QueryRow(
		context.Background(),
		`INSERT INTO tech_trees (name, slug, language, is_published)
		 VALUES ($1, $2, $3, true)
		 RETURNING id`,
		name,
		slug,
		language,
	).Scan(&treeID)

	if err != nil {
		t.Fatalf("Failed to seed tech tree: %v", err)
	}

	return treeID
}

// SeedNode creates a test node and returns the UUID
func (db *TestDatabase) SeedNode(t *testing.T, treeID, name, slug string, orderIndex int) string {
	t.Helper()

	var nodeID string
	err := db.Pool.QueryRow(
		context.Background(),
		`INSERT INTO nodes (tech_tree_id, name, slug, order_index, description)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id`,
		treeID,
		name,
		slug,
		orderIndex,
		"Test node description",
	).Scan(&nodeID)

	if err != nil {
		t.Fatalf("Failed to seed node: %v", err)
	}

	return nodeID
}
