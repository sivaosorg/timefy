package timefy

import (
	"fmt"
	"time"
)

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
	rule := NewRule()
	return &Timex{Time: v, Rule: rule}
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
func (c *Rule) With(v time.Time) *Timex {
	return &Timex{Time: v, Rule: c}
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
func (c *Rule) Parse(s ...string) (time.Time, error) {
	if c.timeLocation == nil {
		return c.With(time.Now()).Parse(s...)
	} else {
		return c.With(time.Now().In(c.timeLocation)).Parse(s...)
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
func (c *Rule) MustParse(s ...string) time.Time {
	if c.timeLocation == nil {
		return c.With(time.Now()).MustParse(s...)
	} else {
		return c.With(time.Now().In(c.timeLocation)).MustParse(s...)
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
	if t.weekStartDay != time.Sunday {
		weekStartDayInt := int(t.weekStartDay)
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

// EndOfHour returns a new time.Time value representing the end of the current hour
// for the given Timex instance.
//
// The function first calculates the start of the hour using `BeginningOfHour()`, then adds one hour
// and subtracts a nanosecond to reach the last nanosecond of the current hour.
//
// Returns:
//   - A `time.Time` value representing the end of the current hour for the Timex instance.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	endOfHour := t.EndOfHour() // This will return the time at the end of the current hour (59:59.999999999).
func (t *Timex) EndOfHour() time.Time {
	return t.BeginningOfHour().Add(time.Hour - time.Nanosecond)
}

// EndOfDay returns a new time.Time value representing the end of the current day
// for the given Timex instance.
//
// The function retrieves the year, month, and day from the underlying time.Time of the Timex struct
// using the `Date()` method. It then constructs a new time.Time value with the same year, month, and day
// but sets the time to 23:59:59.999999999, which is the last nanosecond of the day.
//
// Returns:
//   - A `time.Time` value representing the end of the current day for the Timex instance.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	endOfDay := t.EndOfDay() // This will return the date and time at the end of the current day (23:59:59.999999999).
func (t *Timex) EndOfDay() time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 23, 59, 59, int(time.Second-time.Nanosecond), t.Location())
}

// EndOfWeek returns a new time.Time value representing the end of the current week
// for the given Timex instance, based on the configured `WeekStartDay`.
//
// The function first calculates the start of the current week using `BeginningOfWeek()`,
// then adds 7 days to reach the start of the next week, and subtracts one nanosecond
// to get the last nanosecond of the current week.
//
// Returns:
//   - A `time.Time` value representing the end of the current week for the Timex instance.
//
// Example:
//
//	t := Timex{Time: time.Now(), Config: &Config{WeekStartDay: time.Monday}}
//	endOfWeek := t.EndOfWeek() // This will return the date and time at the end of the current week.
//
// Note:
//
//	If `WeekStartDay` is not set (defaults to Sunday), the function will return the preceding Saturday at 23:59:59.999999999.
func (t *Timex) EndOfWeek() time.Time {
	return t.BeginningOfWeek().AddDate(0, 0, 7).Add(-time.Nanosecond)
}

// EndOfMonth returns a new time.Time value representing the end of the current month
// for the given Timex instance.
//
// The function first calculates the start of the month using `BeginningOfMonth()`,
// then adds one month to reach the beginning of the next month, and subtracts one nanosecond
// to obtain the last nanosecond of the current month.
//
// Returns:
//   - A `time.Time` value representing the end of the current month for the Timex instance.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	endOfMonth := t.EndOfMonth() // This will return the date and time at the end of the current month (e.g., 23:59:59.999999999 on the last day).
func (t *Timex) EndOfMonth() time.Time {
	return t.BeginningOfMonth().AddDate(0, 1, 0).Add(-time.Nanosecond)
}

// EndOfQuarter returns a new time.Time value representing the end of the current quarter
// for the given Timex instance.
//
// The function first calculates the start of the quarter using `BeginningOfQuarter()`,
// then adds three months to reach the beginning of the next quarter, and subtracts one nanosecond
// to obtain the last nanosecond of the current quarter.
//
// Returns:
//   - A `time.Time` value representing the end of the current quarter for the Timex instance.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	endOfQuarter := t.EndOfQuarter() // This will return the date and time at the end of the current quarter (e.g., 23:59:59.999999999 on the last day of the quarter).
func (t *Timex) EndOfQuarter() time.Time {
	return t.BeginningOfQuarter().AddDate(0, 3, 0).Add(-time.Nanosecond)
}

// EndOfHalf returns a new time.Time value representing the end of the current half-year
// for the given Timex instance.
//
// The function first calculates the start of the half-year using `BeginningOfHalf()`, then adds six months
// to reach the beginning of the next half-year, and subtracts one nanosecond to obtain the last nanosecond
// of the current half-year.
//
// Returns:
//   - A `time.Time` value representing the end of the current half-year for the Timex instance.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	endOfHalf := t.EndOfHalf() // This will return the date and time at the end of the current half-year (e.g., 23:59:59.999999999 on June 30th or December 31st).
func (t *Timex) EndOfHalf() time.Time {
	return t.BeginningOfHalf().AddDate(0, 6, 0).Add(-time.Nanosecond)
}

// EndOfYear returns a new time.Time value representing the end of the current year
// for the given Timex instance.
//
// The function first calculates the start of the year using `BeginningOfYear()`, then adds one year
// to reach the beginning of the next year, and subtracts one nanosecond to obtain the last nanosecond
// of the current year.
//
// Returns:
//   - A `time.Time` value representing the end of the current year for the Timex instance.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	endOfYear := t.EndOfYear() // This will return the date and time at the end of the current year (e.g., 23:59:59.999999999 on December 31st).
func (t *Timex) EndOfYear() time.Time {
	return t.BeginningOfYear().AddDate(1, 0, 0).Add(-time.Nanosecond)
}

// Monday returns a new time.Time value representing the most recent Monday at the start of the week
// based on the provided date string(s) or the current date if no date strings are provided.
//
// If a date string is provided, it parses the date using the Timex `Parse` method. If no string is provided,
// the function defaults to the start of the current day. It then determines the weekday of the parsed date.
// If the weekday is Sunday (represented by 0), it adjusts to be the seventh day for calculation purposes.
// The function finally subtracts the necessary days to reach the last Monday.
//
// Parameters:
//   - s ...string: Optional date string(s) to parse; if none are provided, it defaults to today.
//
// Returns:
//   - A `time.Time` value representing the date and time at the start of the last Monday.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	lastMonday := t.Monday() // Returns the start of the most recent Monday.
//
// Note:
//
//	If parsing fails, the function will panic with an error.
func (t *Timex) Monday(s ...string) time.Time {
	var parseTime time.Time
	var err error
	if len(s) > 0 {
		parseTime, err = t.Parse(s...)
		if err != nil {
			panic(err)
		}
	} else {
		parseTime = t.BeginningOfDay()
	}
	weekday := int(parseTime.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return parseTime.AddDate(0, 0, -weekday+1)
}

// Sunday returns a new time.Time value representing the most recent or upcoming Sunday
// based on the provided date string(s) or the current date if no date strings are provided.
//
// If a date string is provided, it parses the date using the Timex `Parse` method. If no string is provided,
// the function defaults to the start of the current day. It then determines the weekday of the parsed date.
// If the weekday is Sunday (represented by 0), it adjusts to be the seventh day for calculation purposes.
// The function finally adds the necessary days to reach the next Sunday.
//
// Parameters:
//   - s ...string: Optional date string(s) to parse; if none are provided, it defaults to today.
//
// Returns:
//   - A `time.Time` value representing the date and time at the start of the most recent or upcoming Sunday.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	lastSunday := t.Sunday() // Returns the date and time at the start of the most recent or upcoming Sunday.
//
// Note:
//
//	If parsing fails, the function will panic with an error.
func (t *Timex) Sunday(s ...string) time.Time {
	var parseTime time.Time
	var err error
	if len(s) > 0 {
		parseTime, err = t.Parse(s...)
		if err != nil {
			panic(err)
		}
	} else {
		parseTime = t.BeginningOfDay()
	}
	weekday := int(parseTime.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return parseTime.AddDate(0, 0, (7 - weekday))
}

// EndOfSunday returns a new time.Time value representing the end of the most recent or upcoming Sunday
// for the given Timex instance.
//
// The function first calculates the start of the most recent or upcoming Sunday using `Sunday()`,
// then constructs a new Timex instance with that date and time. Finally, it calls `EndOfDay()` on this
// instance to retrieve the time at the last nanosecond of the Sunday.
//
// Returns:
//   - A `time.Time` value representing the end of the most recent or upcoming Sunday for the Timex instance.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	endOfSunday := t.EndOfSunday() // Returns the date and time at the end of the most recent or upcoming Sunday (23:59:59.999999999).
func (t *Timex) EndOfSunday() time.Time {
	return New(t.Sunday()).EndOfDay()
}

// Quarter returns the current quarter of the year for the given Timex instance,
// where the year is divided into four quarters: Q1 (January-March), Q2 (April-June),
// Q3 (July-September), and Q4 (October-December).
//
// The function calculates the quarter by taking the current month (retrieved using
// the `Month()` method), subtracting one (to account for zero-based indexing),
// dividing by 3, and then adding 1 to get the quarter number. This results in
// a value between 1 and 4, representing the quarter of the year.
//
// Returns:
//   - A `uint` value representing the current quarter (1 to 4) for the Timex instance.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	quarter := t.Quarter() // Returns the current quarter number (1, 2, 3, or 4).
func (t *Timex) Quarter() uint {
	return (uint(t.Month())-1)/3 + 1
}

// Parse interprets the provided date string(s) and converts them into a time.Time value.
// It attempts to parse each string according to the configured formats, adjusting for the current time
// and location as necessary.
//
// The function checks each string for time information and adjusts the date components accordingly.
// If a date component is missing (e.g., day, month, year) in the parsed string, it will replace
// it with the corresponding component from the current time. If the input string includes only time,
// it sets the day and month to those of the current date. The function handles various cases for
// the input strings and parses them into a valid time.Time value.
//
// Parameters:
//   - `s ...string`: One or more date strings to be parsed. The function will try to parse each string
//     in the order provided and will return the first successful parsed time.
//
// Returns:
//   - `value`: A `time.Time` value representing the parsed date and time.
//   - `err`: An error value indicating any issues that occurred during parsing; if parsing is successful,
//     this will be nil.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	parsedTime, err := t.Parse("2023-10-25 15:04", "15:04") // Attempts to parse the given date strings.
//	if err != nil {
//		// Handle error
//	}
//
// Note:
// - The function modifies the parsed date based on the current time when certain components are missing.
// - It will return the most recent successful parsed value or the zero value of time.Time if none succeed.
func (t *Timex) Parse(s ...string) (value time.Time, err error) {
	var (
		setCurrentTime  bool
		parseTime       []int
		currentLocation = t.Location()
		onlyTimeInStr   = true
		currentTime     = FormatTimex(t.Time)
	)

	for _, str := range s {
		hasTimeInStr := TimeFormatRegexp.MatchString(str) // match 15:04:05, 15
		onlyTimeInStr = hasTimeInStr && onlyTimeInStr && TimeOnlyRegexp.MatchString(str)
		if value, err = t.parseWithFormat(str, currentLocation); err == nil {
			location := value.Location()
			parseTime = FormatTimex(value)

			for i, v := range parseTime {
				// Don't reset hour, minute, second if current time str including time
				if hasTimeInStr && i <= 3 {
					continue
				}

				// If value is zero, replace it with current time
				if v == 0 {
					if setCurrentTime {
						parseTime[i] = currentTime[i]
					}
				} else {
					setCurrentTime = true
				}

				// if current time only includes time, should change day, month to current time
				if onlyTimeInStr {
					if i == 4 || i == 5 {
						parseTime[i] = currentTime[i]
						continue
					}
				}
			}
			value = time.Date(parseTime[6], time.Month(parseTime[5]), parseTime[4], parseTime[3], parseTime[2], parseTime[1], parseTime[0], location)
			currentTime = FormatTimex(value)
		}
	}
	return
}

// MustParse interprets the provided date string(s) and converts them into a time.Time value,
// but it panics if parsing fails. This is useful for scenarios where a valid time is essential,
// and failure to parse should result in a program error.
//
// The function calls the `Parse` method, which attempts to interpret each date string according
// to the configured formats, adjusting for the current time and location as necessary. If
// the parsing succeeds, the resulting time.Time value is returned; if it fails, a panic is triggered
// with the associated error message.
//
// Parameters:
//   - `s ...string`: One or more date strings to be parsed. The function will try to parse each string
//     in the order provided and will panic if none can be successfully parsed.
//
// Returns:
//   - `v`: A `time.Time` value representing the parsed date and time, guaranteed to be valid
//     as the function will panic on failure.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	parsedTime := t.MustParse("2023-10-25 15:04") // Attempts to parse the given date strings.
//	// If parsing fails, it will panic.
//
// Note:
//   - This function is intended for use cases where the caller must have a valid time returned,
//     and any parsing errors should be treated as unrecoverable conditions.
func (t *Timex) MustParse(s ...string) (v time.Time) {
	v, err := t.Parse(s...)
	if err != nil {
		panic(err)
	}
	return v
}

// Between checks if the current Timex instance's time falls within the specified time range,
// defined by the provided `begin` and `end` date strings. The function parses the strings into
// time.Time values and determines if the current time is after the `begin` time and before the `end` time.
//
// Parameters:
//   - `begin`: A string representing the start of the time range; it will be parsed into a time.Time value.
//   - `end`: A string representing the end of the time range; it will also be parsed into a time.Time value.
//
// Returns:
//   - A boolean value indicating whether the current time is within the specified range (exclusive).
//     It returns true if the current time is after `begin` and before `end`; otherwise, it returns false.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	isInRange := t.Between("2023-10-01", "2023-10-31") // Checks if the current time is in October 2023.
//
// Note:
//   - If the parsing of `begin` or `end` fails, the function will panic due to the usage of `MustParse`.
func (t *Timex) Between(begin, end string) bool {
	beginTime := t.MustParse(begin)
	endTime := t.MustParse(end)
	return t.After(beginTime) && t.Before(endTime)
}

// FormatRFC formats the time using the provided layout.
//
// Parameters:
//   - `layout`: A TimeFormatRFC value representing the layout to use for formatting.
//
// Returns:
//   - A string representing the formatted time.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	formatted := t.FormatRFC(TimeFormat20060102T150405) // Returns the formatted time in the format "2023-08-15 13:45:30".
func (t *Timex) FormatRFC(layout TimeFormatRFC) string {
	return t.Time.Format(string(layout))
}

// FormatRFCshort formats the time using the provided layout.
//
// Parameters:
//   - `layout`: A TimeRFC value representing the layout to use for formatting.
//
// Returns:
//   - A string representing the formatted time.
//
// Example:
//
//	t := Timex{Time: time.Now()}
//	formatted := t.FormatRFCshort(TimeRFC01T150405) // Returns the formatted time in the format "13:45:30".
func (t *Timex) FormatRFCshort(layout TimeRFC) string {
	return t.Time.Format(string(layout))
}

// parseWithFormat attempts to parse a given date/time string `s` using a series of predefined formats
// specified in the `TimeFormats` slice of the Timex instance. It tries to parse the string in the provided
// `location` and returns the parsed time value or an error if parsing fails.
//
// The function iterates over each format in `TimeFormats`, attempting to parse the string with each format.
// If the parsing is successful (i.e., no error is returned), it returns the parsed time value. If all formats
// fail to parse the string, it returns an error indicating the failure.
//
// Parameters:
//   - `s`: A string representing the date/time to be parsed.
//   - `location`: A pointer to a time.Location value, which specifies the timezone to use during parsing.
//
// Returns:
//   - `v`: A time.Time value representing the parsed date/time if successful.
//   - `err`: An error value indicating any issues that occurred during parsing; if parsing is successful,
//     this will be nil.
//
// Example:
//
//	t := Timex{TimeFormats: []string{"2006-01-02 15:04:05", "2006-01-02"}}
//	parsedTime, err := t.parseWithFormat("2023-10-25 12:30:00", time.UTC) // Attempts to parse the given string.
//	if err != nil {
//		// Handle parsing error
//	}
//
// Note:
//   - The function will return the first successfully parsed time value and ignore any subsequent formats.
func (t *Timex) parseWithFormat(s string, location *time.Location) (v time.Time, err error) {
	for _, format := range t.timeFormats {
		v, err = time.ParseInLocation(format, s, location)

		if err == nil {
			return
		}
	}
	err = fmt.Errorf("can't parse string as time: %v", s)
	return
}
