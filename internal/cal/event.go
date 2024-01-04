package cal

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type EventTime time.Time
type EventDuration time.Duration
type Events []Event

// TODO(Amr Ojjeh): Add "Added" for exceptional events within the series
type Event struct {
	Id        int           `json:"id"`
	Name      string        `json:"name"`
	Committee string        `json:"committee"`
	Next      NextFunc      `json:"next"`
	Color     string        `json:"color"`
	FirstDate EventTime     `json:"first"`
	Duration  EventDuration `json:"duration"`
	LastDate  EventTime     `json:"last"`
	Cancelled []EventTime   `json:"cancelled"`
}

func (es Events) EventsByCommittees() map[string]Events {
	cs := map[string]Events{}
	for _, e := range es {
		if e.Committee == "" {
			continue
		}
		_, found := cs[e.Committee]
		if !found {
			cs[e.Committee] = Events{}
		}
		cs[e.Committee] = append(cs[e.Committee], e)
	}
	return cs
}

func (es Events) GetEventWithId(id int) (Event, error) {
	for _, e := range es {
		if e.Id == id {
			return e, nil
		}
	}
	return Event{}, errors.New(
		fmt.Sprintf("cal: Events: there's no event with the id %v", id))
}

func (et *EventTime) UnmarshalJSON(data []byte) error {
	t, err := time.Parse("2006/01/02 03:04pm", string(data[1:len(data)-1]))
	if err != nil {
		err = errors.Join(errors.New("cal: EventTime: failed to parse JSON"), err)
		return err
	}

	*et = EventTime(t)
	return nil
}

func (ed *EventDuration) UnmarshalJSON(data []byte) error {
	var x float64
	err := json.Unmarshal(data, &x)
	if err != nil {
		err = errors.Join(errors.New("cal: EventDuration: failed to parse JSON"), err)
		return err
	}
	*ed = EventDuration(float64(time.Hour) * x)
	return nil
}

func NewEventDate(year int, month time.Month, day, hour, min int) EventTime {
	return EventTime(time.Date(year, month, day, hour, min, 0, 0, time.UTC))
}

func NewEventDateFromTime(t time.Time) EventTime {
	return NewEventDate(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute())
}

func NewEvent(name, color string, first time.Time, duration time.Duration,
	last time.Time, next NextFunc) Event {
	return Event{
		Name:      name,
		Next:      next,
		Color:     color,
		FirstDate: NewEventDateFromTime(first),
		LastDate:  NewEventDateFromTime(last),
		Cancelled: []EventTime{},
	}
}

func (e Event) Cancel(year int, month time.Month, day int) {
	e.Cancelled = append(e.Cancelled, NewEventDate(year, month, day,
		e.Hour(), e.Minute()))
}

func (e Event) IsOn(year int, month time.Month, day int) bool {
	t := time.Date(year, time.Month(month), day, e.Hour(),
		e.Minute(), 0, 0, time.UTC)
	if t.Before(time.Time(e.FirstDate)) ||
		(!time.Time(e.LastDate).IsZero() && t.After(time.Time(e.LastDate))) {
		return false
	}

	for _, c := range e.Cancelled {
		ct := time.Time(c)
		if ct.Year() == year && ct.Month() == month && ct.Day() == day {
			return false
		}
	}

	if e.Next == nil {
		return t.Equal(time.Time(e.FirstDate))
	}

	d := time.Time(e.FirstDate)
	for d.Before(t) {
		d = e.Next(d)
	}

	if d.Year() == year && d.Month() == month && d.Day() == day {
		return true
	}
	return false
}

func (e Event) Hour() int {
	return time.Time(e.FirstDate).Hour()
}

func (e Event) Minute() int {
	return time.Time(e.FirstDate).Minute()
}
