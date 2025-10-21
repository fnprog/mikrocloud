package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"

	"github.com/mikrocloud/mikrocloud/internal/api/deps"
	"github.com/mikrocloud/mikrocloud/internal/api/middleware"

	appHandler "github.com/mikrocloud/mikrocloud/internal/domain/applications/handlers"
	dbHandler "github.com/mikrocloud/mikrocloud/internal/domain/databases/handlers"
	disksHandler "github.com/mikrocloud/mikrocloud/internal/domain/disks/handlers"
	envHandler "github.com/mikrocloud/mikrocloud/internal/domain/environments/handlers"
	proxyHandler "github.com/mikrocloud/mikrocloud/internal/domain/proxy/handlers"
)

func RegisterProjectRoutes(r chi.Router, deps *deps.Dependencies) {
	handler := NewProjectHandler(deps.ProjectService)

	r.Route("/projects", func(r chi.Router) {
		r.Use(jwtauth.Authenticator(deps.JwtKeys))
		r.Use(middleware.ExtractUserOrg())

		r.Get("/", handler.List)
		r.Post("/", handler.Create)
		r.Route("/{project_id}", func(r chi.Router) {
			r.Get("/", handler.Get)
			r.Put("/", handler.Update)
			r.Delete("/", handler.Delete)

			envHandler.RegisterEnvironmentRoutes(r, deps)
			appHandler.RegisterApplicationRoutes(r, deps)
			dbHandler.RegisterDatabasesRoutes(r, deps)
			proxyHandler.RegisterProxyRoutes(r, deps)
			disksHandler.RegisterDisksRoutes(r, deps)
		})
	})
}
