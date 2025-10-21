package executor

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// Container represents a single pooled container
type Container struct {
	ID             string
	Language       string
	CreatedAt      time.Time
	LastUsedAt     time.Time
	ExecutionCount int
	MemoryErrors   int
}

// ShouldRotate checks if a container should be rotated based on age, execution count, or errors
func (c *Container) ShouldRotate() bool {
	if c == nil {
		return false
	}

	age := time.Since(c.CreatedAt)
	if age >= MaxContainerAge {
		log.Printf("Container %s needs rotation: age %v exceeds max %v", c.ID[:12], age, MaxContainerAge)
		return true
	}

	if c.ExecutionCount >= MaxExecutionsPerContainer {
		log.Printf("Container %s needs rotation: %d executions exceeds max %d", c.ID[:12], c.ExecutionCount, MaxExecutionsPerContainer)
		return true
	}

	if c.MemoryErrors > 0 {
		log.Printf("Container %s needs rotation: %d memory errors", c.ID[:12], c.MemoryErrors)
		return true
	}

	return false
}

// ContainerPool manages a pool of containers for a specific language
type ContainerPool struct {
	language   string
	config     *LanguageConfig
	containers []*Container
	available  chan int // Indices of available containers
	mutex      sync.RWMutex
}

// GetStatus returns the current status of this container pool
func (p *ContainerPool) GetStatus() PoolStatus {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	containers := make([]ContainerStatus, len(p.containers))
	for i, container := range p.containers {
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
		Language:   p.language,
		PoolSize:   p.config.PoolSize,
		Available:  len(p.available),
		Containers: containers,
	}
}

// ContainerMetrics tracks pool-level metrics
type ContainerMetrics struct {
	PoolUtilization float64
	HealthyCount    int
	TotalCount      int
	RotationEvents  int
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

// newContainerPool creates a new container pool
func newContainerPool(config *LanguageConfig) *ContainerPool {
	return &ContainerPool{
		language:   config.Name,
		config:     config,
		containers: make([]*Container, config.PoolSize),
		available:  make(chan int, config.PoolSize),
	}
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

// GetPoolStatus returns the current status of all container pools
func (s *CodeExecutionService) GetPoolStatus() map[string]PoolStatus {
	status := make(map[string]PoolStatus)

	for language, pool := range s.pools {
		status[language] = pool.GetStatus()
	}

	return status
}
