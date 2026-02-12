package timefy

import "time"

// Config configuration for now package
type Config struct {
	WeekStartDay time.Weekday   `json:"week_start_day,omitempty"`
	TimeLocation *time.Location `json:"time_location,omitempty"`
	TimeFormats  []string       `json:"time_formats,omitempty"`
}

// Timex now struct
type Timex struct {
	time.Time
	*Config
}
