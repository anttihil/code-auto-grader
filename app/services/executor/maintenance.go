package executor

import (
	"log"
	"time"
)

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

		if container.ShouldRotate() {
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
