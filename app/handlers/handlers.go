package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/anttihil/code-auto-grader/internal/logger"
	"github.com/anttihil/code-auto-grader/models"
	"github.com/anttihil/code-auto-grader/repository"
	"github.com/anttihil/code-auto-grader/services"
	"github.com/anttihil/code-auto-grader/services/executor"
	"github.com/google/uuid"
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
	exerciseRepo     repository.ExerciseRepository
	userProgressRepo repository.UserProgressRepository
}

// NewTechTreeHandler creates a new tech tree handler with dependencies
func NewTechTreeHandler(
	techTreeRepo repository.TechTreeRepository,
	nodeRepo repository.NodeRepository,
	exerciseRepo repository.ExerciseRepository,
	userProgressRepo repository.UserProgressRepository,
) *TechTreeHandler {
	if techTreeRepo == nil {
		panic("tech tree repository is required")
	}
	if nodeRepo == nil {
		panic("node repository is required")
	}
	if exerciseRepo == nil {
		panic("exercise repository is required")
	}
	if userProgressRepo == nil {
		panic("user progress repository is required")
	}
	return &TechTreeHandler{
		techTreeRepo:     techTreeRepo,
		nodeRepo:         nodeRepo,
		exerciseRepo:     exerciseRepo,
		userProgressRepo: userProgressRepo,
	}
}

// NodeHandler handles node-related endpoints
type NodeHandler struct {
	nodeRepo         repository.NodeRepository
	exerciseRepo     repository.ExerciseRepository
	submissionRepo   repository.SubmissionRepository
	userProgressRepo repository.UserProgressRepository
}

