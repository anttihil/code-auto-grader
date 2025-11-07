package grading

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/anttihil/code-auto-grader/internal/logger"
	"github.com/anttihil/code-auto-grader/models"
	"github.com/anttihil/code-auto-grader/repository"
	"github.com/anttihil/code-auto-grader/services/executor"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestMain sets up test environment
func TestMain(m *testing.M) {
	// Initialize logger for tests
	logger.Initialize("test")

	// Run tests
	code := m.Run()

	os.Exit(code)
}

// ============================================================================
// Mock Implementations
// ============================================================================

// mockCodeExecutor mocks the CodeExecutor interface
type mockCodeExecutor struct {
	mock.Mock
}

func (m *mockCodeExecutor) ExecuteCodeWithContext(ctx context.Context, code, language string) *executor.ExecutionResult {
	args := m.Called(ctx, code, language)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*executor.ExecutionResult)
}

func (m *mockCodeExecutor) GetPoolStatus() map[string]interface{} {
	args := m.Called()
	return args.Get(0).(map[string]interface{})
}

// mockAIEvaluator mocks the AIEvaluator interface
type mockAIEvaluator struct {
	mock.Mock
}

func (m *mockAIEvaluator) CheckSubmissionWithContext(ctx context.Context, exercise, code, language string) (*Feedback, error) {
	args := m.Called(ctx, exercise, code, language)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Feedback), args.Error(1)
}

// mockSubmissionRepository mocks the SubmissionRepository interface
type mockSubmissionRepository struct {
	mock.Mock
}

func (m *mockSubmissionRepository) Create(ctx context.Context, submission *models.Submission) error {
	args := m.Called(ctx, submission)
	return args.Error(0)
}

func (m *mockSubmissionRepository) GetByID(ctx context.Context, id string) (*models.Submission, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Submission), args.Error(1)
}

func (m *mockSubmissionRepository) GetByExercise(ctx context.Context, userID, exerciseID string) ([]*models.Submission, error) {
	args := m.Called(ctx, userID, exerciseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Submission), args.Error(1)
}

func (m *mockSubmissionRepository) GetByExercisePaginated(ctx context.Context, userID, exerciseID string, limit, offset int) ([]*models.Submission, int, error) {
	args := m.Called(ctx, userID, exerciseID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*models.Submission), args.Int(1), args.Error(2)
}

func (m *mockSubmissionRepository) GetLatestAttempt(ctx context.Context, userID, exerciseID string) (*models.Submission, error) {
	args := m.Called(ctx, userID, exerciseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Submission), args.Error(1)
}

func (m *mockSubmissionRepository) CountAttempts(ctx context.Context, userID, exerciseID string) (int, error) {
	args := m.Called(ctx, userID, exerciseID)
	return args.Int(0), args.Error(1)
}

// mockExerciseRepository mocks the ExerciseRepository interface
type mockExerciseRepository struct {
	mock.Mock
}

func (m *mockExerciseRepository) GetByID(ctx context.Context, id string) (*models.Exercise, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Exercise), args.Error(1)
}

func (m *mockExerciseRepository) GetByNode(ctx context.Context, nodeID string) ([]*models.Exercise, error) {
	args := m.Called(ctx, nodeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Exercise), args.Error(1)
}

func (m *mockExerciseRepository) GetBySlug(ctx context.Context, nodeID, slug string) (*models.Exercise, error) {
	args := m.Called(ctx, nodeID, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Exercise), args.Error(1)
}

func (m *mockExerciseRepository) Create(ctx context.Context, exercise *models.Exercise) error {
	args := m.Called(ctx, exercise)
	return args.Error(0)
}

func (m *mockExerciseRepository) Update(ctx context.Context, exercise *models.Exercise) error {
	args := m.Called(ctx, exercise)
	return args.Error(0)
}

// mockUserProgressRepository mocks the UserProgressRepository interface
type mockUserProgressRepository struct {
	mock.Mock
}

