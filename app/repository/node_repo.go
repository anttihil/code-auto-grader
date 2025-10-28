package repository

import (
	"context"
	"fmt"

	"github.com/anttihil/code-auto-grader/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NodeRepo implements NodeRepository using pgx
type NodeRepo struct {
	pool *pgxpool.Pool
}

// NewNodeRepo creates a new node repository
func NewNodeRepo(pool *pgxpool.Pool) *NodeRepo {
	return &NodeRepo{pool: pool}
}

// GetByID retrieves a node by its ID
func (r *NodeRepo) GetByID(ctx context.Context, id string) (*models.Node, error) {
	query := `
		SELECT id, tech_tree_id, name, slug, description, learning_objectives,
		       order_index, estimated_minutes, created_at, updated_at
		FROM nodes
		WHERE id = $1
	`

	var node models.Node
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&node.ID,
		&node.TechTreeID,
		&node.Name,
		&node.Slug,
		&node.Description,
		&node.LearningObjectives,
		&node.OrderIndex,
		&node.EstimatedMinutes,
		&node.CreatedAt,
		&node.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("node not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get node: %w", err)
	}

	return &node, nil
}

// GetByTechTree retrieves all nodes for a tech tree, ordered by order_index
func (r *NodeRepo) GetByTechTree(ctx context.Context, treeID string) ([]*models.Node, error) {
	query := `
		SELECT id, tech_tree_id, name, slug, description, learning_objectives,
		       order_index, estimated_minutes, created_at, updated_at
		FROM nodes
		WHERE tech_tree_id = $1
		ORDER BY order_index ASC
	`

	rows, err := r.pool.Query(ctx, query, treeID)
	if err != nil {
		return nil, fmt.Errorf("failed to query nodes: %w", err)
	}
	defer rows.Close()

	var nodes []*models.Node
	for rows.Next() {
		var node models.Node
		err := rows.Scan(
			&node.ID,
			&node.TechTreeID,
			&node.Name,
			&node.Slug,
			&node.Description,
			&node.LearningObjectives,
			&node.OrderIndex,
			&node.EstimatedMinutes,
			&node.CreatedAt,
			&node.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan node: %w", err)
		}
		nodes = append(nodes, &node)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating nodes: %w", err)
	}

	return nodes, nil
}

// GetPrerequisites retrieves prerequisite node IDs for a node
func (r *NodeRepo) GetPrerequisites(ctx context.Context, nodeID string) ([]string, error) {
	query := `
		SELECT prerequisite_node_id
		FROM node_prerequisites
		WHERE node_id = $1
		ORDER BY prerequisite_node_id
	`

	rows, err := r.pool.Query(ctx, query, nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to query prerequisites: %w", err)
	}
	defer rows.Close()

	var prerequisites []string
	for rows.Next() {
		var prereqID string
		if err := rows.Scan(&prereqID); err != nil {
			return nil, fmt.Errorf("failed to scan prerequisite: %w", err)
		}
		prerequisites = append(prerequisites, prereqID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating prerequisites: %w", err)
	}

	return prerequisites, nil
}

// GetPrerequisitesBatch retrieves prerequisites for multiple nodes in a single query
// This prevents N+1 query issues when loading prerequisites for multiple nodes
func (r *NodeRepo) GetPrerequisitesBatch(ctx context.Context, nodeIDs []string) (map[string][]string, error) {
	if len(nodeIDs) == 0 {
		return make(map[string][]string), nil
	}

	query := `
		SELECT node_id, prerequisite_node_id
		FROM node_prerequisites
		WHERE node_id = ANY($1)
		ORDER BY node_id, prerequisite_node_id
	`

	rows, err := r.pool.Query(ctx, query, nodeIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to query batch prerequisites: %w", err)
	}
	defer rows.Close()

	// Build map of node_id -> prerequisite IDs
	result := make(map[string][]string)
	for rows.Next() {
		var nodeID, prereqID string
		if err := rows.Scan(&nodeID, &prereqID); err != nil {
			return nil, fmt.Errorf("failed to scan prerequisite: %w", err)
		}
		result[nodeID] = append(result[nodeID], prereqID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating batch prerequisites: %w", err)
	}

	return result, nil
}

// Create creates a new node
func (r *NodeRepo) Create(ctx context.Context, node *models.Node) error {
	query := `
		INSERT INTO nodes (
			id, tech_tree_id, name, slug, description, learning_objectives,
			order_index, estimated_minutes
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at, updated_at
	`

	err := r.pool.QueryRow(ctx, query,
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
		return fmt.Errorf("failed to create node: %w", err)
	}

	return nil
}

// Update updates an existing node
func (r *NodeRepo) Update(ctx context.Context, node *models.Node) error {
	query := `
		UPDATE nodes
		SET name = $2,
		    slug = $3,
		    description = $4,
		    learning_objectives = $5,
		    order_index = $6,
		    estimated_minutes = $7
		WHERE id = $1
		RETURNING updated_at
	`

	err := r.pool.QueryRow(ctx, query,
		node.ID,
		node.Name,
		node.Slug,
		node.Description,
		node.LearningObjectives,
		node.OrderIndex,
		node.EstimatedMinutes,
	).Scan(&node.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("node not found: %s", node.ID)
		}
		return fmt.Errorf("failed to update node: %w", err)
	}

	return nil
}