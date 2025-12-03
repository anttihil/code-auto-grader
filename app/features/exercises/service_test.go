package exercises

import (
	"context"
	"testing"

	"github.com/anttihil/code-auto-grader/models"
	"github.com/anttihil/code-auto-grader/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ===== Test Main Setup =====


// ===== Mock Implementations =====

type mockExerciseGenerator struct {
	mock.Mock
}

func (m *mockExerciseGenerator) GenerateExerciseWithContext(ctx context.Context, difficulty, topic, language string) (*models.Exercise, error) {
	args := m.Called(ctx, difficulty, topic, language)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Exercise), args.Error(1)
}

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

func (m *mockExerciseRepository) GetBySlug(ctx context.Context, nodeID, slug string) (*models.Exercise, error) {
	args := m.Called(ctx, nodeID, slug)
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

func (m *mockExerciseRepository) Create(ctx context.Context, exercise *models.Exercise) error {
	args := m.Called(ctx, exercise)
	return args.Error(0)
}

func (m *mockExerciseRepository) Update(ctx context.Context, exercise *models.Exercise) error {
	args := m.Called(ctx, exercise)
	return args.Error(0)
}

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

func (m *mockSubmissionRepository) CountAttempts(ctx context.Context, userID, exerciseID string) (int, error) {
	args := m.Called(ctx, userID, exerciseID)
	return args.Int(0), args.Error(1)
}

func (m *mockSubmissionRepository) GetLatestAttempt(ctx context.Context, userID, exerciseID string) (*models.Submission, error) {
	args := m.Called(ctx, userID, exerciseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Submission), args.Error(1)
}

// ===== Test Fixtures =====

func createTestExercise() *models.Exercise {
	exerciseID := uuid.New()
	nodeID := uuid.New()
	return &models.Exercise{
		ID:          exerciseID,
		NodeID:      nodeID,
		Title:       "Test Exercise",
		Description: "A test exercise description",
		Difficulty:  "beginner",
		Content: models.ExerciseContent{
			StarterCode: "function test() { }",
			TestCases: []models.TestCase{
				{Name: "test1", Code: "assert(test() === true)", Visible: true},
			},
			Hints: []models.Hint{
				{Content: "Try using a loop"},
			},
			ExpectedOutcomes: []string{"Function should work correctly"},
		},
	}
}

func createTestSubmission(exerciseID uuid.UUID, score int) *models.Submission {
	return &models.Submission{
		ID:         uuid.New(),
		ExerciseID: exerciseID,
		TotalScore: &score,
	}
}

// ===== Constructor Tests =====

func TestNewContentServiceService(t *testing.T) {
	mockGenerator := new(mockExerciseGenerator)
	mockExerciseRepo := new(mockExerciseRepository)
	mockSubmissionRepo := new(mockSubmissionRepository)

	t.Run("creates service with valid dependencies", func(t *testing.T) {
		service := NewContentService(mockGenerator, mockExerciseRepo, mockSubmissionRepo)

		assert.NotNil(t, service)
		assert.Equal(t, mockGenerator, service.generator)
		assert.Equal(t, mockExerciseRepo, service.exerciseRepo)
		assert.Equal(t, mockSubmissionRepo, service.submissionRepo)
	})

	t.Run("panics with nil generator", func(t *testing.T) {
		assert.Panics(t, func() {
			NewContentService(nil, mockExerciseRepo, mockSubmissionRepo)
		})
	})

	t.Run("panics with nil exercise repository", func(t *testing.T) {
		assert.Panics(t, func() {
			NewContentService(mockGenerator, nil, mockSubmissionRepo)
		})
	})

	t.Run("panics with nil submission repository", func(t *testing.T) {
		assert.Panics(t, func() {
			NewContentService(mockGenerator, mockExerciseRepo, nil)
		})
	})
}

// ===== GetExercise Tests =====

func TestGetExercise_HappyPathService(t *testing.T) {
	mockGenerator := new(mockExerciseGenerator)
	mockExerciseRepo := new(mockExerciseRepository)
	mockSubmissionRepo := new(mockSubmissionRepository)

	service := NewContentService(mockGenerator, mockExerciseRepo, mockSubmissionRepo)

	ctx := context.Background()
	exerciseID := "exercise-123"
	userID := "user-123"
	requestID := "request-123"

	exercise := createTestExercise()
	score1, score2, score3 := 75, 90, 85
	submissions := []*models.Submission{
		createTestSubmission(exercise.ID, score1),
		createTestSubmission(exercise.ID, score2),
		createTestSubmission(exercise.ID, score3),
	}

	mockExerciseRepo.On("GetByID", ctx, exerciseID).Return(exercise, nil)
	mockSubmissionRepo.On("GetByExercisePaginated", ctx, userID, exerciseID, 10, 0).
		Return(submissions, 3, nil)

	result, err := service.GetExercise(ctx, exerciseID, userID, requestID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, exercise, result.Exercise)
	assert.Len(t, result.Submissions, 3)
	assert.Equal(t, 90, result.BestScore) // Best score is 90
	assert.Equal(t, 3, result.AttemptCount)

	mockExerciseRepo.AssertExpectations(t)
	mockSubmissionRepo.AssertExpectations(t)
}

func TestGetExercise_NoSubmissionsService(t *testing.T) {
	mockGenerator := new(mockExerciseGenerator)
	mockExerciseRepo := new(mockExerciseRepository)
	mockSubmissionRepo := new(mockSubmissionRepository)

	service := NewContentService(mockGenerator, mockExerciseRepo, mockSubmissionRepo)

	ctx := context.Background()
	exerciseID := "exercise-123"
	userID := "user-123"
	requestID := "request-123"

	exercise := createTestExercise()

	mockExerciseRepo.On("GetByID", ctx, exerciseID).Return(exercise, nil)
	mockSubmissionRepo.On("GetByExercisePaginated", ctx, userID, exerciseID, 10, 0).
		Return([]*models.Submission{}, 0, nil)

	result, err := service.GetExercise(ctx, exerciseID, userID, requestID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, exercise, result.Exercise)
	assert.Len(t, result.Submissions, 0)
	assert.Equal(t, 0, result.BestScore)
	assert.Equal(t, 0, result.AttemptCount)

	mockExerciseRepo.AssertExpectations(t)
	mockSubmissionRepo.AssertExpectations(t)
}

func TestGetExercise_ExerciseNotFoundService(t *testing.T) {
	mockGenerator := new(mockExerciseGenerator)
	mockExerciseRepo := new(mockExerciseRepository)
	mockSubmissionRepo := new(mockSubmissionRepository)

	service := NewContentService(mockGenerator, mockExerciseRepo, mockSubmissionRepo)

	ctx := context.Background()
	exerciseID := "nonexistent"
	userID := "user-123"
	requestID := "request-123"

	mockExerciseRepo.On("GetByID", ctx, exerciseID).Return(nil, repository.ErrNotFound)

	result, err := service.GetExercise(ctx, exerciseID, userID, requestID)

	assert.Error(t, err)
	assert.Nil(t, result)

	mockExerciseRepo.AssertExpectations(t)
}

func TestGetExercise_SubmissionHistoryErrorService(t *testing.T) {
	mockGenerator := new(mockExerciseGenerator)
	mockExerciseRepo := new(mockExerciseRepository)
	mockSubmissionRepo := new(mockSubmissionRepository)

	service := NewContentService(mockGenerator, mockExerciseRepo, mockSubmissionRepo)

	ctx := context.Background()
	exerciseID := "exercise-123"
	userID := "user-123"
	requestID := "request-123"

	exercise := createTestExercise()

	mockExerciseRepo.On("GetByID", ctx, exerciseID).Return(exercise, nil)
	mockSubmissionRepo.On("GetByExercisePaginated", ctx, userID, exerciseID, 10, 0).
		Return(nil, 0, assert.AnError)

	result, err := service.GetExercise(ctx, exerciseID, userID, requestID)

	// Should continue despite submission history error
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, exercise, result.Exercise)
	assert.Nil(t, result.Submissions)
	assert.Equal(t, 0, result.BestScore)
	assert.Equal(t, 0, result.AttemptCount)

	mockExerciseRepo.AssertExpectations(t)
	mockSubmissionRepo.AssertExpectations(t)
}