func (m *mockUserProgressRepository) GetByUserAndNode(ctx context.Context, userID, nodeID string) (*models.UserProgress, error) {
	args := m.Called(ctx, userID, nodeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserProgress), args.Error(1)
}

func (m *mockUserProgressRepository) Upsert(ctx context.Context, progress *models.UserProgress) error {
	args := m.Called(ctx, progress)
	return args.Error(0)
}

func (m *mockUserProgressRepository) GetByUser(ctx context.Context, userID string) ([]*models.UserProgress, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.UserProgress), args.Error(1)
}

func (m *mockUserProgressRepository) UpdateMastery(ctx context.Context, userID, nodeID string, masteryScore float64, isMastered bool) error {
	args := m.Called(ctx, userID, nodeID, masteryScore, isMastered)
	return args.Error(0)
}

// ============================================================================
// Helper Functions for Tests
// ============================================================================

func createTestExercise() *models.Exercise {
	nodeID := uuid.New()
	exerciseID := uuid.New()

	return &models.Exercise{
		ID:          exerciseID,
		NodeID:      nodeID,
		Title:       "Test Exercise",
		Description: "A test exercise for unit testing",
		Difficulty:  "easy",
		Content: models.ExerciseContent{
			StarterCode: "function add(a, b) {\n  // TODO: implement\n  return 0;\n}",
			TestCases: []models.TestCase{
				{
					Name:    "Test 1+1",
					Code:    "assert.strictEqual(add(1, 1), 2);",
					Visible: true,
				},
				{
					Name:    "Test 2+3",
					Code:    "assert.strictEqual(add(2, 3), 5);",
					Visible: true,
				},
			},
			Hints: []models.Hint{
				{Order: 1, Title: "Addition", Content: "Use the + operator"},
			},
			CommonMistakes: []models.CommonMistake{},
			GradingRubric: models.GradingRubric{
				TestWeight: 60,
				AIWeight:   40,
				AICriteria: []models.AICriterion{
					{Name: "Correctness", Points: 50, Description: "Code correctness"},
					{Name: "Code Quality", Points: 50, Description: "Code quality"},
				},
			},
			AIEvaluationPrompt: "Evaluate the code quality and correctness",
		},
	}
}

func createTestSubmission(userID, exerciseID string, score int) *models.Submission {
	submissionID := uuid.New()
	execTime := 150
	execOutput := "Test output"
	aiEvalTime := 2000

	return &models.Submission{
		ID:                 submissionID,
		UserID:             uuid.MustParse(userID),
		ExerciseID:         uuid.MustParse(exerciseID),
		SubmittedCode:      "function add(a, b) { return a + b; }",
		Language:           "javascript",
		AttemptNumber:      1,
		ExecutionTimeMs:    &execTime,
		ExecutionOutput:    &execOutput,
		ExecutionError:     nil,
		AIEvaluationTimeMs: &aiEvalTime,
		TestScore:          &score,
		AIScore:            &score,
		TotalScore:         &score,
		TestResults:        &models.TestResults{TotalTests: 2, PassedTests: 2, FailedTests: 0},
		AIFeedback:         &models.AIFeedback{Passed: true, Score: score, Feedback: "Good job!"},
	}
}

// ============================================================================
// NewGradingService Tests
// ============================================================================

