package timefy

import "time"

// Rule rule for timefy
// weekStartDay set week start day, default is sunday
// timeLocation set time location, default is local
// timeFormats set time formats, default is []string{"2006-01-02 15:04:05"}
type Rule struct {
	weekStartDay time.Weekday
	timeLocation *time.Location
	timeFormats  []string
}

// Timex struct
type Timex struct {
	time.Time
	*Rule
}
