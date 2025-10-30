package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/anttihil/code-auto-grader/internal/content"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	// Parse flags
	shouldSeedUsersOnly := flag.Bool("users", false, "Seed users only")
	shouldSeedTreesOnly := flag.Bool("trees", false, "Seed tech trees only")
	shouldSeedNodesOnly := flag.Bool("nodes", false, "Seed nodes only")
	shouldSeedExercisesOnly := flag.Bool("exercises", false, "Seed exercises only")
	shouldSeedAll := flag.Bool("all", false, "Seed all data")
	force := flag.Bool("force", false, "Force overwrite existing data (delete and re-insert)")
	contentDir := flag.String("dir", "../../../content", "Path to content directory")

	flag.Parse()

	// Default to seeding everything if no specific flags
	if !*shouldSeedUsersOnly && !*shouldSeedTreesOnly && !*shouldSeedNodesOnly && !*shouldSeedExercisesOnly {
		*shouldSeedAll = true
	}

	if *shouldSeedAll {
		*shouldSeedUsersOnly = true
		*shouldSeedTreesOnly = true
		*shouldSeedNodesOnly = true
		*shouldSeedExercisesOnly = true
	}

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Get database URL
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	// Connect to database
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}

	// Resolve content directory path
	absContentDir, err := filepath.Abs(*contentDir)
	if err != nil {
		log.Fatalf("Invalid content directory path: %v", err)
	}

	log.Printf("Loading content from: %s\n", absContentDir)

	// Initialize loader
	loader := content.NewLoader(absContentDir)

	// Load all content
	seedData, exercises, err := loader.LoadAll()
	if err != nil {
		log.Fatalf("Failed to load content: %v", err)
	}

	log.Printf("Loaded: %d users, %d tech trees, %d nodes, %d exercises\n",
		len(seedData.Users), len(seedData.TechTrees), len(seedData.Nodes), len(exercises))

	// Seed in order (respecting foreign key constraints)
	if *shouldSeedUsersOnly {
		fmt.Println("\n📝 Seeding users...")
		if err := seedUsers(ctx, pool, seedData.Users, *force); err != nil {
			log.Fatalf("Failed to seed users: %v", err)
		}
		fmt.Printf("✓ Seeded %d users\n", len(seedData.Users))
	}

	if *shouldSeedTreesOnly {
		fmt.Println("\n🌳 Seeding tech trees...")
		if err := seedTechTrees(ctx, pool, seedData.TechTrees, *force); err != nil {
			log.Fatalf("Failed to seed tech trees: %v", err)
		}
		fmt.Printf("✓ Seeded %d tech trees\n", len(seedData.TechTrees))
	}

	if *shouldSeedNodesOnly {
		fmt.Println("\n📚 Seeding nodes...")
		if err := seedNodesWithPrerequisites(ctx, pool, seedData.Nodes, *force); err != nil {
			log.Fatalf("Failed to seed nodes: %v", err)
		}
		fmt.Printf("✓ Seeded %d nodes\n", len(seedData.Nodes))
	}

	if *shouldSeedExercisesOnly {
		fmt.Println("\n💻 Seeding exercises...")
		if err := seedExercises(ctx, pool, exercises, *force); err != nil {
			log.Fatalf("Failed to seed exercises: %v", err)
		}
		fmt.Printf("✓ Seeded %d exercises\n", len(exercises))
	}

	fmt.Println("\n🎉 Seeding complete!")
}

func seedUsers(ctx context.Context, pool *pgxpool.Pool, users []content.User, force bool) error {
	for _, user := range users {
		if force {
			_, err := pool.Exec(ctx, "DELETE FROM users WHERE id = $1", user.ID)
			if err != nil {
				return fmt.Errorf("deleting user %s: %w", user.Username, err)
			}
		}

		_, err := pool.Exec(ctx, `
			INSERT INTO users (id, username, email, display_name, password_hash)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (id) DO NOTHING
		`, user.ID, user.Username, user.Email, user.DisplayName, user.PasswordHash)

		if err != nil {
			return fmt.Errorf("inserting user %s: %w", user.Username, err)
		}

		fmt.Printf("  ✓ %s (%s)\n", user.Username, user.Email)
	}

	return nil
}

