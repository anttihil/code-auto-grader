package executor

import (
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

// CodeExecutionService handles secure code execution in Docker containers
type CodeExecutionService struct {
	tempDir          string
	queueTimeout     time.Duration
	maxConcurrent    int
	executionSem     chan struct{} // Semaphore for concurrency control
	imagesVerified   bool          // Track if images were verified at startup
	languageConfigs  map[string]*LanguageConfig
	pools            map[string]*ContainerPool
	activeExecutions sync.WaitGroup // Track in-flight executions for graceful shutdown
	shutdownSignal   chan struct{}  // Signal to stop accepting new requests
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

	service := &CodeExecutionService{
		tempDir:         filepath.Join(homeDir, ".code-grader", "executions"),
		queueTimeout:    queueTimeout,
		maxConcurrent:   maxConcurrent,
		executionSem:    make(chan struct{}, maxConcurrent),
		imagesVerified:  false,
		languageConfigs: make(map[string]*LanguageConfig),
		pools:           make(map[string]*ContainerPool),
		shutdownSignal:  make(chan struct{}),
	}

	// Initialize language configurations
	service.initializeLanguageConfigs()

	log.Printf("Code execution service configured: maxConcurrent=%d, queueTimeout=%s", maxConcurrent, queueTimeout)
	logPoolSizes := make([]string, 0, len(service.languageConfigs))
	for _, config := range service.languageConfigs {
		logPoolSizes = append(logPoolSizes, fmt.Sprintf("%s=%d", config.Name, config.PoolSize))
	}
	log.Printf("Container pool sizes: %s", strings.Join(logPoolSizes, ", "))

	// Initialize Docker images
	if err := service.initializeDocker(); err != nil {
		return nil, fmt.Errorf("failed to initialize Docker environment: %w", err)
	}

	// Mark images as verified after successful initialization
	service.imagesVerified = true

	// Create and initialize container pools
	log.Println("Initializing container pools...")
	if err := service.initializePools(); err != nil {
		return nil, err
	}

	log.Println("All container pools initialized successfully")

	// Start background maintenance goroutines
	go service.containerRotationScheduler()
	log.Println("Container rotation scheduler started")

	go service.containerHealthMonitor()
	log.Println("Container health monitor started")

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

	// Get language configuration
	config, exists := s.languageConfigs[language]
	if !exists {
		return ExecutionResult{
			Success:       false,
			Error:         fmt.Sprintf("Language %s is not supported", language),
			ExecutionTime: "error",
		}
	}

	// Get pool for this language (use canonical name, not alias)
	pool, exists := s.pools[config.Name]
	if !exists {
		return ExecutionResult{
			Success:       false,
			Error:         fmt.Sprintf("Pool for language %s not initialized", language),
			ExecutionTime: "error",
		}
	}

	// Acquire container from pool (blocks if all containers are busy)
	select {
	case containerIdx := <-pool.available:
		// Successfully acquired container
		defer func() {
			pool.available <- containerIdx // Release container back to pool
			log.Printf("%s container released (%d/%d available)", pool.language, len(pool.available), pool.config.PoolSize)
		}()
		log.Printf("%s container acquired (%d/%d available)", pool.language, len(pool.available)+1, pool.config.PoolSize)

		// Execute in pooled container
		return s.executeInPooledContainer(code, pool, containerIdx)

	case <-time.After(s.queueTimeout):
		// Timeout waiting for container
		log.Printf("%s pool timeout: all %d containers busy", pool.language, pool.config.PoolSize)
		return ExecutionResult{
			Success:       false,
			Error:         "Server is at maximum capacity for this language, please try again in a few seconds",
			ExecutionTime: "timeout",
		}
	}
}

// executeInPooledContainer executes code in a pre-started pooled container
func (s *CodeExecutionService) executeInPooledContainer(code string, pool *ContainerPool, containerIdx int) ExecutionResult {
	// Check if server is shutting down
	select {
	case <-s.shutdownSignal:
		return ExecutionResult{
			Success:       false,
			Error:         "Server is shutting down, please try again later",
			ExecutionTime: "error",
		}
	default:
	}

	// Track this execution for graceful shutdown
	s.activeExecutions.Add(1)
	defer s.activeExecutions.Done()

	overallStart := time.Now()

	// Get container and update metrics atomically to avoid race conditions
	pool.mutex.Lock()
	container := pool.containers[containerIdx]
	if container == nil {
		pool.mutex.Unlock()
		return ExecutionResult{
			Success:       false,
			Error:         "Container not available",
			ExecutionTime: "error",
		}
	}
	container.LastUsedAt = time.Now()
	container.ExecutionCount++
	containerID := container.ID
	pool.mutex.Unlock()

	config := pool.config

	// Create unique temporary filename to avoid collisions
	uniqueFilename := fmt.Sprintf("%s_%s%s", generateRandomID()[:8], "code", config.FileExtension)
	containerPath := fmt.Sprintf("/workspace/tmp/%s", uniqueFilename)

	// Write code to stdin and execute via docker exec
	ctx, cancel := context.WithTimeout(context.Background(), config.ExecutionTimeout)
	defer cancel()

	var execCmd *exec.Cmd

	if config.RequiresCompile {
		// Compiled languages (e.g., Rust)
		// Use trap to ensure cleanup runs even on failure
		programPath := "/workspace/tmp/program"
		execScript := fmt.Sprintf(`
cleanup() { rm -f %s %s 2>/dev/null; }
trap cleanup EXIT
export TMPDIR=/workspace/tmp
cat > %s && %s %s -o %s 2>&1 && %s 2>&1`,
			containerPath, programPath, containerPath, config.Compiler, containerPath, programPath, programPath)
		execCmd = exec.CommandContext(ctx, "docker", "exec", "-i", containerID, "sh", "-c", execScript)
	} else {
		// Interpreted languages (e.g., JavaScript, Python)
		// Use trap to ensure cleanup runs even on failure
		execScript := fmt.Sprintf(`
cleanup() { rm -f %s 2>/dev/null; }
trap cleanup EXIT
cat > %s && %s %s 2>&1`,
			containerPath, containerPath, config.Interpreter, containerPath)
		execCmd = exec.CommandContext(ctx, "docker", "exec", "-i", containerID, "sh", "-c", execScript)
	}

	// Pipe code to stdin
	execCmd.Stdin = strings.NewReader(code)

	dockerStart := time.Now()
	output, err := execCmd.CombinedOutput()
	dockerDuration := time.Since(dockerStart)
	overallDuration := time.Since(overallStart)

	// Read execution count for logging (thread-safe read)
	pool.mutex.RLock()
	execCount := container.ExecutionCount
	pool.mutex.RUnlock()

	log.Printf("[TIMING] %s execution in pooled container: docker_exec=%v, total=%v, container_id=%s, executions=%d",
		pool.language, dockerDuration, overallDuration, containerID[:12], execCount)

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
		errorMsg := string(output)

		// For compiled languages, detect compilation vs runtime errors
		if config.RequiresCompile {
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

// Cleanup stops all container pools with graceful shutdown (call on shutdown)
func (s *CodeExecutionService) Cleanup() {
	log.Println("Initiating graceful shutdown...")

	// Signal to stop accepting new requests
	close(s.shutdownSignal)
	log.Println("No longer accepting new execution requests")

	// Wait for in-flight executions with timeout
	done := make(chan struct{})
	go func() {
		s.activeExecutions.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("All in-flight executions completed successfully")
	case <-time.After(30 * time.Second):
		log.Println("Shutdown timeout reached (30s), forcing cleanup of remaining executions")
	}

	// Now clean up container pools
	log.Println("Stopping container pools...")
	for _, pool := range s.pools {
		s.stopPool(pool)
	}
	log.Println("Container pools cleaned up. Shutdown complete.")
}
