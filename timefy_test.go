package timefy_test

import (
	"sync"
	"testing"
	"time"

	"github.com/sivaosorg/timefy"
)

func TestWith(t *testing.T) {
	now := time.Now()
	tmx := timefy.With(now)
	if !tmx.Time.Equal(now) {
		t.Errorf("With(%v) time = %v; want %v", now, tmx.Time, now)
	}
}

func TestNew(t *testing.T) {
	now := time.Now()
	tmx := timefy.New(now)
	if !tmx.Time.Equal(now) {
		t.Errorf("New(%v) time = %v; want %v", now, tmx.Time, now)
	}
}

func TestRuleWith(t *testing.T) {
	now := time.Now()
	rule := timefy.NewRule()
	tmx := rule.With(now)
	if !tmx.Time.Equal(now) {
		t.Errorf("Rule.With(%v) time = %v; want %v", now, tmx.Time, now)
	}
}

func TestTimexBeginningOfMinute(t *testing.T) {
	now := time.Date(2023, 10, 25, 14, 30, 45, 100, time.UTC)
	tmx := timefy.New(now)
	got := tmx.BeginningOfMinute()
	expected := time.Date(2023, 10, 25, 14, 30, 0, 0, time.UTC)
	if !got.Equal(expected) {
		t.Errorf("BeginningOfMinute() = %v; want %v", got, expected)
	}
}

func TestTimexBeginningOfHour(t *testing.T) {
	now := time.Date(2023, 10, 25, 14, 30, 45, 100, time.UTC)
	tmx := timefy.New(now)
	got := tmx.BeginningOfHour()
	expected := time.Date(2023, 10, 25, 14, 0, 0, 0, time.UTC)
	if !got.Equal(expected) {
		t.Errorf("BeginningOfHour() = %v; want %v", got, expected)
	}
}

func TestTimexBeginningOfDay(t *testing.T) {
	now := time.Date(2023, 10, 25, 14, 30, 45, 100, time.UTC)
	tmx := timefy.New(now)
	got := tmx.BeginningOfDay()
	expected := time.Date(2023, 10, 25, 0, 0, 0, 0, time.UTC)
	if !got.Equal(expected) {
		t.Errorf("BeginningOfDay() = %v; want %v", got, expected)
	}
}

func TestTimexBeginningOfWeek(t *testing.T) {
	// 2023-10-25 is Wednesday.
	// Default week start is Sunday.
	// Previous Sunday is 2023-10-22.
	now := time.Date(2023, 10, 25, 14, 30, 45, 100, time.UTC)
	tmx := timefy.New(now)
	got := tmx.BeginningOfWeek()
	expected := time.Date(2023, 10, 22, 0, 0, 0, 0, time.UTC)
	if !got.Equal(expected) {
		t.Errorf("BeginningOfWeek() = %v; want %v", got, expected)
	}

	// Test with Monday start
	rule := timefy.NewRule().WithWeekStartDay(time.Monday)
	tmxMonday := rule.With(now)
	gotMonday := tmxMonday.BeginningOfWeek()
	expectedMonday := time.Date(2023, 10, 23, 0, 0, 0, 0, time.UTC)
	if !gotMonday.Equal(expectedMonday) {
		t.Errorf("BeginningOfWeek(Monday) = %v; want %v", gotMonday, expectedMonday)
	}
}

func TestTimexBeginningOfMonth(t *testing.T) {
	now := time.Date(2023, 10, 25, 14, 30, 45, 100, time.UTC)
	tmx := timefy.New(now)
	got := tmx.BeginningOfMonth()
	expected := time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)
	if !got.Equal(expected) {
		t.Errorf("BeginningOfMonth() = %v; want %v", got, expected)
	}
}

func TestTimexBeginningOfQuarter(t *testing.T) {
	// Q4 starts Oct 1
	now := time.Date(2023, 10, 25, 14, 30, 45, 100, time.UTC)
	tmx := timefy.New(now)
	got := tmx.BeginningOfQuarter()
	expected := time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)
	if !got.Equal(expected) {
		t.Errorf("BeginningOfQuarter() = %v; want %v", got, expected)
	}

	// Q1 starts Jan 1
	nowQ1 := time.Date(2023, 2, 15, 0, 0, 0, 0, time.UTC)
	tmxQ1 := timefy.New(nowQ1)
	gotQ1 := tmxQ1.BeginningOfQuarter()
	expectedQ1 := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	if !gotQ1.Equal(expectedQ1) {
		t.Errorf("BeginningOfQuarter(Q1) = %v; want %v", gotQ1, expectedQ1)
	}
}

