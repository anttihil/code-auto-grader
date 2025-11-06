package grading

import "github.com/anttihil/code-auto-grader/models"

// SubmitCodeRequest represents a code submission request
type SubmitCodeRequest struct {
	UserID     string
	ExerciseID string
	Code       string
	Language   string
}

// SubmissionResult represents the complete result of a code submission
type SubmissionResult struct {
	SubmissionID    string
	ExecutionOutput string
	ExecutionError  *string
	TestResults     *models.TestResults
	AIFeedback      *models.AIFeedback
	TestScore       int
	AIScore         int
	TotalScore      int
	AttemptNumber   int
}

// RunCodeRequest represents a code execution request with test cases (no AI evaluation or saving)
type RunCodeRequest struct {
	ExerciseID string
	Code       string
	Language   string
}

// RunCodeResult represents the result of a code execution with test cases
type RunCodeResult struct {
	ExecutionOutput string
	ExecutionError  *string
	TestResults     *models.TestResults
	TestScore       int
}

// CodeSimilarityLevel represents how similar submitted code is to starter code
type CodeSimilarityLevel int

const (
	CodeIdentical       CodeSimilarityLevel = iota // Exact match
	CodeTrivialChanges                             // Only whitespace/comments
	CodeMinimalChanges                             // <10% change
	CodeSubstantial                                // Sufficient changes
)

// MasteryThreshold is the score threshold for achieving mastery (80%)
const MasteryThreshold = 80
