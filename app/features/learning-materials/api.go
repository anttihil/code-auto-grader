package learningmaterials

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/anttihil/code-auto-grader/internal/logger"
	"github.com/labstack/echo/v4"
)

// Handler handles HTTP requests for learning materials
type Handler struct {
	service *Service
}

// NewHandler creates a new learning materials handler
func NewHandler(service *Service) *Handler {
	if service == nil {
		panic("service is required")
	}
	return &Handler{service: service}
}

// GetByNode handles GET /api/nodes/:nodeId/materials
// Retrieves all learning materials for a node
func (h *Handler) GetByNode(c echo.Context) error {
	ctx, requestID := getRequestContext(c)
	nodeID := c.Param("nodeId")

	logger.Logger.Info("fetching learning materials for node",
		slog.String("request_id", requestID),
		slog.String("node_id", nodeID),
	)

	materials, err := h.service.GetByNodeID(ctx, nodeID)
	if err != nil {
		logger.Logger.Error("failed to fetch learning materials",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to fetch learning materials",
			Message: err.Error(),
		})
	}

	// Convert to DTOs
	response := make([]LearningMaterialResponse, len(materials))
	for i, m := range materials {
		response[i] = LearningMaterialResponse{
			ID:         m.ID.String(),
			NodeID:     m.NodeID.String(),
			Title:      m.Title,
			Content:    m.Content,
			OrderIndex: m.OrderIndex,
			CreatedAt:  m.CreatedAt,
			UpdatedAt:  m.UpdatedAt,
		}
	}

	return c.JSON(http.StatusOK, response)
}

// GetByID handles GET /api/materials/:id
// Retrieves a single learning material by ID
func (h *Handler) GetByID(c echo.Context) error {
	ctx, requestID := getRequestContext(c)
	id := c.Param("id")

	material, err := h.service.GetByID(ctx, id)
	if err != nil {
		logger.Logger.Error("failed to fetch learning material",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		return c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Learning material not found",
			Message: err.Error(),
		})
	}

	response := LearningMaterialResponse{
		ID:         material.ID.String(),
		NodeID:     material.NodeID.String(),
		Title:      material.Title,
		Content:    material.Content,
		OrderIndex: material.OrderIndex,
		CreatedAt:  material.CreatedAt,
		UpdatedAt:  material.UpdatedAt,
	}

	return c.JSON(http.StatusOK, response)
}

// Create handles POST /api/nodes/:nodeId/materials
// Creates a new learning material for a node
func (h *Handler) Create(c echo.Context) error {
	ctx, requestID := getRequestContext(c)
	nodeID := c.Param("nodeId")

	var req CreateLearningMaterialRequest
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

	// Basic validation
	if req.Title == "" || req.Content == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Missing required fields",
			Message: "title and content are required",
		})
	}

	material, err := h.service.Create(ctx, nodeID, req)
	if err != nil {
		logger.Logger.Error("failed to create learning material",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to create learning material",
			Message: err.Error(),
		})
	}

	response := LearningMaterialResponse{
		ID:         material.ID.String(),
		NodeID:     material.NodeID.String(),
		Title:      material.Title,
		Content:    material.Content,
		OrderIndex: material.OrderIndex,
		CreatedAt:  material.CreatedAt,
		UpdatedAt:  material.UpdatedAt,
	}

	return c.JSON(http.StatusCreated, response)
}

// Update handles PUT /api/materials/:id
// Updates an existing learning material
func (h *Handler) Update(c echo.Context) error {
	ctx, requestID := getRequestContext(c)
	id := c.Param("id")

	var req UpdateLearningMaterialRequest
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

	material, err := h.service.Update(ctx, id, req)
	if err != nil {
		logger.Logger.Error("failed to update learning material",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to update learning material",
			Message: err.Error(),
		})
	}

	response := LearningMaterialResponse{
		ID:         material.ID.String(),
		NodeID:     material.NodeID.String(),
		Title:      material.Title,
		Content:    material.Content,
		OrderIndex: material.OrderIndex,
		CreatedAt:  material.CreatedAt,
		UpdatedAt:  material.UpdatedAt,
	}

	return c.JSON(http.StatusOK, response)
}

// Delete handles DELETE /api/materials/:id
// Deletes a learning material
func (h *Handler) Delete(c echo.Context) error {
	ctx, requestID := getRequestContext(c)
	id := c.Param("id")

	if err := h.service.Delete(ctx, id); err != nil {
		logger.Logger.Error("failed to delete learning material",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to delete learning material",
			Message: err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent)
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
