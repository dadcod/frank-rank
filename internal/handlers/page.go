package handlers

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/dadcod/frank-rank/internal/templates"
)

type PageHandlerFunc func(w http.ResponseWriter, r *http.Request)

func PageHandler(page templ.Component) PageHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Use-Partial") == "true" {
			page.Render(r.Context(), w)
		} else {
			templates.Layout(page).Render(r.Context(), w)
		}
	}
}
