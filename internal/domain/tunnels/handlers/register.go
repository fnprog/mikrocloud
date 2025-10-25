package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/mikrocloud/mikrocloud/internal/api/deps"
	"github.com/mikrocloud/mikrocloud/internal/api/middleware"
)

func RegisterTunnelRoutes(r chi.Router, deps *deps.Dependencies) {
	tunnelHandler := NewTunnelHandler(deps.TunnelService)

	r.Route("/tunnels", func(r chi.Router) {
		r.Use(middleware.AuthenticateAndExtract())
		r.Get("/", tunnelHandler.ListTunnels)
		r.Post("/", tunnelHandler.CreateTunnel)
		r.Route("/{tunnel_id}", func(r chi.Router) {
			r.Get("/", tunnelHandler.GetTunnel)
			r.Delete("/", tunnelHandler.DeleteTunnel)
			r.Post("/start", tunnelHandler.StartTunnel)
			r.Post("/stop", tunnelHandler.StopTunnel)
			r.Post("/restart", tunnelHandler.RestartTunnel)
		})
	})
}
