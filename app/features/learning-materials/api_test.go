package learningmaterials

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

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

// ===== Handler Tests =====

func TestNewHandler(t *testing.T) {
	t.Run("creates handler with valid service", func(t *testing.T) {
		mockRepo := new(mockLearningMaterialRepository)
		service := NewService(mockRepo)
		handler := NewHandler(service)

		assert.NotNil(t, handler)
	})

	t.Run("panics with nil service", func(t *testing.T) {
		assert.Panics(t, func() {
			NewHandler(nil)
		})
	})
}

// ===== GetByNode Handler Tests =====

func TestGetByNode(t *testing.T) {
	nodeID := uuid.New()

	t.Run("returns materials for node", func(t *testing.T) {
		mockRepo := new(mockLearningMaterialRepository)
		service := NewService(mockRepo)
		handler := NewHandler(service)

		materials := []*models.LearningMaterial{
			{
				ID:         uuid.New(),
				NodeID:     nodeID,
				Title:      "Introduction",
				Content:    "Intro content",
				OrderIndex: 0,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
			{
				ID:         uuid.New(),
				NodeID:     nodeID,
				Title:      "Key Concepts",
				Content:    "Concepts content",
				OrderIndex: 1,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
		}

		mockRepo.On("GetByNodeID", mock.Anything, nodeID.String()).Return(materials, nil)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/nodes/"+nodeID.String()+"/materials", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("nodeId")
		c.SetParamValues(nodeID.String())

		err := handler.GetByNode(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response []LearningMaterialResponse
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Len(t, response, 2)
		assert.Equal(t, "Introduction", response[0].Title)
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns empty array for node with no materials", func(t *testing.T) {
		mockRepo := new(mockLearningMaterialRepository)
		service := NewService(mockRepo)
		handler := NewHandler(service)

		mockRepo.On("GetByNodeID", mock.Anything, nodeID.String()).Return([]*models.LearningMaterial{}, nil)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/nodes/"+nodeID.String()+"/materials", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("nodeId")
		c.SetParamValues(nodeID.String())

		err := handler.GetByNode(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response []LearningMaterialResponse
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Empty(t, response)
		mockRepo.AssertExpectations(t)
	})
}

// ===== GetByID Handler Tests =====

func TestHandler_GetByID(t *testing.T) {
	materialID := uuid.New()
	nodeID := uuid.New()

	t.Run("returns material by ID", func(t *testing.T) {
		mockRepo := new(mockLearningMaterialRepository)
		service := NewService(mockRepo)
		handler := NewHandler(service)

		material := &models.LearningMaterial{
			ID:         materialID,
			NodeID:     nodeID,
			Title:      "Test Material",
			Content:    "Test content",
			OrderIndex: 0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		mockRepo.On("GetByID", mock.Anything, materialID.String()).Return(material, nil)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/materials/"+materialID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(materialID.String())

		err := handler.GetByID(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response LearningMaterialResponse
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, "Test Material", response.Title)
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns 404 when not found", func(t *testing.T) {
		mockRepo := new(mockLearningMaterialRepository)
		service := NewService(mockRepo)
		handler := NewHandler(service)

		mockRepo.On("GetByID", mock.Anything, materialID.String()).Return(nil, errors.New("not found"))

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/materials/"+materialID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(materialID.String())

		err := handler.GetByID(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockRepo.AssertExpectations(t)
	})
}

// ===== Create Handler Tests =====

func TestHandler_Create(t *testing.T) {
	nodeID := uuid.New()

	t.Run("creates material successfully", func(t *testing.T) {
		mockRepo := new(mockLearningMaterialRepository)
		service := NewService(mockRepo)
		handler := NewHandler(service)

		mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(m *models.LearningMaterial) bool {
			return m.Title == "New Material"
		})).Return(nil)

		body := `{"title": "New Material", "content": "Material content", "order_index": 1}`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/api/nodes/"+nodeID.String()+"/materials", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("nodeId")
		c.SetParamValues(nodeID.String())

		err := handler.Create(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var response LearningMaterialResponse
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, "New Material", response.Title)
		assert.Equal(t, "Material content", response.Content)
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns 400 for missing title", func(t *testing.T) {
		mockRepo := new(mockLearningMaterialRepository)
		service := NewService(mockRepo)
		handler := NewHandler(service)

		body := `{"content": "Material content"}`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/api/nodes/"+nodeID.String()+"/materials", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("nodeId")
		c.SetParamValues(nodeID.String())

		err := handler.Create(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("returns 400 for missing content", func(t *testing.T) {
		mockRepo := new(mockLearningMaterialRepository)
		service := NewService(mockRepo)
		handler := NewHandler(service)

		body := `{"title": "New Material"}`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/api/nodes/"+nodeID.String()+"/materials", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("nodeId")
		c.SetParamValues(nodeID.String())

		err := handler.Create(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("returns 400 for invalid JSON", func(t *testing.T) {
		mockRepo := new(mockLearningMaterialRepository)
		service := NewService(mockRepo)
		handler := NewHandler(service)

		body := `{invalid json}`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/api/nodes/"+nodeID.String()+"/materials", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("nodeId")
		c.SetParamValues(nodeID.String())

		err := handler.Create(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

// ===== Update Handler Tests =====

func TestHandler_Update(t *testing.T) {
	materialID := uuid.New()
	nodeID := uuid.New()

	t.Run("updates material successfully", func(t *testing.T) {
		mockRepo := new(mockLearningMaterialRepository)
		service := NewService(mockRepo)
		handler := NewHandler(service)

		existingMaterial := &models.LearningMaterial{
			ID:         materialID,
			NodeID:     nodeID,
			Title:      "Old Title",
			Content:    "Old content",
			OrderIndex: 0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		mockRepo.On("GetByID", mock.Anything, materialID.String()).Return(existingMaterial, nil)
		mockRepo.On("Update", mock.Anything, existingMaterial).Return(nil)

		body := `{"title": "New Title", "content": "New content"}`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPut, "/api/materials/"+materialID.String(), strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(materialID.String())

		err := handler.Update(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response LearningMaterialResponse
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, "New Title", response.Title)
		mockRepo.AssertExpectations(t)
	})
}

// ===== Delete Handler Tests =====

func TestHandler_Delete(t *testing.T) {
	materialID := uuid.New()

	t.Run("deletes material successfully", func(t *testing.T) {
		mockRepo := new(mockLearningMaterialRepository)
		service := NewService(mockRepo)
		handler := NewHandler(service)

		mockRepo.On("Delete", mock.Anything, materialID.String()).Return(nil)

		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/api/materials/"+materialID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(materialID.String())

		err := handler.Delete(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns 500 when delete fails", func(t *testing.T) {
		mockRepo := new(mockLearningMaterialRepository)
		service := NewService(mockRepo)
		handler := NewHandler(service)

		mockRepo.On("Delete", mock.Anything, materialID.String()).Return(errors.New("not found"))

		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/api/materials/"+materialID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(materialID.String())

		err := handler.Delete(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockRepo.AssertExpectations(t)
	})
}