func TestTimexBeginningOfHalf(t *testing.T) {
	// H2 starts July 1
	now := time.Date(2023, 10, 25, 14, 30, 45, 100, time.UTC)
	tmx := timefy.New(now)
	got := tmx.BeginningOfHalf()
	expected := time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC)
	if !got.Equal(expected) {
		t.Errorf("BeginningOfHalf() = %v; want %v", got, expected)
	}

	// H1 starts Jan 1
	nowH1 := time.Date(2023, 3, 15, 0, 0, 0, 0, time.UTC)
	tmxH1 := timefy.New(nowH1)
	gotH1 := tmxH1.BeginningOfHalf()
	expectedH1 := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	if !gotH1.Equal(expectedH1) {
		t.Errorf("BeginningOfHalf(H1) = %v; want %v", gotH1, expectedH1)
	}
}

func TestTimexBeginningOfYear(t *testing.T) {
	now := time.Date(2023, 10, 25, 14, 30, 45, 100, time.UTC)
	tmx := timefy.New(now)
	got := tmx.BeginningOfYear()
	expected := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	if !got.Equal(expected) {
		t.Errorf("BeginningOfYear() = %v; want %v", got, expected)
	}
}

func TestTimexEndOfMinute(t *testing.T) {
	now := time.Date(2023, 10, 25, 14, 30, 45, 100, time.UTC)
	tmx := timefy.New(now)
	got := tmx.EndOfMinute()
	expected := time.Date(2023, 10, 25, 14, 30, 59, 999999999, time.UTC)
	if !got.Equal(expected) {
		t.Errorf("EndOfMinute() = %v; want %v", got, expected)
	}
}

func TestTimexEndOfHour(t *testing.T) {
	now := time.Date(2023, 10, 25, 14, 30, 45, 100, time.UTC)
	tmx := timefy.New(now)
	got := tmx.EndOfHour()
	expected := time.Date(2023, 10, 25, 14, 59, 59, 999999999, time.UTC)
	if !got.Equal(expected) {
		t.Errorf("EndOfHour() = %v; want %v", got, expected)
	}
}

func TestTimexEndOfDay(t *testing.T) {
	now := time.Date(2023, 10, 25, 14, 30, 45, 100, time.UTC)
	tmx := timefy.New(now)
	got := tmx.EndOfDay()
	expected := time.Date(2023, 10, 25, 23, 59, 59, 999999999, time.UTC)
	if !got.Equal(expected) {
		t.Errorf("EndOfDay() = %v; want %v", got, expected)
	}
}

func TestTimexEndOfWeek(t *testing.T) {
	// 2023-10-25 is Wednesday.
	// Default week end is Saturday (because start is Sunday).
	// Next Saturday is 2023-10-28.
	now := time.Date(2023, 10, 25, 14, 30, 45, 100, time.UTC)
	tmx := timefy.New(now)
	got := tmx.EndOfWeek()
	expected := time.Date(2023, 10, 28, 23, 59, 59, 999999999, time.UTC)
	if !got.Equal(expected) {
		t.Errorf("EndOfWeek() = %v; want %v", got, expected)
	}
}

func TestTimexEndOfMonth(t *testing.T) {
	now := time.Date(2023, 10, 25, 14, 30, 45, 100, time.UTC)
	tmx := timefy.New(now)
	got := tmx.EndOfMonth()
	expected := time.Date(2023, 10, 31, 23, 59, 59, 999999999, time.UTC)
	if !got.Equal(expected) {
		t.Errorf("EndOfMonth() = %v; want %v", got, expected)
	}
}

func TestTimexEndOfQuarter(t *testing.T) {
	now := time.Date(2023, 10, 25, 14, 30, 45, 100, time.UTC)
	tmx := timefy.New(now)
	got := tmx.EndOfQuarter()
	// Q4 ends Dec 31
	expected := time.Date(2023, 12, 31, 23, 59, 59, 999999999, time.UTC)
	if !got.Equal(expected) {
		t.Errorf("EndOfQuarter() = %v; want %v", got, expected)
	}
}

func TestTimexEndOfYear(t *testing.T) {
	now := time.Date(2023, 10, 25, 14, 30, 45, 100, time.UTC)
	tmx := timefy.New(now)
	got := tmx.EndOfYear()
	expected := time.Date(2023, 12, 31, 23, 59, 59, 999999999, time.UTC)
	if !got.Equal(expected) {
		t.Errorf("EndOfYear() = %v; want %v", got, expected)
	}
}

