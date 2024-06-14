package main

import (
	"database/sql"
	"fmt"
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
	envFile        = ".env"
	sessionManager *scs.SessionManager
)

func main() {
	fmt.Printf("Starting server on")
	env.LoadEnv(envFile)

	fmt.Printf("$v", os.Getenv("PORT"))

	oauthConfig.ClientID = os.Getenv("GOOGLE_CLIENT_ID")
	oauthConfig.ClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")

	sessionManager = scs.New()

	db, err := sql.Open("sqlite", "frankrank.db")
	if err != nil {
		log.Fatal(err)
	}

	queries := database.New(db)

	router := http.NewServeMux()

	as := handlers.NewAuthSession(oauthConfig, sessionManager, queries)

	router.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./fe/dist/assets"))))

	router.Handle("GET /", templ.Handler(templates.Layout(templates.Home())))
	router.HandleFunc("GET /login", as.HandleLogin)
	router.HandleFunc("GET /callback", as.HandleCallback)
	router.HandleFunc("GET /welcome", as.WelcomeHandler)

	autContext := middleware.NewAuthContext(sessionManager, []string{"/welcome"})

	stack := middleware.CreateStack(middleware.Logging, sessionManager.LoadAndSave, autContext.IsAuthenticated)

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("PORT")),
		Handler: stack(router),
	}
	server.ListenAndServe()
}
