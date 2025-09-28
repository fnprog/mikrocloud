package handlers

import (
	"context"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/mikrocloud/mikrocloud/internal/database"
	"github.com/mikrocloud/mikrocloud/internal/domain/environments"
	"github.com/mikrocloud/mikrocloud/internal/domain/projects"
	"github.com/mikrocloud/mikrocloud/internal/domain/services"
	servicesService "github.com/mikrocloud/mikrocloud/internal/domain/services/service"
	"github.com/mikrocloud/mikrocloud/internal/domain/users"
)

type ServiceHandler struct {
	db             *database.Database
	serviceService *servicesService.ServiceService
}

func NewServiceHandler(db *database.Database, serviceService *servicesService.ServiceService) *ServiceHandler {
	return &ServiceHandler{
		db:             db,
		serviceService: serviceService,
	}
}

// QuickServiceInput for one-off service creation
type QuickServiceInput struct {
	Body struct {
		Name          string                   `json:"name" minLength:"1" maxLength:"50" example:"my-service"`
		GitURL        string                   `json:"git_url" example:"https://github.com/user/repo.git"`
		GitBranch     string                   `json:"git_branch,omitempty" example:"main"`
		ContextRoot   string                   `json:"context_root,omitempty" example:"frontend/"`
		BuildpackType string                   `json:"buildpack_type" example:"nixpacks" enum:"nixpacks,static,dockerfile,docker-compose"`
		Environment   map[string]string        `json:"environment,omitempty"`
		BuildConfig   *ServiceBuildConfigInput `json:"build_config,omitempty"`
	}
}

type ServiceBuildConfigInput struct {
	Nixpacks   *NixpacksConfigInput   `json:"nixpacks,omitempty"`
	Static     *StaticConfigInput     `json:"static,omitempty"`
	Dockerfile *DockerfileConfigInput `json:"dockerfile,omitempty"`
	Compose    *ComposeConfigInput    `json:"compose,omitempty"`
}

type NixpacksConfigInput struct {
	StartCommand string            `json:"start_command,omitempty"`
	BuildCommand string            `json:"build_command,omitempty"`
	Variables    map[string]string `json:"variables,omitempty"`
}

type StaticConfigInput struct {
	BuildCommand string `json:"build_command,omitempty"`
	OutputDir    string `json:"output_dir,omitempty"`
	NginxConfig  string `json:"nginx_config,omitempty"`
}

type DockerfileConfigInput struct {
	DockerfilePath string            `json:"dockerfile_path,omitempty"`
	BuildArgs      map[string]string `json:"build_args,omitempty"`
	Target         string            `json:"target,omitempty"`
}

type ComposeConfigInput struct {
	ComposeFile string `json:"compose_file,omitempty"`
	Service     string `json:"service,omitempty"`
}

type QuickServiceOutput struct {
	Body struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		ProjectID   string `json:"project_id"`
		ProjectName string `json:"project_name"`
		Environment string `json:"environment"`
		GitURL      string `json:"git_url"`
		Status      string `json:"status"`
		CreatedAt   string `json:"created_at"`
	}
}

// CreateServiceInput for creating services in existing project/environment
type CreateServiceInput struct {
	ProjectID     string `path:"project_id"`
	EnvironmentID string `path:"environment_id"`
	Body          struct {
		Name          string                   `json:"name" minLength:"1" maxLength:"50"`
		GitURL        string                   `json:"git_url"`
		GitBranch     string                   `json:"git_branch,omitempty"`
		ContextRoot   string                   `json:"context_root,omitempty"`
		BuildpackType string                   `json:"buildpack_type" enum:"nixpacks,static,dockerfile,docker-compose"`
		Environment   map[string]string        `json:"environment,omitempty"`
		BuildConfig   *ServiceBuildConfigInput `json:"build_config,omitempty"`
	}
}

type ServiceOutput struct {
	Body struct {
		ID            string            `json:"id"`
		Name          string            `json:"name"`
		ProjectID     string            `json:"project_id"`
		EnvironmentID string            `json:"environment_id"`
		GitURL        string            `json:"git_url"`
		GitBranch     string            `json:"git_branch"`
		ContextRoot   string            `json:"context_root"`
		BuildpackType string            `json:"buildpack_type"`
		Environment   map[string]string `json:"environment"`
		Status        string            `json:"status"`
		CreatedAt     string            `json:"created_at"`
		UpdatedAt     string            `json:"updated_at"`
	}
}

type ServiceListItem struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	ProjectID     string `json:"project_id"`
	EnvironmentID string `json:"environment_id"`
	GitURL        string `json:"git_url"`
	Status        string `json:"status"`
	CreatedAt     string `json:"created_at"`
}

type ListServicesOutput struct {
	Body struct {
		Services []ServiceListItem `json:"services"`
	}
}

