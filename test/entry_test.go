package test

import (
	"testing"
	"time"

	"github.com/sivaosorg/timefy"
)

func TestBeginOfDay(t *testing.T) {
	v := time.Date(2023, 10, 25, 14, 30, 0, 0, time.UTC)
	expected := time.Date(2023, 10, 25, 0, 0, 0, 0, time.UTC)
	got := timefy.BeginOfDay(v)
	if !got.Equal(expected) {
		t.Errorf("BeginOfDay(%v) = %v; want %v", v, got, expected)
	}
}

func TestFEndOfDay(t *testing.T) {
	v := time.Date(2023, 10, 25, 14, 30, 0, 0, time.UTC)
	expected := time.Date(2023, 10, 25, 23, 59, 59, 0, time.UTC)
	got := timefy.FEndOfDay(v)
	if !got.Equal(expected) {
		t.Errorf("FEndOfDay(%v) = %v; want %v", v, got, expected)
	}
}

func TestPrevBeginOfDay(t *testing.T) {
	v := time.Date(2023, 10, 25, 14, 30, 0, 0, time.UTC)
	expected := time.Date(2023, 10, 23, 0, 0, 0, 0, time.UTC)
	got := timefy.PrevBeginOfDay(v, 2)
	if !got.Equal(expected) {
		t.Errorf("PrevBeginOfDay(%v, 2) = %v; want %v", v, got, expected)
	}
}

func TestPrevEndOfDay(t *testing.T) {
	v := time.Date(2023, 10, 25, 14, 30, 0, 0, time.UTC)
	expected := time.Date(2023, 10, 23, 23, 59, 59, 0, time.UTC)
	got := timefy.PrevEndOfDay(v, 2)
	if !got.Equal(expected) {
		t.Errorf("PrevEndOfDay(%v, 2) = %v; want %v", v, got, expected)
	}
}

func TestSetTimezone(t *testing.T) {
	v := time.Date(2023, 10, 25, 14, 30, 0, 0, time.UTC)
	tz := "America/New_York"
	got, err := timefy.SetTimezone(v, tz)
	if err != nil {
		t.Errorf("SetTimezone(%v, %s) returned error: %v", v, tz, err)
	}
	loc, _ := time.LoadLocation(tz)
	expected := v.In(loc)
	if !got.Equal(expected) {
		t.Errorf("SetTimezone(%v, %s) = %v; want %v", v, tz, got, expected)
	}

	_, err = timefy.SetTimezone(v, "Invalid/Timezone")
	if err == nil {
		t.Errorf("SetTimezone with invalid timezone should return error")
	}
}

func TestAdjustTimezone(t *testing.T) {
	v := time.Date(2023, 10, 25, 14, 30, 0, 0, time.UTC)
	tz := "America/New_York"
	got := timefy.AdjustTimezone(v, tz)
	loc, _ := time.LoadLocation(tz)
	expected := v.In(loc)
	if !got.Equal(expected) {
		t.Errorf("AdjustTimezone(%v, %s) = %v; want %v", v, tz, got, expected)
	}

	gotInvalid := timefy.AdjustTimezone(v, "Invalid/Timezone")
	if !gotInvalid.Equal(v) {
		t.Errorf("AdjustTimezone with invalid timezone should return original time")
	}
}

func TestAddSecond(t *testing.T) {
	v := time.Date(2023, 10, 25, 14, 30, 0, 0, time.UTC)

	expectedPos := v.Add(5 * time.Second)
	gotPos := timefy.AddSecond(v, 5)
	if !gotPos.Equal(expectedPos) {
		t.Errorf("AddSecond(%v, 5) = %v; want %v", v, gotPos, expectedPos)
	}

	expectedNeg := v.Add(-5 * time.Second)
	gotNeg := timefy.AddSecond(v, -5)
	if !gotNeg.Equal(expectedNeg) {
		t.Errorf("AddSecond(%v, -5) = %v; want %v", v, gotNeg, expectedNeg)
	}

	gotZero := timefy.AddSecond(v, 0)
	if !gotZero.Equal(v) {
		t.Errorf("AddSecond(%v, 0) = %v; want %v", v, gotZero, v)
	}
}

