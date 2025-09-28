package service

import (
	"context"
	"fmt"

	"github.com/mikrocloud/mikrocloud/internal/domain/applications"
	"github.com/mikrocloud/mikrocloud/internal/domain/deployments"
	"github.com/mikrocloud/mikrocloud/internal/domain/users"
)

type DeploymentRepository interface {
	Create(ctx context.Context, deployment *deployments.Deployment) error
	GetByID(ctx context.Context, id deployments.DeploymentID) (*deployments.Deployment, error)
	Update(ctx context.Context, deployment *deployments.Deployment) error
	Delete(ctx context.Context, id deployments.DeploymentID) error
	List(ctx context.Context) ([]*deployments.Deployment, error)
	ListByApplication(ctx context.Context, applicationID applications.ApplicationID) ([]*deployments.Deployment, error)
	GetLatestByApplication(ctx context.Context, applicationID applications.ApplicationID) (*deployments.Deployment, error)
	ListByStatus(ctx context.Context, status deployments.DeploymentStatus) ([]*deployments.Deployment, error)
}

type DeploymentService struct {
	repo DeploymentRepository
}

func NewDeploymentService(repo DeploymentRepository) *DeploymentService {
	return &DeploymentService{
		repo: repo,
	}
}

type CreateDeploymentCommand struct {
	ApplicationID    applications.ApplicationID
	IsProduction     bool
	TriggeredBy      *users.UserID
	TriggerType      deployments.TriggerType
	ImageTag         string
	GitCommitHash    string
	GitCommitMessage string
	GitBranch        string
	GitAuthorName    string
}

func (s *DeploymentService) CreateDeployment(ctx context.Context, cmd CreateDeploymentCommand) (*deployments.Deployment, error) {
	// Get the next deployment number for this application
	deploymentNumber, err := s.getNextDeploymentNumber(ctx, cmd.ApplicationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get next deployment number: %w", err)
	}

	deployment := deployments.NewDeployment(
		cmd.ApplicationID,
		deploymentNumber,
		cmd.IsProduction,
		cmd.TriggeredBy,
		cmd.TriggerType,
		cmd.ImageTag,
	)

	if cmd.GitCommitHash != "" || cmd.GitCommitMessage != "" || cmd.GitBranch != "" || cmd.GitAuthorName != "" {
		deployment.SetGitInfo(cmd.GitCommitHash, cmd.GitCommitMessage, cmd.GitBranch, cmd.GitAuthorName)
	}

	if err := s.repo.Create(ctx, deployment); err != nil {
		return nil, fmt.Errorf("failed to create deployment: %w", err)
	}

	return deployment, nil
}

func (s *DeploymentService) getNextDeploymentNumber(ctx context.Context, applicationID applications.ApplicationID) (int, error) {
	latest, err := s.repo.GetLatestByApplication(ctx, applicationID)
	if err != nil {
		// If no deployments exist, start with 1
		return 1, nil
	}
	return latest.DeploymentNumber() + 1, nil
}

func (s *DeploymentService) GetDeployment(ctx context.Context, id deployments.DeploymentID) (*deployments.Deployment, error) {
	deployment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment: %w", err)
	}
	return deployment, nil
}

func (s *DeploymentService) StartBuild(ctx context.Context, id deployments.DeploymentID) error {
	deployment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("deployment not found: %w", err)
	}

	deployment.StartBuild()

	if err := s.repo.Update(ctx, deployment); err != nil {
		return fmt.Errorf("failed to update deployment: %w", err)
	}

	return nil
}

func (s *DeploymentService) CompleteBuild(ctx context.Context, id deployments.DeploymentID) error {
	deployment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("deployment not found: %w", err)
	}

	deployment.CompleteBuild()

	if err := s.repo.Update(ctx, deployment); err != nil {
		return fmt.Errorf("failed to update deployment: %w", err)
	}

	return nil
}

func (s *DeploymentService) StartDeploy(ctx context.Context, id deployments.DeploymentID) error {
	deployment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("deployment not found: %w", err)
	}

	deployment.StartDeploy()

	if err := s.repo.Update(ctx, deployment); err != nil {
		return fmt.Errorf("failed to update deployment: %w", err)
	}

	return nil
}

