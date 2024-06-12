package middleware

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
)

type AuthUserIDKey string

const AuthUserID = "middleware.auth.userID"

type AuthContext struct {
	SessionManager *scs.SessionManager
	ProtectedPaths []string
}

func NewAuthContext(sessionManager *scs.SessionManager, protectedPaths []string) *AuthContext {
	return &AuthContext{SessionManager: sessionManager, ProtectedPaths: protectedPaths}
}

func (a AuthContext) IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, path := range a.ProtectedPaths {
			if r.URL.Path == path && a.SessionManager.Token(r.Context()) == "" {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
