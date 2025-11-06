package learning

import (
	"context"
	"os"
	"testing"

	"github.com/anttihil/code-auto-grader/internal/logger"
	"github.com/anttihil/code-auto-grader/models"
	"github.com/anttihil/code-auto-grader/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ===== Test Main Setup =====

func TestMain(m *testing.M) {
	logger.Initialize("test")
	code := m.Run()
	os.Exit(code)
}

// ===== Mock Implementations =====

type mockTechTreeRepository struct {
	mock.Mock
}

func (m *mockTechTreeRepository) GetByID(ctx context.Context, id string) (*models.TechTree, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TechTree), args.Error(1)
}

func (m *mockTechTreeRepository) List(ctx context.Context) ([]*models.TechTree, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.TechTree), args.Error(1)
}

func (m *mockTechTreeRepository) Create(ctx context.Context, tree *models.TechTree) error {
	args := m.Called(ctx, tree)
	return args.Error(0)
}

func (m *mockTechTreeRepository) GetBySlug(ctx context.Context, slug string) (*models.TechTree, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TechTree), args.Error(1)
}

func (m *mockTechTreeRepository) Update(ctx context.Context, tree *models.TechTree) error {
	args := m.Called(ctx, tree)
	return args.Error(0)
}

type mockNodeRepository struct {
	mock.Mock
}

func (m *mockNodeRepository) GetByID(ctx context.Context, id string) (*models.Node, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Node), args.Error(1)
}

func (m *mockNodeRepository) GetByTechTree(ctx context.Context, treeID string) ([]*models.Node, error) {
	args := m.Called(ctx, treeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Node), args.Error(1)
}

func (m *mockNodeRepository) Create(ctx context.Context, node *models.Node) error {
	args := m.Called(ctx, node)
	return args.Error(0)
}

func (m *mockNodeRepository) GetPrerequisites(ctx context.Context, nodeID string) ([]string, error) {
	args := m.Called(ctx, nodeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockNodeRepository) GetPrerequisitesBatch(ctx context.Context, nodeIDs []string) (map[string][]string, error) {
	args := m.Called(ctx, nodeIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string][]string), args.Error(1)
}

