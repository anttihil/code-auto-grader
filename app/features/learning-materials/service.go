package learningmaterials

import (
	"context"
	"fmt"

	"github.com/anttihil/code-auto-grader/models"
	"github.com/anttihil/code-auto-grader/repository"
	"github.com/google/uuid"
)

// Service handles business logic for learning materials
type Service struct {
	materialRepo repository.LearningMaterialRepository
}

// NewService creates a new learning materials service
func NewService(materialRepo repository.LearningMaterialRepository) *Service {
	if materialRepo == nil {
		panic("material repository is required")
	}
	return &Service{materialRepo: materialRepo}
}

// GetByNodeID retrieves all learning materials for a node
func (s *Service) GetByNodeID(ctx context.Context, nodeID string) ([]*models.LearningMaterial, error) {
	return s.materialRepo.GetByNodeID(ctx, nodeID)
}

// GetByID retrieves a learning material by ID
func (s *Service) GetByID(ctx context.Context, id string) (*models.LearningMaterial, error) {
	return s.materialRepo.GetByID(ctx, id)
}

// Create creates a new learning material
func (s *Service) Create(ctx context.Context, nodeID string, req CreateLearningMaterialRequest) (*models.LearningMaterial, error) {
	nodeUUID, err := uuid.Parse(nodeID)
	if err != nil {
		return nil, fmt.Errorf("invalid node ID: %w", err)
	}

	material := &models.LearningMaterial{
		ID:         uuid.New(),
		NodeID:     nodeUUID,
		Title:      req.Title,
		Content:    req.Content,
		OrderIndex: req.OrderIndex,
	}

	if err := s.materialRepo.Create(ctx, material); err != nil {
		return nil, fmt.Errorf("failed to create learning material: %w", err)
	}

	return material, nil
}

// Update updates an existing learning material
func (s *Service) Update(ctx context.Context, id string, req UpdateLearningMaterialRequest) (*models.LearningMaterial, error) {
	material, err := s.materialRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("learning material not found: %w", err)
	}

	// Apply updates
	if req.Title != nil {
		material.Title = *req.Title
	}
	if req.Content != nil {
		material.Content = *req.Content
	}
	if req.OrderIndex != nil {
		material.OrderIndex = *req.OrderIndex
	}

	if err := s.materialRepo.Update(ctx, material); err != nil {
		return nil, fmt.Errorf("failed to update learning material: %w", err)
	}

	return material, nil
}

// Delete deletes a learning material
func (s *Service) Delete(ctx context.Context, id string) error {
	return s.materialRepo.Delete(ctx, id)
}
