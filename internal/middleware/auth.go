package middleware

import (
	"context"
	"net/http"
	"os/user"

	"github.com/alexedwards/scs/v2"
)

type AuthUserIDKey string

const AuthUserID AuthUserIDKey = "middleware.auth.userID"

type AuthContext struct {
	SessionManager *scs.SessionManager
	ProtectedPaths []string
}

func NewAuthContext(sessionManager *scs.SessionManager, protectedPaths []string) *AuthContext {
	return &AuthContext{SessionManager: sessionManager, ProtectedPaths: protectedPaths}
}

func (ac AuthContext) IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, path := range ac.ProtectedPaths {
			if r.URL.Path == path && ac.SessionManager.Token(r.Context()) == "" {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func (ac AuthContext) AddUserToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user := ac.SessionManager.GetString(r.Context(), "userID"); user != ""  {
		ctx := context.WithValue(r.Context(), AuthUserID, user)
		next.ServeHTTP(w, r.WithContext(ctx))}
	})
}
