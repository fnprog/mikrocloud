package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mikrocloud/mikrocloud/internal/domain"
	"github.com/mikrocloud/mikrocloud/pkg/cloudprovider"
)

// DeploymentService handles deployment operations
type DeploymentService struct {
	repos        domain.RepositoryManager
	cloudFactory *cloudprovider.CloudProviderFactory
}

// NewDeploymentService creates a new deployment service
func NewDeploymentService(
	repos domain.RepositoryManager,
	cloudFactory *cloudprovider.CloudProviderFactory,
) *DeploymentService {
	return &DeploymentService{
		repos:        repos,
		cloudFactory: cloudFactory,
	}
}

// DeploymentRequest represents a deployment request
type DeploymentRequest struct {
	ProjectID   uuid.UUID   `json:"project_id" validate:"required"`
	CommitSHA   string      `json:"commit_sha" validate:"required"`
	Environment string      `json:"environment" validate:"required,oneof=production staging preview"`
	BuildConfig BuildConfig `json:"build_config,omitempty"`
}

type BuildConfig struct {
	BuildCommand   string            `json:"build_command,omitempty"`
	BuildDir       string            `json:"build_dir,omitempty"`
	OutputDir      string            `json:"output_dir,omitempty"`
	Environment    map[string]string `json:"environment,omitempty"`
	InstallCommand string            `json:"install_command,omitempty"`
}

// DeploymentResponse represents a deployment response
type DeploymentResponse struct {
	ID            uuid.UUID            `json:"id"`
	ProjectID     uuid.UUID            `json:"project_id"`
	CommitSHA     string               `json:"commit_sha"`
	Status        string               `json:"status"`
	Environment   string               `json:"environment"`
	URL           string               `json:"url"`
	CloudProvider domain.CloudProvider `json:"cloud_provider"`
	BuildLogs     string               `json:"build_logs,omitempty"`
	DeployedAt    *time.Time           `json:"deployed_at,omitempty"`
	CreatedAt     time.Time            `json:"created_at"`
	UpdatedAt     time.Time            `json:"updated_at"`
}

