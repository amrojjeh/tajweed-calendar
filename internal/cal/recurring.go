package cal

import (
	"encoding/json"
	"errors"
	"time"
)

// TODO(Amr Ojjeh): Turn NextFunc into an interface such that it
// supports both a "next" method as well as an "IsOn" method

type NextFunc func(time.Time) time.Time

func (n *NextFunc) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		err = errors.Join(errors.New("cal: NextFunc: failed to parse JSON"), err)
		return err
	}
	switch s {
	case "":
		*n = nil
	case "weekly":
		*n = Weekly
	case "monthly by weekday count":
		*n = MonthlyByWeekdayCount
	case "everyday":
		*n = Everyday
	default:
		*n = nil
	}
	return nil
}

func Weekly(t time.Time) time.Time {
	return t.AddDate(0, 0, 7)
}

func MonthlyByWeekdayCount(t time.Time) time.Time {
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

func Everyday(t time.Time) time.Time {
	return t.AddDate(0, 0, 1)
}
