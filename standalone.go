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
//
//   - `v`: A time.Time value representing the date from which the beginning of the day is extracted.
//
// Returns:
//
//   - A time.Time value representing the start of the day (00:00:00) for the provided date.
//
// Example:
//
//	now := time.Now()
//	startOfDay := BeginOfDay(now) // This will set the time to midnight of the current day.
func BeginOfDay(v time.Time) time.Time {
	return time.Date(v.Year(), v.Month(), v.Day(), 0, 0, 0, 0, v.Local().Location())
}

// EndOfDay takes a time value `v` and returns a new time.Time object
// representing the end of the day for that date.
//
// The function uses the time.Date method to set the time to the last possible second
// (23:59:59) of the provided day. It maintains the same year, month, and day values
// from the input time `v`. The location (timezone) of the returned time is the same as
// the input time.
//
// Parameters:
//
//   - `v`: A time.Time value representing the date from which the end of the day is extracted.
//
// Returns:
//
//   - A time.Time value representing the end of the day (23:59:59) for the provided date.
//
// Example:
//
//	now := time.Now()
//	endOfDay := EndOfDay(now) // This will set the time to the last second of the current day.
func EndOfDay(v time.Time) time.Time {
	return time.Date(v.Year(), v.Month(), v.Day(), 23, 59, 59, 0, v.Local().Location())
}

// PrevBeginOfDay takes a time value `v` and an integer `day` representing the number of days to go back.
// It returns a new time.Time object representing the beginning of the day for the date `day` days before the given date.
//
// The function subtracts `day` days from the input time `v` using the time.AddDate method and then calls
// BeginOfDay on the resulting time to set the time to midnight (00:00:00) for that earlier date.
//
// Parameters:
//
//   - `v`: A time.Time value representing the reference date.
//
//   - `day`: An integer representing how many days to go back from the reference date.
//
// Returns:
//
//   - A time.Time value representing the start of the day (00:00:00) for the date that is `day` days before `v`.
//
// Example:
//
//	now := time.Now()
//	twoDaysAgoStart := PrevBeginOfDay(now, 2) // This will return the start of the day two days before the current date.
func PrevBeginOfDay(v time.Time, day int) time.Time {
	last := v.AddDate(0, 0, -day)
	return BeginOfDay(last)
}

// PrevEndOfDay takes a time value `v` and an integer `day` representing the number of days to go back.
// It returns a new time.Time object representing the end of the day for the date `day` days before the given date.
//
// The function subtracts `day` days from the input time `v` using the time.AddDate method and then calls
// EndOfDayX on the resulting time to set the time to the last second (23:59:59) of that earlier date.
//
// Parameters:
//
//   - `v`: A time.Time value representing the reference date.
//
//   - `day`: An integer representing how many days to go back from the reference date.
//
// Returns:
//
//   - A time.Time value representing the end of the day (23:59:59) for the date that is `day` days before `v`.
//
// Example:
//
//	now := time.Now()
//	twoDaysAgoEnd := PrevEndOfDay(now, 2) // This will return the end of the day two days before the current date.
func PrevEndOfDay(v time.Time, day int) time.Time {
	last := v.AddDate(0, 0, -day)
	return EndOfDay(last)
}

// SetTimezone takes a time value `v` and a string `tz` representing the target timezone.
// It returns a new time.Time object with the same time as `v` but converted to the specified timezone `tz`.
//
// The function uses time.LoadLocation to load the location based on the timezone string `tz`.
// It then converts the input time `v` to the specified timezone using the time.In method.
// If an error occurs while loading the timezone (for example, if the timezone string is invalid),
// the function returns the current time value in UTC along with the error.
//
// Parameters:
//
//   - `v`: A time.Time value representing the reference time.
//
//   - `tz`: A string representing the IANA timezone name (e.g., "America/New_York", "Europe/London").
//
// Returns:
//
//   - A time.Time value representing the time `v` converted to the specified timezone.
//
//   - An error value, which will be non-nil if the timezone string is invalid.
//
// Example:
//
//	now := time.Now()
//	nyTime, err := SetTimezone(now, "America/New_York") // This will convert the current time to New York's timezone.
func SetTimezone(v time.Time, tz string) (time.Time, error) {
	loc, err := time.LoadLocation(tz)
	now := v.In(loc)
	return now, err
}

