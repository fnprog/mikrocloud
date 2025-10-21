package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/mikrocloud/mikrocloud/internal/api/deps"
)

func RegisterMaintenanceRoutes(r chi.Router, deps *deps.Dependencies) {
	maintenanceHandler := NewMaintenanceHandler(deps)

	// Maintenance routes (protected)
	r.Route("/maintenance", func(r chi.Router) {
		r.Use(jwtauth.Authenticator(deps.JwtKeys))

		r.Get("/health", maintenanceHandler.HealthCheck)
		r.Get("/status", maintenanceHandler.SystemStatus)
		r.Get("/resources", maintenanceHandler.GetResources)
		r.Get("/info", maintenanceHandler.SystemInfo)

		r.Route("/domains", func(r chi.Router) {
			r.Get("/", maintenanceHandler.ListDomains)
			r.Post("/", maintenanceHandler.AddDomain)
			r.Delete("/{domain_id}", maintenanceHandler.RemoveDomain)
			r.Post("/{domain_id}/ssl", maintenanceHandler.EnableSSL)
		})
	})
}
