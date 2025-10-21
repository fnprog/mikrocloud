package handlers

import (
	"github.com/go-chi/chi/v5"

	"github.com/mikrocloud/mikrocloud/internal/api/deps"
)

func RegisterProxyRoutes(r chi.Router, deps *deps.Dependencies) {
	handler := NewProxyHandler(deps.ProxyService)

	// Proxy routes within project
	r.Route("/proxy", func(r chi.Router) {
		r.Get("/", handler.ListProxyConfigs)
		r.Post("/", handler.CreateProxyConfig)
		r.Route("/{config_id}", func(r chi.Router) {
			r.Get("/", handler.GetProxyConfig)
			r.Put("/", handler.UpdateProxyConfig)
			r.Delete("/", handler.DeleteProxyConfig)
		})
	})
}
