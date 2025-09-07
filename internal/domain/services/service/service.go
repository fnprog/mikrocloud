package services

import (
	"context"
	"fmt"

	"github.com/mikrocloud/mikrocloud/internal/container/build"
	"github.com/mikrocloud/mikrocloud/internal/container/manager"
	"github.com/mikrocloud/mikrocloud/internal/domain/services"
	"github.com/mikrocloud/mikrocloud/internal/domain/services/repository"
)

// ServiceService handles services business logic
type ServiceService struct {
	serviceRepo      repository.Repository
	containerManager manager.ContainerManager
	buildService     *build.BuildService
}

func NewServiceService(
	serviceRepo repository.Repository,
	containerManager manager.ContainerManager,
	buildService *build.BuildService,
) *ServiceService {
	return &ServiceService{
		serviceRepo:      serviceRepo,
		containerManager: containerManager,
		buildService:     buildService,
	}
}

// CreateService creates a new services
func (s *ServiceService) CreateService(ctx context.Context, svc *services.Service) error {
	// Validate services
	if svc == nil {
		return fmt.Errorf("services cannot be nil")
	}

	// Check if services name already exists in the project/environment
	existingServices, err := s.serviceRepo.FindByProjectAndEnvironment(
		svc.ProjectID().String(),
		svc.EnvironmentID().String(),
	)
	if err != nil {
		return fmt.Errorf("failed to check existing services: %w", err)
	}

	for _, existing := range existingServices {
		if existing.Name().String() == svc.Name().String() {
			return fmt.Errorf("services with name '%s' already exists in this environment", svc.Name().String())
		}
	}

	// Save the services
	if err := s.serviceRepo.Save(svc); err != nil {
		return fmt.Errorf("failed to save services: %w", err)
	}

	return nil
}

// GetService retrieves a services by ID
func (s *ServiceService) GetService(ctx context.Context, id services.ServiceID) (*services.Service, error) {
	svc, err := s.serviceRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find services: %w", err)
	}
	return svc, nil
}

// UpdateService updates a services
func (s *ServiceService) UpdateService(ctx context.Context, svc *services.Service) error {
	if err := s.serviceRepo.Update(svc); err != nil {
		return fmt.Errorf("failed to update services: %w", err)
	}
	return nil
}

// DeleteService deletes a services and its containers
func (s *ServiceService) DeleteService(ctx context.Context, id services.ServiceID) error {
	// Get services first
	svc, err := s.serviceRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find services: %w", err)
	}

	// Stop and remove any running containers
	if svc.Status() == services.ServiceStatusRunning {
		if err := s.StopService(ctx, id); err != nil {
			return fmt.Errorf("failed to stop services before deletion: %w", err)
		}
	}

	// Delete from database
	if err := s.serviceRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete services: %w", err)
	}

	return nil
}

// DeployService builds and deploys a services
func (s *ServiceService) DeployService(ctx context.Context, id services.ServiceID) error {
	// Get services
	svc, err := s.serviceRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find services: %w", err)
	}

	// Check if services can be deployed
	if err := svc.CanDeploy(); err != nil {
		return err
	}

	// Update status to building
	svc.ChangeStatus(services.ServiceStatusBuilding)
	if err := s.serviceRepo.Update(svc); err != nil {
		return fmt.Errorf("failed to update services status: %w", err)
	}

	// Create build request
	buildRequest := s.createBuildRequest(svc)

	// Build the image
	buildResult, err := s.buildService.BuildImage(ctx, buildRequest)
	if err != nil {
		svc.ChangeStatus(services.ServiceStatusFailed)
		s.serviceRepo.Update(svc)
		return fmt.Errorf("failed to build services: %w", err)
	}

	if !buildResult.Success {
		svc.ChangeStatus(services.ServiceStatusFailed)
		s.serviceRepo.Update(svc)
		return fmt.Errorf("build failed: %s", buildResult.Error)
	}

	// Update status to deploying
	svc.ChangeStatus(services.ServiceStatusDeploying)
	if err := s.serviceRepo.Update(svc); err != nil {
		return fmt.Errorf("failed to update services status: %w", err)
	}

	// Deploy the container
	if err := s.deployContainer(ctx, svc, buildResult.ImageTag); err != nil {
		svc.ChangeStatus(services.ServiceStatusFailed)
		s.serviceRepo.Update(svc)
		return fmt.Errorf("failed to deploy container: %w", err)
	}

	// Update status to running
	svc.ChangeStatus(services.ServiceStatusRunning)
	if err := s.serviceRepo.Update(svc); err != nil {
		return fmt.Errorf("failed to update services status: %w", err)
	}

	return nil
}

