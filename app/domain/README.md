# Domain Layer Architecture

This directory contains the **domain layer** of the application - pure business logic that is independent of HTTP, databases, or any external frameworks.

## 📐 Architecture Pattern: Domain-Oriented Design

The backend follows a **clean architecture** pattern with three main layers:

```
┌─────────────────────────────────────────────────────┐
│  API Layer (api/)                                    │
│  • HTTP request/response handling                    │
│  • Input validation & parsing                        │
│  • HTTP status codes & error formatting             │
│  • Thin adapters (50-100 LOC per handler)           │
└──────────────────┬──────────────────────────────────┘
                   │ calls
                   ↓
┌─────────────────────────────────────────────────────┐
│  Domain Layer (domain/) ← YOU ARE HERE               │
│  • Business logic & workflows                        │
│  • Domain models & validation rules                  │
│  • Framework-agnostic (no HTTP, no DB)               │
│  • Testable with mocks (no integration needed)       │
└──────────────────┬──────────────────────────────────┘
                   │ uses
                   ↓
┌─────────────────────────────────────────────────────┐
│  Repository Layer (repository/)                      │
│  • Data access abstraction                           │
│  • SQL queries & database operations                 │
│  • Repository pattern interfaces                     │
└─────────────────────────────────────────────────────┘
```

## 🎯 Key Principles

### 1. Framework-Agnostic
Domain services have **zero dependencies** on HTTP frameworks or databases:
- ✅ Can be tested without spinning up a server
- ✅ Can be called from CLI tools, cron jobs, or other services
- ✅ Easy to change web frameworks (Echo → Fiber → Gin) without touching domain logic

### 2. Dependency Inversion
Domain services depend on **interfaces**, not concrete implementations:
```go
// domain/grading/interfaces.go
type CodeExecutor interface {
    ExecuteCodeWithContext(ctx context.Context, code, language string) executor.ExecutionResult
}

// domain/grading/service.go
type GradingService struct {
    executor CodeExecutor  // Interface, not *executor.Service
}
```

This allows:
- Easy mocking in tests
- Swapping implementations without changing domain code
- Clear boundaries between layers

### 3. Pure Business Logic
Domain code focuses on **what** the application does, not **how** it's delivered:
- Scoring algorithms (60% tests + 40% AI)
- Code similarity detection for anti-gaming
- Progress tracking and mastery calculation
- Exercise generation workflows

## 📦 Domain Structure

### Three Domains

#### 1. **grading/** - Submission & Evaluation
**Responsibility**: Handle code submissions, execute tests, evaluate with AI, calculate scores

**Files**:
- `service.go` - Orchestrates submission workflow
- `scorer.go` - Scoring algorithm (60/40 weighted average)
- `similarity.go` - Anti-gaming code similarity checks
- `models.go` - Domain types (`SubmitCodeRequest`, `SubmissionResult`)
- `interfaces.go` - Clean boundaries (`CodeExecutor`, `AIEvaluator`)
- `adapters.go` - Bridge legacy services to domain interfaces

**Key Methods**:
```go
func (s *GradingService) SubmitCode(ctx context.Context, req SubmitCodeRequest, requestID string) (*SubmissionResult, error)
func (s *GradingService) RunCode(ctx context.Context, req RunCodeRequest, requestID string) (*RunCodeResult, error)
```

#### 2. **learning/** - Tech Tree & Progress
**Responsibility**: Navigate tech trees, track progress, calculate completion percentages

**Files**:
- `service.go` - Tech tree navigation and progress calculation
- `models.go` - View models (`TechTreeView`, `NodeView`, `NodeProgress`)
- `errors.go` - Domain-specific errors

**Key Methods**:
```go
func (s *LearningService) GetTechTree(ctx context.Context, treeID, userID, requestID string) (*TechTreeView, error)
func (s *LearningService) GetNode(ctx context.Context, nodeID, userID, requestID string) (*NodeView, error)
func (s *LearningService) CalculateProgress(ctx context.Context, userID, treeID string) (completed, total int, err error)
```

#### 3. **content/** - Exercise Management
**Responsibility**: CRUD operations for exercises, AI-powered exercise generation

**Files**:
- `service.go` - Exercise retrieval and generation
- `adapters.go` - Gemini service adapter for AI generation
- `models.go` - Domain types (`ExerciseView`, `ExerciseGenerationResult`)

**Key Methods**:
```go
func (s *ContentService) GetExercise(ctx context.Context, exerciseID, userID, requestID string) (*ExerciseView, error)
func (s *ContentService) GenerateExercise(ctx context.Context, req GenerateExerciseRequest, requestID string) (*ExerciseGenerationResult, error)
```

## 🔄 Request Flow Example

Let's trace a code submission through the architecture:

