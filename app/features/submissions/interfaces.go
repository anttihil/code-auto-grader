package submissions

import (
	"context"

	codeexecution "github.com/anttihil/code-auto-grader/features/code-execution"
	"github.com/anttihil/code-auto-grader/models"
)

// CodeExecutor defines the interface for code execution
// This allows the grading domain to be independent of the concrete executor implementation
type CodeExecutor interface {
	ExecuteCodeWithContext(ctx context.Context, code, language string) *codeexecution.ExecutionResult
	GetPoolStatus() map[string]interface{}
}

// AIEvaluator defines the interface for AI-based code evaluation
type AIEvaluator interface {
	// CheckSubmissionWithContext evaluates submitted code and returns feedback
	// The exercise parameter should contain evaluation criteria/rubric
	CheckSubmissionWithContext(ctx context.Context, exercise, code, language string) (*Feedback, error)
}

// Feedback represents AI evaluation feedback
// This is defined here to decouple from services package
type Feedback struct {
	Passed       bool     `json:"passed"`
	Score        int      `json:"score"`
	Feedback     string   `json:"feedback"`
	Strengths    []string `json:"strengths"`
	Improvements []string `json:"improvements"`
	Correctness  string   `json:"correctness"`
	CodeQuality  string   `json:"codeQuality"`
	Suggestions  string   `json:"suggestions"`
}

// ToAIFeedback converts domain Feedback to models.AIFeedback
func (f *Feedback) ToAIFeedback() *models.AIFeedback {
	return &models.AIFeedback{
		Passed:       f.Passed,
		Score:        f.Score,
		Feedback:     f.Feedback,
		Strengths:    f.Strengths,
		Improvements: f.Improvements,
		Correctness:  f.Correctness,
		CodeQuality:  f.CodeQuality,
		Suggestions:  f.Suggestions,
	}
}
