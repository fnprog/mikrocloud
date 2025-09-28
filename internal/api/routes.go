package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/mikrocloud/mikrocloud/internal/api/handlers"
	"github.com/mikrocloud/mikrocloud/internal/config"
	"github.com/mikrocloud/mikrocloud/internal/container"
	"github.com/mikrocloud/mikrocloud/internal/database"
	athHanders "github.com/mikrocloud/mikrocloud/internal/domain/auth/handlers"
	pjHandlers "github.com/mikrocloud/mikrocloud/internal/domain/projects/handlers"
	"github.com/mikrocloud/mikrocloud/internal/middleware"
)

func SetupRoutes(api *chi.Mux, db *database.Database, cfg *config.Config, tokenAuth *jwtauth.JWTAuth) error {
	// Create service container
	serviceContainer, err := container.NewContainer(db, cfg)
	if err != nil {
		return err
	}

	// Create handler dependencies
	authHandler := athHanders.NewAuthHandler(db)
	projectHandler := pjHandlers.NewProjectHandler(db)
	environmentHandler := handlers.NewEnvironmentHandler(db)
	serviceHandler := handlers.NewServiceHandler(db, serviceContainer.ServiceService)

	// Maintenance
	// api.Get("/health", handlers.HealthCheck)
	// api.Get("/system/status", handlers.SystemStatus)
	// api.Get("/system/info", handlers.SystemInfo)
	// api.Get("/resources", handlers.GetResources)

	// Auth routes
	api.Post("/auth/login", authHandler.Login)
	api.Post("/auth/register", authHandler.Register)
	api.Post("/auth/logout", authHandler.Logout)

	// Everything down is authed
	api.Use(jwtauth.Authenticator(tokenAuth))
	api.Use(middleware.ExtractUserOrg())

	// Server
	// api.Get("/servers", serverHandler.List)
	// api.Post("/servers", serverHandler.Create)
	// api.Put("/servers/{server_id}", serverHandler.Update)
	// api.Delete("/servers/{server_id}", serverHandler.Delete)
	// api.Get("/servers/{server_id}", serverHandler.Get)
	// api.Get("/servers/{server_id}/domains", serverHandler.GetDomains)           // domain associated with the server
	// api.Get("/servers/{server_id}/destinations", serverHandler.GetDestinations) // Docker/podman, etc...
	// api.Get("/servers/{server_id}/proxy", serverHandler.GetDestinations)        // Server proxy settings
	// api.Get("/servers/{server_id}/resources", serverHandler.GetResources)       // hardware resource data for a server
	// api.Get("/servers/{server_id}/validate", serverHandler.Validate)            // validate server connectivity
	// api.Get("/servers/{server_id}/logs", serverHandler.Validate)                // server logs

	// Project routes
	api.Get("/projects", projectHandler.List)
	api.Post("/projects", projectHandler.Create)
	api.Get("/projects/{project_id}", projectHandler.Get)
	api.Put("/projects/{project_id}", projectHandler.Update)
	api.Delete("/projects/{project_id}", projectHandler.Delete)

	// Environment routes within projects
	api.Get("/projects/{project_id}/environments", environmentHandler.List)
	api.Post("/projects/{project_id}/environments", environmentHandler.Create)
	api.Get("/projects/{project_id}/environments/{environment_id}", environmentHandler.Get)
	api.Put("/projects/{project_id}/environments/{environment_id}", environmentHandler.Update)
	api.Delete("/projects/{project_id}/environments/{environment_id}", environmentHandler.Delete)

	// Applications routes within projects
	api.Get("/projects/{project_id}/environments/{environment_id}/applications", applicationHandler.List)
	api.Post("/projects/{project_id}/environments/{environment_id}/applications", applicationHandler.Create)

	// Service routes within project environments
	// api.Get("/projects/{project_id}/environments/{environment_id}/services", serviceHandler.List)
	// api.Post("/projects/{project_id}/environments/{environment_id}/services", serviceHandler.Create)
	// api.Get("/projects/{project_id}/environments/{environment_id}/services/{service_id}", serviceHandler.Get)
	// api.Delete("/projects/{project_id}/environments/{environment_id}/services/{service_id}", serviceHandler.Delete)

	// Build routes
	// api.Get("/builds", buildHandler.GetAllBuilds) // List all builds
	// api.Get("/builds/{app_id}", buildHandler.Get) // Get Sepcific build info

	// Deployment routes
	// api.Get("/deployments/{deployment_id}", applicationHandler.GetDeployment)

	api.Get("/applications", applicationHandler.GetAllApps) // List all applications
	api.Get("/applications/{app_id}", applicationHandler.Get)
	api.Put("/applications/{app_id}", applicationHandler.Update)
	api.Delete("/applications/{app_id}", applicationHandler.Delete)
	api.Post("/applications/{app_id}/deploy", applicationHandler.Deploy)
	api.Post("/applications/{app_id}/redeploy", applicationHandler.Deploy)
	api.Post("/applications/{app_id}/stop", applicationHandler.Start)
	api.Post("/applications/{app_id}/stop", applicationHandler.Stop)
	api.Post("/applications/{app_id}/restart", applicationHandler.Restart)
	api.Get("/applications/{app_id}/logs", applicationHandler.GetLogs)
	api.Get("/applications/{app_id}/environment-variables", applicationHandler.GetEnvVars)
	api.Get("/applications/{app_id}/deployments", deploymentHandler.ListDeploymentsForApp)
	// Quick creation (creates project + prod env + service/applications in one call)
	// Applications can be from git (public repo or private git repo using a private deploy key or a private github app)
	// Applications can be from Dockerfile
	// api.Post("/applications/quick", serviceHandler.CreateQuickService)

	// Services routes
	// api.Post("/services/quick", serviceHandler.CreateQuickService)
	// api.Get("/services/{service_id}", serviceHandler.Get)
	// api.Get("/services/{service_id}/logs", serviceHandler.GetLogs)
	// api.Delete("/services/{service_id}", serviceHandler.Delete)
	// api.Get("/services/{service_id}/backups", serviceHandler.GetBackups)

	// Security routes
	// api.get("/security/private-keys", securityHandler.GetPrivateKeys)
	// api.get("/security/private-keys/{key_id}", securityHandler.GetPrivateKey)
	// api.get("/security/api-tokens", securityHandler.GetApiTokens)

	// Sources (git, etc...)
	// api.get("/sources", sourceHandler.GetSources)
	// api.get("/sources/{source_id}", sourceHandler.GetSource)

	return nil
}
