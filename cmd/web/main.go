package main

import (
	"encoding/json"
	"flag"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/amrojjeh/tajweed-calendar/data"
	"github.com/amrojjeh/tajweed-calendar/internal/cal"
)

type application struct {
	logger *slog.Logger
	events cal.Events
}

func main() {
	addr := flag.String("addr", ":8080", "Set the HTTP address")
	flag.Parse()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))
	events, err := loadEvents()
	if err != nil {
		logger.Error("could not load events", slog.String("error", err.Error()))
		os.Exit(1)
	}
	logger.Info("loaded events")
	app := application{
		logger: logger,
		events: events,
	}
	server := http.Server{
		Addr:         *addr,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
		Handler:      app.routes(),
	}

	logger.Info("starting server", slog.String("addr", *addr))
	if err := server.ListenAndServe(); err != nil {
		logger.Error("server failed", slog.String("err", err.Error()))
		os.Exit(1)
	}
}

func loadEvents() (cal.Events, error) {
	files, err := fs.Glob(data.Files, "*.json")
	if err != nil {
		return cal.Events{}, nil
	}
	events := []cal.Event{}
	for _, f := range files {
		es := []cal.Event{}
		data, err := data.Files.ReadFile(f)
		if err != nil {
			return cal.Events{}, nil
		}
		err = json.Unmarshal(data, &es)
		if err != nil {
			return cal.Events{}, err
		}

		for _, e := range es {
			for _, d := range e.Details {
				if time.Duration(d.Duration) == 0 {
					d.Duration = e.Duration
				}

				if d.Flyer == "" {
					d.Flyer = e.Flyer
				}
			}
			events = append(events, e)
		}
	}

	for i := 0; i < len(events); i++ {
		events[i].Id = i
	}

	return events, nil
}
