package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/anttihil/code-auto-grader/internal/logger"
	"github.com/anttihil/code-auto-grader/repository"
	"github.com/anttihil/code-auto-grader/services"
	"github.com/anttihil/code-auto-grader/services/executor"
	"github.com/labstack/echo/v4"
)

// SubmissionHandler handles code execution and submission endpoints
type SubmissionHandler struct {
	codeExecutor     *executor.CodeExecutionService
	gemini           *services.GeminiService
	submissionRepo   repository.SubmissionRepository
	exerciseRepo     repository.ExerciseRepository
	userProgressRepo repository.UserProgressRepository
}

// NewSubmissionHandler creates a new submission handler with dependencies
func NewSubmissionHandler(
	exec *executor.CodeExecutionService,
	gemini *services.GeminiService,
	submissionRepo repository.SubmissionRepository,
	exerciseRepo repository.ExerciseRepository,
	userProgressRepo repository.UserProgressRepository,
) *SubmissionHandler {
	if exec == nil {
		panic("code executor is required")
	}
	if gemini == nil {
		panic("gemini service is required")
	}
	if submissionRepo == nil {
		panic("submission repository is required")
	}
	if exerciseRepo == nil {
		panic("exercise repository is required")
	}
	if userProgressRepo == nil {
		panic("user progress repository is required")
	}
	return &SubmissionHandler{
		codeExecutor:     exec,
		gemini:           gemini,
		submissionRepo:   submissionRepo,
		exerciseRepo:     exerciseRepo,
		userProgressRepo: userProgressRepo,
	}
}

// ExerciseHandler handles exercise generation endpoints
type ExerciseHandler struct {
	gemini         *services.GeminiService
	exerciseRepo   repository.ExerciseRepository
	submissionRepo repository.SubmissionRepository
}

// NewExerciseHandler creates a new exercise handler with dependencies
func NewExerciseHandler(
	gemini *services.GeminiService,
	exerciseRepo repository.ExerciseRepository,
	submissionRepo repository.SubmissionRepository,
) *ExerciseHandler {
	if gemini == nil {
		panic("gemini service is required")
	}
	if exerciseRepo == nil {
		panic("exercise repository is required")
	}
	if submissionRepo == nil {
		panic("submission repository is required")
	}
	return &ExerciseHandler{
		gemini:         gemini,
		exerciseRepo:   exerciseRepo,
		submissionRepo: submissionRepo,
	}
}

// TechTreeHandler handles tech tree navigation endpoints
type TechTreeHandler struct {
	techTreeRepo     repository.TechTreeRepository
	nodeRepo         repository.NodeRepository
	userProgressRepo repository.UserProgressRepository
}

// NewTechTreeHandler creates a new tech tree handler with dependencies
func NewTechTreeHandler(
	techTreeRepo repository.TechTreeRepository,
	nodeRepo repository.NodeRepository,
	userProgressRepo repository.UserProgressRepository,
) *TechTreeHandler {
	if techTreeRepo == nil {
		panic("tech tree repository is required")
	}
	if nodeRepo == nil {
		panic("node repository is required")
	}
	if userProgressRepo == nil {
		panic("user progress repository is required")
	}
	return &TechTreeHandler{
		techTreeRepo:     techTreeRepo,
		nodeRepo:         nodeRepo,
		userProgressRepo: userProgressRepo,
	}
}

// NodeHandler handles node-related endpoints
type NodeHandler struct {
	nodeRepo         repository.NodeRepository
	exerciseRepo     repository.ExerciseRepository
	userProgressRepo repository.UserProgressRepository
}

// NewNodeHandler creates a new node handler with dependencies
func NewNodeHandler(
	nodeRepo repository.NodeRepository,
	exerciseRepo repository.ExerciseRepository,
	userProgressRepo repository.UserProgressRepository,
) *NodeHandler {
	if nodeRepo == nil {
		panic("node repository is required")
	}
	if exerciseRepo == nil {
		panic("exercise repository is required")
	}
	if userProgressRepo == nil {
		panic("user progress repository is required")
	}
	return &NodeHandler{
		nodeRepo:         nodeRepo,
		exerciseRepo:     exerciseRepo,
		userProgressRepo: userProgressRepo,
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

	logger.Logger.Info("generating exercise",
		slog.String("request_id", requestID),
		slog.String("difficulty", req.Difficulty),
		slog.String("topic", req.Topic),
		slog.String("language", req.Language),
	)

	exercise, err := h.gemini.GenerateExerciseWithContext(ctx, req.Difficulty, req.Topic, req.Language)
	if err != nil {
		logger.Logger.Error("failed to generate exercise",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to generate exercise",
			"message": err.Error(),
		})
	}

	logger.Logger.Info("exercise generated successfully",
		slog.String("request_id", requestID),
	)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":  true,
		"exercise": exercise,
	})
}

