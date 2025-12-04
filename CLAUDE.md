# CLAUDE.md

Go + Echo backend, PostgreSQL + pgx, SolidJS frontend, Docker code execution, Gemini AI.

## Architecture Rules

**MUST follow:**
- Feature folders (`app/features/`) = self-contained vertical slices
- Each feature has: `api.go`, `service.go`, `models.go`, `dto.go`
- Repository layer (`app/repository/`) = shared DB access via interfaces
- Frontend = monolithic `web/src/App.tsx` (don't split yet)
- User ID = hardcoded `00000000-0000-0000-0000-000000000001`
- Scoring = 60% tests + 40% AI, 80% mastery threshold

**NEVER do:**
- Put business logic in API handlers
- Query DB directly in handlers/services (use repositories)
- Use `context.Background()` in handlers (use `c.Request().Context()`)
- Skip tests (100% coverage required)
- Modify Docker security constraints

## File Structure

```
app/
├── features/                           # Vertical slices
│   ├── submissions/                    # Code submission & grading
│   │   ├── api.go                      # HTTP handlers
│   │   ├── service.go                  # Business logic
│   │   ├── models.go                   # Domain types
│   │   ├── dto.go                      # Request/response types
│   │   ├── interfaces.go               # Dependencies
│   │   ├── adapters.go                 # External service bridges
│   │   ├── scorer.go, similarity.go    # Scoring logic
│   │   └── *_test.go                   # Tests
│   ├── exercises/                      # Exercise management
│   ├── tech-trees/                     # Learning paths
│   ├── progress/                       # User progress tracking
│   └── code-execution/                 # Docker executor & pool
├── repository/
│   ├── repository.go                   # Interfaces
│   └── *_repo.go                       # Implementations
├── shared/                             # Cross-cutting concerns
│   ├── config/                         # Configuration
│   ├── database/                       # DB connection
│   ├── gemini/                         # AI client
│   └── middleware/                     # HTTP middleware
├── models/models.go                    # Shared DB models
├── migrations/                         # SQL migrations
└── main.go                             # Wiring + routes
web/src/App.tsx                         # All frontend UI
```

## Adding a Feature Endpoint

```go
// 1. features/[feature]/service.go - business logic
func (s *Service) DoThing(ctx context.Context, id string) (*Result, error)

// 2. features/[feature]/api.go - HTTP handler
func (h *Handler) DoThing(c echo.Context) error {
    result, err := h.service.DoThing(c.Request().Context(), c.Param("id"))
    return c.JSON(200, result)
}

// 3. main.go - register route
api.POST("/thing/:id", handler.DoThing)

// 4. Tests: features/[feature]/service_test.go + api_test.go
```

## Commands

```bash
cd app && go run main.go          # Backend :5000
cd web && npm run dev             # Frontend :5173
cd app && make test               # All tests
cd app && make db-reset           # Drop + migrate + seed
cd app && migrate create -ext sql -dir migrations -seq name  # New migration
```

## Key References

- `app/features/*/doc.go` - Feature documentation
- `app/repository/repository.go` - All repository interfaces
- `app/main.go` - Dependency injection, all routes

## Learning Mode (Go Training)

**MUST follow when writing Go code:**
- Proactively include `TODO(human)` opportunities for any Go code involving:
  - Error handling patterns (when to wrap, when to return early)
  - Interface design and method signatures
  - Struct definitions and field choices
  - Business logic with multiple valid approaches
  - Concurrency patterns (goroutines, channels, mutexes)
- Present "Learn by Doing" blocks and **wait for human implementation** before continuing
- Focus insights on Go idioms, patterns, and best practices
- Prioritize learning over speed - pause for meaningful contributions even on smaller code blocks
