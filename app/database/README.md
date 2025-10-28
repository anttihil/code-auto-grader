# Database Package

This package provides database connectivity, migrations, and repository interfaces for Code Auto-Grader.

## Quick Start

### 1. Set Up Database

**Using Docker:**
```bash
make docker-db-up
```

**Or install PostgreSQL locally** (macOS):
```bash
brew install postgresql@16
brew services start postgresql@16
createdb code_auto_grader
```

### 2. Configure Environment

Copy `.env.example` to `.env` and update DATABASE_URL:
```bash
cp .env.example .env
# Edit .env and set DATABASE_URL
```

### 3. Run Migrations

```bash
make migrate-up
```

### 4. Seed Database

```bash
make db-seed
```

## Common Commands

| Command | Description |
|---------|-------------|
| `make migrate-up` | Run all pending migrations |
| `make migrate-down` | Rollback last migration |
| `make migrate-version` | Show current migration version |
| `make db-seed` | Seed database with content |
| `make db-seed-force` | Re-seed (deletes and re-inserts) |
| `make db-reset` | Drop, migrate, and seed |
| `make db-fresh` | Migrate and seed (no drop) |
| `make db-console` | Open psql console |
| `make help` | Show all available commands |

## Directory Structure

```
database/
├── database.go          # Database connection and initialization
├── test_helpers.go      # Test database setup utilities
├── database_test.go     # Database integration tests
└── README.md           # This file

migrations/
├── 000001_init_schema.up.sql    # Create all tables
└── 000001_init_schema.down.sql  # Drop all tables

repository/
├── repository.go            # Repository interfaces
├── exercise_repo.go         # Exercise CRUD operations
├── submission_repo.go       # Submission CRUD operations
├── tech_tree_repo.go        # Tech tree CRUD operations
├── node_repo.go             # Node CRUD operations
├── user_repo.go             # User CRUD operations
└── user_progress_repo.go    # Progress tracking operations
```

## Database Schema

### Core Tables

- **users** - User accounts
- **tech_trees** - Learning paths (e.g., "Rust Fundamentals")
- **nodes** - Topics within tech trees
- **node_prerequisites** - DAG structure for node dependencies
- **exercises** - Individual coding exercises
- **submissions** - User submission history
- **user_progress** - Progress tracking at node level

### Relationships

```
tech_trees
    ↓
  nodes ←→ node_prerequisites (DAG)
    ↓
exercises
    ↓
submissions → users
    ↑
user_progress
```

## Usage in Code

### Initialize Database Connection

```go
import "github.com/anttihil/code-auto-grader/database"

// From environment variables
db, err := database.NewFromEnv(ctx)
if err != nil {
    log.Fatal(err)
}
defer db.Close()

// Or with custom config
cfg := database.DefaultConfig()
cfg.URL = "postgres://..."
cfg.MaxConns = 50
db, err := database.New(ctx, cfg)
```

### Use Repositories

```go
import "github.com/anttihil/code-auto-grader/repository"

// Create repositories
exerciseRepo := repository.NewExerciseRepo(db.Pool)
submissionRepo := repository.NewSubmissionRepo(db.Pool)

// Fetch an exercise
exercise, err := exerciseRepo.GetByID(ctx, exerciseID)

// Save a submission
submission := &models.Submission{
    ID:          uuid.NewString(),
    UserID:      userID,
    ExerciseID:  exerciseID,
    // ... other fields
}
err = submissionRepo.Create(ctx, submission)
```

## Testing

### Set Up Test Database

```bash
# One-time setup
make test-db-setup

# Reset between test runs
make test-db-reset
```

### Use Test Helpers

```go
import "github.com/anttihil/code-auto-grader/database"

func TestMyHandler(t *testing.T) {
    db := database.SetupTestDB(t) // Automatically migrates

    // Seed test data
    userID := db.SeedUser(t, "testuser", "test@example.com")
    treeID := db.SeedTechTree(t, "Test Tree", "test-tree", "go")

    // Run tests using db.Pool

    // Cleanup happens automatically via t.Cleanup()
}
```

## Creating New Migrations

```bash
make migrate-create NAME=add_user_badges
```

This creates:
- `migrations/000002_add_user_badges.up.sql`
- `migrations/000002_add_user_badges.down.sql`

**Best practices:**
- Always write both up and down migrations
- Test down migrations before committing
- Use transactions where appropriate (`BEGIN;` ... `COMMIT;`)
- Add comments explaining complex migrations

## Seeding Content

### Seed File Structure

```
content/
├── users.json           # Seed users
├── tech-trees.json      # Tech tree definitions
├── nodes.json           # Node definitions with prerequisites
└── rust/                # Exercise markdown files
    ├── fundamentals/
    ├── ownership/
    └── structs-enums/
```

### Exercise File Format

Each exercise is a markdown file with JSON frontmatter:

```markdown
{
  "id": "uuid",
  "node_id": "uuid",
  "title": "Exercise Title",
  "slug": "exercise-slug",
  "difficulty": "easy",
  "estimated_minutes": 15,
  "order_index": 1
}
---

# Exercise Title

## Problem Description
...

## Starting Code
\`\`\`rust
fn main() { }
\`\`\`

## Reference Solution
\`\`\`rust
fn main() { println!("Hello!"); }
\`\`\`

## Test Cases
\`\`\`rust
#[test]
fn test_something() { }
\`\`\`
```

### Selective Seeding

```bash
# Seed only specific parts
make db-seed-users
make db-seed-trees
make db-seed-nodes
make db-seed-exercises

# Force re-seed (deletes existing data)
make db-seed-force
```

## Troubleshooting

### Migration Version Conflicts

If migrations are in a "dirty" state:

```bash
# Check version
make migrate-version

# Force to specific version
make migrate-force VERSION=1

# Then retry
make migrate-up
```

### Connection Issues

```bash
# Test connection
psql "$DATABASE_URL"

# Check if PostgreSQL is running
ps aux | grep postgres

# Restart PostgreSQL (macOS)
brew services restart postgresql@16
```

### Seed Failures

```bash
# Check what's currently in database
make db-console

# Drop and start fresh
make db-reset
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_URL` | Required | PostgreSQL connection string |
| `TEST_DATABASE_URL` | Required | Test database connection string |
| `DB_AUTO_MIGRATE` | `true` | Run migrations on app startup |
| `DB_MIGRATIONS_PATH` | `file://migrations` | Path to migration files |

## Connection Pool Settings

Default pool configuration (can be customized via `Config`):

- **MaxConns**: 25 (maximum connections)
- **MinConns**: 5 (minimum idle connections)
- **MaxConnLifetime**: 1 hour
- **MaxConnIdleTime**: 30 minutes
- **HealthCheckPeriod**: 1 minute

Adjust based on your workload and database server capacity.