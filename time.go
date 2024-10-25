package timefy

import "time"

// With wraps the provided time value `v` into a Timex object, applying the default configuration.
//
// The function first checks if the global `DefaultConfig` is set. If it is not, it initializes a new
// `Config` with default values for `WeekStartDay` and `TimeFormats`. The configuration and time
// are then used to create a `Timex` object.
//
// Parameters:
//
//   - `v`: A time.Time value representing the time to be wrapped inside the Timex struct.
//
// Returns:
//
//   - A pointer to a `Timex` struct, which includes the provided time and either the default or
//     a custom configuration.
//
// Example:
//
//	t := time.Now()
//	timex := With(t) // This wraps the current time into a Timex object with the default configuration.
func With(v time.Time) *Timex {
	c := DefaultConfig
	if c == nil {
		c = &Config{
			WeekStartDay: WeekStartDay,
			TimeFormats:  TimeFormats,
		}
	}
	return &Timex{Time: v, Config: c}
}

// New creates a new Timex object for the provided time value `v`.
//
// The function calls the `With()` function, which wraps the given time in a `Timex` struct and applies
// the default configuration (or a custom one if set). Essentially, this function is a shorthand for
// creating a `Timex` object.
//
// Parameters:
//
//   - `v`: A time.Time value representing the time to be wrapped inside the Timex struct.
//
// Returns:
//   - A pointer to a `Timex` struct, which includes the provided time and the applicable configuration.
//
// Example:
//
//	t := time.Now()
//	timex := New(t) // This creates a new Timex object for the current time.
func New(v time.Time) *Timex {
	return With(v)
}

// With wraps the provided time value `v` into a Timex object using the current configuration `c`.
//
// This method is called on a `Config` object to create a new `Timex` instance, where the provided
// time and the configuration are bundled together. Unlike the global `With()` function, this method
// uses the specific configuration tied to the `Config` instance `c`.
//
// Parameters:
//
//   - `v`: A time.Time value representing the time to be wrapped inside the Timex struct.
//
// Returns:
//   - A pointer to a `Timex` struct, which includes the provided time and the specific `Config` instance.
//
// Example:
//
//	customConfig := &Config{WeekStartDay: time.Monday}
//	t := time.Now()
//	timex := customConfig.With(t) // This creates a Timex object with the custom configuration.
func (c *Config) With(v time.Time) *Timex {
	return &Timex{Time: v, Config: c}
}

// Parse attempts to parse a given set of string representations of time using the current `Config`.
//
// The function checks whether the `TimeLocation` field in the configuration `c` is set. If it is not set,
// the function defaults to parsing the strings based on the current time (using the system's local time zone).
// If `TimeLocation` is set, it uses the specified location to interpret the time. It delegates the parsing
// task to the `Parse()` method of a new `Timex` object created using `With()`, which wraps the current time
// with the existing configuration.
//
// Parameters:
//
//   - `s`: A variadic list of strings representing dates or times to be parsed.
//
// Returns:
//   - A `time.Time` value if one of the provided strings is successfully parsed.
//   - An error if parsing fails.
//
// Example:
//
//	config := &Config{TimeLocation: time.UTC}
//	parsedTime, err := config.Parse("2023-10-24T12:00:00Z") // Parses the string using UTC time location.
//
// If no `TimeLocation` is set, the function uses the local time zone:
//
//	config := &Config{}
//	parsedTime, err := config.Parse("2023-10-24T12:00:00") // Parses using the local time zone.
func (c *Config) Parse(s ...string) (time.Time, error) {
	if c.TimeLocation == nil {
		return c.With(time.Now()).Parse(s...)
	} else {
		return c.With(time.Now().In(c.TimeLocation)).Parse(s...)
	}
}

