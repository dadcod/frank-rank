package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/alexedwards/scs/v2"
	"github.com/dadcod/frank-rank/pkg/env"
	"github.com/dadcod/frank-rank/internal/middleware"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	oauthConfig = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     google.Endpoint,
	}
	envFile        = "../../.env"
	sessionManager *scs.SessionManager
)

func main() {
	env.LoadEnv(envFile)
	oauthConfig.ClientID = os.Getenv("GOOGLE_CLIENT_ID")
	oauthConfig.ClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	sessionManager = scs.New()

	router := http.NewServeMux()
	stack := middleware.CreateStack(middleware.Logging, sessionManager.LoadAndSave)

	router.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})
	router.HandleFunc("GET /home", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!DOCTYPE html>
							<html lang="en">
								<head>
									<meta charset="UTF-8" />
									<title>Vite + TS</title>
								</head>
								<body>
									<a href="/login">Login with Google</a>
								</body>
							</html>
						`))
	})

	router.HandleFunc("GET /login", HandleLogin)
	router.HandleFunc("GET /callback", HandleCallback)
	router.HandleFunc("GET /welcome", WelcomeHandler)

	server := http.Server{
		Addr:    ":8080",
		Handler: stack(router),
	}
	server.ListenAndServe()
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	url := oauthConfig.AuthCodeURL("state")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
func HandleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	token, err := oauthConfig.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := oauthConfig.Client(context.Background(), token)
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
	sessionManager.Put(r.Context(), "userID", userInfo.Email)
	sessionManager.Put(r.Context(), "userName", userInfo.Name)
	http.Redirect(w, r, "/welcome", http.StatusFound)
}

func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	userID := sessionManager.GetString(r.Context(), "userID")
	userName := sessionManager.GetString(r.Context(), "userName")
	w.Write([]byte("Hello, " + userName + " (" + userID + ")"))
}