func TestGetExercise_BestScoreCalculationService(t *testing.T) {
	mockGenerator := new(mockExerciseGenerator)
	mockExerciseRepo := new(mockExerciseRepository)
	mockSubmissionRepo := new(mockSubmissionRepository)

	service := NewContentService(mockGenerator, mockExerciseRepo, mockSubmissionRepo)

	ctx := context.Background()
	exerciseID := "exercise-123"
	userID := "user-123"
	requestID := "request-123"

	exercise := createTestExercise()

	// Test with nil scores and various values
	score1, score2 := 50, 95
	submissions := []*models.Submission{
		{ID: uuid.New(), ExerciseID: exercise.ID, TotalScore: nil}, // nil score
		createTestSubmission(exercise.ID, score1),
		createTestSubmission(exercise.ID, score2),
		{ID: uuid.New(), ExerciseID: exercise.ID, TotalScore: nil}, // another nil
	}

	mockExerciseRepo.On("GetByID", ctx, exerciseID).Return(exercise, nil)
	mockSubmissionRepo.On("GetByExercisePaginated", ctx, userID, exerciseID, 10, 0).
		Return(submissions, 4, nil)

	result, err := service.GetExercise(ctx, exerciseID, userID, requestID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 95, result.BestScore) // Should pick the highest non-nil score
	assert.Equal(t, 4, result.AttemptCount)

	mockExerciseRepo.AssertExpectations(t)
	mockSubmissionRepo.AssertExpectations(t)
}

// ===== GenerateExercise Tests =====

