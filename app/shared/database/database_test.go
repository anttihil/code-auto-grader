package database

import (
	"context"
	"testing"
)

// TestSetupTestDB verifies that the test database setup works correctly
func TestSetupTestDB(t *testing.T) {
	db := SetupTestDB(t)

	// Verify we can query the database
	var count int
	err := db.Pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query users table: %v", err)
	}

	// Database should be empty initially
	if count != 0 {
		t.Errorf("Expected 0 users, got %d", count)
	}
}

// TestSeedUser verifies the SeedUser helper works
func TestSeedUser(t *testing.T) {
	db := SetupTestDB(t)

	// Seed a test user
	userID := db.SeedUser(t, "testuser", "test@example.com")

	// Verify user was created
	var username string
	err := db.Pool.QueryRow(
		context.Background(),
		"SELECT username FROM users WHERE id = $1",
		userID,
	).Scan(&username)

	if err != nil {
		t.Fatalf("Failed to query user: %v", err)
	}

	if username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", username)
	}
}

// TestTruncateTables verifies that truncation works
func TestTruncateTables(t *testing.T) {
	db := SetupTestDB(t)

	// Seed some data
	db.SeedUser(t, "user1", "user1@example.com")
	db.SeedUser(t, "user2", "user2@example.com")

	// Verify data exists
	var count int
	err := db.Pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query users: %v", err)
	}
	if count != 2 {
		t.Errorf("Expected 2 users, got %d", count)
	}

	// Truncate tables
	db.TruncateTables(t)

	// Verify data was deleted
	err = db.Pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query users after truncate: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected 0 users after truncate, got %d", count)
	}
}

// TestSeedTechTree verifies tech tree seeding
func TestSeedTechTree(t *testing.T) {
	db := SetupTestDB(t)

	treeID := db.SeedTechTree(t, "Rust Fundamentals", "rust-fundamentals", "rust")

	// Verify tech tree was created
	var name, slug string
	err := db.Pool.QueryRow(
		context.Background(),
		"SELECT name, slug FROM tech_trees WHERE id = $1",
		treeID,
	).Scan(&name, &slug)

	if err != nil {
		t.Fatalf("Failed to query tech tree: %v", err)
	}

	if name != "Rust Fundamentals" {
		t.Errorf("Expected name 'Rust Fundamentals', got '%s'", name)
	}
	if slug != "rust-fundamentals" {
		t.Errorf("Expected slug 'rust-fundamentals', got '%s'", slug)
	}
}

// TestSeedNode verifies node seeding
func TestSeedNode(t *testing.T) {
	db := SetupTestDB(t)

	// First create a tech tree
	treeID := db.SeedTechTree(t, "Rust Fundamentals", "rust-fundamentals", "rust")

	// Then create a node
	nodeID := db.SeedNode(t, treeID, "Rust Basics", "rust-basics", 1)

	// Verify node was created
	var name, slug string
	var orderIndex int
	err := db.Pool.QueryRow(
		context.Background(),
		"SELECT name, slug, order_index FROM nodes WHERE id = $1",
		nodeID,
	).Scan(&name, &slug, &orderIndex)

	if err != nil {
		t.Fatalf("Failed to query node: %v", err)
	}

	if name != "Rust Basics" {
		t.Errorf("Expected name 'Rust Basics', got '%s'", name)
	}
	if orderIndex != 1 {
		t.Errorf("Expected order_index 1, got %d", orderIndex)
	}
}

// TestSubtestsWithTruncate demonstrates using TruncateTables between subtests
func TestSubtestsWithTruncate(t *testing.T) {
	db := SetupTestDB(t)

	t.Run("FirstTest", func(t *testing.T) {
		db.TruncateTables(t) // Fresh state
		userID := db.SeedUser(t, "user1", "user1@example.com")

		var count int
		err := db.Pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM users").Scan(&count)
		if err != nil {
			t.Fatalf("Query failed: %v", err)
		}
		if count != 1 {
			t.Errorf("Expected 1 user, got %d", count)
		}

		t.Logf("Created user with ID: %s", userID)
	})

	t.Run("SecondTest", func(t *testing.T) {
		db.TruncateTables(t) // Fresh state again
		userID := db.SeedUser(t, "user2", "user2@example.com")

		var count int
		err := db.Pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM users").Scan(&count)
		if err != nil {
			t.Fatalf("Query failed: %v", err)
		}
		// Should still be 1, not 2, because we truncated
		if count != 1 {
			t.Errorf("Expected 1 user, got %d", count)
		}

		t.Logf("Created user with ID: %s", userID)
	})
}
