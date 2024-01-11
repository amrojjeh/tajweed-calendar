package main

import (
	"fmt"
	"net/http"
	"path"

	"github.com/a-h/templ"
	"github.com/amrojjeh/tajweed-calendar/ui"
)

func (app *application) homeGet() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ui.HomePage(ui.NewHomeViewModel(app.events)).Render(r.Context(), w)
	})
}

func (app *application) eventDetailsGet() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			app.serverError(w, err)
			return
		}

		q := query{
			app:  app,
			form: r.Form,
		}
		month, err := q.month()
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		day, err := q.day()
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		sidebarModel := ui.SidebarViewModel{
			Hide:   false,
			Date:   fmt.Sprintf("%v %v ", month, day),
			Events: []ui.EventDetailsViewModel{},
		}
		events, err := q.events()
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		for _, e := range events {
			m := ui.EventDetailsViewModel{
				Color:        e.Color,
				Name:         e.Name,
				Time:         e.Time(),
				Registration: templ.URL(e.RegistrationURL),
			}
			if e.Flyer != "" {
				m.Flyer = path.Join("/static/flyers", e.Flyer)
			}
			if info, ok := e.EventInfo(2024, month, day); ok {
				m.Flyer = path.Join("/static/flyers", info.Flyer)
				m.Time = info.Time()
			}
			sidebarModel.Events = append(sidebarModel.Events, m)
		}
		ui.Sidebar(sidebarModel).Render(r.Context(), w)
	})
}