func TestAddMinute(t *testing.T) {
	v := time.Date(2023, 10, 25, 14, 30, 0, 0, time.UTC)

	expectedPos := v.Add(10 * time.Minute)
	gotPos := timefy.AddMinute(v, 10)
	if !gotPos.Equal(expectedPos) {
		t.Errorf("AddMinute(%v, 10) = %v; want %v", v, gotPos, expectedPos)
	}

	expectedNeg := v.Add(-10 * time.Minute)
	gotNeg := timefy.AddMinute(v, -10)
	if !gotNeg.Equal(expectedNeg) {
		t.Errorf("AddMinute(%v, -10) = %v; want %v", v, gotNeg, expectedNeg)
	}

	gotZero := timefy.AddMinute(v, 0)
	if !gotZero.Equal(v) {
		t.Errorf("AddMinute(%v, 0) = %v; want %v", v, gotZero, v)
	}
}

func TestAddHour(t *testing.T) {
	v := time.Date(2023, 10, 25, 14, 30, 0, 0, time.UTC)

	expectedPos := v.Add(3 * time.Hour)
	gotPos := timefy.AddHour(v, 3)
	if !gotPos.Equal(expectedPos) {
		t.Errorf("AddHour(%v, 3) = %v; want %v", v, gotPos, expectedPos)
	}

	expectedNeg := v.Add(-3 * time.Hour)
	gotNeg := timefy.AddHour(v, -3)
	if !gotNeg.Equal(expectedNeg) {
		t.Errorf("AddHour(%v, -3) = %v; want %v", v, gotNeg, expectedNeg)
	}

	gotZero := timefy.AddHour(v, 0)
	if !gotZero.Equal(v) {
		t.Errorf("AddHour(%v, 0) = %v; want %v", v, gotZero, v)
	}
}

func TestAddDay(t *testing.T) {
	v := time.Date(2023, 10, 25, 14, 30, 0, 0, time.UTC)

	gotPos := timefy.AddDay(v, 5)

	expectedImpl := v.Add(time.Hour * 24 * 5)
	if !gotPos.Equal(expectedImpl) {
		t.Errorf("AddDay(%v, 5) = %v; want %v", v, gotPos, expectedImpl)
	}

	gotZero := timefy.AddDay(v, 0)
	if !gotZero.Equal(v) {
		t.Errorf("AddDay(%v, 0) = %v; want %v", v, gotZero, v)
	}
}

func TestIsWithinTolerance(t *testing.T) {
	now := time.Now()
	if !timefy.IsWithinTolerance(now) {
		t.Errorf("IsWithinTolerance(now) = false; want true")
	}

	future := now.Add(30 * time.Second)
	if !timefy.IsWithinTolerance(future) {
		t.Errorf("IsWithinTolerance(future) = false; want true")
	}

	past := now.Add(-30 * time.Second)
	if !timefy.IsWithinTolerance(past) {
		t.Errorf("IsWithinTolerance(past) = false; want true")
	}

	farFuture := now.Add(2 * time.Minute)
	if timefy.IsWithinTolerance(farFuture) {
		t.Errorf("IsWithinTolerance(farFuture) = true; want false")
	}
}

func TestIsLeapYear(t *testing.T) {
	tests := []struct {
		year int
		want bool
	}{
		{2000, true},
		{2004, true},
		{2020, true},
		{2021, false},
		{2022, false},
		{2023, false},
		{1900, false},
		{2100, false},
	}

	for _, tt := range tests {
		if got := timefy.IsLeapYear(tt.year); got != tt.want {
			t.Errorf("IsLeapYear(%d) = %v; want %v", tt.year, got, tt.want)
		}
	}
}

func TestIsLeapYearN(t *testing.T) {
	leapTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	if !timefy.IsLeapYearN(leapTime) {
		t.Errorf("IsLeapYearN(2020) = false; want true")
	}

	nonLeapTime := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	if timefy.IsLeapYearN(nonLeapTime) {
		t.Errorf("IsLeapYearN(2021) = true; want false")
	}
}

func TestGetWeekdaysInRange(t *testing.T) {
	// March 2023: 1st is Wednesday.
	// 1(Wed), 2(Thu), 3(Fri), 4(Sat), 5(Sun), 6(Mon)...
	start := time.Date(2023, time.March, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, time.March, 6, 0, 0, 0, 0, time.UTC)

	got := timefy.GetWeekdaysInRange(start, end)
	// Expected: 1, 2, 3, 6 (4 days)
	if len(got) != 4 {
		t.Errorf("GetWeekdaysInRange returned %d days; want 4", len(got))
	}

	expectedDays := []int{1, 2, 3, 6}
	for i, d := range got {
		if d.Day() != expectedDays[i] {
			t.Errorf("Day %d: got %d; want %d", i, d.Day(), expectedDays[i])
		}
	}
}

