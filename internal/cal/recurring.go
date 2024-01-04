package cal

import (
	"encoding/json"
	"errors"
	"time"
)

// TODO(Amr Ojjeh): Turn NextFunc into an interface such that it
// supports both a "next" method as well as an "IsOn" method

type RecurringType struct {
	Next NextFunc
	IsOn IsOnFunc
}

type NextFunc func(time.Time) time.Time
type IsOnFunc func(time.Time, time.Time) bool

func (n *RecurringType) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		err = errors.Join(errors.New("cal: NextFunc: failed to parse JSON"), err)
		return err
	}
	*n = RecurringTypeFromString(s)
	return nil
}

func RecurringTypeFromString(s string) RecurringType {
	switch s {
	case "":
		return RecurringType{}
	case "weekly":
		return RecurringType{
			Next: weekly,
			IsOn: weeklyIsOn,
		}
	case "monthly by weekday count":
		return RecurringType{
			Next: monthlyByWeekdayCount,
			IsOn: monthlyByWeekdayCountIsOn,
		}
	case "everyday":
		return RecurringType{
			Next: everyday,
			IsOn: everydayIsOn,
		}
	default:
		return RecurringType{}
	}
}

func weekly(t time.Time) time.Time {
	return t.AddDate(0, 0, 7)
}

func weeklyIsOn(first time.Time, ison time.Time) bool {
	return first.Weekday() == ison.Weekday()
}

func monthlyByWeekdayCount(t time.Time) time.Time {
	occ_wanted := WeekdayOccurenceNumber(t)
	weekday := t.Weekday()
	y := t.Year()
	var m time.Month
	if t.Month() == 12 {
		y++
		m = time.January
	} else {
		m = t.Month() + 1
	}

	t = time.Date(y, m, 1, t.Hour(), t.Minute(), t.Second(),
		t.Nanosecond(), t.Location())
	occ_actual := 0
	var last time.Time
	for t.Month() == m {
		if t.Weekday() == weekday {
			occ_actual++
			last = t
		}
		if occ_actual == occ_wanted {
			return t
		}
		t = t.AddDate(0, 0, 1)
	}

	return last
}

func monthlyByWeekdayCountIsOn(first time.Time, ison time.Time) bool {
	if first.Weekday() != ison.Weekday() {
		return false
	}

	first_occ := WeekdayOccurenceNumber(first)
	ison_occ := WeekdayOccurenceNumber(ison)
	if first_occ == ison_occ {
		return true
	}

	if first_occ == 5 && ison.AddDate(0, 0, 7).Month() != ison.Month() {
		return true
	}

	return false
}

func everyday(t time.Time) time.Time {
	return t.AddDate(0, 0, 1)
}

func everydayIsOn(first time.Time, ison time.Time) bool {
	return true
}
