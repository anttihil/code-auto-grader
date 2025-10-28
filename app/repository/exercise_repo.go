package repository

import (
	"context"
	"fmt"

	"github.com/anttihil/code-auto-grader/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ExerciseRepo implements ExerciseRepository using pgx
type ExerciseRepo struct {
	pool *pgxpool.Pool
}

// NewExerciseRepo creates a new exercise repository
func NewExerciseRepo(pool *pgxpool.Pool) *ExerciseRepo {
	return &ExerciseRepo{pool: pool}
}

// GetByID retrieves an exercise by its ID
func (r *ExerciseRepo) GetByID(ctx context.Context, id string) (*models.Exercise, error) {
	query := `
		SELECT id, node_id, title, slug, description, difficulty,
		       estimated_minutes, order_index, content, created_at, updated_at
		FROM exercises
		WHERE id = $1
	`

	var exercise models.Exercise
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&exercise.ID,
		&exercise.NodeID,
		&exercise.Title,
		&exercise.Slug,
		&exercise.Description,
		&exercise.Difficulty,
		&exercise.EstimatedMinutes,
		&exercise.OrderIndex,
		&exercise.Content,
		&exercise.CreatedAt,
		&exercise.UpdatedAt,
	)

	if err != nil {
		return nil, WrapDBError(err, fmt.Sprintf("exercise %s", id))
	}

	return &exercise, nil
}

// GetByNode retrieves all exercises for a node, ordered by order_index
func (r *ExerciseRepo) GetByNode(ctx context.Context, nodeID string) ([]*models.Exercise, error) {
	query := `
		SELECT id, node_id, title, slug, description, difficulty,
		       estimated_minutes, order_index, content, created_at, updated_at
		FROM exercises
		WHERE node_id = $1
		ORDER BY order_index ASC
	`

	rows, err := r.pool.Query(ctx, query, nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to query exercises: %w", err)
	}
	defer rows.Close()

	var exercises []*models.Exercise
	for rows.Next() {
		var exercise models.Exercise
		err := rows.Scan(
			&exercise.ID,
			&exercise.NodeID,
			&exercise.Title,
			&exercise.Slug,
			&exercise.Description,
			&exercise.Difficulty,
			&exercise.EstimatedMinutes,
			&exercise.OrderIndex,
			&exercise.Content,
			&exercise.CreatedAt,
			&exercise.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan exercise: %w", err)
		}
		exercises = append(exercises, &exercise)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating exercises: %w", err)
	}

	return exercises, nil
}

// GetBySlug retrieves an exercise by node ID and slug
func (r *ExerciseRepo) GetBySlug(ctx context.Context, nodeID, slug string) (*models.Exercise, error) {
	query := `
		SELECT id, node_id, title, slug, description, difficulty,
		       estimated_minutes, order_index, content, created_at, updated_at
		FROM exercises
		WHERE node_id = $1 AND slug = $2
	`

	var exercise models.Exercise
	err := r.pool.QueryRow(ctx, query, nodeID, slug).Scan(
		&exercise.ID,
		&exercise.NodeID,
		&exercise.Title,
		&exercise.Slug,
		&exercise.Description,
		&exercise.Difficulty,
		&exercise.EstimatedMinutes,
		&exercise.OrderIndex,
		&exercise.Content,
		&exercise.CreatedAt,
		&exercise.UpdatedAt,
	)

	if err != nil {
		return nil, WrapDBError(err, fmt.Sprintf("exercise %s/%s", nodeID, slug))
	}

	return &exercise, nil
}

// Create creates a new exercise
func (r *ExerciseRepo) Create(ctx context.Context, exercise *models.Exercise) error {
	query := `
		INSERT INTO exercises (
			id, node_id, title, slug, description, difficulty,
			estimated_minutes, order_index, content
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING created_at, updated_at
	`

	err := r.pool.QueryRow(ctx, query,
		exercise.ID,
		exercise.NodeID,
		exercise.Title,
		exercise.Slug,
		exercise.Description,
		exercise.Difficulty,
		exercise.EstimatedMinutes,
		exercise.OrderIndex,
		exercise.Content,
	).Scan(&exercise.CreatedAt, &exercise.UpdatedAt)

	if err != nil {
		return WrapDBError(err, fmt.Sprintf("exercise %s", exercise.Title))
	}

	return nil
}

// Update updates an existing exercise
func (r *ExerciseRepo) Update(ctx context.Context, exercise *models.Exercise) error {
	query := `
		UPDATE exercises
		SET title = $2,
		    slug = $3,
		    description = $4,
		    difficulty = $5,
		    estimated_minutes = $6,
		    order_index = $7,
		    content = $8
		WHERE id = $1
		RETURNING updated_at
	`

	err := r.pool.QueryRow(ctx, query,
		exercise.ID,
		exercise.Title,
		exercise.Slug,
		exercise.Description,
		exercise.Difficulty,
		exercise.EstimatedMinutes,
		exercise.OrderIndex,
		exercise.Content,
	).Scan(&exercise.UpdatedAt)

	if err != nil {
		return WrapDBError(err, fmt.Sprintf("exercise %s", exercise.ID))
	}

	return nil
}