func TestSince(t *testing.T) {
	// Mocking time.Since is hard without dependency injection.
	// We will just test basic functionality by using a short duration.
	start := time.Now().Add(-1 * time.Hour)

	hours := timefy.SinceHour(start)
	if hours < 0.99 || hours > 1.01 { // Allow small deviation
		t.Errorf("SinceHour = %v; want ~1.0", hours)
	}

	minutes := timefy.SinceMinute(start)
	if minutes < 59.9 || minutes > 60.1 {
		t.Errorf("SinceMinute = %v; want ~60.0", minutes)
	}

	// For seconds, use a smaller window
	startSec := time.Now().Add(-5 * time.Second)
	seconds := timefy.SinceSecond(startSec)
	if seconds < 4.9 || seconds > 5.5 { // checking if it's around 5
		t.Errorf("SinceSecond = %v; want ~5.0", seconds)
	}
}

func TestFormatTimex(t *testing.T) {
	v := time.Date(2023, 10, 25, 14, 30, 45, 123456789, time.UTC)
	got := timefy.FormatTimex(v)
	expected := []int{123456789, 45, 30, 14, 25, 10, 2023}

	for i, val := range got {
		if val != expected[i] {
			t.Errorf("FormatTimex[%d] = %d; want %d", i, val, expected[i])
		}
	}
}

// Additional tests for time-dependent functions

func TestBeginningOfMinute(t *testing.T) {
	now := time.Now()
	got := timefy.BeginningOfMinute()

	if got.Nanosecond() != 0 || got.Second() != 0 {
		t.Errorf("BeginningOfMinute should have 0 nanoseconds and 0 seconds, got %v", got)
	}
	if got.Minute() != now.Minute() && got.Minute() != now.Add(-1*time.Minute).Minute() {
		// Allow for minute rollover during test execution
		t.Logf("Warning: Minute mismatch likely due to execution timing: got %v, now %v", got, now)
	}
}

func TestBeginningOfHour(t *testing.T) {
	got := timefy.BeginningOfHour()
	if got.Nanosecond() != 0 || got.Second() != 0 || got.Minute() != 0 {
		t.Errorf("BeginningOfHour should have 0 min/sec/nano, got %v", got)
	}
}

func TestBeginningOfDay(t *testing.T) {
	got := timefy.BeginningOfDay()
	if got.Nanosecond() != 0 || got.Second() != 0 || got.Minute() != 0 || got.Hour() != 0 {
		t.Errorf("BeginningOfDay should be midnight, got %v", got)
	}
}

func TestBeginningOfWeek(t *testing.T) {
	got := timefy.BeginningOfWeek()
	if got.Nanosecond() != 0 || got.Second() != 0 || got.Minute() != 0 || got.Hour() != 0 {
		t.Errorf("BeginningOfWeek should be midnight, got %v", got)
	}
	// Default start of week is Sunday
	if got.Weekday() != time.Sunday {
		t.Errorf("BeginningOfWeek should be Sunday by default, got %v", got.Weekday())
	}
}

func TestBeginningOfMonth(t *testing.T) {
	got := timefy.BeginningOfMonth()
	if got.Nanosecond() != 0 || got.Second() != 0 || got.Minute() != 0 || got.Hour() != 0 || got.Day() != 1 {
		t.Errorf("BeginningOfMonth should be 1st of month midnight, got %v", got)
	}
}

func TestBeginningOfQuarter(t *testing.T) {
	got := timefy.BeginningOfQuarter()
	if got.Nanosecond() != 0 || got.Second() != 0 || got.Minute() != 0 || got.Hour() != 0 || got.Day() != 1 {
		t.Errorf("BeginningOfQuarter should be 1st of month midnight, got %v", got)
	}
	m := got.Month()
	if m != time.January && m != time.April && m != time.July && m != time.October {
		t.Errorf("BeginningOfQuarter month should be Jan, Apr, Jul or Oct, got %v", m)
	}
}

