package main

import (
	"net/http"

	"github.com/amrojjeh/tajweed-calendar/ui"
)

func (app application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/", app.homeGet())
	mux.Handle("/event-details", app.eventDetailsGet())
	mux.Handle("/static/", http.FileServer(http.FS(ui.Files)))

	return mux
}