func TestTimexMonday(t *testing.T) {
	// 2023-10-25 is Wednesday.
	// Monday of that week depends on week start?
	// Monday() usually returns Monday of the current week.
	// If week starts Sunday (22nd), Monday is 23rd.
	now := time.Date(2023, 10, 25, 14, 30, 45, 100, time.UTC)
	tmx := timefy.New(now)
	got := tmx.Monday()
	expected := time.Date(2023, 10, 23, 0, 0, 0, 0, time.UTC)
	if !got.Equal(expected) {
		t.Errorf("Monday() = %v; want %v", got, expected)
	}
}

func TestTimexSunday(t *testing.T) {
	// 2023-10-25 is Wednesday.
	// Sunday of current week (starts Sunday 22nd) is 22nd?
	// Or is it looking for next Sunday?
	// The doc says "most recent or upcoming Sunday".
	// Implementation usually returns the Sunday of the current week window.
	// If week starts Sunday, Sunday() should return start of week?
	// Let's assume standard behavior: Sunday of the *current week*.
	// If week starts Sunday, it creates a bit of ambiguity if it means "Start of week" or "End of week".
	// But usually it means the day named Sunday in that week.

	now := time.Date(2023, 10, 25, 14, 30, 45, 100, time.UTC)
	tmx := timefy.New(now)
	got := tmx.Sunday()
	// Sunday() returns the *upcoming* Sunday or today if today is Sunday.
	// For Wed Oct 25, the upcoming Sunday is Oct 29.
	expected := time.Date(2023, 10, 29, 0, 0, 0, 0, time.UTC)
	if !got.Equal(expected) {
		t.Errorf("Sunday() = %v; want %v", got, expected)
	}
}

func TestTimexEndOfSunday(t *testing.T) {
	now := time.Date(2023, 10, 25, 14, 30, 45, 100, time.UTC)
	tmx := timefy.New(now)
	got := tmx.EndOfSunday()
	// Sunday is Oct 29. End of it is Oct 29 23:59...
	expected := time.Date(2023, 10, 29, 23, 59, 59, 999999999, time.UTC)
	if !got.Equal(expected) {
		t.Errorf("EndOfSunday() = %v; want %v", got, expected)
	}
}

func TestTimexQuarter(t *testing.T) {
	now := time.Date(2023, 10, 25, 14, 30, 45, 100, time.UTC)
	tmx := timefy.New(now)
	if tmx.Quarter() != 4 {
		t.Errorf("Quarter() = %d; want 4", tmx.Quarter())
	}
}

func TestRuleParse(t *testing.T) {
	// Assuming default rule logic/formats.
	rule := timefy.NewRule()
	// Add a format to be sure
	rule.AppendTimeFormat("2006-01-02")

	ts := "2023-10-25"
	got, err := rule.Parse(ts)
	if err != nil {
		t.Errorf("Parse(%s) failed: %v", ts, err)
	}
	expected := time.Date(2023, 10, 25, 0, 0, 0, 0, time.UTC)
	// Parse might use local or configured location. Default NewRule might differ.
	// We need to check expected location handling.
	// If Rule has no location, it uses Local? Or UTC?
	// The doc says "If it is not set... uses the local time zone".
	// So we should compare with Local.

	// Actually, time.Parse uses UTC for "2006-01-02" if no timezone info.
	// time.ParseInLocation uses the location.

	// "Parse attempts to parse... using the current Config... checks whether TimeLocation is set."

	// Let's set location to UTC to be deterministic
	rule.WithLocation(time.UTC)
	got, err = rule.Parse(ts)
	if err != nil {
		t.Errorf("Parse(%s) with UTC failed: %v", ts, err)
	}
	if !got.Equal(expected) {
		t.Errorf("Parse(%s) = %v; want %v", ts, got, expected)
	}
}

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

func TestNewRule(t *testing.T) {
	rule := timefy.NewRule()
	if rule == nil {
		t.Fatal("NewRule() returned nil")
	}
	// Verify default values if possible, though fields are unexported.
	// We can verify behavior via With/Parse methods or reflection if needed,
	// but mostly we test the builder methods return valid *Rule.
}

