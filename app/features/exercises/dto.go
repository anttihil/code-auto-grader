// Package content provides HTTP handlers for exercise management endpoints
package exercises

import (
	"time"

	"github.com/anttihil/code-auto-grader/models"
)

// ============================================================================
// REQUEST DTOs
// ============================================================================

// GenerateExerciseRequest represents the request body for exercise generation
type GenerateExerciseDTO struct {
	Difficulty string `json:"difficulty" validate:"required"`
	Topic      string `json:"topic" validate:"required"`
	Language   string `json:"language" validate:"required"`
}

// ============================================================================
// RESPONSE DTOs
// ============================================================================

// ExerciseResponse represents a single exercise with submission history
type ExerciseResponse struct {
	ID               string                  `json:"id"`
	NodeID           string                  `json:"node_id"`
	Title            string                  `json:"title"`
	Slug             string                  `json:"slug"`
	Description      string                  `json:"description"`
	Difficulty       string                  `json:"difficulty"`
	EstimatedMinutes int                     `json:"estimated_minutes"`
	OrderIndex       int                     `json:"order_index"`
	Content          models.ExerciseContent  `json:"content"`
	CreatedAt        time.Time               `json:"created_at"`
	UpdatedAt        time.Time               `json:"updated_at"`
	Submissions      []SubmissionSummary     `json:"submissions"`
	BestScore        float64                 `json:"best_score"`
	AttemptCount     int                     `json:"attempt_count"`
}

// SubmissionSummary represents a submission in the exercise history
type SubmissionSummary struct {
	ID                   string                 `json:"id"`
	ExerciseID           string                 `json:"exercise_id"`
	UserID               string                 `json:"user_id"`
	SubmittedCode        string                 `json:"submitted_code"`
	Language             string                 `json:"language"`
	ExecutionError       *string                `json:"execution_error,omitempty"`
	TestResults          TestResults            `json:"test_results"`
	AIFeedback           *AIFeedback            `json:"ai_feedback,omitempty"`
	AIScore              float64                `json:"ai_score"`
	ExecutionTimeMs      int                    `json:"execution_time_ms"`
	AIEvaluationTimeMs   int                    `json:"ai_evaluation_time_ms"`
	AttemptNumber        int                    `json:"attempt_number"`
	TotalScore           float64                `json:"total_score"`
	CreatedAt            time.Time              `json:"created_at"`
}

// TestResults represents test execution results
type TestResults struct {
	TotalTests  int          `json:"total_tests"`
	PassedTests int          `json:"passed_tests"`
	FailedTests int          `json:"failed_tests"`
	Details     []TestResult `json:"details"`
}

// TestResult represents a single test case result
type TestResult struct {
	Name   string  `json:"name"`
	Passed bool    `json:"passed"`
	Output *string `json:"output,omitempty"`
	Error  *string `json:"error,omitempty"`
}

// AIFeedback represents AI-generated feedback on a submission
type AIFeedback struct {
	Score        int      `json:"score"`
	Strengths    []string `json:"strengths"`
	Improvements []string `json:"improvements"`
	Suggestions  string   `json:"suggestions"`
}

// GenerateExerciseResponse represents the response for exercise generation
type GenerateExerciseResponse struct {
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	Difficulty     string     `json:"difficulty"`
	Language       string     `json:"language"`
	StarterCode    string     `json:"starter_code"`
	Hints          []Hint     `json:"hints"`
	ExpectedOutput string     `json:"expected_output"`
	TestCases      []TestCase `json:"test_cases"`
}

// Hint represents a hint for an exercise
type Hint struct {
	Order   int    `json:"order"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// TestCase represents a test case for an exercise
type TestCase struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

// SampleExerciseResponse represents the response for getting a sample exercise
type SampleExerciseResponse struct {
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	Difficulty     string     `json:"difficulty"`
	Language       string     `json:"language"`
	StarterCode    string     `json:"starter_code"`
	Hints          []Hint     `json:"hints"`
	ExpectedOutput string     `json:"expected_output"`
	TestCases      []TestCase `json:"test_cases"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
