package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/anttihil/code-auto-grader/internal/logger"
	"google.golang.org/genai"
)

// GeminiService handles AI-powered exercise generation and code evaluation
type GeminiService struct {
	client *genai.Client
	models *genai.Models
}

// Exercise represents a coding exercise
type Exercise struct {
	Title          string              `json:"title"`
	Description    string              `json:"description"`
	Difficulty     string              `json:"difficulty"`
	Language       string              `json:"language"`
	StarterCode    string              `json:"starterCode"`
	Hints          []string            `json:"hints"`
	ExpectedOutput string              `json:"expectedOutput"`
	TestCases      []map[string]string `json:"testCases"`
}

// Feedback represents AI feedback on code submission
type Feedback struct {
	Passed       bool     `json:"passed"`
	Score        int      `json:"score"`
	Feedback     string   `json:"feedback"`
	Strengths    []string `json:"strengths"`
	Improvements []string `json:"improvements"`
	Correctness  string   `json:"correctness"`
	CodeQuality  string   `json:"codeQuality"`
	Suggestions  string   `json:"suggestions"`
}

// NewGeminiService creates a new Gemini AI service
func NewGeminiService() *GeminiService {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		logger.Logger.Warn("GEMINI_API_KEY environment variable is not set")
		return nil
	}

	logger.Logger.Info("initializing Gemini service",
		slog.Int("api_key_length", len(apiKey)),
	)

	ctx := context.Background()
	config := &genai.ClientConfig{
		APIKey: apiKey,
	}
	client, err := genai.NewClient(ctx, config)
	if err != nil {
		logger.Logger.Error("failed to create Gemini client",
			slog.String("error", err.Error()),
		)
		return nil
	}

	models := client.Models

	logger.Logger.Info("Gemini service client created successfully")
	return &GeminiService{
		client: client,
		models: models,
	}
}

// GenerateExerciseWithContext generates a coding exercise using Gemini AI with request context
func (s *GeminiService) GenerateExerciseWithContext(ctx context.Context, difficulty, topic, language string) (Exercise, error) {
	requestID := logger.GetRequestID(ctx)
	return s.generateExerciseInternal(ctx, difficulty, topic, language, requestID)
}

// GenerateExercise generates a coding exercise using Gemini AI (backwards compatibility)
func (s *GeminiService) GenerateExercise(difficulty, topic, language string) (Exercise, error) {
	return s.generateExerciseInternal(context.Background(), difficulty, topic, language, "unknown")
}

