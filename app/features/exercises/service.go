// Package content contains business logic for exercise management and generation
package exercises

import (
	"context"
	"log/slog"

	"github.com/anttihil/code-auto-grader/internal/logger"
	"github.com/anttihil/code-auto-grader/models"
	"github.com/anttihil/code-auto-grader/repository"
)

// ExerciseGenerator defines the interface for AI-based exercise generation
type ExerciseGenerator interface {
	GenerateExerciseWithContext(ctx context.Context, difficulty, topic, language string) (*models.Exercise, error)
}

// ContentService handles exercise management and AI generation
type ContentService struct {
	generator      ExerciseGenerator
	exerciseRepo   repository.ExerciseRepository
	submissionRepo repository.SubmissionRepository
}

// NewContentService creates a new content service with required dependencies
func NewContentService(
	generator ExerciseGenerator,
	exerciseRepo repository.ExerciseRepository,
	submissionRepo repository.SubmissionRepository,
) *ContentService {
	if generator == nil {
		panic("exercise generator is required")
	}
	if exerciseRepo == nil {
		panic("exercise repository is required")
	}
	if submissionRepo == nil {
		panic("submission repository is required")
	}

	return &ContentService{
		generator:      generator,
		exerciseRepo:   exerciseRepo,
		submissionRepo: submissionRepo,
	}
}

// GetExercise fetches an exercise with submission history for a user
func (s *ContentService) GetExercise(ctx context.Context, exerciseID, userID, requestID string) (*ExerciseView, error) {
	logger.Logger.Info("fetching exercise",
		slog.String("request_id", requestID),
		slog.String("exercise_id", exerciseID),
		slog.String("user_id", userID),
	)

	// Fetch exercise
	exercise, err := s.exerciseRepo.GetByID(ctx, exerciseID)
	if err != nil {
		logger.Logger.Error("failed to fetch exercise",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	// Fetch submission history (limited to recent submissions)
	submissions, _, err := s.submissionRepo.GetByExercisePaginated(ctx, userID, exerciseID, 10, 0)
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

	logger.Logger.Info("exercise fetched successfully",
		slog.String("request_id", requestID),
		slog.Int("submission_count", len(submissions)),
	)

	return &ExerciseView{
		Exercise:     exercise,
		Submissions:  submissions,
		BestScore:    bestScore,
		AttemptCount: attemptCount,
	}, nil
}

// GenerateExercise generates a new coding exercise using AI
func (s *ContentService) GenerateExercise(ctx context.Context, req GenerateExerciseRequest, requestID string) (*ExerciseGenerationResult, error) {
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

	exercise, err := s.generator.GenerateExerciseWithContext(ctx, req.Difficulty, req.Topic, req.Language)
	if err != nil {
		logger.Logger.Error("failed to generate exercise",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	logger.Logger.Info("exercise generated successfully",
		slog.String("request_id", requestID),
	)

	// Convert to result format
	result := &ExerciseGenerationResult{
		Title:          exercise.Title,
		Description:    exercise.Description,
		Difficulty:     exercise.Difficulty,
		Language:       req.Language,
		StarterCode:    exercise.Content.StarterCode,
		Hints:          make([]string, 0),
		ExpectedOutput: "",
		TestCases:      exercise.Content.TestCases,
	}

	// Extract hints
	for _, hint := range exercise.Content.Hints {
		result.Hints = append(result.Hints, hint.Content)
	}

	// Extract expected outcomes
	if len(exercise.Content.ExpectedOutcomes) > 0 {
		result.ExpectedOutput = exercise.Content.ExpectedOutcomes[0]
	}

	return result, nil
}

// GetSampleExercise returns a sample coding exercise (for demo purposes)
func (s *ContentService) GetSampleExercise() *ExerciseGenerationResult {
	return &ExerciseGenerationResult{
		Title:       "Sum Two Numbers",
		Description: "Write a function that takes two numbers as parameters and returns their sum.",
		Difficulty:  "beginner",
		Language:    "javascript",
		StarterCode: "// Write a function called sum that takes two parameters: a and b\n// Return the sum of a and b\n\nfunction sum(a, b) {\n  // Your code here\n}\n\n// Test your function\nconsole.log(sum(2, 3)); // Should output: 5\nconsole.log(sum(10, 20)); // Should output: 30",
		Hints: []string{
			"Use the + operator to add two numbers",
			"Remember to return the result",
			"The function should work with any two numbers",
		},
		ExpectedOutput: "The function should return the sum of two numbers",
		TestCases: []models.TestCase{
			{Name: "test_2_plus_3", Code: "console.assert(sum(2, 3) === 5)", Visible: true},
			{Name: "test_10_plus_20", Code: "console.assert(sum(10, 20) === 30)", Visible: true},
			{Name: "test_negative", Code: "console.assert(sum(-5, 5) === 0)", Visible: true},
		},
	}
}