// AdjustTimezone takes a time value `v` and a string `tz` representing the target timezone,
// and attempts to adjust `v` to the specified timezone. If the timezone conversion fails (due to an invalid timezone),
// the function returns the original time `v` without any changes.
//
// The function internally calls SetTimezone to convert the time `v` to the specified timezone `tz`.
// If SetTimezone returns an error (e.g., if the timezone string is invalid), the function catches it and returns `v`.
// If no error occurs, it returns the adjusted time in the specified timezone.
//
// Parameters:
//
//   - `v`: A time.Time value representing the reference time.
//
//   - `tz`: A string representing the IANA timezone name (e.g., "America/New_York", "Asia/Tokyo").
//
// Returns:
//
//   - A time.Time value representing the time `v` adjusted to the specified timezone. If an error occurs, it returns the original time.
//
// Example:
//
//	now := time.Now()
//	adjustedTime := AdjustTimezone(now, "Europe/Paris") // This will adjust `now` to Paris' timezone if valid, or return `now` if invalid.
func AdjustTimezone(v time.Time, tz string) time.Time {
	t, err := SetTimezone(v, tz)
	if err != nil {
		return v
	}
	return t
}

// AddSecond takes a time value `v` and an integer `second` representing the number of seconds to add (or subtract if negative).
// It returns a new time.Time object that is adjusted by the specified number of seconds.
//
// The function uses time.Add to add the given number of seconds to `v`. If `second` is 0, the function simply returns the original time `v`.
//
// Parameters:
//
//   - `v`: A time.Time value representing the initial time.
//
//   - `second`: An integer representing the number of seconds to add. If negative, it subtracts the seconds from `v`.
//
// Returns:
//
//   - A time.Time value representing the time `v` adjusted by the specified number of seconds.
//
// Example:
//
// now := time.Now()
//
//	fiveSecondsLater := AddSecond(now, 5)  // This will return the time 5 seconds later.
//	fiveSecondsEarlier := AddSecond(now, -5) // This will return the time 5 seconds earlier.
func AddSecond(v time.Time, second int) time.Time {
	if second == 0 {
		return v
	}
	return v.Add(time.Second * time.Duration(second))
}

// AddMinute takes a time value `v` and an integer `minute` representing the number of minutes to add (or subtract if negative).
// It returns a new time.Time object that is adjusted by the specified number of minutes.
//
// The function uses time.Add to add the given number of minutes to `v`. If `minute` is 0, the function simply returns the original time `v`.
//
// Parameters:
//
//   - `v`: A time.Time value representing the initial time.
//
//   - `minute`: An integer representing the number of minutes to add. If negative, it subtracts the minutes from `v`.
//
// Returns:
//
//   - A time.Time value representing the time `v` adjusted by the specified number of minutes.
//
// Example:
//
//	now := time.Now()
//	tenMinutesLater := AddMinute(now, 10)  // This will return the time 10 minutes later.
//	tenMinutesEarlier := AddMinute(now, -10) // This will return the time 10 minutes earlier.
func AddMinute(v time.Time, minute int) time.Time {
	if minute == 0 {
		return v
	}
	return v.Add(time.Minute * time.Duration(minute))
}

// AddHour takes a time value `v` and an integer `hour` representing the number of hours to add (or subtract if negative).
// It returns a new time.Time object that is adjusted by the specified number of hours.
//
// The function uses time.Add to add the given number of hours to `v`. If `hour` is 0, the function simply returns the original time `v`.
//
// Parameters:
//
//   - `v`: A time.Time value representing the initial time.
//
//   - `hour`: An integer representing the number of hours to add. If negative, it subtracts the hours from `v`.
//
// Returns:
//
//   - A time.Time value representing the time `v` adjusted by the specified number of hours.
//
// Example:
//
//	now := time.Now()
//	threeHoursLater := AddHour(now, 3)  // This will return the time 3 hours later.
//	threeHoursEarlier := AddHour(now, -3) // This will return the time 3 hours earlier.
func AddHour(v time.Time, hour int) time.Time {
	if hour == 0 {
		return v
	}
	return v.Add(time.Hour * time.Duration(hour))
}

