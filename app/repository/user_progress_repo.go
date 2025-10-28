package repository

import (
	"context"
	"fmt"

	"github.com/anttihil/code-auto-grader/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserProgressRepo implements UserProgressRepository using pgx
type UserProgressRepo struct {
	pool *pgxpool.Pool
}

// NewUserProgressRepo creates a new user progress repository
func NewUserProgressRepo(pool *pgxpool.Pool) *UserProgressRepo {
	return &UserProgressRepo{pool: pool}
}

// GetByUserAndNode retrieves progress for a specific user and node
func (r *UserProgressRepo) GetByUserAndNode(ctx context.Context, userID, nodeID string) (*models.UserProgress, error) {
	query := `
		SELECT id, user_id, node_id, status, exercises_completed, exercises_total,
		       mastery_score, mastery_threshold, is_mastered, started_at,
		       completed_at, last_activity_at, created_at, updated_at
		FROM user_progress
		WHERE user_id = $1 AND node_id = $2
	`

	var progress models.UserProgress
	err := r.pool.QueryRow(ctx, query, userID, nodeID).Scan(
		&progress.ID,
		&progress.UserID,
		&progress.NodeID,
		&progress.Status,
		&progress.ExercisesCompleted,
		&progress.ExercisesTotal,
		&progress.MasteryScore,
		&progress.MasteryThreshold,
		&progress.IsMastered,
		&progress.StartedAt,
		&progress.CompletedAt,
		&progress.LastActivityAt,
		&progress.CreatedAt,
		&progress.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // No progress yet, not an error
		}
		return nil, fmt.Errorf("failed to get user progress: %w", err)
	}

	return &progress, nil
}

// GetByUser retrieves all progress records for a user
func (r *UserProgressRepo) GetByUser(ctx context.Context, userID string) ([]*models.UserProgress, error) {
	query := `
		SELECT id, user_id, node_id, status, exercises_completed, exercises_total,
		       mastery_score, mastery_threshold, is_mastered, started_at,
		       completed_at, last_activity_at, created_at, updated_at
		FROM user_progress
		WHERE user_id = $1
		ORDER BY last_activity_at DESC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user progress: %w", err)
	}
	defer rows.Close()

	var progressList []*models.UserProgress
	for rows.Next() {
		var progress models.UserProgress
		err := rows.Scan(
			&progress.ID,
			&progress.UserID,
			&progress.NodeID,
			&progress.Status,
			&progress.ExercisesCompleted,
			&progress.ExercisesTotal,
			&progress.MasteryScore,
			&progress.MasteryThreshold,
			&progress.IsMastered,
			&progress.StartedAt,
			&progress.CompletedAt,
			&progress.LastActivityAt,
			&progress.CreatedAt,
			&progress.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user progress: %w", err)
		}
		progressList = append(progressList, &progress)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user progress: %w", err)
	}

	return progressList, nil
}

// Upsert creates or updates a user progress record
func (r *UserProgressRepo) Upsert(ctx context.Context, progress *models.UserProgress) error {
	query := `
		INSERT INTO user_progress (
			id, user_id, node_id, status, exercises_completed, exercises_total,
			mastery_score, mastery_threshold, is_mastered, started_at,
			completed_at, last_activity_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (user_id, node_id)
		DO UPDATE SET
			status = EXCLUDED.status,
			exercises_completed = EXCLUDED.exercises_completed,
			exercises_total = EXCLUDED.exercises_total,
			mastery_score = EXCLUDED.mastery_score,
			mastery_threshold = EXCLUDED.mastery_threshold,
			is_mastered = EXCLUDED.is_mastered,
			completed_at = EXCLUDED.completed_at,
			last_activity_at = EXCLUDED.last_activity_at
		RETURNING created_at, updated_at
	`

	err := r.pool.QueryRow(ctx, query,
		progress.ID,
		progress.UserID,
		progress.NodeID,
		progress.Status,
		progress.ExercisesCompleted,
		progress.ExercisesTotal,
		progress.MasteryScore,
		progress.MasteryThreshold,
		progress.IsMastered,
		progress.StartedAt,
		progress.CompletedAt,
		progress.LastActivityAt,
	).Scan(&progress.CreatedAt, &progress.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to upsert user progress: %w", err)
	}

	return nil
}

// UpdateMastery updates mastery-related fields for a user's progress
func (r *UserProgressRepo) UpdateMastery(ctx context.Context, userID, nodeID string, masteryScore float64, isMastered bool) error {
	query := `
		UPDATE user_progress
		SET mastery_score = $3,
		    is_mastered = $4,
		    last_activity_at = current_timestamp
		WHERE user_id = $1 AND node_id = $2
		RETURNING updated_at
	`

	var updatedAt string
	err := r.pool.QueryRow(ctx, query, userID, nodeID, masteryScore, isMastered).Scan(&updatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("user progress not found: user=%s node=%s", userID, nodeID)
		}
		return fmt.Errorf("failed to update mastery: %w", err)
	}

	return nil
}