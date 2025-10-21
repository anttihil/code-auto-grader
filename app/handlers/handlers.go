package handlers

import (
	"net/http"
	"time"

	"github.com/anttihil/code-auto-grader/services"
	"github.com/anttihil/code-auto-grader/services/executor"
	"github.com/labstack/echo/v4"
)

// SubmissionHandler handles code execution and submission endpoints
type SubmissionHandler struct {
	codeExecutor *executor.CodeExecutionService
}

// NewSubmissionHandler creates a new submission handler with dependencies
func NewSubmissionHandler(exec *executor.CodeExecutionService) *SubmissionHandler {
	if exec == nil {
		panic("code executor is required")
	}
	return &SubmissionHandler{
		codeExecutor: exec,
	}
}

// ExerciseHandler handles exercise generation endpoints
type ExerciseHandler struct {
	gemini *services.GeminiService
}

// NewExerciseHandler creates a new exercise handler with dependencies
func NewExerciseHandler(gemini *services.GeminiService) *ExerciseHandler {
	if gemini == nil {
		panic("gemini service is required")
	}
	return &ExerciseHandler{
		gemini: gemini,
	}
}

// HealthCheck returns the health status of the API
func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// GetSampleExercise returns a sample coding exercise
func (h *ExerciseHandler) GetSampleExercise(c echo.Context) error {
	sampleExercise := map[string]interface{}{
		"title":       "Sum Two Numbers",
		"description": "Write a function that takes two numbers as parameters and returns their sum.",
		"difficulty":  "beginner",
		"language":    "javascript",
		"starterCode": "// Write a function called sum that takes two parameters: a and b\n// Return the sum of a and b\n\nfunction sum(a, b) {\n  // Your code here\n}\n\n// Test your function\nconsole.log(sum(2, 3)); // Should output: 5\nconsole.log(sum(10, 20)); // Should output: 30",
		"hints": []string{
			"Use the + operator to add two numbers",
			"Remember to return the result",
			"The function should work with any two numbers",
		},
		"expectedOutput": "The function should return the sum of two numbers",
		"testCases": []map[string]string{
			{"input": "sum(2, 3)", "expectedOutput": "5"},
			{"input": "sum(10, 20)", "expectedOutput": "30"},
			{"input": "sum(-5, 5)", "expectedOutput": "0"},
		},
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":  true,
		"exercise": sampleExercise,
	})
}

// GenerateExercise generates a new coding exercise using AI
func (h *ExerciseHandler) GenerateExercise(c echo.Context) error {
	var req struct {
		Difficulty string `json:"difficulty"`
		Topic      string `json:"topic"`
		Language   string `json:"language"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	// Set defaults
	if req.Difficulty == "" {
		req.Difficulty = "beginner"
	}
	if req.Topic == "" {
		req.Topic = "general"
	}
	if req.Language == "" {
		req.Language = "javascript"
	}

	exercise, err := h.gemini.GenerateExercise(req.Difficulty, req.Topic, req.Language)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to generate exercise",
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":  true,
		"exercise": exercise,
	})
}

// RunCode executes code without evaluation
func (h *SubmissionHandler) RunCode(c echo.Context) error {
	var req struct {
		Code     string `json:"code"`
		Language string `json:"language"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	if req.Code == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Code is required",
		})
	}

	if req.Language == "" {
		req.Language = "javascript"
	}

	result := h.codeExecutor.ExecuteCode(req.Code, req.Language)

	// Check if execution timed out due to queue capacity
	if !result.Success && result.ExecutionTime == "timeout" &&
		(result.Error == "Server is at maximum capacity, please try again in a few seconds") {
		c.Response().Header().Set("Retry-After", "5")
		return c.JSON(http.StatusServiceUnavailable, map[string]interface{}{
			"success":         false,
			"executionResult": result,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":         result.Success,
		"executionResult": result,
	})
}

// SubmitCode submits code for AI evaluation
func (h *SubmissionHandler) SubmitCode(c echo.Context) error {
	var req struct {
		Code     string `json:"code"`
		Exercise string `json:"exercise"`
		Language string `json:"language"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	if req.Code == "" || req.Exercise == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Code and exercise are required",
		})
	}

	if req.Language == "" {
		req.Language = "javascript"
	}

	// Execute the code
	executionResult := h.codeExecutor.ExecuteCode(req.Code, req.Language)

	// Check if execution timed out due to queue capacity
	if !executionResult.Success && executionResult.ExecutionTime == "timeout" &&
		(executionResult.Error == "Server is at maximum capacity, please try again in a few seconds") {
		c.Response().Header().Set("Retry-After", "5")
		return c.JSON(http.StatusServiceUnavailable, map[string]interface{}{
			"success":         false,
			"executionResult": executionResult,
		})
	}

	// Get AI feedback
	gemini := services.NewGeminiService()
	feedback, err := gemini.CheckSubmission(req.Exercise, req.Code, req.Language)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to process submission",
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":         executionResult.Success,
		"executionResult": executionResult,
		"feedback":        feedback,
	})
}

// GetPoolStatus returns the status of container pools
func (h *SubmissionHandler) GetPoolStatus(c echo.Context) error {
	status := h.codeExecutor.GetPoolStatus()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"pools":   status,
	})
}
