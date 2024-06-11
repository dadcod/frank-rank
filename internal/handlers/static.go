package handlers

import "net/http"

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", fs)
	http.ServeFile(w, r, "static/index.html")
}
