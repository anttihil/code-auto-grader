package executor

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/anttihil/code-auto-grader/internal/logger"
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
			logger.Logger.Info("using MAX_CONCURRENT_EXECUTIONS from environment",
				slog.Int("max_concurrent", val),
			)
		}
	}

	// Get queue timeout from environment or use default
	queueTimeout := DefaultQueueTimeout
	if envTimeout := os.Getenv("EXECUTION_QUEUE_TIMEOUT"); envTimeout != "" {
		if duration, err := time.ParseDuration(envTimeout); err == nil {
			queueTimeout = duration
			logger.Logger.Info("using EXECUTION_QUEUE_TIMEOUT from environment",
				slog.Duration("timeout", duration),
			)
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

	logger.Logger.Info("code execution service configured",
		slog.Int("max_concurrent", maxConcurrent),
		slog.Duration("queue_timeout", queueTimeout),
	)

	// Log pool sizes
	for _, config := range service.languageConfigs {
		logger.Logger.Info("container pool configured",
			slog.String("language", config.Name),
			slog.Int("pool_size", config.PoolSize),
		)
	}

	// Initialize Docker images
	if err := service.initializeDocker(); err != nil {
		return nil, fmt.Errorf("failed to initialize Docker environment: %w", err)
	}

	// Mark images as verified after successful initialization
	service.imagesVerified = true

	// Create and initialize container pools
	logger.Logger.Info("initializing container pools",
		slog.Int("pool_count", len(service.languageConfigs)),
	)
	if err := service.initializePools(); err != nil {
		return nil, err
	}

	logger.Logger.Info("all container pools initialized successfully")

	// Start background maintenance goroutines
	go service.containerRotationScheduler()
	logger.Logger.Info("container rotation scheduler started")

	go service.containerHealthMonitor()
	logger.Logger.Info("container health monitor started")

	return service, nil
}

// ExecuteCodeWithContext executes code based on the language with request context
func (s *CodeExecutionService) ExecuteCodeWithContext(ctx context.Context, code, language string) ExecutionResult {
	requestID := logger.GetRequestID(ctx)
	return s.executeCodeInternal(ctx, code, language, requestID)
}

// ExecuteCode executes code based on the language (backwards compatibility)
func (s *CodeExecutionService) ExecuteCode(code, language string) ExecutionResult {
	return s.executeCodeInternal(context.Background(), code, language, "unknown")
}

// executeCodeInternal is the internal implementation with context support
func (s *CodeExecutionService) executeCodeInternal(ctx context.Context, code, language, requestID string) ExecutionResult {
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
			logger.Logger.Debug("container released",
				slog.String("request_id", requestID),
				slog.String("language", pool.language),
				slog.Int("pool_available", len(pool.available)),
				slog.Int("pool_size", pool.config.PoolSize),
			)
		}()
		logger.Logger.Debug("container acquired",
			slog.String("request_id", requestID),
			slog.String("language", pool.language),
			slog.Int("pool_available", len(pool.available)+1),
			slog.Int("pool_size", pool.config.PoolSize),
		)

		// Execute in pooled container
		return s.executeInPooledContainer(ctx, code, pool, containerIdx, requestID)

	case <-time.After(s.queueTimeout):
		// Timeout waiting for container
		logger.Logger.Warn("pool timeout - all containers busy",
			slog.String("request_id", requestID),
			slog.String("language", pool.language),
			slog.Int("pool_size", pool.config.PoolSize),
			slog.Duration("timeout", s.queueTimeout),
		)
		return ExecutionResult{
			Success:       false,
			Error:         "Server is at maximum capacity for this language, please try again in a few seconds",
			ExecutionTime: "timeout",
		}
	}
}

