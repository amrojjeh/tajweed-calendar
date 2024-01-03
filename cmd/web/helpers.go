package main

import (
	"log/slog"
	"net/http"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	http.Error(w, http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError)
	app.logger.Error("server error",
		slog.String("error", err.Error()))
}

func (app *application) clientError(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}
