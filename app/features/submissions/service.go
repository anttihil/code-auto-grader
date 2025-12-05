// Package grading contains the business logic for code submission, evaluation, and grading.
package submissions

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/anttihil/code-auto-grader/features/progress"
	"github.com/anttihil/code-auto-grader/internal/logger"
	"github.com/anttihil/code-auto-grader/models"
	"github.com/anttihil/code-auto-grader/repository"
	"github.com/google/uuid"
)

// GradingService handles code submission, test execution, AI evaluation, and scoring
type GradingService struct {
	executor                 CodeExecutor
	evaluator                AIEvaluator
	submissionRepo           repository.SubmissionRepository
	exerciseRepo             repository.ExerciseRepository
	userExerciseProgressRepo repository.UserExerciseProgressRepository
	progressService          *progress.Service
}

// NewGradingService creates a new grading service with required dependencies
func NewGradingService(
	executor CodeExecutor,
	evaluator AIEvaluator,
	submissionRepo repository.SubmissionRepository,
	exerciseRepo repository.ExerciseRepository,
	userExerciseProgressRepo repository.UserExerciseProgressRepository,
	progressService *progress.Service,
) *GradingService {
	if executor == nil {
		panic("code executor is required")
	}
	if evaluator == nil {
		panic("AI evaluator is required")
	}
	if submissionRepo == nil {
		panic("submission repository is required")
	}
	if exerciseRepo == nil {
		panic("exercise repository is required")
	}
	if userExerciseProgressRepo == nil {
		panic("user exercise progress repository is required")
	}
	if progressService == nil {
		panic("progress service is required")
	}

	return &GradingService{
		executor:                 executor,
		evaluator:                evaluator,
		submissionRepo:           submissionRepo,
		exerciseRepo:             exerciseRepo,
		userExerciseProgressRepo: userExerciseProgressRepo,
		progressService:          progressService,
	}
}