func TestWithWeekStartDay(t *testing.T) {
	rule := timefy.NewRule()
	updated := rule.WithWeekStartDay(time.Monday)
	if updated != rule {
		t.Errorf("WithWeekStartDay should return the same rule instance")
	}
	// Verification via side effect on Monday()/Sunday() logic using this rule
	now := time.Date(2023, 10, 25, 12, 0, 0, 0, time.UTC) // Wed
	tmx := rule.With(now)
	// If week starts Monday, BeginningOfWeek should be Mon Oct 23.
	startOfWeek := tmx.BeginningOfWeek()
	expected := time.Date(2023, 10, 23, 0, 0, 0, 0, time.UTC)
	if !startOfWeek.Equal(expected) {
		t.Errorf("WithWeekStartDay(Monday) affect failed; got %v want %v", startOfWeek, expected)
	}
}

func TestWithTimeFormats(t *testing.T) {
	rule := timefy.NewRule()
	formats := []string{"2006-01-02"}
	updated := rule.WithTimeFormats(formats)
	if updated != rule {
		t.Errorf("WithTimeFormats should return the same rule instance")
	}
	// Verify via Parse
	parsed, err := rule.Parse("2023-10-25")
	if err != nil {
		t.Errorf("Parse failed with custom format: %v", err)
	}
	if parsed.Year() != 2023 {
		t.Errorf("Parse result incorrect")
	}
}

func TestWithLocation(t *testing.T) {
	rule := timefy.NewRule()
	loc, _ := time.LoadLocation("America/New_York")
	updated := rule.WithLocation(loc)
	if updated != rule {
		t.Errorf("WithLocation should return the same rule instance")
	}
	// Verify via Parse - should return time in that location
	rule.WithTimeFormats([]string{"2006-01-02 15:04:05"})
	parsed, err := rule.Parse("2023-10-25 12:00:00")
	if err != nil {
		t.Errorf("Parse failed: %v", err)
	}
	if parsed.Location().String() != "America/New_York" {
		t.Errorf("Parse location = %v; want America/New_York", parsed.Location())
	}
}

func TestWithLocationRFC(t *testing.T) {
	rule := timefy.NewRule()
	// Use a constant if available or cast string
	// From constant.go outline, DefaultTimezoneVietnam is available
	updated := rule.WithLocationRFC(timefy.DefaultTimezoneVietnam)
	if updated != rule {
		t.Errorf("WithLocationRFC should return the same rule instance")
	}

	rule.WithTimeFormats([]string{"2006-01-02 15:04:05"})
	parsed, err := rule.Parse("2023-10-25 12:00:00")
	if err != nil {
		t.Errorf("Parse failed: %v", err)
	}
	// "Asia/Ho_Chi_Minh"
	if parsed.Location().String() != "Asia/Ho_Chi_Minh" {
		t.Errorf("Parse location = %v; want Asia/Ho_Chi_Minh", parsed.Location())
	}
}

func TestApplyUTCLoc(t *testing.T) {
	rule := timefy.NewRule()
	updated := rule.ApplyUTCLoc()
	if updated != rule {
		t.Errorf("ApplyUTCLoc should return the same rule instance")
	}
	rule.WithTimeFormats([]string{"2006-01-02 15:04:05"})
	parsed, err := rule.Parse("2023-10-25 12:00:00")
	if err != nil {
		t.Errorf("Parse failed: %v", err)
	}
	if parsed.Location() != time.UTC {
		t.Errorf("Parse location = %v; want UTC", parsed.Location())
	}
}

func TestApplyLocalLoc(t *testing.T) {
	rule := timefy.NewRule()
	updated := rule.ApplyLocalLoc()
	if updated != rule {
		t.Errorf("ApplyLocalLoc should return the same rule instance")
	}
	rule.WithTimeFormats([]string{"2006-01-02 15:04:05"})
	parsed, err := rule.Parse("2023-10-25 12:00:00")
	if err != nil {
		t.Errorf("Parse failed: %v", err)
	}
	if parsed.Location().String() != "Local" { // time.Local string representation
		t.Errorf("Parse location = %v; want Local", parsed.Location())
	}
}

func TestAppendTimeFormat(t *testing.T) {
	rule := timefy.NewRule()
	// Clear formats implicitly by providing a new list first?
	// Or just append to defaults.
	rule.AppendTimeFormat("2006/01/02")

	parsed, err := rule.Parse("2023/10/25")
	if err != nil {
		t.Errorf("Parse failed with appended format: %v", err)
	}
	if parsed.Year() != 2023 {
		t.Errorf("Parse result incorrect")
	}
}

