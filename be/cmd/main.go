package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/alexedwards/scs/v2"
	"github.com/dadcod/frank-rank/internal/database"
	"github.com/dadcod/frank-rank/internal/middleware"
	"github.com/dadcod/frank-rank/internal/templates"
	"github.com/dadcod/frank-rank/pkg/env"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	_ "modernc.org/sqlite"
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

func createRouter(config *http.Server, middleware middleware.Middleware) *http.Server {
	router := http.NewServeMux()

	config.Handler = middleware(router)
	return config
}

func main() {
	env.LoadEnv(envFile)
	oauthConfig.ClientID = os.Getenv("GOOGLE_CLIENT_ID")
	oauthConfig.ClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	sessionManager = scs.New()
	autContext := middleware.AuthContext{SessionManager: sessionManager, AllowedPaths: []string{"/login", "/callback", "/home"}}
	router := http.NewServeMux()
	stack := middleware.CreateStack(middleware.Logging, sessionManager.LoadAndSave, autContext.IsAuthenticated)

	router.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})
	router.Handle("GET /home", templ.Handler(templates.Home()))

	router.HandleFunc("GET /login", HandleLogin)
	router.HandleFunc("GET /callback", HandleCallback)
	router.HandleFunc("GET /welcome", WelcomeHandler)
	router.HandleFunc("GET /test", TestDbHandler)

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

var ddl string = "../frankrank.db"

func TestDbHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite", ddl)
	if err != nil {
		log.Fatal(err)
	}

	queries := database.New(db)

	ctx := context.Background()

	if name := sessionManager.GetString(r.Context(), "userName"); name != "" {

		// Create a new user
		err = queries.CreateUser(ctx, database.CreateUserParams{
			Name:  sessionManager.GetString(r.Context(), "userName"),
			Email: sessionManager.GetString(r.Context(), "userID"),
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	// Get the user
	user, err := queries.GetUserByEmail(ctx, sessionManager.GetString(r.Context(), "userID"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("User: %v\n", user)
}
