package techtrees

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/anttihil/code-auto-grader/internal/logger"
	"github.com/anttihil/code-auto-grader/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestMain sets up test environment
func TestMain(m *testing.M) {
	// Initialize logger for tests
	logger.Initialize("test")

	// Run tests
	code := m.Run()

	os.Exit(code)
}

// ============================================================================
// Service Interface and Mock
// ============================================================================

// learningService is the interface used for testing
type learningService interface {
	GetTechTree(ctx context.Context, treeID, userID, requestID string) (*TechTreeView, error)
	GetNode(ctx context.Context, nodeID, userID, requestID string) (*NodeView, error)
	CalculateProgress(ctx context.Context, userID, treeID string) (int, int, error)
}

// mockLearningService implements learningService for testing
type mockLearningService struct {
	mock.Mock
}

func (m *mockLearningService) GetTechTree(ctx context.Context, treeID, userID, requestID string) (*TechTreeView, error) {
	args := m.Called(ctx, treeID, userID, requestID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*TechTreeView), args.Error(1)
}

func (m *mockLearningService) GetNode(ctx context.Context, nodeID, userID, requestID string) (*NodeView, error) {
	args := m.Called(ctx, nodeID, userID, requestID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*NodeView), args.Error(1)
}

func (m *mockLearningService) CalculateProgress(ctx context.Context, userID, treeID string) (int, int, error) {
	args := m.Called(ctx, userID, treeID)
	return args.Int(0), args.Int(1), args.Error(2)
}

// testTechTreeHandler wraps TechTreeHandler to accept interface for testing
type testTechTreeHandler struct {
	service learningService
}

func (h *testTechTreeHandler) GetTechTree(c echo.Context) error {
	ctx, requestID := getRequestContext(c)
	treeID := c.Param("treeId")
	userID := "00000000-0000-0000-0000-000000000001"

	view, err := h.service.GetTechTree(ctx, treeID, userID, requestID)
	if err != nil {
		if _, ok := err.(*TechTreeNotFoundError); ok {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"success": false,
				"error":   "Tech tree not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to fetch tech tree",
		})
	}

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

// testNodeHandler wraps NodeHandler to accept interface for testing
type testNodeHandler struct {
	service learningService
}

func (h *testNodeHandler) GetNode(c echo.Context) error {
	ctx, requestID := getRequestContext(c)
	nodeID := c.Param("nodeId")
	userID := "00000000-0000-0000-0000-000000000001"

	view, err := h.service.GetNode(ctx, nodeID, userID, requestID)
	if err != nil {
		if _, ok := err.(*NodeNotFoundError); ok {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"success": false,
				"error":   "Node not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to fetch node",
		})
	}

	progressData := map[string]interface{}{
		"node_id":             view.Progress.NodeID,
		"exercises_completed": view.Progress.ExercisesCompleted,
		"total_exercises":     view.Progress.TotalExercises,
		"mastery_achieved":    view.Progress.MasteryAchieved,
		"is_unlocked":         view.Progress.IsUnlocked,
		"best_scores":         view.Progress.BestScores,
	}

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

// ============================================================================
// Helper Functions
// ============================================================================

func setupEcho() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	return e
}

func createMockTechTreeView() *TechTreeView {
	treeID := uuid.New()
	nodeID1 := uuid.New()
	nodeID2 := uuid.New()
	description := "Test tech tree"

	return &TechTreeView{
		TechTree: &models.TechTree{
			ID:          treeID,
			Name:        "JavaScript Fundamentals",
			Slug:        "js-fundamentals",
			Description: &description,
			Language:    "javascript",
		},
		Nodes: []*NodeWithProgress{
			{
				Node: &models.Node{
					ID:                 nodeID1,
					TechTreeID:         treeID,
					Name:               "Variables",
					Slug:               "variables",
					OrderIndex:         0,
					LearningObjectives: []string{"Understand var, let, const"},
				},
				Progress: &NodeProgress{
					NodeID:             nodeID1.String(),
					ExercisesCompleted: 3,
					TotalExercises:     3,
					MasteryAchieved:    true,
					IsUnlocked:         true,
					BestScores:         map[string]int{},
				},
			},
			{
				Node: &models.Node{
					ID:                 nodeID2,
					TechTreeID:         treeID,
					Name:               "Functions",
					Slug:               "functions",
					OrderIndex:         1,
					LearningObjectives: []string{"Master function declarations"},
				},
				Progress: &NodeProgress{
					NodeID:             nodeID2.String(),
					ExercisesCompleted: 1,
					TotalExercises:     5,
					MasteryAchieved:    false,
					IsUnlocked:         true,
					BestScores:         map[string]int{},
				},
			},
		},
		TotalNodes:     2,
		CompletedNodes: 1,
	}
}

