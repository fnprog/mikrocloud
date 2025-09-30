package api

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/mikrocloud/mikrocloud/internal/api/middleware"
	"github.com/mikrocloud/mikrocloud/internal/config"
	"github.com/mikrocloud/mikrocloud/internal/database"
	appHandlers "github.com/mikrocloud/mikrocloud/internal/domain/applications/handlers"
	applicationsService "github.com/mikrocloud/mikrocloud/internal/domain/applications/service"
	authHandlers "github.com/mikrocloud/mikrocloud/internal/domain/auth/handlers"
	authService "github.com/mikrocloud/mikrocloud/internal/domain/auth/service"
	databaseHandlers "github.com/mikrocloud/mikrocloud/internal/domain/databases/handlers"
	databaseService "github.com/mikrocloud/mikrocloud/internal/domain/databases/service"
	deploymentHandlers "github.com/mikrocloud/mikrocloud/internal/domain/deployments/handlers"
	deploymentService "github.com/mikrocloud/mikrocloud/internal/domain/deployments/service"
	envHandlers "github.com/mikrocloud/mikrocloud/internal/domain/environments/handlers"
	environmentService "github.com/mikrocloud/mikrocloud/internal/domain/environments/service"
	projectHandlers "github.com/mikrocloud/mikrocloud/internal/domain/projects/handlers"
	projectService "github.com/mikrocloud/mikrocloud/internal/domain/projects/service"
	serviceHandlers "github.com/mikrocloud/mikrocloud/internal/domain/services/handlers"
	"github.com/mikrocloud/mikrocloud/internal/domain/services/repository"
	servicesService "github.com/mikrocloud/mikrocloud/internal/domain/services/service"
	buildService "github.com/mikrocloud/mikrocloud/pkg/containers/build"
	databaseContainers "github.com/mikrocloud/mikrocloud/pkg/containers/database"
	"github.com/mikrocloud/mikrocloud/pkg/containers/manager"
)

