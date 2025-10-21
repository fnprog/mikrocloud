package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/mikrocloud/mikrocloud/internal/api/deps"
	"github.com/mikrocloud/mikrocloud/internal/api/middleware"
)

func RegisterServersRoutes(r chi.Router, deps *deps.Dependencies) {
	handler := NewServersHandler(deps.ServerService)
	r.Route("/servers", func(r chi.Router) {
		r.Use(jwtauth.Authenticator(deps.JwtKeys))
		r.Use(middleware.ExtractUserOrg())

		r.Get("/", handler.ListServers)
		r.Post("/", handler.CreateServer)
		r.Route("/{server_id}", func(r chi.Router) {
			r.Get("/", handler.GetServer)
			r.Put("/", handler.UpdateServer)
			r.Delete("/", handler.DeleteServer)
		})
	})
}
