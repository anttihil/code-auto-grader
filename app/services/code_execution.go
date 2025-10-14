package services

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	// DefaultMaxConcurrent is the default maximum number of concurrent executions
	DefaultMaxConcurrent = 15
	// DefaultQueueTimeout is the default timeout for waiting in queue
	DefaultQueueTimeout = 10 * time.Second
	// Container rotation policies
	MaxExecutionsPerContainer = 1000
	MaxContainerAge           = 1 * time.Hour
)

// CodeExecutionService handles secure code execution in Docker containers
type CodeExecutionService struct {
	tempDir         string
	timeout         time.Duration
	rustTimeout     time.Duration
	queueTimeout    time.Duration
	maxConcurrent   int
	jsDockerImage   string
	pyDockerImage   string
	rustDockerImage string
	executionSem    chan struct{} // Semaphore for concurrency control
	imagesVerified  bool          // Track if images were verified at startup
	jsPool          *ContainerPool
	pyPool          *ContainerPool
	rustPool        *ContainerPool
}

// ExecutionResult represents the result of code execution
type ExecutionResult struct {
	Success       bool   `json:"success"`
	Output        string `json:"output"`
	Error         string `json:"error"`
	ExecutionTime string `json:"executionTime"`
}

// Container represents a single pooled container
type Container struct {
	ID             string
	Language       string
	CreatedAt      time.Time
	LastUsedAt     time.Time
	ExecutionCount int
	MemoryErrors   int
}

// ContainerPool manages a pool of containers for a specific language
type ContainerPool struct {
	language   string
	image      string
	poolSize   int
	memoryMB   int
	cpus       string
	containers []*Container
	available  chan int // Indices of available containers
	mutex      sync.RWMutex
}

// ContainerMetrics tracks pool-level metrics
type ContainerMetrics struct {
	PoolUtilization float64
	HealthyCount    int
	TotalCount      int
	RotationEvents  int
}