func TestNewGradingService(t *testing.T) {
	mockExecutor := new(mockCodeExecutor)
	mockEvaluator := new(mockAIEvaluator)
	mockSubmissionRepo := new(mockSubmissionRepository)
	mockExerciseRepo := new(mockExerciseRepository)
	mockProgressRepo := new(mockUserProgressRepository)

	t.Run("creates service with all dependencies", func(t *testing.T) {
		service := NewGradingService(mockExecutor, mockEvaluator, mockSubmissionRepo, mockExerciseRepo, mockProgressRepo)
		assert.NotNil(t, service)
		assert.Equal(t, mockExecutor, service.executor)
		assert.Equal(t, mockEvaluator, service.evaluator)
		assert.Equal(t, mockSubmissionRepo, service.submissionRepo)
		assert.Equal(t, mockExerciseRepo, service.exerciseRepo)
		assert.Equal(t, mockProgressRepo, service.userProgressRepo)
	})

	t.Run("panics when executor is nil", func(t *testing.T) {
		assert.Panics(t, func() {
			NewGradingService(nil, mockEvaluator, mockSubmissionRepo, mockExerciseRepo, mockProgressRepo)
		})
	})

	t.Run("panics when evaluator is nil", func(t *testing.T) {
		assert.Panics(t, func() {
			NewGradingService(mockExecutor, nil, mockSubmissionRepo, mockExerciseRepo, mockProgressRepo)
		})
	})

	t.Run("panics when submissionRepo is nil", func(t *testing.T) {
		assert.Panics(t, func() {
			NewGradingService(mockExecutor, mockEvaluator, nil, mockExerciseRepo, mockProgressRepo)
		})
	})

	t.Run("panics when exerciseRepo is nil", func(t *testing.T) {
		assert.Panics(t, func() {
			NewGradingService(mockExecutor, mockEvaluator, mockSubmissionRepo, nil, mockProgressRepo)
		})
	})

	t.Run("panics when progressRepo is nil", func(t *testing.T) {
		assert.Panics(t, func() {
			NewGradingService(mockExecutor, mockEvaluator, mockSubmissionRepo, mockExerciseRepo, nil)
		})
	})
}

// ============================================================================
// SubmitCode Tests
// ============================================================================

func TestSubmitCode_HappyPath(t *testing.T) {
	// Setup mocks
	mockExecutor := new(mockCodeExecutor)
	mockEvaluator := new(mockAIEvaluator)
	mockSubmissionRepo := new(mockSubmissionRepository)
	mockExerciseRepo := new(mockExerciseRepository)
	mockProgressRepo := new(mockUserProgressRepository)

	service := NewGradingService(mockExecutor, mockEvaluator, mockSubmissionRepo, mockExerciseRepo, mockProgressRepo)

	ctx := context.Background()
	userID := "00000000-0000-0000-0000-000000000001"
	exercise := createTestExercise()
	exerciseID := exercise.ID.String()
	code := "function add(a, b) { return a + b; }"

	// Mock exercise repository
	mockExerciseRepo.On("GetByID", ctx, exerciseID).Return(exercise, nil)

	// Mock code execution  - return value matches interface return type
	execResult := &executor.ExecutionResult{
		Output: "2\n5",
		Error:  "",
	}
	mockExecutor.On("ExecuteCodeWithContext", ctx, mock.Anything, "javascript").Return(execResult)

	// Mock AI evaluation
	mockEvaluator.On("CheckSubmissionWithContext", ctx, mock.Anything, code, "javascript").Return(&Feedback{
		Passed:       true,
		Score:        90,
		Feedback:     "Excellent implementation!",
		Strengths:    []string{"Clean code", "Correct logic"},
		Improvements: []string{},
		Correctness:  "Perfect",
		CodeQuality:  "Great",
		Suggestions:  "None",
	}, nil)

	// Mock submission count (for attempt number)
	mockSubmissionRepo.On("CountAttempts", ctx, userID, exerciseID).Return(0, nil)

	// Mock submission creation
	mockSubmissionRepo.On("Create", ctx, mock.AnythingOfType("*models.Submission")).Return(nil)

	// Mock progress repository (no existing progress, no error - will create new)
	mockProgressRepo.On("GetByUserAndNode", ctx, userID, exercise.NodeID.String()).Return(
		nil, nil,
	)
	mockProgressRepo.On("Upsert", ctx, mock.AnythingOfType("*models.UserProgress")).Return(nil)

	// Execute
	req := SubmitCodeRequest{
		UserID:     userID,
		ExerciseID: exerciseID,
		Code:       code,
		Language:   "javascript",
	}

	result, err := service.SubmitCode(ctx, req, "test-request-id")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.SubmissionID)
	assert.Equal(t, 1, result.AttemptNumber)
	assert.NotNil(t, result.TestResults)
	assert.Equal(t, 2, result.TestResults.PassedTests)
	assert.Equal(t, 2, result.TestResults.TotalTests)
	assert.NotNil(t, result.AIFeedback)
	assert.True(t, result.AIFeedback.Passed)
	assert.Greater(t, result.TotalScore, 0)

	// Verify all mocks were called as expected
	mockExerciseRepo.AssertExpectations(t)
	mockExecutor.AssertExpectations(t)
	mockEvaluator.AssertExpectations(t)
	mockSubmissionRepo.AssertExpectations(t)
	mockProgressRepo.AssertExpectations(t)
}

