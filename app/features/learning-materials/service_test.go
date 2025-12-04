package learningmaterials

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/anttihil/code-auto-grader/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ===== Mock Implementations =====

type mockLearningMaterialRepository struct {
	mock.Mock
}

func (m *mockLearningMaterialRepository) GetByNodeID(ctx context.Context, nodeID string) ([]*models.LearningMaterial, error) {
	args := m.Called(ctx, nodeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.LearningMaterial), args.Error(1)
}

func (m *mockLearningMaterialRepository) GetByID(ctx context.Context, id string) (*models.LearningMaterial, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LearningMaterial), args.Error(1)
}

func (m *mockLearningMaterialRepository) Create(ctx context.Context, material *models.LearningMaterial) error {
	args := m.Called(ctx, material)
	return args.Error(0)
}

func (m *mockLearningMaterialRepository) Update(ctx context.Context, material *models.LearningMaterial) error {
	args := m.Called(ctx, material)
	return args.Error(0)
}

func (m *mockLearningMaterialRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ===== Test Fixtures =====

func createTestMaterial(nodeID uuid.UUID, title string, orderIndex int) *models.LearningMaterial {
	return &models.LearningMaterial{
		ID:         uuid.New(),
		NodeID:     nodeID,
		Title:      title,
		Content:    "Test content for " + title,
		OrderIndex: orderIndex,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

// ===== Constructor Tests =====

func TestNewService(t *testing.T) {
	t.Run("creates service with valid repository", func(t *testing.T) {
		mockRepo := new(mockLearningMaterialRepository)
		service := NewService(mockRepo)

		assert.NotNil(t, service)
		assert.Equal(t, mockRepo, service.materialRepo)
	})

	t.Run("panics with nil repository", func(t *testing.T) {
		assert.Panics(t, func() {
			NewService(nil)
		})
	})
}

// ===== GetByNodeID Tests =====

func TestGetByNodeID(t *testing.T) {
	ctx := context.Background()
	nodeID := uuid.New()

	t.Run("returns materials for node", func(t *testing.T) {
		mockRepo := new(mockLearningMaterialRepository)
		service := NewService(mockRepo)

		materials := []*models.LearningMaterial{
			createTestMaterial(nodeID, "Introduction", 0),
			createTestMaterial(nodeID, "Key Concepts", 1),
			createTestMaterial(nodeID, "Examples", 2),
		}

		mockRepo.On("GetByNodeID", ctx, nodeID.String()).Return(materials, nil)

		result, err := service.GetByNodeID(ctx, nodeID.String())

		assert.NoError(t, err)
		assert.Len(t, result, 3)
		assert.Equal(t, "Introduction", result[0].Title)
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns empty slice for node with no materials", func(t *testing.T) {
		mockRepo := new(mockLearningMaterialRepository)
		service := NewService(mockRepo)

		mockRepo.On("GetByNodeID", ctx, nodeID.String()).Return([]*models.LearningMaterial{}, nil)

		result, err := service.GetByNodeID(ctx, nodeID.String())

		assert.NoError(t, err)
		assert.Empty(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns error from repository", func(t *testing.T) {
		mockRepo := new(mockLearningMaterialRepository)
		service := NewService(mockRepo)

		mockRepo.On("GetByNodeID", ctx, nodeID.String()).Return(nil, errors.New("database error"))

		result, err := service.GetByNodeID(ctx, nodeID.String())

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

// ===== GetByID Tests =====

func TestService_GetByID(t *testing.T) {
	ctx := context.Background()
	materialID := uuid.New()
	nodeID := uuid.New()

	t.Run("returns material by ID", func(t *testing.T) {
		mockRepo := new(mockLearningMaterialRepository)
		service := NewService(mockRepo)

		material := createTestMaterial(nodeID, "Test Material", 0)
		material.ID = materialID

		mockRepo.On("GetByID", ctx, materialID.String()).Return(material, nil)

		result, err := service.GetByID(ctx, materialID.String())

		assert.NoError(t, err)
		assert.Equal(t, materialID, result.ID)
		assert.Equal(t, "Test Material", result.Title)
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns error when not found", func(t *testing.T) {
		mockRepo := new(mockLearningMaterialRepository)
		service := NewService(mockRepo)

		mockRepo.On("GetByID", ctx, materialID.String()).Return(nil, errors.New("not found"))

		result, err := service.GetByID(ctx, materialID.String())

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

// ===== Create Tests =====

func TestService_Create(t *testing.T) {
	ctx := context.Background()
	nodeID := uuid.New()

	t.Run("creates material successfully", func(t *testing.T) {
		mockRepo := new(mockLearningMaterialRepository)
		service := NewService(mockRepo)

		req := CreateLearningMaterialRequest{
			Title:      "New Material",
			Content:    "Content for the material",
			OrderIndex: 1,
		}

		mockRepo.On("Create", ctx, mock.AnythingOfType("*models.LearningMaterial")).Return(nil)

		result, err := service.Create(ctx, nodeID.String(), req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "New Material", result.Title)
		assert.Equal(t, "Content for the material", result.Content)
		assert.Equal(t, 1, result.OrderIndex)
		assert.Equal(t, nodeID, result.NodeID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns error for invalid node ID", func(t *testing.T) {
		mockRepo := new(mockLearningMaterialRepository)
		service := NewService(mockRepo)

		req := CreateLearningMaterialRequest{
			Title:   "New Material",
			Content: "Content",
		}

		result, err := service.Create(ctx, "invalid-uuid", req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid node ID")
	})

	t.Run("returns error from repository", func(t *testing.T) {
		mockRepo := new(mockLearningMaterialRepository)
		service := NewService(mockRepo)

		req := CreateLearningMaterialRequest{
			Title:   "New Material",
			Content: "Content",
		}

		mockRepo.On("Create", ctx, mock.AnythingOfType("*models.LearningMaterial")).
			Return(errors.New("database error"))

		result, err := service.Create(ctx, nodeID.String(), req)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

// ===== Update Tests =====

func TestService_Update(t *testing.T) {
	ctx := context.Background()
	materialID := uuid.New()
	nodeID := uuid.New()

	t.Run("updates material successfully", func(t *testing.T) {
		mockRepo := new(mockLearningMaterialRepository)
		service := NewService(mockRepo)

		existingMaterial := createTestMaterial(nodeID, "Old Title", 0)
		existingMaterial.ID = materialID

		newTitle := "New Title"
		newContent := "New Content"
		newOrder := 5
		req := UpdateLearningMaterialRequest{
			Title:      &newTitle,
			Content:    &newContent,
			OrderIndex: &newOrder,
		}

		mockRepo.On("GetByID", ctx, materialID.String()).Return(existingMaterial, nil)
		mockRepo.On("Update", ctx, existingMaterial).Return(nil)

		result, err := service.Update(ctx, materialID.String(), req)

		assert.NoError(t, err)
		assert.Equal(t, "New Title", result.Title)
		assert.Equal(t, "New Content", result.Content)
		assert.Equal(t, 5, result.OrderIndex)
		mockRepo.AssertExpectations(t)
	})

	t.Run("updates only specified fields", func(t *testing.T) {
		mockRepo := new(mockLearningMaterialRepository)
		service := NewService(mockRepo)

		existingMaterial := createTestMaterial(nodeID, "Original Title", 3)
		existingMaterial.ID = materialID
		existingMaterial.Content = "Original Content"

		newTitle := "Updated Title Only"
		req := UpdateLearningMaterialRequest{
			Title: &newTitle,
			// Content and OrderIndex are nil - should not change
		}

		mockRepo.On("GetByID", ctx, materialID.String()).Return(existingMaterial, nil)
		mockRepo.On("Update", ctx, existingMaterial).Return(nil)

		result, err := service.Update(ctx, materialID.String(), req)

		assert.NoError(t, err)
		assert.Equal(t, "Updated Title Only", result.Title)
		assert.Equal(t, "Original Content", result.Content)
		assert.Equal(t, 3, result.OrderIndex)
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns error when material not found", func(t *testing.T) {
		mockRepo := new(mockLearningMaterialRepository)
		service := NewService(mockRepo)

		req := UpdateLearningMaterialRequest{}

		mockRepo.On("GetByID", ctx, materialID.String()).Return(nil, errors.New("not found"))

		result, err := service.Update(ctx, materialID.String(), req)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

// ===== Delete Tests =====

func TestService_Delete(t *testing.T) {
	ctx := context.Background()
	materialID := uuid.New()

	t.Run("deletes material successfully", func(t *testing.T) {
		mockRepo := new(mockLearningMaterialRepository)
		service := NewService(mockRepo)

		mockRepo.On("Delete", ctx, materialID.String()).Return(nil)

		err := service.Delete(ctx, materialID.String())

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns error from repository", func(t *testing.T) {
		mockRepo := new(mockLearningMaterialRepository)
		service := NewService(mockRepo)

		mockRepo.On("Delete", ctx, materialID.String()).Return(errors.New("not found"))

		err := service.Delete(ctx, materialID.String())

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
