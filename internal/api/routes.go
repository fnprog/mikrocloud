package api

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/mikrocloud/mikrocloud/internal/api/handlers"
	"github.com/mikrocloud/mikrocloud/internal/config"
	"github.com/mikrocloud/mikrocloud/internal/container"
	"github.com/mikrocloud/mikrocloud/internal/database"
)

func SetupRoutes(api huma.API, db *database.Database, cfg *config.Config) error {
	// Create service container
	serviceContainer, err := container.NewContainer(db, cfg)
	if err != nil {
		return err
	}

	// Create handler dependencies
	projectHandler := handlers.NewProjectHandler(db)
	environmentHandler := handlers.NewEnvironmentHandler(db)
	serviceHandler := handlers.NewServiceHandler(db, serviceContainer.ServiceService)

	// Health check
	huma.Get(api, "/health", handlers.HealthCheck)

	// Project routes
	huma.Get(api, "/v1/projects", projectHandler.ListProjects)
	huma.Post(api, "/v1/projects", projectHandler.CreateProject)
	huma.Get(api, "/v1/projects/{project_id}", projectHandler.GetProject)
	huma.Delete(api, "/v1/projects/{project_id}", projectHandler.DeleteProject)

	// Environment routes within projects
	huma.Get(api, "/v1/projects/{project_id}/environments", environmentHandler.ListEnvironments)
	huma.Post(api, "/v1/projects/{project_id}/environments", environmentHandler.CreateEnvironment)
	huma.Get(api, "/v1/projects/{project_id}/environments/{environment_id}", environmentHandler.GetEnvironment)
	huma.Put(api, "/v1/projects/{project_id}/environments/{environment_id}", environmentHandler.UpdateEnvironment)
	huma.Delete(api, "/v1/projects/{project_id}/environments/{environment_id}", environmentHandler.DeleteEnvironment)

	// Service routes within project environments
	huma.Get(api, "/v1/projects/{project_id}/environments/{environment_id}/services", serviceHandler.ListServices)
	huma.Post(api, "/v1/projects/{project_id}/environments/{environment_id}/services", serviceHandler.CreateService)
	huma.Get(api, "/v1/projects/{project_id}/environments/{environment_id}/services/{service_id}", serviceHandler.GetService)
	huma.Delete(api, "/v1/projects/{project_id}/environments/{environment_id}/services/{service_id}", serviceHandler.DeleteService)

	// Service action routes
	huma.Post(api, "/v1/projects/{project_id}/environments/{environment_id}/services/{service_id}/deploy", serviceHandler.DeployService)
	huma.Post(api, "/v1/projects/{project_id}/environments/{environment_id}/services/{service_id}/stop", serviceHandler.StopService)
	huma.Post(api, "/v1/projects/{project_id}/environments/{environment_id}/services/{service_id}/restart", serviceHandler.RestartService)

	// Quick service creation (creates project + prod env + service in one call)
	huma.Post(api, "/v1/services/quick", serviceHandler.CreateQuickService)

	// System routes
	huma.Get(api, "/v1/system/status", handlers.SystemStatus)
	huma.Get(api, "/v1/system/info", handlers.SystemInfo)

	// Domain management routes (placeholder)
	huma.Get(api, "/v1/domains", handlers.ListDomains)
	huma.Post(api, "/v1/domains", handlers.AddDomain)
	huma.Delete(api, "/v1/domains/{domain}", handlers.RemoveDomain)
	huma.Post(api, "/v1/domains/{domain}/ssl", handlers.EnableSSL)

	// Auth routes (placeholder)
	huma.Post(api, "/v1/auth/login", handlers.Login)
	huma.Post(api, "/v1/auth/register", handlers.Register)
	huma.Post(api, "/v1/auth/logout", handlers.Logout)

	return nil
}
