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

// LanguageConfig defines all language-specific configuration
type LanguageConfig struct {
	Name              string        // e.g., "javascript", "python", "rust"
	FileExtension     string        // e.g., ".js", ".py", ".rs"
	DockerImage       string        // Docker image name
	DockerfilePath    string        // Path to Dockerfile
	ImageBuildTimeout time.Duration // Timeout for building Docker image
	ExecutionTimeout  time.Duration // Timeout for code execution
	MemoryMB          int           // Memory limit in MB
	CPUs              string        // CPU limit (e.g., "0.5", "1.0")
	TmpfsSize         string        // Tmpfs size (e.g., "8m", "64m")
	PIDsLimit         int           // Process limit
	Interpreter       string        // Interpreter command (e.g., "node", "python3", "")
	RequiresCompile   bool          // Whether code needs compilation
	Compiler          string        // Compiler command (e.g., "rustc")
	HealthCheckCmd    []string      // Command for health check
	PoolSize          int           // Container pool size
}

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
	config     *LanguageConfig
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

// initializeLanguageConfigs sets up all supported language configurations
func (s *CodeExecutionService) initializeLanguageConfigs() {
	// JavaScript configuration
	jsPoolSize := getEnvInt("JS_POOL_SIZE", 5)
	s.languageConfigs["javascript"] = &LanguageConfig{
		Name:              "javascript",
		FileExtension:     ".js",
		DockerImage:       "code-grader-javascript",
		DockerfilePath:    filepath.Join("docker", "javascript.Dockerfile"),
		ImageBuildTimeout: 120 * time.Second,
		ExecutionTimeout:  10 * time.Second,
		MemoryMB:          128,
		CPUs:              "0.5",
		TmpfsSize:         "8m",
		PIDsLimit:         50,
		Interpreter:       "node",
		RequiresCompile:   false,
		Compiler:          "",
		HealthCheckCmd:    []string{"node", "-e", "console.log('ok')"},
		PoolSize:          jsPoolSize,
	}
	// Alias for "js"
	s.languageConfigs["js"] = s.languageConfigs["javascript"]

	// Python configuration
	pyPoolSize := getEnvInt("PY_POOL_SIZE", 5)
	s.languageConfigs["python"] = &LanguageConfig{
		Name:              "python",
		FileExtension:     ".py",
		DockerImage:       "code-grader-python",
		DockerfilePath:    filepath.Join("docker", "python.Dockerfile"),
		ImageBuildTimeout: 120 * time.Second,
		ExecutionTimeout:  10 * time.Second,
		MemoryMB:          128,
		CPUs:              "0.5",
		TmpfsSize:         "8m",
		PIDsLimit:         50,
		Interpreter:       "python3 -u",
		RequiresCompile:   false,
		Compiler:          "",
		HealthCheckCmd:    []string{"python3", "-c", "print('ok')"},
		PoolSize:          pyPoolSize,
	}
	// Alias for "py"
	s.languageConfigs["py"] = s.languageConfigs["python"]

	// Rust configuration
	rustPoolSize := getEnvInt("RUST_POOL_SIZE", 2)
	s.languageConfigs["rust"] = &LanguageConfig{
		Name:              "rust",
		FileExtension:     ".rs",
		DockerImage:       "code-grader-rust",
		DockerfilePath:    filepath.Join("docker", "rust.Dockerfile"),
		ImageBuildTimeout: 300 * time.Second, // Rust build takes longer
		ExecutionTimeout:  30 * time.Second,  // Compilation + execution
		MemoryMB:          512,
		CPUs:              "1.0",
		TmpfsSize:         "64m",
		PIDsLimit:         100,
		Interpreter:       "",
		RequiresCompile:   true,
		Compiler:          "rustc",
		HealthCheckCmd:    []string{"rustc", "--version"},
		PoolSize:          rustPoolSize,
	}
	// Alias for "rs"
	s.languageConfigs["rs"] = s.languageConfigs["rust"]
}

