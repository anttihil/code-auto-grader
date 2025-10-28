# Transaction Architecture Refactor Notes

## Current Issue

The transaction helpers in `app/repository/transaction_helpers.go` break the repository abstraction:
- They accept raw `pgx.Tx` and duplicate SQL queries
- Callers must know about pgx implementation details
- Can't be tested or mocked easily
- Not idiomatic Go/pgx pattern

## Recommended Pattern: Querier Interface

**Key Insight**: Both `*pgxpool.Pool` and `pgx.Tx` implement the same query methods. Use an interface they both satisfy.

```go
// Querier is satisfied by both pool and transaction
type Querier interface {
    Query(ctx, sql, args...) (pgx.Rows, error)
    QueryRow(ctx, sql, args...) pgx.Row
    Exec(ctx, sql, args...) (pgconn.CommandTag, error)
}
```

## Refactoring Steps

### 1. Add Querier Interface
**File**: `app/repository/querier.go`
- Define Querier interface
- Create txWrapper type to adapt pgx.Tx

### 2. Refactor Repositories
**Files**: All `*_repo.go` files

Add `WithTx()` method to each repository:
```go
func (r *SubmissionRepo) WithTx(tx pgx.Tx) *SubmissionRepo {
    return &SubmissionRepo{pool: txWrapper{tx: tx}}
}
```

**No other changes needed** - existing methods work with both pool and tx!

### 3. Create Transaction Manager
**File**: `app/repository/transaction.go`
- Move `WithTransaction` helper from database.go
- Returns TxManager struct with pool

### 4. Create Service Layer
**File**: `app/services/submission_service.go`

Services orchestrate repositories within transactions:
```go
type SubmissionService struct {
    txManager      *TxManager
    submissionRepo *SubmissionRepo
    progressRepo   *UserProgressRepo
}

func (s *SubmissionService) SubmitExercise(...) error {
    return s.txManager.WithTransaction(ctx, func(tx pgx.Tx) error {
        // Get transaction-aware repos
        txSub := s.submissionRepo.WithTx(tx)
        txProg := s.progressRepo.WithTx(tx)

        // Use normal Create methods - they work in transaction!
        txSub.Create(ctx, submission)
        txProg.Upsert(ctx, progress)
        return nil
    })
}
```

### 5. Update Transaction Helpers
**File**: `app/repository/transaction_helpers.go`

Refactor to accept repository instances:
```go
func CreateSubmissionWithProgress(
    ctx context.Context,
    tx pgx.Tx,
    submissionRepo *SubmissionRepo,  // Accept repo, not raw SQL
    progressRepo *UserProgressRepo,
    submission *models.Submission,
    progress *models.UserProgress,
) error {
    txSub := submissionRepo.WithTx(tx)
    txProg := progressRepo.WithTx(tx)

    // Reuse existing repository methods!
    if err := txSub.Create(ctx, submission); err != nil {
        return err
    }
    return txProg.Upsert(ctx, progress)
}
```

### 6. Update main.go Dependency Injection

```go
// Create repos
submissionRepo := repository.NewSubmissionRepo(pool)
progressRepo := repository.NewUserProgressRepo(pool)

// Create transaction manager
txManager := repository.NewTxManager(pool)

// Create services (they handle transactions)
submissionService := services.NewSubmissionService(
    txManager,
    submissionRepo,
    progressRepo,
)

// Handlers call services (not repos directly)
submissionHandler := handlers.NewSubmissionHandler(submissionService)
```

## Architecture Layers

```
Handlers → Services → Repositories → Database
           (transactions)  (data access)
```

- **Handlers**: HTTP, validation, call services
- **Services**: Business logic, manage transactions
- **Repositories**: Data access (work with/without tx via Querier)

## Benefits

1. **No duplicate SQL** - Repositories work in both modes
2. **Testable** - Mock Querier interface
3. **Clean separation** - Each layer has clear responsibility
4. **Idiomatic** - Follows pgx and Go best practices
5. **Type safe** - Compile-time checking

## Implementation Order

1. ✅ **Keep current working code** - No rush to refactor
2. Create Querier interface + txWrapper
3. Add WithTx() to one repository (test it works)
4. Create simple service for one use case
5. Expand pattern to other repositories
6. Update handlers to use services

## Files to Create/Modify

**New files**:
- `app/repository/querier.go` - Querier interface
- `app/repository/transaction.go` - TxManager
- `app/services/submission_service.go` - Business logic

**Modify**:
- `app/repository/*_repo.go` - Add WithTx() methods
- `app/repository/transaction_helpers.go` - Accept repo instances
- `app/handlers/submission.go` - Call services not repos
- `app/main.go` - Wire up services

## References

- pgx official examples use this pattern
- sqlc generates code with this pattern
- Ent (Facebook) uses this pattern

## Notes

- Current code **works** - this is about clean architecture
- Refactor incrementally - one repository at a time
- Services are optional but recommended for complex transactions
- Keep transaction helpers as high-level orchestrators