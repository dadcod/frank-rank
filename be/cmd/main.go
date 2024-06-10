package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/alexedwards/scs/v2"
	"github.com/dadcod/frank-rank/internal/database"
	"github.com/dadcod/frank-rank/internal/handlers"
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

func main() {
	env.LoadEnv(envFile)

	oauthConfig.ClientID = os.Getenv("GOOGLE_CLIENT_ID")
	oauthConfig.ClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")

	sessionManager = scs.New()

	db, err := sql.Open("sqlite", "../frankrank.db")
	if err != nil {
		log.Fatal(err)
	}

	queries := database.New(db)

	router := http.NewServeMux()

	router.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	router.Handle("GET /home", templ.Handler(templates.Home()))

	as := handlers.NewAuthSession(oauthConfig, sessionManager, queries)

	router.HandleFunc("GET /login", as.HandleLogin)
	router.HandleFunc("GET /callback", as.HandleCallback)
	router.HandleFunc("GET /welcome", as.WelcomeHandler)
	router.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	autContext := middleware.NewAuthContext(sessionManager, []string{"/login", "/callback", "/home", "/static"})

	stack := middleware.CreateStack(middleware.Logging, sessionManager.LoadAndSave, autContext.IsAuthenticated)
	server := http.Server{
		Addr:    ":8080",
		Handler: stack(router),
	}
	server.ListenAndServe()
}