// getEnvInt retrieves an integer from environment or returns default
func getEnvInt(key string, defaultValue int) int {
	if envValue := os.Getenv(key); envValue != "" {
		if val, err := strconv.Atoi(envValue); err == nil && val > 0 {
			log.Printf("Using %s=%d from environment", key, val)
			return val
		}
	}
	return defaultValue
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

// initializeDocker builds Docker images for code execution
func (s *CodeExecutionService) initializeDocker() error {
	fmt.Println("Initializing Docker execution environment...")

	// Build images for all unique configurations (skip aliases)
	builtImages := make(map[string]bool)
	for _, config := range s.languageConfigs {
		if builtImages[config.Name] {
			continue
		}
		if err := s.buildDockerImage(config); err != nil {
			return fmt.Errorf("failed to build %s Docker image: %w", config.Name, err)
		}
		builtImages[config.Name] = true
	}

	fmt.Println("Docker execution environment ready")
	return nil
}

// buildDockerImage builds a Docker image from the Dockerfile using language config
func (s *CodeExecutionService) buildDockerImage(config *LanguageConfig) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.ImageBuildTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "build", "-f", config.DockerfilePath, "-t", config.DockerImage, ".")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker build failed: %v\nstdout: %s\nstderr: %s", err, stdout.String(), stderr.String())
	}

	fmt.Printf("Successfully built Docker image: %s\n", config.DockerImage)
	return nil
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

// newContainerPool creates a new container pool
func newContainerPool(config *LanguageConfig) *ContainerPool {
	return &ContainerPool{
		language:   config.Name,
		config:     config,
		containers: make([]*Container, config.PoolSize),
		available:  make(chan int, config.PoolSize),
	}
}

// healthCheckContainer verifies a container is running and can execute code
func (s *CodeExecutionService) healthCheckContainer(containerID string, config *LanguageConfig) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Run health check command from config
	args := append([]string{"exec", containerID}, config.HealthCheckCmd...)
	healthCmd := exec.CommandContext(ctx, "docker", args...)

	output, err := healthCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("health check command failed: %v, output: %s", err, string(output))
	}

	log.Printf("Health check passed for %s container: %s", config.Name, containerID[:12])
	return nil
}

// startContainer starts a single Docker container and returns its ID
func (s *CodeExecutionService) startContainer(config *LanguageConfig) (string, error) {
	containerName := fmt.Sprintf("code-grader-%s-%s", config.Name, generateRandomID()[:8])

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Start container with security constraints
	args := []string{
		"run", "-d",
		"--name", containerName,
		"--network", "none",
		"--memory", fmt.Sprintf("%dm", config.MemoryMB),
		"--cpus", config.CPUs,
		"--pids-limit", strconv.Itoa(config.PIDsLimit),
		"--read-only",
		"--tmpfs", fmt.Sprintf("/workspace/tmp:size=%s,exec,mode=1777", config.TmpfsSize),
		"--cap-drop", "ALL",
		"--security-opt", "no-new-privileges",
		config.DockerImage,
		"sh", "-c", "sleep infinity",
	}

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

	log.Printf("Started %s container: %s (memory=%dMB, cpus=%s)", config.Name, containerID[:12], config.MemoryMB, config.CPUs)
	return containerID, nil
}

// initializePools creates and initializes all container pools
func (s *CodeExecutionService) initializePools() error {
	// Track which pools we've initialized for cleanup on error
	initialized := make([]string, 0, len(s.languageConfigs))

	// Initialize pools for unique configurations (skip aliases)
	createdPools := make(map[string]bool)
	for _, config := range s.languageConfigs {
		if createdPools[config.Name] {
			continue
		}

		log.Printf("Initializing %s container pool (size=%d)...", config.Name, config.PoolSize)

		pool := newContainerPool(config)
		if err := s.initializePool(pool); err != nil {
			// Clean up any pools we've already initialized
			for _, lang := range initialized {
				s.stopPool(s.pools[lang])
			}
			return fmt.Errorf("failed to initialize %s pool: %w", config.Name, err)
		}

		s.pools[config.Name] = pool
		initialized = append(initialized, config.Name)
		createdPools[config.Name] = true

		log.Printf("%s container pool ready: %d containers", config.Name, config.PoolSize)
	}

	return nil
}

