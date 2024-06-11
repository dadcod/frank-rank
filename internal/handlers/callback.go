package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/dadcod/frank-rank/internal/database"
	"golang.org/x/oauth2"
)

type AuthSession struct {
	oauthConfig    *oauth2.Config
	sessionManager *scs.SessionManager
	queries        *database.Queries
}

func NewAuthSession(oauthConfig *oauth2.Config, sessionManager *scs.SessionManager, queries *database.Queries) *AuthSession {
	return &AuthSession{oauthConfig: oauthConfig, sessionManager: sessionManager, queries: queries}
}

func (as *AuthSession) HandleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	token, err := as.oauthConfig.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := as.oauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userInfo struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	as.sessionManager.Put(r.Context(), "userID", userInfo.Email)
	as.sessionManager.Put(r.Context(), "userName", userInfo.Name)

	// Create a new user
	err = as.queries.CreateUser(r.Context(), database.CreateUserParams{
		Name:  as.sessionManager.GetString(r.Context(), "userName"),
		Email: as.sessionManager.GetString(r.Context(), "userID"),
	})
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/welcome", http.StatusFound)
}
