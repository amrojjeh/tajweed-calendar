package main

import (
	"net/http"

	"github.com/amrojjeh/yaag/ui"
)

func (app application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/", app.homeGet())
	mux.Handle("/static/", http.FileServer(http.FS(ui.Files)))

	return mux
}
