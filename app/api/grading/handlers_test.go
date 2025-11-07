package grading

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	domainGrading "github.com/anttihil/code-auto-grader/domain/grading"
	"github.com/anttihil/code-auto-grader/internal/logger"
	"github.com/anttihil/code-auto-grader/models"
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
// Mock Grading Service
// ============================================================================

type mockGradingService struct {
	mock.Mock
}

func (m *mockGradingService) SubmitCode(ctx context.Context, req domainGrading.SubmitCodeRequest, requestID string) (*domainGrading.SubmissionResult, error) {
	args := m.Called(ctx, req, requestID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainGrading.SubmissionResult), args.Error(1)
}

func (m *mockGradingService) RunCode(ctx context.Context, req domainGrading.RunCodeRequest, requestID string) (*domainGrading.RunCodeResult, error) {
	args := m.Called(ctx, req, requestID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainGrading.RunCodeResult), args.Error(1)
}

func (m *mockGradingService) GetPoolStatus() map[string]interface{} {
	args := m.Called()
	return args.Get(0).(map[string]interface{})
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

func createMockSubmissionResult() *domainGrading.SubmissionResult {
	return &domainGrading.SubmissionResult{
		SubmissionID:    "submission-123",
		ExecutionOutput: "Test output",
		ExecutionError:  nil,
		TestResults: &models.TestResults{
			TotalTests:  2,
			PassedTests: 2,
			FailedTests: 0,
			Results:     []models.TestResult{},
		},
		AIFeedback: &models.AIFeedback{
			Passed:       true,
			Score:        90,
			Feedback:     "Great job!",
			Strengths:    []string{"Clean code"},
			Improvements: []string{},
		},
		TestScore:     100,
		AIScore:       90,
		TotalScore:    96,
		AttemptNumber: 1,
	}
}

func createMockRunCodeResult() *domainGrading.RunCodeResult {
	return &domainGrading.RunCodeResult{
		ExecutionOutput: "Test output",
		ExecutionError:  nil,
		TestResults: &models.TestResults{
			TotalTests:  2,
			PassedTests: 2,
			FailedTests: 0,
			Results:     []models.TestResult{},
		},
		TestScore: 100,
	}
}

// ============================================================================
// NewSubmissionHandler Tests
// ============================================================================

func TestNewSubmissionHandler(t *testing.T) {
	t.Run("creates handler with service", func(t *testing.T) {
		mockService := new(mockGradingService)
		handler := NewSubmissionHandler(mockService)
		assert.NotNil(t, handler)
	})

	t.Run("panics when service is nil", func(t *testing.T) {
		assert.Panics(t, func() {
			NewSubmissionHandler(nil)
		})
	})
}

// ============================================================================
// SubmitToExercise Tests
// ============================================================================

func TestSubmitToExercise_HappyPath(t *testing.T) {
	mockService := new(mockGradingService)
	handler := NewSubmissionHandler(mockService)

	// Setup mock expectation
	expectedResult := createMockSubmissionResult()
	mockService.On("SubmitCode", mock.Anything, mock.MatchedBy(func(req domainGrading.SubmitCodeRequest) bool {
		return req.Code == "function add(a, b) { return a + b; }" &&
			req.Language == "javascript" &&
			req.ExerciseID == "exercise-123"
	}), mock.Anything).Return(expectedResult, nil)

	// Create request
	e := setupEcho()
	reqBody := `{"code":"function add(a, b) { return a + b; }","language":"javascript"}`
	req := httptest.NewRequest(http.MethodPost, "/api/submission/exercise-123", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("exerciseId")
	c.SetParamValues("exercise-123")
	c.Set("request_id", "test-request-id")

	// Execute
	err := handler.SubmitToExercise(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response SubmissionResponse
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, "submission-123", response.SubmissionID)
	assert.Equal(t, 96, response.TotalScore)
	assert.Equal(t, 1, response.AttemptNumber)

	mockService.AssertExpectations(t)
}

func TestSubmitToExercise_DefaultLanguage(t *testing.T) {
	mockService := new(mockGradingService)
	handler := NewSubmissionHandler(mockService)

	expectedResult := createMockSubmissionResult()
	mockService.On("SubmitCode", mock.Anything, mock.MatchedBy(func(req domainGrading.SubmitCodeRequest) bool {
		return req.Language == "javascript" // Should default to javascript
	}), mock.Anything).Return(expectedResult, nil)

	e := setupEcho()
	reqBody := `{"code":"function add(a, b) { return a + b; }"}` // No language specified
	req := httptest.NewRequest(http.MethodPost, "/api/submission/exercise-123", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("exerciseId")
	c.SetParamValues("exercise-123")
	c.Set("request_id", "test-request-id")

	err := handler.SubmitToExercise(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}

func TestSubmitToExercise_InvalidJSON(t *testing.T) {
	mockService := new(mockGradingService)
	handler := NewSubmissionHandler(mockService)

	e := setupEcho()
	reqBody := `{"code":` // Invalid JSON
	req := httptest.NewRequest(http.MethodPost, "/api/submission/exercise-123", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("exerciseId")
	c.SetParamValues("exercise-123")
	c.Set("request_id", "test-request-id")

	err := handler.SubmitToExercise(c)

	assert.NoError(t, err) // Handler doesn't return error, sends JSON error response
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response ErrorResponse
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, "Invalid request body", response.Error)
}

func TestSubmitToExercise_MissingCode(t *testing.T) {
	mockService := new(mockGradingService)
	handler := NewSubmissionHandler(mockService)

	e := setupEcho()
	reqBody := `{"language":"javascript"}` // Missing code
	req := httptest.NewRequest(http.MethodPost, "/api/submission/exercise-123", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("exerciseId")
	c.SetParamValues("exercise-123")
	c.Set("request_id", "test-request-id")

	err := handler.SubmitToExercise(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response ErrorResponse
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, "Code is required", response.Error)
}

func TestSubmitToExercise_IdenticalToStarter(t *testing.T) {
	mockService := new(mockGradingService)
	handler := NewSubmissionHandler(mockService)

	// Mock returns similarity error
	mockService.On("SubmitCode", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, errors.New("IDENTICAL_TO_STARTER: submission is identical to the starter code"))

	e := setupEcho()
	reqBody := `{"code":"// TODO: implement"}`
	req := httptest.NewRequest(http.MethodPost, "/api/submission/exercise-123", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("exerciseId")
	c.SetParamValues("exercise-123")
	c.Set("request_id", "test-request-id")

	err := handler.SubmitToExercise(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response ErrorResponse
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, "IDENTICAL_TO_STARTER", response.Error)
	assert.Contains(t, response.Message, "identical")

	mockService.AssertExpectations(t)
}

func TestSubmitToExercise_TrivialChanges(t *testing.T) {
	mockService := new(mockGradingService)
	handler := NewSubmissionHandler(mockService)

	mockService.On("SubmitCode", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, errors.New("TRIVIAL_CHANGES: submission only contains trivial changes"))

	e := setupEcho()
	reqBody := `{"code":"// My solution\n// TODO: implement"}`
	req := httptest.NewRequest(http.MethodPost, "/api/submission/exercise-123", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("exerciseId")
	c.SetParamValues("exercise-123")
	c.Set("request_id", "test-request-id")

	err := handler.SubmitToExercise(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response ErrorResponse
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, "TRIVIAL_CHANGES", response.Error)
	assert.Contains(t, response.Message, "trivial")

	mockService.AssertExpectations(t)
}

func TestSubmitToExercise_ServiceError(t *testing.T) {
	mockService := new(mockGradingService)
	handler := NewSubmissionHandler(mockService)

	mockService.On("SubmitCode", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, errors.New("database connection failed"))

	e := setupEcho()
	reqBody := `{"code":"function add(a, b) { return a + b; }"}`
	req := httptest.NewRequest(http.MethodPost, "/api/submission/exercise-123", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("exerciseId")
	c.SetParamValues("exercise-123")
	c.Set("request_id", "test-request-id")

	err := handler.SubmitToExercise(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var response ErrorResponse
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, "Failed to process submission", response.Error)

	mockService.AssertExpectations(t)
}

// ============================================================================
// RunCode Tests
// ============================================================================

func TestRunCode_HappyPath(t *testing.T) {
	mockService := new(mockGradingService)
	handler := NewSubmissionHandler(mockService)

	expectedResult := createMockRunCodeResult()
	mockService.On("RunCode", mock.Anything, mock.MatchedBy(func(req domainGrading.RunCodeRequest) bool {
		return req.Code == "function add(a, b) { return a + b; }" &&
			req.Language == "javascript" &&
			req.ExerciseID == "exercise-123"
	}), mock.Anything).Return(expectedResult, nil)

	e := setupEcho()
	reqBody := `{"code":"function add(a, b) { return a + b; }","language":"javascript"}`
	req := httptest.NewRequest(http.MethodPost, "/api/code/run/exercise-123", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("exerciseId")
	c.SetParamValues("exercise-123")
	c.Set("request_id", "test-request-id")

	err := handler.RunCode(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response RunCodeResponse
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, 100, response.TestScore)
	assert.NotNil(t, response.TestResults)

	mockService.AssertExpectations(t)
}

func TestRunCode_MissingCode(t *testing.T) {
	mockService := new(mockGradingService)
	handler := NewSubmissionHandler(mockService)

	e := setupEcho()
	reqBody := `{"language":"javascript"}`
	req := httptest.NewRequest(http.MethodPost, "/api/code/run/exercise-123", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("exerciseId")
	c.SetParamValues("exercise-123")
	c.Set("request_id", "test-request-id")

	err := handler.RunCode(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response ErrorResponse
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, "Code is required", response.Error)
}

func TestRunCode_InvalidJSON(t *testing.T) {
	mockService := new(mockGradingService)
	handler := NewSubmissionHandler(mockService)

	e := setupEcho()
	reqBody := `{invalid json}`
	req := httptest.NewRequest(http.MethodPost, "/api/code/run/exercise-123", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("exerciseId")
	c.SetParamValues("exercise-123")
	c.Set("request_id", "test-request-id")

	err := handler.RunCode(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response ErrorResponse
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, "Invalid request body", response.Error)
}

func TestRunCode_ServiceError(t *testing.T) {
	mockService := new(mockGradingService)
	handler := NewSubmissionHandler(mockService)

	mockService.On("RunCode", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, errors.New("exercise not found"))

	e := setupEcho()
	reqBody := `{"code":"function add(a, b) { return a + b; }"}`
	req := httptest.NewRequest(http.MethodPost, "/api/code/run/exercise-123", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("exerciseId")
	c.SetParamValues("exercise-123")
	c.Set("request_id", "test-request-id")

	err := handler.RunCode(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var response ErrorResponse
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, "Failed to run code", response.Error)

	mockService.AssertExpectations(t)
}

func TestRunCode_DefaultLanguage(t *testing.T) {
	mockService := new(mockGradingService)
	handler := NewSubmissionHandler(mockService)

	expectedResult := createMockRunCodeResult()
	mockService.On("RunCode", mock.Anything, mock.MatchedBy(func(req domainGrading.RunCodeRequest) bool {
		return req.Language == "javascript" // Should default
	}), mock.Anything).Return(expectedResult, nil)

	e := setupEcho()
	reqBody := `{"code":"function add(a, b) { return a + b; }"}` // No language
	req := httptest.NewRequest(http.MethodPost, "/api/code/run/exercise-123", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("exerciseId")
	c.SetParamValues("exercise-123")
	c.Set("request_id", "test-request-id")

	err := handler.RunCode(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}

// ============================================================================
// GetPoolStatus Tests
// ============================================================================

func TestGetPoolStatus_Success(t *testing.T) {
	mockService := new(mockGradingService)
	handler := NewSubmissionHandler(mockService)

	expectedStatus := map[string]interface{}{
		"javascript": map[string]interface{}{
			"size":      5,
			"available": 5,
		},
		"python": map[string]interface{}{
			"size":      5,
			"available": 3,
		},
	}
	mockService.On("GetPoolStatus").Return(expectedStatus)

	e := setupEcho()
	req := httptest.NewRequest(http.MethodGet, "/api/pool/status", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GetPoolStatus(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NotNil(t, response["javascript"])
	assert.NotNil(t, response["python"])

	mockService.AssertExpectations(t)
}

func TestGetPoolStatus_EmptyPools(t *testing.T) {
	mockService := new(mockGradingService)
	handler := NewSubmissionHandler(mockService)

	expectedStatus := map[string]interface{}{}
	mockService.On("GetPoolStatus").Return(expectedStatus)

	e := setupEcho()
	req := httptest.NewRequest(http.MethodGet, "/api/pool/status", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GetPoolStatus(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &response)
	// Response should be an empty map when no pools exist
	assert.Empty(t, response)

	mockService.AssertExpectations(t)
}
