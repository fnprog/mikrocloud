package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/mikrocloud/mikrocloud/internal/api/deps"
	"github.com/mikrocloud/mikrocloud/internal/api/middleware"
	gitHandlers "github.com/mikrocloud/mikrocloud/internal/domain/git/handlers"
)

func RegisterOAuthRoutes(r chi.Router, deps *deps.Dependencies) {
	handler := gitHandlers.NewOAuthHandlers(deps.GitService, deps.Config)

	r.Group(func(r chi.Router) {
		r.Use(middleware.CookieTokenInjector())
		r.Use(jwtauth.Authenticator(deps.JwtKeys))
		r.Use(middleware.ExtractUserOrg())
		r.Get("/git/oauth/start", handler.StartOAuth)
	})

	// Git OAuth callback (public route - OAuth provider redirects here)
	r.Get("/git/oauth/callback", handler.Callback)
}