```go
// 1. HTTP Request arrives at API layer
// POST /api/submission/:exerciseId
// Body: {"code": "function add(a,b){return a+b}", "language": "javascript"}

// 2. API Handler (api/grading/handlers.go)
func (h *SubmissionHandler) SubmitToExercise(c echo.Context) error {
    // Parse HTTP request
    var req dto.SubmitCodeRequest
    c.Bind(&req)

    // Call domain service (no HTTP knowledge)
    result, err := h.service.SubmitCode(c.Request().Context(), domainReq, requestID)

    // Format HTTP response
    return c.JSON(http.StatusOK, result)
}

// 3. Domain Service (domain/grading/service.go)
func (s *GradingService) SubmitCode(ctx context.Context, req SubmitCodeRequest, requestID string) (*SubmissionResult, error) {
    // Step 1: Get exercise
    exercise, err := s.exerciseRepo.GetByID(ctx, req.ExerciseID)

    // Step 2: Check similarity to prevent gaming
    if checkCodeSimilarity(req.Code, exercise.StarterCode) == Identical {
        return nil, ErrIdenticalToStarterCode
    }

    // Step 3: Execute code in Docker container
    execResult := s.executor.ExecuteCodeWithContext(ctx, req.Code, req.Language)

    // Step 4: Calculate test score
    testScore := calculateTestScore(execResult, exercise.TestCases)

    // Step 5: Get AI evaluation
    aiScore, feedback := s.evaluator.EvaluateCode(ctx, req.Code, exercise)

    // Step 6: Calculate final weighted score (60% tests + 40% AI)
    finalScore := CalculateScore(testScore, aiScore, exercise.GradingRubric)

    // Step 7: Determine if mastery achieved (≥80%)
    passed := DeterminePassStatus(finalScore)

    // Step 8: Save submission
    submission, err := s.submissionRepo.Create(ctx, ...)

    // Step 9: Update progress if passed
    if passed {
        s.progressRepo.UpdateProgress(ctx, req.UserID, exercise.NodeID, finalScore)
    }

    return &SubmissionResult{
        SubmissionID: submission.ID,
        TotalScore:   finalScore,
        TestScore:    testScore,
        AIScore:      aiScore,
        AIFeedback:   feedback,
        Passed:       passed,
    }, nil
}

// 4. Repository (repository/submission_repo.go)
//    Handles database operations (not shown)

// 5. External Services (services/executor, services/gemini)
//    Wrapped by adapters (domain/grading/adapters.go)
```

## 🧪 Testing Strategy

### Domain Tests Are Fast & Isolated

Domain services can be tested **without**:
- ❌ Starting an HTTP server
- ❌ Connecting to a real database
- ❌ Running Docker containers
- ❌ Calling external AI APIs

**Example test pattern:**
```go
// domain/grading/service_test.go
func TestSubmitCode_Success(t *testing.T) {
    // Create mocks (no real dependencies!)
    mockExecutor := new(mockCodeExecutor)
    mockEvaluator := new(mockAIEvaluator)
    mockSubmissionRepo := new(mockSubmissionRepository)
    mockExerciseRepo := new(mockExerciseRepository)
    mockProgressRepo := new(mockUserProgressRepository)

    // Configure mock expectations
    mockExerciseRepo.On("GetByID", mock.Anything, "exercise-123").
        Return(&models.Exercise{
            ID:          "exercise-123",
            StarterCode: "function add(a,b){/* TODO */}",
            TestCases:   []models.TestCase{...},
        }, nil)

    mockExecutor.On("ExecuteCodeWithContext", mock.Anything, mock.Anything, "javascript").
        Return(executor.ExecutionResult{Success: true, Output: "All tests passed"})

    // ... more mock setup ...

    // Create service with mocks
    service := NewGradingService(
        mockExecutor,
        mockEvaluator,
        mockSubmissionRepo,
        mockExerciseRepo,
        mockProgressRepo,
    )

    // Execute domain logic
    result, err := service.SubmitCode(context.Background(), SubmitCodeRequest{
        UserID:     "user-123",
        ExerciseID: "exercise-123",
        Code:       "function add(a,b){return a+b}",
        Language:   "javascript",
    }, "request-id")

    // Assert business logic results
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.True(t, result.TotalScore >= 0)
    assert.True(t, result.Passed)

    // Verify all mocks were called correctly
    mockExerciseRepo.AssertExpectations(t)
    mockExecutor.AssertExpectations(t)
}
```

## ✨ Adding New Features

### Example: Adding a "Get Submission History" Feature

#### Step 1: Add Domain Method
```go
// domain/grading/service.go
func (s *GradingService) GetSubmissionHistory(ctx context.Context, exerciseID, userID string, limit int) ([]*SubmissionSummary, error) {
    // Business logic here
    submissions, err := s.submissionRepo.GetByUserAndExercise(ctx, userID, exerciseID, limit)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch submissions: %w", err)
    }

    // Transform to domain model
    summaries := make([]*SubmissionSummary, len(submissions))
    for i, sub := range submissions {
        summaries[i] = &SubmissionSummary{
            ID:         sub.ID,
            Score:      sub.Score,
            SubmittedAt: sub.CreatedAt,
            Passed:     sub.Score >= 80,
        }
    }

    return summaries, nil
}
```