// executeInPooledContainer executes code in a pre-started pooled container
func (s *CodeExecutionService) executeInPooledContainer(ctx context.Context, code string, pool *ContainerPool, containerIdx int, requestID string) ExecutionResult {
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
		// Check if this is Rust test code
		isRustTest := pool.language == "rust" && (strings.Contains(code, "#[test]") || strings.Contains(code, "#[cfg(test)]"))

		programPath := "/workspace/tmp/program"
		var execScript string

		if isRustTest {
			// Rust test mode: compile with --test flag and run test binary
			execScript = fmt.Sprintf(`
cleanup() { rm -f %s %s 2>/dev/null; }
trap cleanup EXIT
export TMPDIR=/workspace/tmp
cat > %s && %s --test %s -o %s 2>&1 && %s 2>&1; test_exit=$?; exit $test_exit`,
				containerPath, programPath, containerPath, config.Compiler, containerPath, programPath, programPath)
		} else {
			// Regular compilation mode
			execScript = fmt.Sprintf(`
cleanup() { rm -f %s %s 2>/dev/null; }
trap cleanup EXIT
export TMPDIR=/workspace/tmp
cat > %s && %s %s -o %s 2>&1 && %s 2>&1`,
				containerPath, programPath, containerPath, config.Compiler, containerPath, programPath, programPath)
		}

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

	logger.Logger.Info("code execution in pooled container",
		slog.String("request_id", requestID),
		slog.String("language", pool.language),
		slog.Duration("docker_exec_ms", dockerDuration),
		slog.Duration("total_ms", overallDuration),
		slog.String("container_id", containerID[:12]),
		slog.Int("execution_count", execCount),
	)

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

		// For Rust tests, exit code 101 means "tests ran but some failed"
		// This is actually success from an execution perspective
		isRustTest := pool.language == "rust" && (strings.Contains(code, "#[test]") || strings.Contains(code, "#[cfg(test)]"))
		if isRustTest && strings.Contains(errorMsg, "test result:") {
			// Test runner completed successfully, return output even though exit code != 0
			return ExecutionResult{
				Success:       true,
				Output:        errorMsg,
				ExecutionTime: fmt.Sprintf("%.2fs", dockerDuration.Seconds()),
			}
		}

		// Sanitize container paths from error messages
		errorMsg = sanitizeErrorMessage(errorMsg, containerPath)

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

// sanitizeErrorMessage removes container-specific paths from error messages
// Converts "/workspace/tmp/ad894b93_code.rs:26:28" to "line 26, column 28"
func sanitizeErrorMessage(errorMsg, containerPath string) string {
	if errorMsg == "" {
		return errorMsg
	}

	// Extract the base filename pattern (e.g., "/workspace/tmp/ad894b93_code.rs")
	// containerPath is like "/workspace/tmp/ad894b93_code.rs"

	// Replace full paths with just "line X, column Y" or "line X"
	// Pattern: /workspace/tmp/HASH_code.EXT:LINE:COLUMN
	lines := strings.Split(errorMsg, "\n")
	sanitized := make([]string, 0, len(lines))

	for _, line := range lines {
		// Replace container path references
		if strings.Contains(line, containerPath) {
			// Extract line:column from paths like "/workspace/tmp/ad894b93_code.rs:26:28"
			line = strings.ReplaceAll(line, containerPath, "(your code)")
		}
		sanitized = append(sanitized, line)
	}

	return strings.Join(sanitized, "\n")
}

// Cleanup stops all container pools with graceful shutdown (call on shutdown)
func (s *CodeExecutionService) Cleanup() {
	logger.Logger.Info("initiating graceful shutdown")

	// Signal to stop accepting new requests
	close(s.shutdownSignal)
	logger.Logger.Info("no longer accepting new execution requests")

	// Wait for in-flight executions with timeout
	done := make(chan struct{})
	go func() {
		s.activeExecutions.Wait()
		close(done)
	}()

	select {
	case <-done:
		logger.Logger.Info("all in-flight executions completed successfully")
	case <-time.After(30 * time.Second):
		logger.Logger.Warn("shutdown timeout reached, forcing cleanup",
			slog.Duration("timeout", 30*time.Second),
		)
	}

	// Now clean up container pools
	logger.Logger.Info("stopping container pools",
		slog.Int("pool_count", len(s.pools)),
	)
	for _, pool := range s.pools {
		s.stopPool(pool)
	}
	logger.Logger.Info("container pools cleaned up - shutdown complete")
}
