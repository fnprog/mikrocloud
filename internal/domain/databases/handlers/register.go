package handlers

import (
	"github.com/go-chi/chi/v5"

	"github.com/mikrocloud/mikrocloud/internal/api/deps"
)

func RegisterDatabasesRoutes(r chi.Router, deps *deps.Dependencies) {
	databaseHandler := NewDatabaseHandler(deps.DatabaseService, deps.ContainerService)

	// Database routes within project
	r.Route("/databases", func(r chi.Router) {
		r.Get("/", databaseHandler.ListDatabases)
		r.Post("/", databaseHandler.CreateDatabase)
		r.Get("/types", databaseHandler.GetDatabaseTypes)
		r.Get("/types/{type}/config", databaseHandler.GetDefaultDatabaseConfig)
		r.Route("/{database_id}", func(r chi.Router) {
			r.Get("/", databaseHandler.GetDatabase)
			r.Put("/", databaseHandler.UpdateDatabase)
			r.Delete("/", databaseHandler.DeleteDatabase)
			r.Post("/action", databaseHandler.DatabaseAction)
			r.Get("/logs", databaseHandler.GetDatabaseLogs)
			r.Get("/terminal", databaseHandler.HandleTerminal)

			RegisterDatabaseStudioRoutes(r, deps)
		})
	})
}

func RegisterDatabaseStudioRoutes(r chi.Router, deps *deps.Dependencies) {
	handler := NewStudioHandler(deps.DatabaseService)

	// Database studio routes
	r.Route("/studio", func(r chi.Router) {
		r.Get("/info", handler.GetDatabaseInfo)
		r.Get("/schemas", handler.ListSchemas)
		r.Get("/tables", handler.ListTables)
		r.Get("/tables/{table_name}/schema", handler.GetTableSchema)
		r.Get("/tables/{table_name}/data", handler.GetTableData)
		r.Post("/tables/{table_name}/data", handler.GetTableData)
		r.Post("/query", handler.ExecuteQuery)
		r.Post("/tables/{table_name}/rows", handler.InsertRow)
		r.Put("/tables/{table_name}/rows", handler.UpdateRow)
		r.Delete("/tables/{table_name}/rows", handler.DeleteRow)
	})
}