func (s *DeploymentService) CompleteDeploy(ctx context.Context, id deployments.DeploymentID) error {
	deployment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("deployment not found: %w", err)
	}

	deployment.CompleteDeploy()

	if err := s.repo.Update(ctx, deployment); err != nil {
		return fmt.Errorf("failed to update deployment: %w", err)
	}

	return nil
}

func (s *DeploymentService) FailDeployment(ctx context.Context, id deployments.DeploymentID, errorMessage string) error {
	deployment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("deployment not found: %w", err)
	}

	deployment.Fail(errorMessage)

	if err := s.repo.Update(ctx, deployment); err != nil {
		return fmt.Errorf("failed to update deployment: %w", err)
	}

	return nil
}

func (s *DeploymentService) CancelDeployment(ctx context.Context, id deployments.DeploymentID) error {
	deployment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("deployment not found: %w", err)
	}

	deployment.Cancel()

	if err := s.repo.Update(ctx, deployment); err != nil {
		return fmt.Errorf("failed to update deployment: %w", err)
	}

	return nil
}

func (s *DeploymentService) StopDeployment(ctx context.Context, id deployments.DeploymentID) error {
	deployment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("deployment not found: %w", err)
	}

	deployment.Stop()

	if err := s.repo.Update(ctx, deployment); err != nil {
		return fmt.Errorf("failed to update deployment: %w", err)
	}

	return nil
}

func (s *DeploymentService) SetContainerID(ctx context.Context, id deployments.DeploymentID, containerID string) error {
	deployment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("deployment not found: %w", err)
	}

	deployment.SetContainerID(containerID)

	if err := s.repo.Update(ctx, deployment); err != nil {
		return fmt.Errorf("failed to update deployment: %w", err)
	}

	return nil
}

func (s *DeploymentService) SetImageDigest(ctx context.Context, id deployments.DeploymentID, digest string) error {
	deployment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("deployment not found: %w", err)
	}

	deployment.SetImageDigest(digest)

	if err := s.repo.Update(ctx, deployment); err != nil {
		return fmt.Errorf("failed to update deployment: %w", err)
	}

	return nil
}

func (s *DeploymentService) AppendBuildLogs(ctx context.Context, id deployments.DeploymentID, logs string) error {
	deployment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("deployment not found: %w", err)
	}

	deployment.AppendBuildLogs(logs)

	if err := s.repo.Update(ctx, deployment); err != nil {
		return fmt.Errorf("failed to update deployment: %w", err)
	}

	return nil
}

func (s *DeploymentService) AppendDeployLogs(ctx context.Context, id deployments.DeploymentID, logs string) error {
	deployment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("deployment not found: %w", err)
	}

	deployment.AppendDeployLogs(logs)

	if err := s.repo.Update(ctx, deployment); err != nil {
		return fmt.Errorf("failed to update deployment: %w", err)
	}

	return nil
}

func (s *DeploymentService) ListDeployments(ctx context.Context) ([]*deployments.Deployment, error) {
	deployments, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments: %w", err)
	}
	return deployments, nil
}

func (s *DeploymentService) ListDeploymentsByApplication(ctx context.Context, applicationID applications.ApplicationID) ([]*deployments.Deployment, error) {
	deployments, err := s.repo.ListByApplication(ctx, applicationID)
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments by application: %w", err)
	}
	return deployments, nil
}

func (s *DeploymentService) GetLatestDeploymentByApplication(ctx context.Context, applicationID applications.ApplicationID) (*deployments.Deployment, error) {
	deployment, err := s.repo.GetLatestByApplication(ctx, applicationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest deployment: %w", err)
	}
	return deployment, nil
}

func (s *DeploymentService) ListDeploymentsByStatus(ctx context.Context, status deployments.DeploymentStatus) ([]*deployments.Deployment, error) {
	deployments, err := s.repo.ListByStatus(ctx, status)
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments by status: %w", err)
	}
	return deployments, nil
}

func (s *DeploymentService) DeleteDeployment(ctx context.Context, id deployments.DeploymentID) error {
	// Check if deployment exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("deployment not found: %w", err)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete deployment: %w", err)
	}

	return nil
}