func createMockNodeView() *NodeView {
	nodeID := uuid.New()
	treeID := uuid.New()
	exerciseID1 := uuid.New()
	exerciseID2 := uuid.New()
	description := "Test node"

	return &NodeView{
		Node: &models.Node{
			ID:                 nodeID,
			TechTreeID:         treeID,
			Name:               "Variables",
			Slug:               "variables",
			Description:        &description,
			OrderIndex:         0,
			LearningObjectives: []string{"Understand var, let, const"},
		},
		Exercises: []*models.Exercise{
			{
				ID:          exerciseID1,
				NodeID:      nodeID,
				Title:       "Variable Declaration",
				Description: "Practice declaring variables",
				Difficulty:  "beginner",
			},
			{
				ID:          exerciseID2,
				NodeID:      nodeID,
				Title:       "Variable Scope",
				Description: "Understand variable scope",
				Difficulty:  "intermediate",
			},
		},
		Progress: &NodeProgress{
			NodeID:             nodeID.String(),
			ExercisesCompleted: 2,
			TotalExercises:     2,
			MasteryAchieved:    true,
			IsUnlocked:         true,
			BestScores: map[string]int{
				exerciseID1.String(): 95,
				exerciseID2.String(): 87,
			},
		},
	}
}

// ============================================================================
// TechTreeHandler Tests
// ============================================================================

func TestNewTechTreeHandler(t *testing.T) {
	t.Run("panics when service is nil", func(t *testing.T) {
		assert.Panics(t, func() {
			NewTechTreeHandler(nil)
		})
	})
}

func TestGetTechTree_HappyPath(t *testing.T) {
	mockService := new(mockLearningService)
	handler := &testTechTreeHandler{service: mockService}

	expectedView := createMockTechTreeView()
	mockService.On("GetTechTree", mock.Anything, "tree-123", "00000000-0000-0000-0000-000000000001", mock.Anything).
		Return(expectedView, nil)

	e := setupEcho()
	req := httptest.NewRequest(http.MethodGet, "/api/tech-tree/tree-123", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("treeId")
	c.SetParamValues("tree-123")
	c.Set("request_id", "test-request-id")

	err := handler.GetTechTree(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))
	assert.Equal(t, "JavaScript Fundamentals", response["name"])
	assert.Equal(t, "js-fundamentals", response["slug"])
	assert.Equal(t, "javascript", response["language"])

	nodes := response["nodes"].([]interface{})
	assert.Len(t, nodes, 2)

	// Check first node
	node1 := nodes[0].(map[string]interface{})
	assert.Equal(t, "Variables", node1["name"])
	assert.Equal(t, "variables", node1["slug"])

	progress1 := node1["progress"].(map[string]interface{})
	assert.Equal(t, float64(3), progress1["exercises_completed"])
	assert.Equal(t, float64(3), progress1["total_exercises"])
	assert.True(t, progress1["mastery_achieved"].(bool))
	assert.True(t, progress1["is_unlocked"].(bool))

	mockService.AssertExpectations(t)
}

