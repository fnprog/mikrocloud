package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/mikrocloud/mikrocloud/internal/api/deps"
	"github.com/mikrocloud/mikrocloud/internal/api/middleware"
)

func RegisterGitRoutes(r chi.Router, deps *deps.Dependencies) {
	gitHandler := NewGitHandler(deps.GitService)

	gitOAuthHandler := NewOAuthHandlers(deps.GitService, deps.Config)
	gitHubAppHandler := NewGitHubAppHandlers(deps.GitService, deps.Config, gitOAuthHandler.GetStateStore())
	gitWebhookHandler := NewWebhookHandlers(deps.GitService)

	// Git routes
	r.Route("/git", func(r chi.Router) {
		r.Post("/validate", gitHandler.ValidateRepository)
		r.Post("/branches", gitHandler.ListBranches)
		r.Post("/detect-build", gitHandler.DetectBuildMethod)
		r.Post("/sources", gitHandler.CreateGitSource)
		r.Get("/sources", gitHandler.ListGitSources)
		r.Route("/sources/{source_id}", func(r chi.Router) {
			r.Get("/", gitHandler.GetGitSource)
			r.Put("/", gitHandler.UpdateGitSource)
			r.Delete("/", gitHandler.DeleteGitSource)
		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.CookieTokenInjector())
			r.Use(jwtauth.Authenticator(deps.JwtKeys))
			r.Use(middleware.ExtractUserOrg())
			r.Get("/github-app/manifest", gitHubAppHandler.GenerateManifest)
		})

		// GitHub App callbacks (public routes - GitHub redirects here)
		r.Get("/github-app/callback", gitHubAppHandler.Callback)
		r.Get("/github-app/install", gitHubAppHandler.InstallCallback)
	})

	// Public webhook routes (no authentication required)
	RegisterWebhookRoutes(r, gitWebhookHandler)
}