// generateExerciseInternal is the internal implementation with context support
func (s *GeminiService) generateExerciseInternal(ctx context.Context, difficulty, topic, language, requestID string) (Exercise, error) {
	if s == nil {
		logger.Logger.Error("GenerateExercise called but GeminiService is nil",
			slog.String("request_id", requestID),
		)
		return Exercise{}, fmt.Errorf("Gemini service not initialized (missing API key)")
	}

	start := time.Now()
	aiLogger := logger.WithAIContext(ctx, requestID, "generate_exercise")

	aiLogger.Info("generating exercise",
		slog.String("difficulty", difficulty),
		slog.String("topic", topic),
		slog.String("language", language),
	)

	prompt := fmt.Sprintf(`Generate a coding exercise with the following requirements:
- Difficulty: %s
- Topic: %s
- Programming Language: %s

Make the exercise educational and engaging. Include:
- A clear title and detailed description of the problem
- Starter code template for the student
- 3-5 helpful hints
- Description of expected output or behavior
- At least 2 test cases with input and expected output

Be creative and ensure the exercise is appropriate for the specified difficulty level.`,
		difficulty, topic, language)

	var temperature float32

	switch difficulty {
		case "easy":
			temperature = 0.5
		case "hard":
			temperature = 0.9
		default:
			temperature = 0.7
	}

	// Define response schema for structured output
	responseSchema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"title":          {Type: genai.TypeString},
			"description":    {Type: genai.TypeString},
			"difficulty":     {Type: genai.TypeString},
			"language":       {Type: genai.TypeString},
			"starterCode":    {Type: genai.TypeString},
			"hints":          {Type: genai.TypeArray, Items: &genai.Schema{Type: genai.TypeString}},
			"expectedOutput": {Type: genai.TypeString},
			"testCases": {
				Type: genai.TypeArray,
				Items: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"input":          {Type: genai.TypeString},
						"expectedOutput": {Type: genai.TypeString},
					},
					Required: []string{"input", "expectedOutput"},
				},
			},
		},
		Required: []string{"title", "description", "difficulty", "language", "starterCode", "hints", "expectedOutput", "testCases"},
	}

	aiLogger.Debug("calling Gemini API for exercise generation")
	resp, err := s.models.GenerateContent(ctx, "gemini-2.5-flash", genai.Text(prompt), &genai.GenerateContentConfig{
		MaxOutputTokens:  20000,
		Temperature:      &temperature,
		ResponseMIMEType: "application/json",
		ResponseSchema:   responseSchema,
	})
	latency := time.Since(start)

	if err != nil {
		aiLogger.Error("Gemini API call failed for exercise generation",
			slog.Duration("latency_ms", latency),
			slog.String("error", err.Error()),
		)
		return Exercise{}, fmt.Errorf("failed to generate exercise: %w", err)
	}

	aiLogger.Info("Gemini API call successful",
		slog.Duration("latency_ms", latency),
	)

	// Check if response is nil or has no candidates
	if resp == nil {
		aiLogger.Error("received nil response from Gemini API")
		return Exercise{}, fmt.Errorf("received nil response from Gemini API")
	}

	// Check response details
	if len(resp.Candidates) == 0 {
		aiLogger.Warn("no candidates in Gemini response",
			slog.Any("prompt_feedback", resp.PromptFeedback),
		)
		return Exercise{}, fmt.Errorf("no candidates in Gemini response - content may have been blocked or filtered")
	}

	// Check for content filtering or blocking
	if resp.Candidates[0].FinishReason != "" {
		aiLogger.Debug("response finish reason",
			slog.String("finish_reason", string(resp.Candidates[0].FinishReason)),
		)
	}

	// Extract and parse JSON from response
	text := resp.Text()

	if text == "" {
		aiLogger.Error("empty response text from Gemini API")
		return Exercise{}, fmt.Errorf("empty response text from Gemini API")
	}

	aiLogger.Debug("received response",
		slog.Int("response_length", len(text)),
	)

	var exercise Exercise
	if err := json.Unmarshal([]byte(text), &exercise); err != nil {
		aiLogger.Error("failed to parse exercise JSON",
			slog.String("error", err.Error()),
		)
		return Exercise{}, fmt.Errorf("failed to parse exercise JSON: %w. %s", err, text)
	}

	aiLogger.Info("exercise generated successfully",
		slog.String("title", exercise.Title),
	)
	return exercise, nil
}

// CheckSubmissionWithContext checks submitted code using Gemini AI with request context
func (s *GeminiService) CheckSubmissionWithContext(ctx context.Context, exercise, submittedCode, language string) (Feedback, error) {
	requestID := logger.GetRequestID(ctx)
	return s.checkSubmissionInternal(ctx, exercise, submittedCode, language, requestID)
}

// CheckSubmission checks submitted code using Gemini AI (backwards compatibility)
func (s *GeminiService) CheckSubmission(exercise, submittedCode, language string) (Feedback, error) {
	return s.checkSubmissionInternal(context.Background(), exercise, submittedCode, language, "unknown")
}

