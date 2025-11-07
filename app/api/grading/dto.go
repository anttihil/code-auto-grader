package grading

import "github.com/anttihil/code-auto-grader/models"

// SubmitCodeRequest represents the HTTP request body for code submission
type SubmitCodeRequest struct {
	Code     string `json:"code"`
	Language string `json:"language"`
}

// RunCodeRequest represents the HTTP request body for quick code execution
type RunCodeRequest struct {
	Code     string `json:"code"`
	Language string `json:"language"`
}

// SubmissionResponse represents the HTTP response for a code submission
type SubmissionResponse struct {
	SubmissionID    string              `json:"submission_id"`
	ExecutionOutput string              `json:"execution_output"`
	ExecutionError  *string             `json:"execution_error,omitempty"`
	TestResults     *models.TestResults `json:"test_results"`
	AIFeedback      *models.AIFeedback  `json:"ai_feedback"`
	TestScore       int                 `json:"test_score"`
	AIScore         int                 `json:"ai_score"`
	TotalScore      int                 `json:"total_score"`
	AttemptNumber   int                 `json:"attempt_number"`
}

// RunCodeResponse represents the HTTP response for code execution with test cases
type RunCodeResponse struct {
	ExecutionOutput string              `json:"execution_output"`
	ExecutionError  *string             `json:"execution_error,omitempty"`
	TestResults     *models.TestResults `json:"test_results"`
	TestScore       int                 `json:"test_score"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
