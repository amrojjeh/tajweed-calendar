package main

import (
	"fmt"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/a-h/templ"
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
		sidebarModel := ui.SidebarViewModel{
			Hide:   false,
			Date:   fmt.Sprintf("%v %v ", time.Month(month), day),
			Events: []ui.EventDetailsViewModel{},
		}
		idstrs := r.Form["id"]
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
			m := ui.EventDetailsViewModel{
				Color:        e.Color,
				Name:         e.Name,
				Time:         e.Time(),
				Registration: templ.URL(e.RegistrationURL),
			}
			if e.Flyer != "" {
				m.Flyer = path.Join("/static/flyers", e.Flyer)
			}
			sidebarModel.Events = append(sidebarModel.Events, m)
		}
		ui.Sidebar(sidebarModel).Render(r.Context(), w)
	})
}