func TestSubmitCode_SimilarityRejection_Identical(t *testing.T) {
	// Setup mocks
	mockExecutor := new(mockCodeExecutor)
	mockEvaluator := new(mockAIEvaluator)
	mockSubmissionRepo := new(mockSubmissionRepository)
	mockExerciseRepo := new(mockExerciseRepository)
	mockProgressRepo := new(mockUserProgressRepository)

	service := NewGradingService(mockExecutor, mockEvaluator, mockSubmissionRepo, mockExerciseRepo, mockProgressRepo)

	ctx := context.Background()
	userID := "00000000-0000-0000-0000-000000000001"
	exercise := createTestExercise()
	exerciseID := exercise.ID.String()

	// Submit code that is identical to starter code
	code := exercise.Content.StarterCode

	// Mock exercise repository
	mockExerciseRepo.On("GetByID", ctx, exerciseID).Return(exercise, nil)

	// Execute
	req := SubmitCodeRequest{
		UserID:     userID,
		ExerciseID: exerciseID,
		Code:       code,
		Language:   "javascript",
	}

	result, err := service.SubmitCode(ctx, req, "test-request-id")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "IDENTICAL_TO_STARTER")

	// Verify only exercise fetch was called
	mockExerciseRepo.AssertExpectations(t)
	mockExecutor.AssertNotCalled(t, "ExecuteCodeWithContext")
	mockEvaluator.AssertNotCalled(t, "CheckSubmissionWithContext")
	mockSubmissionRepo.AssertNotCalled(t, "Create")
}

func TestSubmitCode_SimilarityRejection_TrivialChanges(t *testing.T) {
	// Setup mocks
	mockExecutor := new(mockCodeExecutor)
	mockEvaluator := new(mockAIEvaluator)
	mockSubmissionRepo := new(mockSubmissionRepository)
	mockExerciseRepo := new(mockExerciseRepository)
	mockProgressRepo := new(mockUserProgressRepository)

	service := NewGradingService(mockExecutor, mockEvaluator, mockSubmissionRepo, mockExerciseRepo, mockProgressRepo)

	ctx := context.Background()
	userID := "00000000-0000-0000-0000-000000000001"
	exercise := createTestExercise()
	exerciseID := exercise.ID.String()

	// Submit code with only comments added (trivial changes)
	code := "// This is my solution\n" + exercise.Content.StarterCode

	// Mock exercise repository
	mockExerciseRepo.On("GetByID", ctx, exerciseID).Return(exercise, nil)

	// Execute
	req := SubmitCodeRequest{
		UserID:     userID,
		ExerciseID: exerciseID,
		Code:       code,
		Language:   "javascript",
	}

	result, err := service.SubmitCode(ctx, req, "test-request-id")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "TRIVIAL_CHANGES")

	// Verify only exercise fetch was called
	mockExerciseRepo.AssertExpectations(t)
	mockExecutor.AssertNotCalled(t, "ExecuteCodeWithContext")
	mockEvaluator.AssertNotCalled(t, "CheckSubmissionWithContext")
	mockSubmissionRepo.AssertNotCalled(t, "Create")
}

