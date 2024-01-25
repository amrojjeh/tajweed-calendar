package main

import (
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/amrojjeh/tajweed-calendar/internal/cal"
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

type query struct {
	app  *application
	form url.Values
}

func (q query) events() ([]cal.Event, error) {
	ids := q.form["id"]
	events := make([]cal.Event, len(ids))
	for i, idstr := range ids {
		id, err := strconv.Atoi(idstr)
		if err != nil {
			return events, err
		}

		event, err := q.app.events.GetEventWithId(id)
		if err != nil {
			return events, err
		}

		events[i] = event
	}

	return events, nil
}

func (q query) month() (time.Month, error) {
	i, err := strconv.Atoi(q.form.Get("m"))
	if err != nil {
		return 0, nil
	}

	return time.Month(i), nil
}

func (q query) day() (int, error) {
	i, err := strconv.Atoi(q.form.Get("d"))
	if err != nil {
		return 0, nil
	}

	return i, nil
}
