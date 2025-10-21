package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/mikrocloud/mikrocloud/internal/utils"
)

func JWTCookieVerifier(ja *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var tokenString string

			if cookie, err := r.Cookie("auth_token"); err == nil && cookie.Value != "" {
				tokenString = cookie.Value
			} else {
				tokenString = jwtauth.TokenFromHeader(r)
			}

			if tokenString == "" {
				ctx := jwtauth.NewContext(r.Context(), nil, nil)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			token, err := ja.Decode(tokenString)
			if err != nil {
				ctx := jwtauth.NewContext(r.Context(), nil, err)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			if token == nil || jwt.Validate(token) != nil {
				ctx := jwtauth.NewContext(r.Context(), token, jwtauth.ErrExpired)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			ctx := jwtauth.NewContext(r.Context(), token, nil)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AuthenticateAndExtract() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, claims, err := jwtauth.FromContext(r.Context())

			if err != nil {
				utils.SendError(w, http.StatusUnauthorized, "unauthorized", "Invalid or missing token")
				return
			}

			if token == nil {
				utils.SendError(w, http.StatusUnauthorized, "unauthorized", "No token provided")
				return
			}

			if jwt.Validate(token) != nil {
				utils.SendError(w, http.StatusUnauthorized, "unauthorized", "Token is invalid or expired")
				return
			}

			userID, ok := claims["user_id"].(string)
			if !ok || userID == "" {
				utils.SendError(w, http.StatusUnauthorized, "unauthorized", "Missing user_id in token")
				return
			}

			orgID, _ := claims["org_id"].(string)

			if orgID == "" {
				orgID = userID
			}

			ctx := context.WithValue(r.Context(), CtxUserID, userID)
			ctx = context.WithValue(ctx, CtxOrgID, orgID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