// NewNodeHandler creates a new node handler with dependencies
func NewNodeHandler(
	nodeRepo repository.NodeRepository,
	exerciseRepo repository.ExerciseRepository,
	submissionRepo repository.SubmissionRepository,
	userProgressRepo repository.UserProgressRepository,
) *NodeHandler {
	if nodeRepo == nil {
		panic("node repository is required")
	}
	if exerciseRepo == nil {
		panic("exercise repository is required")
	}
	if submissionRepo == nil {
		panic("submission repository is required")
	}
	if userProgressRepo == nil {
		panic("user progress repository is required")
	}
	return &NodeHandler{
		nodeRepo:         nodeRepo,
		exerciseRepo:     exerciseRepo,
		submissionRepo:   submissionRepo,
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
	sampleExercise := map[string]any{
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

	return c.JSON(http.StatusOK, map[string]any{
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
	progressMap := make(map[string]*models.UserProgress)
	for i := range userProgress {
		progressMap[userProgress[i].NodeID.String()] = userProgress[i]
	}

	// Fetch exercise counts for each node
	exerciseCounts := make(map[string]int)
	for _, node := range nodes {
		exercises, err := h.exerciseRepo.GetByNode(ctx, node.ID.String())
		if err != nil {
			logger.Logger.Warn("failed to count exercises for node",
				slog.String("request_id", requestID),
				slog.String("node_id", node.ID.String()),
				slog.String("error", err.Error()),
			)
			exerciseCounts[node.ID.String()] = 0
		} else {
			exerciseCounts[node.ID.String()] = len(exercises)
		}
	}

	// Combine nodes with embedded progress
	nodesWithProgress := make([]map[string]interface{}, 0, len(nodes))
	for i, node := range nodes {
		nodeData := map[string]interface{}{
			"id":                  node.ID,
			"tech_tree_id":        node.TechTreeID,
			"name":                node.Name,
			"slug":                node.Slug,
			"description":         node.Description,
			"learning_objectives": node.LearningObjectives,
			"order_index":         node.OrderIndex,
			"estimated_minutes":   node.EstimatedMinutes,
			"created_at":          node.CreatedAt,
			"updated_at":          node.UpdatedAt,
		}

		// Add progress if it exists, otherwise create default progress
		if progress, exists := progressMap[node.ID.String()]; exists {
			// Update exercises_total if not set
			if progress.ExercisesTotal == 0 {
				progress.ExercisesTotal = exerciseCounts[node.ID.String()]
			}

			// Build progress object
			nodeData["progress"] = map[string]interface{}{
				"node_id":             progress.NodeID,
				"exercises_completed": progress.ExercisesCompleted,
				"total_exercises":     progress.ExercisesTotal,
				"mastery_achieved":    progress.IsMastered,
				"is_unlocked":         true, // If progress exists, node is unlocked
				"best_scores":         make(map[string]int), // TODO: Fetch from submissions
			}
		} else {
			// Create default progress for nodes without progress
			exerciseCount := exerciseCounts[node.ID.String()]
			nodeData["progress"] = map[string]interface{}{
				"node_id":             node.ID,
				"exercises_completed": 0,
				"total_exercises":     exerciseCount,
				"mastery_achieved":    false,
				"is_unlocked":         i == 0, // First node is always unlocked
				"best_scores":         make(map[string]int),
			}
		}

		nodesWithProgress = append(nodesWithProgress, nodeData)
	}

	// Flatten tech tree fields to response root level with nodes embedded
	response := map[string]interface{}{
		"success":     true,
		"id":          tree.ID,
		"name":        tree.Name,
		"slug":        tree.Slug,
		"description": tree.Description,
		"language":    tree.Language,
		"created_at":  tree.CreatedAt,
		"nodes":       nodesWithProgress,
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

	// Fetch best scores for each exercise from submissions
	bestScores := make(map[string]int)
	for _, exercise := range exercises {
		submissions, _, err := h.submissionRepo.GetByExercisePaginated(ctx, userID, exercise.ID.String(), 100, 0)
		if err == nil && len(submissions) > 0 {
			// Find best score
			for _, sub := range submissions {
				if sub.TotalScore != nil {
					if bestScore, exists := bestScores[exercise.ID.String()]; !exists || *sub.TotalScore > bestScore {
						bestScores[exercise.ID.String()] = *sub.TotalScore
					}
				}
			}
		}
	}

	// Build progress object with best scores
	var progressData map[string]interface{}
	if progress != nil {
		// Update exercises_total if not set
		if progress.ExercisesTotal == 0 {
			progress.ExercisesTotal = len(exercises)
		}

		progressData = map[string]interface{}{
			"node_id":             progress.NodeID,
			"exercises_completed": progress.ExercisesCompleted,
			"total_exercises":     progress.ExercisesTotal,
			"mastery_achieved":    progress.IsMastered,
			"is_unlocked":         true, // If progress exists, node is unlocked
			"best_scores":         bestScores,
		}
	} else {
		// Create default progress
		progressData = map[string]interface{}{
			"node_id":             nodeID,
			"exercises_completed": 0,
			"total_exercises":     len(exercises),
			"mastery_achieved":    false,
			"is_unlocked":         true, // Assume unlocked if viewing it
			"best_scores":         bestScores,
		}
	}

	// Flatten node fields to response root level with embedded progress
	response := map[string]interface{}{
		"success":             true,
		"id":                  node.ID,
		"tech_tree_id":        node.TechTreeID,
		"name":                node.Name,
		"slug":                node.Slug,
		"description":         node.Description,
		"learning_objectives": node.LearningObjectives,
		"order_index":         node.OrderIndex,
		"estimated_minutes":   node.EstimatedMinutes,
		"created_at":          node.CreatedAt,
		"updated_at":          node.UpdatedAt,
		"progress":            progressData,
		"exercises":           exercises,
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

	// Calculate best score and attempt count
	var bestScore int
	attemptCount := len(submissions)
	if len(submissions) > 0 {
		for _, sub := range submissions {
			if sub.TotalScore != nil {
				if *sub.TotalScore > bestScore {
					bestScore = *sub.TotalScore
				}
			}
		}
	}

	// Flatten exercise fields to response root level with submission data
	response := map[string]interface{}{
		"success":      true,
		"exercise":     exercise,
		"submissions":  submissions,  // Frontend expects "submissions" not "submission_history"
		"best_score":   bestScore,
		"attempt_count": attemptCount,
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
	ctx, requestID := getRequestContext(c)

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

	// Set default language if not provided
	if req.Language == "" {
		req.Language = "javascript"
	}

	// 1. Fetch exercise from database
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

	logger.Logger.Info("exercise fetched for submission",
		slog.String("request_id", requestID),
		slog.String("exercise_title", exercise.Title),
		slog.Int("test_case_count", len(exercise.Content.TestCases)),
	)

	// 2. Run test cases
	testResults, executionOutput, executionError, executionTimeMs := h.runTestCases(ctx, req.Code, req.Language, exercise.Content.TestCases, requestID)

	// Calculate test score (0-100)
	testScore := 0
	if testResults.TotalTests > 0 {
		testScore = (testResults.PassedTests * 100) / testResults.TotalTests
	}

	logger.Logger.Info("test cases executed",
		slog.String("request_id", requestID),
		slog.Int("passed", testResults.PassedTests),
		slog.Int("total", testResults.TotalTests),
		slog.Int("test_score", testScore),
	)

	// 3. Get AI feedback using exercise's AI evaluation prompt
	startAI := time.Now()
	aiFeedback, aiScore, err := h.evaluateWithAI(ctx, exercise, req.Code, req.Language, testResults, requestID)
	if err != nil {
		logger.Logger.Warn("AI evaluation failed, continuing without it",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		// Continue without AI feedback rather than failing the whole submission
		aiFeedback = &models.AIFeedback{
			Passed:       false,
			Score:        0,
			Feedback:     "AI evaluation temporarily unavailable",
			Strengths:    []string{},
			Improvements: []string{},
			Correctness:  "",
			CodeQuality:  "",
			Suggestions:  "",
		}
		aiScore = 0
	}
	aiEvaluationTimeMs := int(time.Since(startAI).Milliseconds())

	logger.Logger.Info("AI evaluation completed",
		slog.String("request_id", requestID),
		slog.Int("ai_score", aiScore),
		slog.Int("duration_ms", aiEvaluationTimeMs),
	)

	// 4. Calculate total score using grading rubric weights
	testWeight := exercise.Content.GradingRubric.TestWeight
	aiWeight := exercise.Content.GradingRubric.AIWeight

	// Normalize weights if they don't sum to 100
	totalWeight := testWeight + aiWeight
	if totalWeight == 0 {
		totalWeight = 100
		testWeight = 60
		aiWeight = 40
	}

	totalScore := (testScore*testWeight + aiScore*aiWeight) / totalWeight

	// Scale AI score to match the weight (0-40) for database storage
	// AI returns 0-100, but database expects it scaled to the AI weight
	aiScoreScaled := (aiScore * aiWeight) / 100
	testScoreScaled := (testScore * testWeight) / 100

	logger.Logger.Info("final score calculated",
		slog.String("request_id", requestID),
		slog.Int("test_score", testScore),
		slog.Int("ai_score", aiScore),
		slog.Int("test_score_scaled", testScoreScaled),
		slog.Int("ai_score_scaled", aiScoreScaled),
		slog.Int("total_score", totalScore),
		slog.Int("test_weight", testWeight),
		slog.Int("ai_weight", aiWeight),
	)

	// 5. Save submission to database
	submission, err := h.saveSubmission(ctx, userID, exerciseID, req.Code, req.Language,
		&executionOutput, executionError, &executionTimeMs,
		testResults, aiFeedback, &aiScoreScaled, &aiEvaluationTimeMs,
		&testScoreScaled, &totalScore, requestID)
	if err != nil {
		logger.Logger.Error("failed to save submission",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to save submission",
		})
	}

	logger.Logger.Info("submission saved successfully",
		slog.String("request_id", requestID),
		slog.String("submission_id", submission.ID.String()),
		slog.Int("attempt_number", submission.AttemptNumber),
	)

	// 6. Update user progress and mastery tracking
	if err := h.updateUserProgress(ctx, userID, exercise.NodeID.String(), totalScore, requestID); err != nil {
		logger.Logger.Warn("failed to update user progress",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		// Don't fail the submission if progress update fails
	}

	// Return submission result
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":          true,
		"submission_id":    submission.ID,
		"execution_output": executionOutput,
		"execution_error":  executionError,
		"test_results":     testResults,
		"ai_feedback":      aiFeedback,
		"test_score":       testScore,
		"ai_score":         aiScore,
		"total_score":      totalScore,
		"attempt_number":   submission.AttemptNumber,
	})
}

// ============================================
// Helper Methods for Submission Processing
// ============================================

// runTestCases executes all test cases and returns results
func (h *SubmissionHandler) runTestCases(ctx context.Context, submittedCode, language string, testCases []models.TestCase, requestID string) (*models.TestResults, string, *string, int) {
	// For Rust, run all tests together in one compilation
	if strings.ToLower(language) == "rust" {
		return h.runRustTests(ctx, submittedCode, testCases, requestID)
	}

	// For other languages, run tests individually (existing behavior)
	results := &models.TestResults{
		TotalTests:  len(testCases),
		PassedTests: 0,
		FailedTests: 0,
		Results:     make([]models.TestResult, 0, len(testCases)),
	}

	var combinedOutput strings.Builder
	var totalExecutionTime int
	var lastError *string

	for _, testCase := range testCases {
		// Combine submitted code with test code
		codeWithTest := submittedCode + "\n\n" + testCase.Code

		logger.Logger.Debug("running test case",
			slog.String("request_id", requestID),
			slog.String("test_name", testCase.Name),
		)

		// Execute the combined code
		startTime := time.Now()
		execResult := h.codeExecutor.ExecuteCodeWithContext(ctx, codeWithTest, language)
		duration := int(time.Since(startTime).Milliseconds())
		totalExecutionTime += duration

		testResult := models.TestResult{
			Name:       testCase.Name,
			DurationMs: &duration,
		}

		if execResult.Success {
			testResult.Status = "passed"
			results.PassedTests++
			combinedOutput.WriteString(fmt.Sprintf("✓ %s: PASSED\n", testCase.Name))
			if execResult.Output != "" {
				combinedOutput.WriteString(fmt.Sprintf("  Output: %s\n", execResult.Output))
			}
		} else {
			testResult.Status = "failed"
			results.FailedTests++
			errorMsg := execResult.Error
			testResult.Error = &errorMsg
			lastError = &errorMsg
			combinedOutput.WriteString(fmt.Sprintf("✗ %s: FAILED\n", testCase.Name))
			combinedOutput.WriteString(fmt.Sprintf("  Error: %s\n", errorMsg))
		}

		results.Results = append(results.Results, testResult)
	}

	output := combinedOutput.String()
	return results, output, lastError, totalExecutionTime
}

// runRustTests runs all Rust tests in a single compilation
func (h *SubmissionHandler) runRustTests(ctx context.Context, submittedCode string, testCases []models.TestCase, requestID string) (*models.TestResults, string, *string, int) {
	results := &models.TestResults{
		TotalTests:  len(testCases),
		PassedTests: 0,
		FailedTests: 0,
		Results:     make([]models.TestResult, 0, len(testCases)),
	}

	// Build complete Rust test file
	var codeBuilder strings.Builder
	codeBuilder.WriteString(submittedCode)
	codeBuilder.WriteString("\n\n#[cfg(test)]\nmod tests {\n    use super::*;\n\n")

	// Add all test cases
	for _, testCase := range testCases {
		// Extract test function from testCase.Code (remove #[test] and outer wrapping if present)
		testCode := testCase.Code
		// Simple approach: just append the test code
		codeBuilder.WriteString("    ")
		codeBuilder.WriteString(strings.ReplaceAll(testCode, "\n", "\n    "))
		codeBuilder.WriteString("\n\n")
	}

	codeBuilder.WriteString("}\n")
	completeCode := codeBuilder.String()

	logger.Logger.Debug("running Rust tests",
		slog.String("request_id", requestID),
		slog.Int("test_count", len(testCases)),
	)

	// Execute all tests together
	startTime := time.Now()
	execResult := h.codeExecutor.ExecuteCodeWithContext(ctx, completeCode, "rust")
	totalExecutionTime := int(time.Since(startTime).Milliseconds())

	// Parse test output
	if !execResult.Success {
		// Compilation error or runtime panic
		errorMsg := execResult.Error
		return results, execResult.Error, &errorMsg, totalExecutionTime
	}

	// Parse Rust test output
	output := execResult.Output
	testResults := parseRustTestOutput(output, testCases)

	// Build results
	var combinedOutput strings.Builder
	combinedOutput.WriteString(output)
	combinedOutput.WriteString("\n\n=== Individual Test Results ===\n")

	var lastError *string
	for _, tr := range testResults {
		results.Results = append(results.Results, tr)
		if tr.Status == "passed" {
			results.PassedTests++
			combinedOutput.WriteString(fmt.Sprintf("✓ %s: PASSED\n", tr.Name))
		} else {
			results.FailedTests++
			combinedOutput.WriteString(fmt.Sprintf("✗ %s: FAILED\n", tr.Name))
			if tr.Error != nil {
				combinedOutput.WriteString(fmt.Sprintf("  Error: %s\n", *tr.Error))
				lastError = tr.Error
			}
		}
	}

	return results, combinedOutput.String(), lastError, totalExecutionTime
}

// parseRustTestOutput parses Rust test runner output to extract individual test results
func parseRustTestOutput(output string, testCases []models.TestCase) []models.TestResult {
	results := make([]models.TestResult, 0, len(testCases))

	// Parse lines like: "test tests::test_name ... ok" or "test tests::test_name ... FAILED"
	lines := strings.Split(output, "\n")
	testResultsMap := make(map[string]models.TestResult)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "test ") {
			continue
		}

		// Extract test name and status
		// Format: "test tests::test_name ... ok" or "test tests::test_name ... FAILED"
		parts := strings.Fields(line)
		if len(parts) < 4 {
			continue
		}

		testFullName := parts[1] // e.g., "tests::test_name"
		status := parts[len(parts)-1] // e.g., "ok" or "FAILED"

		// Extract short name (after "::")
		nameParts := strings.Split(testFullName, "::")
		shortName := nameParts[len(nameParts)-1]

		result := models.TestResult{
			Name: shortName,
		}

		switch status {
			case "ok":
				result.Status = "passed"
			case "FAILED":
				result.Status = "failed"
				// TODO: Extract error message from failure output
				errorMsg := "Test assertion failed"
				result.Error = &errorMsg
		}

		testResultsMap[shortName] = result
	}

	// Match test cases with results
	for _, testCase := range testCases {
		if result, found := testResultsMap[testCase.Name]; found {
			results = append(results, result)
		} else {
			// Test didn't run (compilation error?)
			errorMsg := "Test did not run"
			results = append(results, models.TestResult{
				Name:   testCase.Name,
				Status: "failed",
				Error:  &errorMsg,
			})
		}
	}

	return results
}

// evaluateWithAI gets AI feedback and score for the submission
func (h *SubmissionHandler) evaluateWithAI(ctx context.Context, exercise *models.Exercise, code, language string, testResults *models.TestResults, requestID string) (*models.AIFeedback, int, error) {
	// Build evaluation prompt using exercise's AI evaluation prompt
	prompt := exercise.Content.AIEvaluationPrompt

	// If no custom prompt, use a default one
	if prompt == "" {
		prompt = buildDefaultAIPrompt(exercise, code, language, testResults)
	} else {
		// Replace placeholders in the custom prompt
		prompt = replacePlaceholders(prompt, exercise, code, language, testResults)
	}

	logger.Logger.Debug("sending code to AI for evaluation",
		slog.String("request_id", requestID),
		slog.Int("prompt_length", len(prompt)),
	)

	// Get AI feedback from Gemini (returns services.Feedback struct)
	serviceFeedback, err := h.gemini.CheckSubmissionWithContext(ctx, prompt, code, language)
	if err != nil {
		return nil, 0, fmt.Errorf("AI evaluation failed: %w", err)
	}

	// Convert services.Feedback to models.AIFeedback
	feedback := &models.AIFeedback{
		Passed:       serviceFeedback.Passed,
		Score:        serviceFeedback.Score,
		Feedback:     serviceFeedback.Feedback,
		Strengths:    serviceFeedback.Strengths,
		Improvements: serviceFeedback.Improvements,
		Correctness:  serviceFeedback.Correctness,
		CodeQuality:  serviceFeedback.CodeQuality,
		Suggestions:  serviceFeedback.Suggestions,
	}

	// Use the score from the AI feedback, or derive one if not provided
	aiScore := feedback.Score
	if aiScore == 0 {
		// Fallback: derive score from feedback text if AI didn't provide one
		aiScore = deriveAIScore(feedback.Feedback, testResults)
		feedback.Score = aiScore
	}

	return feedback, aiScore, nil
}

// buildDefaultAIPrompt creates a default evaluation prompt
func buildDefaultAIPrompt(exercise *models.Exercise, code, language string, testResults *models.TestResults) string {
	return fmt.Sprintf(`Evaluate this code submission for the following exercise:

Exercise: %s
Description: %s
Difficulty: %s
Language: %s

Test Results: %d/%d tests passed

Submitted Code:
%s

Please provide feedback on:
1. Code correctness and logic
2. Code quality and style
3. Efficiency and best practices
4. Suggestions for improvement

Be constructive and educational in your feedback.`,
		exercise.Title,
		exercise.Description,
		exercise.Difficulty,
		language,
		testResults.PassedTests,
		testResults.TotalTests,
		code)
}

// replacePlaceholders replaces template variables in AI prompt
func replacePlaceholders(prompt string, exercise *models.Exercise, code, language string, testResults *models.TestResults) string {
	replacer := strings.NewReplacer(
		"{exercise_title}", exercise.Title,
		"{exercise_description}", exercise.Description,
		"{difficulty}", exercise.Difficulty,
		"{language}", language,
		"{tests_passed}", fmt.Sprintf("%d", testResults.PassedTests),
		"{tests_total}", fmt.Sprintf("%d", testResults.TotalTests),
		"{submitted_code}", code,
	)
	return replacer.Replace(prompt)
}

// deriveAIScore calculates a score based on AI feedback and test results
func deriveAIScore(feedback string, testResults *models.TestResults) int {
	// Simple heuristic: if all tests pass, base score is 80-100
	// If some tests fail, base score is 40-80
	// Adjust based on feedback sentiment (basic implementation)

	baseScore := 50
	if testResults.PassedTests == testResults.TotalTests {
		baseScore = 85
	} else if testResults.PassedTests > 0 {
		ratio := float64(testResults.PassedTests) / float64(testResults.TotalTests)
		baseScore = int(50 + ratio*35)
	}

	// Adjust based on feedback keywords (simple sentiment analysis)
	feedbackLower := strings.ToLower(feedback)
	if strings.Contains(feedbackLower, "excellent") || strings.Contains(feedbackLower, "perfect") {
		baseScore = min(100, baseScore+15)
	} else if strings.Contains(feedbackLower, "good") || strings.Contains(feedbackLower, "well done") {
		baseScore = min(100, baseScore+10)
	} else if strings.Contains(feedbackLower, "needs improvement") || strings.Contains(feedbackLower, "consider") {
		baseScore = max(0, baseScore-10)
	} else if strings.Contains(feedbackLower, "incorrect") || strings.Contains(feedbackLower, "wrong") {
		baseScore = max(0, baseScore-20)
	}

	return baseScore
}

// saveSubmission saves the submission to the database
func (h *SubmissionHandler) saveSubmission(
	ctx context.Context,
	userID, exerciseID, code, language string,
	executionOutput *string,
	executionError *string,
	executionTimeMs *int,
	testResults *models.TestResults,
	aiFeedback *models.AIFeedback,
	aiScore *int,
	aiEvaluationTimeMs *int,
	testScore *int,
	totalScore *int,
	requestID string,
) (*models.Submission, error) {
	// Get attempt number
	attemptCount, err := h.submissionRepo.CountAttempts(ctx, userID, exerciseID)
	if err != nil {
		return nil, fmt.Errorf("failed to count attempts: %w", err)
	}
	attemptNumber := attemptCount + 1

	// Create submission
	submission := &models.Submission{
		ID:                 uuid.New(),
		UserID:             uuid.MustParse(userID),
		ExerciseID:         uuid.MustParse(exerciseID),
		SubmittedCode:      code,
		Language:           language,
		ExecutionOutput:    executionOutput,
		ExecutionError:     executionError,
		ExecutionTimeMs:    executionTimeMs,
		TestResults:        testResults,
		AIFeedback:         aiFeedback,
		AIScore:            aiScore,
		AIEvaluationTimeMs: aiEvaluationTimeMs,
		TestScore:          testScore,
		TotalScore:         totalScore,
		AttemptNumber:      attemptNumber,
	}

	// Save to database
	if err := h.submissionRepo.Create(ctx, submission); err != nil {
		return nil, fmt.Errorf("failed to create submission: %w", err)
	}

	return submission, nil
}

// updateUserProgress updates user progress and mastery tracking
func (h *SubmissionHandler) updateUserProgress(ctx context.Context, userID, nodeID string, score int, requestID string) error {
	// For dogfooding MVP, implement basic mastery tracking
	// A score of 80+ counts as passing the exercise

	// Get current progress
	progress, err := h.userProgressRepo.GetByUserAndNode(ctx, userID, nodeID)
	if err != nil {
		return fmt.Errorf("failed to get user progress: %w", err)
	}

	// If no progress exists, create initial progress record
	if progress == nil {
		now := time.Now()
		progress = &models.UserProgress{
			ID:               uuid.New(),
			UserID:           uuid.MustParse(userID),
			NodeID:           uuid.MustParse(nodeID),
			Status:           "in_progress",
			ExercisesCompleted: 0,
			ExercisesTotal:   0, // TODO: Count exercises for this node
			MasteryThreshold: 80.0,
			IsMastered:       false,
			StartedAt:        &now,
			LastActivityAt:   &now,
		}
	}

	// Update last activity
	now := time.Now()
	progress.LastActivityAt = &now

	// If score >= 80, increment exercises completed
	if score >= 80 {
		progress.ExercisesCompleted++
	}

	// Calculate mastery score (average of completed exercises)
	// For MVP, we'll use a simplified approach
	// TODO: Implement proper mastery calculation based on all submissions
	masteryScore := float64(score)
	progress.MasteryScore = &masteryScore

	// Check if mastered (for MVP, just check if score >= threshold)
	if masteryScore >= progress.MasteryThreshold {
		progress.IsMastered = true
		progress.Status = "completed"
		if progress.CompletedAt == nil {
			progress.CompletedAt = &now
		}
	}

	// Save updated progress
	if err := h.userProgressRepo.Upsert(ctx, progress); err != nil {
		return fmt.Errorf("failed to update progress: %w", err)
	}

	logger.Logger.Info("user progress updated",
		slog.String("request_id", requestID),
		slog.String("user_id", userID),
		slog.String("node_id", nodeID),
		slog.Int("exercises_completed", progress.ExercisesCompleted),
		slog.Bool("is_mastered", progress.IsMastered),
	)

	return nil
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
