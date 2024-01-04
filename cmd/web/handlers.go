package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

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
		ms := []ui.EventViewModel{}
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
			ms = append(ms, ui.EventViewModel{
				Color: e.Color,
				Name:  e.Name,
			})
		}
		monthstr := r.Form.Get("m")
		month, err := strconv.Atoi(monthstr)
		if err != nil {
			app.serverError(w, err)
			return
		}
		daystr := r.Form.Get("d")
		day, err := strconv.Atoi(daystr)
		if err != nil {
			app.serverError(w, err)
			return
		}
		m := ui.EventDetailsViewModel{
			Events: ms,
			Time:   fmt.Sprintf("%v %v ", time.Month(month), day),
		}
		ui.EventDetails(m).Render(r.Context(), w)
	})
}
