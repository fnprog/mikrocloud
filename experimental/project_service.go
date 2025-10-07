package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/mikrocloud/mikrocloud/internal/domain"
	"github.com/mikrocloud/mikrocloud/pkg/cloudprovider"
)

// ProjectService handles project management operations
type ProjectService struct {
	repos         domain.RepositoryManager
	cloudFactory  *cloudprovider.CloudProviderFactory
	deploymentSvc *DeploymentService
}

// NewProjectService creates a new project service
func NewProjectService(
	repos domain.RepositoryManager,
	cloudFactory *cloudprovider.CloudProviderFactory,
	deploymentSvc *DeploymentService,
) *ProjectService {
	return &ProjectService{
		repos:         repos,
		cloudFactory:  cloudFactory,
		deploymentSvc: deploymentSvc,
	}
}

// CreateProjectRequest represents a project creation request
type CreateProjectRequest struct {
	UserID        uuid.UUID            `json:"user_id" validate:"required"`
	Name          string               `json:"name" validate:"required,min=1,max=50"`
	GitRepo       string               `json:"git_repo" validate:"url"`
	Framework     string               `json:"framework" validate:"required"`
	CloudProvider domain.CloudProvider `json:"cloud_provider" validate:"required,oneof=aws azure gcp"`
	Region        string               `json:"region" validate:"required"`
	Environment   map[string]string    `json:"environment,omitempty"`
}

// ProjectResponse represents a project response
type ProjectResponse struct {
	ID            uuid.UUID                  `json:"id"`
	UserID        uuid.UUID                  `json:"user_id"`
	Name          string                     `json:"name"`
	GitRepo       string                     `json:"git_repo"`
	Framework     string                     `json:"framework"`
	CloudProvider domain.CloudProviderConfig `json:"cloud_provider"`
	Status        string                     `json:"status"`
	ResourceIDs   map[string]interface{}     `json:"resource_ids"`
	CreatedAt     string                     `json:"created_at"`
	UpdatedAt     string                     `json:"updated_at"`
}