func TestGenerateExercise_HappyPathService(t *testing.T) {
	mockGenerator := new(mockExerciseGenerator)
	mockExerciseRepo := new(mockExerciseRepository)
	mockSubmissionRepo := new(mockSubmissionRepository)

	service := NewContentService(mockGenerator, mockExerciseRepo, mockSubmissionRepo)

	ctx := context.Background()
	requestID := "request-123"
	req := GenerateExerciseRequest{
		Difficulty: "intermediate",
		Topic:      "algorithms",
		Language:   "python",
	}

	generatedExercise := createTestExercise()
	generatedExercise.Title = "Binary Search"
	generatedExercise.Description = "Implement binary search"
	generatedExercise.Difficulty = "intermediate"

	mockGenerator.On("GenerateExerciseWithContext", ctx, "intermediate", "algorithms", "python").
		Return(generatedExercise, nil)

	result, err := service.GenerateExercise(ctx, req, requestID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Binary Search", result.Title)
	assert.Equal(t, "Implement binary search", result.Description)
	assert.Equal(t, "intermediate", result.Difficulty)
	assert.Equal(t, "python", result.Language)
	assert.NotEmpty(t, result.StarterCode)
	assert.Len(t, result.Hints, 1)
	assert.Equal(t, "Try using a loop", result.Hints[0])
	assert.Equal(t, "Function should work correctly", result.ExpectedOutput)
	assert.Len(t, result.TestCases, 1)

	mockGenerator.AssertExpectations(t)
}

func TestGenerateExercise_DefaultValuesService(t *testing.T) {
	mockGenerator := new(mockExerciseGenerator)
	mockExerciseRepo := new(mockExerciseRepository)
	mockSubmissionRepo := new(mockSubmissionRepository)

	service := NewContentService(mockGenerator, mockExerciseRepo, mockSubmissionRepo)

	ctx := context.Background()
	requestID := "request-123"
	req := GenerateExerciseRequest{
		// All fields empty - should use defaults
	}

	generatedExercise := createTestExercise()

	mockGenerator.On("GenerateExerciseWithContext", ctx, "beginner", "general", "javascript").
		Return(generatedExercise, nil)

	result, err := service.GenerateExercise(ctx, req, requestID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "javascript", result.Language)

	mockGenerator.AssertExpectations(t)
}

func TestGenerateExercise_GeneratorErrorService(t *testing.T) {
	mockGenerator := new(mockExerciseGenerator)
	mockExerciseRepo := new(mockExerciseRepository)
	mockSubmissionRepo := new(mockSubmissionRepository)

	service := NewContentService(mockGenerator, mockExerciseRepo, mockSubmissionRepo)

	ctx := context.Background()
	requestID := "request-123"
	req := GenerateExerciseRequest{
		Difficulty: "expert",
		Topic:      "quantum",
		Language:   "rust",
	}

	mockGenerator.On("GenerateExerciseWithContext", ctx, "expert", "quantum", "rust").
		Return(nil, assert.AnError)

	result, err := service.GenerateExercise(ctx, req, requestID)

	assert.Error(t, err)
	assert.Nil(t, result)

	mockGenerator.AssertExpectations(t)
}

func TestGenerateExercise_EmptyHintsService(t *testing.T) {
	mockGenerator := new(mockExerciseGenerator)
	mockExerciseRepo := new(mockExerciseRepository)
	mockSubmissionRepo := new(mockSubmissionRepository)

	service := NewContentService(mockGenerator, mockExerciseRepo, mockSubmissionRepo)

	ctx := context.Background()
	requestID := "request-123"
	req := GenerateExerciseRequest{
		Difficulty: "beginner",
		Topic:      "basics",
		Language:   "javascript",
	}

	generatedExercise := createTestExercise()
	generatedExercise.Content.Hints = []models.Hint{} // Empty hints
	generatedExercise.Content.ExpectedOutcomes = []string{} // Empty outcomes

	mockGenerator.On("GenerateExerciseWithContext", ctx, "beginner", "basics", "javascript").
		Return(generatedExercise, nil)

	result, err := service.GenerateExercise(ctx, req, requestID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Hints, 0)
	assert.Equal(t, "", result.ExpectedOutput)

	mockGenerator.AssertExpectations(t)
}

// ===== GetSampleExercise Tests =====

func TestGetSampleExerciseService(t *testing.T) {
	mockGenerator := new(mockExerciseGenerator)
	mockExerciseRepo := new(mockExerciseRepository)
	mockSubmissionRepo := new(mockSubmissionRepository)

	service := NewContentService(mockGenerator, mockExerciseRepo, mockSubmissionRepo)

	result := service.GetSampleExercise()

	assert.NotNil(t, result)
	assert.Equal(t, "Sum Two Numbers", result.Title)
	assert.Equal(t, "beginner", result.Difficulty)
	assert.Equal(t, "javascript", result.Language)
	assert.NotEmpty(t, result.StarterCode)
	assert.Len(t, result.Hints, 3)
	assert.Len(t, result.TestCases, 3)
	assert.Contains(t, result.Description, "sum")
}
