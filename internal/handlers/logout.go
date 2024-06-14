package handlers

import "net/http"

func (as *AuthSession) HandleLogout(w http.ResponseWriter, r *http.Request) {
	err := as.sessionManager.Destroy(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