func TestSubmitCode_ExerciseNotFound(t *testing.T) {
	// Setup mocks
	mockExecutor := new(mockCodeExecutor)
	mockEvaluator := new(mockAIEvaluator)
	mockSubmissionRepo := new(mockSubmissionRepository)
	mockExerciseRepo := new(mockExerciseRepository)
	mockProgressRepo := new(mockUserProgressRepository)

	service := NewGradingService(mockExecutor, mockEvaluator, mockSubmissionRepo, mockExerciseRepo, mockProgressRepo)

	ctx := context.Background()
	userID := "00000000-0000-0000-0000-000000000001"
	exerciseID := "nonexistent-exercise-id"

	// Mock exercise repository to return error
	mockExerciseRepo.On("GetByID", ctx, exerciseID).Return(nil, repository.ErrNotFound)

	// Execute
	req := SubmitCodeRequest{
		UserID:     userID,
		ExerciseID: exerciseID,
		Code:       "some code",
		Language:   "javascript",
	}

	result, err := service.SubmitCode(ctx, req, "test-request-id")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "exercise not found")

	// Verify
	mockExerciseRepo.AssertExpectations(t)
}

func TestSubmitCode_AIEvaluationFailure_ContinuesWithFallback(t *testing.T) {
	// Setup mocks
	mockExecutor := new(mockCodeExecutor)
	mockEvaluator := new(mockAIEvaluator)
	mockSubmissionRepo := new(mockSubmissionRepository)
	mockExerciseRepo := new(mockExerciseRepository)
	mockProgressRepo := new(mockUserProgressRepository)

	service := NewGradingService(mockExecutor, mockEvaluator, mockSubmissionRepo, mockExerciseRepo, mockProgressRepo)

	ctx := context.Background()
	userID := "00000000-0000-0000-0000-000000000001"
	exercise := createTestExercise()
	exerciseID := exercise.ID.String()
	code := "function add(a, b) { return a + b; }"

	// Mock exercise repository
	mockExerciseRepo.On("GetByID", ctx, exerciseID).Return(exercise, nil)

	// Mock code execution (successful) - return value matches interface return type
	execResult := &executor.ExecutionResult{
		Output: "2\n5",
		Error:  "",
	}
	mockExecutor.On("ExecuteCodeWithContext", ctx, mock.Anything, "javascript").Return(execResult)

	// Mock AI evaluation to fail
	mockEvaluator.On("CheckSubmissionWithContext", ctx, mock.Anything, code, "javascript").Return(
		nil, errors.New("AI service unavailable"),
	)

	// Mock submission count
	mockSubmissionRepo.On("CountAttempts", ctx, userID, exerciseID).Return(0, nil)

	// Mock submission creation
	mockSubmissionRepo.On("Create", ctx, mock.AnythingOfType("*models.Submission")).Return(nil)

	// Mock progress repository (no existing progress, no error - will create new)
	mockProgressRepo.On("GetByUserAndNode", ctx, userID, exercise.NodeID.String()).Return(
		nil, nil,
	)
	mockProgressRepo.On("Upsert", ctx, mock.AnythingOfType("*models.UserProgress")).Return(nil)

	// Execute
	req := SubmitCodeRequest{
		UserID:     userID,
		ExerciseID: exerciseID,
		Code:       code,
		Language:   "javascript",
	}

	result, err := service.SubmitCode(ctx, req, "test-request-id")

	// Assert - should succeed despite AI failure
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.AIFeedback)
	assert.Equal(t, 0, result.AIScore) // AI score should be 0 when AI fails
	assert.Equal(t, "AI evaluation temporarily unavailable", result.AIFeedback.Feedback)
	assert.Equal(t, 60, result.TotalScore) // Should be based only on test score (100% tests passed, 60% weight)

	// Verify
	mockExerciseRepo.AssertExpectations(t)
	mockExecutor.AssertExpectations(t)
	mockEvaluator.AssertExpectations(t)
	mockSubmissionRepo.AssertExpectations(t)
	mockProgressRepo.AssertExpectations(t)
}