// initializePool starts all containers in a pool
func (s *CodeExecutionService) initializePool(pool *ContainerPool) error {
	for i := 0; i < pool.config.PoolSize; i++ {
		containerID, err := s.startContainer(pool.config)
		if err != nil {
			// Clean up any containers we've already started
			s.stopPool(pool)
			return fmt.Errorf("failed to start container %d for %s: %w", i, pool.language, err)
		}

		// Health check the container
		if err := s.healthCheckContainer(containerID, pool.config); err != nil {
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
	newContainerID, err := s.startContainer(pool.config)
	if err != nil {
		log.Printf("Failed to start replacement container: %v", err)
		return err // Keep old container running
	}

	// Step 2: Health check new container
	if err := s.healthCheckContainer(newContainerID, pool.config); err != nil {
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

// containerHealthMonitor runs in the background and periodically health checks all containers
func (s *CodeExecutionService) containerHealthMonitor() {
	ticker := time.NewTicker(2 * time.Minute)
	defer ticker.Stop()

	log.Println("Container health monitor running (checking every 2 minutes)")

	for range ticker.C {
		log.Println("Running container health checks...")
		for _, pool := range s.pools {
			s.healthCheckPool(pool)
		}
		log.Println("Container health checks completed")
	}
}

// healthCheckPool performs health checks on all containers in a pool
func (s *CodeExecutionService) healthCheckPool(pool *ContainerPool) {
	if pool == nil {
		return
	}

	for i := 0; i < pool.config.PoolSize; i++ {
		pool.mutex.RLock()
		container := pool.containers[i]
		pool.mutex.RUnlock()

		if container == nil {
			continue
		}

		// Perform health check
		if err := s.healthCheckContainer(container.ID, pool.config); err != nil {
			log.Printf("HEALTH CHECK FAILED for %s container %s: %v", pool.language, container.ID[:12], err)

			// Mark container for rotation by incrementing memory errors
			pool.mutex.Lock()
			container.MemoryErrors++
			pool.mutex.Unlock()

			log.Printf("Container %s marked for rotation (health check failures: %d)", container.ID[:12], container.MemoryErrors)
		}
	}
}

// containerRotationScheduler runs in the background and periodically checks for containers that need rotation
func (s *CodeExecutionService) containerRotationScheduler() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	log.Println("Container rotation scheduler running (checking every 5 minutes)")

	for range ticker.C {
		log.Println("Running container rotation check...")
		for _, pool := range s.pools {
			s.checkAndRotateContainersInPool(pool)
		}
		log.Println("Container rotation check completed")
	}
}

// checkAndRotateContainersInPool checks all containers in a pool and rotates if needed
func (s *CodeExecutionService) checkAndRotateContainersInPool(pool *ContainerPool) {
	if pool == nil {
		return
	}

	for i := 0; i < pool.config.PoolSize; i++ {
		pool.mutex.RLock()
		container := pool.containers[i]
		pool.mutex.RUnlock()

		if container == nil {
			log.Printf("Skipping nil container at index %d in %s pool", i, pool.language)
			continue
		}

		if s.shouldRotateContainer(container) {
			// Try to acquire the container from the pool (non-blocking)
			select {
			case idx := <-pool.available:
				if idx == i {
					// We successfully acquired the container that needs rotation
					log.Printf("Rotating %s container at index %d", pool.language, i)
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
				log.Printf("Skipping rotation for busy %s container %s (index %d)", pool.language, container.ID[:12], i)
			}
		}
	}
}

// PoolStatus represents the status of a single pool
type PoolStatus struct {
	Language   string            `json:"language"`
	PoolSize   int               `json:"poolSize"`
	Available  int               `json:"available"`
	Containers []ContainerStatus `json:"containers"`
}

// ContainerStatus represents the status of a single container
type ContainerStatus struct {
	ID             string `json:"id"`
	ExecutionCount int    `json:"executionCount"`
	Age            string `json:"age"`
	LastUsed       string `json:"lastUsed"`
}

// GetPoolStatus returns the current status of all container pools
func (s *CodeExecutionService) GetPoolStatus() map[string]PoolStatus {
	status := make(map[string]PoolStatus)

	for language, pool := range s.pools {
		status[language] = s.getPoolStatusForPool(pool)
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
		PoolSize:   pool.config.PoolSize,
		Available:  len(pool.available),
		Containers: containers,
	}
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
