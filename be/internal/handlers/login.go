package handlers

import "net/http"

func (as *AuthSession) HandleLogin(w http.ResponseWriter, r *http.Request) {
	url := as.oauthConfig.AuthCodeURL("state")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}