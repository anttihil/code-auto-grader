package exercises

import "github.com/anttihil/code-auto-grader/models"

// ExerciseView represents an exercise with submission history
type ExerciseView struct {
	Exercise     *models.Exercise
	Submissions  []*models.Submission
	BestScore    int
	AttemptCount int
}

// GenerateExerciseRequest represents a request to generate a new exercise using AI
type GenerateExerciseRequest struct {
	Difficulty string
	Topic      string
	Language   string
}

// ExerciseGenerationResult represents the result of AI exercise generation
type ExerciseGenerationResult struct {
	Title          string                   `json:"title"`
	Description    string                   `json:"description"`
	Difficulty     string                   `json:"difficulty"`
	Language       string                   `json:"language"`
	StarterCode    string                   `json:"starterCode"`
	Hints          []string                 `json:"hints"`
	ExpectedOutput string                   `json:"expectedOutput"`
	TestCases      []models.TestCase        `json:"testCases"`
}
