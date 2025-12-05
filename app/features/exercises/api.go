// Package content provides HTTP handlers for exercise management endpoints
package exercises

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/anttihil/code-auto-grader/internal/logger"
	"github.com/anttihil/code-auto-grader/models"
	"github.com/labstack/echo/v4"
)

// ExerciseHandler handles HTTP requests for exercise endpoints
type ExerciseHandler struct {
	contentService *ContentService
}

// NewExerciseHandler creates a new exercise handler
func NewExerciseHandler(contentService *ContentService) *ExerciseHandler {
	if contentService == nil {
		panic("content service is required")
	}
	return &ExerciseHandler{
		contentService: contentService,
	}
}

// GetExercise handles GET /api/exercise/:exerciseId
// Fetches an exercise with submission history
func (h *ExerciseHandler) GetExercise(c echo.Context) error {
	ctx, requestID := getRequestContext(c)

	exerciseID := c.Param("exerciseId")
	// TODO: Get real user ID from auth context
	userID := "00000000-0000-0000-0000-000000000001" // Hardcoded for MVP

	view, err := h.contentService.GetExercise(ctx, exerciseID, userID, requestID)
	if err != nil {
		logger.Logger.Error("failed to fetch exercise",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		return c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Exercise not found",
			Message: err.Error(),
		})
	}

	// Convert submissions to DTO format
	submissions := make([]SubmissionSummary, len(view.Submissions))
	for i, sub := range view.Submissions {
		aiScore := 0.0
		if sub.AIScore != nil {
			aiScore = float64(*sub.AIScore)
		}

		execTimeMs := 0
		if sub.ExecutionTimeMs != nil {
			execTimeMs = *sub.ExecutionTimeMs
		}

		aiEvalTimeMs := 0
		if sub.AIEvaluationTimeMs != nil {
			aiEvalTimeMs = *sub.AIEvaluationTimeMs
		}

		totalScore := 0.0
		if sub.TotalScore != nil {
			totalScore = float64(*sub.TotalScore)
		}

		submissions[i] = SubmissionSummary{
			ID:                 sub.ID.String(),
			ExerciseID:         sub.ExerciseID.String(),
			UserID:             sub.UserID.String(),
			SubmittedCode:      sub.SubmittedCode,
			Language:           sub.Language,
			ExecutionError:     sub.ExecutionError,
			TestResults:        convertTestResults(sub.TestResults),
			AIFeedback:         convertAIFeedback(sub.AIFeedback),
			AIScore:            aiScore,
			ExecutionTimeMs:    execTimeMs,
			AIEvaluationTimeMs: aiEvalTimeMs,
			TotalScore:         totalScore,
			CreatedAt:          sub.CreatedAt,
		}
	}

	estimatedMins := 0
	if view.Exercise.EstimatedMinutes != nil {
		estimatedMins = *view.Exercise.EstimatedMinutes
	}

	response := ExerciseResponse{
		ID:               view.Exercise.ID.String(),
		NodeID:           view.Exercise.NodeID.String(),
		Title:            view.Exercise.Title,
		Slug:             view.Exercise.Slug,
		Description:      view.Exercise.Description,
		Difficulty:       view.Exercise.Difficulty,
		EstimatedMinutes: estimatedMins,
		OrderIndex:       view.Exercise.OrderIndex,
		Content:          view.Exercise.Content,
		CreatedAt:        view.Exercise.CreatedAt,
		UpdatedAt:        view.Exercise.UpdatedAt,
		Submissions:      submissions,
		BestScore:        float64(view.BestScore),
		AttemptCount:     view.AttemptCount,
	}

	return c.JSON(http.StatusOK, response)
}

// GenerateExercise handles POST /api/exercise/generate
// Generates a new coding exercise using AI
func (h *ExerciseHandler) GenerateExercise(c echo.Context) error {
	ctx, requestID := getRequestContext(c)

	var req GenerateExerciseDTO

	if err := c.Bind(&req); err != nil {
		logger.Logger.Warn("invalid request body",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	result, err := h.contentService.GenerateExercise(ctx, GenerateExerciseRequest{
		Difficulty: req.Difficulty,
		Topic:      req.Topic,
		Language:   req.Language,
	}, requestID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to generate exercise",
			Message: err.Error(),
		})
	}

	// Convert hints (from strings to structured hints)
	hints := make([]Hint, len(result.Hints))
	for i, hintText := range result.Hints {
		hints[i] = Hint{
			Order:   i + 1,
			Title:   "",
			Content: hintText,
		}
	}

	testCases := make([]TestCase, len(result.TestCases))
	for i, tc := range result.TestCases {
		testCases[i] = TestCase{
			Name: tc.Name,
			Code: tc.Code,
		}
	}

	response := GenerateExerciseResponse{
		Title:          result.Title,
		Description:    result.Description,
		Difficulty:     result.Difficulty,
		Language:       result.Language,
		StarterCode:    result.StarterCode,
		Hints:          hints,
		ExpectedOutput: result.ExpectedOutput,
		TestCases:      testCases,
	}

	return c.JSON(http.StatusOK, response)
}

// GetSampleExercise handles GET /api/exercise/sample
// Returns a sample coding exercise
func (h *ExerciseHandler) GetSampleExercise(c echo.Context) error {
	sample := h.contentService.GetSampleExercise()

	// Convert hints (from strings to structured hints)
	hints := make([]Hint, len(sample.Hints))
	for i, hintText := range sample.Hints {
		hints[i] = Hint{
			Order:   i + 1,
			Title:   "",
			Content: hintText,
		}
	}

	testCases := make([]TestCase, len(sample.TestCases))
	for i, tc := range sample.TestCases {
		testCases[i] = TestCase{
			Name: tc.Name,
			Code: tc.Code,
		}
	}

	response := SampleExerciseResponse{
		Title:          sample.Title,
		Description:    sample.Description,
		Difficulty:     sample.Difficulty,
		Language:       sample.Language,
		StarterCode:    sample.StarterCode,
		Hints:          hints,
		ExpectedOutput: sample.ExpectedOutput,
		TestCases:      testCases,
	}

	return c.JSON(http.StatusOK, response)
}

// Helper functions to convert domain models to DTOs

func convertTestResults(tr *models.TestResults) TestResults {
	if tr == nil {
		return TestResults{}
	}

	details := make([]TestResult, len(tr.Results))
	for i, r := range tr.Results {
		details[i] = TestResult{
			Name:   r.Name,
			Passed: r.Status == "passed",
			Error:  r.Error,
		}
	}

	return TestResults{
		TotalTests:  tr.TotalTests,
		PassedTests: tr.PassedTests,
		FailedTests: tr.FailedTests,
		Details:     details,
	}
}

func convertAIFeedback(af *models.AIFeedback) *AIFeedback {
	if af == nil {
		return nil
	}

	return &AIFeedback{
		Score:        af.Score,
		Strengths:    af.Strengths,
		Improvements: af.Improvements,
		Suggestions:  af.Suggestions,
	}
}

// getRequestContext extracts request context with request ID from Echo context
func getRequestContext(c echo.Context) (context.Context, string) {
	requestID := "unknown"
	if id, ok := c.Get("request_id").(string); ok {
		requestID = id
	}
	ctx := context.WithValue(c.Request().Context(), "request_id", requestID)
	return ctx, requestID
}
