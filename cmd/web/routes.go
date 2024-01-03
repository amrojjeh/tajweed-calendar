package main

import "net/http"

func (app application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/", app.homeGet())
	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(
		http.Dir("./ui/static/"))))

	return mux
}