// checkSubmissionInternal is the internal implementation with context support
func (s *GeminiService) checkSubmissionInternal(ctx context.Context, exercise, submittedCode, language, requestID string) (Feedback, error) {
	if s == nil {
		logger.Logger.Error("CheckSubmission called but GeminiService is nil",
			slog.String("request_id", requestID),
		)
		return Feedback{}, fmt.Errorf("Gemini service not initialized (missing API key)")
	}

	start := time.Now()
	aiLogger := logger.WithAIContext(ctx, requestID, "check_submission")

	aiLogger.Info("checking submission",
		slog.String("language", language),
		slog.Int("exercise_length", len(exercise)),
		slog.Int("code_length", len(submittedCode)),
	)

	prompt := fmt.Sprintf(`You are a code reviewer and teacher. Review the following code submission:

EXERCISE:
%s

PROGRAMMING LANGUAGE: %s

SUBMITTED CODE:
%s%s%s

Provide a thorough review that includes:
- Whether the code passes the exercise requirements (passed: true/false)
- A numerical score from 0-100
- Overall feedback on the submission
- List of strengths in the code
- List of areas for improvement
- Assessment of correctness
- Assessment of code quality and style
- Specific suggestions for improvement

Be constructive, educational, and encouraging in your feedback.`,
		exercise, language, "```"+language+"\n", submittedCode, "\n```")

	// Define response schema for structured output
	responseSchema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"passed":       {Type: genai.TypeBoolean},
			"score":        {Type: genai.TypeInteger},
			"feedback":     {Type: genai.TypeString},
			"strengths":    {Type: genai.TypeArray, Items: &genai.Schema{Type: genai.TypeString}},
			"improvements": {Type: genai.TypeArray, Items: &genai.Schema{Type: genai.TypeString}},
			"correctness":  {Type: genai.TypeString},
			"codeQuality":  {Type: genai.TypeString},
			"suggestions":  {Type: genai.TypeString},
		},
		Required: []string{"passed", "score", "feedback", "strengths", "improvements", "correctness", "codeQuality", "suggestions"},
	}

	aiLogger.Debug("calling Gemini API for code submission feedback")
	resp, err := s.models.GenerateContent(ctx, "gemini-2.5-flash", genai.Text(prompt), &genai.GenerateContentConfig{
		MaxOutputTokens:  3000,
		ResponseMIMEType: "application/json",
		ResponseSchema:   responseSchema,
	})
	latency := time.Since(start)

	if err != nil {
		aiLogger.Error("Gemini API call failed",
			slog.Duration("latency_ms", latency),
			slog.String("error", err.Error()),
		)
		return Feedback{}, fmt.Errorf("failed to check submission: %w", err)
	}

	aiLogger.Info("Gemini API call successful",
		slog.Duration("latency_ms", latency),
	)

	// Check if response is nil or has no candidates
	if resp == nil {
		aiLogger.Error("received nil response from Gemini API")
		return Feedback{}, fmt.Errorf("received nil response from Gemini API")
	}

	// Check response details
	if len(resp.Candidates) == 0 {
		aiLogger.Warn("no candidates in Gemini response",
			slog.Any("prompt_feedback", resp.PromptFeedback),
		)
		return Feedback{}, fmt.Errorf("no candidates in Gemini response - content may have been blocked or filtered")
	}

	// Check for content filtering or blocking
	if resp.Candidates[0].FinishReason != "" {
		finishReason := string(resp.Candidates[0].FinishReason)
		aiLogger.Debug("response finish reason",
			slog.String("finish_reason", finishReason),
		)
		if finishReason == "MAX_TOKENS" {
			aiLogger.Warn("response truncated due to MAX_TOKENS limit")
		}
	}

	// Extract and parse JSON from response
	text := resp.Text()

	if text == "" {
		aiLogger.Error("empty response text from Gemini API")
		return Feedback{}, fmt.Errorf("empty response text from Gemini API")
	}

	aiLogger.Debug("received response",
		slog.Int("response_length", len(text)),
	)

	var feedback Feedback
	if err := json.Unmarshal([]byte(text), &feedback); err != nil {
		aiLogger.Error("failed to parse feedback JSON",
			slog.String("error", err.Error()),
		)
		return Feedback{}, fmt.Errorf("failed to parse feedback JSON: %w", err)
	}

	aiLogger.Info("submission checked successfully",
		slog.Bool("passed", feedback.Passed),
		slog.Int("score", feedback.Score),
	)
	return feedback, nil
}
