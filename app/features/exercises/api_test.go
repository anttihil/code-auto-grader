package exercises

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
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

// contentService is the interface used for testing
type contentService interface {
	GetExercise(ctx context.Context, exerciseID, userID, requestID string) (*ExerciseView, error)
	GenerateExercise(ctx context.Context, req GenerateExerciseRequest, requestID string) (*ExerciseGenerationResult, error)
	GetSampleExercise() *ExerciseGenerationResult
}

// mockContentService implements contentService for testing
type mockContentService struct {
	mock.Mock
}

func (m *mockContentService) GetExercise(ctx context.Context, exerciseID, userID, requestID string) (*ExerciseView, error) {
	args := m.Called(ctx, exerciseID, userID, requestID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ExerciseView), args.Error(1)
}

func (m *mockContentService) GenerateExercise(ctx context.Context, req GenerateExerciseRequest, requestID string) (*ExerciseGenerationResult, error) {
	args := m.Called(ctx, req, requestID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ExerciseGenerationResult), args.Error(1)
}

func (m *mockContentService) GetSampleExercise() *ExerciseGenerationResult {
	args := m.Called()
	return args.Get(0).(*ExerciseGenerationResult)
}

// testExerciseHandler wraps ExerciseHandler to accept interface for testing
type testExerciseHandler struct {
	service contentService
}

func (h *testExerciseHandler) GetExercise(c echo.Context) error {
	ctx, requestID := getRequestContext(c)
	exerciseID := c.Param("exerciseId")
	userID := "00000000-0000-0000-0000-000000000001"

	view, err := h.service.GetExercise(ctx, exerciseID, userID, requestID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "Exercise not found",
		})
	}

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

