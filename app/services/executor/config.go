package executor

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
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
