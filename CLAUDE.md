# CLAUDE.md

This file provides guidance to Claude Code when working with this repository.

## Project Overview

LLM-powered educational platform for generating coding exercises and providing AI-driven feedback.

**Current Phase**: Dogfooding MVP - Week 3 (Learning Loop Integration)

- Backend: ✅ Complete (APIs, test execution, AI evaluation, mastery tracking)
- Frontend: 🔄 Next (Tech tree navigation UI, exercise history)
- Goal: Personal Rust learning validation

**Tech Stack**: Go + Echo, PostgreSQL + pgx, SolidJS, Docker (code execution), Gemini AI

## Critical Architectural Decisions

Before coding, understand these constraints:

1. **Single Handler File**: All HTTP handlers in `app/handlers/handlers.go` (not split by domain)

   - Add methods to existing structs: `SubmissionHandler`, `ExerciseHandler`, `TechTreeHandler`, `NodeHandler`

2. **Hardcoded User ID**: `00000000-0000-0000-0000-000000000001` for MVP

   - Search for: `TODO: Get real user ID from auth context`

3. **Repository Pattern Required**: Handlers never query database directly

   - Interfaces: `app/repository/repository.go`
   - Implementations: `app/repository/*_repo.go`
   - Wiring: `app/main.go` (dependency injection)

4. **JSONB Exercise Content**: `exercises.content` stores `{starter_code, test_cases[], hints[], evaluation_rubric}`

   - Parse in application layer, use JSONB operators in queries

