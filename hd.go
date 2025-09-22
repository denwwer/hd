// Humanize Duration (hd) – Go package that works like Duration.String() but returns
// a calendar-accurate difference (years, months, days, hours, minutes, seconds) for human-friendly output.
package hd

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Duration represents the difference in calendar units.
type Duration struct {
	Years, Months, Days     int
	Hours, Minutes, Seconds int
}

// String returns a string representing the duration in format "1y 6m 8d 5h 8m 17s".
func (d Duration) String() string {
	var s strings.Builder

	write := func(format string, v int) {
		if v != 0 {
			if s.Len() > 0 {
				s.WriteByte(' ')
			}
			fmt.Fprintf(&s, format, v)
		}
	}

	write("%dy", d.Years)
	write("%dm", d.Months)
	write("%dd", d.Days)
	write("%dh", d.Hours)
	write("%dm", d.Minutes)
	write("%ds", d.Seconds)

	if s.Len() == 0 {
		return "0s"
	}
	return s.String()
}

// MarshalJSON implements json.Marshaler for Duration.
func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

// Since returns a calendar-accurate duration from t to time.Now() in the given location zone.
func Since(t time.Time, loc *time.Location) Duration {
	return Between(t, time.Now(), loc)
}

// Between calculates the calendar-accurate duration
// between any two times (start -> end) in the given location. If loc is nil, time.UTC is used.
func Between(start, end time.Time, loc *time.Location) Duration {
	if loc == nil {
		loc = time.UTC
	}

	start = start.In(loc)
	end = end.In(loc)

	if start.After(end) {
		start, end = end, start
	}

	years, months, days := 0, 0, 0

	for !start.AddDate(years+1, 0, 0).After(end) {
		years++
	}
	for !start.AddDate(years, months+1, 0).After(end) {
		months++
	}
	for !start.AddDate(years, months, days+1).After(end) {
		days++
	}

	partial := end.Sub(start.AddDate(years, months, days))
	hours := int(partial.Hours())
	minutes := int(partial.Minutes()) % 60
	seconds := int(partial.Seconds()) % 60

	return Duration{
		Years:   years,
		Months:  months,
		Days:    days,
		Hours:   hours,
		Minutes: minutes,
		Seconds: seconds,
	}
}
