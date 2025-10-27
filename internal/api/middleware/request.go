package middleware

import (
	"net/http"
)

// keys for context
type ctxKey string

const (
	CtxUserID ctxKey = "userID"
	CtxOrgID  ctxKey = "orgID"
)

func CookieTokenInjector() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Authorization") == "" {
				if cookie, err := r.Cookie("auth_token"); err == nil {
					r.Header.Set("Authorization", "Bearer "+cookie.Value)
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

// GetUserID returns the user_id from the jwt token
func GetUserID(r *http.Request) string {
	v, _ := r.Context().Value(CtxUserID).(string)
	return v
}

// GetOrgID returns the user_id from the jwt token
func GetOrgID(r *http.Request) string {
	v, _ := r.Context().Value(CtxOrgID).(string)
	return v
}

func WebSocketTokenInjector() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.URL.Query().Get("token")
			if token != "" {
				r.Header.Set("Authorization", "Bearer "+token)
			}
			next.ServeHTTP(w, r)
		})
	}
}