// ============================================================================
// GetPoolStatus Tests
// ============================================================================

func TestGetPoolStatus(t *testing.T) {
	mockExecutor := new(mockCodeExecutor)
	mockEvaluator := new(mockAIEvaluator)
	mockSubmissionRepo := new(mockSubmissionRepository)
	mockExerciseRepo := new(mockExerciseRepository)
	mockProgressRepo := new(mockUserProgressRepository)

	service := NewGradingService(mockExecutor, mockEvaluator, mockSubmissionRepo, mockExerciseRepo, mockProgressRepo)

	expectedStatus := map[string]interface{}{
		"javascript": map[string]interface{}{
			"size":      5,
			"available": 5,
		},
		"python": map[string]interface{}{
			"size":      5,
			"available": 4,
		},
	}

	mockExecutor.On("GetPoolStatus").Return(expectedStatus)

	status := service.GetPoolStatus()

	assert.Equal(t, expectedStatus, status)
	mockExecutor.AssertExpectations(t)
}

// ============================================================================
// RunCode Tests
// ============================================================================

func TestRunCode_HappyPath(t *testing.T) {
	// Setup mocks
	mockExecutor := new(mockCodeExecutor)
	mockEvaluator := new(mockAIEvaluator)
	mockSubmissionRepo := new(mockSubmissionRepository)
	mockExerciseRepo := new(mockExerciseRepository)
	mockProgressRepo := new(mockUserProgressRepository)

	service := NewGradingService(mockExecutor, mockEvaluator, mockSubmissionRepo, mockExerciseRepo, mockProgressRepo)

	ctx := context.Background()
	exercise := createTestExercise()
	exerciseID := exercise.ID.String()
	code := "function add(a, b) { return a + b; }"

	// Mock exercise repository
	mockExerciseRepo.On("GetByID", ctx, exerciseID).Return(exercise, nil)

	// Mock code execution - all tests pass
	execResult := &executor.ExecutionResult{
		Output: "2\n5",
		Error:  "",
	}
	mockExecutor.On("ExecuteCodeWithContext", ctx, mock.Anything, "javascript").Return(execResult)

	// Execute
	req := RunCodeRequest{
		ExerciseID: exerciseID,
		Code:       code,
		Language:   "javascript",
	}

	result, err := service.RunCode(ctx, req, "test-request-id")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Contains(t, result.ExecutionOutput, "PASSED") // Contains formatted test output
	assert.Contains(t, result.ExecutionOutput, "Test 1+1")
	assert.Nil(t, result.ExecutionError)
	assert.NotNil(t, result.TestResults)
	assert.Equal(t, 2, result.TestResults.PassedTests)
	assert.Equal(t, 2, result.TestResults.TotalTests)
	assert.Equal(t, 100, result.TestScore)

	// Verify mocks - note that AI and submission repos should NOT be called
	mockExerciseRepo.AssertExpectations(t)
	mockExecutor.AssertExpectations(t)
	mockEvaluator.AssertNotCalled(t, "CheckSubmissionWithContext")
	mockSubmissionRepo.AssertNotCalled(t, "Create")
	mockProgressRepo.AssertNotCalled(t, "Upsert")
}

