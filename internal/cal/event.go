package cal

import (
	"encoding/json"
	"errors"
	"time"
)

type EventTime time.Time
type EventDuration time.Duration

// TODO(Amr Ojjeh): Add "Added" for exceptional events within the series
type Event struct {
	Name      string        `json:"name"`
	Next      NextFunc      `json:"next"`
	Color     string        `json:"color"`
	FirstDate EventTime     `json:"first"`
	Duration  EventDuration `json:"duration"`
	LastDate  EventTime     `json:"last"`
	Cancelled []EventTime   `json:"cancelled"`
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
		e.NormalHour(), e.NormalMinute()))
}

func (e Event) IsOn(year int, month time.Month, day int) bool {
	t := time.Date(year, time.Month(month), day, e.NormalHour(),
		e.NormalMinute(), 0, 0, time.UTC)
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

func (e Event) NormalHour() int {
	return time.Time(e.FirstDate).Hour()
}

func (e Event) NormalMinute() int {
	return time.Time(e.FirstDate).Minute()
}
