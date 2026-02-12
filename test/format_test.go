package test

import (
	"testing"
	"time"

	"github.com/sivaosorg/timefy"
)

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