// CreateProject creates a new project with cloud resources
func (s *ProjectService) CreateProject(ctx context.Context, req *CreateProjectRequest) (*ProjectResponse, error) {
	// Validate user exists
	_, err := s.repos.User().GetByID(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Get or create cloud provider config
	cloudConfig, err := s.getOrCreateCloudConfig(ctx, req.UserID, req.CloudProvider, req.Region)
	if err != nil {
		return nil, fmt.Errorf("failed to get cloud config: %w", err)
	}

	// Create project in database
	project := &domain.Project{
		UserID:          req.UserID,
		CloudProviderID: cloudConfig.ID,
		Name:            req.Name,
		GitRepo:         req.GitRepo,
		Framework:       req.Framework,
		Status:          "creating",
		ResourceIDs:     "{}",
	}

	if err := s.repos.Project().Create(ctx, project); err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	// Set up environment variables if provided
	if len(req.Environment) > 0 {
		if err := s.setEnvironmentVariables(ctx, project.ID, req.Environment, "production"); err != nil {
			return nil, fmt.Errorf("failed to set environment variables: %w", err)
		}
	}

	// Initialize cloud resources asynchronously
	go func() {
		if err := s.initializeCloudResources(context.Background(), project, cloudConfig); err != nil {
			// Log error and update project status
			project.Status = "failed"
			s.repos.Project().Update(context.Background(), project)
		}
	}()

	return s.projectToResponse(project, cloudConfig), nil
}

// GetProject retrieves a project by ID
func (s *ProjectService) GetProject(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (*ProjectResponse, error) {
	project, err := s.repos.Project().GetByID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	// Verify ownership
	if project.UserID != userID {
		return nil, fmt.Errorf("access denied")
	}

	return s.projectToResponse(project, &project.CloudProvider), nil
}

// ListProjects retrieves projects for a user
func (s *ProjectService) ListProjects(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*ProjectResponse, error) {
	projects, err := s.repos.Project().GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}

	responses := make([]*ProjectResponse, len(projects))
	for i, project := range projects {
		responses[i] = s.projectToResponse(project, &project.CloudProvider)
	}

	return responses, nil
}

// UpdateProject updates project configuration
func (s *ProjectService) UpdateProject(ctx context.Context, projectID uuid.UUID, userID uuid.UUID, updates map[string]interface{}) (*ProjectResponse, error) {
	project, err := s.repos.Project().GetByID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	// Verify ownership
	if project.UserID != userID {
		return nil, fmt.Errorf("access denied")
	}

	// Apply updates
	if name, ok := updates["name"].(string); ok {
		project.Name = name
	}
	if gitRepo, ok := updates["git_repo"].(string); ok {
		project.GitRepo = gitRepo
	}
	if framework, ok := updates["framework"].(string); ok {
		project.Framework = framework
	}

	if err := s.repos.Project().Update(ctx, project); err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	return s.projectToResponse(project, &project.CloudProvider), nil
}

// DeleteProject deletes a project and cleans up cloud resources
func (s *ProjectService) DeleteProject(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) error {
	project, err := s.repos.Project().GetByID(ctx, projectID)
	if err != nil {
		return fmt.Errorf("project not found: %w", err)
	}

	// Verify ownership
	if project.UserID != userID {
		return fmt.Errorf("access denied")
	}

	// Clean up cloud resources
	if err := s.cleanupCloudResources(ctx, project); err != nil {
		return fmt.Errorf("failed to cleanup cloud resources: %w", err)
	}

	// Delete from database
	if err := s.repos.Project().Delete(ctx, projectID); err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	return nil
}

// Helper methods

func (s *ProjectService) getOrCreateCloudConfig(ctx context.Context, userID uuid.UUID, provider domain.CloudProvider, region string) (*domain.CloudProviderConfig, error) {
	// Try to get existing config
	config, err := s.repos.CloudProvider().GetByUserAndProvider(ctx, userID, provider)
	if err == nil {
		return config, nil
	}

	// Create new config with placeholder credentials
	config = &domain.CloudProviderConfig{
		UserID:               userID,
		Provider:             provider,
		Region:               region,
		CredentialsEncrypted: "placeholder-credentials", // TODO: Implement proper credential management
		IsDefault:            false,
	}

	if err := s.repos.CloudProvider().Create(ctx, config); err != nil {
		return nil, err
	}

	return config, nil
}

func (s *ProjectService) setEnvironmentVariables(ctx context.Context, projectID uuid.UUID, envVars map[string]string, environment string) error {
	for key, value := range envVars {
		envVar := &domain.EnvironmentVariable{
			ProjectID:      projectID,
			Key:            key,
			ValueEncrypted: value, // TODO: Implement encryption
			Environment:    environment,
		}

		if err := s.repos.EnvironmentVariable().Create(ctx, envVar); err != nil {
			return err
		}
	}
	return nil
}

func (s *ProjectService) initializeCloudResources(ctx context.Context, project *domain.Project, cloudConfig *domain.CloudProviderConfig) error {
	// Create cloud provider instance
	credentials := map[string]string{
		"region": cloudConfig.Region,
		// TODO: Decrypt and provide actual credentials
	}

	provider, err := s.cloudFactory.Create(cloudConfig.Provider, credentials)
	if err != nil {
		return fmt.Errorf("failed to create cloud provider: %w", err)
	}

	// Set up networking
	networkConfig := cloudprovider.NetworkConfig{
		ProjectID:        project.ID.String(),
		VPCName:          fmt.Sprintf("%s-vpc", project.Name),
		CIDRBlock:        "10.0.0.0/16",
		PublicSubnets:    []string{"10.0.1.0/24", "10.0.2.0/24"},
		PrivateSubnets:   []string{"10.0.10.0/24", "10.0.20.0/24"},
		EnableNATGateway: true,
	}

	network, err := provider.SetupNetworking(ctx, networkConfig)
	if err != nil {
		return fmt.Errorf("failed to setup networking: %w", err)
	}

	// Update project with resource IDs
	resourceIDs := map[string]interface{}{
		"vpc_id":          network.VPCID,
		"public_subnets":  network.PublicSubnetIDs,
		"private_subnets": network.PrivateSubnetIDs,
	}

	project.ResourceIDs = fmt.Sprintf("%v", resourceIDs) // TODO: Proper JSON serialization
	project.Status = "active"

	return s.repos.Project().Update(ctx, project)
}

func (s *ProjectService) cleanupCloudResources(ctx context.Context, project *domain.Project) error {
	// TODO: Implement cloud resource cleanup
	// This would involve:
	// 1. Parse resource IDs from project
	// 2. Create cloud provider instance
	// 3. Delete all associated resources (VPC, subnets, etc.)
	return nil
}

func (s *ProjectService) projectToResponse(project *domain.Project, cloudConfig *domain.CloudProviderConfig) *ProjectResponse {
	return &ProjectResponse{
		ID:            project.ID,
		UserID:        project.UserID,
		Name:          project.Name,
		GitRepo:       project.GitRepo,
		Framework:     project.Framework,
		CloudProvider: *cloudConfig,
		Status:        project.Status,
		ResourceIDs:   make(map[string]interface{}), // TODO: Parse JSON
		CreatedAt:     project.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:     project.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
