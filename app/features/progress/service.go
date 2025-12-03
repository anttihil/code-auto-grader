package progress

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/anttihil/code-auto-grader/internal/logger"
	"github.com/anttihil/code-auto-grader/models"
	"github.com/google/uuid"
)

// Service handles user progress tracking and mastery calculations
type Service struct {
	repo Repository
}

// NewService creates a new progress service
func NewService(repo Repository) *Service {
	if repo == nil {
		panic("progress repository is required")
	}
	return &Service{
		repo: repo,
	}
}

// UpdateAfterSubmission updates user progress after a code submission
// This should be called when a user submits code for an exercise
func (s *Service) UpdateAfterSubmission(ctx context.Context, userID, nodeID string, score int, requestID string) error {
	// Get current progress
	progress, err := s.repo.GetByUserAndNode(ctx, userID, nodeID)
	if err != nil {
		return fmt.Errorf("failed to get user progress: %w", err)
	}

	// If no progress exists, create initial progress record
	if progress == nil {
		now := time.Now()
		progress = &models.UserProgress{
			ID:                 uuid.New(),
			UserID:             uuid.MustParse(userID),
			NodeID:             uuid.MustParse(nodeID),
			Status:             StatusInProgress,
			ExercisesCompleted: 0,
			ExercisesTotal:     0, // TODO: Count exercises for this node
			MasteryThreshold:   DefaultMasteryThreshold,
			IsMastered:         false,
			StartedAt:          &now,
			LastActivityAt:     &now,
		}
	}

	// Update last activity
	now := time.Now()
	progress.LastActivityAt = &now

	// If score >= mastery threshold, increment exercises completed
	if float64(score) >= progress.MasteryThreshold {
		progress.ExercisesCompleted++
	}

	// Calculate mastery score
	masteryScore := float64(score)
	progress.MasteryScore = &masteryScore

	// Check if mastered
	if masteryScore >= progress.MasteryThreshold {
		progress.IsMastered = true
		progress.Status = StatusCompleted
		if progress.CompletedAt == nil {
			progress.CompletedAt = &now
		}
	}

	// Save updated progress
	if err := s.repo.Upsert(ctx, progress); err != nil {
		return fmt.Errorf("failed to update progress: %w", err)
	}

	logger.Logger.Info("user progress updated",
		slog.String("request_id", requestID),
		slog.String("user_id", userID),
		slog.String("node_id", nodeID),
		slog.Int("exercises_completed", progress.ExercisesCompleted),
		slog.Bool("is_mastered", progress.IsMastered),
	)

	return nil
}

// GetByUserAndNode retrieves progress for a specific user and node
func (s *Service) GetByUserAndNode(ctx context.Context, userID, nodeID string) (*models.UserProgress, error) {
	return s.repo.GetByUserAndNode(ctx, userID, nodeID)
}

// GetByUser retrieves all progress records for a user
func (s *Service) GetByUser(ctx context.Context, userID string) ([]*models.UserProgress, error) {
	return s.repo.GetByUser(ctx, userID)
}

// CalculateCompletionStats calculates how many nodes a user has completed
// Returns (completedCount, totalCount, error)
func CalculateCompletionStats(progressList []*models.UserProgress, totalNodes int) (int, int) {
	completed := 0
	for _, p := range progressList {
		if p.IsMastered {
			completed++
		}
	}
	return completed, totalNodes
}