type ServiceActionInput struct {
	ProjectID     string `path:"project_id"`
	EnvironmentID string `path:"environment_id"`
	ServiceID     string `path:"service_id"`
}

type ServiceActionOutput struct {
	Body struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}
}

// CreateQuickService creates a new service with auto-generated project and prod environment
func (h *ServiceHandler) CreateQuickService(ctx context.Context, input *QuickServiceInput) (*QuickServiceOutput, error) {
	// Create project with same name as service (or add suffix)
	projectName, err := projects.NewProjectName(input.Body.Name)
	if err != nil {
		return nil, huma.Error400BadRequest("invalid project name", err)
	}

	// Create service name
	serviceName, err := services.NewServiceName(input.Body.Name)
	if err != nil {
		return nil, huma.Error400BadRequest("invalid service name", err)
	}

	// Create project
	// TODO: Get actual user ID from authentication context
	userID, _ := users.UserIDFromString("00000000-0000-0000-0000-000000000000")
	orgID, _ := users.OrganizationIDFromString("00000000-0000-0000-0000-000000000000")
	proj := projects.NewProject(projectName, fmt.Sprintf("Auto-generated project for %s", input.Body.Name), userID, orgID, userID)

	// Create default prod environment
	env := environments.NewEnvironment(environments.EnvironmentProduction, proj.ID().UUID(), "Production environment", true)

	// Parse Git URL
	gitURL, err := services.NewGitURL(input.Body.GitURL, input.Body.GitBranch, input.Body.ContextRoot)
	if err != nil {
		return nil, huma.Error400BadRequest("invalid git URL", err)
	}

	// Create build config
	buildConfig := h.createBuildConfig(input.Body.BuildpackType, input.Body.BuildConfig)

	// Create service
	svc := services.NewService(serviceName, proj.ID().UUID(), env.ID().UUID(), gitURL, buildConfig)

	// Set environment variables
	for key, value := range input.Body.Environment {
		if err := svc.SetEnvironmentVariable(key, value); err != nil {
			return nil, huma.Error400BadRequest("invalid environment variable", err)
		}
	}

	// TODO: Save to database through service layer
	// This would involve:
	// 1. Saving the project
	// 2. Saving the environment
	// 3. Saving the service
	// For now, we'll return a mock response

	return &QuickServiceOutput{
		Body: struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			ProjectID   string `json:"project_id"`
			ProjectName string `json:"project_name"`
			Environment string `json:"environment"`
			GitURL      string `json:"git_url"`
			Status      string `json:"status"`
			CreatedAt   string `json:"created_at"`
		}{
			ID:          svc.ID().String(),
			Name:        svc.Name().String(),
			ProjectID:   proj.ID().String(),
			ProjectName: proj.Name().String(),
			Environment: env.Name().String(),
			GitURL:      gitURL.URL(),
			Status:      svc.Status().String(),
			CreatedAt:   svc.CreatedAt().Format("2006-01-02T15:04:05Z"),
		},
	}, nil
}

// CreateService creates a new service in an existing project/environment
func (h *ServiceHandler) CreateService(ctx context.Context, input *CreateServiceInput) (*ServiceOutput, error) {
	// Validate project and environment exist
	projectID, err := uuid.Parse(input.ProjectID)
	if err != nil {
		return nil, huma.Error400BadRequest("invalid project ID", err)
	}

	environmentID, err := uuid.Parse(input.EnvironmentID)
	if err != nil {
		return nil, huma.Error400BadRequest("invalid environment ID", err)
	}

	// Create service name
	serviceName, err := services.NewServiceName(input.Body.Name)
	if err != nil {
		return nil, huma.Error400BadRequest("invalid service name", err)
	}

	// Parse Git URL
	gitURL, err := services.NewGitURL(input.Body.GitURL, input.Body.GitBranch, input.Body.ContextRoot)
	if err != nil {
		return nil, huma.Error400BadRequest("invalid git URL", err)
	}

	// Create build config
	buildConfig := h.createBuildConfig(input.Body.BuildpackType, input.Body.BuildConfig)

	// Create service
	svc := services.NewService(serviceName, projectID, environmentID, gitURL, buildConfig)

	// Set environment variables
	for key, value := range input.Body.Environment {
		if err := svc.SetEnvironmentVariable(key, value); err != nil {
			return nil, huma.Error400BadRequest("invalid environment variable", err)
		}
	}

	// TODO: Save to database through service layer

	return &ServiceOutput{
		Body: struct {
			ID            string            `json:"id"`
			Name          string            `json:"name"`
			ProjectID     string            `json:"project_id"`
			EnvironmentID string            `json:"environment_id"`
			GitURL        string            `json:"git_url"`
			GitBranch     string            `json:"git_branch"`
			ContextRoot   string            `json:"context_root"`
			BuildpackType string            `json:"buildpack_type"`
			Environment   map[string]string `json:"environment"`
			Status        string            `json:"status"`
			CreatedAt     string            `json:"created_at"`
			UpdatedAt     string            `json:"updated_at"`
		}{
			ID:            svc.ID().String(),
			Name:          svc.Name().String(),
			ProjectID:     svc.ProjectID().String(),
			EnvironmentID: svc.EnvironmentID().String(),
			GitURL:        gitURL.URL(),
			GitBranch:     gitURL.Branch(),
			ContextRoot:   gitURL.ContextRoot(),
			BuildpackType: string(buildConfig.BuildpackType()),
			Environment:   svc.Environment(),
			Status:        svc.Status().String(),
			CreatedAt:     svc.CreatedAt().Format("2006-01-02T15:04:05Z"),
			UpdatedAt:     svc.UpdatedAt().Format("2006-01-02T15:04:05Z"),
		},
	}, nil
}

