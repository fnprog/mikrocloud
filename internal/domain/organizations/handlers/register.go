package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/mikrocloud/mikrocloud/internal/api/deps"
	"github.com/mikrocloud/mikrocloud/internal/api/middleware"
)

func RegisterOrganizationsRoutes(r chi.Router, deps *deps.Dependencies) {
	handler := NewOrganizationHandler(deps.OrganizationService)

	// Organizations routes (protected)
	r.Route("/organizations", func(r chi.Router) {
		r.Use(middleware.AuthenticateAndExtract())

		r.Get("/", handler.ListOrganizations)
		r.Post("/", handler.CreateOrganization)
		r.Route("/{organization_id}", func(r chi.Router) {
			r.Get("/", handler.GetOrganization)
			r.Put("/", handler.UpdateOrganization)
			r.Delete("/", handler.DeleteOrganization)

			r.Route("/members", func(r chi.Router) {
				r.Get("/", handler.ListOrganizationMembers)
				r.Post("/", handler.InviteMember)
				r.Route("/{member_id}", func(r chi.Router) {
					r.Put("/", handler.UpdateMemberRole)
					r.Delete("/", handler.RemoveMember)
				})
			})

			r.Post("/transfer-ownership", handler.TransferOwnership)
		})
	})
}
