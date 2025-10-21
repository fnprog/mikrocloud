package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/mikrocloud/mikrocloud/internal/api/deps"
)

func RegisterTemplatesRoutes(r chi.Router, deps *deps.Dependencies) {
	handler := NewTemplateHandler(deps.TemplateService)

	r.Route("/templates", func(r chi.Router) {
		r.Get("/", handler.ListTemplates)
		r.Post("/", handler.CreateTemplate)
		r.Get("/official", handler.ListOfficialTemplates)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handler.GetTemplate)
			r.Put("/", handler.UpdateTemplate)
			r.Delete("/", handler.DeleteTemplate)
			r.Post("/deploy", handler.DeployTemplate)
			r.Post("/preview", handler.PreviewDeployment)
		})
	})
}
