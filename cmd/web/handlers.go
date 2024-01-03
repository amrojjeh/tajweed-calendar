package main

import (
	"log/slog"
	"net/http"

	"github.com/amrojjeh/yaag/ui"
)

func (app application) homeGet() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := ui.HomePage(ui.NewHomeViewModel(app.events)).Render(r.Context(), w)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			app.logger.Error("couldn't render HomePage",
				slog.String("error", err.Error()))
			return
		}
	})
}
