package timefy

import "time"

// NewRule creates a new rule with the default configuration.
//
// Returns:
//   - A pointer to a `Config` struct, which includes the default configuration.
//
// Example:
//
//	rule := NewRule() // This creates a new rule with the default configuration.
//	rule.WithWeekStartDay(time.Monday) // This sets the week start day to Monday.
//	rule.WithTimeFormats([]string{"2006-01-02 15:04:05"}) // This sets the time formats to the default time formats.
//	rule.WithLocation(time.UTC) // This sets the time location to UTC.
//	return rule
func NewRule() *Rule {
	return &Rule{
		weekStartDay: WeekStartDay,
		timeFormats:  TimeFormats,
	}
}

// WithWeekStartDay sets the week start day.
//
// Parameters:
//   - `weekStartDay`: A time.Weekday value representing the week start day.
//
// Returns:
//   - A pointer to a `Config` struct, which includes the updated configuration.
//
// Example:
//
//	rule := NewRule() // This creates a new rule with the default configuration.
//	rule.WithWeekStartDay(time.Monday) // This sets the week start day to Monday.
//	return rule
func (c *Rule) WithWeekStartDay(weekStartDay time.Weekday) *Rule {
	c.weekStartDay = weekStartDay
	return c
}

// WithTimeFormats sets the time formats.
//
// Parameters:
//   - `timeFormats`: A slice of strings representing the time formats.
//
// Returns:
//   - A pointer to a `Config` struct, which includes the updated configuration.
//
// Example:
//
//	rule := NewRule() // This creates a new rule with the default configuration.
//	rule.WithTimeFormats([]string{"2006-01-02 15:04:05"}) // This sets the time formats to the default time formats.
//	return rule
func (c *Rule) WithTimeFormats(timeFormats []string) *Rule {
	c.timeFormats = timeFormats
	return c
}

// WithLocation sets the time location.
//
// Parameters:
//   - `timeLocation`: A pointer to a time.Location value representing the time location.
//
// Returns:
//   - A pointer to a `Config` struct, which includes the updated configuration.
//
// Example:
//
//	rule := NewRule() // This creates a new rule with the default configuration.
//	rule.WithLocation(time.UTC) // This sets the time location to UTC.
//	return rule
func (c *Rule) WithLocation(location *time.Location) *Rule {
	c.timeLocation = location
	return c
}

// WithLocationRFC sets the time location using the ZoneRFC value.
//
// Parameters:
//   - `location`: A ZoneRFC value representing the time location.
//
// Returns:
//   - A pointer to a `Config` struct, which includes the updated configuration.
//
// Example:
//
//	rule := NewRule() // This creates a new rule with the default configuration.
//	rule.WithLocationRFC(DefaultTimezoneVietnam) // This sets the time location to Vietnam.
//	return rule
func (c *Rule) WithLocationRFC(location ZoneRFC) *Rule {
	c.timeLocation, _ = time.LoadLocation(string(location))
	return c
}

// ApplyUTCLoc applies the UTC location to the configuration.
//
// Returns:
//   - A pointer to a `Config` struct, which includes the updated configuration.
//
// Example:
//
//	rule := NewRule() // This creates a new rule with the default configuration.
//	rule.ApplyUTCLoc() // This sets the time location to UTC.
//	return rule
func (c *Rule) ApplyUTCLoc() *Rule {
	c.timeLocation = time.UTC
	return c
}

// ApplyLocalLoc applies the local location to the configuration.
//
// Returns:
//   - A pointer to a `Config` struct, which includes the updated configuration.
//
// Example:
//
//	rule := NewRule() // This creates a new rule with the default configuration.
//	rule.ApplyLocalLoc() // This sets the time location to local.
//	return rule
func (c *Rule) ApplyLocalLoc() *Rule {
	c.timeLocation = time.Local
	return c
}

// AppendTimeFormat appends the time formats.
//
// Parameters:
//   - `format`: A variadic list of strings representing the time formats to append.
//
// Returns:
//   - A pointer to a `Rule` struct, which includes the updated configuration.
//
// Example:
//
//	rule := NewRule() // This creates a new rule with the default configuration.
//	rule.AppendTimeFormat("2006-01-02 15:04:05") // This appends the time format "2006-01-02 15:04:05".
//	return rule
func (c *Rule) AppendTimeFormat(format ...string) *Rule {
	if len(format) == 0 {
		return c
	}
	c.timeFormats = append(c.timeFormats, format...)
	return c
}

// AppendTimeFormatRFC appends the time formats RFC.
//
// Parameters:
//   - `format`: A variadic list of TimeFormatRFC values representing the time formats to append.
//
// Returns:
//   - A pointer to a `Rule` struct, which includes the updated configuration.
//
// Example:
//
//	rule := NewRule() // This creates a new rule with the default configuration.
//	rule.AppendTimeFormatRFC(TimeFormat20060102T150405999999, TimeFormat20060102T150405) // This appends the time formats "2006-01-02T15:04:05.999999" and "2006-01-02T15:04:05".
//	return rule
func (c *Rule) AppendTimeFormatRFC(format ...TimeFormatRFC) *Rule {
	if len(format) == 0 {
		return c
	}
	for _, f := range format {
		c.timeFormats = append(c.timeFormats, string(f))
	}
	return c
}

// AppendTimeFormatRFCshort appends the time formats RFC short.
//
// Parameters:
//   - `format`: A variadic list of strings representing the time formats to append.
//
// Returns:
//   - A pointer to a `Rule` struct, which includes the updated configuration.
//
// Example:
//
//	rule := NewRule() // This creates a new rule with the default configuration.
//	rule.AppendTimeFormatRFCshort(TimeRFC01T150405) // This appends the time format "15:04:05".
//	return rule
func (c *Rule) AppendTimeFormatRFCshort(format ...TimeRFC) *Rule {
	if len(format) == 0 {
		return c
	}
	for _, f := range format {
		c.timeFormats = append(c.timeFormats, string(f))
	}
	return c
}