// MustParse attempts to parse a given set of string representations of time using the current `Config`,
// and panics if parsing fails.
//
// This function checks if the `TimeLocation` field in the configuration `c` is set. If it is not set,
// it uses the system's local time zone for parsing. If `TimeLocation` is set, it utilizes the specified
// location to interpret the time. The actual parsing task is delegated to the `MustParse()` method
// of a new `Timex` object created using `With()`, which wraps the current time with the existing
// configuration.
//
// Parameters:
//
//   - `s`: A variadic list of strings representing dates or times to be parsed.
//
// Returns:
//   - A `time.Time` value representing the parsed time. This function will panic if parsing fails.
//
// Example:
//
//	config := &Config{TimeLocation: time.UTC}
//	parsedTime := config.MustParse("2023-10-24T12:00:00Z") // Parses the string using UTC time location.
//
//	 If no `TimeLocation` is set, the function uses the local time zone:
//
//	config := &Config{}
//	parsedTime := config.MustParse("2023-10-24T12:00:00") // Parses using the local time zone, panicking on failure.
func (c *Config) MustParse(s ...string) time.Time {
	if c.TimeLocation == nil {
		return c.With(time.Now()).MustParse(s...)
	} else {
		return c.With(time.Now().In(c.TimeLocation)).MustParse(s...)
	}
}

// BeginningOfMinute returns a new time.Time value representing the start of the minute for the
// given Timex instance.
//
// The function uses the `Truncate()` method on the underlying time.Time of the Timex struct,
// truncating the time to the nearest minute. This effectively resets the seconds and
// nanoseconds to zero, providing the time corresponding to the beginning of the current minute.
//
// Returns:
//   - A `time.Time` value representing the start of the minute for the current Timex instance.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	startOfMinute := t.BeginningOfMinute() // This will return the current time truncated to the start of the minute.
func (t *Timex) BeginningOfMinute() time.Time {
	return t.Truncate(time.Minute)
}

// BeginningOfHour returns a new time.Time value representing the start of the hour for the
// given Timex instance.
//
// The function retrieves the year, month, and day from the underlying time.Time of the Timex struct
// using the `Date()` method. It then constructs a new time.Time value with the same year, month, and
// day but sets the hour to the current hour and resets the minutes, seconds, and nanoseconds to zero.
// This effectively provides the time corresponding to the beginning of the current hour.
//
// Returns:
//   - A `time.Time` value representing the start of the hour for the current Timex instance.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	startOfHour := t.BeginningOfHour() // This will return the current time truncated to the start of the hour.
func (t *Timex) BeginningOfHour() time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, t.Time.Hour(), 0, 0, 0, t.Time.Location())
}

// BeginningOfDay returns a new time.Time value representing the start of the day for the
// given Timex instance.
//
// The function retrieves the year, month, and day from the underlying time.Time of the Timex struct
// using the `Date()` method. It then constructs a new time.Time value with the same year, month, and
// day but sets the hour, minute, second, and nanosecond to zero. This effectively resets the time to
// midnight (00:00:00) of the same date, providing the start of the day.
//
// Returns:
//   - A `time.Time` value representing the start of the day for the current Timex instance.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	startOfDay := t.BeginningOfDay() // This will return the current date at 00:00:00 (midnight).
func (t *Timex) BeginningOfDay() time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Time.Location())
}

// BeginningOfWeek returns a new time.Time value representing the start of the week for the
//
//	given Timex instance, based on the configured `WeekStartDay`.
//
// The function first calculates the beginning of the current day using `BeginningOfDay()`. Then, it determines
// the current weekday and calculates how many days to subtract in order to reach the start of the week, based on
// the configured start day (`WeekStartDay`). If `WeekStartDay` is set to a day other than Sunday, the function
// adjusts the number of days to move back accordingly.
//
// Returns:
//
//   - A `time.Time` value representing the start of the week for the current Timex instance.
//
// Example:
//
//	t := Timex{Time: time.Now(), Config: &Config{WeekStartDay: time.Monday}}
//	startOfWeek := t.BeginningOfWeek() // This will return the date and time at the start of the current week.
//
// Note:
//
//	If `WeekStartDay` is not set (defaults to Sunday), the function will return the preceding Sunday.
func (t *Timex) BeginningOfWeek() time.Time {
	day := t.BeginningOfDay()
	weekday := int(day.Weekday())
	if t.WeekStartDay != time.Sunday {
		weekStartDayInt := int(t.WeekStartDay)
		if weekday < weekStartDayInt {
			weekday = weekday + 7 - weekStartDayInt
		} else {
			weekday = weekday - weekStartDayInt
		}
	}
	return day.AddDate(0, 0, -weekday)
}

