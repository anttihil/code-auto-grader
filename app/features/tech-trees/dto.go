// Package learning provides HTTP handlers for tech tree navigation and progress tracking endpoints
package techtrees

import (
	"time"

	"github.com/anttihil/code-auto-grader/models"
)

// ============================================================================
// RESPONSE DTOs
// ============================================================================

// TechTreeResponse represents a tech tree with all nodes and user progress
type TechTreeResponse struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	Slug        string               `json:"slug"`
	Description string               `json:"description"`
	Language    string               `json:"language"`
	CreatedAt   time.Time            `json:"created_at"`
	Nodes       []NodeWithProgressDTO `json:"nodes"`
}

// NodeWithProgress represents a node with embedded progress data
type NodeWithProgressDTO struct {
	ID                 string            `json:"id"`
	TechTreeID         string            `json:"tech_tree_id"`
	Name               string            `json:"name"`
	Slug               string            `json:"slug"`
	Description        string            `json:"description"`
	LearningObjectives []string          `json:"learning_objectives"`
	OrderIndex         int               `json:"order_index"`
	EstimatedMinutes   int               `json:"estimated_minutes"`
	CreatedAt          time.Time         `json:"created_at"`
	UpdatedAt          time.Time         `json:"updated_at"`
	Progress           ProgressData      `json:"progress"`
}

// ProgressData represents user progress for a node
type ProgressData struct {
	NodeID              string             `json:"node_id"`
	ExercisesCompleted  int                `json:"exercises_completed"`
	TotalExercises      int                `json:"total_exercises"`
	MasteryAchieved     bool               `json:"mastery_achieved"`
	IsUnlocked          bool               `json:"is_unlocked"`
	BestScores          map[string]float64 `json:"best_scores"`
}

// NodeResponse represents a node with its exercises and user progress
type NodeResponse struct {
	ID                 string                 `json:"id"`
	TechTreeID         string                 `json:"tech_tree_id"`
	Name               string                 `json:"name"`
	Slug               string                 `json:"slug"`
	Description        string                 `json:"description"`
	LearningObjectives []string               `json:"learning_objectives"`
	OrderIndex         int                    `json:"order_index"`
	EstimatedMinutes   int                    `json:"estimated_minutes"`
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at"`
	Progress           ProgressData           `json:"progress"`
	Exercises          []ExerciseSummary      `json:"exercises"`
}

// ExerciseSummary represents a summary of an exercise
type ExerciseSummary struct {
	ID               string                 `json:"id"`
	NodeID           string                 `json:"node_id"`
	Title            string                 `json:"title"`
	Slug             string                 `json:"slug"`
	Description      string                 `json:"description"`
	Difficulty       string                 `json:"difficulty"`
	EstimatedMinutes int                    `json:"estimated_minutes"`
	OrderIndex       int                    `json:"order_index"`
	Content          models.ExerciseContent `json:"content"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
