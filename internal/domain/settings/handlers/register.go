package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/mikrocloud/mikrocloud/internal/api/deps"
)

func RegisterSettingsRoutes(r chi.Router, deps *deps.Dependencies) {
	handler := NewSettingsHandler(deps.SettingsService)

	r.Route("/settings", func(r chi.Router) {
		r.Get("/general", handler.GetGeneralSettings)

		r.Route("/", func(r chi.Router) {
			r.Use(jwtauth.Authenticator(deps.JwtKeys))
			r.Get("/advanced", handler.GetAdvancedSettings)
			r.Post("/general", handler.SaveGeneralSettings)
			r.Post("/advanced", handler.SaveAdvancedSettings)
			r.Get("/updates", handler.GetUpdateSettings)
			r.Post("/updates", handler.SaveUpdateSettings)
			r.Get("/instance", handler.GetInstanceInfo)
			r.Get("/detect-ips", handler.DetectIPAddresses)
			r.Post("/backup", handler.CreateBackup)
			r.Post("/restore", handler.RestoreBackup)
		})
	})
}
