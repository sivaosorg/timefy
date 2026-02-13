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
	return time.Date(v.Year(), v.Month(), v.Day(), 0, 0, 0, 0, v.Location())
}

// FEndOfDay takes a time value `v` and returns a new time.Time object
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
//	endOfDay := FEndOfDay(now) // This will set the time to the last second of the current day.
func FEndOfDay(v time.Time) time.Time {
	return time.Date(v.Year(), v.Month(), v.Day(), 23, 59, 59, 0, v.Location())
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
	return FEndOfDay(last)
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
	if err != nil {
		return v, err
	}
	now := v.In(loc)
	return now, nil
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

// BeginningOfMinute returns the current time rounded down to the beginning of the current minute.
// It utilizes the With() function to achieve this. The resulting time will have seconds and nanoseconds set to zero.
//
// Returns:
//   - A time.Time value representing the start of the current minute.
//
// Example:
//
//	beginning := BeginningOfMinute() // This will return the current time set to the start of the minute (e.g., 12:30:00).
func BeginningOfMinute() time.Time {
	return With(time.Now()).BeginningOfMinute()
}

// BeginningOfHour returns the current time rounded down to the beginning of the current hour.
// This function resets the minute, second, and nanosecond components of the time to zero, providing
// a time value that represents the exact start of the hour.
//
// It utilizes the With() function to obtain the current time and then applies the BeginningOfHour()
// method to round it down.
//
// Returns:
//   - A time.Time value representing the start of the current hour.
//
// Example:
//
//	beginning := BeginningOfHour() // This will return the current time set to the start of the hour (e.g., 12:00:00).
func BeginningOfHour() time.Time {
	return With(time.Now()).BeginningOfHour()
}

// BeginningOfDay returns the current time rounded down to the beginning of the current day.
// This function resets the hour, minute, second, and nanosecond components of the time to zero,
// providing a time value that represents the exact start of the day (midnight).
//
// It utilizes the With() function to obtain the current time and then applies the BeginningOfDay()
// method to achieve this rounding.
//
// Returns:
//   - A time.Time value representing the start of the current day (e.g., 00:00:00).
//
// Example:
//
//	beginning := BeginningOfDay() // This will return the current time set to the start of the day (e.g., 2023-10-25 00:00:00).
func BeginningOfDay() time.Time {
	return With(time.Now()).BeginningOfDay()
}

// BeginningOfWeek returns the current time rounded down to the beginning of the current week.
// This function resets the hour, minute, second, and nanosecond components of the time to zero,
// providing a time value that represents the exact start of the week (usually Sunday or Monday
// depending on the locale).
//
// It utilizes the With() function to obtain the current time and then applies the BeginningOfWeek()
// method to achieve this rounding.
//
// Returns:
//   - A time.Time value representing the start of the current week (e.g., 00:00:00 on the first day of the week).
//
// Example:
//
//	beginning := BeginningOfWeek() // This will return the current time set to the start of the week (e.g., 2023-10-22 00:00:00 if Sunday is the start of the week).
func BeginningOfWeek() time.Time {
	return With(time.Now()).BeginningOfWeek()
}

// BeginningOfMonth returns the current time rounded down to the beginning of the current month.
// This function resets the day, hour, minute, second, and nanosecond components of the time to zero,
// providing a time value that represents the exact start of the month (the first day at midnight).
//
// It utilizes the With() function to obtain the current time and then applies the BeginningOfMonth()
// method to achieve this rounding.
//
// Returns:
//   - A time.Time value representing the start of the current month (e.g., 2023-10-01 00:00:00).
//
// Example:
//
//	beginning := BeginningOfMonth() // This will return the current time set to the start of the month (e.g., 2023-10-01 00:00:00).
func BeginningOfMonth() time.Time {
	return With(time.Now()).BeginningOfMonth()
}

// BeginningOfQuarter returns the current time rounded down to the beginning of the current quarter.
// A quarter is a three-month period within the year, specifically:
//   - Q1: January to March
//   - Q2: April to June
//   - Q3: July to September
//   - Q4: October to December
//
// This function resets the day, hour, minute, second, and nanosecond components of the time to zero,
// providing a time value that represents the exact start of the current quarter (the first day of the
// quarter at midnight).
//
// It utilizes the With() function to obtain the current time and then applies the BeginningOfQuarter()
// method to achieve this rounding.
//
// Returns:
//   - A time.Time value representing the start of the current quarter (e.g., 2023-10-01 00:00:00
//     if the current date is in Q4).
//
// Example:
//
//	beginning := BeginningOfQuarter() // This will return the current time set to the start of the current quarter (e.g., 2023-10-01 00:00:00 if it's the fourth quarter).
func BeginningOfQuarter() time.Time {
	return With(time.Now()).BeginningOfQuarter()
}

// BeginningOfYear returns the current time rounded down to the beginning of the current year.
// This function resets the month, day, hour, minute, second, and nanosecond components of the time to zero,
// providing a time value that represents the exact start of the year (January 1st at midnight).
//
// It utilizes the With() function to obtain the current time and then applies the BeginningOfYear()
// method to achieve this rounding.
//
// Returns:
//   - A time.Time value representing the start of the current year (e.g., 2023-01-01 00:00:00).
//
// Example:
//
//	beginning := BeginningOfYear() // This will return the current time set to the start of the year (e.g., 2023-01-01 00:00:00).
func BeginningOfYear() time.Time {
	return With(time.Now()).BeginningOfYear()
}

// EndOfMinute returns the current time rounded up to the end of the current minute.
// This function resets the second and nanosecond components of the time to zero and then adds one minute,
// providing a time value that represents the last moment of the current minute (59 seconds and 999999999 nanoseconds).
//
// It utilizes the With() function to obtain the current time and then applies the EndOfMinute()
// method to achieve this rounding.
//
// Returns:
//   - A time.Time value representing the end of the current minute (e.g., 12:30:59.999999999).
//
// Example:
//
//	end := EndOfMinute() // This will return the current time set to the end of the minute (e.g., 12:30:59.999999999).
func EndOfMinute() time.Time {
	return With(time.Now()).EndOfMinute()
}

// EndOfHour returns the current time rounded up to the end of the current hour.
// This function resets the minute, second, and nanosecond components of the time to zero and then adds one hour,
// providing a time value that represents the last moment of the current hour (59 minutes, 59 seconds, and 999999999 nanoseconds).
//
// It utilizes the With() function to obtain the current time and then applies the EndOfHour()
// method to achieve this rounding.
//
// Returns:
//   - A time.Time value representing the end of the current hour (e.g., 12:59:59.999999999).
//
// Example:
//
//	end := EndOfHour() // This will return the current time set to the end of the hour (e.g., 12:59:59.999999999).
func EndOfHour() time.Time {
	return With(time.Now()).EndOfHour()
}

// EndOfDay returns the current time rounded up to the end of the current day.
// This function resets the hour, minute, second, and nanosecond components of the time to zero and then adds one day,
// providing a time value that represents the last moment of the current day (23 hours, 59 minutes, 59 seconds, and 999999999 nanoseconds).
//
// It utilizes the With() function to obtain the current time and then applies the EndOfDay()
// method to achieve this rounding.
//
// Returns:
//   - A time.Time value representing the end of the current day (e.g., 2023-10-25 23:59:59.999999999).
//
// Example:
//
//	end := EndOfDay() // This will return the current time set to the end of the day (e.g., 2023-10-25 23:59:59.999999999).
func EndOfDay() time.Time {
	return With(time.Now()).EndOfDay()
}

// EndOfWeek returns the current time rounded up to the end of the current week.
// This function resets the hour, minute, second, and nanosecond components of the time to zero and then adds one week,
// providing a time value that represents the last moment of the current week (e.g., 23 hours, 59 minutes, 59 seconds, and 999999999 nanoseconds)
// on the last day of the week (usually Saturday or Sunday, depending on the locale).
//
// It utilizes the With() function to obtain the current time and then applies the EndOfWeek()
// method to achieve this rounding.
//
// Returns:
//   - A time.Time value representing the end of the current week (e.g., 2023-10-29 23:59:59.999999999 if Sunday is the last day of the week).
//
// Example:
//
//	end := EndOfWeek() // This will return the current time set to the end of the week (e.g., 2023-10-29 23:59:59.999999999).
func EndOfWeek() time.Time {
	return With(time.Now()).EndOfWeek()
}

// EndOfMonth returns the current time rounded up to the end of the current month.
// This function resets the day, hour, minute, second, and nanosecond components of the time to zero
// and then adds one month, providing a time value that represents the last moment of the current month
// (e.g., 23 hours, 59 minutes, 59 seconds, and 999999999 nanoseconds) on the last day of the month.
//
// It utilizes the With() function to obtain the current time and then applies the EndOfMonth()
// method to achieve this rounding.
//
// Returns:
//   - A time.Time value representing the end of the current month (e.g., 2023-10-31 23:59:59.999999999).
//
// Example:
//
//	end := EndOfMonth() // This will return the current time set to the end of the month (e.g., 2023-10-31 23:59:59.999999999).
func EndOfMonth() time.Time {
	return With(time.Now()).EndOfMonth()
}

// EndOfQuarter returns the current time rounded up to the end of the current quarter.
// A quarter is defined as a three-month period within the year, specifically:
//   - Q1: January to March
//   - Q2: April to June
//   - Q3: July to September
//   - Q4: October to December
//
// This function resets the day, hour, minute, second, and nanosecond components of the time to zero
// and then adds the number of months needed to reach the first day of the next quarter,
// providing a time value that represents the last moment of the current quarter
// (e.g., 23 hours, 59 minutes, 59 seconds, and 999999999 nanoseconds) on the last day of the quarter.
//
// It utilizes the With() function to obtain the current time and then applies the EndOfQuarter()
// method to achieve this rounding.
//
// Returns:
//   - A time.Time value representing the end of the current quarter (e.g., 2023-03-31 23:59:59.999999999
//     if the current date is in Q1).
//
// Example:
//
//	end := EndOfQuarter() // This will return the current time set to the end of the current quarter (e.g., 2023-12-31 23:59:59.999999999 if it's the fourth quarter).
func EndOfQuarter() time.Time {
	return With(time.Now()).EndOfQuarter()
}

// EndOfYear returns the current time rounded up to the end of the current year.
// This function resets the month, day, hour, minute, second, and nanosecond components of the time to zero
// and then adds one year, providing a time value that represents the last moment of the current year
// (e.g., 23 hours, 59 minutes, 59 seconds, and 999999999 nanoseconds) on December 31st.
//
// It utilizes the With() function to obtain the current time and then applies the EndOfYear()
// method to achieve this rounding.
//
// Returns:
//   - A time.Time value representing the end of the current year (e.g., 2023-12-31 23:59:59.999999999).
//
// Example:
//
//	end := EndOfYear() // This will return the current time set to the end of the current year (e.g., 2023-12-31 23:59:59.999999999).
func EndOfYear() time.Time {
	return With(time.Now()).EndOfYear()
}

// Monday returns the date and time of the most recent or upcoming Monday relative to the current time.
// This function can take an optional string parameter to specify the desired format for the output,
// but it defaults to the standard representation of time if no arguments are provided.
//
// It utilizes the With() function to obtain the current time and then applies the Monday()
// method to determine the appropriate Monday date and time.
//
// Returns:
//   - A time.Time value representing the most recent or upcoming Monday based on the current date and time.
//     The output can vary depending on the optional format parameter.
//
// Example:
//
//	monday := Monday() // This will return the date and time for the next upcoming Monday (e.g., 2023-10-30 00:00:00).
//	mondayFormatted := Monday("2006-01-02") // This will return the next Monday formatted as "YYYY-MM-DD".
func Monday(s ...string) time.Time {
	return With(time.Now()).Monday(s...)
}

// Tuesday returns the date and time of the most recent or upcoming Tuesday relative to the current time.
// This function can take an optional string parameter to specify the desired format for the output,
// but it defaults to the standard representation of time if no arguments are provided.
//
// It utilizes the With() function to obtain the current time and then applies the Tuesday()
// method to determine the appropriate Tuesday date and time.
//
// Returns:
//   - A time.Time value representing the most recent or upcoming Tuesday based on the current date and time.
//     The output can vary depending on the optional format parameter.
//
// Example:
//
//	tuesday := Tuesday() // This will return the date and time for the next upcoming Tuesday (e.g., 2023-10-31 00:00:00).
//	tuesdayFormatted := Tuesday("2006-01-02") // This will return the next Tuesday formatted as "YYYY-MM-DD".
func Tuesday(s ...string) time.Time {
	return With(time.Now()).Tuesday(s...)
}

// Wednesday returns the date and time of the most recent or upcoming Wednesday relative to the current time.
// This function can take an optional string parameter to specify the desired format for the output,
// but it defaults to the standard representation of time if no arguments are provided.
//
// It utilizes the With() function to obtain the current time and then applies the Wednesday()
// method to determine the appropriate Wednesday date and time.
//
// Returns:
//   - A time.Time value representing the most recent or upcoming Wednesday based on the current date and time.
//     The output can vary depending on the optional format parameter.
//
// Example:
//
//	wednesday := Wednesday() // This will return the date and time for the next upcoming Wednesday (e.g., 2023-11-01 00:00:00).
//	wednesdayFormatted := Wednesday("2006-01-02") // This will return the next Wednesday formatted as "YYYY-MM-DD".
func Wednesday(s ...string) time.Time {
	return With(time.Now()).Wednesday(s...)
}

// Thursday returns the date and time of the most recent or upcoming Thursday relative to the current time.
// This function can take an optional string parameter to specify the desired format for the output,
// but it defaults to the standard representation of time if no arguments are provided.
//
// It utilizes the With() function to obtain the current time and then applies the Thursday()
// method to determine the appropriate Thursday date and time.
//
// Returns:
//   - A time.Time value representing the most recent or upcoming Thursday based on the current date and time.
//     The output can vary depending on the optional format parameter.
//
// Example:
//
//	thursday := Thursday() // This will return the date and time for the next upcoming Thursday (e.g., 2023-11-02 00:00:00).
//	thursdayFormatted := Thursday("2006-01-02") // This will return the next Thursday formatted as "YYYY-MM-DD".
func Thursday(s ...string) time.Time {
	return With(time.Now()).Thursday(s...)
}

// Friday returns the date and time of the most recent or upcoming Friday relative to the current time.
// This function can take an optional string parameter to specify the desired format for the output,
// but it defaults to the standard representation of time if no arguments are provided.
//
// It utilizes the With() function to obtain the current time and then applies the Friday()
// method to determine the appropriate Friday date and time.
//
// Returns:
//   - A time.Time value representing the most recent or upcoming Friday based on the current date and time.
//     The output can vary depending on the optional format parameter.
//
// Example:
//
//	friday := Friday() // This will return the date and time for the next upcoming Friday (e.g., 2023-11-03 00:00:00).
//	fridayFormatted := Friday("2006-01-02") // This will return the next Friday formatted as "YYYY-MM-DD".
func Friday(s ...string) time.Time {
	return With(time.Now()).Friday(s...)
}

// Saturday returns the date and time of the most recent or upcoming Saturday relative to the current time.
// This function can take an optional string parameter to specify the desired format for the output,
// but it defaults to the standard representation of time if no arguments are provided.
//
// It utilizes the With() function to obtain the current time and then applies the Saturday()
// method to determine the appropriate Saturday date and time.
//
// Returns:
//   - A time.Time value representing the most recent or upcoming Saturday based on the current date and time.
//     The output can vary depending on the optional format parameter.
//
// Example:
//
//	friday := Friday() // This will return the date and time for the next upcoming Friday (e.g., 2023-11-03 00:00:00).
//	fridayFormatted := Friday("2006-01-02") // This will return the next Friday formatted as "YYYY-MM-DD".
func Saturday(s ...string) time.Time {
	return With(time.Now()).Saturday(s...)
}

// Sunday returns the date and time of the most recent or upcoming Sunday relative to the current time.
// This function can take an optional string parameter to specify the desired format for the output,
// but it defaults to the standard representation of time if no arguments are provided.
//
// It utilizes the With() function to obtain the current time and then applies the Sunday()
// method to determine the appropriate Sunday date and time.
//
// Returns:
//   - A time.Time value representing the most recent or upcoming Sunday based on the current date and time.
//     The output can vary depending on the optional format parameter.
//
// Example:
//
//	sunday := Sunday() // This will return the date and time for the next upcoming Sunday (e.g., 2023-10-29 00:00:00).
//	sundayFormatted := Sunday("2006-01-02") // This will return the next Sunday formatted as "YYYY-MM-DD".
func Sunday(s ...string) time.Time {
	return With(time.Now()).Sunday(s...)
}

// EndOfSunday returns the date and time representing the end of the most recent or upcoming Sunday
// relative to the current time. This function resets the time to 23 hours, 59 minutes, 59 seconds,
// and 999999999 nanoseconds, providing a time value that represents the last moment of Sunday.
//
// It utilizes the With() function to obtain the current time and then applies the EndOfSunday()
// method to achieve this rounding.
//
// Returns:
//   - A time.Time value representing the end of the most recent or upcoming Sunday based on the
//     current date and time (e.g., 2023-10-29 23:59:59.999999999 if the current date is within that week).
//
// Example:
//
//	end := EndOfSunday() // This will return the date and time set to the end of the next Sunday (e.g., 2023-10-29 23:59:59.999999999).
func EndOfSunday() time.Time {
	return With(time.Now()).EndOfSunday()
}

// Quarter returns the current quarter of the year based on the current date and time.
// A quarter is defined as a three-month period within the year, specifically:
//   - Q1: January to March
//   - Q2: April to June
//   - Q3: July to September
//   - Q4: October to December
//
// This function utilizes the With() function to obtain the current time and then applies the
// Quarter() method to determine the current quarter.
//
// Returns:
//   - A uint value representing the current quarter of the year (1, 2, 3, or 4).
//
// Example:
//
//	quarter := Quarter() // This will return the current quarter (e.g., 4 for October).
func Quarter() uint {
	return With(time.Now()).Quarter()
}

// Parse takes a variable number of string inputs and attempts to parse them into a time.Time value.
// This function uses the With() function to obtain the current time as a reference point and then
// applies the Parse() method to interpret the provided string(s) as time.
//
// The function can handle multiple formats and will return the first successfully parsed time value
// along with any potential error encountered during parsing.
//
// Returns:
//   - A time.Time value representing the parsed time if successful.
//   - An error indicating any issues encountered during the parsing process, or nil if parsing was successful.
//
// Example:
//
//	timeValue, err := Parse("2023-10-25") // This will return the parsed time if the input string is in a valid format.
//	if err != nil {
//		// Handle the parsing error
//	}
func Parse(s ...string) (time.Time, error) {
	return With(time.Now()).Parse(s...)
}

// ParseInLocation takes a variable number of string inputs and attempts to parse them into a time.Time value
// based on a specified time zone location. This function utilizes the With() function to obtain the current
// time in the provided location as a reference point and then applies the Parse() method to interpret
// the provided string(s) as time.
//
// The function can handle multiple formats and will return the first successfully parsed time value
// in the specified location, along with any potential error encountered during parsing.
//
// Parameters:
//   - loc: A pointer to a time.Location struct that specifies the desired time zone for parsing.
//   - s: A variadic parameter that accepts one or more strings to be parsed into a time.Time value.
//
// Returns:
//   - A time.Time value representing the parsed time in the specified location if successful.
//   - An error indicating any issues encountered during the parsing process, or nil if parsing was successful.
//
// Example:
//
//	timeValue, err := ParseInLocation(time.UTC, "2023-10-25") // This will return the parsed time in UTC if the input string is in a valid format.
//	if err != nil {
//		// Handle the parsing error
//	}
func ParseInLocation(loc *time.Location, s ...string) (time.Time, error) {
	return With(time.Now().In(loc)).Parse(s...)
}

// MustParse takes a variable number of string inputs and attempts to parse them into a time.Time value.
// This function uses the With() function to obtain the current time as a reference point and then
// applies the MustParse() method to interpret the provided string(s) as time.
//
// Unlike the Parse function, MustParse will panic if the parsing fails, making it suitable for scenarios
// where a valid time is expected and errors are not anticipated. If the input strings are in valid formats,
// it returns the corresponding time.Time value.
//
// Parameters:
//   - s: A variadic parameter that accepts one or more strings to be parsed into a time.Time value.
//
// Returns:
//   - A time.Time value representing the parsed time if successful.
//   - This function does not return an error; instead, it will panic if parsing fails.
//
// Example:
//
//	timeValue := MustParse("2023-10-25") // This will return the parsed time if the input string is in a valid format.
//	// If the input is invalid, it will cause a panic.
func MustParse(s ...string) time.Time {
	return With(time.Now()).MustParse(s...)
}

// MustParseInLocation takes a variable number of string inputs and attempts to parse them into a time.Time value
// based on a specified time zone location. This function utilizes the With() function to obtain the current
// time in the provided location as a reference point and then applies the MustParse() method to interpret
// the provided string(s) as time.
//
// Similar to MustParse, this function will panic if the parsing fails, making it suitable for scenarios
// where a valid time is expected and errors are not anticipated. If the input strings are in valid formats,
// it returns the corresponding time.Time value adjusted to the specified location.
//
// Parameters:
//   - loc: A pointer to a time.Location struct that specifies the desired time zone for parsing.
//   - s: A variadic parameter that accepts one or more strings to be parsed into a time.Time value.
//
// Returns:
//   - A time.Time value representing the parsed time in the specified location if successful.
//   - This function does not return an error; instead, it will panic if parsing fails.
//
// Example:
//
//	timeValue := MustParseInLocation(time.UTC, "2023-10-25") // This will return the parsed time in UTC if the input string is in a valid format.
//	// If the input is invalid, it will cause a panic.
func MustParseInLocation(loc *time.Location, s ...string) time.Time {
	return With(time.Now().In(loc)).MustParse(s...)
}

// Between takes two string inputs representing time values and checks if the current time falls
// within the range defined by those two times. This function utilizes the With() function to obtain
// the current time as a reference point and then applies the Between() method to evaluate the range.
//
// The function assumes that the provided time strings are in a valid format. If the current time
// is greater than or equal to time1 and less than or equal to time2, it returns true; otherwise, it returns false.
//
// Parameters:
//   - time1: A string representing the start of the time range.
//   - time2: A string representing the end of the time range.
//
// Returns:
//   - A boolean value indicating whether the current time is within the specified range (inclusive).
//
// Example:
//
//	isWithin := Between("2023-10-20", "2023-10-30") // This will return true if the current time is between these two dates.
//	isWithin := Between("2023-10-25", "2023-10-26") // This will return true if the current date is exactly 2023-10-25.
func Between(time1, time2 string) bool {
	return With(time.Now()).Between(time1, time2)
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
//	formatted := FormatRFC(TimeFormat20060102T150405) // Returns the formatted time in the format "2023-08-15 13:45:30".
func FormatRFC(layout TimeFormatRFC) string {
	return With(time.Now()).FormatRFC(layout)
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
//	formatted := FormatRFCshort(TimeRFC01T150405) // Returns the formatted time in the format "13:45:30".
func FormatRFCshort(layout TimeRFC) string {
	return With(time.Now()).FormatRFCshort(layout)
}

// FormatRFC formats the time using the provided layout.
//
// Parameters:
//   - `v`: A time.Time value representing the time to format.
//   - `layout`: A TimeFormatRFC value representing the layout to use for formatting.
//
// Returns:
//   - A string representing the formatted time.
//
// Example:
//
//	formatted := FormatRFC(v, TimeFormat20060102T150405) // Returns the formatted time in the format "2023-08-15 13:45:30".
func FFormatRFC(v time.Time, layout TimeFormatRFC) string {
	return With(v).FormatRFC(layout)
}

// FormatRFCshort formats the time using the provided layout.
//
// Parameters:
//   - `v`: A time.Time value representing the time to format.
//   - `layout`: A TimeRFC value representing the layout to use for formatting.
//
// Returns:
//   - A string representing the formatted time.
//
// Example:
//
//	formatted := FormatRFCshort(v, TimeRFC01T150405) // Returns the formatted time in the format "13:45:30".
func FFormatRFCshort(v time.Time, layout TimeRFC) string {
	return With(v).FormatRFCshort(layout)
}

// DefaultFormatRFC formats the time using the default layout.
//
// Parameters:
//   - `v`: A time.Time value representing the time to format.
//
// Returns:
//   - A string representing the formatted time.
//
// Example:
//
//	v := time.Date(2023, time.August, 15, 13, 45, 30, 0, time.UTC)
//	formatted := DefaultFormatRFC(v) // Returns the formatted time in the default format.
func FDefaultFormatRFC(v time.Time) string {
	return With(v).DefaultFormatRFC()
}

// DefaultFormatRFC formats the current time using the default layout.
//
// Returns:
//   - A string representing the formatted time.
//
// Example:
//
//	formatted := DefaultFormatRFC() // Returns the formatted current time in the default format.
func DefaultFormatRFC() string {
	return FDefaultFormatRFC(time.Now())
}