// AddDay takes a time value `v` and an integer `day` representing the number of days to add (or subtract if negative).
// It returns a new time.Time object that is adjusted by the specified number of days.
//
// The function uses time.Add to add the given number of days to `v`. Since a day has 24 hours, it multiplies 24 by the number of days
// to convert the days into hours. If `day` is 0, the function simply returns the original time `v`.
//
// Parameters:
//
//   - `v`: A time.Time value representing the initial time
//
//   - `day`: An integer representing the number of days to add. If negative, it subtracts the days from `v`.
//
// Returns:
//
//   - A time.Time value representing the time `v` adjusted by the specified number of days.
//
// Example:
//
//	now := time.Now()
//	fiveDaysLater := AddDay(now, 5)  // This will return the time 5 days later.
//	fiveDaysEarlier := AddDay(now, -5) // This will return the time 5 days earlier.
func AddDay(v time.Time, day int) time.Time {
	if day == 0 {
		return v
	}
	return v.Add(time.Hour * 24 * time.Duration(day))
}

// IsWithinTolerance checks if the provided time `v` is within a one-minute tolerance window around the current time.
//
// The function calculates the difference between the provided time `v` and the current time using time.Sub.
// It then checks if the difference is within a range of plus or minus one minute (`tolerance`). If the difference falls
// within this range, the function returns true, indicating that `v` is within the tolerance window.
// Otherwise, it returns false.
//
// Parameters:
//
//   - `v`: A time.Time value representing the time to check.
//
// Returns:
//
// - A boolean value:
//
//   - true if `v` is within one minute (before or after) the current time;
//
//   - false if `v` is outside this one-minute window.
//
// Example:
//
//	now := time.Now()
//	checkTime := now.Add(time.Second * 30)
//	isOnTime := IsWithinTolerance(checkTime) // This will return true since checkTime is within 1 minute of now.
func IsWithinTolerance(v time.Time) bool {
	target := time.Now()
	tolerance := time.Minute
	diff := v.Sub(target)
	return diff >= -tolerance && diff <= tolerance
}

// IsLeapYear determines if the specified year is a leap year.
//
// A year is considered a leap year if:
//   - It is divisible by 4; and
//   - It is not divisible by 100, unless it is also divisible by 400.
//
// Parameters:
//
//   - `year`: An integer representing the year to check.
//
// Returns:
// - A boolean value:
//   - true if the specified year is a leap year;
//   - false otherwise.
//
// Example:
//
//	leapYear := IsLeapYear(2020) // This will return true since 2020 is a leap year.
//	nonLeapYear := IsLeapYear(2021) // This will return false since 2021 is not a leap year.
func IsLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// IsLeapYearN checks if the year of the provided time value `v` is a leap year.
//
// The function retrieves the year from the time.Time object using the Year() method
// and calls IsLeapYear to determine if that year is a leap year.
//
// Parameters:
//
//   - `v`: A time.Time value representing the date from which the year is extracted.
//
// Returns:
// - A boolean value:
//   - true if the year from the date `v` is a leap year;
//   - false otherwise.
//
// Example:
//
//	now := time.Now()
//	isLeap := IsLeapYearN(now) // This will return true if the current year is a leap year.
func IsLeapYearN(v time.Time) bool {
	return IsLeapYear(v.Year())
}

// GetWeekdaysInRange returns a slice of time.Time objects representing all weekdays (Monday to Friday)
// between the specified start and end dates, inclusive.
//
// The function iterates through each date from `start` to `end`, checking if each date is a weekday.
// It excludes Saturdays and Sundays. It also handles leap years correctly by ensuring that February 29
// is included only in leap years. If the year is not a leap year, it checks if the day is valid for the month.
//
// Parameters:
//
//   - `start`: A time.Time value representing the start date of the range.
//
//   - `end`: A time.Time value representing the end date of the range.
//
// Returns:
//
//   - A slice of time.Time values representing all weekdays between `start` and `end`, inclusive.
//
// Example:
//
//	start := time.Date(2023, time.March, 1, 0, 0, 0, 0, time.UTC)
//	end := time.Date(2023, time.March, 10, 0, 0, 0, 0, time.UTC)
//	weekdays := GetWeekdaysInRange(start, end) // This will return all weekdays between March 1 and March 10, 2023.
func GetWeekdaysInRange(start time.Time, end time.Time) []time.Time {
	var weekdays []time.Time
	for current := start; current.Before(end) || current.Equal(end); current = current.AddDate(0, 0, 1) {
		d := current.Weekday()
		if d != time.Sunday && d != time.Saturday {
			y := current.Year()
			if IsLeapYear(y) && current.Month() == time.February && current.Day() == 29 {
				weekdays = append(weekdays, current)
			} else if !IsLeapYear(y) && current.Day() <= time.Date(y, time.December, 31, 0, 0, 0, 0, time.UTC).Day() {
				weekdays = append(weekdays, current)
			}
		}
	}
	return weekdays
}

