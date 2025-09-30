package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

// keys for context
type ctxKey string

const (
	CtxUserID ctxKey = "userID"
	CtxOrgID  ctxKey = "orgID"
)

func ExtractUserOrg() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())

			userID, _ := claims["user_id"].(string)
			orgID, _ := claims["org_id"].(string)

			if userID == "" {
				http.Error(w, "missing user in token", http.StatusUnauthorized)
				return
			}

			// For now, use user_id as org_id if no org_id is provided (temporary fix)
			if orgID == "" {
				orgID = userID
			}

			// put into context
			ctx := context.WithValue(r.Context(), CtxUserID, userID)
			ctx = context.WithValue(ctx, CtxOrgID, orgID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// helpers for handlers
func GetUserID(r *http.Request) string {
	v, _ := r.Context().Value(CtxUserID).(string)
	return v
}

func GetOrgID(r *http.Request) string {
	v, _ := r.Context().Value(CtxOrgID).(string)
	return v
}