func (m *mockNodeRepository) Update(ctx context.Context, node *models.Node) error {
	args := m.Called(ctx, node)
	return args.Error(0)
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

type mockUserProgressRepository struct {
	mock.Mock
}

func (m *mockUserProgressRepository) GetByUser(ctx context.Context, userID string) ([]*models.UserProgress, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.UserProgress), args.Error(1)
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

func (m *mockUserProgressRepository) UpdateMastery(ctx context.Context, userID, nodeID string, masteryScore float64, isMastered bool) error {
	args := m.Called(ctx, userID, nodeID, masteryScore, isMastered)
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

func createTestTechTree() *models.TechTree {
	treeID := uuid.New()
	description := "A test tech tree for unit tests"
	return &models.TechTree{
		ID:          treeID,
		Name:        "Test Tech Tree",
		Description: &description,
		Slug:        "test-tree",
		Language:    "javascript",
		IsPublished: true,
	}
}

func createTestNode(treeID uuid.UUID, order int) *models.Node {
	nodeID := uuid.New()
	description := "A test node"
	return &models.Node{
		ID:                 nodeID,
		TechTreeID:         treeID,
		Name:               "Test Node",
		Description:        &description,
		Slug:               "test-node",
		OrderIndex:         order,
		LearningObjectives: []string{"Learn Go", "Master testing"},
	}
}

func createTestExercise(nodeID uuid.UUID) *models.Exercise {
	exerciseID := uuid.New()
	return &models.Exercise{
		ID:          exerciseID,
		NodeID:      nodeID,
		Title:       "Test Exercise",
		Description: "A test exercise",
		Difficulty:  "beginner",
	}
}

func createTestUserProgress(nodeID uuid.UUID, completed, total int, mastered bool) *models.UserProgress {
	return &models.UserProgress{
		NodeID:             nodeID,
		ExercisesCompleted: completed,
		ExercisesTotal:     total,
		IsMastered:         mastered,
	}
}

// ===== Constructor Tests =====

func TestNewLearningService(t *testing.T) {
	mockTechTreeRepo := new(mockTechTreeRepository)
	mockNodeRepo := new(mockNodeRepository)
	mockExerciseRepo := new(mockExerciseRepository)
	mockProgressRepo := new(mockUserProgressRepository)
	mockSubmissionRepo := new(mockSubmissionRepository)

	t.Run("creates service with valid dependencies", func(t *testing.T) {
		service := NewLearningService(
			mockTechTreeRepo,
			mockNodeRepo,
			mockExerciseRepo,
			mockProgressRepo,
			mockSubmissionRepo,
		)

		assert.NotNil(t, service)
		assert.Equal(t, mockTechTreeRepo, service.techTreeRepo)
		assert.Equal(t, mockNodeRepo, service.nodeRepo)
		assert.Equal(t, mockExerciseRepo, service.exerciseRepo)
		assert.Equal(t, mockProgressRepo, service.userProgressRepo)
		assert.Equal(t, mockSubmissionRepo, service.submissionRepo)
	})

	t.Run("panics with nil tech tree repository", func(t *testing.T) {
		assert.Panics(t, func() {
			NewLearningService(nil, mockNodeRepo, mockExerciseRepo, mockProgressRepo, mockSubmissionRepo)
		})
	})

	t.Run("panics with nil node repository", func(t *testing.T) {
		assert.Panics(t, func() {
			NewLearningService(mockTechTreeRepo, nil, mockExerciseRepo, mockProgressRepo, mockSubmissionRepo)
		})
	})

	t.Run("panics with nil exercise repository", func(t *testing.T) {
		assert.Panics(t, func() {
			NewLearningService(mockTechTreeRepo, mockNodeRepo, nil, mockProgressRepo, mockSubmissionRepo)
		})
	})

	t.Run("panics with nil user progress repository", func(t *testing.T) {
		assert.Panics(t, func() {
			NewLearningService(mockTechTreeRepo, mockNodeRepo, mockExerciseRepo, nil, mockSubmissionRepo)
		})
	})

	t.Run("panics with nil submission repository", func(t *testing.T) {
		assert.Panics(t, func() {
			NewLearningService(mockTechTreeRepo, mockNodeRepo, mockExerciseRepo, mockProgressRepo, nil)
		})
	})
}

// ===== GetTechTree Tests =====

func TestGetTechTree_HappyPath(t *testing.T) {
	mockTechTreeRepo := new(mockTechTreeRepository)
	mockNodeRepo := new(mockNodeRepository)
	mockExerciseRepo := new(mockExerciseRepository)
	mockProgressRepo := new(mockUserProgressRepository)
	mockSubmissionRepo := new(mockSubmissionRepository)

	service := NewLearningService(
		mockTechTreeRepo,
		mockNodeRepo,
		mockExerciseRepo,
		mockProgressRepo,
		mockSubmissionRepo,
	)

	ctx := context.Background()
	treeID := "tree-123"
	userID := "user-123"
	requestID := "request-123"

	tree := createTestTechTree()
	node1 := createTestNode(tree.ID, 0)
	node2 := createTestNode(tree.ID, 1)
	exercises1 := []*models.Exercise{createTestExercise(node1.ID)}
	exercises2 := []*models.Exercise{createTestExercise(node2.ID), createTestExercise(node2.ID)}
	progress1 := createTestUserProgress(node1.ID, 1, 1, true)

	mockTechTreeRepo.On("GetByID", ctx, treeID).Return(tree, nil)
	mockNodeRepo.On("GetByTechTree", ctx, treeID).Return([]*models.Node{node1, node2}, nil)
	mockProgressRepo.On("GetByUser", ctx, userID).Return([]*models.UserProgress{progress1}, nil)
	mockExerciseRepo.On("GetByNode", ctx, node1.ID.String()).Return(exercises1, nil)
	mockExerciseRepo.On("GetByNode", ctx, node2.ID.String()).Return(exercises2, nil)

	result, err := service.GetTechTree(ctx, treeID, userID, requestID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, tree, result.TechTree)
	assert.Len(t, result.Nodes, 2)
	assert.Equal(t, 2, result.TotalNodes)
	assert.Equal(t, 1, result.CompletedNodes)

	// Check node 1 progress (has progress)
	assert.Equal(t, node1, result.Nodes[0].Node)
	assert.True(t, result.Nodes[0].Progress.MasteryAchieved)
	assert.Equal(t, 1, result.Nodes[0].Progress.ExercisesCompleted)
	assert.Equal(t, 1, result.Nodes[0].Progress.TotalExercises)
	assert.True(t, result.Nodes[0].Progress.IsUnlocked)

	// Check node 2 progress (no progress)
	assert.Equal(t, node2, result.Nodes[1].Node)
	assert.False(t, result.Nodes[1].Progress.MasteryAchieved)
	assert.Equal(t, 0, result.Nodes[1].Progress.ExercisesCompleted)
	assert.Equal(t, 2, result.Nodes[1].Progress.TotalExercises)
	assert.False(t, result.Nodes[1].Progress.IsUnlocked) // Not first node, no progress

	mockTechTreeRepo.AssertExpectations(t)
	mockNodeRepo.AssertExpectations(t)
	mockExerciseRepo.AssertExpectations(t)
	mockProgressRepo.AssertExpectations(t)
}

func TestGetTechTree_TreeNotFound(t *testing.T) {
	mockTechTreeRepo := new(mockTechTreeRepository)
	mockNodeRepo := new(mockNodeRepository)
	mockExerciseRepo := new(mockExerciseRepository)
	mockProgressRepo := new(mockUserProgressRepository)
	mockSubmissionRepo := new(mockSubmissionRepository)

	service := NewLearningService(
		mockTechTreeRepo,
		mockNodeRepo,
		mockExerciseRepo,
		mockProgressRepo,
		mockSubmissionRepo,
	)

	ctx := context.Background()
	treeID := "nonexistent"
	userID := "user-123"
	requestID := "request-123"

	mockTechTreeRepo.On("GetByID", ctx, treeID).Return(nil, repository.ErrNotFound)

	result, err := service.GetTechTree(ctx, treeID, userID, requestID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.IsType(t, &TechTreeNotFoundError{}, err)

	mockTechTreeRepo.AssertExpectations(t)
}

func TestGetTechTree_NoProgress(t *testing.T) {
	mockTechTreeRepo := new(mockTechTreeRepository)
	mockNodeRepo := new(mockNodeRepository)
	mockExerciseRepo := new(mockExerciseRepository)
	mockProgressRepo := new(mockUserProgressRepository)
	mockSubmissionRepo := new(mockSubmissionRepository)

	service := NewLearningService(
		mockTechTreeRepo,
		mockNodeRepo,
		mockExerciseRepo,
		mockProgressRepo,
		mockSubmissionRepo,
	)

	ctx := context.Background()
	treeID := "tree-123"
	userID := "user-123"
	requestID := "request-123"

	tree := createTestTechTree()
	node1 := createTestNode(tree.ID, 0)
	exercises := []*models.Exercise{createTestExercise(node1.ID)}

	mockTechTreeRepo.On("GetByID", ctx, treeID).Return(tree, nil)
	mockNodeRepo.On("GetByTechTree", ctx, treeID).Return([]*models.Node{node1}, nil)
	mockProgressRepo.On("GetByUser", ctx, userID).Return(nil, repository.ErrNotFound)
	mockExerciseRepo.On("GetByNode", ctx, node1.ID.String()).Return(exercises, nil)

	result, err := service.GetTechTree(ctx, treeID, userID, requestID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, result.CompletedNodes)
	assert.Len(t, result.Nodes, 1)
	assert.True(t, result.Nodes[0].Progress.IsUnlocked) // First node always unlocked
	assert.False(t, result.Nodes[0].Progress.MasteryAchieved)

	mockTechTreeRepo.AssertExpectations(t)
	mockNodeRepo.AssertExpectations(t)
	mockExerciseRepo.AssertExpectations(t)
	mockProgressRepo.AssertExpectations(t)
}

func TestGetTechTree_EmptyNodes(t *testing.T) {
	mockTechTreeRepo := new(mockTechTreeRepository)
	mockNodeRepo := new(mockNodeRepository)
	mockExerciseRepo := new(mockExerciseRepository)
	mockProgressRepo := new(mockUserProgressRepository)
	mockSubmissionRepo := new(mockSubmissionRepository)

	service := NewLearningService(
		mockTechTreeRepo,
		mockNodeRepo,
		mockExerciseRepo,
		mockProgressRepo,
		mockSubmissionRepo,
	)

	ctx := context.Background()
	treeID := "tree-123"
	userID := "user-123"
	requestID := "request-123"

	tree := createTestTechTree()

	mockTechTreeRepo.On("GetByID", ctx, treeID).Return(tree, nil)
	mockNodeRepo.On("GetByTechTree", ctx, treeID).Return([]*models.Node{}, nil)
	mockProgressRepo.On("GetByUser", ctx, userID).Return([]*models.UserProgress{}, nil)

	result, err := service.GetTechTree(ctx, treeID, userID, requestID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, result.TotalNodes)
	assert.Equal(t, 0, result.CompletedNodes)
	assert.Len(t, result.Nodes, 0)

	mockTechTreeRepo.AssertExpectations(t)
	mockNodeRepo.AssertExpectations(t)
	mockProgressRepo.AssertExpectations(t)
}

func TestGetTechTree_ExerciseCountError(t *testing.T) {
	mockTechTreeRepo := new(mockTechTreeRepository)
	mockNodeRepo := new(mockNodeRepository)
	mockExerciseRepo := new(mockExerciseRepository)
	mockProgressRepo := new(mockUserProgressRepository)
	mockSubmissionRepo := new(mockSubmissionRepository)

	service := NewLearningService(
		mockTechTreeRepo,
		mockNodeRepo,
		mockExerciseRepo,
		mockProgressRepo,
		mockSubmissionRepo,
	)

	ctx := context.Background()
	treeID := "tree-123"
	userID := "user-123"
	requestID := "request-123"

	tree := createTestTechTree()
	node1 := createTestNode(tree.ID, 0)

	mockTechTreeRepo.On("GetByID", ctx, treeID).Return(tree, nil)
	mockNodeRepo.On("GetByTechTree", ctx, treeID).Return([]*models.Node{node1}, nil)
	mockProgressRepo.On("GetByUser", ctx, userID).Return(nil, repository.ErrNotFound)
	mockExerciseRepo.On("GetByNode", ctx, node1.ID.String()).Return(nil, assert.AnError)

	result, err := service.GetTechTree(ctx, treeID, userID, requestID)

	// Should continue despite exercise count error
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Nodes, 1)
	assert.Equal(t, 0, result.Nodes[0].Progress.TotalExercises) // Defaults to 0 on error

	mockTechTreeRepo.AssertExpectations(t)
	mockNodeRepo.AssertExpectations(t)
	mockExerciseRepo.AssertExpectations(t)
	mockProgressRepo.AssertExpectations(t)
}

// ===== GetNode Tests =====

func TestGetNode_HappyPath(t *testing.T) {
	mockTechTreeRepo := new(mockTechTreeRepository)
	mockNodeRepo := new(mockNodeRepository)
	mockExerciseRepo := new(mockExerciseRepository)
	mockProgressRepo := new(mockUserProgressRepository)
	mockSubmissionRepo := new(mockSubmissionRepository)

	service := NewLearningService(
		mockTechTreeRepo,
		mockNodeRepo,
		mockExerciseRepo,
		mockProgressRepo,
		mockSubmissionRepo,
	)

	ctx := context.Background()
	userID := "user-123"
	requestID := "request-123"

	node := createTestNode(uuid.New(), 0)
	nodeID := node.ID.String()
	exercises := []*models.Exercise{createTestExercise(node.ID), createTestExercise(node.ID)}
	progress := createTestUserProgress(node.ID, 2, 2, true)

	// Create submissions with scores
	score1 := 85
	score2 := 90
	submission1 := &models.Submission{
		ID:         uuid.New(),
		ExerciseID: exercises[0].ID,
		TotalScore: &score1,
	}
	submission2 := &models.Submission{
		ID:         uuid.New(),
		ExerciseID: exercises[1].ID,
		TotalScore: &score2,
	}

	mockNodeRepo.On("GetByID", ctx, nodeID).Return(node, nil)
	mockExerciseRepo.On("GetByNode", ctx, nodeID).Return(exercises, nil)
	mockProgressRepo.On("GetByUserAndNode", ctx, userID, nodeID).Return(progress, nil)
	mockSubmissionRepo.On("GetByExercisePaginated", ctx, userID, exercises[0].ID.String(), 100, 0).
		Return([]*models.Submission{submission1}, 1, nil)
	mockSubmissionRepo.On("GetByExercisePaginated", ctx, userID, exercises[1].ID.String(), 100, 0).
		Return([]*models.Submission{submission2}, 1, nil)

	result, err := service.GetNode(ctx, nodeID, userID, requestID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, node, result.Node)
	assert.Len(t, result.Exercises, 2)
	assert.NotNil(t, result.Progress)
	assert.True(t, result.Progress.MasteryAchieved)
	assert.Equal(t, 2, result.Progress.ExercisesCompleted)
	assert.Equal(t, 2, result.Progress.TotalExercises)
	assert.True(t, result.Progress.IsUnlocked)

	// Check best scores
	assert.Equal(t, 85, result.Progress.BestScores[exercises[0].ID.String()])
	assert.Equal(t, 90, result.Progress.BestScores[exercises[1].ID.String()])

	mockNodeRepo.AssertExpectations(t)
	mockExerciseRepo.AssertExpectations(t)
	mockProgressRepo.AssertExpectations(t)
	mockSubmissionRepo.AssertExpectations(t)
}

func TestGetNode_NodeNotFound(t *testing.T) {
	mockTechTreeRepo := new(mockTechTreeRepository)
	mockNodeRepo := new(mockNodeRepository)
	mockExerciseRepo := new(mockExerciseRepository)
	mockProgressRepo := new(mockUserProgressRepository)
	mockSubmissionRepo := new(mockSubmissionRepository)

	service := NewLearningService(
		mockTechTreeRepo,
		mockNodeRepo,
		mockExerciseRepo,
		mockProgressRepo,
		mockSubmissionRepo,
	)

	ctx := context.Background()
	nodeID := "nonexistent"
	userID := "user-123"
	requestID := "request-123"

	mockNodeRepo.On("GetByID", ctx, nodeID).Return(nil, repository.ErrNotFound)

	result, err := service.GetNode(ctx, nodeID, userID, requestID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.IsType(t, &NodeNotFoundError{}, err)

	mockNodeRepo.AssertExpectations(t)
}

func TestGetNode_NoProgress(t *testing.T) {
	mockTechTreeRepo := new(mockTechTreeRepository)
	mockNodeRepo := new(mockNodeRepository)
	mockExerciseRepo := new(mockExerciseRepository)
	mockProgressRepo := new(mockUserProgressRepository)
	mockSubmissionRepo := new(mockSubmissionRepository)

	service := NewLearningService(
		mockTechTreeRepo,
		mockNodeRepo,
		mockExerciseRepo,
		mockProgressRepo,
		mockSubmissionRepo,
	)

	ctx := context.Background()
	userID := "user-123"
	requestID := "request-123"

	node := createTestNode(uuid.New(), 0)
	nodeID := node.ID.String()
	exercises := []*models.Exercise{createTestExercise(node.ID)}

	mockNodeRepo.On("GetByID", ctx, nodeID).Return(node, nil)
	mockExerciseRepo.On("GetByNode", ctx, nodeID).Return(exercises, nil)
	mockProgressRepo.On("GetByUserAndNode", ctx, userID, nodeID).Return(nil, repository.ErrNotFound)
	mockSubmissionRepo.On("GetByExercisePaginated", ctx, userID, exercises[0].ID.String(), 100, 0).
		Return([]*models.Submission{}, 0, nil)

	result, err := service.GetNode(ctx, nodeID, userID, requestID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, node, result.Node)
	assert.NotNil(t, result.Progress)
	assert.False(t, result.Progress.MasteryAchieved)
	assert.Equal(t, 0, result.Progress.ExercisesCompleted)
	assert.Equal(t, 1, result.Progress.TotalExercises)
	assert.True(t, result.Progress.IsUnlocked) // Assume unlocked if viewing

	mockNodeRepo.AssertExpectations(t)
	mockExerciseRepo.AssertExpectations(t)
	mockProgressRepo.AssertExpectations(t)
	mockSubmissionRepo.AssertExpectations(t)
}

func TestGetNode_BestScoreCalculation(t *testing.T) {
	mockTechTreeRepo := new(mockTechTreeRepository)
	mockNodeRepo := new(mockNodeRepository)
	mockExerciseRepo := new(mockExerciseRepository)
	mockProgressRepo := new(mockUserProgressRepository)
	mockSubmissionRepo := new(mockSubmissionRepository)

	service := NewLearningService(
		mockTechTreeRepo,
		mockNodeRepo,
		mockExerciseRepo,
		mockProgressRepo,
		mockSubmissionRepo,
	)

	ctx := context.Background()
	userID := "user-123"
	requestID := "request-123"

	node := createTestNode(uuid.New(), 0)
	nodeID := node.ID.String()
	exercise := createTestExercise(node.ID)

	// Multiple submissions with different scores - should pick the best
	score1 := 70
	score2 := 95
	score3 := 80
	submissions := []*models.Submission{
		{ID: uuid.New(), ExerciseID: exercise.ID, TotalScore: &score1},
		{ID: uuid.New(), ExerciseID: exercise.ID, TotalScore: &score2},
		{ID: uuid.New(), ExerciseID: exercise.ID, TotalScore: &score3},
	}

	mockNodeRepo.On("GetByID", ctx, nodeID).Return(node, nil)
	mockExerciseRepo.On("GetByNode", ctx, nodeID).Return([]*models.Exercise{exercise}, nil)
	mockProgressRepo.On("GetByUserAndNode", ctx, userID, nodeID).Return(nil, repository.ErrNotFound)
	mockSubmissionRepo.On("GetByExercisePaginated", ctx, userID, exercise.ID.String(), 100, 0).
		Return(submissions, 3, nil)

	result, err := service.GetNode(ctx, nodeID, userID, requestID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 95, result.Progress.BestScores[exercise.ID.String()]) // Best score

	mockNodeRepo.AssertExpectations(t)
	mockExerciseRepo.AssertExpectations(t)
	mockProgressRepo.AssertExpectations(t)
	mockSubmissionRepo.AssertExpectations(t)
}

func TestGetNode_NoExercises(t *testing.T) {
	mockTechTreeRepo := new(mockTechTreeRepository)
	mockNodeRepo := new(mockNodeRepository)
	mockExerciseRepo := new(mockExerciseRepository)
	mockProgressRepo := new(mockUserProgressRepository)
	mockSubmissionRepo := new(mockSubmissionRepository)

	service := NewLearningService(
		mockTechTreeRepo,
		mockNodeRepo,
		mockExerciseRepo,
		mockProgressRepo,
		mockSubmissionRepo,
	)

	ctx := context.Background()
	userID := "user-123"
	requestID := "request-123"

	node := createTestNode(uuid.New(), 0)
	nodeID := node.ID.String()

	mockNodeRepo.On("GetByID", ctx, nodeID).Return(node, nil)
	mockExerciseRepo.On("GetByNode", ctx, nodeID).Return([]*models.Exercise{}, nil)
	mockProgressRepo.On("GetByUserAndNode", ctx, userID, nodeID).Return(nil, repository.ErrNotFound)

	result, err := service.GetNode(ctx, nodeID, userID, requestID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Exercises, 0)
	assert.Equal(t, 0, result.Progress.TotalExercises)

	mockNodeRepo.AssertExpectations(t)
	mockExerciseRepo.AssertExpectations(t)
	mockProgressRepo.AssertExpectations(t)
}

// ===== CalculateProgress Tests =====

func TestCalculateProgress_HappyPath(t *testing.T) {
	mockTechTreeRepo := new(mockTechTreeRepository)
	mockNodeRepo := new(mockNodeRepository)
	mockExerciseRepo := new(mockExerciseRepository)
	mockProgressRepo := new(mockUserProgressRepository)
	mockSubmissionRepo := new(mockSubmissionRepository)

	service := NewLearningService(
		mockTechTreeRepo,
		mockNodeRepo,
		mockExerciseRepo,
		mockProgressRepo,
		mockSubmissionRepo,
	)

	ctx := context.Background()
	userID := "user-123"
	treeID := "tree-123"

	treeUUID := uuid.New()
	node1 := createTestNode(treeUUID, 0)
	node2 := createTestNode(treeUUID, 1)
	node3 := createTestNode(treeUUID, 2)

	progress1 := createTestUserProgress(node1.ID, 3, 3, true)
	progress2 := createTestUserProgress(node2.ID, 2, 3, false)

	mockNodeRepo.On("GetByTechTree", ctx, treeID).Return([]*models.Node{node1, node2, node3}, nil)
	mockProgressRepo.On("GetByUser", ctx, userID).Return([]*models.UserProgress{progress1, progress2}, nil)

	completedNodes, totalNodes, err := service.CalculateProgress(ctx, userID, treeID)

	assert.NoError(t, err)
	assert.Equal(t, 1, completedNodes) // Only node1 is mastered
	assert.Equal(t, 3, totalNodes)

	mockNodeRepo.AssertExpectations(t)
	mockProgressRepo.AssertExpectations(t)
}

func TestCalculateProgress_NoProgress(t *testing.T) {
	mockTechTreeRepo := new(mockTechTreeRepository)
	mockNodeRepo := new(mockNodeRepository)
	mockExerciseRepo := new(mockExerciseRepository)
	mockProgressRepo := new(mockUserProgressRepository)
	mockSubmissionRepo := new(mockSubmissionRepository)

	service := NewLearningService(
		mockTechTreeRepo,
		mockNodeRepo,
		mockExerciseRepo,
		mockProgressRepo,
		mockSubmissionRepo,
	)

	ctx := context.Background()
	userID := "user-123"
	treeID := "tree-123"

	treeUUID := uuid.New()
	node1 := createTestNode(treeUUID, 0)
	node2 := createTestNode(treeUUID, 1)

	mockNodeRepo.On("GetByTechTree", ctx, treeID).Return([]*models.Node{node1, node2}, nil)
	mockProgressRepo.On("GetByUser", ctx, userID).Return(nil, repository.ErrNotFound)

	completedNodes, totalNodes, err := service.CalculateProgress(ctx, userID, treeID)

	assert.NoError(t, err)
	assert.Equal(t, 0, completedNodes)
	assert.Equal(t, 2, totalNodes)

	mockNodeRepo.AssertExpectations(t)
	mockProgressRepo.AssertExpectations(t)
}

func TestCalculateProgress_EmptyTree(t *testing.T) {
	mockTechTreeRepo := new(mockTechTreeRepository)
	mockNodeRepo := new(mockNodeRepository)
	mockExerciseRepo := new(mockExerciseRepository)
	mockProgressRepo := new(mockUserProgressRepository)
	mockSubmissionRepo := new(mockSubmissionRepository)

	service := NewLearningService(
		mockTechTreeRepo,
		mockNodeRepo,
		mockExerciseRepo,
		mockProgressRepo,
		mockSubmissionRepo,
	)

	ctx := context.Background()
	userID := "user-123"
	treeID := "tree-123"

	mockNodeRepo.On("GetByTechTree", ctx, treeID).Return([]*models.Node{}, nil)
	mockProgressRepo.On("GetByUser", ctx, userID).Return([]*models.UserProgress{}, nil)

	completedNodes, totalNodes, err := service.CalculateProgress(ctx, userID, treeID)

	assert.NoError(t, err)
	assert.Equal(t, 0, completedNodes)
	assert.Equal(t, 0, totalNodes)

	mockNodeRepo.AssertExpectations(t)
}

func TestCalculateProgress_NodeRepositoryError(t *testing.T) {
	mockTechTreeRepo := new(mockTechTreeRepository)
	mockNodeRepo := new(mockNodeRepository)
	mockExerciseRepo := new(mockExerciseRepository)
	mockProgressRepo := new(mockUserProgressRepository)
	mockSubmissionRepo := new(mockSubmissionRepository)

	service := NewLearningService(
		mockTechTreeRepo,
		mockNodeRepo,
		mockExerciseRepo,
		mockProgressRepo,
		mockSubmissionRepo,
	)

	ctx := context.Background()
	userID := "user-123"
	treeID := "tree-123"

	mockNodeRepo.On("GetByTechTree", ctx, treeID).Return(nil, assert.AnError)

	completedNodes, totalNodes, err := service.CalculateProgress(ctx, userID, treeID)

	assert.Error(t, err)
	assert.Equal(t, 0, completedNodes)
	assert.Equal(t, 0, totalNodes)

	mockNodeRepo.AssertExpectations(t)
}