// RunCode executes code without evaluation
func (h *SubmissionHandler) RunCode(c echo.Context) error {
	ctx, requestID := getRequestContext(c)

	var req struct {
		Code     string `json:"code"`
		Language string `json:"language"`
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

	if req.Code == "" {
		logger.Logger.Warn("missing code",
			slog.String("request_id", requestID),
		)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Code is required",
		})
	}

	if req.Language == "" {
		req.Language = "javascript"
	}

	logger.Logger.Info("running code",
		slog.String("request_id", requestID),
		slog.String("language", req.Language),
		slog.Int("code_length", len(req.Code)),
	)

	result := h.codeExecutor.ExecuteCodeWithContext(ctx, req.Code, req.Language)

	// Check if execution timed out due to queue capacity
	if !result.Success && result.ExecutionTime == "timeout" &&
		(result.Error == "Server is at maximum capacity, please try again in a few seconds") {
		logger.Logger.Warn("server at capacity",
			slog.String("request_id", requestID),
			slog.String("language", req.Language),
		)
		c.Response().Header().Set("Retry-After", "5")
		return c.JSON(http.StatusServiceUnavailable, map[string]interface{}{
			"success":         false,
			"executionResult": result,
		})
	}

	logger.Logger.Info("code execution completed",
		slog.String("request_id", requestID),
		slog.Bool("success", result.Success),
	)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":         result.Success,
		"executionResult": result,
	})
}

// SubmitCode submits code for AI evaluation
func (h *SubmissionHandler) SubmitCode(c echo.Context) error {
	ctx, requestID := getRequestContext(c)

	logger.Logger.Info("submit code endpoint called",
		slog.String("request_id", requestID),
	)

	var req struct {
		Code     string `json:"code"`
		Exercise string `json:"exercise"`
		Language string `json:"language"`
	}

	if err := c.Bind(&req); err != nil {
		logger.Logger.Error("failed to bind request body",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	if req.Code == "" || req.Exercise == "" {
		logger.Logger.Warn("missing required fields",
			slog.String("request_id", requestID),
			slog.Bool("has_code", req.Code != ""),
			slog.Bool("has_exercise", req.Exercise != ""),
		)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Code and exercise are required",
		})
	}

	if req.Language == "" {
		req.Language = "javascript"
	}

	logger.Logger.Info("submitting code for evaluation",
		slog.String("request_id", requestID),
		slog.String("language", req.Language),
		slog.Int("code_length", len(req.Code)),
		slog.Int("exercise_length", len(req.Exercise)),
	)

	// Execute the code
	logger.Logger.Debug("executing code",
		slog.String("request_id", requestID),
	)
	executionResult := h.codeExecutor.ExecuteCodeWithContext(ctx, req.Code, req.Language)
	logger.Logger.Info("code execution completed",
		slog.String("request_id", requestID),
		slog.Bool("success", executionResult.Success),
	)

	// Check if execution timed out due to queue capacity
	if !executionResult.Success && executionResult.ExecutionTime == "timeout" &&
		(executionResult.Error == "Server is at maximum capacity, please try again in a few seconds") {
		logger.Logger.Warn("server at capacity during submission",
			slog.String("request_id", requestID),
			slog.String("language", req.Language),
		)
		c.Response().Header().Set("Retry-After", "5")
		return c.JSON(http.StatusServiceUnavailable, map[string]interface{}{
			"success":         false,
			"executionResult": executionResult,
		})
	}

	// Get AI feedback
	logger.Logger.Debug("requesting AI feedback from Gemini",
		slog.String("request_id", requestID),
	)
	feedback, err := h.gemini.CheckSubmissionWithContext(ctx, req.Exercise, req.Code, req.Language)
	if err != nil {
		logger.Logger.Error("failed to get AI feedback",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to process submission",
			"message": err.Error(),
		})
	}

	logger.Logger.Info("submit code completed successfully",
		slog.String("request_id", requestID),
	)
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

// ============================================
// Tech Tree Handler Methods
// ============================================

// GetTechTree fetches a tech tree with all nodes and user progress
func (h *TechTreeHandler) GetTechTree(c echo.Context) error {
	ctx, requestID := getRequestContext(c)

	treeID := c.Param("treeId")
	// TODO: Get real user ID from auth context
	// For dogfooding MVP, use hardcoded user ID
	userID := "00000000-0000-0000-0000-000000000001" // Hardcoded for MVP

	logger.Logger.Info("fetching tech tree",
		slog.String("request_id", requestID),
		slog.String("tree_id", treeID),
		slog.String("user_id", userID),
	)

	// Fetch tech tree
	tree, err := h.techTreeRepo.GetByID(ctx, treeID)
	if err != nil {
		logger.Logger.Error("failed to fetch tech tree",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "Tech tree not found",
		})
	}

	// Fetch all nodes for this tree
	nodes, err := h.nodeRepo.GetByTechTree(ctx, treeID)
	if err != nil {
		logger.Logger.Error("failed to fetch nodes",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to fetch nodes",
		})
	}

	// Fetch user progress for all nodes
	userProgress, err := h.userProgressRepo.GetByUser(ctx, userID)
	if err != nil {
		logger.Logger.Warn("failed to fetch user progress, continuing without it",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
	}

	// Create progress map for efficient lookup
	progressMap := make(map[string]interface{})
	for _, progress := range userProgress {
		progressMap[progress.NodeID.String()] = progress
	}

	// Combine nodes with progress
	// TODO: Implement prerequisite checking to determine IsLocked status
	response := map[string]interface{}{
		"success":     true,
		"tech_tree":   tree,
		"nodes":       nodes,
		"progress":    progressMap,
	}

	logger.Logger.Info("tech tree fetched successfully",
		slog.String("request_id", requestID),
		slog.Int("node_count", len(nodes)),
	)

	return c.JSON(http.StatusOK, response)
}

