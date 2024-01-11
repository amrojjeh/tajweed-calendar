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

var id_counter = 0

type Event struct {
	Id              int           `json:"-"`
	Name            string        `json:"name"`
	Committee       string        `json:"committee"`
	Recurring       RecurringType `json:"next"`
	Color           string        `json:"color"`
	RegistrationURL string        `json:"registration"`
	FirstDate       EventTime     `json:"first"`
	Duration        EventDuration `json:"duration"`
	LastDate        EventTime     `json:"last"`
	Cancelled       []EventTime   `json:"cancelled"`
	Details         []EventInfo   `json:"details"`
}

type EventInfo struct {
	Date     EventTime     `json:"date"`
	Flyer    string        `json:"flyer"`
	Duration EventDuration `json:"duration"`
}

func (es Events) EventsByCommittee() map[string]Events {
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

func (ed EventDuration) Hours() float64 {
	return time.Duration(ed).Hours()
}

func NewEventDate(year int, month time.Month, day, hour, min int) EventTime {
	return EventTime(time.Date(year, month, day, hour, min, 0, 0, time.UTC))
}

func NewEventDateFromTime(t time.Time) EventTime {
	return NewEventDate(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute())
}

func NewEvent(name, color string, first time.Time, duration time.Duration,
	last time.Time, next string) Event {
	return Event{
		Name:      name,
		Recurring: RecurringTypeFromString(next),
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

func (e Event) EventInfo(year int, month time.Month, day int) (EventInfo, bool) {
	for _, d := range e.Details {
		dTime := time.Time(d.Date)
		if dTime.Year() == year && dTime.Month() == month && dTime.Day() == day {
			return d, true
		}
	}

	return EventInfo{}, false
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

	if e.Recurring.IsOn != nil {
		return e.Recurring.IsOn(time.Time(e.FirstDate), t)
	}

	return t.Equal(time.Time(e.FirstDate))
}

func (e Event) Hour() int {
	return time.Time(e.FirstDate).Hour()
}

func (e Event) Minute() int {
	return time.Time(e.FirstDate).Minute()
}

func (e Event) Time() string {
	if e.Duration.Hours() == 24.0 {
		return "All day"
	}

	begins := time.Time(e.FirstDate)
	ends := begins.Add(time.Duration(e.Duration))
	return fmt.Sprintf("%v - %v", begins.Format("03:04"), ends.Format("03:04pm"))
}

func (e EventInfo) Time() string {
	begins := time.Time(e.Date)
	ends := begins.Add(time.Duration(e.Duration))
	return fmt.Sprintf("%v - %v", begins.Format("03:04"), ends.Format("03:04pm"))
}
