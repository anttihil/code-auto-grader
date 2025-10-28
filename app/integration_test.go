package main

import (
	"context"
	"testing"

	"github.com/anttihil/code-auto-grader/database"
	"github.com/anttihil/code-auto-grader/repository"
)

// TestDatabaseIntegration verifies the complete database setup works
func TestDatabaseIntegration(t *testing.T) {
	// Set up test database with migrations
	db := database.SetupTestDB(t)

	// Verify tables exist by querying each one
	ctx := context.Background()

	t.Run("VerifyTablesExist", func(t *testing.T) {
		tables := []string{
			"users", "tech_trees", "nodes", "node_prerequisites",
			"exercises", "submissions", "user_progress",
		}

		for _, table := range tables {
			var exists bool
			query := `
				SELECT EXISTS (
					SELECT FROM information_schema.tables
					WHERE table_name = $1
				)
			`
			err := db.Pool.QueryRow(ctx, query, table).Scan(&exists)
			if err != nil {
				t.Fatalf("Failed to check if table %s exists: %v", table, err)
			}
			if !exists {
				t.Errorf("Table %s does not exist", table)
			}
		}
	})

	t.Run("UserRepository", func(t *testing.T) {
		db.TruncateTables(t)
		userRepo := repository.NewUserRepo(db.Pool)

		// Create a user
		userID := db.SeedUser(t, "testuser", "test@example.com")

		// Test GetByID
		user, err := userRepo.GetByID(ctx, userID)
		if err != nil {
			t.Fatalf("Failed to get user by ID: %v", err)
		}
		if user.Username != "testuser" {
			t.Errorf("Expected username 'testuser', got '%s'", user.Username)
		}

		// Test GetByUsername
		user2, err := userRepo.GetByUsername(ctx, "testuser")
		if err != nil {
			t.Fatalf("Failed to get user by username: %v", err)
		}
		if user2.ID.String() != userID {
			t.Errorf("Expected user ID %s, got %s", userID, user2.ID)
		}

		// Test GetByEmail
		user3, err := userRepo.GetByEmail(ctx, "test@example.com")
		if err != nil {
			t.Fatalf("Failed to get user by email: %v", err)
		}
		if user3.ID.String() != userID {
			t.Errorf("Expected user ID %s, got %s", userID, user3.ID)
		}
	})

	t.Run("TechTreeRepository", func(t *testing.T) {
		db.TruncateTables(t)
		treeRepo := repository.NewTechTreeRepo(db.Pool)

		// Create a tech tree
		treeID := db.SeedTechTree(t, "Rust Fundamentals", "rust-fundamentals", "rust")

		// Test GetByID
		tree, err := treeRepo.GetByID(ctx, treeID)
		if err != nil {
			t.Fatalf("Failed to get tech tree by ID: %v", err)
		}
		if tree.Name != "Rust Fundamentals" {
			t.Errorf("Expected name 'Rust Fundamentals', got '%s'", tree.Name)
		}

		// Test GetBySlug
		tree2, err := treeRepo.GetBySlug(ctx, "rust-fundamentals")
		if err != nil {
			t.Fatalf("Failed to get tech tree by slug: %v", err)
		}
		if tree2.ID.String() != treeID {
			t.Errorf("Expected tree ID %s, got %s", treeID, tree2.ID)
		}

		// Test List
		trees, err := treeRepo.List(ctx)
		if err != nil {
			t.Fatalf("Failed to list tech trees: %v", err)
		}
		if len(trees) != 1 {
			t.Errorf("Expected 1 tech tree, got %d", len(trees))
		}
	})

	t.Run("NodeRepository", func(t *testing.T) {
		db.TruncateTables(t)
		nodeRepo := repository.NewNodeRepo(db.Pool)

		// Create tech tree and node
		treeID := db.SeedTechTree(t, "Rust Fundamentals", "rust-fundamentals", "rust")
		nodeID := db.SeedNode(t, treeID, "Rust Basics", "rust-basics", 1)

		// Test GetByID
		node, err := nodeRepo.GetByID(ctx, nodeID)
		if err != nil {
			t.Fatalf("Failed to get node by ID: %v", err)
		}
		if node.Name != "Rust Basics" {
			t.Errorf("Expected name 'Rust Basics', got '%s'", node.Name)
		}

		// Test GetByTechTree
		nodes, err := nodeRepo.GetByTechTree(ctx, treeID)
		if err != nil {
			t.Fatalf("Failed to get nodes by tech tree: %v", err)
		}
		if len(nodes) != 1 {
			t.Errorf("Expected 1 node, got %d", len(nodes))
		}

		// Test GetPrerequisites (empty for first node)
		prereqs, err := nodeRepo.GetPrerequisites(ctx, nodeID)
		if err != nil {
			t.Fatalf("Failed to get prerequisites: %v", err)
		}
		if len(prereqs) != 0 {
			t.Errorf("Expected 0 prerequisites, got %d", len(prereqs))
		}
	})
}

// TestDatabaseConnection verifies basic database connectivity
func TestDatabaseConnection(t *testing.T) {
	db := database.SetupTestDB(t)

	ctx := context.Background()
	if err := db.Pool.Ping(ctx); err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}

	// Query database version
	var version string
	err := db.Pool.QueryRow(ctx, "SELECT version()").Scan(&version)
	if err != nil {
		t.Fatalf("Failed to query database version: %v", err)
	}

	t.Logf("Connected to: %s", version)
}

// TestDatabaseStats verifies pool statistics
func TestDatabaseStats(t *testing.T) {
	db := database.SetupTestDB(t)

	stats := db.Pool.Stat()

	t.Logf("Pool Stats:")
	t.Logf("  Total Connections: %d", stats.TotalConns())
	t.Logf("  Idle Connections: %d", stats.IdleConns())
	t.Logf("  Acquired Connections: %d", stats.AcquiredConns())
	t.Logf("  Max Connections: %d", stats.MaxConns())

	if stats.MaxConns() <= 0 {
		t.Error("Expected max connections > 0")
	}
}
