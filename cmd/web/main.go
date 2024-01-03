package main

import (
	"encoding/json"
	"flag"
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

func loadEvents() ([]cal.Event, error) {
	data, err := data.Files.ReadFile("events.json")
	if err != nil {
		return []cal.Event{}, nil
	}
	events := []cal.Event{}
	err = json.Unmarshal(data, &events)
	if err != nil {
		return []cal.Event{}, err
	}

	return events, nil
}