// SinceHour calculates the number of hours that have passed since the provided time value `v`.
//
// The function computes the time difference between the current time and `v` using time.Since().
// The resulting duration is then converted into hours using the Hours() method.
//
// Parameters:
//
//   - `v`: A time.Time value representing the starting time to calculate the elapsed hours from.
//
// Returns:
//
//   - A float64 value representing the number of hours that have passed since the time `v`.
//
// Example:
//
//	start := time.Date(2023, time.March, 15, 8, 0, 0, 0, time.UTC)
//	elapsedHours := SinceHour(start) // This will return the hours passed since March 15, 2023, 8:00 AM.
func SinceHour(v time.Time) float64 {
	duration := time.Since(v)
	hours := duration.Hours()
	return hours
}

// SinceMinute calculates the number of minutes that have passed since the provided time value `v`.
//
// The function computes the time difference between the current time and `v` using time.Since().
// The resulting duration is then converted into minutes using the Minutes() method.
//
// Parameters:
//
//   - `v`: A time.Time value representing the starting time to calculate the elapsed minutes from.
//
// Returns:
//
//   - A float64 value representing the number of minutes that have passed since the time `v`.
//
// Example:
//
//	start := time.Date(2023, time.March, 15, 8, 0, 0, 0, time.UTC)
//	elapsedMinutes := SinceMinute(start) // This will return the minutes passed since March 15, 2023, 8:00 AM.
func SinceMinute(v time.Time) float64 {
	duration := time.Since(v)
	minutes := duration.Minutes()
	return minutes
}

// SinceSecond calculates the number of seconds that have passed since the provided time value `v`.
//
// The function computes the time difference between the current time and `v` using time.Since().
// The resulting duration is then converted into seconds using the Seconds() method.
//
// Parameters:
//
//   - `v`: A time.Time value representing the starting time to calculate the elapsed seconds from.
//
// Returns:
//
//   - A float64 value representing the number of seconds that have passed since the time `v`.
//
// Example:
//
//	start := time.Date(2023, time.March, 15, 8, 0, 0, 0, time.UTC)
//	elapsedSeconds := SinceSecond(start) // This will return the seconds passed since March 15, 2023, 8:00 AM.
func SinceSecond(v time.Time) float64 {
	duration := time.Since(v)
	seconds := duration.Seconds()
	return seconds
}

// FormatTimex converts a given time.Time value into a slice of integers representing various time components.
//
// The function extracts the hour, minute, second, nanosecond, day, month, and year from the provided
// time.Time value. It then returns these components in a specific order as a slice of integers.
//
// The order of the returned slice is as follows:
//   - [0]: Nanosecond
//   - [1]: Second
//   - [2]: Minute
//   - [3]: Hour
//   - [4]: Day
//   - [5]: Month (as an integer, where January is 1)
//   - [6]: Year
//
// Parameters:
//
//   - `t`: A time.Time value representing the time to format.
//
// Returns:
//
//   - A slice of integers containing the nanosecond, second, minute, hour, day, month, and year
//     components of the provided time.Time value.
//
// Example:
//
//	t := time.Date(2023, time.March, 15, 8, 0, 0, 0, time.UTC)
//	formatted := FormatTimex(t) // This will return a slice with the components [0, 0, 0, 8, 15, 3, 2023].
func FormatTimex(t time.Time) []int {
	hour, min, sec := t.Clock()
	year, month, day := t.Date()
	return []int{t.Nanosecond(), sec, min, hour, day, int(month), year}
}
