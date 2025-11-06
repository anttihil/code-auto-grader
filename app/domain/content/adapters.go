package content

import (
	"context"

	"github.com/anttihil/code-auto-grader/models"
	"github.com/anttihil/code-auto-grader/services"
	"github.com/google/uuid"
)

// GeminiAdapter adapts services.GeminiService to the ExerciseGenerator interface
type GeminiAdapter struct {
	gemini *services.GeminiService
}

// NewGeminiAdapter creates a new Gemini adapter
func NewGeminiAdapter(gemini *services.GeminiService) *GeminiAdapter {
	return &GeminiAdapter{gemini: gemini}
}

// GenerateExerciseWithContext implements the ExerciseGenerator interface
func (a *GeminiAdapter) GenerateExerciseWithContext(ctx context.Context, difficulty, topic, language string) (*models.Exercise, error) {
	// Call the underlying service
	serviceExercise, err := a.gemini.GenerateExerciseWithContext(ctx, difficulty, topic, language)
	if err != nil {
		return nil, err
	}

	// Convert services.Exercise to models.Exercise
	// Parse test cases
	testCases := make([]models.TestCase, 0, len(serviceExercise.TestCases))
	for i, tc := range serviceExercise.TestCases {
		testCases = append(testCases, models.TestCase{
			Name:    tc["name"],
			Code:    tc["code"],
			Visible: true, // Default to visible for generated exercises
		})
		// Give default name if missing
		if testCases[i].Name == "" {
			testCases[i].Name = "test_case_" + string(rune(i+1))
		}
	}

	// Parse hints
	hints := make([]models.Hint, 0, len(serviceExercise.Hints))
	for i, hint := range serviceExercise.Hints {
		hints = append(hints, models.Hint{
			Order:   i + 1,
			Title:   "Hint " + string(rune(i+1)),
			Content: hint,
		})
	}

	// Create exercise content
	content := models.ExerciseContent{
		StarterCode:       serviceExercise.StarterCode,
		ReferenceSolution: "", // Not provided by AI generation
		TestCases:         testCases,
		Hints:             hints,
		GradingRubric: models.GradingRubric{
			TestWeight: 60,
			AIWeight:   40,
			AICriteria: []models.AICriterion{},
		},
		AIEvaluationPrompt: "", // Use default
		CommonMistakes:     []models.CommonMistake{},
		ExpectedOutcomes:   []string{serviceExercise.ExpectedOutput},
	}

	// Create full exercise model
	exercise := &models.Exercise{
		ID:               uuid.New(),
		NodeID:           uuid.Nil, // Will be set when saved to database
		Title:            serviceExercise.Title,
		Slug:             "",       // Will be generated when saved
		Description:      serviceExercise.Description,
		Difficulty:       serviceExercise.Difficulty,
		EstimatedMinutes: nil,
		OrderIndex:       0,
		Content:          content,
	}

	return exercise, nil
}
