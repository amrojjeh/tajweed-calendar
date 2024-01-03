package cal

import (
	"time"
)

var Months = map[time.Month]int{time.January: 31, time.February: 29,
	time.March: 31, time.April: 30, time.May: 31, time.June: 30, time.July: 31,
	time.August: 31, time.September: 30, time.October: 31, time.November: 30,
	time.December: 31}

func FirstWeekdayInMonth(y int, m time.Month) time.Weekday {
	d := time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)
	return d.Weekday()
}

func WeekdayOccurenceNumber(t time.Time) int {
	w := 1
	tt := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	for t.Day() > tt.Day() {
		if tt.Weekday() == t.Weekday() {
			w++
		}
		tt = tt.Add(time.Hour * 24)
	}

	return w
}
