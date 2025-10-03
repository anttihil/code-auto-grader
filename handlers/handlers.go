package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/anttihil/code-auto-grader/services"
)

// HealthCheck returns the health status of the API
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetSampleExercise returns a sample coding exercise
func GetSampleExercise(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sampleExercise := map[string]interface{}{
		"title":       "Sum Two Numbers",
		"description": "Write a function that takes two numbers as parameters and returns their sum.",
		"difficulty":  "beginner",
		"language":    "javascript",
		"starterCode": "// Write a function called sum that takes two parameters: a and b\n// Return the sum of a and b\n\nfunction sum(a, b) {\n  // Your code here\n}\n\n// Test your function\nconsole.log(sum(2, 3)); // Should output: 5\nconsole.log(sum(10, 20)); // Should output: 30",
		"hints": []string{
			"Use the + operator to add two numbers",
			"Remember to return the result",
			"The function should work with any two numbers",
		},
		"expectedOutput": "The function should return the sum of two numbers",
		"testCases": []map[string]string{
			{"input": "sum(2, 3)", "expectedOutput": "5"},
			{"input": "sum(10, 20)", "expectedOutput": "30"},
			{"input": "sum(-5, 5)", "expectedOutput": "0"},
		},
	}

	response := map[string]interface{}{
		"success":  true,
		"exercise": sampleExercise,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GenerateExercise generates a new coding exercise using AI
func GenerateExercise(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Difficulty string `json:"difficulty"`
		Topic      string `json:"topic"`
		Language   string `json:"language"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}

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

	gemini := services.NewGeminiService()
	exercise, err := gemini.GenerateExercise(req.Difficulty, req.Topic, req.Language)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Failed to generate exercise",
			"message": err.Error(),
		})
		return
	}

	response := map[string]interface{}{
		"success":  true,
		"exercise": exercise,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RunCode executes code without evaluation
func RunCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Code     string `json:"code"`
		Language string `json:"language"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}

	if req.Code == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Code is required",
		})
		return
	}

	if req.Language == "" {
		req.Language = "javascript"
	}

	executor := services.NewCodeExecutionService()
	result := executor.ExecuteCode(req.Code, req.Language)

	response := map[string]interface{}{
		"success":         true,
		"executionResult": result,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// SubmitCode submits code for AI evaluation
func SubmitCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Code     string `json:"code"`
		Exercise string `json:"exercise"`
		Language string `json:"language"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}

	if req.Code == "" || req.Exercise == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Code and exercise are required",
		})
		return
	}

	if req.Language == "" {
		req.Language = "javascript"
	}

	// Execute the code
	executor := services.NewCodeExecutionService()
	executionResult := executor.ExecuteCode(req.Code, req.Language)

	// Get AI feedback
	gemini := services.NewGeminiService()
	feedback, err := gemini.CheckSubmission(req.Exercise, req.Code, req.Language)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Failed to process submission",
			"message": err.Error(),
		})
		return
	}

	response := map[string]interface{}{
		"success":         true,
		"executionResult": executionResult,
		"feedback":        feedback,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
