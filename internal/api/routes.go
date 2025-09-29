package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	environmentHandlers "github.com/mikrocloud/mikrocloud/internal/api/handlers"
	"github.com/mikrocloud/mikrocloud/internal/config"
	"github.com/mikrocloud/mikrocloud/internal/database"
	authHandlers "github.com/mikrocloud/mikrocloud/internal/domain/auth/handlers"
	authService "github.com/mikrocloud/mikrocloud/internal/domain/auth/service"
	environmentService "github.com/mikrocloud/mikrocloud/internal/domain/environments/service"
	"github.com/mikrocloud/mikrocloud/internal/domain/projects/handlers"
	projectService "github.com/mikrocloud/mikrocloud/internal/domain/projects/service"
	"github.com/mikrocloud/mikrocloud/internal/middleware"
)

func SetupRoutes(api chi.Router, db *database.Database, cfg *config.Config, tokenAuth *jwtauth.JWTAuth) error {
	// Create service instances
	projSvc := projectService.NewProjectService(db.ProjectRepository)
	authSvc := authService.NewAuthService(db.SessionRepository, db.AuthRepository, cfg.Auth.JWTSecret)
	envSvc := environmentService.NewEnvironmentService(db.EnvironmentRepository)

	// Create handler dependencies
	authHandler := authHandlers.NewAuthHandler(authSvc)
	projectHandler := handlers.NewProjectHandler(projSvc)
	environmentHandler := environmentHandlers.NewEnvironmentHandler(envSvc)

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
			})
		})

		// Application routes - TODO: Implement application handlers
		/*
			r.Route("/applications", func(r chi.Router) {
				r.Get("/", applicationHandler.List)
				r.Post("/", applicationHandler.Create)
			})
		*/

		// Service routes - TODO: Implement service handlers
		/*
			r.Route("/services", func(r chi.Router) {
				r.Get("/", serviceHandler.List)
				r.Post("/", serviceHandler.Create)
				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", serviceHandler.Get)
					r.Put("/", serviceHandler.Update)
					r.Delete("/", serviceHandler.Delete)
					r.Post("/start", serviceHandler.Start)
					r.Post("/stop", serviceHandler.Stop)
					r.Post("/restart", serviceHandler.Restart)
					r.Get("/logs", serviceHandler.Logs)
					r.Get("/stats", serviceHandler.Stats)
					r.Post("/scale", serviceHandler.Scale)
				})
			})
		*/

		// Deployment routes - TODO: Implement deployment handlers
		/*
			r.Route("/deployments", func(r chi.Router) {
				r.Get("/", deploymentHandler.List)
				r.Post("/", deploymentHandler.Create)
				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", deploymentHandler.Get)
					r.Put("/", deploymentHandler.Update)
					r.Delete("/", deploymentHandler.Delete)
					r.Post("/start", deploymentHandler.Start)
					r.Post("/stop", deploymentHandler.Stop)
					r.Post("/restart", deploymentHandler.Restart)
					r.Get("/logs", deploymentHandler.Logs)
				})
			})
		*/
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
