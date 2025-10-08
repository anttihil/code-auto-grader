package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"

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

Please provide the response in the following JSON format:
{
  "title": "Exercise title",
  "description": "Detailed description of the problem",
  "difficulty": "%s",
  "language": "%s",
  "starterCode": "Initial code template for the student",
  "hints": ["hint1", "hint2", "hint3"],
  "expectedOutput": "Description of expected output or behavior",
  "testCases": [
    {"input": "test input 1", "expectedOutput": "expected output 1"},
    {"input": "test input 2", "expectedOutput": "expected output 2"}
  ]
}

Make the exercise educational and engaging. Include clear instructions and appropriate test cases.`, 
		difficulty, topic, language, difficulty, language)

	var temperature float32 = 0.7
	if difficulty == "easy" {
		temperature = 0.5
	} else if difficulty == "hard" {
		temperature = 0.9
	}

	ctx := context.Background()
	resp, err := s.models.GenerateContent(ctx, "gemini-2.5", genai.Text(prompt), &genai.GenerateContentConfig{
		MaxOutputTokens: 1000,
		Temperature:     &temperature,
	})
	if err != nil {
		return Exercise{}, fmt.Errorf("failed to generate exercise: %w", err)
	}

	// Extract text from response
	var text string
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				text += fmt.Sprintf("%v", part)
			}
		}
	}

	// Try to extract JSON from the response
	re := regexp.MustCompile(`\{[\s\S]*\}`)
	jsonMatch := re.FindString(text)

	var exercise Exercise
	if jsonMatch != "" {
		if err := json.Unmarshal([]byte(jsonMatch), &exercise); err == nil {
			return exercise, nil
		}
	}

	// Fallback if JSON parsing fails
	return Exercise{
		Title:          "Generated Exercise",
		Description:    text,
		Difficulty:     difficulty,
		Language:       language,
		StarterCode:    fmt.Sprintf("// Write your %s code here\n", language),
		Hints:          []string{},
		ExpectedOutput: "See description for details",
		TestCases:      []map[string]string{},
	}, nil
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

Please provide feedback in the following JSON format:
{
  "passed": true/false,
  "score": 0-100,
  "feedback": "Overall feedback on the submission",
  "strengths": ["strength1", "strength2"],
  "improvements": ["improvement1", "improvement2"],
  "correctness": "Assessment of correctness",
  "codeQuality": "Assessment of code quality and style",
  "suggestions": "Specific suggestions for improvement"
}

Be constructive, educational, and encouraging in your feedback.`,
		exercise, language, "```"+language+"\n", submittedCode, "\n```")

	ctx := context.Background()

	resp, err := s.models.GenerateContent(ctx, "gemini-2.5-flash", genai.Text(prompt), &genai.GenerateContentConfig{
		MaxOutputTokens: 1000,
	})
	if err != nil {
		return Feedback{}, fmt.Errorf("failed to check submission: %w", err)
	}

	// Extract text from response
	var text string
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				text += fmt.Sprintf("%v", part)
			}
		}
	}

	// Try to extract JSON from the response
	re := regexp.MustCompile(`\{[\s\S]*\}`)
	jsonMatch := re.FindString(text)

	var feedback Feedback
	if jsonMatch != "" {
		if err := json.Unmarshal([]byte(jsonMatch), &feedback); err == nil {
			return feedback, nil
		}
	}

	// Fallback if JSON parsing fails
	return Feedback{
		Passed:       false,
		Score:        50,
		Feedback:     text,
		Strengths:    []string{},
		Improvements: []string{},
		Correctness:  "Unable to determine",
		CodeQuality:  "Unable to determine",
		Suggestions:  "Please review the feedback above",
	}, nil
}