// StopService stops a running services
func (s *ServiceService) StopService(ctx context.Context, id services.ServiceID) error {
	// Get services
	svc, err := s.serviceRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find services: %w", err)
	}

	// Check if services can be stopped
	if err := svc.CanStop(); err != nil {
		return err
	}

	// Stop the container (using services name as container name)
	containerName := s.getContainerName(svc)
	if err := s.containerManager.Stop(ctx, containerName); err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}

	// Update status
	svc.ChangeStatus(services.ServiceStatusStopped)
	if err := s.serviceRepo.Update(svc); err != nil {
		return fmt.Errorf("failed to update services status: %w", err)
	}

	return nil
}

// RestartService restarts a services
func (s *ServiceService) RestartService(ctx context.Context, id services.ServiceID) error {
	// Get services
	svc, err := s.serviceRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find services: %w", err)
	}

	// Restart the container
	containerName := s.getContainerName(svc)
	if err := s.containerManager.Restart(ctx, containerName); err != nil {
		return fmt.Errorf("failed to restart container: %w", err)
	}

	// Update status
	svc.ChangeStatus(services.ServiceStatusRunning)
	if err := s.serviceRepo.Update(svc); err != nil {
		return fmt.Errorf("failed to update services status: %w", err)
	}

	return nil
}

// ListServices lists services by project and environment
func (s *ServiceService) ListServices(ctx context.Context, projectID, environmentID string) ([]*services.Service, error) {
	if projectID != "" && environmentID != "" {
		return s.serviceRepo.FindByProjectAndEnvironment(projectID, environmentID)
	} else if projectID != "" {
		return s.serviceRepo.FindByProjectID(projectID)
	} else if environmentID != "" {
		return s.serviceRepo.FindByEnvironmentID(environmentID)
	}
	return s.serviceRepo.List()
}

// Helper methods

func (s *ServiceService) createBuildRequest(svc *services.Service) build.BuildRequest {
	buildRequest := build.BuildRequest{
		ID:          svc.ID().String(),
		GitRepo:     svc.GitURL().URL(),
		GitBranch:   svc.GitURL().Branch(),
		ContextRoot: svc.GitURL().ContextRoot(),
		ImageTag:    s.getImageTag(svc),
		Environment: svc.Environment(),
	}

	buildConfig := svc.BuildConfig()
	switch buildConfig.BuildpackType() {
	case services.BuildpackNixpacks:
		buildRequest.BuildpackType = build.Nixpacks
		if config := buildConfig.NixpacksConfig(); config != nil {
			buildRequest.NixpacksConfig = &build.NixpacksConfig{
				StartCommand: config.StartCommand,
				BuildCommand: config.BuildCommand,
				Variables:    config.Variables,
			}
		}
	case services.BuildpackStatic:
		buildRequest.BuildpackType = build.Static
		if config := buildConfig.StaticConfig(); config != nil {
			buildRequest.StaticConfig = &build.StaticConfig{
				BuildCommand: config.BuildCommand,
				OutputDir:    config.OutputDir,
				NginxConfig:  config.NginxConfig,
			}
		}
	case services.BuildpackDockerfile:
		buildRequest.BuildpackType = build.DockerfileType
		if config := buildConfig.DockerfileConfig(); config != nil {
			buildRequest.DockerfileConfig = &build.DockerfileConfig{
				DockerfilePath: config.DockerfilePath,
				BuildArgs:      config.BuildArgs,
				Target:         config.Target,
			}
		}
	case services.BuildpackDockerCompose:
		buildRequest.BuildpackType = build.DockerCompose
		if config := buildConfig.ComposeConfig(); config != nil {
			buildRequest.ComposeConfig = &build.ComposeConfig{
				ComposeFile: config.ComposeFile,
				Service:     config.Service,
			}
		}
	}

	return buildRequest
}

func (s *ServiceService) deployContainer(ctx context.Context, svc *services.Service, imageTag string) error {
	containerConfig := manager.ContainerConfig{
		Image:       imageTag,
		Name:        s.getContainerName(svc),
		Environment: svc.Environment(),
		// TODO: Add port mappings, volumes, etc. based on services configuration
		RestartPolicy: "unless-stopped",
	}

	// Create and start the container
	containerID, err := s.containerManager.Create(ctx, containerConfig)
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	if err := s.containerManager.Start(ctx, containerID); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	return nil
}

func (s *ServiceService) getContainerName(svc *services.Service) string {
	// Create a unique container name based on project, environment, and services
	return fmt.Sprintf("mikrocloud-%s-%s-%s",
		svc.ProjectID().String()[:8],
		svc.EnvironmentID().String()[:8],
		svc.Name().String())
}

func (s *ServiceService) getImageTag(svc *services.Service) string {
	// Create image tag based on services
	return fmt.Sprintf("mikrocloud/%s:%s", svc.Name().String(), "latest")
}
