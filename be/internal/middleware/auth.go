package middleware

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
)

type AuthUserIDKey string

const AuthUserID = "middleware.auth.userID"

type AuthContext struct {
	SessionManager *scs.SessionManager
	AllowedPaths   []string
}

func (a AuthContext) IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, path := range a.AllowedPaths {
			if r.URL.Path == path {
				next.ServeHTTP(w, r)
				return
			}
		}
		if a.SessionManager.Token(r.Context()) == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}
