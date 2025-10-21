package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/mikrocloud/mikrocloud/internal/api/deps"
	oauthHandler "github.com/mikrocloud/mikrocloud/internal/domain/oauth/handlers"
)

func RegisterAuthRoutes(r chi.Router, deps *deps.Dependencies) {
	handler := NewAuthHandler(deps.AuthService)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", handler.Login)
		r.Post("/register", handler.Register)
		r.Post("/refresh", handler.RefreshToken)
		r.Get("/setup", handler.GetSetupStatus)

		// Oauth Routes
		oauthHandler.RegisterOAuthRoutes(r, deps)

		// Protected auth routes
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Authenticator(deps.JwtKeys))
			r.Post("/logout", handler.Logout)
			r.Get("/profile", handler.GetProfile)
			r.Put("/profile", handler.UpdateProfile)
			r.Post("/avatar", handler.UploadAvatar)
			r.Delete("/profile", handler.DeleteAccount)
			r.Put("/email", handler.UpdateEmail)
			r.Put("/password", handler.UpdatePassword)
		})
	})
}
