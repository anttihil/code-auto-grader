// Package grading provides the domain logic for code submission, execution, and evaluation.
//
// This package orchestrates the complete submission workflow:
//   - Validating submitted code against starter code to prevent gaming
//   - Executing code in isolated Docker containers with test cases
//   - Evaluating code quality with AI-powered feedback
//   - Calculating weighted scores (60% test results + 40% AI evaluation)
//   - Determining mastery achievement (≥80% threshold)
//   - Updating user progress upon successful completion
//
// The package is framework-agnostic and can be tested without HTTP or database dependencies.
// All external dependencies (code executor, AI evaluator, repositories) are accessed through
// interfaces, enabling easy mocking in tests.
//
// # Key Components
//
// Service: GradingService orchestrates the submission workflow
//   - SubmitCode: Complete submission with tests, AI evaluation, and scoring
//   - RunCode: Execute code without saving (for testing/preview)
//
// Scorer: Implements the weighted scoring algorithm
//   - CalculateScore: Combines test and AI scores with configurable weights
//   - DeterminePassStatus: Applies the 80% mastery threshold
//
// Similarity: Anti-gaming code similarity detection
//   - checkCodeSimilarity: Detects identical or trivial changes to starter code
//   - normalizeCode: Normalizes whitespace for fair comparison
//   - stripComments: Removes comments to detect substantive changes
//
// # Example Usage
//
//	service := grading.NewGradingService(executor, evaluator, submissionRepo, exerciseRepo, progressRepo)
//
//	result, err := service.SubmitCode(ctx, grading.SubmitCodeRequest{
//	    UserID:     "user-123",
//	    ExerciseID: "exercise-456",
//	    Code:       "function add(a, b) { return a + b; }",
//	    Language:   "javascript",
//	}, "request-id")
//
//	if err != nil {
//	    // Handle error (e.g., ErrIdenticalToStarterCode, ErrTrivialChanges)
//	}
//
//	fmt.Printf("Score: %d%%, Passed: %v\n", result.TotalScore, result.Passed)
//
// # Testing
//
// The package provides interfaces for all external dependencies, making it easy to test
// with mocks. See service_test.go for comprehensive examples of testing the domain logic
// without real executors, AI services, or databases.
package grading
