package content

import (
	"github.com/anttihil/code-auto-grader/models"
	"github.com/google/uuid"
	"time"
)

// Re-export models for convenience
type Exercise = models.Exercise
type ExerciseContent = models.ExerciseContent
type TestCase = models.TestCase
type Hint = models.Hint
type CommonMistake = models.CommonMistake
type GradingRubric = models.GradingRubric

// User represents seed data for users
type User struct {
	ID           uuid.UUID  `json:"id"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	DisplayName  *string    `json:"display_name,omitempty"`
	PasswordHash *string    `json:"password_hash,omitempty"`
	CreatedAt    time.Time  `json:"-"`
	UpdatedAt    time.Time  `json:"-"`
}

// TechTree represents seed data for tech trees
type TechTree struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Slug            string    `json:"slug"`
	Description     *string   `json:"description,omitempty"`
	Language        string    `json:"language"`
	DifficultyLevel *string   `json:"difficulty_level,omitempty"`
	EstimatedHours  *int      `json:"estimated_hours,omitempty"`
	IsPublished     bool      `json:"is_published"`
	CreatedAt       time.Time `json:"-"`
	UpdatedAt       time.Time `json:"-"`
}

// Node represents seed data for nodes
type Node struct {
	ID                 uuid.UUID `json:"id"`
	TechTreeID         uuid.UUID `json:"tech_tree_id"`
	Name               string    `json:"name"`
	Slug               string    `json:"slug"`
	Description        *string   `json:"description,omitempty"`
	LearningObjectives []string  `json:"learning_objectives"`
	OrderIndex         int       `json:"order_index"`
	EstimatedMinutes   *int      `json:"estimated_minutes,omitempty"`
	Prerequisites      []string  `json:"prerequisites"` // UUID strings
	CreatedAt          time.Time `json:"-"`
	UpdatedAt          time.Time `json:"-"`
}

// SeedData represents all seed data from JSON files
type SeedData struct {
	Users     []User     `json:"users"`
	TechTrees []TechTree `json:"tech_trees"`
	Nodes     []Node     `json:"nodes"`
}
