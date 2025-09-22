package hd

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// helper to create times in UTC quickly
func utc(y int, m time.Month, d, h, min, s int) time.Time {
	return time.Date(y, m, d, h, min, s, 0, time.UTC)
}

func TestBetween(t *testing.T) {
	t.Parallel()

	start := utc(2020, 1, 1, 0, 0, 0)
	end := utc(2023, 4, 15, 6, 30, 15)

	d := Between(start, end, nil)
	assert.Equal(t, Duration{
		Years:   3,
		Months:  3,
		Days:    14,
		Hours:   6,
		Minutes: 30,
		Seconds: 15,
	}, d)
}

func TestBetween_NilLocationUTC(t *testing.T) {
	t.Parallel()

	start := utc(2022, 3, 14, 10, 0, 0)
	end := utc(2023, 3, 14, 12, 0, 0)

	d := Between(start, end, nil)
	assert.Equal(t, 1, d.Years)
	assert.Equal(t, 0, d.Months)
	assert.Equal(t, 0, d.Days)
	assert.Equal(t, 2, d.Hours)
}

func TestBetween_WithLocation(t *testing.T) {
	loc, err := time.LoadLocation("America/New_York")
	require.Nil(t, err)

	start := utc(2025, 9, 22, 0, 0, 0)
	end := utc(2025, 9, 23, 12, 30, 15)

	d := Between(start, end, loc)
	assert.Equal(t, 0, d.Years)
	assert.Equal(t, 0, d.Months)
	assert.Equal(t, 1, d.Days)
	assert.Equal(t, 12, d.Hours)
	assert.Equal(t, 30, d.Minutes)
	assert.Equal(t, 15, d.Seconds)
}

func TestBetween_ReverseOrder(t *testing.T) {
	t.Parallel()
	// Start is after end; should still produce positive difference
	start := utc(2025, 1, 1, 0, 0, 0)
	end := utc(2023, 1, 1, 0, 0, 0)

	d := Between(start, end, nil)
	assert.Equal(t, 2, d.Years)
}

func TestBetween_LeapYear(t *testing.T) {
	t.Parallel()
	// Feb 29, 2020 → Feb 28, 2021 should be 1 year minus 1 day
	start := utc(2020, 2, 29, 0, 0, 0)
	end := utc(2021, 2, 28, 0, 0, 0)

	d := Between(start, end, time.UTC)
	assert.Equal(t, 0, d.Years)
	assert.Equal(t, 11, d.Months)
}

func TestSince(t *testing.T) {
	t.Parallel()

	d := Since(time.Now().AddDate(0, 0, -2).Add(-time.Hour), nil)
	assert.Equal(t, "2d 1h", d.String())
}

func TestDuration_String(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		dur    Duration
		expect string
	}{
		{
			name:   "All non-zero",
			dur:    Duration{Years: 1, Months: 2, Days: 3, Hours: 4, Minutes: 5, Seconds: 6},
			expect: "1y 2m 3d 4h 5m 6s",
		},
		{
			name:   "Some zeros",
			dur:    Duration{Years: 0, Months: 0, Days: 3, Hours: 0, Minutes: 5, Seconds: 0},
			expect: "3d 5m",
		},
		{
			name:   "All zeros",
			dur:    Duration{},
			expect: "0s",
		},
		{
			name:   "Single unit",
			dur:    Duration{Hours: 7},
			expect: "7h",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expect, tc.dur.String())
		})
	}
}

func TestDuration_MarshalJSON(t *testing.T) {
	t.Parallel()

	start := utc(2020, 1, 1, 0, 0, 0)
	end := utc(2023, 4, 15, 6, 30, 15)

	data := struct {
		Duration Duration `json:"duration"`
	}{
		Duration: Between(start, end, nil),
	}

	b, err := json.Marshal(data)
	require.Nil(t, err)
	assert.Equal(t, `{"duration":"3y 3m 14d 6h 30m 15s"}`, string(b))
}