// NewCodeExecutionService creates a new code execution service
func NewCodeExecutionService() (*CodeExecutionService, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "/tmp"
	}

	// Get max concurrent from environment or use default
	maxConcurrent := DefaultMaxConcurrent
	if envMax := os.Getenv("MAX_CONCURRENT_EXECUTIONS"); envMax != "" {
		if val, err := strconv.Atoi(envMax); err == nil && val > 0 {
			maxConcurrent = val
			log.Printf("Using MAX_CONCURRENT_EXECUTIONS=%d from environment", val)
		}
	}

	// Get queue timeout from environment or use default
	queueTimeout := DefaultQueueTimeout
	if envTimeout := os.Getenv("EXECUTION_QUEUE_TIMEOUT"); envTimeout != "" {
		if duration, err := time.ParseDuration(envTimeout); err == nil {
			queueTimeout = duration
			log.Printf("Using EXECUTION_QUEUE_TIMEOUT=%s from environment", duration)
		}
	}

	// Get pool sizes from environment or use defaults
	jsPoolSize := 5
	if envSize := os.Getenv("JS_POOL_SIZE"); envSize != "" {
		if val, err := strconv.Atoi(envSize); err == nil && val > 0 {
			jsPoolSize = val
			log.Printf("Using JS_POOL_SIZE=%d from environment", val)
		}
	}

	pyPoolSize := 5
	if envSize := os.Getenv("PY_POOL_SIZE"); envSize != "" {
		if val, err := strconv.Atoi(envSize); err == nil && val > 0 {
			pyPoolSize = val
			log.Printf("Using PY_POOL_SIZE=%d from environment", val)
		}
	}

	rustPoolSize := 2
	if envSize := os.Getenv("RUST_POOL_SIZE"); envSize != "" {
		if val, err := strconv.Atoi(envSize); err == nil && val > 0 {
			rustPoolSize = val
			log.Printf("Using RUST_POOL_SIZE=%d from environment", val)
		}
	}

	service := &CodeExecutionService{
		tempDir:         filepath.Join(homeDir, ".code-grader", "executions"),
		timeout:         10 * time.Second, // Increased from 5s to handle concurrent load
		rustTimeout:     30 * time.Second,
		queueTimeout:    queueTimeout,
		maxConcurrent:   maxConcurrent,
		jsDockerImage:   "code-grader-javascript",
		pyDockerImage:   "code-grader-python",
		rustDockerImage: "code-grader-rust",
		executionSem:    make(chan struct{}, maxConcurrent),
		imagesVerified:  false,
	}

	log.Printf("Code execution service configured: maxConcurrent=%d, queueTimeout=%s", maxConcurrent, queueTimeout)
	log.Printf("Container pool sizes: JS=%d, Python=%d, Rust=%d", jsPoolSize, pyPoolSize, rustPoolSize)

	// Initialize Docker images
	if err := service.initializeDocker(); err != nil {
		return nil, fmt.Errorf("failed to initialize Docker environment: %w", err)
	}

	// Mark images as verified after successful initialization
	service.imagesVerified = true

	// Create and initialize container pools
	log.Println("Initializing container pools...")

	// JavaScript pool
	service.jsPool = newContainerPool("javascript", service.jsDockerImage, jsPoolSize, 128, "0.5")
	if err := service.initializePool(service.jsPool); err != nil {
		return nil, fmt.Errorf("failed to initialize JavaScript pool: %w", err)
	}

	// Python pool
	service.pyPool = newContainerPool("python", service.pyDockerImage, pyPoolSize, 128, "0.5")
	if err := service.initializePool(service.pyPool); err != nil {
		service.stopPool(service.jsPool) // Clean up JS pool
		return nil, fmt.Errorf("failed to initialize Python pool: %w", err)
	}

	// Rust pool
	service.rustPool = newContainerPool("rust", service.rustDockerImage, rustPoolSize, 512, "1.0")
	if err := service.initializePool(service.rustPool); err != nil {
		service.stopPool(service.jsPool) // Clean up JS pool
		service.stopPool(service.pyPool) // Clean up Python pool
		return nil, fmt.Errorf("failed to initialize Rust pool: %w", err)
	}

	log.Println("All container pools initialized successfully")

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

	// Route to appropriate pool
	var pool *ContainerPool
	var filename string

	switch language {
	case "javascript", "js":
		pool = s.jsPool
		filename = "code.js"
	case "python", "py":
		pool = s.pyPool
		filename = "code.py"
	case "rust", "rs":
		pool = s.rustPool
		filename = "code.rs"
	default:
		return ExecutionResult{
			Success:       false,
			Error:         fmt.Sprintf("Language %s is not supported", language),
			ExecutionTime: "error",
		}
	}

	// Acquire container from pool (blocks if all containers are busy)
	select {
	case containerIdx := <-pool.available:
		// Successfully acquired container
		defer func() {
			pool.available <- containerIdx // Release container back to pool
			log.Printf("%s container released (%d/%d available)", pool.language, len(pool.available), pool.poolSize)
		}()
		log.Printf("%s container acquired (%d/%d available)", pool.language, len(pool.available)+1, pool.poolSize)

		// Execute in pooled container
		return s.executeInPooledContainer(code, pool, containerIdx, filename)

	case <-time.After(s.queueTimeout):
		// Timeout waiting for container
		log.Printf("%s pool timeout: all %d containers busy", pool.language, pool.poolSize)
		return ExecutionResult{
			Success:       false,
			Error:         "Server is at maximum capacity for this language, please try again in a few seconds",
			ExecutionTime: "timeout",
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

	if err := s.buildRustDockerImage(); err != nil {
		return fmt.Errorf("failed to build Rust Docker image: %w", err)
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

// buildRustDockerImage builds the Rust Docker image with extended timeout
func (s *CodeExecutionService) buildRustDockerImage() error {
	dockerfilePath := filepath.Join("docker", "rust.Dockerfile")

	// Rust image build takes longer due to stdlib pre-compilation cache
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "build", "-f", dockerfilePath, "-t", s.rustDockerImage, ".")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker build failed: %v\nstdout: %s\nstderr: %s", err, stdout.String(), stderr.String())
	}

	fmt.Printf("Successfully built Docker image: %s\n", s.rustDockerImage)
	return nil
}

// executeInDocker executes code in a Docker container
func (s *CodeExecutionService) executeInDocker(code, image, filename string) ExecutionResult {
	overallStart := time.Now()

	// Skip image verification if already verified at startup (optimization)
	var verifyDuration time.Duration
	if !s.imagesVerified {
		verifyStart := time.Now()
		if err := s.verifyDockerImage(image); err != nil {
			return ExecutionResult{
				Success:       false,
				Error:         fmt.Sprintf("Docker image not available: %v", err),
				ExecutionTime: "error",
			}
		}
		verifyDuration = time.Since(verifyStart)
		log.Printf("[TIMING] Docker image verify: %v", verifyDuration)
	} else {
		// Images verified at startup, skip per-request verification
		verifyDuration = 0
	}

	filePath := filepath.Join(s.tempDir, generateRandomID()+"_"+filename)

	mkdirStart := time.Now()
	if err := os.MkdirAll(s.tempDir, 0755); err != nil {
		return ExecutionResult{
			Success:       false,
			Error:         fmt.Sprintf("Failed to create temp directory: %v", err),
			ExecutionTime: "error",
		}
	}
	mkdirDuration := time.Since(mkdirStart)
	log.Printf("[TIMING] MkdirAll: %v", mkdirDuration)

	writeStart := time.Now()
	if err := os.WriteFile(filePath, []byte(code), 0644); err != nil {
		return ExecutionResult{
			Success:       false,
			Error:         fmt.Sprintf("Failed to write code file: %v", err),
			ExecutionTime: "error",
		}
	}
	writeDuration := time.Since(writeStart)
	log.Printf("[TIMING] Write file: %v", writeDuration)
	defer os.Remove(filePath)

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	dockerStart := time.Now()
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
	dockerDuration := time.Since(dockerStart)
	overallDuration := time.Since(overallStart)

	log.Printf("[TIMING] Docker execution: %v", dockerDuration)
	log.Printf("[TIMING] Total execution: %v (verify=%v, mkdir=%v, write=%v, docker=%v)",
		overallDuration, verifyDuration, mkdirDuration, writeDuration, dockerDuration)

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

// executeRustInDocker executes Rust code with compilation phase
func (s *CodeExecutionService) executeRustInDocker(code string) ExecutionResult {
	// Verify Docker image exists
	if err := s.verifyDockerImage(s.rustDockerImage); err != nil {
		return ExecutionResult{
			Success:       false,
			Error:         fmt.Sprintf("Docker image not available: %v", err),
			ExecutionTime: "error",
		}
	}

	// Create unique filename
	randomID := generateRandomID()
	filePath := filepath.Join(s.tempDir, randomID+"_code.rs")

	if err := os.MkdirAll(s.tempDir, 0755); err != nil {
		return ExecutionResult{
			Success:       false,
			Error:         fmt.Sprintf("Failed to create temp directory: %v", err),
			ExecutionTime: "error",
		}
	}

	// Write Rust source file
	if err := os.WriteFile(filePath, []byte(code), 0644); err != nil {
		return ExecutionResult{
			Success:       false,
			Error:         fmt.Sprintf("Failed to write code file: %v", err),
			ExecutionTime: "error",
		}
	}
	defer os.Remove(filePath)

	// Use extended timeout for Rust compilation
	ctx, cancel := context.WithTimeout(context.Background(), s.rustTimeout)
	defer cancel()

	// Two-phase execution script:
	// Phase 1: Compile with rustc (output to /workspace/tmp/)
	// Phase 2: Execute if compilation succeeds
	// Set TMPDIR to writable tmpfs for rustc intermediate files
	execScript := `#!/bin/sh
export TMPDIR=/workspace/tmp
rustc /workspace/code.rs -o /workspace/tmp/program 2>&1
if [ $? -ne 0 ]; then
    exit 1
fi
/workspace/tmp/program 2>&1
`

	cmd := exec.CommandContext(ctx, "docker", "run",
		"--rm",
		"--network", "none",
		"--memory", "512m",
		"--cpus", "1.0",
		"--pids-limit", "100",
		"--read-only",
		"--tmpfs", "/workspace/tmp:size=64m,exec,mode=1777",
		"--cap-drop", "ALL",
		"--security-opt", "no-new-privileges",
		"-v", filePath+":/workspace/code.rs:ro",
		s.rustDockerImage,
		"/bin/sh", "-c", execScript,
	)

	startTime := time.Now()
	output, err := cmd.CombinedOutput()
	executionTime := time.Since(startTime)

	// Handle timeout
	if ctx.Err() == context.DeadlineExceeded {
		return ExecutionResult{
			Success:       false,
			Error:         "Compilation or execution timed out (30s limit)",
			ExecutionTime: fmt.Sprintf("%.2fs", executionTime.Seconds()),
		}
	}

	// Handle compilation or execution errors
	if err != nil {
		errorMsg := string(output)

		// Detect compilation errors vs runtime errors
		if strings.Contains(errorMsg, "error:") || strings.Contains(errorMsg, "error[E") {
			errorMsg = "Compilation failed:\n" + errorMsg
		} else {
			errorMsg = "Runtime error:\n" + errorMsg
		}

		return ExecutionResult{
			Success:       false,
			Error:         errorMsg,
			ExecutionTime: fmt.Sprintf("%.2fs", executionTime.Seconds()),
		}
	}

	return ExecutionResult{
		Success:       true,
		Output:        string(output),
		ExecutionTime: fmt.Sprintf("%.2fs", executionTime.Seconds()),
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

// executeInPooledContainer executes code in a pre-started pooled container
func (s *CodeExecutionService) executeInPooledContainer(code string, pool *ContainerPool, containerIdx int, filename string) ExecutionResult {
	overallStart := time.Now()

	pool.mutex.RLock()
	container := pool.containers[containerIdx]
	pool.mutex.RUnlock()

	// Update container metrics
	pool.mutex.Lock()
	container.LastUsedAt = time.Now()
	container.ExecutionCount++
	pool.mutex.Unlock()

	// Create unique temporary filename to avoid collisions
	uniqueFilename := fmt.Sprintf("%s_%s", generateRandomID()[:8], filename)
	containerPath := fmt.Sprintf("/workspace/tmp/%s", uniqueFilename)

	// Write code to stdin and execute via docker exec
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	if pool.language == "rust" {
		ctx, cancel = context.WithTimeout(context.Background(), s.rustTimeout)
	}
	defer cancel()

	var execCmd *exec.Cmd

	if pool.language == "rust" {
		// Rust requires compilation step with TMPDIR set to writable location
		execScript := fmt.Sprintf(`export TMPDIR=/workspace/tmp && cat > %s && rustc %s -o /workspace/tmp/program 2>&1 && /workspace/tmp/program 2>&1; rm -f %s /workspace/tmp/program`,
			containerPath, containerPath, containerPath)
		execCmd = exec.CommandContext(ctx, "docker", "exec", "-i", container.ID, "sh", "-c", execScript)
	} else {
		// JavaScript and Python can execute directly
		var interpreter string
		if pool.language == "javascript" {
			interpreter = "node"
		} else {
			interpreter = "python3 -u"
		}

		execScript := fmt.Sprintf(`cat > %s && %s %s 2>&1; rm -f %s`, containerPath, interpreter, containerPath, containerPath)
		execCmd = exec.CommandContext(ctx, "docker", "exec", "-i", container.ID, "sh", "-c", execScript)
	}

	// Pipe code to stdin
	execCmd.Stdin = strings.NewReader(code)

	dockerStart := time.Now()
	output, err := execCmd.CombinedOutput()
	dockerDuration := time.Since(dockerStart)
	overallDuration := time.Since(overallStart)

	log.Printf("[TIMING] %s execution in pooled container: docker_exec=%v, total=%v, container_id=%s, executions=%d",
		pool.language, dockerDuration, overallDuration, container.ID[:12], container.ExecutionCount)

	// Check for timeout
	if ctx.Err() == context.DeadlineExceeded {
		return ExecutionResult{
			Success:       false,
			Error:         "Execution timed out",
			ExecutionTime: "timeout",
		}
	}

	// Check for execution errors
	if err != nil {
		// For Rust, detect compilation vs runtime errors
		errorMsg := string(output)
		if pool.language == "rust" {
			if strings.Contains(errorMsg, "error:") || strings.Contains(errorMsg, "error[E") {
				errorMsg = "Compilation failed:\n" + errorMsg
			} else {
				errorMsg = "Runtime error:\n" + errorMsg
			}
		}

		return ExecutionResult{
			Success:       false,
			Error:         errorMsg,
			ExecutionTime: fmt.Sprintf("%.2fs", dockerDuration.Seconds()),
		}
	}

	return ExecutionResult{
		Success:       true,
		Output:        string(output),
		ExecutionTime: fmt.Sprintf("%.2fs", dockerDuration.Seconds()),
	}
}

// generateRandomID generates a random hex ID
func generateRandomID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// newContainerPool creates a new container pool
func newContainerPool(language, image string, poolSize, memoryMB int, cpus string) *ContainerPool {
	return &ContainerPool{
		language:   language,
		image:      image,
		poolSize:   poolSize,
		memoryMB:   memoryMB,
		cpus:       cpus,
		containers: make([]*Container, poolSize),
		available:  make(chan int, poolSize),
	}
}

// healthCheckContainer verifies a container is running and can execute code
func (s *CodeExecutionService) healthCheckContainer(containerID, language string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Run a simple health check command based on language
	var healthCmd *exec.Cmd
	switch language {
	case "javascript":
		healthCmd = exec.CommandContext(ctx, "docker", "exec", containerID, "node", "-e", "console.log('ok')")
	case "python":
		healthCmd = exec.CommandContext(ctx, "docker", "exec", containerID, "python3", "-c", "print('ok')")
	case "rust":
		healthCmd = exec.CommandContext(ctx, "docker", "exec", containerID, "rustc", "--version")
	default:
		return fmt.Errorf("unknown language: %s", language)
	}

	output, err := healthCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("health check command failed: %v, output: %s", err, string(output))
	}

	log.Printf("Health check passed for %s container: %s", language, containerID[:12])
	return nil
}

// startContainer starts a single Docker container and returns its ID
func (s *CodeExecutionService) startContainer(language, image string, memoryMB int, cpus string) (string, error) {
	containerName := fmt.Sprintf("code-grader-%s-%s", language, generateRandomID()[:8])

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Start container with security constraints
	args := []string{
		"run", "-d",
		"--name", containerName,
		"--network", "none",
		"--memory", fmt.Sprintf("%dm", memoryMB),
		"--cpus", cpus,
		"--pids-limit", "50",
		"--read-only",
		"--cap-drop", "ALL",
		"--security-opt", "no-new-privileges",
	}

	// Add tmpfs for code execution (Rust needs more for compilation)
	if language == "rust" {
		args = append(args, "--tmpfs", "/workspace/tmp:size=64m,exec,mode=1777")
	} else {
		args = append(args, "--tmpfs", "/workspace/tmp:size=8m,exec,mode=1777")
	}

	args = append(args, image)

	// Keep container running with sleep infinity
	args = append(args, "sh", "-c", "sleep infinity")

	cmd := exec.CommandContext(ctx, "docker", args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to start container: %v\nstdout: %s\nstderr: %s", err, stdout.String(), stderr.String())
	}

	// Get container ID
	containerID := strings.TrimSpace(stdout.String())
	if containerID == "" {
		// If stdout is empty, use container name to get ID
		getIDCmd := exec.Command("docker", "inspect", "-f", "{{.Id}}", containerName)
		idOutput, err := getIDCmd.Output()
		if err != nil {
			return "", fmt.Errorf("failed to get container ID: %v", err)
		}
		containerID = strings.TrimSpace(string(idOutput))
	}

	log.Printf("Started %s container: %s (memory=%dMB, cpus=%s)", language, containerID[:12], memoryMB, cpus)
	return containerID, nil
}

// initializePool starts all containers in a pool
func (s *CodeExecutionService) initializePool(pool *ContainerPool) error {
	log.Printf("Initializing %s container pool (size=%d)...", pool.language, pool.poolSize)

	for i := 0; i < pool.poolSize; i++ {
		containerID, err := s.startContainer(pool.language, pool.image, pool.memoryMB, pool.cpus)
		if err != nil {
			// Clean up any containers we've already started
			s.stopPool(pool)
			return fmt.Errorf("failed to start container %d for %s: %w", i, pool.language, err)
		}

		// Health check the container
		if err := s.healthCheckContainer(containerID, pool.language); err != nil {
			s.stopPool(pool)
			return fmt.Errorf("health check failed for container %d (%s): %w", i, containerID[:12], err)
		}

		// Create container metadata
		container := &Container{
			ID:             containerID,
			Language:       pool.language,
			CreatedAt:      time.Now(),
			LastUsedAt:     time.Now(),
			ExecutionCount: 0,
			MemoryErrors:   0,
		}

		pool.containers[i] = container
		pool.available <- i // Mark as available
	}

	log.Printf("%s container pool ready: %d containers", pool.language, pool.poolSize)
	return nil
}

// stopContainer stops and removes a container
func (s *CodeExecutionService) stopContainer(containerID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Stop container
	stopCmd := exec.CommandContext(ctx, "docker", "stop", containerID)
	if err := stopCmd.Run(); err != nil {
		log.Printf("Warning: failed to stop container %s: %v", containerID[:12], err)
	}

	// Remove container
	rmCmd := exec.CommandContext(ctx, "docker", "rm", "-f", containerID)
	if err := rmCmd.Run(); err != nil {
		log.Printf("Warning: failed to remove container %s: %v", containerID[:12], err)
		return err
	}

	log.Printf("Stopped and removed container: %s", containerID[:12])
	return nil
}

// stopPool stops all containers in a pool
func (s *CodeExecutionService) stopPool(pool *ContainerPool) {
	if pool == nil {
		return
	}

	log.Printf("Stopping %s container pool...", pool.language)

	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	for _, container := range pool.containers {
		if container != nil && container.ID != "" {
			s.stopContainer(container.ID)
		}
	}
}

// resetContainerState cleans up any leftover files and processes in a container
func (s *CodeExecutionService) resetContainerState(containerID, language string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Kill any lingering processes (best effort, may fail if no processes)
	killCmd := exec.CommandContext(ctx, "docker", "exec", containerID, "sh", "-c", "pkill -9 -u coderunner || true")
	killCmd.Run()

	// Clean up workspace files
	cleanCmd := exec.CommandContext(ctx, "docker", "exec", containerID, "sh", "-c", "rm -rf /workspace/*.js /workspace/*.py /workspace/*.rs /workspace/tmp/* 2>/dev/null || true")
	if err := cleanCmd.Run(); err != nil {
		log.Printf("Warning: failed to clean container %s: %v", containerID[:12], err)
		return err
	}

	return nil
}

// shouldRotateContainer checks if a container should be rotated based on age or execution count
func (s *CodeExecutionService) shouldRotateContainer(container *Container) bool {
	if container == nil {
		return false
	}

	age := time.Since(container.CreatedAt)
	if age >= MaxContainerAge {
		log.Printf("Container %s needs rotation: age %v exceeds max %v", container.ID[:12], age, MaxContainerAge)
		return true
	}

	if container.ExecutionCount >= MaxExecutionsPerContainer {
		log.Printf("Container %s needs rotation: %d executions exceeds max %d", container.ID[:12], container.ExecutionCount, MaxExecutionsPerContainer)
		return true
	}

	if container.MemoryErrors > 0 {
		log.Printf("Container %s needs rotation: %d memory errors", container.ID[:12], container.MemoryErrors)
		return true
	}

	return false
}

// rotateContainer replaces a container in the pool with a new one (zero-downtime)
func (s *CodeExecutionService) rotateContainer(pool *ContainerPool, containerIdx int) error {
	pool.mutex.Lock()
	oldContainer := pool.containers[containerIdx]
	pool.mutex.Unlock()

	log.Printf("Rotating %s container %s (executions=%d, age=%v)", pool.language, oldContainer.ID[:12], oldContainer.ExecutionCount, time.Since(oldContainer.CreatedAt))

	// Step 1: Start new container
	newContainerID, err := s.startContainer(pool.language, pool.image, pool.memoryMB, pool.cpus)
	if err != nil {
		log.Printf("Failed to start replacement container: %v", err)
		return err // Keep old container running
	}

	// Step 2: Health check new container
	if err := s.healthCheckContainer(newContainerID, pool.language); err != nil {
		log.Printf("Health check failed for replacement container: %v", err)
		s.stopContainer(newContainerID)
		return err // Keep old container running
	}

	// Step 3: Create new container metadata
	newContainer := &Container{
		ID:             newContainerID,
		Language:       pool.language,
		CreatedAt:      time.Now(),
		LastUsedAt:     time.Now(),
		ExecutionCount: 0,
		MemoryErrors:   0,
	}

	// Step 4: Swap in pool (atomic operation)
	pool.mutex.Lock()
	pool.containers[containerIdx] = newContainer
	pool.mutex.Unlock()

	// Step 5: Stop old container (async, non-blocking)
	go func() {
		if err := s.stopContainer(oldContainer.ID); err != nil {
			log.Printf("Warning: failed to stop old container %s: %v", oldContainer.ID[:12], err)
		}
	}()

	log.Printf("Successfully rotated %s container: %s -> %s", pool.language, oldContainer.ID[:12], newContainerID[:12])
	return nil
}

// checkAndRotateContainersInPool checks all containers in a pool and rotates if needed
func (s *CodeExecutionService) checkAndRotateContainersInPool(pool *ContainerPool) {
	if pool == nil {
		return
	}

	for i := 0; i < pool.poolSize; i++ {
		pool.mutex.RLock()
		container := pool.containers[i]
		pool.mutex.RUnlock()

		if s.shouldRotateContainer(container) {
			// Try to acquire the container from the pool (non-blocking)
			select {
			case idx := <-pool.available:
				if idx == i {
					// We successfully acquired the container that needs rotation
					if err := s.rotateContainer(pool, i); err != nil {
						log.Printf("Failed to rotate container: %v", err)
					}
					// Return container to pool
					pool.available <- i
				} else {
					// We acquired a different container, return it
					pool.available <- idx
				}
			default:
				// Container is busy, skip rotation for now
				log.Printf("Skipping rotation for busy container %s", container.ID[:12])
			}
		}
	}
}

// PoolStatus represents the status of a single pool
type PoolStatus struct {
	Language    string              `json:"language"`
	PoolSize    int                 `json:"poolSize"`
	Available   int                 `json:"available"`
	Containers  []ContainerStatus   `json:"containers"`
}

// ContainerStatus represents the status of a single container
type ContainerStatus struct {
	ID             string  `json:"id"`
	ExecutionCount int     `json:"executionCount"`
	Age            string  `json:"age"`
	LastUsed       string  `json:"lastUsed"`
}

// GetPoolStatus returns the current status of all container pools
func (s *CodeExecutionService) GetPoolStatus() map[string]PoolStatus {
	status := make(map[string]PoolStatus)

	if s.jsPool != nil {
		status["javascript"] = s.getPoolStatusForPool(s.jsPool)
	}

	if s.pyPool != nil {
		status["python"] = s.getPoolStatusForPool(s.pyPool)
	}

	if s.rustPool != nil {
		status["rust"] = s.getPoolStatusForPool(s.rustPool)
	}

	return status
}

// getPoolStatusForPool gets status for a single pool
func (s *CodeExecutionService) getPoolStatusForPool(pool *ContainerPool) PoolStatus {
	pool.mutex.RLock()
	defer pool.mutex.RUnlock()

	containers := make([]ContainerStatus, len(pool.containers))
	for i, container := range pool.containers {
		if container != nil {
			containers[i] = ContainerStatus{
				ID:             container.ID[:12],
				ExecutionCount: container.ExecutionCount,
				Age:            time.Since(container.CreatedAt).Round(time.Second).String(),
				LastUsed:       time.Since(container.LastUsedAt).Round(time.Second).String(),
			}
		}
	}

	return PoolStatus{
		Language:   pool.language,
		PoolSize:   pool.poolSize,
		Available:  len(pool.available),
		Containers: containers,
	}
}

// Cleanup stops all container pools (call on shutdown)
func (s *CodeExecutionService) Cleanup() {
	log.Println("Cleaning up container pools...")
	s.stopPool(s.jsPool)
	s.stopPool(s.pyPool)
	s.stopPool(s.rustPool)
	log.Println("Container pools cleaned up")
}
