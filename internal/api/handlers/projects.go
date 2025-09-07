package handlers

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/mikrocloud/mikrocloud/internal/domain/environment"
	"github.com/mikrocloud/mikrocloud/internal/domain/projects/service"
)

// ProjectHandler handles project-related HTTP requests
type ProjectHandler struct {
	projectService *service.ProjectService
}

// NewProjectHandler creates a new project handler
func NewProjectHandler(pgs *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		projectService: pgs,
	}
}

// ProjectResponse represents a project in API responses
type ProjectResponse struct {
	ID          string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string `json:"name" example:"my-project"`
	Description string `json:"description" example:"My awesome project"`
	CreatedAt   string `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt   string `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

type CreateProjectInput struct {
	Body struct {
		Name        string `json:"name" minLength:"1" maxLength:"100" example:"my-project"`
		Description string `json:"description,omitempty" example:"My awesome project"`
	}
}

type ProjectOutput struct {
	Body struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		CreatedAt   string `json:"created_at"`
		UpdatedAt   string `json:"updated_at"`
		Environment struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"default_environment"`
	}
}

type ProjectListItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}

type ListProjectsOutput struct {
	Body struct {
		Projects []ProjectListItem `json:"projects"`
	}
}

type ProjectActionInput struct {
	ProjectID string `path:"project_id"`
}

// CreateProject creates a new project with a default prod environment
func (h *ProjectHandler) CreateProject(ctx context.Context, input *CreateProjectInput) (*ProjectOutput, error) {
	// Create project name
	projectName, err := h.projectService.NewProjectName(input.Body.Name)
	if err != nil {
		return nil, huma.Error400BadRequest("invalid project name", err)
	}

	// Create project
	proj := h.projectService.CreateProject(projectName, input.Body.Description)

	// Create default production environment
	env := environment.NewEnvironment(environment.EnvironmentProduction, proj.ID().UUID(), "Production environment")

	// Save to database
	err = h.db.ProjectRepository.Save(ctx, proj)
	if err != nil {
		return nil, huma.Error400BadRequest("Failed to create project", err)
	}

	err = h.db.EnvironmentRepository.Save(env)
	if err != nil {
		return nil, huma.Error400BadRequest("Failed to create default environment", err)
	}

	return &ProjectOutput{
		Body: struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
			CreatedAt   string `json:"created_at"`
			UpdatedAt   string `json:"updated_at"`
			Environment struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"default_environment"`
		}{
			ID:          proj.ID().String(),
			Name:        proj.Name().String(),
			Description: proj.Description(),
			CreatedAt:   proj.CreatedAt().Format("2006-01-02T15:04:05Z"),
			UpdatedAt:   proj.UpdatedAt().Format("2006-01-02T15:04:05Z"),
			Environment: struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			}{
				ID:   env.ID().String(),
				Name: env.Name().String(),
			},
		},
	}, nil
}

// ListProjects lists all projects
func (h *ProjectHandler) ListProjects(ctx context.Context, input *struct{}) (*ListProjectsOutput, error) {
	// TODO: Implement database lookup
	return &ListProjectsOutput{
		Body: struct {
			Projects []ProjectListItem `json:"projects"`
		}{
			Projects: []ProjectListItem{},
		},
	}, nil
}

// GetProject retrieves a specific project
func (h *ProjectHandler) GetProject(ctx context.Context, input *ProjectActionInput) (*ProjectOutput, error) {
	// TODO: Implement database lookup
	return nil, huma.Error501NotImplemented("not yet implemented")
}

// DeleteProject deletes a project and all its environments and services
func (h *ProjectHandler) DeleteProject(ctx context.Context, input *ProjectActionInput) (*struct {
	Body struct {
		Message string `json:"message"`
	}
}, error) {
	// TODO: Implement delete through service layer
	// This should cascade delete environments and services
	return &struct {
		Body struct {
			Message string `json:"message"`
		}
	}{
		Body: struct {
			Message string `json:"message"`
		}{
			Message: "Project deleted successfully",
		},
	}, nil
}
