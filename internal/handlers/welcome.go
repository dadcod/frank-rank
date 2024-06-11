package handlers

import "net/http"

func (as *AuthSession) WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	userID := as.sessionManager.GetString(r.Context(), "userID")
	userName := as.sessionManager.GetString(r.Context(), "userName")
	w.Write([]byte("Hello, " + userName + " (" + userID + ")"))
}
