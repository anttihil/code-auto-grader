package repository

import (
	"context"
	"fmt"

	"github.com/anttihil/code-auto-grader/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

// LearningMaterialRepo implements LearningMaterialRepository using pgx
type LearningMaterialRepo struct {
	pool *pgxpool.Pool
}

// NewLearningMaterialRepo creates a new learning material repository
func NewLearningMaterialRepo(pool *pgxpool.Pool) *LearningMaterialRepo {
	return &LearningMaterialRepo{pool: pool}
}

// GetByNodeID retrieves all learning materials for a node, ordered by order_index
func (r *LearningMaterialRepo) GetByNodeID(ctx context.Context, nodeID string) ([]*models.LearningMaterial, error) {
	query := `
		SELECT id, node_id, title, content, order_index, created_at, updated_at
		FROM learning_materials
		WHERE node_id = $1
		ORDER BY order_index ASC
	`

	rows, err := r.pool.Query(ctx, query, nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to query learning materials: %w", err)
	}
	defer rows.Close()

	var materials []*models.LearningMaterial
	for rows.Next() {
		var material models.LearningMaterial
		err := rows.Scan(
			&material.ID,
			&material.NodeID,
			&material.Title,
			&material.Content,
			&material.OrderIndex,
			&material.CreatedAt,
			&material.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan learning material: %w", err)
		}
		materials = append(materials, &material)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating learning materials: %w", err)
	}

	return materials, nil
}

// GetByID retrieves a learning material by its ID
func (r *LearningMaterialRepo) GetByID(ctx context.Context, id string) (*models.LearningMaterial, error) {
	query := `
		SELECT id, node_id, title, content, order_index, created_at, updated_at
		FROM learning_materials
		WHERE id = $1
	`

	var material models.LearningMaterial
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&material.ID,
		&material.NodeID,
		&material.Title,
		&material.Content,
		&material.OrderIndex,
		&material.CreatedAt,
		&material.UpdatedAt,
	)

	if err != nil {
		return nil, WrapDBError(err, fmt.Sprintf("learning material %s", id))
	}

	return &material, nil
}

// Create creates a new learning material
func (r *LearningMaterialRepo) Create(ctx context.Context, material *models.LearningMaterial) error {
	query := `
		INSERT INTO learning_materials (id, node_id, title, content, order_index)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at, updated_at
	`

	err := r.pool.QueryRow(ctx, query,
		material.ID,
		material.NodeID,
		material.Title,
		material.Content,
		material.OrderIndex,
	).Scan(&material.CreatedAt, &material.UpdatedAt)

	if err != nil {
		return WrapDBError(err, fmt.Sprintf("learning material %s", material.Title))
	}

	return nil
}

// Update updates an existing learning material
func (r *LearningMaterialRepo) Update(ctx context.Context, material *models.LearningMaterial) error {
	query := `
		UPDATE learning_materials
		SET title = $2,
		    content = $3,
		    order_index = $4
		WHERE id = $1
		RETURNING updated_at
	`

	err := r.pool.QueryRow(ctx, query,
		material.ID,
		material.Title,
		material.Content,
		material.OrderIndex,
	).Scan(&material.UpdatedAt)

	if err != nil {
		return WrapDBError(err, fmt.Sprintf("learning material %s", material.ID))
	}

	return nil
}

// Delete deletes a learning material by its ID
func (r *LearningMaterialRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM learning_materials WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return WrapDBError(err, fmt.Sprintf("learning material %s", id))
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("learning material %s not found", id)
	}

	return nil
}
