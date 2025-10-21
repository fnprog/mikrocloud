package handlers

import (
	"github.com/go-chi/chi/v5"

	"github.com/mikrocloud/mikrocloud/internal/api/deps"
)

func RegisterDisksRoutes(r chi.Router, deps *deps.Dependencies) {
	diskHandler := NewDiskHandler(deps.DiskService, deps.DatabaseService, deps.ContainerService)

	r.Route("/disks", func(r chi.Router) {
		r.Get("/", diskHandler.ListDisks)
		r.Post("/", diskHandler.CreateDisk)
		r.Route("/{disk_id}", func(r chi.Router) {
			r.Get("/", diskHandler.GetDisk)
			r.Put("/resize", diskHandler.ResizeDisk)
			r.Delete("/", diskHandler.DeleteDisk)
			r.Post("/attach", diskHandler.AttachDisk)
			r.Post("/detach", diskHandler.DetachDisk)
		})
	})
}
