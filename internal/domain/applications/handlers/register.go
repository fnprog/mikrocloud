package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/mikrocloud/mikrocloud/internal/api/deps"
	deploymentsHandler "github.com/mikrocloud/mikrocloud/internal/domain/deployments/handlers"
)

func RegisterApplicationRoutes(r chi.Router, deps *deps.Dependencies) {
	applicationHandler := NewApplicationHandler(deps.ApplicationService, deps.DeploymentService, deps.ContainerService)

	// Application routes within project
	r.Route("/applications", func(r chi.Router) {
		r.Get("/", applicationHandler.ListApplications)
		r.Post("/", applicationHandler.CreateApplication)
		r.Route("/{application_id}", func(r chi.Router) {
			r.Get("/", applicationHandler.GetApplication)
			r.Put("/", applicationHandler.UpdateApplication)
			r.Delete("/", applicationHandler.DeleteApplication)
			r.Post("/deploy", applicationHandler.DeployApplication)
			r.Post("/start", applicationHandler.StartApplication)
			r.Post("/stop", applicationHandler.StopApplication)
			r.Post("/restart", applicationHandler.RestartApplication)
			r.Get("/logs", applicationHandler.GetApplicationLogs)
			r.Patch("/general", applicationHandler.UpdateGeneral)
			r.Post("/domain/generate", applicationHandler.GenerateDomain)
			r.Put("/domain", applicationHandler.AssignDomain)
			r.Put("/ports", applicationHandler.UpdatePorts)

			// Deployment routes within application
			deploymentsHandler.RegisterDeploymentsRoutes(r, deps)
		})
	})
}
