package repository

import (
	"context"
	"fmt"

	"github.com/anttihil/code-auto-grader/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TechTreeRepo implements TechTreeRepository using pgx
type TechTreeRepo struct {
	pool *pgxpool.Pool
}

// NewTechTreeRepo creates a new tech tree repository
func NewTechTreeRepo(pool *pgxpool.Pool) *TechTreeRepo {
	return &TechTreeRepo{pool: pool}
}

// GetByID retrieves a tech tree by its ID
func (r *TechTreeRepo) GetByID(ctx context.Context, id string) (*models.TechTree, error) {
	query := `
		SELECT id, name, slug, description, language, difficulty_level,
		       estimated_hours, is_published, created_at, updated_at
		FROM tech_trees
		WHERE id = $1
	`

	var tree models.TechTree
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&tree.ID,
		&tree.Name,
		&tree.Slug,
		&tree.Description,
		&tree.Language,
		&tree.DifficultyLevel,
		&tree.EstimatedHours,
		&tree.IsPublished,
		&tree.CreatedAt,
		&tree.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("tech tree not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get tech tree: %w", err)
	}

	return &tree, nil
}

// GetBySlug retrieves a tech tree by its slug
func (r *TechTreeRepo) GetBySlug(ctx context.Context, slug string) (*models.TechTree, error) {
	query := `
		SELECT id, name, slug, description, language, difficulty_level,
		       estimated_hours, is_published, created_at, updated_at
		FROM tech_trees
		WHERE slug = $1
	`

	var tree models.TechTree
	err := r.pool.QueryRow(ctx, query, slug).Scan(
		&tree.ID,
		&tree.Name,
		&tree.Slug,
		&tree.Description,
		&tree.Language,
		&tree.DifficultyLevel,
		&tree.EstimatedHours,
		&tree.IsPublished,
		&tree.CreatedAt,
		&tree.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("tech tree not found: %s", slug)
		}
		return nil, fmt.Errorf("failed to get tech tree: %w", err)
	}

	return &tree, nil
}

// List retrieves all published tech trees
func (r *TechTreeRepo) List(ctx context.Context) ([]*models.TechTree, error) {
	query := `
		SELECT id, name, slug, description, language, difficulty_level,
		       estimated_hours, is_published, created_at, updated_at
		FROM tech_trees
		WHERE is_published = true
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query tech trees: %w", err)
	}
	defer rows.Close()

	var trees []*models.TechTree
	for rows.Next() {
		var tree models.TechTree
		err := rows.Scan(
			&tree.ID,
			&tree.Name,
			&tree.Slug,
			&tree.Description,
			&tree.Language,
			&tree.DifficultyLevel,
			&tree.EstimatedHours,
			&tree.IsPublished,
			&tree.CreatedAt,
			&tree.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tech tree: %w", err)
		}
		trees = append(trees, &tree)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tech trees: %w", err)
	}

	return trees, nil
}

// Create creates a new tech tree
func (r *TechTreeRepo) Create(ctx context.Context, tree *models.TechTree) error {
	query := `
		INSERT INTO tech_trees (
			id, name, slug, description, language, difficulty_level,
			estimated_hours, is_published
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at, updated_at
	`

	err := r.pool.QueryRow(ctx, query,
		tree.ID,
		tree.Name,
		tree.Slug,
		tree.Description,
		tree.Language,
		tree.DifficultyLevel,
		tree.EstimatedHours,
		tree.IsPublished,
	).Scan(&tree.CreatedAt, &tree.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create tech tree: %w", err)
	}

	return nil
}

// Update updates an existing tech tree
func (r *TechTreeRepo) Update(ctx context.Context, tree *models.TechTree) error {
	query := `
		UPDATE tech_trees
		SET name = $2,
		    slug = $3,
		    description = $4,
		    language = $5,
		    difficulty_level = $6,
		    estimated_hours = $7,
		    is_published = $8
		WHERE id = $1
		RETURNING updated_at
	`

	err := r.pool.QueryRow(ctx, query,
		tree.ID,
		tree.Name,
		tree.Slug,
		tree.Description,
		tree.Language,
		tree.DifficultyLevel,
		tree.EstimatedHours,
		tree.IsPublished,
	).Scan(&tree.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("tech tree not found: %s", tree.ID)
		}
		return fmt.Errorf("failed to update tech tree: %w", err)
	}

	return nil
}