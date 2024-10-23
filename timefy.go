package timefy

import "time"

// BeginOfDay takes a time value `v` and returns a new time.Time object
// representing the beginning of the day for that date.
//
// The function uses the time.Date method to set the time to midnight (00:00:00)
// while keeping the year, month, and day values from the input time `v`. The
// location (timezone) of the returned time is the same as the input time.
//
// Parameters:
// - `v`: A time.Time value representing the date from which the beginning of the day is extracted.
//
// Returns:
// - A time.Time value representing the start of the day (00:00:00) for the provided date.
//
// Example:
//
// now := time.Now()
// startOfDay := BeginOfDay(now) // This will set the time to midnight of the current day.
func BeginOfDay(v time.Time) time.Time {
	return time.Date(v.Year(), v.Month(), v.Day(), 0, 0, 0, 0, v.Local().Location())
}

// EndOfDayX takes a time value `v` and returns a new time.Time object
// representing the end of the day for that date.
//
// The function uses the time.Date method to set the time to the last possible second
// (23:59:59) of the provided day. It maintains the same year, month, and day values
// from the input time `v`. The location (timezone) of the returned time is the same as
// the input time.
//
// Parameters:
// - `v`: A time.Time value representing the date from which the end of the day is extracted.
//
// Returns:
// - A time.Time value representing the end of the day (23:59:59) for the provided date.
//
// Example:
//
// now := time.Now()
// endOfDay := EndOfDayX(now) // This will set the time to the last second of the current day.
func EndOfDayX(v time.Time) time.Time {
	return time.Date(v.Year(), v.Month(), v.Day(), 23, 59, 59, 0, v.Local().Location())
}