5. **Monolithic Frontend**: All UI in `web/src/App.tsx` (don't split components yet)

6. **Scoring Algorithm**: 60% test results + 40% AI evaluation, 80% threshold for mastery

## Common Task Patterns

### Adding API Endpoint

1. Read `app/handlers/handlers.go` → add method to handler struct
2. Read `app/main.go` → add route with handler method
3. Use repositories for database access (never direct queries)
4. Add integration test to `app/integration_test.go`

### Modifying Database Schema

```bash
cd app && migrate create -ext sql -dir migrations -seq description
# Edit .up.sql and .down.sql
cd app && make db-reset  # Test migration
# Update app/models/models.go and repository methods
```

### Frontend Development

- Edit: `web/src/App.tsx` (single-page app)
- API: fetch to `/api/*` (Vite proxies to :5000)
- Keep monolithic until patterns emerge

### Database Seeding

```bash
cd app && make db-seed  # or: go run cmd/seed/main.go
# Content: app/content/*.md
# Parser: app/internal/content/parser.go
```

## File Structure Quick Reference

**Key Backend Files** (read these first):

- `app/handlers/handlers.go` - ALL HTTP handlers
- `app/models/models.go` - Database models
- `app/repository/repository.go` - Repository interfaces
- `app/main.go` - Server setup, DI, routing

**Multi-file Components** (use Glob to explore):

- `app/repository/*_repo.go` - Repository implementations
- `app/services/executor/*.go` - Code execution
- `app/migrations/*.sql` - Database migrations

## Request Lifecycle (Submission Flow)

```
POST /api/submission/:exerciseId
  → SubmissionHandler.SubmitExercise (handlers.go)
  → ExerciseRepository.GetByID (fetch test cases from DB)
  → CodeExecutionService.Execute (Docker container, 5s timeout)
  → GeminiService.EvaluateCode (AI feedback with rubric)
  → Calculate Score (60% tests + 40% AI)
  → SubmissionRepository.Create (save to DB)
  → UserProgressRepository.UpdateProgress (update mastery if ≥80%)
  → Return JSON {submission, testResults, aiFeedback, score}
```

## Common Issues & Quick Fixes

**"Failed to execute code"**

```bash
docker ps -a | grep code-grader  # Check containers
curl localhost:5000/api/pool/status  # Check pool
cd app && docker build -f docker/*.Dockerfile -t code-grader-* .  # Rebuild
```

**Database connection errors**

```bash
make docker-db-up  # Start PostgreSQL
make migrate-up    # Run migrations
cat app/.env | grep DATABASE_URL  # Verify config
```

**"Dirty database" migration errors**

```bash
make db-reset  # Nuclear option: drop + migrate + seed
```

**Frontend CORS/404 errors**

```bash
curl localhost:5000/api/health  # Backend must be on :5000
# Frontend must be on :5173 (Vite proxy config)
```

## Development Commands

**Setup** (one-time):

```bash
cd app && cp .env.example .env  # Add GEMINI_API_KEY
cd app && make install-tools    # Install migrate CLI
cd app && make docker-db-up     # Start PostgreSQL
cd app && make migrate-up       # Run migrations
cd app && make db-seed          # Seed exercises
```

**Daily workflow**:

```bash
# Terminal 1: Backend
cd app && go run main.go        # :5000

# Terminal 2: Frontend
cd web && npm run dev           # :5173 → proxies to :5000

# Database
cd app && make db-reset         # Fresh start (drop/migrate/seed)
cd app && make test             # Run all tests
```

**Production build**:

```bash
cd web && npm run build         # Outputs to web/dist/
cd app && go build -o code-auto-grader
./app/code-auto-grader          # Serves both API + static frontend
```

## What NOT to Do

- ❌ Split handlers into multiple files (keep in `handlers.go`)
- ❌ Query database directly in handlers (use repositories)
- ❌ Add user authentication yet (MVP uses hardcoded ID)
- ❌ Modify Docker security constraints (no network, read-only FS, resource limits)
- ❌ Split frontend into components (keep `App.tsx` monolithic)
- ❌ Add dependencies without checking (minimal: Echo, pgx, Gemini, Docker SDK)
- ❌ Use `context.Background()` in handlers (use `c.Request().Context()`)

## Testing

```bash
cd app && make test              # All tests
cd app && go test ./repository/  # Specific package
cd app && go test -v -run Integration  # Integration only
```

**Test database**: Uses `TEST_DATABASE_URL` (port 5433)
**Pattern**: Use `database.NewTestDatabase(t)` in tests for auto-setup/cleanup
**Integration tests**: `app/integration_test.go` tests full HTTP lifecycle

## Environment Variables

**Required in `app/.env`**:

```bash
GEMINI_API_KEY=your_key_here
DATABASE_URL=postgres://code_grader:dev_password@localhost:5432/code_auto_grader?sslmode=disable
```

**Optional**:

```bash
PORT=5000                    # Server port
ENVIRONMENT=development      # or "production"
JS_POOL_SIZE=3              # Container pool sizes
PYTHON_POOL_SIZE=3
RUST_POOL_SIZE=2
```

## Code Execution Security

Docker containers are heavily sandboxed:

- No network access
- Read-only filesystem
- 5-second timeout
- 128MB memory, 0.5 CPU cores, 50 PIDs max
- No capabilities, no privilege escalation

Pre-warmed pools per language reduce latency. Check status: `GET /api/pool/status`

## Key API Endpoints

**Implemented**:

- `GET /api/tech-tree/:treeId` - Tech tree with nodes and progress
- `GET /api/node/:nodeId` - Node with exercises and progress
- `GET /api/exercise/:exerciseId` - Exercise with submission history
- `POST /api/submission/:exerciseId` - Submit code for evaluation
  - Body: `{code, language}`
  - Returns: test results, AI feedback, score, mastery updates

**Monitoring**:

- `GET /api/health` - Health check
- `GET /api/pool/status` - Container pool status

## Database Schema

**Core tables**: `users`, `tech_trees`, `nodes`, `exercises`, `submissions`, `user_progress`

**Key features**:

- JSONB content fields for flexible exercise structure
- DAG for node prerequisites
- Automatic `updated_at` triggers
- Mastery threshold: 80%

See `app/database/README.md` for detailed schema docs.

## Next Steps (Current Priority)

Week 3 UI tasks:

1. Tech tree navigation UI
2. Exercise history view
3. "Next Exercise" button
4. Progress indicators

See `IMPLEMENTATION_PLAN.md` for roadmap.
