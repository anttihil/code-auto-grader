// Package content provides HTTP handlers for exercise management endpoints
package content

import (
	"context"
	"log/slog"
	"net/http"

	domainContent "github.com/anttihil/code-auto-grader/domain/content"
	"github.com/anttihil/code-auto-grader/internal/logger"
	"github.com/labstack/echo/v4"
)

// ExerciseHandler handles HTTP requests for exercise endpoints
type ExerciseHandler struct {
	contentService *domainContent.ContentService
}

// NewExerciseHandler creates a new exercise handler
func NewExerciseHandler(contentService *domainContent.ContentService) *ExerciseHandler {
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
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "Exercise not found",
		})
	}

	// Spread exercise fields to response root level with submission data
	response := map[string]interface{}{
		"success":           true,
		"id":                view.Exercise.ID,
		"node_id":           view.Exercise.NodeID,
		"title":             view.Exercise.Title,
		"slug":              view.Exercise.Slug,
		"description":       view.Exercise.Description,
		"difficulty":        view.Exercise.Difficulty,
		"estimated_minutes": view.Exercise.EstimatedMinutes,
		"order_index":       view.Exercise.OrderIndex,
		"content":           view.Exercise.Content,
		"created_at":        view.Exercise.CreatedAt,
		"updated_at":        view.Exercise.UpdatedAt,
		"submissions":       view.Submissions,
		"best_score":        view.BestScore,
		"attempt_count":     view.AttemptCount,
	}

	return c.JSON(http.StatusOK, response)
}

// GenerateExercise handles POST /api/exercise/generate
// Generates a new coding exercise using AI
func (h *ExerciseHandler) GenerateExercise(c echo.Context) error {
	ctx, requestID := getRequestContext(c)

	var req struct {
		Difficulty string `json:"difficulty"`
		Topic      string `json:"topic"`
		Language   string `json:"language"`
	}

	if err := c.Bind(&req); err != nil {
		logger.Logger.Warn("invalid request body",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	result, err := h.contentService.GenerateExercise(ctx, domainContent.GenerateExerciseRequest{
		Difficulty: req.Difficulty,
		Topic:      req.Topic,
		Language:   req.Language,
	}, requestID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to generate exercise",
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":        true,
		"title":          result.Title,
		"description":    result.Description,
		"difficulty":     result.Difficulty,
		"language":       result.Language,
		"starterCode":    result.StarterCode,
		"hints":          result.Hints,
		"expectedOutput": result.ExpectedOutput,
		"testCases":      result.TestCases,
	})
}

// GetSampleExercise handles GET /api/exercise/sample
// Returns a sample coding exercise
func (h *ExerciseHandler) GetSampleExercise(c echo.Context) error {
	sample := h.contentService.GetSampleExercise()

	response := map[string]interface{}{
		"success":        true,
		"title":          sample.Title,
		"description":    sample.Description,
		"difficulty":     sample.Difficulty,
		"language":       sample.Language,
		"starterCode":    sample.StarterCode,
		"hints":          sample.Hints,
		"expectedOutput": sample.ExpectedOutput,
		"testCases":      sample.TestCases,
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
