package test

import (
	"sync"
	"testing"
	"time"

	"github.com/sivaosorg/timefy"
)

// TestGetSetWeekStartDay tests the thread-safe getter and setter for WeekStartDay.
func TestGetSetWeekStartDay(t *testing.T) {
	// Store original value to restore later
	original := timefy.GetWeekStartDay()
	defer timefy.SetWeekStartDay(original)

	// Test setting to Monday
	timefy.SetWeekStartDay(time.Monday)
	if got := timefy.GetWeekStartDay(); got != time.Monday {
		t.Errorf("GetWeekStartDay() = %v; want Monday", got)
	}

	// Test setting to Sunday
	timefy.SetWeekStartDay(time.Sunday)
	if got := timefy.GetWeekStartDay(); got != time.Sunday {
		t.Errorf("GetWeekStartDay() = %v; want Sunday", got)
	}
}

// TestWeekStartDayConcurrency tests that GetWeekStartDay and SetWeekStartDay
// are safe to call concurrently.
func TestWeekStartDayConcurrency(t *testing.T) {
	// Store original value to restore later
	original := timefy.GetWeekStartDay()
	defer timefy.SetWeekStartDay(original)

	const goroutines = 100
	const iterations = 100

	var wg sync.WaitGroup
	wg.Add(goroutines * 2)

	// Concurrent writers
	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				day := time.Weekday(id % 7)
				timefy.SetWeekStartDay(day)
			}
		}(i)
	}

	// Concurrent readers
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				_ = timefy.GetWeekStartDay()
			}
		}()
	}

	wg.Wait()
	// Test passes if no race conditions detected
}

// TestSafeWeekdayMethods tests the safe weekday methods that return errors
// instead of panicking.
func TestSafeWeekdayMethods(t *testing.T) {
	// Oct 25, 2023 is a Wednesday (weekday 3)
	now := time.Date(2023, 10, 25, 12, 0, 0, 0, time.UTC) // Wednesday
	tmx := timefy.New(now)

	// The weekday methods calculate the target day in the current "week window"
	// using the formula: parseTime.AddDate(0, 0, -weekday+offset)
	// For Oct 25 (Wed, weekday=3):
	// - Monday (offset=1): -3+1 = -2 → Oct 23
	// - Tuesday (offset=2): -3+2 = -1 → Oct 24
	// - Wednesday (offset=3): -3+3 = 0 → Oct 25
	// - Thursday (offset=4): -3+4 = +1 → Oct 26
	// - Friday (offset=5): -3+5 = +2 → Oct 27
	// - Saturday (offset=6): -3+6 = +3 → Oct 28
	// - Sunday uses (7-weekday) = 7-3 = +4 → Oct 29
	tests := []struct {
		name       string
		method     func(...string) (time.Time, error)
		wantDay    time.Weekday
		wantMonth  time.Month
		wantDayNum int
	}{
		{"MondaySafe", tmx.MondaySafe, time.Monday, time.October, 23},
		{"TuesdaySafe", tmx.TuesdaySafe, time.Tuesday, time.October, 24},
		{"WednesdaySafe", tmx.WednesdaySafe, time.Wednesday, time.October, 25},
		{"ThursdaySafe", tmx.ThursdaySafe, time.Thursday, time.October, 26},
		{"FridaySafe", tmx.FridaySafe, time.Friday, time.October, 27},
		{"SaturdaySafe", tmx.SaturdaySafe, time.Saturday, time.October, 28},
		{"SundaySafe", tmx.SundaySafe, time.Sunday, time.October, 29},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.method()
			if err != nil {
				t.Errorf("%s() unexpected error: %v", tt.name, err)
				return
			}
			if got.Weekday() != tt.wantDay {
				t.Errorf("%s() weekday = %v; want %v", tt.name, got.Weekday(), tt.wantDay)
			}
			if got.Month() != tt.wantMonth {
				t.Errorf("%s() month = %v; want %v", tt.name, got.Month(), tt.wantMonth)
			}
			if got.Day() != tt.wantDayNum {
				t.Errorf("%s() day = %v; want %v", tt.name, got.Day(), tt.wantDayNum)
			}
		})
	}
}

