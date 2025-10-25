package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mikrocloud/mikrocloud/internal/api/deps"
	"github.com/mikrocloud/mikrocloud/internal/api/middleware"
	activitiesHandler "github.com/mikrocloud/mikrocloud/internal/domain/activities/handlers"
	authHandler "github.com/mikrocloud/mikrocloud/internal/domain/auth/handlers"
	deploymentsHandler "github.com/mikrocloud/mikrocloud/internal/domain/deployments/handlers"
	gitHandler "github.com/mikrocloud/mikrocloud/internal/domain/git/handlers"
	maintenanceHandler "github.com/mikrocloud/mikrocloud/internal/domain/maintenance/handlers"
	organizationsHandler "github.com/mikrocloud/mikrocloud/internal/domain/organizations/handlers"
	projectsHandler "github.com/mikrocloud/mikrocloud/internal/domain/projects/handlers"
	serversHandler "github.com/mikrocloud/mikrocloud/internal/domain/servers/handlers"
	templatesHandler "github.com/mikrocloud/mikrocloud/internal/domain/services/handlers"
	settingsHandler "github.com/mikrocloud/mikrocloud/internal/domain/settings/handlers"
	tunnelsHandler "github.com/mikrocloud/mikrocloud/internal/domain/tunnels/handlers"
)

func SetupRoutes(api chi.Router, dependencies *deps.Dependencies) {
	api.Use(middleware.JWTCookieVerifier(dependencies.JwtKeys))

	activitiesHandler.RegisterActivitiesRoutes(api, dependencies)
	authHandler.RegisterAuthRoutes(api, dependencies)
	deploymentsHandler.RegisterDeploymentsRoutes(api, dependencies)
	gitHandler.RegisterGitRoutes(api, dependencies)
	maintenanceHandler.RegisterMaintenanceRoutes(api, dependencies)
	organizationsHandler.RegisterOrganizationsRoutes(api, dependencies)
	projectsHandler.RegisterProjectRoutes(api, dependencies)
	serversHandler.RegisterServersRoutes(api, dependencies)
	settingsHandler.RegisterSettingsRoutes(api, dependencies)
	templatesHandler.RegisterTemplatesRoutes(api, dependencies)
	tunnelsHandler.RegisterTunnelRoutes(api, dependencies)

	// Serve storage files (public access)
	storageDir := "./storage"
	api.Handle("/storage/*", http.StripPrefix("/storage/", http.FileServer(http.Dir(storageDir))))

	// Protected routes that require authentication
	// 	r.Use(middleware.WebSocketTokenInjector())
}
