# Database Migrations

This directory contains SQL migration files for the Code Auto-Grader database schema.

## Migration Files

Migrations are numbered sequentially and come in pairs (up/down):

- **000001_init_schema**: Create all core database tables
- **000002_seed_dogfooding_data**: Insert initial data for MVP (user, tech tree, nodes)

## Running Migrations

### Option 1: Using golang-migrate CLI

Install golang-migrate:
```bash
# macOS
brew install golang-migrate

# Linux
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/

# Or using Go
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Run migrations:
```bash
# Set database URL
export DATABASE_URL="postgres://user:password@localhost:5432/code_auto_grader?sslmode=disable"

# Migrate up (apply all pending migrations)
migrate -database ${DATABASE_URL} -path ./app/migrations up

# Migrate down (revert last migration)
migrate -database ${DATABASE_URL} -path ./app/migrations down 1

# Check current version
migrate -database ${DATABASE_URL} -path ./app/migrations version

# Force a specific version (if migrations are dirty)
migrate -database ${DATABASE_URL} -path ./app/migrations force VERSION
```

### Option 2: Using psql directly

```bash
# Run up migration
psql -U user -d code_auto_grader -f app/migrations/000001_init_schema.up.sql
psql -U user -d code_auto_grader -f app/migrations/000002_seed_dogfooding_data.up.sql

# Run down migration
psql -U user -d code_auto_grader -f app/migrations/000002_seed_dogfooding_data.down.sql
psql -U user -d code_auto_grader -f app/migrations/000001_init_schema.down.sql
```

### Option 3: Using custom Go script (recommended for development)

```bash
# Create database
createdb code_auto_grader

# Run migrations
go run scripts/migrate.go up

# Revert migrations
go run scripts/migrate.go down
```

## Seeding Exercise Data

After running migrations, seed exercises from markdown files:

```bash
go run scripts/seed_exercises.go
```

This script:
1. Parses exercise markdown files in `/exercises/` directory
2. Extracts content (starter code, tests, hints, etc.)
3. Inserts exercises with JSONB content into database

## Database Setup (PostgreSQL)

### Install PostgreSQL

```bash
# macOS
brew install postgresql@15
brew services start postgresql@15

# Ubuntu/Debian
sudo apt-get install postgresql-15
sudo systemctl start postgresql

# Create user and database
createuser -s code_auto_grader_user
createdb -O code_auto_grader_user code_auto_grader
```

### Environment Variables

Set these in `app/.env`:

```bash
# Database connection
DATABASE_URL="postgres://code_auto_grader_user:password@localhost:5432/code_auto_grader?sslmode=disable"

# Alternative: separate variables
DB_HOST=localhost
DB_PORT=5432
DB_USER=code_auto_grader_user
DB_PASSWORD=password
DB_NAME=code_auto_grader
DB_SSLMODE=disable
```

## Migration Workflow

### Creating a New Migration

1. Create numbered files (increment from last):
   ```bash
   touch app/migrations/000003_add_achievements.up.sql
   touch app/migrations/000003_add_achievements.down.sql
   ```

2. Write the up migration (schema changes)
3. Write the down migration (revert changes)
4. Test both directions:
   ```bash
   migrate up    # Apply changes
   migrate down 1  # Revert
   migrate up    # Re-apply
   ```

### Migration Best Practices

1. **Always create both up and down migrations**
2. **Test migrations on a copy of production data**
3. **Make migrations idempotent when possible** (use IF NOT EXISTS)
4. **Avoid data migrations in schema migrations** (create separate data migrations)
5. **Use transactions** (migrations run in a transaction by default)
6. **Document complex migrations** with comments

### Example: Adding a Column

**000003_add_user_timezone.up.sql**:
```sql
ALTER TABLE users
ADD COLUMN timezone VARCHAR(50) DEFAULT 'UTC';

COMMENT ON COLUMN users.timezone IS 'User timezone for scheduling';
```

**000003_add_user_timezone.down.sql**:
```sql
ALTER TABLE users
DROP COLUMN timezone;
```

## Troubleshooting

### Dirty Migration State

If a migration fails partway through:

```bash
# Check migration version and dirty status
migrate version

# Force version (be careful!)
migrate force VERSION
```

### Connection Issues

```bash
# Test database connection
psql $DATABASE_URL -c "SELECT 1"

# Check PostgreSQL is running
pg_isready
```

### Reset Database (Development Only)

```bash
# Drop and recreate database
dropdb code_auto_grader
createdb code_auto_grader

# Run all migrations
migrate up
```

## Schema Validation

After running migrations, validate the schema:

```bash
# Check tables exist
psql $DATABASE_URL -c "\dt"

# Check indexes
psql $DATABASE_URL -c "\di"

# Verify constraints
psql $DATABASE_URL -c "
SELECT
    conname AS constraint_name,
    contype AS constraint_type,
    conrelid::regclass AS table_name
FROM pg_constraint
WHERE connamespace = 'public'::regnamespace
ORDER BY conrelid::regclass::text, contype;
"
```

## See Also

- [DATABASE_SCHEMA.md](../../DATABASE_SCHEMA.md) - Full schema documentation
- [scripts/migrate.go](../../scripts/migrate.go) - Migration runner
- [scripts/seed_exercises.go](../../scripts/seed_exercises.go) - Exercise seeder