// GetService retrieves a specific service
func (h *ServiceHandler) GetService(ctx context.Context, input *ServiceActionInput) (*ServiceOutput, error) {
	// TODO: Implement database lookup
	return nil, huma.Error501NotImplemented("not yet implemented")
}

// ListServices lists services in a project/environment
func (h *ServiceHandler) ListServices(ctx context.Context, input *struct {
	ProjectID     string `path:"project_id"`
	EnvironmentID string `path:"environment_id"`
}) (*ListServicesOutput, error) {
	// TODO: Implement database lookup
	return &ListServicesOutput{
		Body: struct {
			Services []ServiceListItem `json:"services"`
		}{
			Services: []ServiceListItem{},
		},
	}, nil
}

// DeployService triggers deployment of a service
func (h *ServiceHandler) DeployService(ctx context.Context, input *ServiceActionInput) (*ServiceActionOutput, error) {
	// TODO: Implement deployment through container service
	return &ServiceActionOutput{
		Body: struct {
			Message string `json:"message"`
			Status  string `json:"status"`
		}{
			Message: "Service deployment started",
			Status:  "deploying",
		},
	}, nil
}

// StopService stops a running service
func (h *ServiceHandler) StopService(ctx context.Context, input *ServiceActionInput) (*ServiceActionOutput, error) {
	// TODO: Implement stop through container service
	return &ServiceActionOutput{
		Body: struct {
			Message string `json:"message"`
			Status  string `json:"status"`
		}{
			Message: "Service stopped",
			Status:  "stopped",
		},
	}, nil
}

// RestartService restarts a service
func (h *ServiceHandler) RestartService(ctx context.Context, input *ServiceActionInput) (*ServiceActionOutput, error) {
	// TODO: Implement restart through container service
	return &ServiceActionOutput{
		Body: struct {
			Message string `json:"message"`
			Status  string `json:"status"`
		}{
			Message: "Service restarted",
			Status:  "running",
		},
	}, nil
}

// DeleteService deletes a service
func (h *ServiceHandler) DeleteService(ctx context.Context, input *ServiceActionInput) (*ServiceActionOutput, error) {
	// TODO: Implement delete through service layer
	return &ServiceActionOutput{
		Body: struct {
			Message string `json:"message"`
			Status  string `json:"status"`
		}{
			Message: "Service deleted",
			Status:  "deleted",
		},
	}, nil
}

// Helper method to create build config from input
func (h *ServiceHandler) createBuildConfig(buildpackType string, configInput *ServiceBuildConfigInput) *services.BuildConfig {
	buildConfig := services.NewBuildConfig(services.BuildpackType(buildpackType))

	if configInput == nil {
		return buildConfig
	}

	switch buildpackType {
	case "nixpacks":
		if configInput.Nixpacks != nil {
			buildConfig.SetNixpacksConfig(&services.NixpacksConfig{
				StartCommand: configInput.Nixpacks.StartCommand,
				BuildCommand: configInput.Nixpacks.BuildCommand,
				Variables:    configInput.Nixpacks.Variables,
			})
		}
	case "static":
		if configInput.Static != nil {
			buildConfig.SetStaticConfig(&services.StaticConfig{
				BuildCommand: configInput.Static.BuildCommand,
				OutputDir:    configInput.Static.OutputDir,
				NginxConfig:  configInput.Static.NginxConfig,
			})
		}
	case "dockerfile":
		if configInput.Dockerfile != nil {
			buildConfig.SetDockerfileConfig(&services.DockerfileConfig{
				DockerfilePath: configInput.Dockerfile.DockerfilePath,
				BuildArgs:      configInput.Dockerfile.BuildArgs,
				Target:         configInput.Dockerfile.Target,
			})
		}
	case "docker-compose":
		if configInput.Compose != nil {
			buildConfig.SetComposeConfig(&services.ComposeConfig{
				ComposeFile: configInput.Compose.ComposeFile,
				Service:     configInput.Compose.Service,
			})
		}
	}

	return buildConfig
}
