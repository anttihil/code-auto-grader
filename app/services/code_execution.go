package services

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// CodeExecutionService handles secure code execution in Docker containers
type CodeExecutionService struct {
	tempDir       string
	timeout       time.Duration
	jsDockerImage string
	pyDockerImage string
}

// ExecutionResult represents the result of code execution
type ExecutionResult struct {
	Success       bool   `json:"success"`
	Output        string `json:"output"`
	Error         string `json:"error"`
	ExecutionTime string `json:"executionTime"`
}

// NewCodeExecutionService creates a new code execution service
func NewCodeExecutionService() (*CodeExecutionService, error) {
	service := &CodeExecutionService{
		tempDir:       "/tmp/code-executions",
		timeout:       5 * time.Second,
		jsDockerImage: "code-grader-javascript",
		pyDockerImage: "code-grader-python",
	}

	if err := service.initializeDocker(); err != nil {
		return nil, fmt.Errorf("failed to initialize Docker environment: %w", err)
	}
	return service, nil
}

// ExecuteCode executes code based on the language
func (s *CodeExecutionService) ExecuteCode(code, language string) ExecutionResult {
	if strings.TrimSpace(code) == "" {
		return ExecutionResult{
			Success:       false,
			Error:         "No code provided",
			ExecutionTime: "error",
		}
	}

	language = strings.ToLower(language)
	switch language {
	case "javascript", "js":
		return s.executeInDocker(code, s.jsDockerImage, "code.js")
	case "python", "py":
		return s.executeInDocker(code, s.pyDockerImage, "code.py")
	default:
		return ExecutionResult{
			Success:       false,
			Error:         fmt.Sprintf("Language %s is not supported", language),
			ExecutionTime: "error",
		}
	}
}

// initializeDocker builds Docker images for code execution
func (s *CodeExecutionService) initializeDocker() error {
	fmt.Println("Initializing Docker execution environment...")

	if err := s.buildDockerImage("javascript", s.jsDockerImage); err != nil {
		return fmt.Errorf("failed to build JavaScript Docker image: %w", err)
	}

	if err := s.buildDockerImage("python", s.pyDockerImage); err != nil {
		return fmt.Errorf("failed to build Python Docker image: %w", err)
	}

	fmt.Println("Docker execution environment ready")
	return nil
}

// buildDockerImage builds a Docker image from the Dockerfile
func (s *CodeExecutionService) buildDockerImage(language, imageName string) error {
	dockerfilePath := filepath.Join("docker", language+".Dockerfile")

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "build", "-f", dockerfilePath, "-t", imageName, ".")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker build failed: %v\nstdout: %s\nstderr: %s", err, stdout.String(), stderr.String())
	}

	fmt.Printf("Successfully built Docker image: %s\n", imageName)
	return nil
}

// executeInDocker executes code in a Docker container
func (s *CodeExecutionService) executeInDocker(code, image, filename string) ExecutionResult {
	// Verify Docker image exists
	if err := s.verifyDockerImage(image); err != nil {
		return ExecutionResult{
			Success:       false,
			Error:         fmt.Sprintf("Docker image not available: %v", err),
			ExecutionTime: "error",
		}
	}

	filePath := filepath.Join(s.tempDir, generateRandomID()+"_"+filename)

	if err := os.MkdirAll(s.tempDir, 0755); err != nil {
		return ExecutionResult{
			Success:       false,
			Error:         fmt.Sprintf("Failed to create temp directory: %v", err),
			ExecutionTime: "error",
		}
	}

	if err := os.WriteFile(filePath, []byte(code), 0644); err != nil {
		return ExecutionResult{
			Success:       false,
			Error:         fmt.Sprintf("Failed to write code file: %v", err),
			ExecutionTime: "error",
		}
	}
	defer os.Remove(filePath)

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "run",
		"--rm",
		"--network", "none",
		"--memory", "128m",
		"--cpus", "0.5",
		"--pids-limit", "50",
		"--read-only",
		"--cap-drop", "ALL",
		"--security-opt", "no-new-privileges",
		"-v", filePath+":/workspace/"+filename+":ro",
		image,
		filename,
	)

	output, err := cmd.CombinedOutput()

	if ctx.Err() == context.DeadlineExceeded {
		return ExecutionResult{
			Success:       false,
			Error:         "Execution timed out",
			ExecutionTime: "timeout",
		}
	}

	if err != nil {
		return ExecutionResult{
			Success:       false,
			Error:         string(output),
			ExecutionTime: "error",
		}
	}

	return ExecutionResult{
		Success:       true,
		Output:        string(output),
		ExecutionTime: "completed",
	}
}

// verifyDockerImage checks if a Docker image exists locally
func (s *CodeExecutionService) verifyDockerImage(imageName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "image", "inspect", imageName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("image '%s' not found - please ensure Docker images are built", imageName)
	}

	return nil
}

// generateRandomID generates a random hex ID
func generateRandomID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
