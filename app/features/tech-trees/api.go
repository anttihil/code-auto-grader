// Package learning provides HTTP handlers for tech tree navigation and progress tracking endpoints
package techtrees

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/anttihil/code-auto-grader/internal/logger"
	"github.com/labstack/echo/v4"
)

// TechTreeHandler handles HTTP requests for tech tree endpoints
type TechTreeHandler struct {
	learningService *LearningService
}

// NewTechTreeHandler creates a new tech tree handler
func NewTechTreeHandler(learningService *LearningService) *TechTreeHandler {
	if learningService == nil {
		panic("learning service is required")
	}
	return &TechTreeHandler{
		learningService: learningService,
	}
}

// GetTechTree handles GET /api/tech-tree/:treeId
// Fetches a tech tree with all nodes and user progress
func (h *TechTreeHandler) GetTechTree(c echo.Context) error {
	ctx, requestID := getRequestContext(c)

	treeID := c.Param("treeId")
	// TODO: Get real user ID from auth context
	userID := "00000000-0000-0000-0000-000000000001" // Hardcoded for MVP

	view, err := h.learningService.GetTechTree(ctx, treeID, userID, requestID)
	if err != nil {
		if _, ok := err.(*TechTreeNotFoundError); ok {
			return c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Tech tree not found",
			})
		}
		logger.Logger.Error("failed to fetch tech tree",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to fetch tech tree",
			Message: err.Error(),
		})
	}

	// Convert nodes with progress to DTO
	nodes := make([]NodeWithProgressDTO, len(view.Nodes))
	for i, nodeWithProg := range view.Nodes {
		// Convert best scores from map[string]int to map[string]float64
		bestScores := make(map[string]float64)
		for k, v := range nodeWithProg.Progress.BestScores {
			bestScores[k] = float64(v)
		}

		estimatedMins := 0
		if nodeWithProg.Node.EstimatedMinutes != nil {
			estimatedMins = *nodeWithProg.Node.EstimatedMinutes
		}

		nodes[i] = NodeWithProgressDTO{
			ID:                 nodeWithProg.Node.ID.String(),
			TechTreeID:         nodeWithProg.Node.TechTreeID.String(),
			Name:               nodeWithProg.Node.Name,
			Slug:               nodeWithProg.Node.Slug,
			Description:        getStringValue(nodeWithProg.Node.Description),
			LearningObjectives: nodeWithProg.Node.LearningObjectives,
			OrderIndex:         nodeWithProg.Node.OrderIndex,
			EstimatedMinutes:   estimatedMins,
			CreatedAt:          nodeWithProg.Node.CreatedAt,
			UpdatedAt:          nodeWithProg.Node.UpdatedAt,
			Progress: ProgressData{
				NodeID:             nodeWithProg.Progress.NodeID,
				ExercisesCompleted: nodeWithProg.Progress.ExercisesCompleted,
				TotalExercises:     nodeWithProg.Progress.TotalExercises,
				MasteryAchieved:    nodeWithProg.Progress.MasteryAchieved,
				IsUnlocked:         nodeWithProg.Progress.IsUnlocked,
				BestScores:         bestScores,
			},
		}
	}

	response := TechTreeResponse{
		ID:          view.TechTree.ID.String(),
		Name:        view.TechTree.Name,
		Slug:        view.TechTree.Slug,
		Description: getStringValue(view.TechTree.Description),
		Language:    view.TechTree.Language,
		CreatedAt:   view.TechTree.CreatedAt,
		Nodes:       nodes,
	}

	return c.JSON(http.StatusOK, response)
}

// NodeHandler handles HTTP requests for node endpoints
type NodeHandler struct {
	learningService *LearningService
}

// NewNodeHandler creates a new node handler
func NewNodeHandler(learningService *LearningService) *NodeHandler {
	if learningService == nil {
		panic("learning service is required")
	}
	return &NodeHandler{
		learningService: learningService,
	}
}

// GetNode handles GET /api/node/:nodeId
// Fetches a node with its exercises and user progress
func (h *NodeHandler) GetNode(c echo.Context) error {
	ctx, requestID := getRequestContext(c)

	nodeID := c.Param("nodeId")
	// TODO: Get real user ID from auth context
	userID := "00000000-0000-0000-0000-000000000001" // Hardcoded for MVP

	view, err := h.learningService.GetNode(ctx, nodeID, userID, requestID)
	if err != nil {
		if _, ok := err.(*NodeNotFoundError); ok {
			return c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Node not found",
			})
		}
		logger.Logger.Error("failed to fetch node",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to fetch node",
			Message: err.Error(),
		})
	}

	// Convert best scores from map[string]int to map[string]float64
	bestScores := make(map[string]float64)
	for k, v := range view.Progress.BestScores {
		bestScores[k] = float64(v)
	}

	progressData := ProgressData{
		NodeID:             view.Progress.NodeID,
		ExercisesCompleted: view.Progress.ExercisesCompleted,
		TotalExercises:     view.Progress.TotalExercises,
		MasteryAchieved:    view.Progress.MasteryAchieved,
		IsUnlocked:         view.Progress.IsUnlocked,
		BestScores:         bestScores,
	}

	// Convert exercises to DTO
	exercises := make([]ExerciseSummary, len(view.Exercises))
	for i, ex := range view.Exercises {
		estimatedMins := 0
		if ex.EstimatedMinutes != nil {
			estimatedMins = *ex.EstimatedMinutes
		}

		exercises[i] = ExerciseSummary{
			ID:               ex.ID.String(),
			NodeID:           ex.NodeID.String(),
			Title:            ex.Title,
			Slug:             ex.Slug,
			Description:      ex.Description,
			Difficulty:       ex.Difficulty,
			EstimatedMinutes: estimatedMins,
			OrderIndex:       ex.OrderIndex,
			Content:          ex.Content,
			CreatedAt:        ex.CreatedAt,
			UpdatedAt:        ex.UpdatedAt,
		}
	}

	estimatedMins := 0
	if view.Node.EstimatedMinutes != nil {
		estimatedMins = *view.Node.EstimatedMinutes
	}

	response := NodeResponse{
		ID:                 view.Node.ID.String(),
		TechTreeID:         view.Node.TechTreeID.String(),
		Name:               view.Node.Name,
		Slug:               view.Node.Slug,
		Description:        getStringValue(view.Node.Description),
		LearningObjectives: view.Node.LearningObjectives,
		OrderIndex:         view.Node.OrderIndex,
		EstimatedMinutes:   estimatedMins,
		CreatedAt:          view.Node.CreatedAt,
		UpdatedAt:          view.Node.UpdatedAt,
		Progress:           progressData,
		Exercises:          exercises,
	}

	return c.JSON(http.StatusOK, response)
}

// Helper functions

func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
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
