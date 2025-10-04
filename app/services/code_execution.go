package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// CodeExecutionService handles secure code execution
type CodeExecutionService struct {
	tempDir string
	timeout time.Duration
}

// ExecutionResult represents the result of code execution
type ExecutionResult struct {
	Success       bool   `json:"success"`
	Output        string `json:"output"`
	Error         string `json:"error"`
	ExecutionTime string `json:"executionTime"`
}

// NewCodeExecutionService creates a new code execution service
func NewCodeExecutionService() *CodeExecutionService {
	return &CodeExecutionService{
		tempDir: "/tmp/code-executions",
		timeout: 5 * time.Second,
	}
}

// ExecuteCode executes code based on the language
func (s *CodeExecutionService) ExecuteCode(code, language string) ExecutionResult {
	// Initialize temp directory
	if err := os.MkdirAll(s.tempDir, 0755); err != nil {
		return ExecutionResult{
			Success:       false,
			Output:        "",
			Error:         fmt.Sprintf("Failed to create temp directory: %v", err),
			ExecutionTime: "error",
		}
	}

	// Validate code
	if strings.TrimSpace(code) == "" {
		return ExecutionResult{
			Success:       false,
			Output:        "",
			Error:         "No code provided",
			ExecutionTime: "error",
		}
	}

	// Check for dangerous patterns
	if err := s.checkDangerousPatterns(code); err != nil {
		return ExecutionResult{
			Success:       false,
			Output:        "",
			Error:         err.Error(),
			ExecutionTime: "blocked",
		}
	}

	// Execute based on language
	language = strings.ToLower(language)
	switch language {
	case "javascript", "js":
		return s.executeJavaScript(code)
	case "python", "py":
		return s.executePython(code)
	default:
		return ExecutionResult{
			Success:       false,
			Output:        "",
			Error:         fmt.Sprintf("Language %s is not supported yet", language),
			ExecutionTime: "error",
		}
	}
}

// checkDangerousPatterns checks for potentially dangerous code patterns
func (s *CodeExecutionService) checkDangerousPatterns(code string) error {
	dangerousPatterns := []struct {
		pattern *regexp.Regexp
		desc    string
	}{
		{regexp.MustCompile(`require\s*\(\s*['"]child_process['"]\s*\)`), "child_process"},
		{regexp.MustCompile(`require\s*\(\s*['"]fs['"]\s*\)`), "fs"},
		{regexp.MustCompile(`require\s*\(\s*['"]net['"]\s*\)`), "net"},
		{regexp.MustCompile(`import\s+os`), "os"},
		{regexp.MustCompile(`import\s+subprocess`), "subprocess"},
		{regexp.MustCompile(`import\s+socket`), "socket"},
		{regexp.MustCompile(`eval\s*\(`), "eval"},
		{regexp.MustCompile(`exec\s*\(`), "exec"},
		{regexp.MustCompile(`system\s*\(`), "system"},
	}

	for _, dp := range dangerousPatterns {
		if dp.pattern.MatchString(code) {
			return fmt.Errorf("Code contains potentially dangerous operations")
		}
	}

	return nil
}

// executeJavaScript executes JavaScript code
func (s *CodeExecutionService) executeJavaScript(code string) ExecutionResult {
	executionID := generateRandomID()
	filePath := filepath.Join(s.tempDir, executionID+".js")

	// Write code to file
	if err := os.WriteFile(filePath, []byte(code), 0644); err != nil {
		return ExecutionResult{
			Success:       false,
			Output:        "",
			Error:         fmt.Sprintf("Failed to write code file: %v", err),
			ExecutionTime: "error",
		}
	}
	defer os.Remove(filePath)

	// Execute with Node.js
	nodePath, err := exec.LookPath("node")
	if err != nil {
		nodePath = "/usr/local/bin/node"
	}

	cmd := exec.Command("timeout", "5s", nodePath, "--max-old-space-size=128", filePath)
	cmd.Env = append(os.Environ(), "NODE_ENV=sandbox")

	output, err := cmd.CombinedOutput()
	if err != nil {
		// Check if it was a timeout
		if strings.Contains(err.Error(), "killed") || strings.Contains(string(output), "timeout") {
			return ExecutionResult{
				Success:       false,
				Output:        "",
				Error:         "Execution timed out",
				ExecutionTime: "timeout",
			}
		}
		return ExecutionResult{
			Success:       false,
			Output:        "",
			Error:         string(output),
			ExecutionTime: "error",
		}
	}

	return ExecutionResult{
		Success:       true,
		Output:        string(output),
		Error:         "",
		ExecutionTime: "completed",
	}
}

// executePython executes Python code
func (s *CodeExecutionService) executePython(code string) ExecutionResult {
	executionID := generateRandomID()
	filePath := filepath.Join(s.tempDir, executionID+".py")

	// Write code to file
	if err := os.WriteFile(filePath, []byte(code), 0644); err != nil {
		return ExecutionResult{
			Success:       false,
			Output:        "",
			Error:         fmt.Sprintf("Failed to write code file: %v", err),
			ExecutionTime: "error",
		}
	}
	defer os.Remove(filePath)

	// Execute with Python
	cmd := exec.Command("timeout", "5s", "python3", "-u", filePath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		// Check if it was a timeout
		if strings.Contains(err.Error(), "killed") || strings.Contains(string(output), "timeout") {
			return ExecutionResult{
				Success:       false,
				Output:        "",
				Error:         "Execution timed out",
				ExecutionTime: "timeout",
			}
		}
		return ExecutionResult{
			Success:       false,
			Output:        "",
			Error:         string(output),
			ExecutionTime: "error",
		}
	}

	return ExecutionResult{
		Success:       true,
		Output:        string(output),
		Error:         "",
		ExecutionTime: "completed",
	}
}

// generateRandomID generates a random hex ID
func generateRandomID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