// ============================================
// Node Handler Methods
// ============================================

// GetNode fetches a node with its exercises and user progress
func (h *NodeHandler) GetNode(c echo.Context) error {
	ctx, requestID := getRequestContext(c)

	nodeID := c.Param("nodeId")
	// TODO: Get real user ID from auth context
	userID := "00000000-0000-0000-0000-000000000001" // Hardcoded for MVP

	logger.Logger.Info("fetching node",
		slog.String("request_id", requestID),
		slog.String("node_id", nodeID),
		slog.String("user_id", userID),
	)

	// Fetch node
	node, err := h.nodeRepo.GetByID(ctx, nodeID)
	if err != nil {
		logger.Logger.Error("failed to fetch node",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "Node not found",
		})
	}

	// Fetch exercises for this node
	exercises, err := h.exerciseRepo.GetByNode(ctx, nodeID)
	if err != nil {
		logger.Logger.Error("failed to fetch exercises",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to fetch exercises",
		})
	}

	// Fetch user progress for this node
	progress, err := h.userProgressRepo.GetByUserAndNode(ctx, userID, nodeID)
	if err != nil {
		logger.Logger.Warn("no progress found for user and node",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		progress = nil
	}

	response := map[string]interface{}{
		"success":   true,
		"node":      node,
		"exercises": exercises,
		"progress":  progress,
	}

	logger.Logger.Info("node fetched successfully",
		slog.String("request_id", requestID),
		slog.Int("exercise_count", len(exercises)),
	)

	return c.JSON(http.StatusOK, response)
}

// ============================================
// Exercise Handler Methods (New)
// ============================================

// GetExercise fetches an exercise with submission history
func (h *ExerciseHandler) GetExercise(c echo.Context) error {
	ctx, requestID := getRequestContext(c)

	exerciseID := c.Param("exerciseId")
	// TODO: Get real user ID from auth context
	userID := "00000000-0000-0000-0000-000000000001" // Hardcoded for MVP

	logger.Logger.Info("fetching exercise",
		slog.String("request_id", requestID),
		slog.String("exercise_id", exerciseID),
		slog.String("user_id", userID),
	)

	// Fetch exercise
	exercise, err := h.exerciseRepo.GetByID(ctx, exerciseID)
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

	// Fetch submission history (limited to recent submissions)
	submissions, _, err := h.submissionRepo.GetByExercisePaginated(ctx, userID, exerciseID, 10, 0)
	if err != nil {
		logger.Logger.Warn("failed to fetch submission history",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		submissions = nil
	}

	// Calculate best score
	var bestScore *int
	if submissions != nil && len(submissions) > 0 {
		for _, sub := range submissions {
			if sub.TotalScore != nil {
				if bestScore == nil || *sub.TotalScore > *bestScore {
					bestScore = sub.TotalScore
				}
			}
		}
	}

	response := map[string]interface{}{
		"success":            true,
		"exercise":           exercise,
		"submission_history": submissions,
		"best_score":         bestScore,
	}

	logger.Logger.Info("exercise fetched successfully",
		slog.String("request_id", requestID),
		slog.Int("submission_count", len(submissions)),
	)

	return c.JSON(http.StatusOK, response)
}

// ============================================
// Submission Handler Methods (New)
// ============================================

// SubmitToExercise submits code to a specific exercise with full evaluation and scoring
func (h *SubmissionHandler) SubmitToExercise(c echo.Context) error {
	_, requestID := getRequestContext(c)

	exerciseID := c.Param("exerciseId")
	// TODO: Get real user ID from auth context
	userID := "00000000-0000-0000-0000-000000000001" // Hardcoded for MVP

	logger.Logger.Info("submit to exercise endpoint called",
		slog.String("request_id", requestID),
		slog.String("exercise_id", exerciseID),
		slog.String("user_id", userID),
	)

	var req struct {
		Code     string `json:"code"`
		Language string `json:"language"`
	}

	if err := c.Bind(&req); err != nil {
		logger.Logger.Error("failed to bind request body",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	if req.Code == "" {
		logger.Logger.Warn("missing required fields",
			slog.String("request_id", requestID),
		)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Code is required",
		})
	}

	// TODO: Implement full exercise submission with:
	// 1. Fetch exercise and run test cases
	// 2. Execute code and compare results
	// 3. Get AI feedback based on exercise rubric
	// 4. Calculate scores (60% tests + 40% AI)
	// 5. Save submission to database
	// 6. Update user progress and mastery tracking

	return c.JSON(http.StatusNotImplemented, map[string]interface{}{
		"success": false,
		"error":   "Submission with full evaluation not yet implemented",
		"message": "This endpoint will support test case execution, AI evaluation, and mastery tracking",
	})
}
