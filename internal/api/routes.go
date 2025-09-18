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

	// Maintenance
	huma.Get(api, "/health", handlers.HealthCheck)
	huma.Get(api, "/system/status", handlers.SystemStatus)
	huma.Get(api, "/system/info", handlers.SystemInfo)
	huma.Get(api, "/resources", handlers.GetResources) // list all ressources (dbs, servers, etc....)

	// Auth routes
	huma.Post(api, "/auth/login", handlers.Login)
	huma.Post(api, "/auth/register", handlers.Register)
	huma.Post(api, "/auth/logout", handlers.Logout)

	// Server
	huma.Get(api, "/servers", serverHandler.List)
	huma.Post(api, "/servers", serverHandler.Create)
	huma.Put(api, "/servers/{server_id}", serverHandler.Update)
	huma.Delete(api, "/servers/{server_id}", serverHandler.Delete)
	huma.Get(api, "/servers/{server_id}", serverHandler.Get)
	huma.Get(api, "/servers/{server_id}/domains", serverHandler.GetDomains)           // domain associated with the server
	huma.Get(api, "/servers/{server_id}/destinations", serverHandler.GetDestinations) // Docker/podman, etc...
	huma.Get(api, "/servers/{server_id}/proxy", serverHandler.GetDestinations)        // Server proxy settings
	huma.Get(api, "/servers/{server_id}/resources", serverHandler.GetResources)       // hardware resource data for a server
	huma.Get(api, "/servers/{server_id}/validate", serverHandler.Validate)            // validate server connectivity
	huma.Get(api, "/servers/{server_id}/logs", serverHandler.Validate)                // server logs

	// Project routes
	huma.Get(api, "/projects", projectHandler.List)
	huma.Post(api, "/projects", projectHandler.Create)
	huma.Get(api, "/projects/{project_id}", projectHandler.Get)
	huma.Put(api, "/projects/{project_id}", projectHandler.Update)
	huma.Delete(api, "/projects/{project_id}", projectHandler.Delete)
	// Quick creation from a docker compose file (all items in the docker file are components of the project)
	huma.Post(api, "/projects/quick", serviceHandler.CreateQuickService)

	// Environment routes within projects
	huma.Get(api, "/projects/{project_id}/environments", environmentHandler.List)
	huma.Post(api, "/projects/{project_id}/environments", environmentHandler.Create)
	huma.Get(api, "/projects/{project_id}/environments/{environment_id}", environmentHandler.Get)
	huma.Put(api, "/projects/{project_id}/environments/{environment_id}", environmentHandler.Update)
	huma.Delete(api, "/projects/{project_id}/environments/{environment_id}", environmentHandler.Delete)

	// Applications routes within projects
	huma.Get(api, "/projects/{project_id}/environments/{environment_id}/applications", applicationHandler.List)
	huma.Post(api, "/projects/{project_id}/environments/{environment_id}/applications", applicationHandler.Create)

	// Service routes within project environments
	huma.Get(api, "/projects/{project_id}/environments/{environment_id}/services", serviceHandler.List)
	huma.Post(api, "/projects/{project_id}/environments/{environment_id}/services", serviceHandler.Create)
	huma.Get(api, "/projects/{project_id}/environments/{environment_id}/services/{service_id}", serviceHandler.Get)
	huma.Delete(api, "/projects/{project_id}/environments/{environment_id}/services/{service_id}", serviceHandler.Delete)

	// Build routes
	huma.Get(api, "/builds", buildHandler.GetAllBuilds) // List all builds
	huma.Get(api, "/builds/{app_id}", buildHandler.Get) // Get Sepcific build info

	// Deployment routes
	huma.Get(api, "/deployments/{deployment_id}", applicationHandler.GetDeployment)

	huma.Get(api, "/applications", applicationHandler.GetAllApps) // List all applications
	huma.Get(api, "/applications/{app_id}", applicationHandler.Get)
	huma.Put(api, "/applications/{app_id}", applicationHandler.Update)
	huma.Delete(api, "/applications/{app_id}", applicationHandler.Delete)
	huma.Post(api, "/applications/{app_id}/deploy", applicationHandler.Deploy)
	huma.Post(api, "/applications/{app_id}/redeploy", applicationHandler.Deploy)
	huma.Post(api, "/applications/{app_id}/stop", applicationHandler.Start)
	huma.Post(api, "/applications/{app_id}/stop", applicationHandler.Stop)
	huma.Post(api, "/applications/{app_id}/restart", applicationHandler.Restart)
	huma.Get(api, "/applications/{app_id}/logs", applicationHandler.GetLogs)
	huma.Get(api, "/applications/{app_id}/environment-variables", applicationHandler.GetEnvVars)
	huma.Get(api, "/applications/{app_id}/deployments", deploymentHandler.ListDeploymentsForApp)
	// Quick creation (creates project + prod env + service/applications in one call)
	// Applications can be from git (public repo or private git repo using a private deploy key or a private github app)
	// Applications can be from Dockerfile
	huma.Post(api, "/applications/quick", serviceHandler.CreateQuickService)

	// Services routes
	huma.Post(api, "/services/quick", serviceHandler.CreateQuickService)
	huma.Get(api, "/services/{service_id}", serviceHandler.Get)
	huma.Get(api, "/services/{service_id}/logs", serviceHandler.GetLogs)
	huma.Delete(api, "/services/{service_id}", serviceHandler.Delete)
	huma.Get(api, "/services/{service_id}/backups", serviceHandler.GetBackups)

	// Security routes
	huma.get(api, "/security/private-keys", securityHandler.GetPrivateKeys)
	huma.get(api, "/security/private-keys/{key_id}", securityHandler.GetPrivateKey)
	huma.get(api, "/security/api-tokens", securityHandler.GetApiTokens)

	// Sources (git, etc...)
	huma.get(api, "/sources", sourceHandler.GetSources)
	huma.get(api, "/sources/{source_id}", sourceHandler.GetSource)

	return nil
}