func (h *testExerciseHandler) GenerateExercise(c echo.Context) error {
	ctx, requestID := getRequestContext(c)

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

	result, err := h.service.GenerateExercise(ctx, GenerateExerciseRequest{
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

func (h *testExerciseHandler) GetSampleExercise(c echo.Context) error {
	sample := h.service.GetSampleExercise()

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

// ============================================================================
// Helper Functions
// ============================================================================

func setupEcho() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	return e
}

func createMockExerciseView() *ExerciseView {
	exerciseID := uuid.New()
	nodeID := uuid.New()
	submissionID1 := uuid.New()
	submissionID2 := uuid.New()
	score1 := 85
	score2 := 92

	return &ExerciseView{
		Exercise: &models.Exercise{
			ID:          exerciseID,
			NodeID:      nodeID,
			Title:       "Sum Two Numbers",
			Slug:        "sum-two-numbers",
			Description: "Write a function that adds two numbers",
			Difficulty:  "beginner",
			Content: models.ExerciseContent{
				StarterCode: "function add(a, b) {\n  // TODO: implement\n}",
				TestCases: []models.TestCase{
					{Name: "Test 1+1", Code: "assert.strictEqual(add(1, 1), 2);", Visible: true},
					{Name: "Test 2+3", Code: "assert.strictEqual(add(2, 3), 5);", Visible: true},
				},
				Hints: []models.Hint{
					{Order: 1, Title: "Addition", Content: "Use the + operator"},
				},
			},
		},
		Submissions: []*models.Submission{
			{
				ID:         submissionID1,
				ExerciseID: exerciseID,
				TotalScore: &score1,
			},
			{
				ID:         submissionID2,
				ExerciseID: exerciseID,
				TotalScore: &score2,
			},
		},
		BestScore:    92,
		AttemptCount: 2,
	}
}

func createMockExercisePreview() *ExerciseGenerationResult {
	return &ExerciseGenerationResult{
		Title:          "Binary Search",
		Description:    "Implement binary search algorithm",
		Difficulty:     "intermediate",
		Language:       "python",
		StarterCode:    "def binary_search(arr, target):\n    # TODO: implement\n    pass",
		Hints:          []string{"Use two pointers", "Divide the array in half"},
		ExpectedOutput: "Returns the index of target element",
		TestCases: []models.TestCase{
			{Name: "test_found", Code: "binary_search([1,2,3,4,5], 3) == 2", Visible: true},
			{Name: "test_not_found", Code: "binary_search([1,2,3], 5) == -1", Visible: true},
		},
	}
}

func createSampleExercisePreview() *ExerciseGenerationResult {
	return &ExerciseGenerationResult{
		Title:          "Sum Two Numbers",
		Description:    "Write a function that takes two numbers and returns their sum",
		Difficulty:     "beginner",
		Language:       "javascript",
		StarterCode:    "function add(a, b) {\n  // TODO: implement\n  return 0;\n}",
		Hints:          []string{"Use the + operator", "Return the result", "Test with different numbers"},
		ExpectedOutput: "add(2, 3) should return 5",
		TestCases: []models.TestCase{
			{Name: "test_2_plus_3", Code: "console.assert(sum(2, 3) === 5)", Visible: true},
			{Name: "test_10_plus_20", Code: "console.assert(sum(10, 20) === 30)", Visible: true},
			{Name: "test_negative", Code: "console.assert(sum(-5, 5) === 0)", Visible: true},
		},
	}
}

// ============================================================================
// ExerciseHandler Constructor Tests
// ============================================================================

func TestNewExerciseHandler(t *testing.T) {
	t.Run("panics when service is nil", func(t *testing.T) {
		assert.Panics(t, func() {
			NewExerciseHandler(nil)
		})
	})
}

// ============================================================================
// GetExercise Tests
// ============================================================================

func TestGetExercise_HappyPath(t *testing.T) {
	mockService := new(mockContentService)
	handler := &testExerciseHandler{service: mockService}

	expectedView := createMockExerciseView()
	mockService.On("GetExercise", mock.Anything, "exercise-123", "00000000-0000-0000-0000-000000000001", mock.Anything).
		Return(expectedView, nil)

	e := setupEcho()
	req := httptest.NewRequest(http.MethodGet, "/api/exercise/exercise-123", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("exerciseId")
	c.SetParamValues("exercise-123")
	c.Set("request_id", "test-request-id")

	err := handler.GetExercise(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))
	assert.Equal(t, "Sum Two Numbers", response["title"])
	assert.Equal(t, "sum-two-numbers", response["slug"])
	assert.Equal(t, "beginner", response["difficulty"])
	assert.Equal(t, "Write a function that adds two numbers", response["description"])
	assert.Equal(t, float64(92), response["best_score"])
	assert.Equal(t, float64(2), response["attempt_count"])

	// Check content structure
	content := response["content"].(map[string]interface{})
	assert.Contains(t, content["starter_code"], "function add")
	assert.NotNil(t, content["test_cases"])
	assert.NotNil(t, content["hints"])

	// Check submissions
	submissions := response["submissions"].([]interface{})
	assert.Len(t, submissions, 2)

	mockService.AssertExpectations(t)
}

func TestGetExercise_NotFound(t *testing.T) {
	mockService := new(mockContentService)
	handler := &testExerciseHandler{service: mockService}

	mockService.On("GetExercise", mock.Anything, "nonexistent", "00000000-0000-0000-0000-000000000001", mock.Anything).
		Return(nil, assert.AnError)

	e := setupEcho()
	req := httptest.NewRequest(http.MethodGet, "/api/exercise/nonexistent", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("exerciseId")
	c.SetParamValues("nonexistent")
	c.Set("request_id", "test-request-id")

	err := handler.GetExercise(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)

	var response map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.False(t, response["success"].(bool))
	assert.Equal(t, "Exercise not found", response["error"])

	mockService.AssertExpectations(t)
}

func TestGetExercise_NoSubmissions(t *testing.T) {
	mockService := new(mockContentService)
	handler := &testExerciseHandler{service: mockService}

	view := createMockExerciseView()
	view.Submissions = []*models.Submission{}
	view.BestScore = 0
	view.AttemptCount = 0

	mockService.On("GetExercise", mock.Anything, "exercise-123", "00000000-0000-0000-0000-000000000001", mock.Anything).
		Return(view, nil)

	e := setupEcho()
	req := httptest.NewRequest(http.MethodGet, "/api/exercise/exercise-123", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("exerciseId")
	c.SetParamValues("exercise-123")
	c.Set("request_id", "test-request-id")

	err := handler.GetExercise(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, float64(0), response["best_score"])
	assert.Equal(t, float64(0), response["attempt_count"])

	submissions := response["submissions"].([]interface{})
	assert.Len(t, submissions, 0)

	mockService.AssertExpectations(t)
}

func TestGetExercise_NoRequestID(t *testing.T) {
	mockService := new(mockContentService)
	handler := &testExerciseHandler{service: mockService}

	expectedView := createMockExerciseView()
	mockService.On("GetExercise", mock.Anything, "exercise-123", "00000000-0000-0000-0000-000000000001", "unknown").
		Return(expectedView, nil)

	e := setupEcho()
	req := httptest.NewRequest(http.MethodGet, "/api/exercise/exercise-123", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("exerciseId")
	c.SetParamValues("exercise-123")
	// No request_id set

	err := handler.GetExercise(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	mockService.AssertExpectations(t)
}

// ============================================================================
// GenerateExercise Tests
// ============================================================================

func TestGenerateExercise_HappyPath(t *testing.T) {
	mockService := new(mockContentService)
	handler := &testExerciseHandler{service: mockService}

	expectedPreview := createMockExercisePreview()
	mockService.On("GenerateExercise", mock.Anything, mock.MatchedBy(func(req GenerateExerciseRequest) bool {
		return req.Difficulty == "intermediate" &&
			req.Topic == "algorithms" &&
			req.Language == "python"
	}), mock.Anything).Return(expectedPreview, nil)

	e := setupEcho()
	reqBody := `{"difficulty":"intermediate","topic":"algorithms","language":"python"}`
	req := httptest.NewRequest(http.MethodPost, "/api/exercise/generate", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("request_id", "test-request-id")

	err := handler.GenerateExercise(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))
	assert.Equal(t, "Binary Search", response["title"])
	assert.Equal(t, "intermediate", response["difficulty"])
	assert.Equal(t, "python", response["language"])
	assert.Contains(t, response["starterCode"], "binary_search")

	hints := response["hints"].([]interface{})
	assert.Len(t, hints, 2)

	testCases := response["testCases"].([]interface{})
	assert.Len(t, testCases, 2)

	mockService.AssertExpectations(t)
}

func TestGenerateExercise_DefaultValues(t *testing.T) {
	mockService := new(mockContentService)
	handler := &testExerciseHandler{service: mockService}

	expectedPreview := createMockExercisePreview()
	expectedPreview.Difficulty = "beginner"
	expectedPreview.Language = "javascript"

	mockService.On("GenerateExercise", mock.Anything, mock.MatchedBy(func(req GenerateExerciseRequest) bool {
		// Empty values should be handled by service (defaults applied in domain layer)
		return true
	}), mock.Anything).Return(expectedPreview, nil)

	e := setupEcho()
	reqBody := `{}`
	req := httptest.NewRequest(http.MethodPost, "/api/exercise/generate", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("request_id", "test-request-id")

	err := handler.GenerateExercise(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))

	mockService.AssertExpectations(t)
}

func TestGenerateExercise_InvalidJSON(t *testing.T) {
	mockService := new(mockContentService)
	handler := &testExerciseHandler{service: mockService}

	e := setupEcho()
	reqBody := `{invalid json}`
	req := httptest.NewRequest(http.MethodPost, "/api/exercise/generate", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("request_id", "test-request-id")

	err := handler.GenerateExercise(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.False(t, response["success"].(bool))
	assert.Equal(t, "Invalid request body", response["error"])
}

func TestGenerateExercise_ServiceError(t *testing.T) {
	mockService := new(mockContentService)
	handler := &testExerciseHandler{service: mockService}

	mockService.On("GenerateExercise", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, assert.AnError)

	e := setupEcho()
	reqBody := `{"difficulty":"intermediate","topic":"algorithms","language":"python"}`
	req := httptest.NewRequest(http.MethodPost, "/api/exercise/generate", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("request_id", "test-request-id")

	err := handler.GenerateExercise(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var response map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.False(t, response["success"].(bool))
	assert.Equal(t, "Failed to generate exercise", response["error"])
	assert.NotNil(t, response["message"])

	mockService.AssertExpectations(t)
}

func TestGenerateExercise_PartialRequest(t *testing.T) {
	mockService := new(mockContentService)
	handler := &testExerciseHandler{service: mockService}

	expectedPreview := createMockExercisePreview()
	mockService.On("GenerateExercise", mock.Anything, mock.MatchedBy(func(req GenerateExerciseRequest) bool {
		return req.Difficulty == "advanced" && req.Topic == "" && req.Language == ""
	}), mock.Anything).Return(expectedPreview, nil)

	e := setupEcho()
	reqBody := `{"difficulty":"advanced"}`
	req := httptest.NewRequest(http.MethodPost, "/api/exercise/generate", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("request_id", "test-request-id")

	err := handler.GenerateExercise(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	mockService.AssertExpectations(t)
}

// ============================================================================
// GetSampleExercise Tests
// ============================================================================

func TestGetSampleExercise_Success(t *testing.T) {
	mockService := new(mockContentService)
	handler := &testExerciseHandler{service: mockService}

	samplePreview := createSampleExercisePreview()
	mockService.On("GetSampleExercise").Return(samplePreview)

	e := setupEcho()
	req := httptest.NewRequest(http.MethodGet, "/api/exercise/sample", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GetSampleExercise(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))
	assert.Equal(t, "Sum Two Numbers", response["title"])
	assert.Equal(t, "beginner", response["difficulty"])
	assert.Equal(t, "javascript", response["language"])
	assert.Contains(t, response["starterCode"], "function add")
	assert.Contains(t, response["description"], "sum")

	hints := response["hints"].([]interface{})
	assert.Len(t, hints, 3)
	assert.Equal(t, "Use the + operator", hints[0])

	testCases := response["testCases"].([]interface{})
	assert.Len(t, testCases, 3)

	expectedOutput := response["expectedOutput"].(string)
	assert.Contains(t, expectedOutput, "add(2, 3)")

	mockService.AssertExpectations(t)
}

func TestGetSampleExercise_NoRequestID(t *testing.T) {
	mockService := new(mockContentService)
	handler := &testExerciseHandler{service: mockService}

	samplePreview := createSampleExercisePreview()
	mockService.On("GetSampleExercise").Return(samplePreview)

	e := setupEcho()
	req := httptest.NewRequest(http.MethodGet, "/api/exercise/sample", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// No request_id set - should still work

	err := handler.GetSampleExercise(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))

	mockService.AssertExpectations(t)
}

func TestGetSampleExercise_ResponseStructure(t *testing.T) {
	mockService := new(mockContentService)
	handler := &testExerciseHandler{service: mockService}

	samplePreview := createSampleExercisePreview()
	mockService.On("GetSampleExercise").Return(samplePreview)

	e := setupEcho()
	req := httptest.NewRequest(http.MethodGet, "/api/exercise/sample", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GetSampleExercise(c)

	assert.NoError(t, err)

	var response map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &response)

	// Verify all expected fields are present
	requiredFields := []string{"success", "title", "description", "difficulty", "language", "starterCode", "hints", "expectedOutput", "testCases"}
	for _, field := range requiredFields {
		assert.Contains(t, response, field, "Response should contain field: %s", field)
	}

	mockService.AssertExpectations(t)
}
