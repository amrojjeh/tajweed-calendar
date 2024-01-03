package cal

import (
	"testing"
	"time"

	"github.com/amrojjeh/tajweed-calendar/internal/assert"
)

func TestEvent_IsOn(t *testing.T) {
	tests := []struct {
		name  string
		event Event
		year  int
		month time.Month
		day   int
		want  bool
	}{
		{
			name: "Single Event - True",
			event: Event{
				Name:      "Soccer",
				Next:      nil,
				Color:     "soccer",
				FirstDate: NewEventDate(2024, time.January, 13, 13, 30),
				Duration:  1.0,
			},
			year:  2024,
			month: time.January,
			day:   13,
			want:  true,
		},
		{
			name: "Single Event - False",
			event: Event{
				Name:      "Soccer",
				Next:      nil,
				Color:     "soccer",
				FirstDate: NewEventDate(2024, time.January, 13, 13, 30),
				Duration:  1.0,
			},
			year:  2024,
			month: time.February,
			day:   13,
			want:  false,
		},
		{
			name: "Weekly event - First",
			event: Event{
				Name:      "Soccer",
				Next:      Weekly,
				Color:     "soccer",
				FirstDate: NewEventDate(2024, time.January, 13, 13, 30),
				Duration:  1.0,
			},
			year:  2024,
			month: time.January,
			day:   13,
			want:  true,
		},
		{
			name: "Weekly event - Second",
			event: Event{
				Name:      "Soccer",
				Next:      Weekly,
				Color:     "soccer",
				FirstDate: NewEventDate(2024, time.January, 13, 13, 30),
				Duration:  1.0,
			},
			year:  2024,
			month: time.January,
			day:   20,
			want:  true,
		},
		{
			name: "Weekly event - Next Month",
			event: Event{
				Name:      "Soccer",
				Next:      Weekly,
				Color:     "soccer",
				FirstDate: NewEventDate(2024, time.January, 13, 13, 30),
				Duration:  1.0,
			},
			year:  2024,
			month: time.February,
			day:   3,
			want:  true,
		},
		{
			name: "Weekly event - False",
			event: Event{
				Name:      "Soccer",
				Next:      Weekly,
				Color:     "soccer",
				FirstDate: NewEventDate(2024, time.January, 13, 13, 30),
				Duration:  1.0,
			},
			year:  2024,
			month: time.February,
			day:   4,
			want:  false,
		},
		{
			name: "Monthly by weekend count event - First",
			event: Event{
				Name:      "Soccer",
				Next:      MonthlyByWeekdayCount,
				Color:     "soccer",
				FirstDate: NewEventDate(2024, time.January, 13, 13, 30),
				Duration:  1.0,
			},
			year:  2024,
			month: time.January,
			day:   13,
			want:  true,
		},
		{
			name: "Monthly by weekend count event - Second",
			event: Event{
				Name:      "Soccer",
				Next:      MonthlyByWeekdayCount,
				Color:     "soccer",
				FirstDate: NewEventDate(2024, time.January, 13, 13, 30),
				Duration:  1.0,
			},
			year:  2024,
			month: time.February,
			day:   10,
			want:  true,
		},
		{
			name: "Monthly by weekend count event - Fifth Day",
			event: Event{
				Name:      "Soccer",
				Next:      MonthlyByWeekdayCount,
				Color:     "soccer",
				FirstDate: NewEventDate(2024, time.January, 31, 13, 30),
				Duration:  1.0,
			},
			year:  2024,
			month: time.February,
			day:   28,
			want:  true,
		},
		{
			name: "Monthly by weekend count event - False",
			event: Event{
				Name:      "Soccer",
				Next:      MonthlyByWeekdayCount,
				Color:     "soccer",
				FirstDate: NewEventDate(2024, time.January, 13, 13, 30),
				Duration:  1.0,
			},
			year:  2024,
			month: time.February,
			day:   28,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.event.IsOn(tt.year, tt.month, tt.day)
			assert.Equal(t, actual, tt.want)
		})
	}
}