// BeginningOfMonth returns a new time.Time value representing the start of the month for the
// given Timex instance.
//
// The function retrieves the year and month from the underlying time.Time of the Timex struct using the `Date()`
// method. It then constructs a new time.Time value with the same year and month, but sets the day to the first of
// the month and resets the hour, minute, second, and nanosecond to zero, providing the start of the month.
//
// Returns:
//   - A `time.Time` value representing the start of the month for the current Timex instance.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	startOfMonth := t.BeginningOfMonth() // This will return the date and time at the start of the current month.
func (t *Timex) BeginningOfMonth() time.Time {
	y, m, _ := t.Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, t.Location())
}

// BeginningOfQuarter returns a new time.Time value representing the start of the quarter
// for the given Timex instance.
//
// The function first calculates the start of the current month using `BeginningOfMonth()`. It then determines
// how far the current month is into the quarter by calculating the offset (based on a 3-month quarter system).
// The offset is subtracted from the current month to return a time value corresponding to the start of the quarter
// (e.g., if the current month is April, it will return the start of April; if it's June, it will return the start of April).
//
// Returns:
//   - A `time.Time` value representing the start of the quarter for the current Timex instance.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	startOfQuarter := t.BeginningOfQuarter() // This will return the date and time at the start of the current quarter.
func (t *Timex) BeginningOfQuarter() time.Time {
	month := t.BeginningOfMonth()
	offset := (int(month.Month()) - 1) % 3
	return month.AddDate(0, -offset, 0)
}

// BeginningOfHalf returns a new time.Time value representing the start of the half-year
// for the given Timex instance.
//
// The function first calculates the start of the current month using `BeginningOfMonth()`. It then determines
// how far the current month is into the half-year by calculating the offset (based on a 6-month system).
// The offset is subtracted from the current month to return a time value corresponding to the start of the half-year
// (e.g., if the current month is July, it will return the start of July; if it's October, it will return the start of July).
//
// Returns:
//   - A `time.Time` value representing the start of the half-year for the current Timex instance.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	startOfHalf := t.BeginningOfHalf() // This will return the date and time at the start of the current half-year.
func (t *Timex) BeginningOfHalf() time.Time {
	month := t.BeginningOfMonth()
	offset := (int(month.Month()) - 1) % 6
	return month.AddDate(0, -offset, 0)
}

// BeginningOfYear returns a new time.Time value representing the start of the year
// for the given Timex instance.
//
// The function retrieves the year from the underlying time.Time of the Timex struct using the `Date()` method.
// It then constructs a new time.Time value with the same year but sets the month to January, and the day, hour,
// minute, second, and nanosecond to zero. This effectively resets the time to midnight (00:00:00) on January 1st
// of the current year.
//
// Returns:
//   - A `time.Time` value representing the start of the year for the current Timex instance.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	startOfYear := t.BeginningOfYear() // This will return January 1st of the current year at 00:00:00.
func (t *Timex) BeginningOfYear() time.Time {
	y, _, _ := t.Date()
	return time.Date(y, time.January, 1, 0, 0, 0, 0, t.Location())
}

// EndOfMinute returns a new time.Time value representing the end of the current minute
// for the given Timex instance.
//
// The function first calculates the start of the minute using `BeginningOfMinute()`, then adds one minute
// and subtracts a nanosecond. This adjustment effectively gives the time at the last nanosecond of the current minute.
//
// Returns:
//   - A `time.Time` value representing the end of the current minute for the Timex instance.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	endOfMinute := t.EndOfMinute() // This will return the time at the end of the current minute (59.999999999 nanoseconds).
func (t *Timex) EndOfMinute() time.Time {
	return t.BeginningOfMinute().Add(time.Minute - time.Nanosecond)
}
