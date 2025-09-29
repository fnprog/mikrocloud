package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/mikrocloud/mikrocloud/internal/domain/applications"
)

type ApplicationRepository interface {
	Save(ctx context.Context, app *applications.Application) error
	FindByID(ctx context.Context, id applications.ApplicationID) (*applications.Application, error)
	FindByName(ctx context.Context, projectID uuid.UUID, name applications.ApplicationName) (*applications.Application, error)
	FindByProject(ctx context.Context, projectID uuid.UUID) ([]*applications.Application, error)
	FindByEnvironment(ctx context.Context, environmentID uuid.UUID) ([]*applications.Application, error)
	FindAll(ctx context.Context) ([]*applications.Application, error)
	Delete(ctx context.Context, id applications.ApplicationID) error
	Exists(ctx context.Context, projectID uuid.UUID, name applications.ApplicationName) (bool, error)
}

type ApplicationService struct {
	repo ApplicationRepository
}

func NewApplicationService(repo ApplicationRepository) *ApplicationService {
	return &ApplicationService{
		repo: repo,
	}
}

type CreateApplicationCommand struct {
	Name             string
	Description      string
	ProjectID        uuid.UUID
	EnvironmentID    uuid.UUID
	DeploymentSource applications.DeploymentSource
	BuildpackConfig  applications.BuildpackConfig
	EnvVars          map[string]string
}

func (s *ApplicationService) CreateApplication(ctx context.Context, cmd CreateApplicationCommand) (*applications.Application, error) {
	name, err := applications.NewApplicationName(cmd.Name)
	if err != nil {
		return nil, fmt.Errorf("invalid application name: %w", err)
	}

	// Check if application already exists
	exists, err := s.repo.Exists(ctx, cmd.ProjectID, name)
	if err != nil {
		return nil, fmt.Errorf("failed to check if application exists: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("application with name %s already exists in project", name.String())
	}

	app := applications.NewApplication(
		name,
		cmd.Description,
		cmd.ProjectID,
		cmd.EnvironmentID,
		cmd.DeploymentSource,
		cmd.BuildpackConfig,
	)

	// Set env vars if provided
	if cmd.EnvVars != nil {
		app.SetEnvVars(cmd.EnvVars)
	}

	if err := s.repo.Save(ctx, app); err != nil {
		return nil, fmt.Errorf("failed to create application: %w", err)
	}

	return app, nil
}

func (s *ApplicationService) GetApplication(ctx context.Context, id applications.ApplicationID) (*applications.Application, error) {
	app, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get application: %w", err)
	}
	return app, nil
}

type UpdateApplicationCommand struct {
	ID               applications.ApplicationID
	Description      *string
	DeploymentSource *applications.DeploymentSource
	RepoURL          *string // For backward compatibility
	RepoBranch       *string // For backward compatibility
	RepoPath         *string // For backward compatibility
	Domain           *string
	BuildpackType    *applications.BuildpackType // For backward compatibility
	BuildpackConfig  *applications.BuildpackConfig
	Config           *string // For backward compatibility
	EnvVars          map[string]string
	AutoDeploy       *bool
}

func (s *ApplicationService) UpdateApplication(ctx context.Context, cmd UpdateApplicationCommand) (*applications.Application, error) {
	app, err := s.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("application not found: %w", err)
	}

	if cmd.Description != nil {
		app.UpdateDescription(*cmd.Description)
	}

	// Handle deployment source updates
	if cmd.DeploymentSource != nil {
		app.SetDeploymentSource(*cmd.DeploymentSource)
	} else {
		// Handle backward compatibility fields
		if cmd.RepoURL != nil {
			app.SetRepoURL(*cmd.RepoURL)
		}
		if cmd.RepoBranch != nil {
			app.SetRepoBranch(*cmd.RepoBranch)
		}
		if cmd.RepoPath != nil {
			app.SetRepoPath(*cmd.RepoPath)
		}
	}

	if cmd.Domain != nil {
		app.SetDomain(*cmd.Domain)
	}

	// Handle buildpack updates
	if cmd.BuildpackConfig != nil {
		app.SetBuildpack(*cmd.BuildpackConfig)
	} else {
		// Handle backward compatibility fields
		if cmd.BuildpackType != nil {
			app.SetBuildpackType(*cmd.BuildpackType)
		}
		if cmd.Config != nil {
			app.UpdateConfig(*cmd.Config)
		}
	}

	// Handle env vars
	if cmd.EnvVars != nil {
		app.SetEnvVars(cmd.EnvVars)
	}

	if cmd.AutoDeploy != nil {
		app.SetAutoDeploy(*cmd.AutoDeploy)
	}

	if err := s.repo.Save(ctx, app); err != nil {
		return nil, fmt.Errorf("failed to update application: %w", err)
	}

	return app, nil
}

func (s *ApplicationService) DeleteApplication(ctx context.Context, id applications.ApplicationID) error {
	// Check if application exists
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("application not found: %w", err)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete application: %w", err)
	}

	return nil
}

func (s *ApplicationService) ListApplications(ctx context.Context) ([]*applications.Application, error) {
	apps, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list applications: %w", err)
	}
	return apps, nil
}

func (s *ApplicationService) ListApplicationsByProject(ctx context.Context, projectID uuid.UUID) ([]*applications.Application, error) {
	apps, err := s.repo.FindByProject(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list applications by project: %w", err)
	}
	return apps, nil
}

func (s *ApplicationService) ListApplicationsByEnvironment(ctx context.Context, environmentID uuid.UUID) ([]*applications.Application, error) {
	apps, err := s.repo.FindByEnvironment(ctx, environmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to list applications by environment: %w", err)
	}
	return apps, nil
}

func (s *ApplicationService) StartDeployment(ctx context.Context, id applications.ApplicationID) error {
	app, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("application not found: %w", err)
	}

	if err := app.CanDeploy(); err != nil {
		return fmt.Errorf("cannot deploy application: %w", err)
	}

	app.ChangeStatus(applications.ApplicationStatusDeploying)

	if err := s.repo.Save(ctx, app); err != nil {
		return fmt.Errorf("failed to update application status: %w", err)
	}

	return nil
}

func (s *ApplicationService) StopApplication(ctx context.Context, id applications.ApplicationID) error {
	app, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("application not found: %w", err)
	}

	if err := app.CanStop(); err != nil {
		return fmt.Errorf("cannot stop application: %w", err)
	}

	app.ChangeStatus(applications.ApplicationStatusStopped)

	if err := s.repo.Save(ctx, app); err != nil {
		return fmt.Errorf("failed to update application status: %w", err)
	}

	return nil
}

func (s *ApplicationService) UpdateApplicationStatus(ctx context.Context, id applications.ApplicationID, status applications.ApplicationStatus) error {
	app, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("application not found: %w", err)
	}

	app.ChangeStatus(status)

	if err := s.repo.Save(ctx, app); err != nil {
		return fmt.Errorf("failed to update application status: %w", err)
	}

	return nil
}
