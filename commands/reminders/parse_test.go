// SPDX-License-Identifier: AGPL-3.0-only
package reminders

import (
	"testing"
	"time"
)

// all these are rounded to second precision to make sure they don't fail due to any slowdown
func TestParseTime(t *testing.T) {
	tests := []struct {
		Name     string
		Input    []string
		WantTime time.Time
		WantPos  int
		Timezone *time.Location
	}{
		{
			"test tomorrow",
			[]string{"tomorrow"},
			time.Now().UTC().Add(24 * time.Hour),
			0, time.UTC,
		},
		{
			"test random time 1, utc",
			[]string{"january", "5", "2022", "at", "3:00pm"},
			time.Date(2022, 1, 5, 15, 0, 0, 0, time.UTC),
			4, time.UTC,
		},
		{
			"test random time 1, timezone",
			[]string{"january", "5", "2022", "at", "3:20pm"},
			time.Date(2022, 1, 5, 15, 20, 0, 0, mustLoadLocation("Europe/Amsterdam")),
			4, mustLoadLocation("Europe/Amsterdam"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			got, gotPos, err := ParseTime(tt.Input, tt.Timezone)
			if err != nil {
				t.Errorf("want nil error, got %v", err)
			}

			if !got.Round(time.Second).Equal(tt.WantTime.Round(time.Second)) {
				t.Errorf("want time %s, got time %s", tt.WantTime.Round(time.Second), got.Round(time.Second))
			}

			if gotPos != tt.WantPos {
				t.Errorf("want position %d, got position %d", gotPos, tt.WantPos)
			}
		})
	}
}

func mustLoadLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic(err)
	}
	return loc
}
