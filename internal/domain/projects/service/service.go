package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/mikrocloud/mikrocloud/internal/domain/activities"
	activitiesService "github.com/mikrocloud/mikrocloud/internal/domain/activities/service"
	"github.com/mikrocloud/mikrocloud/internal/domain/environments"
	envRepo "github.com/mikrocloud/mikrocloud/internal/domain/environments/repository"
	"github.com/mikrocloud/mikrocloud/internal/domain/projects"
	"github.com/mikrocloud/mikrocloud/internal/domain/projects/repository"
	"github.com/mikrocloud/mikrocloud/internal/domain/users"
)

// ProjectService handles projects-related business operations
type ProjectService struct {
	projectRepo       repository.Repository
	envRepo           envRepo.Repository
	activitiesService *activitiesService.ActivitiesService
}

// NewProjectService creates a new projects service
func NewProjectService(projectRepo repository.Repository, envRepo envRepo.Repository, activitiesSvc *activitiesService.ActivitiesService) *ProjectService {
	return &ProjectService{
		projectRepo:       projectRepo,
		envRepo:           envRepo,
		activitiesService: activitiesSvc,
	}
}

// CreateProjectCommand represents the data needed to create a projects
type CreateProjectCommand struct {
	Name           string
	Description    *string
	UserID         string
	OrganisationID string
}

// CreateProject creates a new projects following business rules
func (s *ProjectService) CreateProject(ctx context.Context, cmd CreateProjectCommand) (*projects.Project, error) {
	// Validate projects name
	projectName, err := projects.NewProjectName(cmd.Name)
	if err != nil {
		return nil, fmt.Errorf("invalid projects name: %w", err)
	}

	// Check if projects with name already exists
	exists, err := s.projectRepo.Exists(ctx, projectName)
	if err != nil {
		return nil, fmt.Errorf("failed to check projects existence: %w", err)
	}

	if exists {
		return nil, fmt.Errorf("projects '%s' already exists", cmd.Name)
	}

	orgID, err := users.OrganizationIDFromString(cmd.OrganisationID)
	if err != nil {
		return nil, fmt.Errorf("failed to translate organisationID: %w", err)
	}

	userID, err := users.UserIDFromString(cmd.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to translate userID: %w", err)
	}

	// Create the projects
	proj := projects.NewProject(projectName, cmd.Description, userID, orgID)

	// Save the projects
	if err := s.projectRepo.Save(ctx, proj); err != nil {
		return nil, fmt.Errorf("failed to save projects: %w", err)
	}

	// Create default "production" environment for the project
	envName, err := environments.NewEnvironmentName("production")
	if err != nil {
		return nil, fmt.Errorf("failed to create default environment name: %w", err)
	}

	env := environments.NewEnvironment(envName, proj.ID().UUID(), "Default production environment", true)
	if err := s.envRepo.Save(ctx, env); err != nil {
		return nil, fmt.Errorf("failed to create default environment: %w", err)
	}

	// Log activity
	if s.activitiesService != nil {
		projIDUUID := proj.ID().UUID()
		projNameStr := proj.Name().String()
		resourceType := "project"
		userUUID, _ := uuid.Parse(userID.String())
		orgUUID, _ := uuid.Parse(orgID.String())
		_ = s.activitiesService.LogActivity(
			activities.EventTypeProjectCreated,
			fmt.Sprintf("Project '%s' created", projNameStr),
			&userUUID,
			&resourceType,
			&projIDUUID,
			&projNameStr,
			nil,
			orgUUID,
		)
	}

	return proj, nil
}

// ListProjects retrieves all projects
func (s *ProjectService) ListProjects(ctx context.Context, orgID users.OrganizationID) ([]*projects.Project, error) {
	return s.projectRepo.FindAll(ctx, orgID)
}

// GetProject retrieves a projects by ID
func (s *ProjectService) GetProject(ctx context.Context, id string) (*projects.Project, error) {
	projectID, err := projects.ProjectIDFromString(id)
	if err != nil {
		return nil, fmt.Errorf("invalid projects identifier: %w", err)
	}

	return s.projectRepo.FindByID(ctx, projectID)
}

// GetProjectByName retrieves a projects by name
func (s *ProjectService) GetProjectByName(ctx context.Context, name string) (*projects.Project, error) {
	projectName, err := projects.NewProjectName(name)
	if err != nil {
		return nil, fmt.Errorf("invalid projects name: %w", err)
	}

	return s.projectRepo.FindByName(ctx, projectName)
}

// UpdateProject updates an existing projects
func (s *ProjectService) UpdateProject(ctx context.Context, id string, description *string) (*projects.Project, error) {
	// Get existing projects
	proj, err := s.GetProject(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("projects not found: %w", err)
	}

	// Update description
	proj.UpdateDescription(description)

	// Save the updated projects
	if err := s.projectRepo.Save(ctx, proj); err != nil {
		return nil, fmt.Errorf("failed to update projects: %w", err)
	}

	return proj, nil
}

// DeleteProject removes a projects
func (s *ProjectService) DeleteProject(ctx context.Context, id string) error {
	// Get existing projects to validate it exists
	proj, err := s.GetProject(ctx, id)
	if err != nil {
		return fmt.Errorf("projects not found: %w", err)
	}

	// Business rule: Cannot delete default projects
	if proj.Name().String() == "default" {
		return fmt.Errorf("cannot delete default projects")
	}

	// TODO: Business rule: Cannot delete projects with applications
	// This would require checking if any applications exist in this projects

	return s.projectRepo.Delete(ctx, proj.ID())
}
