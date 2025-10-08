package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

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
		return nil
	}

	ctx := context.Background()
	config := &genai.ClientConfig{
		APIKey: apiKey,
	}
	client, err := genai.NewClient(ctx, config)
	if err != nil {
		return nil
	}

	models := client.Models

	return &GeminiService{
		client: client,
		models: models,
	}
}

// GenerateExercise generates a coding exercise using Gemini AI
func (s *GeminiService) GenerateExercise(difficulty, topic, language string) (Exercise, error) {
	if s == nil {
		return Exercise{}, fmt.Errorf("Gemini service not initialized (missing API key)")
	}

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

	ctx := context.Background()

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

	resp, err := s.models.GenerateContent(ctx, "gemini-2.5", genai.Text(prompt), &genai.GenerateContentConfig{
		MaxOutputTokens:  1000,
		Temperature:      &temperature,
		ResponseMIMEType: "application/json",
		ResponseSchema:   responseSchema,
	})
	if err != nil {
		return Exercise{}, fmt.Errorf("failed to generate exercise: %w", err)
	}

	// Extract and parse JSON from response
	text := resp.Text()
	var exercise Exercise
	if err := json.Unmarshal([]byte(text), &exercise); err != nil {
		return Exercise{}, fmt.Errorf("failed to parse exercise JSON: %w", err)
	}

	return exercise, nil
}

// CheckSubmission checks submitted code using Gemini AI
func (s *GeminiService) CheckSubmission(exercise, submittedCode, language string) (Feedback, error) {
	if s == nil {
		return Feedback{}, fmt.Errorf("Gemini service not initialized (missing API key)")
	}

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

	ctx := context.Background()

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

	resp, err := s.models.GenerateContent(ctx, "gemini-2.5-flash", genai.Text(prompt), &genai.GenerateContentConfig{
		MaxOutputTokens:  1000,
		ResponseMIMEType: "application/json",
		ResponseSchema:   responseSchema,
	})
	if err != nil {
		return Feedback{}, fmt.Errorf("failed to check submission: %w", err)
	}

	// Extract and parse JSON from response
	text := resp.Text()
	var feedback Feedback
	if err := json.Unmarshal([]byte(text), &feedback); err != nil {
		return Feedback{}, fmt.Errorf("failed to parse feedback JSON: %w", err)
	}

	return feedback, nil
}