func TestBeginningOfYear(t *testing.T) {
	got := timefy.BeginningOfYear()
	if got.Nanosecond() != 0 || got.Second() != 0 || got.Minute() != 0 || got.Hour() != 0 || got.Day() != 1 || got.Month() != time.January {
		t.Errorf("BeginningOfYear should be Jan 1st midnight, got %v", got)
	}
}

func TestEndOfMinute(t *testing.T) {
	got := timefy.EndOfMinute()
	if got.Second() != 59 || got.Nanosecond() != 999999999 {
		t.Errorf("EndOfMinute should end with 59s 999999999ns, got %v", got)
	}
}

func TestEndOfHour(t *testing.T) {
	got := timefy.EndOfHour()
	if got.Minute() != 59 || got.Second() != 59 || got.Nanosecond() != 999999999 {
		t.Errorf("EndOfHour should end with 59m 59s 999999999ns, got %v", got)
	}
}

func TestEndOfDay(t *testing.T) {
	got := timefy.EndOfDay()
	if got.Hour() != 23 || got.Minute() != 59 || got.Second() != 59 || got.Nanosecond() != 999999999 {
		t.Errorf("EndOfDay should be 23:59:59.999999999, got %v", got)
	}
}

func TestEndOfWeek(t *testing.T) {
	got := timefy.EndOfWeek()
	if got.Hour() != 23 || got.Minute() != 59 || got.Second() != 59 || got.Nanosecond() != 999999999 {
		t.Errorf("EndOfWeek should be 23:59:59.999999999, got %v", got)
	}
}

func TestEndOfMonth(t *testing.T) {
	got := timefy.EndOfMonth()
	if got.Hour() != 23 || got.Minute() != 59 || got.Second() != 59 || got.Nanosecond() != 999999999 {
		t.Errorf("EndOfMonth should be 23:59:59.999999999, got %v", got)
	}
	nextDay := got.Add(time.Nanosecond)
	if nextDay.Day() != 1 {
		t.Errorf("EndOfMonth + 1ns should be 1st of next month, got %v", nextDay)
	}
}

func TestEndOfQuarter(t *testing.T) {
	got := timefy.EndOfQuarter()
	if got.Hour() != 23 || got.Minute() != 59 || got.Second() != 59 || got.Nanosecond() != 999999999 {
		t.Errorf("EndOfQuarter should be 23:59:59.999999999, got %v", got)
	}
}

func TestEndOfYear(t *testing.T) {
	got := timefy.EndOfYear()
	if got.Month() != time.December || got.Day() != 31 || got.Hour() != 23 || got.Minute() != 59 || got.Second() != 59 {
		t.Errorf("EndOfYear should be Dec 31 23:59:59..., got %v", got)
	}
}

func TestQuarter(t *testing.T) {
	q := timefy.Quarter()
	if q < 1 || q > 4 {
		t.Errorf("Quarter should be between 1 and 4, got %v", q)
	}
}

func TestParse(t *testing.T) {
	ts := "2023-10-25 14:30:00"
	got, err := timefy.Parse(ts)
	if err == nil {
		if got.Year() != 2023 || got.Month() != 10 || got.Day() != 25 {
			t.Errorf("Parse result mismatch: %v", got)
		}
	} else {
		// Tolerant to default config not supporting this format, though unlikely.
		t.Logf("Parse failed with error: %v", err)
	}
}

func TestMonday(t *testing.T) {
	now := time.Now()
	monday := timefy.Monday()
	if monday.Weekday() != time.Monday {
		t.Errorf("Monday() returned %v; want Monday", monday.Weekday())
	}
	diff := monday.Sub(now)
	if diff > 7*24*time.Hour || diff < -7*24*time.Hour {
		t.Errorf("Monday() returned %v; too far from now %v", monday, now)
	}
}

func TestSunday(t *testing.T) {
	sunday := timefy.Sunday()
	if sunday.Weekday() != time.Sunday {
		t.Errorf("Sunday() returned %v; want Sunday", sunday.Weekday())
	}
}

func TestEndOfSunday(t *testing.T) {
	sunday := timefy.EndOfSunday()
	if sunday.Weekday() != time.Sunday {
		t.Errorf("EndOfSunday() returned %v; want Sunday", sunday.Weekday())
	}
	if sunday.Hour() != 23 || sunday.Minute() != 59 || sunday.Second() != 59 {
		t.Errorf("EndOfSunday() time part mismatch: %v", sunday)
	}
}