// TriggerDeployment starts a new deployment
func (s *DeploymentService) TriggerDeployment(ctx context.Context, userID uuid.UUID, req *DeploymentRequest) (*DeploymentResponse, error) {
	// Get project and verify ownership
	project, err := s.repos.Project().GetByID(ctx, req.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	if project.UserID != userID {
		return nil, fmt.Errorf("access denied")
	}

	// Create deployment record
	deployment := &domain.Deployment{
		ProjectID:     req.ProjectID,
		CommitSHA:     req.CommitSHA,
		Status:        "pending",
		Environment:   req.Environment,
		CloudProvider: project.CloudProvider.Provider,
		BuildLogs:     "",
	}

	if err := s.repos.Deployment().Create(ctx, deployment); err != nil {
		return nil, fmt.Errorf("failed to create deployment: %w", err)
	}

	// Start deployment process asynchronously
	go func() {
		if err := s.executeDeployment(context.Background(), deployment, project, req.BuildConfig); err != nil {
			deployment.Status = "failed"
			deployment.BuildLogs += fmt.Sprintf("\nDeployment failed: %v", err)
			s.repos.Deployment().Update(context.Background(), deployment)
		}
	}()

	return s.deploymentToResponse(deployment), nil
}

// GetDeployment retrieves a deployment by ID
func (s *DeploymentService) GetDeployment(ctx context.Context, deploymentID uuid.UUID, userID uuid.UUID) (*DeploymentResponse, error) {
	deployment, err := s.repos.Deployment().GetByID(ctx, deploymentID)
	if err != nil {
		return nil, fmt.Errorf("deployment not found: %w", err)
	}

	// Verify project ownership
	project, err := s.repos.Project().GetByID(ctx, deployment.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	if project.UserID != userID {
		return nil, fmt.Errorf("access denied")
	}

	return s.deploymentToResponse(deployment), nil
}

// ListDeployments retrieves deployments for a project
func (s *DeploymentService) ListDeployments(ctx context.Context, projectID uuid.UUID, userID uuid.UUID, limit, offset int) ([]*DeploymentResponse, error) {
	// Verify project ownership
	project, err := s.repos.Project().GetByID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	if project.UserID != userID {
		return nil, fmt.Errorf("access denied")
	}

	deployments, err := s.repos.Deployment().GetByProjectID(ctx, projectID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments: %w", err)
	}

	responses := make([]*DeploymentResponse, len(deployments))
	for i, deployment := range deployments {
		responses[i] = s.deploymentToResponse(deployment)
	}

	return responses, nil
}

// GetDeploymentLogs retrieves deployment logs
func (s *DeploymentService) GetDeploymentLogs(ctx context.Context, deploymentID uuid.UUID, userID uuid.UUID) (string, error) {
	deployment, err := s.repos.Deployment().GetByID(ctx, deploymentID)
	if err != nil {
		return "", fmt.Errorf("deployment not found: %w", err)
	}

	// Verify project ownership
	project, err := s.repos.Project().GetByID(ctx, deployment.ProjectID)
	if err != nil {
		return "", fmt.Errorf("project not found: %w", err)
	}

	if project.UserID != userID {
		return "", fmt.Errorf("access denied")
	}

	return deployment.BuildLogs, nil
}

// PromoteDeployment promotes a deployment to production
func (s *DeploymentService) PromoteDeployment(ctx context.Context, deploymentID uuid.UUID, userID uuid.UUID) (*DeploymentResponse, error) {
	deployment, err := s.repos.Deployment().GetByID(ctx, deploymentID)
	if err != nil {
		return nil, fmt.Errorf("deployment not found: %w", err)
	}

	// Verify project ownership
	project, err := s.repos.Project().GetByID(ctx, deployment.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	if project.UserID != userID {
		return nil, fmt.Errorf("access denied")
	}

	if deployment.Status != "success" {
		return nil, fmt.Errorf("can only promote successful deployments")
	}

	if deployment.Environment == "production" {
		return nil, fmt.Errorf("deployment is already in production")
	}

	// Create new production deployment
	prodDeployment := &domain.Deployment{
		ProjectID:     deployment.ProjectID,
		CommitSHA:     deployment.CommitSHA,
		Status:        "pending",
		Environment:   "production",
		CloudProvider: deployment.CloudProvider,
		BuildLogs:     "Promoted from " + deployment.Environment,
	}

	if err := s.repos.Deployment().Create(ctx, prodDeployment); err != nil {
		return nil, fmt.Errorf("failed to create production deployment: %w", err)
	}

	// Start promotion process asynchronously
	go func() {
		if err := s.executePromotion(context.Background(), prodDeployment, deployment, project); err != nil {
			prodDeployment.Status = "failed"
			prodDeployment.BuildLogs += fmt.Sprintf("\nPromotion failed: %v", err)
			s.repos.Deployment().Update(context.Background(), prodDeployment)
		}
	}()

	return s.deploymentToResponse(prodDeployment), nil
}

// DeleteDeployment deletes a deployment
func (s *DeploymentService) DeleteDeployment(ctx context.Context, deploymentID uuid.UUID, userID uuid.UUID) error {
	deployment, err := s.repos.Deployment().GetByID(ctx, deploymentID)
	if err != nil {
		return fmt.Errorf("deployment not found: %w", err)
	}

	// Verify project ownership
	project, err := s.repos.Project().GetByID(ctx, deployment.ProjectID)
	if err != nil {
		return fmt.Errorf("project not found: %w", err)
	}

	if project.UserID != userID {
		return fmt.Errorf("access denied")
	}

	// Don't allow deletion of production deployments
	if deployment.Environment == "production" {
		return fmt.Errorf("cannot delete production deployment")
	}

	// Clean up cloud resources if needed
	if err := s.cleanupDeploymentResources(ctx, deployment, project); err != nil {
		return fmt.Errorf("failed to cleanup deployment resources: %w", err)
	}

	return s.repos.Deployment().Delete(ctx, deploymentID)
}

// Helper methods

func (s *DeploymentService) executeDeployment(ctx context.Context, deployment *domain.Deployment, project *domain.Project, buildConfig BuildConfig) error {
	// Update status to building
	deployment.Status = "building"
	deployment.BuildLogs = "Starting deployment...\n"
	s.repos.Deployment().Update(ctx, deployment)

	// Get cloud provider
	provider, err := s.getCloudProvider(ctx, project)
	if err != nil {
		return fmt.Errorf("failed to get cloud provider: %w", err)
	}

	// Determine deployment type based on framework
	deploymentType := s.determineDeploymentType(project.Framework)

	switch deploymentType {
	case "static":
		return s.deployStaticSite(ctx, deployment, project, provider, buildConfig)
	case "serverless":
		return s.deployServerlessApp(ctx, deployment, project, provider, buildConfig)
	case "container":
		return s.deployContainerApp(ctx, deployment, project, provider, buildConfig)
	default:
		return fmt.Errorf("unsupported deployment type: %s", deploymentType)
	}
}

func (s *DeploymentService) deployStaticSite(ctx context.Context, deployment *domain.Deployment, project *domain.Project, provider cloudprovider.CloudProvider, buildConfig BuildConfig) error {
	deployment.BuildLogs += "Deploying static site...\n"
	s.repos.Deployment().Update(ctx, deployment)

	config := cloudprovider.StaticSiteConfig{
		ProjectID:    project.ID.String(),
		Name:         project.Name,
		BuildCommand: buildConfig.BuildCommand,
		BuildDir:     buildConfig.OutputDir,
		Environment:  buildConfig.Environment,
	}

	result, err := provider.DeployStaticSite(ctx, config)
	if err != nil {
		return err
	}

	// Update deployment with success
	deployment.Status = "success"
	deployment.URL = result.URL
	now := time.Now()
	deployment.DeployedAt = &now
	deployment.BuildLogs += fmt.Sprintf("Deployment successful! URL: %s\n", result.URL)

	return s.repos.Deployment().Update(ctx, deployment)
}

func (s *DeploymentService) deployServerlessApp(ctx context.Context, deployment *domain.Deployment, project *domain.Project, provider cloudprovider.CloudProvider, buildConfig BuildConfig) error {
	deployment.BuildLogs += "Deploying serverless application...\n"
	s.repos.Deployment().Update(ctx, deployment)

	// For now, create a simple function deployment
	config := cloudprovider.FunctionConfig{
		ProjectID:   project.ID.String(),
		Name:        "main",
		Runtime:     s.getRuntimeFromFramework(project.Framework),
		Handler:     "index.handler",
		Code:        []byte("placeholder code"), // TODO: Get actual code
		Environment: buildConfig.Environment,
		Timeout:     30,
		MemorySize:  256,
	}

	result, err := provider.DeployServerlessFunction(ctx, config)
	if err != nil {
		return err
	}

	// Update deployment with success
	deployment.Status = "success"
	deployment.URL = result.URL
	now := time.Now()
	deployment.DeployedAt = &now
	deployment.BuildLogs += fmt.Sprintf("Deployment successful! URL: %s\n", result.URL)

	return s.repos.Deployment().Update(ctx, deployment)
}

func (s *DeploymentService) deployContainerApp(ctx context.Context, deployment *domain.Deployment, project *domain.Project, provider cloudprovider.CloudProvider, buildConfig BuildConfig) error {
	deployment.BuildLogs += "Deploying container application...\n"
	s.repos.Deployment().Update(ctx, deployment)

	config := cloudprovider.ContainerConfig{
		ProjectID:   project.ID.String(),
		Name:        project.Name,
		Image:       fmt.Sprintf("%s:latest", project.Name), // TODO: Build actual image
		Port:        8080,
		Environment: buildConfig.Environment,
		CPU:         0.5,
		Memory:      512,
		Instances:   1,
	}

	result, err := provider.DeployContainer(ctx, config)
	if err != nil {
		return err
	}

	// Update deployment with success
	deployment.Status = "success"
	deployment.URL = result.URL
	now := time.Now()
	deployment.DeployedAt = &now
	deployment.BuildLogs += fmt.Sprintf("Deployment successful! URL: %s\n", result.URL)

	return s.repos.Deployment().Update(ctx, deployment)
}

func (s *DeploymentService) executePromotion(ctx context.Context, prodDeployment *domain.Deployment, sourceDeployment *domain.Deployment, project *domain.Project) error {
	// For promotion, we typically just copy the configuration and redeploy to production
	// This is a simplified implementation
	prodDeployment.Status = "success"
	prodDeployment.URL = sourceDeployment.URL // TODO: Update to production URL
	now := time.Now()
	prodDeployment.DeployedAt = &now
	prodDeployment.BuildLogs += "Promotion completed successfully\n"

	return s.repos.Deployment().Update(ctx, prodDeployment)
}

func (s *DeploymentService) cleanupDeploymentResources(ctx context.Context, deployment *domain.Deployment, project *domain.Project) error {
	// TODO: Implement cleanup of cloud resources for the specific deployment
	return nil
}

func (s *DeploymentService) getCloudProvider(ctx context.Context, project *domain.Project) (cloudprovider.CloudProvider, error) {
	// Get cloud configuration
	cloudConfig, err := s.repos.CloudProvider().GetByID(ctx, project.CloudProviderID)
	if err != nil {
		return nil, err
	}

	// Create credentials map
	credentials := map[string]string{
		"region": cloudConfig.Region,
		// TODO: Decrypt actual credentials
	}

	return s.cloudFactory.Create(cloudConfig.Provider, credentials)
}

func (s *DeploymentService) determineDeploymentType(framework string) string {
	staticFrameworks := map[string]bool{
		"nextjs-static": true,
		"react":         true,
		"vue":           true,
		"angular":       true,
		"svelte":        true,
		"gatsby":        true,
		"nuxt-static":   true,
	}

	containerFrameworks := map[string]bool{
		"docker": true,
		"python": true,
		"go":     true,
		"java":   true,
		"dotnet": true,
		"php":    true,
		"ruby":   true,
	}

	if staticFrameworks[framework] {
		return "static"
	}

	if containerFrameworks[framework] {
		return "container"
	}

	// Default to serverless for Node.js and similar
	return "serverless"
}

func (s *DeploymentService) getRuntimeFromFramework(framework string) string {
	runtimeMap := map[string]string{
		"nextjs": "nodejs18.x",
		"nodejs": "nodejs18.x",
		"python": "python3.9",
		"go":     "provided.al2",
		"java":   "java17",
		"dotnet": "dotnet6",
	}

	if runtime, ok := runtimeMap[framework]; ok {
		return runtime
	}

	return "nodejs18.x" // default
}

func (s *DeploymentService) deploymentToResponse(deployment *domain.Deployment) *DeploymentResponse {
	var deployedAt *time.Time
	if deployment.DeployedAt != nil && !deployment.DeployedAt.IsZero() {
		deployedAt = deployment.DeployedAt
	}

	return &DeploymentResponse{
		ID:          deployment.ID,
		ProjectID:   deployment.ProjectID,
		CommitSHA:   deployment.CommitSHA,
		Status:      deployment.Status,
		Environment: deployment.Environment,
		URL:         deployment.URL,
		BuildLogs:   deployment.BuildLogs,
		DeployedAt:  deployedAt,
		CreatedAt:   deployment.CreatedAt,
		UpdatedAt:   deployment.UpdatedAt,
	}
}

// CreateDeploymentRequest represents a deployment creation request
type CreateDeploymentRequest struct {
	ProjectID   uuid.UUID         `json:"project_id" validate:"required"`
	CommitSHA   string            `json:"commit_sha" validate:"required"`
	Environment string            `json:"environment" validate:"required,oneof=production staging preview"`
	BuildConfig map[string]string `json:"build_config,omitempty"`
}

// CreateDeployment creates a new deployment
func (s *DeploymentService) CreateDeployment(ctx context.Context, req *CreateDeploymentRequest, userID uuid.UUID) (*DeploymentResponse, error) {
	// Verify project exists and user has access
	project, err := s.repos.Project().GetByID(ctx, req.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	if project.UserID != userID {
		return nil, fmt.Errorf("access denied")
	}

	// Create deployment record
	deployment := &domain.Deployment{
		ProjectID:   req.ProjectID,
		CommitSHA:   req.CommitSHA,
		Status:      "building",
		Environment: req.Environment,
		BuildLogs:   "",
		URL:         "",
	}

	if err := s.repos.Deployment().Create(ctx, deployment); err != nil {
		return nil, fmt.Errorf("failed to create deployment: %w", err)
	}

	// Start deployment process asynchronously
	go func() {
		s.executeDeployment(context.Background(), deployment, project, BuildConfig{})
	}()

	return s.deploymentToResponse(deployment), nil
}

// RestartDeployment restarts a deployment
func (s *DeploymentService) RestartDeployment(ctx context.Context, deploymentID uuid.UUID, userID uuid.UUID) error {
	deployment, err := s.repos.Deployment().GetByID(ctx, deploymentID)
	if err != nil {
		return fmt.Errorf("deployment not found: %w", err)
	}

	// Verify access through project
	project, err := s.repos.Project().GetByID(ctx, deployment.ProjectID)
	if err != nil {
		return fmt.Errorf("project not found: %w", err)
	}

	if project.UserID != userID {
		return fmt.Errorf("access denied")
	}

	// TODO: Implement actual restart logic with cloud provider
	deployment.Status = "restarting"
	if err := s.repos.Deployment().Update(ctx, deployment); err != nil {
		return fmt.Errorf("failed to update deployment: %w", err)
	}

	// Simulate restart process
	go func() {
		time.Sleep(5 * time.Second)
		deployment.Status = "active"
		s.repos.Deployment().Update(context.Background(), deployment)
	}()

	return nil
}

// ScaleDeployment scales a deployment
func (s *DeploymentService) ScaleDeployment(ctx context.Context, deploymentID uuid.UUID, userID uuid.UUID, replicas int) error {
	deployment, err := s.repos.Deployment().GetByID(ctx, deploymentID)
	if err != nil {
		return fmt.Errorf("deployment not found: %w", err)
	}

	// Verify access through project
	project, err := s.repos.Project().GetByID(ctx, deployment.ProjectID)
	if err != nil {
		return fmt.Errorf("project not found: %w", err)
	}

	if project.UserID != userID {
		return fmt.Errorf("access denied")
	}

	// TODO: Implement actual scaling logic with cloud provider
	// For container deployments, this would scale ECS services, Azure Container Instances, etc.
	
	deployment.Status = "scaling"
	if err := s.repos.Deployment().Update(ctx, deployment); err != nil {
		return fmt.Errorf("failed to update deployment: %w", err)
	}

	// Simulate scaling process
	go func() {
		time.Sleep(3 * time.Second)
		deployment.Status = "active"
		s.repos.Deployment().Update(context.Background(), deployment)
	}()

	return nil
}
