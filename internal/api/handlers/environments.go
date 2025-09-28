package handlers

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/mikrocloud/mikrocloud/internal/database"
	"github.com/mikrocloud/mikrocloud/internal/domain/environments"
)

type EnvironmentHandler struct {
	db *database.Database
}

func NewEnvironmentHandler(db *database.Database) *EnvironmentHandler {
	return &EnvironmentHandler{
		db: db,
	}
}

// CreateEnvironmentInput for creating environments in a project
type CreateEnvironmentInput struct {
	ProjectID string `path:"project_id"`
	Body      struct {
		Name        string            `json:"name" minLength:"1" maxLength:"50" example:"dev"`
		Description string            `json:"description,omitempty" example:"Development environment"`
		Variables   map[string]string `json:"variables,omitempty"`
	}
}

type EnvironmentOutput struct {
	Body struct {
		ID          string            `json:"id"`
		Name        string            `json:"name"`
		ProjectID   string            `json:"project_id"`
		Description string            `json:"description"`
		Variables   map[string]string `json:"variables"`
		CreatedAt   string            `json:"created_at"`
		UpdatedAt   string            `json:"updated_at"`
	}
}

type EnvironmentListItem struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	ProjectID   string            `json:"project_id"`
	Description string            `json:"description"`
	Variables   map[string]string `json:"variables"`
	CreatedAt   string            `json:"created_at"`
}

type ListEnvironmentsOutput struct {
	Body struct {
		Environments []EnvironmentListItem `json:"environments"`
	}
}

type EnvironmentActionInput struct {
	ProjectID     string `path:"project_id"`
	EnvironmentID string `path:"environment_id"`
}

type UpdateEnvironmentInput struct {
	ProjectID     string `path:"project_id"`
	EnvironmentID string `path:"environment_id"`
	Body          struct {
		Description string            `json:"description,omitempty"`
		Variables   map[string]string `json:"variables,omitempty"`
	}
}

// CreateEnvironment creates a new environment in a project
func (h *EnvironmentHandler) CreateEnvironment(ctx context.Context, input *CreateEnvironmentInput) (*EnvironmentOutput, error) {
	// Validate project ID
	projectID, err := uuid.Parse(input.ProjectID)
	if err != nil {
		return nil, huma.Error400BadRequest("invalid project ID", err)
	}

	// Create environment name
	envName, err := environments.NewEnvironmentName(input.Body.Name)
	if err != nil {
		return nil, huma.Error400BadRequest("invalid environment name", err)
	}

	// Create environment
	env := environments.NewEnvironment(envName, projectID, input.Body.Description, false)

	// Set variables
	for key, value := range input.Body.Variables {
		if err := env.SetVariable(key, value); err != nil {
			return nil, huma.Error400BadRequest("invalid environment variable", err)
		}
	}

	// TODO: Save to database through service layer

	return &EnvironmentOutput{
		Body: struct {
			ID          string            `json:"id"`
			Name        string            `json:"name"`
			ProjectID   string            `json:"project_id"`
			Description string            `json:"description"`
			Variables   map[string]string `json:"variables"`
			CreatedAt   string            `json:"created_at"`
			UpdatedAt   string            `json:"updated_at"`
		}{
			ID:          env.ID().String(),
			Name:        env.Name().String(),
			ProjectID:   env.ProjectID().String(),
			Description: env.Description(),
			Variables:   env.Variables(),
			CreatedAt:   env.CreatedAt().Format("2006-01-02T15:04:05Z"),
			UpdatedAt:   env.UpdatedAt().Format("2006-01-02T15:04:05Z"),
		},
	}, nil
}

// GetEnvironment retrieves a specific environment
func (h *EnvironmentHandler) GetEnvironment(ctx context.Context, input *EnvironmentActionInput) (*EnvironmentOutput, error) {
	// TODO: Implement database lookup
	return nil, huma.Error501NotImplemented("not yet implemented")
}

// ListEnvironments lists all environments in a project
func (h *EnvironmentHandler) ListEnvironments(ctx context.Context, input *struct {
	ProjectID string `path:"project_id"`
}) (*ListEnvironmentsOutput, error) {
	// TODO: Implement database lookup
	return &ListEnvironmentsOutput{
		Body: struct {
			Environments []EnvironmentListItem `json:"environments"`
		}{
			Environments: []EnvironmentListItem{},
		},
	}, nil
}

// UpdateEnvironment updates an environment
func (h *EnvironmentHandler) UpdateEnvironment(ctx context.Context, input *UpdateEnvironmentInput) (*EnvironmentOutput, error) {
	// TODO: Implement update through service layer
	return nil, huma.Error501NotImplemented("not yet implemented")
}

// DeleteEnvironment deletes an environment
func (h *EnvironmentHandler) DeleteEnvironment(ctx context.Context, input *EnvironmentActionInput) (*struct {
	Body struct {
		Message string `json:"message"`
	}
}, error) {
	// TODO: Implement delete through service layer
	return &struct {
		Body struct {
			Message string `json:"message"`
		}
	}{
		Body: struct {
			Message string `json:"message"`
		}{
			Message: "Environment deleted successfully",
		},
	}, nil
}
