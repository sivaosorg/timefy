package test

import (
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