func SetupRoutes(api chi.Router, db *database.Database, cfg *config.Config, tokenAuth *jwtauth.JWTAuth) error {
	// Create container manager
	containerManager, err := createContainerManager(cfg)
	if err != nil {
		return fmt.Errorf("failed to create container manager: %w", err)
	}

	// Create service instances
	projSvc := projectService.NewProjectService(db.ProjectRepository)
	authSvc := authService.NewAuthService(db.SessionRepository, db.AuthRepository, db.UserRepository, cfg.Auth.JWTSecret)
	envSvc := environmentService.NewEnvironmentService(db.EnvironmentRepository)
	appSvc := applicationsService.NewApplicationService(db.ApplicationRepository)

	// Create database container deployment service
	dbImageResolver := databaseContainers.NewDefaultImageResolver()
	dbConfigBuilder := databaseContainers.NewDefaultContainerConfigBuilder(dbImageResolver)
	dbDeploymentSvc := databaseContainers.NewDatabaseDeploymentService(containerManager, dbImageResolver, dbConfigBuilder)
	dbSvc := databaseService.NewDatabaseService(db.DatabaseRepository, dbDeploymentSvc)

	// Create and start database status sync service
	logger := slog.Default()
	statusSyncSvc := databaseService.NewStatusSyncService(dbSvc, containerManager, logger, 30*time.Second)
	go statusSyncSvc.Start(context.Background())

	// Create QuickDeployService wrapper for ApplicationService
	quickDeployService := repository.NewQuickDeployService(db.TemplateRepository, appSvc)
	templateSvc := servicesService.NewTemplateService(db.TemplateRepository, quickDeployService)

	// Create build service
	buildSvc := buildService.NewBuildService(containerManager, cfg.Docker.SocketPath)

	// Create deployment service
	deploymentSvc := deploymentService.NewDeploymentService(
		db.DeploymentRepository,
		buildSvc,
	)

	// Create handler dependencies
	authHandler := authHandlers.NewAuthHandler(authSvc)
	projectHandler := projectHandlers.NewProjectHandler(projSvc)
	environmentHandler := envHandlers.NewEnvironmentHandler(envSvc)
	applicationHandler := appHandlers.NewApplicationHandler(appSvc)
	databaseHandler := databaseHandlers.NewDatabaseHandler(dbSvc, containerManager)
	deploymentHandler := deploymentHandlers.NewDeploymentHandler(deploymentSvc, appSvc)
	templateHandler := serviceHandlers.NewTemplateHandler(templateSvc)

	// Protected routes that require authentication
	api.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))
		r.Use(middleware.ExtractUserOrg())

		// Project routes
		r.Route("/projects", func(r chi.Router) {
			r.Get("/", projectHandler.List)
			r.Post("/", projectHandler.Create)
			r.Route("/{project_id}", func(r chi.Router) {
				r.Get("/", projectHandler.Get)
				r.Put("/", projectHandler.Update)
				r.Delete("/", projectHandler.Delete)

				// Environment routes within project
				r.Route("/environments", func(r chi.Router) {
					r.Get("/", environmentHandler.ListEnvironments)
					r.Post("/", environmentHandler.CreateEnvironment)
					r.Route("/{environment_id}", func(r chi.Router) {
						r.Get("/", environmentHandler.GetEnvironment)
						r.Put("/", environmentHandler.UpdateEnvironment)
						r.Delete("/", environmentHandler.DeleteEnvironment)
					})
				})

				// Application routes within project
				r.Route("/applications", func(r chi.Router) {
					r.Get("/", applicationHandler.ListApplications)
					r.Post("/", applicationHandler.CreateApplication)
					r.Route("/{application_id}", func(r chi.Router) {
						r.Get("/", applicationHandler.GetApplication)
						r.Put("/", applicationHandler.UpdateApplication)
						r.Delete("/", applicationHandler.DeleteApplication)
						r.Post("/deploy", applicationHandler.DeployApplication)

						// Deployment routes within application
						r.Route("/deployments", func(r chi.Router) {
							r.Get("/", deploymentHandler.ListDeployments)
							r.Post("/", deploymentHandler.CreateDeployment)
							r.Route("/{deployment_id}", func(r chi.Router) {
								r.Get("/", deploymentHandler.GetDeployment)
								r.Post("/stop", deploymentHandler.StopDeployment)
								r.Post("/cancel", deploymentHandler.CancelDeployment)
								r.Get("/logs", deploymentHandler.GetDeploymentLogs)
							})
						})
					})
				})

				// Database routes within project
				r.Route("/databases", func(r chi.Router) {
					r.Get("/", databaseHandler.ListDatabases)
					r.Post("/", databaseHandler.CreateDatabase)
					r.Get("/types", databaseHandler.GetDatabaseTypes)
					r.Get("/types/{type}/config", databaseHandler.GetDefaultDatabaseConfig)
					r.Route("/{database_id}", func(r chi.Router) {
						r.Get("/", databaseHandler.GetDatabase)
						r.Put("/", databaseHandler.UpdateDatabase)
						r.Delete("/", databaseHandler.DeleteDatabase)
						r.Post("/action", databaseHandler.DatabaseAction)
						r.Get("/logs", databaseHandler.GetDatabaseLogs)
					})
				})
			})
		})

		// Template routes
		r.Route("/templates", func(r chi.Router) {
			r.Get("/", templateHandler.ListTemplates)
			r.Post("/", templateHandler.CreateTemplate)
			r.Get("/official", templateHandler.ListOfficialTemplates)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", templateHandler.GetTemplate)
				r.Put("/", templateHandler.UpdateTemplate)
				r.Delete("/", templateHandler.DeleteTemplate)
				r.Post("/deploy", templateHandler.DeployTemplate)
				r.Post("/preview", templateHandler.PreviewDeployment)
			})
		})
	})

	// Public routes (no authentication required)
	api.Route("/auth", func(r chi.Router) {
		r.Get("/setup", authHandler.GetSetupStatus)
		r.Post("/login", authHandler.Login)
		r.Post("/register", authHandler.Register)
		r.Post("/refresh", authHandler.RefreshToken)

		// Protected auth routes
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator(tokenAuth))
			r.Post("/logout", authHandler.Logout)
			r.Get("/profile", authHandler.GetProfile)
		})
	})

	return nil
}

func createContainerManager(cfg *config.Config) (manager.ContainerManager, error) {
	switch cfg.Docker.Runtime {
	case "docker":
		return manager.NewDockerManager()
	case "podman":
		return manager.NewPodmanManager()
	default:
		return manager.NewDockerManager() // Default to Docker
	}
}
