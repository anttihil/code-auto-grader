package repository

import (
	"context"
	"fmt"

	"github.com/anttihil/code-auto-grader/models"
	"github.com/jackc/pgx/v5"
)

// CreateSubmissionWithProgress creates a submission and updates user progress atomically
// This ensures data consistency when a user submits an exercise
func CreateSubmissionWithProgress(ctx context.Context, tx pgx.Tx, submission *models.Submission, progress *models.UserProgress) error {
	// Insert submission
	submissionQuery := `
		INSERT INTO submissions (
			id, user_id, exercise_id, submitted_code, language,
			execution_output, execution_error, execution_time_ms,
			test_results, ai_feedback, ai_score, ai_evaluation_time_ms,
			test_score, total_score, attempt_number
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING created_at
	`

	err := tx.QueryRow(ctx, submissionQuery,
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
		submission.AttemptNumber,
	).Scan(&submission.CreatedAt)

	if err != nil {
		return WrapDBError(err, "submission creation")
	}

	// Upsert user progress
	progressQuery := `
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
			mastery_score = EXCLUDED.mastery_score,
			is_mastered = EXCLUDED.is_mastered,
			completed_at = EXCLUDED.completed_at,
			last_activity_at = EXCLUDED.last_activity_at
		RETURNING created_at, updated_at
	`

	err = tx.QueryRow(ctx, progressQuery,
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
		return WrapDBError(err, "user progress update")
	}

	return nil
}

// CreateNodeWithPrerequisites creates a node and its prerequisites atomically
// This ensures referential integrity when creating node relationships
func CreateNodeWithPrerequisites(ctx context.Context, tx pgx.Tx, node *models.Node, prerequisiteIDs []string) error {
	// Insert node
	nodeQuery := `
		INSERT INTO nodes (
			id, tech_tree_id, name, slug, description, learning_objectives,
			order_index, estimated_minutes
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at, updated_at
	`

	err := tx.QueryRow(ctx, nodeQuery,
		node.ID,
		node.TechTreeID,
		node.Name,
		node.Slug,
		node.Description,
		node.LearningObjectives,
		node.OrderIndex,
		node.EstimatedMinutes,
	).Scan(&node.CreatedAt, &node.UpdatedAt)

	if err != nil {
		return WrapDBError(err, fmt.Sprintf("node %s", node.Name))
	}

	// Insert prerequisites
	if len(prerequisiteIDs) > 0 {
		prereqQuery := `
			INSERT INTO node_prerequisites (node_id, prerequisite_node_id)
			VALUES ($1, $2)
		`

		for _, prereqID := range prerequisiteIDs {
			_, err = tx.Exec(ctx, prereqQuery, node.ID, prereqID)
			if err != nil {
				return WrapDBError(err, fmt.Sprintf("prerequisite for node %s", node.Name))
			}
		}
	}

	return nil
}

// RecalculateMastery recalculates a user's mastery score for a node based on all submissions
// This runs in a transaction to ensure isolation
func RecalculateMastery(ctx context.Context, tx pgx.Tx, userID, nodeID string) error {
	// Calculate average score from all submissions for exercises in this node
	calcQuery := `
		SELECT
			AVG(s.total_score) as avg_score,
			COUNT(DISTINCT s.exercise_id) as completed_exercises
		FROM submissions s
		JOIN exercises e ON s.exercise_id = e.id
		WHERE s.user_id = $1
		  AND e.node_id = $2
		  AND s.total_score IS NOT NULL
		GROUP BY s.user_id, e.node_id
	`

	var avgScore *float64
	var completedExercises int
	err := tx.QueryRow(ctx, calcQuery, userID, nodeID).Scan(&avgScore, &completedExercises)
	if err != nil {
		// If no submissions yet, that's okay - just return
		if err == pgx.ErrNoRows {
			return nil
		}
		return fmt.Errorf("failed to calculate mastery: %w", err)
	}

	// Update user progress with calculated mastery
	updateQuery := `
		UPDATE user_progress
		SET mastery_score = $1,
		    exercises_completed = $2,
		    is_mastered = (COALESCE($1, 0) >= mastery_threshold),
		    last_activity_at = current_timestamp
		WHERE user_id = $3 AND node_id = $4
	`

	_, err = tx.Exec(ctx, updateQuery, avgScore, completedExercises, userID, nodeID)
	if err != nil {
		return WrapDBError(err, fmt.Sprintf("mastery update for user %s node %s", userID, nodeID))
	}

	return nil
}