#### Step 2: Add API Handler
```go
// api/grading/handlers.go
func (h *SubmissionHandler) GetSubmissionHistory(c echo.Context) error {
    exerciseID := c.Param("exerciseId")
    userID := getUserID(c) // Hardcoded for MVP

    history, err := h.service.GetSubmissionHistory(
        c.Request().Context(),
        exerciseID,
        userID,
        50, // limit
    )
    if err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
    }

    return c.JSON(http.StatusOK, echo.Map{"submissions": history})
}
```

#### Step 3: Register Route
```go
// main.go
api.GET("/exercise/:exerciseId/submissions", submissionHandler.GetSubmissionHistory)
```

#### Step 4: Write Tests
```go
// domain/grading/service_test.go
func TestGetSubmissionHistory_Success(t *testing.T) {
    mockRepo := new(mockSubmissionRepository)
    mockRepo.On("GetByUserAndExercise", mock.Anything, "user-123", "exercise-123", 50).
        Return([]*models.Submission{...}, nil)

    service := NewGradingService(nil, nil, mockRepo, nil, nil)

    result, err := service.GetSubmissionHistory(context.Background(), "exercise-123", "user-123", 50)

    assert.NoError(t, err)
    assert.Len(t, result, 3)
    mockRepo.AssertExpectations(t)
}

// api/grading/handlers_test.go
func TestGetSubmissionHistory_HappyPath(t *testing.T) {
    mockService := new(mockGradingService)
    mockService.On("GetSubmissionHistory", mock.Anything, "exercise-123", mock.Anything, 50).
        Return([]*grading.SubmissionSummary{...}, nil)

    handler := NewSubmissionHandler(mockService)

    e := echo.New()
    req := httptest.NewRequest(http.MethodGet, "/api/exercise/exercise-123/submissions", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetParamNames("exerciseId")
    c.SetParamValues("exercise-123")

    err := handler.GetSubmissionHistory(c)

    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, rec.Code)
}
```

## 📊 Current Test Coverage

**100% domain + API test coverage** - All business logic is tested!

```
Domain Tests (Pure Unit Tests - Fast):
✅ domain/grading/scorer_test.go       - 35 tests covering scoring algorithm
✅ domain/grading/similarity_test.go   - 15 tests for anti-gaming checks
✅ domain/grading/service_test.go      - 20 tests for submission workflow
✅ domain/learning/service_test.go     - 15 tests for tech tree navigation
✅ domain/content/service_test.go      - 12 tests for exercise management

API Tests (HTTP Handler Tests - Integration):
✅ api/grading/handlers_test.go        - 14 tests for HTTP layer
✅ api/learning/handlers_test.go       - 11 tests for HTTP layer
✅ api/content/handlers_test.go        - 12 tests for HTTP layer

Total: 134 test cases, ~5,200 LOC of test code
```

## 🏗️ Architecture Benefits

### Before Refactoring
- 🔴 **1,719 LOC monolithic handlers file** (`handlers/handlers.go`)
- 🔴 HTTP, business logic, and data access mixed together
- 🔴 Hard to test (required HTTP mocks and database setup)
- 🔴 Tight coupling to Echo framework

### After Refactoring
- ✅ **Domain layer: 1,200+ LOC** of pure, testable business logic
- ✅ **API layer: 560+ LOC** of thin HTTP adapters
- ✅ Clear separation of concerns
- ✅ **Fast tests**: Domain tests run in milliseconds without any external dependencies
- ✅ **Reusable**: Domain services can be used in CLI tools, cron jobs, or other frameworks
- ✅ **Maintainable**: Each domain evolves independently

## 🎓 Further Reading

**Clean Architecture Concepts**:
- [The Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) - Uncle Bob Martin
- [Domain-Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html) - Martin Fowler
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/) - Alistair Cockburn

**Go-Specific Resources**:
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- [Go Best Practices](https://go.dev/doc/effective_go)

**Testing**:
- [Testing in Go](https://go.dev/doc/tutorial/add-a-test)
- [testify/mock Documentation](https://pkg.go.dev/github.com/stretchr/testify/mock)

## 💡 Key Takeaways

1. **Domain layer = Business logic only** - No HTTP, no database, no frameworks
2. **Use interfaces for dependencies** - Enables mocking and flexibility
3. **Test at the domain level first** - Fast, isolated, comprehensive
4. **API layer is thin** - Just HTTP request/response mapping (50-100 LOC)
5. **Each domain is independent** - Can evolve without affecting others

---

**Questions?** See `REFACTORING_STATUS.md` for the complete refactoring journey and test strategy guide.
