package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/mikrocloud/mikrocloud/internal/api/deps"
	"github.com/mikrocloud/mikrocloud/internal/api/middleware"
)

func RegisterActivitiesRoutes(r chi.Router, deps *deps.Dependencies) {
	handlers := NewActivitiesHandlers(deps.ActivitiesService)

	r.Route("/activities", func(r chi.Router) {
		r.Use(jwtauth.Authenticator(deps.JwtKeys))
		r.Use(middleware.ExtractUserOrg())

		r.Get("/", handlers.GetRecentActivities)
		r.Get("/{resource_type}/{resource_id}", handlers.GetResourceActivities)
	})
}