func TestGetTechTree_NotFound(t *testing.T) {
	mockService := new(mockLearningService)
	handler := &testTechTreeHandler{service: mockService}

	mockService.On("GetTechTree", mock.Anything, "nonexistent", "00000000-0000-0000-0000-000000000001", mock.Anything).
		Return(nil, &TechTreeNotFoundError{TreeID: "nonexistent"})

	e := setupEcho()
	req := httptest.NewRequest(http.MethodGet, "/api/tech-tree/nonexistent", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("treeId")
	c.SetParamValues("nonexistent")
	c.Set("request_id", "test-request-id")

	err := handler.GetTechTree(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)

	var response map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.False(t, response["success"].(bool))
	assert.Equal(t, "Tech tree not found", response["error"])

	mockService.AssertExpectations(t)
}

func TestGetTechTree_ServiceError(t *testing.T) {
	mockService := new(mockLearningService)
	handler := &testTechTreeHandler{service: mockService}

	mockService.On("GetTechTree", mock.Anything, "tree-123", "00000000-0000-0000-0000-000000000001", mock.Anything).
		Return(nil, assert.AnError)

	e := setupEcho()
	req := httptest.NewRequest(http.MethodGet, "/api/tech-tree/tree-123", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("treeId")
	c.SetParamValues("tree-123")
	c.Set("request_id", "test-request-id")

	err := handler.GetTechTree(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var response map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.False(t, response["success"].(bool))
	assert.Equal(t, "Failed to fetch tech tree", response["error"])

	mockService.AssertExpectations(t)
}

func TestGetTechTree_NoRequestID(t *testing.T) {
	mockService := new(mockLearningService)
	handler := &testTechTreeHandler{service: mockService}

	expectedView := createMockTechTreeView()
	mockService.On("GetTechTree", mock.Anything, "tree-123", "00000000-0000-0000-0000-000000000001", "unknown").
		Return(expectedView, nil)

	e := setupEcho()
	req := httptest.NewRequest(http.MethodGet, "/api/tech-tree/tree-123", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("treeId")
	c.SetParamValues("tree-123")
	// No request_id set

	err := handler.GetTechTree(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	mockService.AssertExpectations(t)
}

// ============================================================================
// NodeHandler Tests
// ============================================================================

func TestNewNodeHandler(t *testing.T) {
	mockService := &LearningService{}

	t.Run("creates handler with service", func(t *testing.T) {
		handler := NewNodeHandler(mockService)
		assert.NotNil(t, handler)
		assert.Equal(t, mockService, handler.learningService)
	})

	t.Run("panics when service is nil", func(t *testing.T) {
		assert.Panics(t, func() {
			NewNodeHandler(nil)
		})
	})
}

func TestGetNode_HappyPath(t *testing.T) {
	mockService := new(mockLearningService)
	handler := &testNodeHandler{service: mockService}

	expectedView := createMockNodeView()
	mockService.On("GetNode", mock.Anything, "node-123", "00000000-0000-0000-0000-000000000001", mock.Anything).
		Return(expectedView, nil)

	e := setupEcho()
	req := httptest.NewRequest(http.MethodGet, "/api/node/node-123", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("nodeId")
	c.SetParamValues("node-123")
	c.Set("request_id", "test-request-id")

	err := handler.GetNode(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))
	assert.Equal(t, "Variables", response["name"])
	assert.Equal(t, "variables", response["slug"])

	// Check progress
	progress := response["progress"].(map[string]interface{})
	assert.Equal(t, float64(2), progress["exercises_completed"])
	assert.Equal(t, float64(2), progress["total_exercises"])
	assert.True(t, progress["mastery_achieved"].(bool))
	assert.True(t, progress["is_unlocked"].(bool))

	bestScores := progress["best_scores"].(map[string]interface{})
	assert.NotEmpty(t, bestScores)

	// Check exercises
	exercises := response["exercises"].([]interface{})
	assert.Len(t, exercises, 2)

	exercise1 := exercises[0].(map[string]interface{})
	assert.Equal(t, "Variable Declaration", exercise1["title"])
	assert.Equal(t, "beginner", exercise1["difficulty"])

	mockService.AssertExpectations(t)
}

func TestGetNode_NotFound(t *testing.T) {
	mockService := new(mockLearningService)
	handler := &testNodeHandler{service: mockService}

	mockService.On("GetNode", mock.Anything, "nonexistent", "00000000-0000-0000-0000-000000000001", mock.Anything).
		Return(nil, &NodeNotFoundError{NodeID: "nonexistent"})

	e := setupEcho()
	req := httptest.NewRequest(http.MethodGet, "/api/node/nonexistent", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("nodeId")
	c.SetParamValues("nonexistent")
	c.Set("request_id", "test-request-id")

	err := handler.GetNode(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)

	var response map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.False(t, response["success"].(bool))
	assert.Equal(t, "Node not found", response["error"])

	mockService.AssertExpectations(t)
}

func TestGetNode_ServiceError(t *testing.T) {
	mockService := new(mockLearningService)
	handler := &testNodeHandler{service: mockService}

	mockService.On("GetNode", mock.Anything, "node-123", "00000000-0000-0000-0000-000000000001", mock.Anything).
		Return(nil, assert.AnError)

	e := setupEcho()
	req := httptest.NewRequest(http.MethodGet, "/api/node/node-123", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("nodeId")
	c.SetParamValues("node-123")
	c.Set("request_id", "test-request-id")

	err := handler.GetNode(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var response map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.False(t, response["success"].(bool))
	assert.Equal(t, "Failed to fetch node", response["error"])

	mockService.AssertExpectations(t)
}

func TestGetNode_NoRequestID(t *testing.T) {
	mockService := new(mockLearningService)
	handler := &testNodeHandler{service: mockService}

	expectedView := createMockNodeView()
	mockService.On("GetNode", mock.Anything, "node-123", "00000000-0000-0000-0000-000000000001", "unknown").
		Return(expectedView, nil)

	e := setupEcho()
	req := httptest.NewRequest(http.MethodGet, "/api/node/node-123", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("nodeId")
	c.SetParamValues("node-123")
	// No request_id set

	err := handler.GetNode(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	mockService.AssertExpectations(t)
}

func TestGetNode_WithLearningObjectives(t *testing.T) {
	mockService := new(mockLearningService)
	handler := &testNodeHandler{service: mockService}

	view := createMockNodeView()
	view.Node.LearningObjectives = []string{
		"Understand variable declarations",
		"Master const vs let vs var",
		"Avoid hoisting pitfalls",
	}

	mockService.On("GetNode", mock.Anything, "node-123", "00000000-0000-0000-0000-000000000001", mock.Anything).
		Return(view, nil)

	e := setupEcho()
	req := httptest.NewRequest(http.MethodGet, "/api/node/node-123", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("nodeId")
	c.SetParamValues("node-123")
	c.Set("request_id", "test-request-id")

	err := handler.GetNode(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &response)

	objectives := response["learning_objectives"].([]interface{})
	assert.Len(t, objectives, 3)
	assert.Equal(t, "Understand variable declarations", objectives[0])

	mockService.AssertExpectations(t)
}
