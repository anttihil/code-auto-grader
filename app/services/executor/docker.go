package executor

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

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
