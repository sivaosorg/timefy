package test

import (
	"testing"
	"time"

	"github.com/sivaosorg/timefy"
)

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
