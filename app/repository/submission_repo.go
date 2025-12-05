package repository

import (
	"context"
	"fmt"

	"github.com/anttihil/code-auto-grader/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SubmissionRepo implements SubmissionRepository using pgx
type SubmissionRepo struct {
	pool *pgxpool.Pool
}

// NewSubmissionRepo creates a new submission repository
func NewSubmissionRepo(pool *pgxpool.Pool) *SubmissionRepo {
	return &SubmissionRepo{pool: pool}
}

// GetByID retrieves a submission by its ID
func (r *SubmissionRepo) GetByID(ctx context.Context, id string) (*models.Submission, error) {
	query := `
		SELECT id, user_id, exercise_id, submitted_code, language,
		       execution_output, execution_error, execution_time_ms,
		       test_results, ai_feedback, ai_score, ai_evaluation_time_ms,
		       test_score, total_score, created_at
		FROM submissions
		WHERE id = $1
	`

	var submission models.Submission
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&submission.ID,
		&submission.UserID,
		&submission.ExerciseID,
		&submission.SubmittedCode,
		&submission.Language,
		&submission.ExecutionOutput,
		&submission.ExecutionError,
		&submission.ExecutionTimeMs,
		&submission.TestResults,
		&submission.AIFeedback,
		&submission.AIScore,
		&submission.AIEvaluationTimeMs,
		&submission.TestScore,
		&submission.TotalScore,
		&submission.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("submission not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get submission: %w", err)
	}

	return &submission, nil
}

// GetByExercise retrieves all submissions for an exercise by a user, ordered by attempt number
func (r *SubmissionRepo) GetByExercise(ctx context.Context, userID, exerciseID string) ([]*models.Submission, error) {
	query := `
		SELECT id, user_id, exercise_id, submitted_code, language,
		       execution_output, execution_error, execution_time_ms,
		       test_results, ai_feedback, ai_score, ai_evaluation_time_ms,
		       test_score, total_score, created_at
		FROM submissions
		WHERE user_id = $1 AND exercise_id = $2
		ORDER BY created_at ASC
	`

	rows, err := r.pool.Query(ctx, query, userID, exerciseID)
	if err != nil {
		return nil, fmt.Errorf("failed to query submissions: %w", err)
	}
	defer rows.Close()

	var submissions []*models.Submission
	for rows.Next() {
		var submission models.Submission
		err := rows.Scan(
			&submission.ID,
			&submission.UserID,
			&submission.ExerciseID,
			&submission.SubmittedCode,
			&submission.Language,
			&submission.ExecutionOutput,
			&submission.ExecutionError,
			&submission.ExecutionTimeMs,
			&submission.TestResults,
			&submission.AIFeedback,
			&submission.AIScore,
			&submission.AIEvaluationTimeMs,
			&submission.TestScore,
			&submission.TotalScore,
			&submission.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan submission: %w", err)
		}
		submissions = append(submissions, &submission)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating submissions: %w", err)
	}

	return submissions, nil
}

// GetLatestAttempt gets the most recent submission for a user+exercise
func (r *SubmissionRepo) GetLatestAttempt(ctx context.Context, userID, exerciseID string) (*models.Submission, error) {
	query := `
		SELECT id, user_id, exercise_id, submitted_code, language,
		       execution_output, execution_error, execution_time_ms,
		       test_results, ai_feedback, ai_score, ai_evaluation_time_ms,
		       test_score, total_score, created_at
		FROM submissions
		WHERE user_id = $1 AND exercise_id = $2
		ORDER BY created_at DESC
		LIMIT 1
	`

	var submission models.Submission
	err := r.pool.QueryRow(ctx, query, userID, exerciseID).Scan(
		&submission.ID,
		&submission.UserID,
		&submission.ExerciseID,
		&submission.SubmittedCode,
		&submission.Language,
		&submission.ExecutionOutput,
		&submission.ExecutionError,
		&submission.ExecutionTimeMs,
		&submission.TestResults,
		&submission.AIFeedback,
		&submission.AIScore,
		&submission.AIEvaluationTimeMs,
		&submission.TestScore,
		&submission.TotalScore,
		&submission.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // No previous attempts, not an error
		}
		return nil, fmt.Errorf("failed to get latest attempt: %w", err)
	}

	return &submission, nil
}

// Create creates a new submission
func (r *SubmissionRepo) Create(ctx context.Context, submission *models.Submission) error {
	query := `
		INSERT INTO submissions (
			id, user_id, exercise_id, submitted_code, language,
			execution_output, execution_error, execution_time_ms,
			test_results, ai_feedback, ai_score, ai_evaluation_time_ms,
			test_score, total_score
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING created_at
	`

	err := r.pool.QueryRow(ctx, query,
		submission.ID,
		submission.UserID,
		submission.ExerciseID,
		submission.SubmittedCode,
		submission.Language,
		submission.ExecutionOutput,
		submission.ExecutionError,
		submission.ExecutionTimeMs,
		submission.TestResults,
		submission.AIFeedback,
		submission.AIScore,
		submission.AIEvaluationTimeMs,
		submission.TestScore,
		submission.TotalScore,
	).Scan(&submission.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create submission: %w", err)
	}

	return nil
}

// CountAttempts counts how many attempts a user has made for an exercise
func (r *SubmissionRepo) CountAttempts(ctx context.Context, userID, exerciseID string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM submissions
		WHERE user_id = $1 AND exercise_id = $2
	`

	var count int
	err := r.pool.QueryRow(ctx, query, userID, exerciseID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count attempts: %w", err)
	}

	return count, nil
}

// GetByExercisePaginated retrieves submissions for an exercise with pagination
// Returns submissions, total count, and error
func (r *SubmissionRepo) GetByExercisePaginated(ctx context.Context, userID, exerciseID string, limit, offset int) ([]*models.Submission, int, error) {
	// First, get the total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM submissions
		WHERE user_id = $1 AND exercise_id = $2
	`
	err := r.pool.QueryRow(ctx, countQuery, userID, exerciseID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count submissions: %w", err)
	}

	// Then get the paginated results
	query := `
		SELECT id, user_id, exercise_id, submitted_code, language,
		       execution_output, execution_error, execution_time_ms,
		       test_results, ai_feedback, ai_score, ai_evaluation_time_ms,
		       test_score, total_score, created_at
		FROM submissions
		WHERE user_id = $1 AND exercise_id = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.pool.Query(ctx, query, userID, exerciseID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query paginated submissions: %w", err)
	}
	defer rows.Close()

	var submissions []*models.Submission
	for rows.Next() {
		var submission models.Submission
		err := rows.Scan(
			&submission.ID,
			&submission.UserID,
			&submission.ExerciseID,
			&submission.SubmittedCode,
			&submission.Language,
			&submission.ExecutionOutput,
			&submission.ExecutionError,
			&submission.ExecutionTimeMs,
			&submission.TestResults,
			&submission.AIFeedback,
			&submission.AIScore,
			&submission.AIEvaluationTimeMs,
			&submission.TestScore,
			&submission.TotalScore,
			&submission.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan submission: %w", err)
		}
		submissions = append(submissions, &submission)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating submissions: %w", err)
	}

	return submissions, total, nil
}