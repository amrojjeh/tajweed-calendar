package cal

import (
	"testing"
	"time"

	"github.com/amrojjeh/yaag/internal/assert"
)

func TestFirstWeekdayInMonth(t *testing.T) {
	tests := []struct {
		name  string
		year  int
		month time.Month
		want  time.Weekday
	}{
		{
			name:  "January",
			year:  2024,
			month: time.January,
			want:  time.Monday,
		},
		{
			name:  "March",
			year:  2024,
			month: time.March,
			want:  time.Friday,
		},
		{
			name:  "August",
			year:  2024,
			month: time.August,
			want:  time.Thursday,
		},
		{
			name:  "December",
			year:  2024,
			month: time.December,
			want:  time.Sunday,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, FirstWeekdayInMonth(tt.year, tt.month), tt.want)
		})
	}
}