// TestSafeWeekdayMethodsWithInvalidInput tests that safe weekday methods
// return errors on invalid input instead of panicking.
func TestSafeWeekdayMethodsWithInvalidInput(t *testing.T) {
	now := time.Date(2023, 10, 25, 12, 0, 0, 0, time.UTC)
	rule := timefy.NewRule().WithTimeFormats([]string{"2006-01-02"}) // Only accept this format
	tmx := rule.With(now)

	invalidInput := "not-a-valid-date"

	methods := []struct {
		name   string
		method func(...string) (time.Time, error)
	}{
		{"MondaySafe", tmx.MondaySafe},
		{"TuesdaySafe", tmx.TuesdaySafe},
		{"WednesdaySafe", tmx.WednesdaySafe},
		{"ThursdaySafe", tmx.ThursdaySafe},
		{"FridaySafe", tmx.FridaySafe},
		{"SaturdaySafe", tmx.SaturdaySafe},
		{"SundaySafe", tmx.SundaySafe},
	}

	for _, tt := range methods {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.method(invalidInput)
			if err == nil {
				t.Errorf("%s(%q) expected error, got nil", tt.name, invalidInput)
			}
		})
	}
}

// TestTimexWithNilRule tests that Timex methods handle nil Rule gracefully.
func TestTimexWithNilRule(t *testing.T) {
	// Create a Timex directly without a Rule (not recommended but possible)
	now := time.Date(2023, 10, 25, 12, 0, 0, 0, time.UTC)
	tmx := &timefy.Timex{Time: now}

	// These methods should use defaults when Rule is nil
	t.Run("BeginningOfWeek", func(t *testing.T) {
		got := tmx.BeginningOfWeek()
		// Default week starts on Sunday (Oct 22)
		expected := time.Date(2023, 10, 22, 0, 0, 0, 0, time.UTC)
		if !got.Equal(expected) {
			t.Errorf("BeginningOfWeek() with nil Rule = %v; want %v", got, expected)
		}
	})

	t.Run("BeginningOfDay", func(t *testing.T) {
		got := tmx.BeginningOfDay()
		expected := time.Date(2023, 10, 25, 0, 0, 0, 0, time.UTC)
		if !got.Equal(expected) {
			t.Errorf("BeginningOfDay() with nil Rule = %v; want %v", got, expected)
		}
	})
}

// TestGetWeekdaysInRangeLeapYear tests that GetWeekdaysInRange works correctly
// with leap years (previously had a bug).
func TestGetWeekdaysInRangeLeapYear(t *testing.T) {
	// February 2020 was a leap year
	// Feb 24 (Mon), 25 (Tue), 26 (Wed), 27 (Thu), 28 (Fri), 29 (Sat), March 1 (Sun), March 2 (Mon)
	start := time.Date(2020, time.February, 24, 0, 0, 0, 0, time.UTC) // Monday
	end := time.Date(2020, time.March, 2, 0, 0, 0, 0, time.UTC)       // Monday

	got := timefy.GetWeekdaysInRange(start, end)

	// Weekdays: Feb 24 (Mon), 25 (Tue), 26 (Wed), 27 (Thu), 28 (Fri), March 2 (Mon) = 6 days
	// Feb 29 is Saturday, March 1 is Sunday - both skipped
	if len(got) != 6 {
		t.Errorf("GetWeekdaysInRange returned %d days; want 6", len(got))
	}

	// Verify Feb 28 (Friday) is included
	hasFeb28 := false
	for _, d := range got {
		if d.Month() == time.February && d.Day() == 28 {
			hasFeb28 = true
			break
		}
	}
	if !hasFeb28 {
		t.Error("GetWeekdaysInRange should include Feb 28 (Friday)")
	}
}

// TestGetWeekdaysInRangeNonLeapYear tests that GetWeekdaysInRange works correctly
// with non-leap years (previously had a bug where non-leap year days might be skipped).
func TestGetWeekdaysInRangeNonLeapYear(t *testing.T) {
	// February 2023 was not a leap year
	start := time.Date(2023, time.February, 27, 0, 0, 0, 0, time.UTC) // Monday
	end := time.Date(2023, time.March, 3, 0, 0, 0, 0, time.UTC)       // Friday

	got := timefy.GetWeekdaysInRange(start, end)

	// Weekdays: Feb 27 (Mon), 28 (Tue), March 1 (Wed), 2 (Thu), 3 (Fri) = 5 days
	if len(got) != 5 {
		t.Errorf("GetWeekdaysInRange returned %d days; want 5", len(got))
		for i, d := range got {
			t.Logf("Day %d: %v", i, d)
		}
	}
}
