// Package grading provides HTTP handlers for code submission and evaluation endpoints
package grading

import (
	"context"
	"log/slog"
	"net/http"

	domainGrading "github.com/anttihil/code-auto-grader/domain/grading"
	"github.com/anttihil/code-auto-grader/internal/logger"
	"github.com/labstack/echo/v4"
)

// SubmissionHandler handles HTTP requests for code submission and execution
type SubmissionHandler struct {
	gradingService *domainGrading.GradingService
}

// NewSubmissionHandler creates a new submission handler
func NewSubmissionHandler(gradingService *domainGrading.GradingService) *SubmissionHandler {
	if gradingService == nil {
		panic("grading service is required")
	}
	return &SubmissionHandler{
		gradingService: gradingService,
	}
}

// SubmitToExercise handles POST /api/submission/:exerciseId
// Submits code to a specific exercise with full evaluation and scoring
func (h *SubmissionHandler) SubmitToExercise(c echo.Context) error {
	ctx, requestID := getRequestContext(c)

	exerciseID := c.Param("exerciseId")
	// TODO: Get real user ID from auth context
	userID := "00000000-0000-0000-0000-000000000001" // Hardcoded for MVP

	logger.Logger.Info("submit to exercise endpoint called",
		slog.String("request_id", requestID),
		slog.String("exercise_id", exerciseID),
		slog.String("user_id", userID),
	)

	var req SubmitCodeRequest
	if err := c.Bind(&req); err != nil {
		logger.Logger.Error("failed to bind request body",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	if req.Code == "" {
		logger.Logger.Warn("missing required fields",
			slog.String("request_id", requestID),
		)
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "Code is required",
		})
	}

	// Set default language if not provided
	if req.Language == "" {
		req.Language = "javascript"
	}

	// Call domain service
	result, err := h.gradingService.SubmitCode(ctx, domainGrading.SubmitCodeRequest{
		UserID:     userID,
		ExerciseID: exerciseID,
		Code:       req.Code,
		Language:   req.Language,
	}, requestID)

	if err != nil {
		// Check for specific error types
		if err.Error() == "IDENTICAL_TO_STARTER: submission is identical to the starter code" {
			return c.JSON(http.StatusBadRequest, ErrorResponse{
				Success: false,
				Error:   "IDENTICAL_TO_STARTER",
				Message: "Your submission is identical to the starter code. Please implement the solution described in the exercise.",
			})
		}
		if err.Error() == "TRIVIAL_CHANGES: submission only contains trivial changes" {
			return c.JSON(http.StatusBadRequest, ErrorResponse{
				Success: false,
				Error:   "TRIVIAL_CHANGES",
				Message: "Your submission only contains trivial changes (comments or formatting). Please implement a solution to the exercise.",
			})
		}

		logger.Logger.Error("submission failed",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error:   "Failed to process submission",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, SubmissionResponse{
		Success:         true,
		SubmissionID:    result.SubmissionID,
		ExecutionOutput: result.ExecutionOutput,
		ExecutionError:  result.ExecutionError,
		TestResults:     result.TestResults,
		AIFeedback:      result.AIFeedback,
		TestScore:       result.TestScore,
		AIScore:         result.AIScore,
		TotalScore:      result.TotalScore,
		AttemptNumber:   result.AttemptNumber,
	})
}

// RunCode handles POST /api/code/run/:exerciseId
// Executes code with test cases (quick test without AI evaluation or saving)
func (h *SubmissionHandler) RunCode(c echo.Context) error {
	ctx, requestID := getRequestContext(c)

	exerciseID := c.Param("exerciseId")

	var req RunCodeRequest
	if err := c.Bind(&req); err != nil {
		logger.Logger.Warn("invalid request body",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	if req.Code == "" {
		logger.Logger.Warn("missing code",
			slog.String("request_id", requestID),
		)
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "Code is required",
		})
	}

	if req.Language == "" {
		req.Language = "javascript"
	}

	logger.Logger.Info("running code with test cases",
		slog.String("request_id", requestID),
		slog.String("exercise_id", exerciseID),
		slog.String("language", req.Language),
		slog.Int("code_length", len(req.Code)),
	)

	result, err := h.gradingService.RunCode(ctx, domainGrading.RunCodeRequest{
		ExerciseID: exerciseID,
		Code:       req.Code,
		Language:   req.Language,
	}, requestID)

	if err != nil {
		logger.Logger.Error("code run failed",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error:   "Failed to run code",
			Message: err.Error(),
		})
	}

	logger.Logger.Info("code run completed",
		slog.String("request_id", requestID),
		slog.Int("test_score", result.TestScore),
	)

	return c.JSON(http.StatusOK, RunCodeResponse{
		Success:         true,
		ExecutionOutput: result.ExecutionOutput,
		ExecutionError:  result.ExecutionError,
		TestResults:     result.TestResults,
		TestScore:       result.TestScore,
	})
}

// GetPoolStatus handles GET /api/pool/status
// Returns the status of code execution container pools
func (h *SubmissionHandler) GetPoolStatus(c echo.Context) error {
	status := h.gradingService.GetPoolStatus()

	// Spread pool status into response
	response := map[string]interface{}{
		"success": true,
	}
	for lang, poolStatus := range status {
		response[lang] = poolStatus
	}

	return c.JSON(http.StatusOK, response)
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
