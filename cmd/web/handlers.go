package main

import (
	"net/http"
	"strconv"

	"github.com/amrojjeh/tajweed-calendar/ui"
)

func (app application) homeGet() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ui.HomePage(ui.NewHomeViewModel(app.events)).Render(r.Context(), w)
	})
}

func (app application) eventDetailsGet() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			app.serverError(w, err)
			return
		}
		idstrs := r.Form["id"]
		ms := []ui.EventDetailsViewModel{}
		for _, str := range idstrs {
			id, err := strconv.Atoi(str)
			if err != nil {
				app.clientError(w, http.StatusBadRequest)
				return
			}
			e, err := app.events.GetEventWithId(id)
			if err != nil {
				app.clientError(w, http.StatusBadRequest)
				return
			}
			ms = append(ms, ui.EventDetailsViewModel{
				Color: e.Color,
				Name:  e.Name,
			})
		}
		ui.EventDetails(ms...).Render(r.Context(), w)
	})
}