func TestRunCode_PartialTestPass(t *testing.T) {
	// Setup mocks
	mockExecutor := new(mockCodeExecutor)
	mockEvaluator := new(mockAIEvaluator)
	mockSubmissionRepo := new(mockSubmissionRepository)
	mockExerciseRepo := new(mockExerciseRepository)
	mockProgressRepo := new(mockUserProgressRepository)

	service := NewGradingService(mockExecutor, mockEvaluator, mockSubmissionRepo, mockExerciseRepo, mockProgressRepo)

	ctx := context.Background()
	exercise := createTestExercise()
	exerciseID := exercise.ID.String()
	code := "function add(a, b) { return a + b + 1; }" // Wrong implementation

	// Mock exercise repository
	mockExerciseRepo.On("GetByID", ctx, exerciseID).Return(exercise, nil)

	// Mock code execution - only first test passes
	execResult := &executor.ExecutionResult{
		Output: "2\nAssertionError: expected 6 to equal 5",
		Error:  "Test failed",
	}
	mockExecutor.On("ExecuteCodeWithContext", ctx, mock.Anything, "javascript").Return(execResult)

	// Execute
	req := RunCodeRequest{
		ExerciseID: exerciseID,
		Code:       code,
		Language:   "javascript",
	}

	result, err := service.RunCode(ctx, req, "test-request-id")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.ExecutionError)
	assert.NotNil(t, result.TestResults)
	assert.Less(t, result.TestScore, 100) // Not all tests passed

	// Verify
	mockExerciseRepo.AssertExpectations(t)
	mockExecutor.AssertExpectations(t)
}

func TestRunCode_ExerciseNotFound(t *testing.T) {
	// Setup mocks
	mockExecutor := new(mockCodeExecutor)
	mockEvaluator := new(mockAIEvaluator)
	mockSubmissionRepo := new(mockSubmissionRepository)
	mockExerciseRepo := new(mockExerciseRepository)
	mockProgressRepo := new(mockUserProgressRepository)

	service := NewGradingService(mockExecutor, mockEvaluator, mockSubmissionRepo, mockExerciseRepo, mockProgressRepo)

	ctx := context.Background()
	exerciseID := "nonexistent-exercise-id"

	// Mock exercise repository to return error
	mockExerciseRepo.On("GetByID", ctx, exerciseID).Return(nil, repository.ErrNotFound)

	// Execute
	req := RunCodeRequest{
		ExerciseID: exerciseID,
		Code:       "some code",
		Language:   "javascript",
	}

	result, err := service.RunCode(ctx, req, "test-request-id")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "exercise not found")

	// Verify
	mockExerciseRepo.AssertExpectations(t)
	mockExecutor.AssertNotCalled(t, "ExecuteCodeWithContext")
}

func TestRunCode_ExecutionError(t *testing.T) {
	// Setup mocks
	mockExecutor := new(mockCodeExecutor)
	mockEvaluator := new(mockAIEvaluator)
	mockSubmissionRepo := new(mockSubmissionRepository)
	mockExerciseRepo := new(mockExerciseRepository)
	mockProgressRepo := new(mockUserProgressRepository)

	service := NewGradingService(mockExecutor, mockEvaluator, mockSubmissionRepo, mockExerciseRepo, mockProgressRepo)

	ctx := context.Background()
	exercise := createTestExercise()
	exerciseID := exercise.ID.String()
	code := "function add(a, b) { return a +" // Syntax error

	// Mock exercise repository
	mockExerciseRepo.On("GetByID", ctx, exerciseID).Return(exercise, nil)

	// Mock code execution with syntax error
	execResult := &executor.ExecutionResult{
		Output: "",
		Error:  "SyntaxError: Unexpected end of input",
	}
	mockExecutor.On("ExecuteCodeWithContext", ctx, mock.Anything, "javascript").Return(execResult)

	// Execute
	req := RunCodeRequest{
		ExerciseID: exerciseID,
		Code:       code,
		Language:   "javascript",
	}

	result, err := service.RunCode(ctx, req, "test-request-id")

	// Assert - should succeed but with error in result
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.ExecutionError)
	assert.Contains(t, *result.ExecutionError, "SyntaxError")
	assert.Equal(t, 0, result.TestScore) // No tests passed

	// Verify
	mockExerciseRepo.AssertExpectations(t)
	mockExecutor.AssertExpectations(t)
}