func seedTechTrees(ctx context.Context, pool *pgxpool.Pool, trees []content.TechTree, force bool) error {
	for _, tree := range trees {
		if force {
			_, err := pool.Exec(ctx, "DELETE FROM tech_trees WHERE id = $1", tree.ID)
			if err != nil {
				return fmt.Errorf("deleting tech tree %s: %w", tree.Name, err)
			}
		}

		_, err := pool.Exec(ctx, `
			INSERT INTO tech_trees (
				id, name, slug, description, language,
				difficulty_level, estimated_hours, is_published
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			ON CONFLICT (id) DO NOTHING
		`, tree.ID, tree.Name, tree.Slug, tree.Description, tree.Language,
			tree.DifficultyLevel, tree.EstimatedHours, tree.IsPublished)

		if err != nil {
			return fmt.Errorf("inserting tech tree %s: %w", tree.Name, err)
		}

		fmt.Printf("  ✓ %s (%s)\n", tree.Name, tree.Language)
	}

	return nil
}

func seedNodesWithPrerequisites(ctx context.Context, pool *pgxpool.Pool, nodes []content.Node, force bool) error {
	// First pass: insert nodes
	for _, node := range nodes {
		if force {
			_, err := pool.Exec(ctx, "DELETE FROM nodes WHERE id = $1", node.ID)
			if err != nil {
				return fmt.Errorf("deleting node %s: %w", node.Name, err)
			}
		}

		_, err := pool.Exec(ctx, `
			INSERT INTO nodes (
				id, tech_tree_id, name, slug, description,
				learning_objectives, order_index, estimated_minutes
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			ON CONFLICT (id) DO NOTHING
		`, node.ID, node.TechTreeID, node.Name, node.Slug, node.Description,
			node.LearningObjectives, node.OrderIndex, node.EstimatedMinutes)

		if err != nil {
			return fmt.Errorf("inserting node %s: %w", node.Name, err)
		}

		fmt.Printf("  ✓ %s (order %d)\n", node.Name, node.OrderIndex)
	}

	// Second pass: insert prerequisites
	for _, node := range nodes {
		if force {
			_, err := pool.Exec(ctx, "DELETE FROM node_prerequisites WHERE node_id = $1", node.ID)
			if err != nil {
				return fmt.Errorf("deleting prerequisites for %s: %w", node.Name, err)
			}
		}

		for _, prereqID := range node.Prerequisites {
			_, err := pool.Exec(ctx, `
				INSERT INTO node_prerequisites (node_id, prerequisite_node_id)
				VALUES ($1, $2)
				ON CONFLICT DO NOTHING
			`, node.ID, prereqID)

			if err != nil {
				return fmt.Errorf("inserting prerequisite for %s: %w", node.Name, err)
			}
		}
	}

	return nil
}

func seedExercises(ctx context.Context, pool *pgxpool.Pool, exercises []*content.Exercise, force bool) error {
	for _, exercise := range exercises {
		if force {
			_, err := pool.Exec(ctx, "DELETE FROM exercises WHERE id = $1", exercise.ID)
			if err != nil {
				return fmt.Errorf("deleting exercise %s: %w", exercise.Title, err)
			}
		}

		_, err := pool.Exec(ctx, `
			INSERT INTO exercises (
				id, node_id, title, slug, description, difficulty,
				estimated_minutes, order_index, content
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			ON CONFLICT (id) DO NOTHING
		`, exercise.ID, exercise.NodeID, exercise.Title, exercise.Slug,
			exercise.Description, exercise.Difficulty,
			exercise.EstimatedMinutes, exercise.OrderIndex, exercise.Content)

		if err != nil {
			return fmt.Errorf("inserting exercise %s: %w", exercise.Title, err)
		}

		fmt.Printf("  ✓ %s (%s)\n", exercise.Title, exercise.Difficulty)
	}

	return nil
}
