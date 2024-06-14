package handlers

import (
	"net/http"

	"github.com/dadcod/frank-rank/internal/templates"
)

func (as *AuthSession) WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	userID := as.sessionManager.GetString(r.Context(), "userID")
	userName := as.sessionManager.GetString(r.Context(), "userName")
	templates.Layout(templates.Welcome(userName, userID)).Render(r.Context(), w)
}
