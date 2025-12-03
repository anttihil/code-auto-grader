// Package content provides the domain logic for exercise management and AI-powered generation.
//
// This package handles exercise-related operations:
//   - Retrieving exercises with submission history and best scores
//   - Generating new exercises using AI based on topic, difficulty, and language
//   - Providing sample exercises for demonstrations
//   - Enriching exercise data with user-specific submission information
//
// The package integrates with an AI service (via adapter) to generate contextually
// appropriate coding exercises with starter code, test cases, and evaluation rubrics.
//
// # Key Components
//
// Service: ContentService provides exercise management functionality
//   - GetExercise: Returns an exercise with submission history for a user
//   - GenerateExercise: Creates a new exercise using AI based on parameters
//   - GetSampleExercise: Returns a pre-generated sample exercise
//
// Models: Domain types for exercise content
//   - ExerciseView: Exercise with enriched submission data
//   - ExerciseGenerationResult: AI-generated exercise with all components
//   - GenerateExerciseRequest: Parameters for exercise generation
//
// Adapters: Bridge to external AI services
//   - ExerciseGenerator: Interface for AI-powered exercise generation
//   - GeminiGeneratorAdapter: Concrete implementation using Gemini AI
//
// # Example Usage
//
//	service := content.NewContentService(exerciseRepo, submissionRepo, generator)
//
//	// Retrieve existing exercise with submission history
//	exercise, err := service.GetExercise(ctx, "exercise-id", "user-id", "request-id")
//	if err != nil {
//	    // Handle error (e.g., ErrExerciseNotFound)
//	}
//
//	fmt.Printf("Exercise: %s (Best score: %d)\n", exercise.Title, exercise.BestScore)
//
//	// Generate new exercise with AI
//	generated, err := service.GenerateExercise(ctx, content.GenerateExerciseRequest{
//	    Topic:      "array manipulation",
//	    Difficulty: "medium",
//	    Language:   "javascript",
//	}, "request-id")
//
//	if err != nil {
//	    // Handle error
//	}
//
//	fmt.Printf("Generated: %s\n", generated.Title)
//	fmt.Printf("Starter code: %s\n", generated.StarterCode)
//
// # AI-Powered Generation
//
// The package uses an AI service to generate complete exercises including:
//   - Title and description
//   - Learning objectives
//   - Starter code with TODO markers
//   - Multiple test cases (visible and hidden)
//   - Progressive hints for struggling students
//   - Evaluation rubric for AI assessment
//
// Generation parameters allow customization:
//   - Topic: The programming concept to teach (e.g., "recursion", "closures")
//   - Difficulty: "easy", "medium", or "hard"
//   - Language: "javascript", "python", "rust", or "go"
//
// # Testing
//
// The package uses interfaces for all external dependencies (repositories and AI generator),
// enabling comprehensive testing with mocks. See service_test.go for examples of testing
// exercise retrieval and generation without real databases or AI API calls.
package exercises
