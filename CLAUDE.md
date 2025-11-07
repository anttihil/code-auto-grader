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

1. **Domain-Oriented Architecture**: Backend follows clean architecture with separation of concerns

   - **Domain Layer** (`app/domain/`): Pure business logic, framework-agnostic
   - **API Layer** (`app/api/`): Thin HTTP handlers for each domain
   - **Repository Layer** (`app/repository/`): Data access abstraction
   - See `app/domain/README.md` for detailed architecture guide

2. **Hardcoded User ID**: `00000000-0000-0000-0000-000000000001` for MVP

   - Search for: `TODO: Get real user ID from auth context`

3. **Repository Pattern Required**: Handlers and services never query database directly

   - Interfaces: `app/repository/repository.go`
   - Implementations: `app/repository/*_repo.go`
   - Wiring: `app/main.go` (dependency injection)

4. **JSONB Exercise Content**: `exercises.content` stores `{starter_code, test_cases[], hints[], evaluation_rubric}`

   - Parse in application layer, use JSONB operators in queries

5. **Monolithic Frontend**: All UI in `web/src/App.tsx` (don't split components yet)

6. **Scoring Algorithm**: 60% test results + 40% AI evaluation, 80% threshold for mastery

## Common Task Patterns

### Adding API Endpoint

**For a new endpoint in an existing domain:**

1. **Domain Service** (`app/domain/[domain]/service.go`):
   - Add business logic method to the service
   - Return domain models, not database models
   - Keep logic framework-agnostic (testable without HTTP)

2. **API Handler** (`app/api/[domain]/handlers.go`):
   - Add thin HTTP handler method (50-100 LOC)
   - Parse request, call domain service, format response
   - Handle HTTP-specific concerns (status codes, validation)

3. **Routes** (`app/main.go`):
   - Register new route with handler method
   - Use appropriate Echo group (e.g., `api.POST("/endpoint/:id", handler.Method)`)

4. **Tests**:
   - Add domain service test to `app/domain/[domain]/service_test.go` (with mocks)
   - Add HTTP handler test to `app/api/[domain]/handlers_test.go` (with httptest)

**Example workflow:**
```go
// 1. Add to domain/grading/service.go
func (s *GradingService) GetSubmission(ctx context.Context, submissionID string) (*SubmissionResult, error) {
    // Business logic here
}

// 2. Add to api/grading/handlers.go
func (h *SubmissionHandler) GetSubmission(c echo.Context) error {
    result, err := h.service.GetSubmission(c.Request().Context(), c.Param("id"))
    // HTTP handling here
}

// 3. Add to main.go
api.GET("/submission/:id", submissionHandler.GetSubmission)
```

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

**Domain-Oriented Architecture** (read these first):

```
app/
├── domain/                      # Business Logic Layer
│   ├── grading/                 # Submission & evaluation domain
│   │   ├── service.go          # Submission workflow orchestration
│   │   ├── scorer.go           # Scoring algorithm (60/40 split)
│   │   ├── similarity.go       # Anti-gaming validation
│   │   ├── models.go           # Domain types & constants
│   │   ├── interfaces.go       # Clean boundaries (executor, AI)
│   │   └── adapters.go         # Bridge to external services
│   ├── learning/               # Tech tree & progress domain
│   │   ├── service.go          # Navigation & progress tracking
│   │   ├── models.go           # View models (TechTreeView, NodeView)
│   │   └── errors.go           # Domain-specific errors
│   └── content/                # Exercise management domain
│       ├── service.go          # CRUD & AI generation
│       ├── adapters.go         # Gemini service adapter
│       └── models.go           # Domain types
│
├── api/                        # HTTP Layer (Thin handlers)
│   ├── grading/handlers.go     # Submission endpoints
│   ├── learning/handlers.go    # Tech tree & node endpoints
│   └── content/handlers.go     # Exercise endpoints
│
├── repository/                 # Data Access Layer
│   ├── repository.go           # Repository interfaces
│   └── *_repo.go              # Repository implementations
│
├── models/models.go            # Database models (shared)
├── services/                   # External services (executor, Gemini)
└── main.go                     # Wiring & routing
```

**Key Files to Read:**
- `app/domain/README.md` - Architecture patterns and examples
- `app/main.go` - Dependency injection and routing
- `app/repository/repository.go` - Data access interfaces

**Multi-file Components** (use Glob to explore):
- `app/domain/*/service.go` - Domain services
- `app/api/*/handlers.go` - HTTP handlers
- `app/repository/*_repo.go` - Repository implementations
- `app/migrations/*.sql` - Database migrations

## Request Lifecycle (Submission Flow)

```
POST /api/submission/:exerciseId
  ↓
  HTTP Layer (api/grading/handlers.go)
  → SubmissionHandler.SubmitToExercise
    - Parse request body & params
    - Extract user ID from context (hardcoded for MVP)
    ↓
  Domain Layer (domain/grading/service.go)
  → GradingService.SubmitCode
    - Fetch exercise from repository
    - Check code similarity vs starter code (similarity.go)
    - Execute code in Docker container (via CodeExecutor interface)
    - Run AI evaluation with rubric (via AIEvaluator interface)
    - Calculate weighted score: 60% tests + 40% AI (scorer.go)
    - Save submission to repository
    - Update user progress if score ≥ 80% (mastery threshold)
    ↓
  HTTP Response
  → Return JSON {submission, testResults, aiFeedback, score, masteryAchieved}
```

**Key Insight**: The domain layer orchestrates business logic without knowing about HTTP. The API layer handles HTTP concerns. This separation enables:
- Testing domain logic without HTTP mocks
- Reusing domain logic in CLI tools, cron jobs, etc.
- Changing HTTP frameworks without touching business logic

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

- ❌ Put business logic in API handlers (use domain services)
- ❌ Query database directly in handlers or services (use repositories)
- ❌ Add HTTP dependencies to domain layer (keep it framework-agnostic)
- ❌ Skip writing tests (we have 100% domain + API test coverage)
- ❌ Add user authentication yet (MVP uses hardcoded ID)
- ❌ Modify Docker security constraints (no network, read-only FS, resource limits)
- ❌ Split frontend into components (keep `App.tsx` monolithic)
- ❌ Add dependencies without checking (minimal: Echo, pgx, Gemini, Docker SDK)
- ❌ Use `context.Background()` in handlers (use `c.Request().Context()`)

## Testing

```bash
cd app && make test                    # All tests
cd app && go test ./domain/...         # All domain tests (fast, no DB)
cd app && go test ./api/...            # All API handler tests
cd app && go test ./repository/...     # Repository tests (requires DB)
cd app && go test -v -run Integration  # Integration tests
```

**Test Coverage**: 100% coverage of domain and API layers
- **Domain Tests** (`domain/*/service_test.go`): Pure unit tests with mocks, no HTTP/DB
- **API Tests** (`api/*/handlers_test.go`): HTTP handler tests with httptest
- **Repository Tests**: Integration tests with real test database
- **Test Database**: Uses `TEST_DATABASE_URL` (port 5433)
- **Pattern**: Use `database.NewTestDatabase(t)` in tests for auto-setup/cleanup

**Test Structure**:
```
✅ domain/grading/scorer_test.go       - Scoring algorithm (35 tests)
✅ domain/grading/similarity_test.go   - Code similarity (15 tests)
✅ domain/grading/service_test.go      - Submission workflow (20 tests)
✅ domain/learning/service_test.go     - Tech tree navigation (15 tests)
✅ domain/content/service_test.go      - Exercise management (12 tests)
✅ api/grading/handlers_test.go        - Grading endpoints (14 tests)
✅ api/learning/handlers_test.go       - Learning endpoints (11 tests)
✅ api/content/handlers_test.go        - Content endpoints (12 tests)
```

See `REFACTORING_STATUS.md` for complete test strategy guide.

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
