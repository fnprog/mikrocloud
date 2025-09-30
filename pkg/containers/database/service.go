package database

import (
	"context"
	"fmt"
	"strconv"

	"github.com/mikrocloud/mikrocloud/internal/domain/databases"
	"github.com/mikrocloud/mikrocloud/pkg/containers/manager"
)

// Service implementation
type Service struct {
	containerManager manager.ContainerManager
	imageResolver    DatabaseImageResolver
	configBuilder    ContainerConfigBuilder
}

func NewDatabaseDeploymentService(
	containerManager manager.ContainerManager,
	imageResolver DatabaseImageResolver,
	configBuilder ContainerConfigBuilder,
) DatabaseDeploymentService {
	return &Service{
		containerManager: containerManager,
		imageResolver:    imageResolver,
		configBuilder:    configBuilder,
	}
}

// Deploy creates and starts a database container
func (s *Service) Deploy(ctx context.Context, database *databases.Database) (*DeploymentResult, error) {
	// Build container configuration
	config, err := s.configBuilder.BuildConfig(database)
	if err != nil {
		return nil, fmt.Errorf("failed to build container config: %w", err)
	}

	// Convert to manager.ContainerConfig
	containerConfig := manager.ContainerConfig{
		Image:         config.Image,
		Name:          config.ContainerName,
		Ports:         map[string]string{config.Port: config.Port},
		Environment:   config.Environment,
		Volumes:       config.Volumes,
		RestartPolicy: "unless-stopped",
		Command:       config.Command,
	}

	// Create and start the container
	containerID, err := s.containerManager.Create(ctx, containerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create container: %w", err)
	}

	// Start the container
	if err := s.containerManager.Start(ctx, containerID); err != nil {
		// Clean up on failure
		_ = s.containerManager.Delete(ctx, containerID)
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	return &DeploymentResult{
		ContainerID:   containerID,
		ContainerName: config.ContainerName,
		Port:          config.Port,
		Status:        "running",
		CreatedAt:     "", // Will be populated by inspection
	}, nil
}

// Start starts an existing database container
func (s *Service) Start(ctx context.Context, database *databases.Database) error {
	containerName := s.buildContainerName(database)

	// Find container by name
	containers, err := s.containerManager.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list containers: %w", err)
	}

	for _, container := range containers {
		if container.Name == containerName {
			return s.containerManager.Start(ctx, container.ID)
		}
	}

	return fmt.Errorf("container %s not found", containerName)
}

// Stop stops a running database container
func (s *Service) Stop(ctx context.Context, database *databases.Database) error {
	containerName := s.buildContainerName(database)

	// Find container by name
	containers, err := s.containerManager.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list containers: %w", err)
	}

	for _, container := range containers {
		if container.Name == containerName {
			return s.containerManager.Stop(ctx, container.ID)
		}
	}

	return fmt.Errorf("container %s not found", containerName)
}

// Remove removes a database container completely
func (s *Service) Remove(ctx context.Context, database *databases.Database) error {
	containerName := s.buildContainerName(database)

	// Find container by name
	containers, err := s.containerManager.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list containers: %w", err)
	}

	for _, container := range containers {
		if container.Name == containerName {
			// Stop first if running
			_ = s.containerManager.Stop(ctx, container.ID)
			return s.containerManager.Delete(ctx, container.ID)
		}
	}

	return fmt.Errorf("container %s not found", containerName)
}

// GetStatus returns the current status of a database container
func (s *Service) GetStatus(ctx context.Context, database *databases.Database) (*ContainerStatus, error) {
	containerName := s.buildContainerName(database)

	// Find container by name
	containers, err := s.containerManager.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	for _, container := range containers {
		if container.Name == containerName {
			info, err := s.containerManager.Inspect(ctx, container.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to inspect container: %w", err)
			}

			// Extract port information
			var port string
			for hostPort := range info.Ports {
				port = hostPort
				break // Take first port
			}

			return &ContainerStatus{
				ID:        info.ID,
				Name:      info.Name,
				State:     info.State,
				Status:    info.Status,
				Port:      port,
				Health:    "unknown", // TODO: Add health check support
				StartedAt: "",        // TODO: Add started timestamp
			}, nil
		}
	}

	return nil, fmt.Errorf("container %s not found", containerName)
}

// GetLogs retrieves logs from a database container
func (s *Service) GetLogs(ctx context.Context, database *databases.Database, follow bool) ([]byte, error) {
	containerName := s.buildContainerName(database)

	// Find container by name
	containers, err := s.containerManager.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	for _, container := range containers {
		if container.Name == containerName {
			logStream, err := s.containerManager.StreamLogs(ctx, container.ID, follow)
			if err != nil {
				return nil, fmt.Errorf("failed to get logs: %w", err)
			}
			defer logStream.Close()

			// Read all logs (for non-follow mode)
			if !follow {
				buf := make([]byte, 64*1024) // 64KB buffer
				n, err := logStream.Read(buf)
				if err != nil && err.Error() != "EOF" {
					return nil, fmt.Errorf("failed to read logs: %w", err)
				}
				return buf[:n], nil
			}

			return nil, fmt.Errorf("follow mode not implemented yet")
		}
	}

	return nil, fmt.Errorf("container %s not found", containerName)
}

// buildContainerName creates a consistent container name for the database
func (s *Service) buildContainerName(database *databases.Database) string {
	return fmt.Sprintf("mikrocloud-%s-%s-%s",
		database.ProjectID(),
		database.EnvironmentID(),
		database.Name().String())
}

// Helper to build port mapping
func buildPortMapping(port int) string {
	return strconv.Itoa(port)
}
