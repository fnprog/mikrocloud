package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/mikrocloud/mikrocloud/internal/api/deps"
)

func RegisterDeploymentsRoutes(r chi.Router, deps *deps.Dependencies) {
	deploymentHandler := NewDeploymentHandler(deps.DeploymentService, deps.ApplicationService)

	r.Route("/deployments", func(r chi.Router) {
		r.Get("/", deploymentHandler.ListDeployments)
		r.Post("/", deploymentHandler.CreateDeployment)
		r.Route("/{deployment_id}", func(r chi.Router) {
			r.Get("/", deploymentHandler.GetDeployment)
			r.Post("/stop", deploymentHandler.StopDeployment)
			r.Post("/cancel", deploymentHandler.CancelDeployment)
			r.Get("/logs", deploymentHandler.GetDeploymentLogs)
			r.Get("/logs/stream", deploymentHandler.StreamDeploymentLogs)
		})
	})
}
