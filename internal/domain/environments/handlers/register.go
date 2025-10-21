package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/mikrocloud/mikrocloud/internal/api/deps"
)

func RegisterEnvironmentRoutes(r chi.Router, deps *deps.Dependencies) {
	environmentHandler := NewEnvironmentHandler(deps.EnvironmentService)

	r.Route("/environments", func(r chi.Router) {
		r.Get("/", environmentHandler.ListEnvironments)
		r.Post("/", environmentHandler.CreateEnvironment)
		r.Route("/{environment_id}", func(r chi.Router) {
			r.Get("/", environmentHandler.GetEnvironment)
			r.Put("/", environmentHandler.UpdateEnvironment)
			r.Delete("/", environmentHandler.DeleteEnvironment)
		})
	})
}