// SubmitCode handles the complete submission workflow:
// 1. Fetch exercise
// 2. Check code similarity to starter code
// 3. Run test cases
// 4. Get AI evaluation
// 5. Calculate scores
// 6. Save submission
// 7. Update user progress
func (s *GradingService) SubmitCode(ctx context.Context, req SubmitCodeRequest, requestID string) (*SubmissionResult, error) {
	// 1. Fetch exercise from database
	exercise, err := s.exerciseRepo.GetByID(ctx, req.ExerciseID)
	if err != nil {
		return nil, fmt.Errorf("exercise not found: %w", err)
	}

	logger.Logger.Info("exercise fetched for submission",
		slog.String("request_id", requestID),
		slog.String("exercise_title", exercise.Title),
		slog.Int("test_case_count", len(exercise.Content.TestCases)),
	)

	// 2. Check code similarity to starter code (anti-gaming validation)
	similarityLevel := s.checkCodeSimilarity(req.Code, exercise.Content.StarterCode, req.Language)

	switch similarityLevel {
	case CodeIdentical:
		logger.Logger.Warn("submission rejected: identical to starter code",
			slog.String("request_id", requestID),
			slog.String("exercise_id", req.ExerciseID),
			slog.String("user_id", req.UserID),
		)
		return nil, fmt.Errorf("IDENTICAL_TO_STARTER: submission is identical to the starter code")

	case CodeTrivialChanges:
		logger.Logger.Warn("submission rejected: trivial changes only",
			slog.String("request_id", requestID),
			slog.String("exercise_id", req.ExerciseID),
			slog.String("user_id", req.UserID),
		)
		return nil, fmt.Errorf("TRIVIAL_CHANGES: submission only contains trivial changes")

	case CodeMinimalChanges:
		logger.Logger.Info("submission flagged: minimal changes detected",
			slog.String("request_id", requestID),
			slog.String("exercise_id", req.ExerciseID),
			slog.String("user_id", req.UserID),
		)
		// Continue but flag for stricter AI evaluation

	case CodeSubstantial:
		logger.Logger.Debug("submission accepted: substantial changes detected",
			slog.String("request_id", requestID),
		)
		// Normal flow
	}

	// 3. Run test cases
	testResults, executionOutput, executionError, executionTimeMs := s.runTestCases(ctx, req.Code, req.Language, exercise.Content.TestCases, requestID)

	// Calculate test score (0-100)
	testScore := CalculateTestScore(testResults)

	logger.Logger.Info("test cases executed",
		slog.String("request_id", requestID),
		slog.Int("passed", testResults.PassedTests),
		slog.Int("total", testResults.TotalTests),
		slog.Int("test_score", testScore),
	)

	// 4. Get AI feedback (conditional based on exercise configuration)
	var aiFeedback *models.AIFeedback
	var aiScore int
	var aiEvaluationTimeMs int

	if s.shouldRequestAIFeedback(exercise) {
		startAI := time.Now()
		aiFeedback, aiScore, err = s.evaluateWithAI(ctx, exercise, req.Code, req.Language, testResults, requestID)
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
		aiEvaluationTimeMs = int(time.Since(startAI).Milliseconds())

		logger.Logger.Info("AI evaluation completed",
			slog.String("request_id", requestID),
			slog.Int("ai_score", aiScore),
			slog.Int("duration_ms", aiEvaluationTimeMs),
		)
	} else {
		logger.Logger.Info("AI evaluation skipped (weight=0 or no prompt configured)",
			slog.String("request_id", requestID),
			slog.Int("ai_weight", exercise.Content.GradingRubric.AIWeight),
		)
		// No AI feedback needed - score based purely on tests
		aiScore = 0
		aiEvaluationTimeMs = 0
	}

	// 5. Calculate total score using grading rubric weights
	totalScore := CalculateScore(testScore, aiScore, exercise.Content.GradingRubric)

	// Scale scores for database storage
	testWeight := exercise.Content.GradingRubric.TestWeight
	aiWeight := exercise.Content.GradingRubric.AIWeight
	if testWeight+aiWeight == 0 {
		testWeight = 60
		aiWeight = 40
	}

	aiScoreScaled := ScaleScoreToWeight(aiScore, aiWeight)
	testScoreScaled := ScaleScoreToWeight(testScore, testWeight)

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

	// 6. Save submission to database
	submission, err := s.saveSubmission(ctx, req.UserID, req.ExerciseID, req.Code, req.Language,
		&executionOutput, executionError, &executionTimeMs,
		testResults, aiFeedback, &aiScoreScaled, &aiEvaluationTimeMs,
		&testScoreScaled, &totalScore, requestID)
	if err != nil {
		return nil, fmt.Errorf("failed to save submission: %w", err)
	}

	logger.Logger.Info("submission saved successfully",
		slog.String("request_id", requestID),
		slog.String("submission_id", submission.ID.String()),
	)

	// 7. Update per-exercise progress (increment attempts count)
	if err := s.userExerciseProgressRepo.UpdateHelpCounts(ctx, req.UserID, req.ExerciseID, 0, 1); err != nil {
		logger.Logger.Warn("failed to update exercise progress attempts",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		// Don't fail the submission if progress update fails
	}

	// 8. Update user progress and mastery tracking
	if err := s.updateUserProgress(ctx, req.UserID, exercise.NodeID.String(), totalScore, requestID); err != nil {
		logger.Logger.Warn("failed to update user progress",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		// Don't fail the submission if progress update fails
	}

	// Return submission result
	return &SubmissionResult{
		SubmissionID:    submission.ID.String(),
		ExecutionOutput: executionOutput,
		ExecutionError:  executionError,
		TestResults:     testResults,
		AIFeedback:      aiFeedback,
		TestScore:       testScore,
		AIScore:         aiScore,
		TotalScore:      totalScore,
	}, nil
}

// RunCode executes code with test cases (quick test execution without AI evaluation or saving)
func (s *GradingService) RunCode(ctx context.Context, req RunCodeRequest, requestID string) (*RunCodeResult, error) {
	// 1. Fetch exercise from database to get test cases
	exercise, err := s.exerciseRepo.GetByID(ctx, req.ExerciseID)
	if err != nil {
		return nil, fmt.Errorf("exercise not found: %w", err)
	}

	logger.Logger.Info("exercise fetched for code run",
		slog.String("request_id", requestID),
		slog.String("exercise_title", exercise.Title),
		slog.Int("test_case_count", len(exercise.Content.TestCases)),
	)

	// 2. Run test cases using the same logic as SubmitCode
	testResults, executionOutput, executionError, _ := s.runTestCases(ctx, req.Code, req.Language, exercise.Content.TestCases, requestID)

	// Calculate test score
	testScore := CalculateTestScore(testResults)

	logger.Logger.Info("code run completed",
		slog.String("request_id", requestID),
		slog.Int("passed", testResults.PassedTests),
		slog.Int("total", testResults.TotalTests),
		slog.Int("test_score", testScore),
	)

	return &RunCodeResult{
		ExecutionOutput: executionOutput,
		ExecutionError:  executionError,
		TestResults:     testResults,
		TestScore:       testScore,
	}, nil
}

// GetPoolStatus returns the status of code execution container pools
func (s *GradingService) GetPoolStatus() map[string]interface{} {
	return s.executor.GetPoolStatus()
}

// ============================================
// Helper Methods
// ============================================

// shouldRequestAIFeedback determines if AI evaluation should be performed based on exercise configuration.
// Returns false when AI weight is 0, allowing exercises to skip AI evaluation entirely.
func (s *GradingService) shouldRequestAIFeedback(exercise *models.Exercise) bool {
	rubric := exercise.Content.GradingRubric

	// Skip AI if weight is explicitly set to 0
	if rubric.AIWeight == 0 {
		return false
	}

	// Skip AI if there's no AI evaluation prompt configured
	if exercise.Content.AIEvaluationPrompt == "" && len(rubric.AICriteria) == 0 {
		return false
	}

	return true
}

// runTestCases executes all test cases and returns results
func (s *GradingService) runTestCases(ctx context.Context, submittedCode, language string, testCases []models.TestCase, requestID string) (*models.TestResults, string, *string, int) {
	// For Rust, run all tests together in one compilation
	if strings.ToLower(language) == "rust" {
		return s.runRustTests(ctx, submittedCode, testCases, requestID)
	}

	// For other languages, run tests individually
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
		execResult := s.executor.ExecuteCodeWithContext(ctx, codeWithTest, language)
		duration := int(time.Since(startTime).Milliseconds())
		totalExecutionTime += duration

		testResult := models.TestResult{
			Name:       testCase.Name,
			DurationMs: &duration,
		}

		if execResult.Error == "" {
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
func (s *GradingService) runRustTests(ctx context.Context, submittedCode string, testCases []models.TestCase, requestID string) (*models.TestResults, string, *string, int) {
	results := &models.TestResults{
		TotalTests:  len(testCases),
		PassedTests: 0,
		FailedTests: 0,
		Results:     make([]models.TestResult, 0, len(testCases)),
	}

	// Build complete Rust test file
	var codeBuilder strings.Builder
	codeBuilder.WriteString(submittedCode)
	codeBuilder.WriteString("\n\n#[cfg(test)]\nmod tests {\n    #[allow(unused_imports)]\n    use super::*;\n\n")

	// Add all test cases
	for _, testCase := range testCases {
		testCode := testCase.Code
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
	execResult := s.executor.ExecuteCodeWithContext(ctx, completeCode, "rust")
	totalExecutionTime := int(time.Since(startTime).Milliseconds())

	// Parse test output
	if execResult.Error != "" {
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
func (s *GradingService) evaluateWithAI(ctx context.Context, exercise *models.Exercise, code, language string, testResults *models.TestResults, requestID string) (*models.AIFeedback, int, error) {
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

	// Get AI feedback from evaluator
	domainFeedback, err := s.evaluator.CheckSubmissionWithContext(ctx, prompt, code, language)
	if err != nil {
		return nil, 0, fmt.Errorf("AI evaluation failed: %w", err)
	}

	// Convert to models.AIFeedback
	feedback := domainFeedback.ToAIFeedback()

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

STARTER CODE PROVIDED:
%s

SUBMITTED CODE:
%s

CRITICAL: Compare the submitted code against the starter code. Evaluate:

1. LEARNING PROGRESS:
   - What changes did the student make from the starter code?
   - Are the changes substantial and demonstrate understanding?
   - Did the student just add comments or trivial modifications?

2. SOLUTION QUALITY:
   - Code correctness and logic
   - Code quality and style
   - Efficiency and best practices

3. EFFORT ASSESSMENT:
   - If changes are minimal or cosmetic only, assign a low score (0-20)
   - If solution shows genuine problem-solving, score appropriately
   - Even if tests pass, minimal effort should receive minimal credit

4. FEEDBACK:
   - Provide constructive, educational feedback
   - If submission lacks effort, explain what's missing
   - Suggest specific improvements

Be strict about effort: unchanged or trivially changed code should receive very low scores, even if tests pass.`,
		exercise.Title,
		exercise.Description,
		exercise.Difficulty,
		language,
		testResults.PassedTests,
		testResults.TotalTests,
		exercise.Content.StarterCode,
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
		"{starter_code}", exercise.Content.StarterCode,
	)
	return replacer.Replace(prompt)
}

// deriveAIScore calculates a score based on AI feedback and test results
func deriveAIScore(feedback string, testResults *models.TestResults) int {
	// Simple heuristic: if all tests pass, base score is 80-100
	// If some tests fail, base score is 40-80

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
func (s *GradingService) saveSubmission(
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
	}

	// Save to database
	if err := s.submissionRepo.Create(ctx, submission); err != nil {
		return nil, fmt.Errorf("failed to create submission: %w", err)
	}

	return submission, nil
}

// updateUserProgress updates user progress and mastery tracking
func (s *GradingService) updateUserProgress(ctx context.Context, userID, nodeID string, score int, requestID string) error {
	// Delegate to progress service
	return s.progressService.UpdateAfterSubmission(ctx, userID, nodeID, score, requestID)
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
