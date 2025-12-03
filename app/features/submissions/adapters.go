package submissions

import (
	"context"

	codeexecution "github.com/anttihil/code-auto-grader/features/code-execution"
	"github.com/anttihil/code-auto-grader/shared/gemini"
)

// GeminiAdapter adapts gemini.GeminiService to the AIEvaluator interface
type GeminiAdapter struct {
	gemini *gemini.GeminiService
}

// NewGeminiAdapter creates a new Gemini adapter
func NewGeminiAdapter(geminiService *gemini.GeminiService) *GeminiAdapter {
	return &GeminiAdapter{gemini: geminiService}
}

// CheckSubmissionWithContext implements the AIEvaluator interface
func (a *GeminiAdapter) CheckSubmissionWithContext(ctx context.Context, exercise, code, language string) (*Feedback, error) {
	// Call the underlying service
	serviceFeedback, err := a.gemini.CheckSubmissionWithContext(ctx, exercise, code, language)
	if err != nil {
		return nil, err
	}

	// Convert services.Feedback to domain Feedback
	return &Feedback{
		Passed:       serviceFeedback.Passed,
		Score:        serviceFeedback.Score,
		Feedback:     serviceFeedback.Feedback,
		Strengths:    serviceFeedback.Strengths,
		Improvements: serviceFeedback.Improvements,
		Correctness:  serviceFeedback.Correctness,
		CodeQuality:  serviceFeedback.CodeQuality,
		Suggestions:  serviceFeedback.Suggestions,
	}, nil
}

// ExecutorAdapter adapts codeexecution.CodeExecutionService to the CodeExecutor interface
type ExecutorAdapter struct {
	executor *codeexecution.CodeExecutionService
}

// NewExecutorAdapter creates a new executor adapter
func NewExecutorAdapter(exec *codeexecution.CodeExecutionService) *ExecutorAdapter {
	return &ExecutorAdapter{executor: exec}
}

// ExecuteCodeWithContext implements the CodeExecutor interface
func (a *ExecutorAdapter) ExecuteCodeWithContext(ctx context.Context, code, language string) *codeexecution.ExecutionResult {
	result := a.executor.ExecuteCodeWithContext(ctx, code, language)
	return &result
}

// GetPoolStatus implements the CodeExecutor interface
func (a *ExecutorAdapter) GetPoolStatus() map[string]interface{} {
	poolStatus := a.executor.GetPoolStatus()
	// Convert map[string]executor.PoolStatus to map[string]interface{}
	result := make(map[string]interface{}, len(poolStatus))
	for k, v := range poolStatus {
		result[k] = v
	}
	return result
}
