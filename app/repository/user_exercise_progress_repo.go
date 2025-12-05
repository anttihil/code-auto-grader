package repository

import (
	"context"
	"fmt"

	"github.com/anttihil/code-auto-grader/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserExerciseProgressRepo implements UserExerciseProgressRepository using pgx
type UserExerciseProgressRepo struct {
	pool *pgxpool.Pool
}

// NewUserExerciseProgressRepo creates a new user exercise progress repository
func NewUserExerciseProgressRepo(pool *pgxpool.Pool) *UserExerciseProgressRepo {
	return &UserExerciseProgressRepo{pool: pool}
}

// GetByUserAndExercise retrieves progress for a specific user and exercise
func (r *UserExerciseProgressRepo) GetByUserAndExercise(ctx context.Context, userID, exerciseID string) (*models.UserExerciseProgress, error) {
	query := `
		SELECT id, user_id, exercise_id, ai_help_count, attempts_count, created_at, updated_at
		FROM user_exercise_progress
		WHERE user_id = $1 AND exercise_id = $2
	`

	var progress models.UserExerciseProgress
	err := r.pool.QueryRow(ctx, query, userID, exerciseID).Scan(
		&progress.ID,
		&progress.UserID,
		&progress.ExerciseID,
		&progress.AIHelpCount,
		&progress.AttemptsCount,
		&progress.CreatedAt,
		&progress.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // No progress yet, not an error
		}
		return nil, fmt.Errorf("failed to get user exercise progress: %w", err)
	}

	return &progress, nil
}

// GetByUserAndNode retrieves all exercise progress for a user within a specific node
func (r *UserExerciseProgressRepo) GetByUserAndNode(ctx context.Context, userID, nodeID string) ([]*models.UserExerciseProgress, error) {
	// TODO(human): Write the SQL query that joins user_exercise_progress with exercises
	// to filter by node_id. The exercises table has the node_id column.
	query := `
		SELECT uep.id, uep.user_id, uep.exercise_id, uep.ai_help_count, uep.attempts_count,
		       uep.created_at, uep.updated_at
		FROM user_exercise_progress uep
		INNER JOIN exercises e ON uep.exercise_id = e.id
		WHERE uep.user_id = $1 AND e.node_id = $2
		ORDER BY e.order_index
	`

	rows, err := r.pool.Query(ctx, query, userID, nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user exercise progress: %w", err)
	}
	defer rows.Close()

	var progressList []*models.UserExerciseProgress
	for rows.Next() {
		var progress models.UserExerciseProgress
		err := rows.Scan(
			&progress.ID,
			&progress.UserID,
			&progress.ExerciseID,
			&progress.AIHelpCount,
			&progress.AttemptsCount,
			&progress.CreatedAt,
			&progress.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user exercise progress: %w", err)
		}
		progressList = append(progressList, &progress)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user exercise progress: %w", err)
	}

	return progressList, nil
}

// UpsertProgress creates or updates a user exercise progress record
func (r *UserExerciseProgressRepo) UpsertProgress(ctx context.Context, progress *models.UserExerciseProgress) error {
	query := `
		INSERT INTO user_exercise_progress (id, user_id, exercise_id, ai_help_count, attempts_count)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, exercise_id)
		DO UPDATE SET
			ai_help_count = EXCLUDED.ai_help_count,
			attempts_count = EXCLUDED.attempts_count,
			updated_at = current_timestamp
		RETURNING created_at, updated_at
	`

	err := r.pool.QueryRow(ctx, query,
		progress.ID,
		progress.UserID,
		progress.ExerciseID,
		progress.AIHelpCount,
		progress.AttemptsCount,
	).Scan(&progress.CreatedAt, &progress.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to upsert user exercise progress: %w", err)
	}

	return nil
}

// UpdateHelpCounts atomically increments AI help and attempts counters
func (r *UserExerciseProgressRepo) UpdateHelpCounts(ctx context.Context, userID, exerciseID string, addedAICount, addedAttempts int) error {
	query := `
		INSERT INTO user_exercise_progress AS u (id, user_id, exercise_id, ai_help_count, attempts_count)
		VALUES (gen_random_uuid(), $1, $2, $3, $4)
		ON CONFLICT (user_id, exercise_id)
		DO UPDATE SET
			ai_help_count = u.ai_help_count + EXCLUDED.ai_help_count,
			attempts_count = u.attempts_count + EXCLUDED.attempts_count,
			updated_at = current_timestamp
	`

	_, err := r.pool.Exec(ctx, query, userID, exerciseID, addedAICount, addedAttempts)
	if err != nil {
		return fmt.Errorf("failed to update user exercise progress: %w", err)
	}

	return nil
}
