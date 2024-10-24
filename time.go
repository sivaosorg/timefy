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
// t := time.Now()
// timex := With(t) // This wraps the current time into a Timex object with the default configuration.
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
// t := time.Now()
// timex := New(t) // This creates a new Timex object for the current time.
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
// customConfig := &Config{WeekStartDay: time.Monday}
//
// t := time.Now()
//
// timex := customConfig.With(t) // This creates a Timex object with the custom configuration.
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
// config := &Config{TimeLocation: time.UTC}
//
// parsedTime, err := config.Parse("2023-10-24T12:00:00Z") // Parses the string using UTC time location.
//
// If no `TimeLocation` is set, the function uses the local time zone:
//
// config := &Config{}
//
// parsedTime, err := config.Parse("2023-10-24T12:00:00") // Parses using the local time zone.
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
// config := &Config{TimeLocation: time.UTC}
//
// parsedTime := config.MustParse("2023-10-24T12:00:00Z") // Parses the string using UTC time location.
//
// If no `TimeLocation` is set, the function uses the local time zone:
//
// config := &Config{}
//
// parsedTime := config.MustParse("2023-10-24T12:00:00") // Parses using the local time zone, panicking on failure.
func (c *Config) MustParse(s ...string) time.Time {
	if c.TimeLocation == nil {
		return c.With(time.Now()).MustParse(s...)
	} else {
		return c.With(time.Now().In(c.TimeLocation)).MustParse(s...)
	}
}
