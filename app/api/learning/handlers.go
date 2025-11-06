// Package learning provides HTTP handlers for tech tree navigation and progress tracking endpoints
package learning

import (
	"context"
	"log/slog"
	"net/http"

	domainLearning "github.com/anttihil/code-auto-grader/domain/learning"
	"github.com/anttihil/code-auto-grader/internal/logger"
	"github.com/labstack/echo/v4"
)

// TechTreeHandler handles HTTP requests for tech tree endpoints
type TechTreeHandler struct {
	learningService *domainLearning.LearningService
}

// NewTechTreeHandler creates a new tech tree handler
func NewTechTreeHandler(learningService *domainLearning.LearningService) *TechTreeHandler {
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
		if _, ok := err.(*domainLearning.TechTreeNotFoundError); ok {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"success": false,
				"error":   "Tech tree not found",
			})
		}
		logger.Logger.Error("failed to fetch tech tree",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to fetch tech tree",
		})
	}

	// Build nodes with embedded progress
	nodesWithProgress := make([]map[string]interface{}, 0, len(view.Nodes))
	for _, nodeWithProg := range view.Nodes {
		nodeData := map[string]interface{}{
			"id":                  nodeWithProg.Node.ID,
			"tech_tree_id":        nodeWithProg.Node.TechTreeID,
			"name":                nodeWithProg.Node.Name,
			"slug":                nodeWithProg.Node.Slug,
			"description":         nodeWithProg.Node.Description,
			"learning_objectives": nodeWithProg.Node.LearningObjectives,
			"order_index":         nodeWithProg.Node.OrderIndex,
			"estimated_minutes":   nodeWithProg.Node.EstimatedMinutes,
			"created_at":          nodeWithProg.Node.CreatedAt,
			"updated_at":          nodeWithProg.Node.UpdatedAt,
			"progress": map[string]interface{}{
				"node_id":             nodeWithProg.Progress.NodeID,
				"exercises_completed": nodeWithProg.Progress.ExercisesCompleted,
				"total_exercises":     nodeWithProg.Progress.TotalExercises,
				"mastery_achieved":    nodeWithProg.Progress.MasteryAchieved,
				"is_unlocked":         nodeWithProg.Progress.IsUnlocked,
				"best_scores":         nodeWithProg.Progress.BestScores,
			},
		}
		nodesWithProgress = append(nodesWithProgress, nodeData)
	}

	// Flatten tech tree fields to response root level with nodes embedded
	response := map[string]interface{}{
		"success":     true,
		"id":          view.TechTree.ID,
		"name":        view.TechTree.Name,
		"slug":        view.TechTree.Slug,
		"description": view.TechTree.Description,
		"language":    view.TechTree.Language,
		"created_at":  view.TechTree.CreatedAt,
		"nodes":       nodesWithProgress,
	}

	return c.JSON(http.StatusOK, response)
}

// NodeHandler handles HTTP requests for node endpoints
type NodeHandler struct {
	learningService *domainLearning.LearningService
}

// NewNodeHandler creates a new node handler
func NewNodeHandler(learningService *domainLearning.LearningService) *NodeHandler {
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
		if _, ok := err.(*domainLearning.NodeNotFoundError); ok {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"success": false,
				"error":   "Node not found",
			})
		}
		logger.Logger.Error("failed to fetch node",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to fetch node",
		})
	}

	// Build progress data
	progressData := map[string]interface{}{
		"node_id":             view.Progress.NodeID,
		"exercises_completed": view.Progress.ExercisesCompleted,
		"total_exercises":     view.Progress.TotalExercises,
		"mastery_achieved":    view.Progress.MasteryAchieved,
		"is_unlocked":         view.Progress.IsUnlocked,
		"best_scores":         view.Progress.BestScores,
	}

	// Flatten node fields to response root level with embedded progress
	response := map[string]interface{}{
		"success":             true,
		"id":                  view.Node.ID,
		"tech_tree_id":        view.Node.TechTreeID,
		"name":                view.Node.Name,
		"slug":                view.Node.Slug,
		"description":         view.Node.Description,
		"learning_objectives": view.Node.LearningObjectives,
		"order_index":         view.Node.OrderIndex,
		"estimated_minutes":   view.Node.EstimatedMinutes,
		"created_at":          view.Node.CreatedAt,
		"updated_at":          view.Node.UpdatedAt,
		"progress":            progressData,
		"exercises":           view.Exercises,
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