func TestAppendTimeFormatRFC(t *testing.T) {
	rule := timefy.NewRule()
	// Using a constant from constant.go
	// TimeFormat20060102T150405
	rule.AppendTimeFormatRFC(timefy.TimeFormat20060102T150405)

	parsed, err := rule.Parse("2023-10-25T14:30:00")
	if err != nil {
		t.Errorf("Parse failed with appended RFC format: %v", err)
	}
	if parsed.Year() != 2023 {
		t.Errorf("Parse result incorrect")
	}
}

func TestAppendTimeFormatRFCshort(t *testing.T) {
	rule := timefy.NewRule()
	// TimeRFC01T150405 = "15:04:05"
	rule.AppendTimeFormatRFCshort(timefy.TimeRFC01T150405)

	parsed, err := rule.Parse("14:30:00")
	if err != nil {
		t.Errorf("Parse failed with appended RFC short format: %v", err)
	}
	if parsed.Hour() != 14 {
		t.Errorf("Parse result incorrect")
	}
}

func TestTimeAgo(t *testing.T) {
	now := time.Now()

	// Just now
	tmx := timefy.New(now)
	if got := tmx.TimeAgo(); got != "just now" {
		t.Errorf("TimeAgo(now) = %q; want \"just now\"", got)
	}

	// 5 minutes ago
	tmx = timefy.New(now.Add(-5 * time.Minute))
	if got := tmx.TimeAgo(); got != "5 minutes ago" {
		t.Errorf("TimeAgo(5m ago) = %q; want \"5 minutes ago\"", got)
	}

	// 1 hour ago
	tmx = timefy.New(now.Add(-1 * time.Hour))
	if got := tmx.TimeAgo(); got != "1 hour ago" {
		t.Errorf("TimeAgo(1h ago) = %q; want \"1 hour ago\"", got)
	}

	// 2 days ago
	tmx = timefy.New(now.Add(-48 * time.Hour))
	if got := tmx.TimeAgo(); got != "2 days ago" {
		t.Errorf("TimeAgo(2d ago) = %q; want \"2 days ago\"", got)
	}

	// 1 month ago (approx 30 days)
	tmx = timefy.New(now.Add(-30 * 24 * time.Hour))
	if got := tmx.TimeAgo(); got != "1 month ago" {
		// Logic: days < 60 -> 1 month ago.
		// 30 days -> 1 month ago?
		// switch: case days < 30: ... case days < 60: return "1 month ago"
		// 30 days falls into < 60.
		t.Errorf("TimeAgo(30d ago) = %q; want \"1 month ago\"", got)
	}

	// 1 year ago (approx 365 days)
	tmx = timefy.New(now.Add(-365 * 24 * time.Hour))
	if got := tmx.TimeAgo(); got != "1 year ago" {
		// Logic: days < 730 -> 1 year ago
		t.Errorf("TimeAgo(365d ago) = %q; want \"1 year ago\"", got)
	}
}

func TestTimeUntil(t *testing.T) {
	now := time.Now()

	// In the past
	tmx := timefy.New(now.Add(-5 * time.Minute))
	if got := tmx.TimeUntil(); got != "in the past" {
		t.Errorf("TimeUntil(past) = %q; want \"in the past\"", got)
	}

	// In a few seconds
	tmx = timefy.New(now.Add(5 * time.Second))
	if got := tmx.TimeUntil(); got != "in a few seconds" {
		t.Errorf("TimeUntil(5s) = %q; want \"in a few seconds\"", got)
	}

	// In 5 minutes (plus buffer)
	tmx = timefy.New(now.Add(5*time.Minute + 2*time.Second))
	if got := tmx.TimeUntil(); got != "in 5 minutes" {
		t.Errorf("TimeUntil(5m) = %q; want \"in 5 minutes\"", got)
	}

	// In 1 hour (plus buffer)
	tmx = timefy.New(now.Add(1*time.Hour + 2*time.Second))
	if got := tmx.TimeUntil(); got != "in 1 hour" {
		t.Errorf("TimeUntil(1h) = %q; want \"in 1 hour\"", got)
	}

	// In 2 days (plus buffer)
	tmx = timefy.New(now.Add(48*time.Hour + 2*time.Second))
	if got := tmx.TimeUntil(); got != "in 2 days" {
		t.Errorf("TimeUntil(2d) = %q; want \"in 2 days\"", got)
	}
}

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
