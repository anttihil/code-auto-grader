// Package learningmaterials provides HTTP handlers and business logic for learning materials
package learningmaterials

import "time"

// LearningMaterialResponse is the DTO for a learning material
type LearningMaterialResponse struct {
	ID         string    `json:"id"`
	NodeID     string    `json:"node_id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	OrderIndex int       `json:"order_index"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// CreateLearningMaterialRequest is the DTO for creating a learning material
type CreateLearningMaterialRequest struct {
	Title      string `json:"title" validate:"required"`
	Content    string `json:"content" validate:"required"`
	OrderIndex int    `json:"order_index"`
}

// UpdateLearningMaterialRequest is the DTO for updating a learning material
type UpdateLearningMaterialRequest struct {
	Title      *string `json:"title,omitempty"`
	Content    *string `json:"content,omitempty"`
	OrderIndex *int    `json:"order_index,omitempty"`
}

// ErrorResponse is a standard error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
