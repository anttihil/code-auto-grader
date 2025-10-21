package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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
		log.Println("WARNING: GEMINI_API_KEY environment variable is not set")
		return nil
	}

	log.Printf("Initializing Gemini service with API key (length: %d characters)", len(apiKey))

	ctx := context.Background()
	config := &genai.ClientConfig{
		APIKey: apiKey,
	}
	client, err := genai.NewClient(ctx, config)
	if err != nil {
		log.Printf("ERROR: Failed to create Gemini client: %v", err)
		return nil
	}

	models := client.Models

	log.Println("Gemini service client created successfully")
	return &GeminiService{
		client: client,
		models: models,
	}
}

// GenerateExercise generates a coding exercise using Gemini AI
func (s *GeminiService) GenerateExercise(difficulty, topic, language string) (Exercise, error) {
	if s == nil {
		log.Println("ERROR: GenerateExercise called but GeminiService is nil")
		return Exercise{}, fmt.Errorf("Gemini service not initialized (missing API key)")
	}

	log.Printf("GenerateExercise called - Difficulty: %s, Topic: %s, Language: %s", difficulty, topic, language)

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

	log.Println("Calling Gemini API for exercise generation...")
	resp, err := s.models.GenerateContent(ctx, "gemini-2.5-flash", genai.Text(prompt), &genai.GenerateContentConfig{
		MaxOutputTokens:  20000,
		Temperature:      &temperature,
		ResponseMIMEType: "application/json",
		ResponseSchema:   responseSchema,
	})
	if err != nil {
		log.Printf("ERROR: Gemini API call failed for exercise generation: %v", err)
		return Exercise{}, fmt.Errorf("failed to generate exercise: %w", err)
	}

	log.Println("Gemini API call successful, checking response...")

	// Check if response is nil or has no candidates
	if resp == nil {
		log.Println("ERROR: Received nil response from Gemini API")
		return Exercise{}, fmt.Errorf("received nil response from Gemini API")
	}

	// Log response details for debugging
	log.Printf("Response candidates count: %d", len(resp.Candidates))
	if len(resp.Candidates) == 0 {
		log.Println("ERROR: No candidates in response")
		return Exercise{}, fmt.Errorf("no candidates in Gemini response - content may have been blocked or filtered")
	}

	// Check for content filtering or blocking
	if resp.Candidates[0].FinishReason != "" {
		log.Printf("Response finish reason: %s", resp.Candidates[0].FinishReason)
	}
	if resp.PromptFeedback != nil {
		log.Printf("Prompt feedback: %+v", resp.PromptFeedback)
	}

	// Extract and parse JSON from response
	text := resp.Text()
	log.Printf("Response text length: %d characters", len(text))

	if text == "" {
		log.Println("ERROR: Empty response text from Gemini API")
		log.Printf("Full response structure: %+v", resp)
		return Exercise{}, fmt.Errorf("empty response text from Gemini API")
	}

	var exercise Exercise
	if err := json.Unmarshal([]byte(text), &exercise); err != nil {
		log.Printf("ERROR: Failed to parse exercise JSON: %v. Response text: %s", err, text)
		return Exercise{}, fmt.Errorf("failed to parse exercise JSON: %w. %s", err, text)
	}

	log.Printf("GenerateExercise completed successfully. Title: %s", exercise.Title)
	return exercise, nil
}

// CheckSubmission checks submitted code using Gemini AI
func (s *GeminiService) CheckSubmission(exercise, submittedCode, language string) (Feedback, error) {
	if s == nil {
		log.Println("ERROR: CheckSubmission called but GeminiService is nil")
		return Feedback{}, fmt.Errorf("Gemini service not initialized (missing API key)")
	}

	log.Printf("CheckSubmission called for language: %s, exercise length: %d, code length: %d",
		language, len(exercise), len(submittedCode))

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

	log.Println("Calling Gemini API for code submission feedback...")
	resp, err := s.models.GenerateContent(ctx, "gemini-2.5-flash", genai.Text(prompt), &genai.GenerateContentConfig{
		MaxOutputTokens:  3000,
		ResponseMIMEType: "application/json",
		ResponseSchema:   responseSchema,
	})
	if err != nil {
		log.Printf("ERROR: Gemini API call failed: %v", err)
		return Feedback{}, fmt.Errorf("failed to check submission: %w", err)
	}

	log.Println("Gemini API call successful, checking response...")

	// Check if response is nil or has no candidates
	if resp == nil {
		log.Println("ERROR: Received nil response from Gemini API")
		return Feedback{}, fmt.Errorf("received nil response from Gemini API")
	}

	// Log response details for debugging
	log.Printf("Response candidates count: %d", len(resp.Candidates))
	if len(resp.Candidates) == 0 {
		log.Println("ERROR: No candidates in response")
		return Feedback{}, fmt.Errorf("no candidates in Gemini response - content may have been blocked or filtered")
	}

	// Check for content filtering or blocking
	if resp.Candidates[0].FinishReason != "" {
		log.Printf("Response finish reason: %s", resp.Candidates[0].FinishReason)
		if resp.Candidates[0].FinishReason == "MAX_TOKENS" {
			log.Println("WARNING: Response was truncated due to MAX_TOKENS limit. Consider increasing MaxOutputTokens.")
		}
	}
	if resp.PromptFeedback != nil {
		log.Printf("Prompt feedback: %+v", resp.PromptFeedback)
	}

	// Extract and parse JSON from response
	text := resp.Text()
	log.Printf("Response text length: %d characters", len(text))

	if text == "" {
		log.Println("ERROR: Empty response text from Gemini API")
		log.Printf("Full response structure: %+v", resp)
		return Feedback{}, fmt.Errorf("empty response text from Gemini API")
	}

	var feedback Feedback
	if err := json.Unmarshal([]byte(text), &feedback); err != nil {
		log.Printf("ERROR: Failed to parse Gemini response JSON: %v. Response text: %s", err, text)
		return Feedback{}, fmt.Errorf("failed to parse feedback JSON: %w", err)
	}

	log.Printf("CheckSubmission completed successfully. Score: %d, Passed: %v", feedback.Score, feedback.Passed)
	return feedback, nil
}
