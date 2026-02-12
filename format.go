package timefy

import (
	"fmt"
	"time"
)

// TimeAgo returns a human-readable string representing the time difference between the current time and the specified time.
// It uses the local time zone for calculations and returns a relative time string (e.g., "just now", "5 minutes ago", "2 weeks ago").
//
// Parameters:
//   - `t`: A pointer to the Timex instance containing the time to compare.
//
// Returns:
//   - A string representing the relative time difference.
//
// Example:
//
//	t := Timex{Time: time.Now().Add(-5 * time.Minute)}
//	fmt.Println(t.TimeAgo()) // Output: "5 minutes ago"
func (t *Timex) TimeAgo() string {
	now := time.Now()
	diff := now.Sub(t.Time)

	seconds := int(diff.Seconds())
	minutes := int(diff.Minutes())
	hours := int(diff.Hours())
	days := int(diff.Hours() / 24)

	switch {
	case seconds < 60:
		return "just now"
	case minutes < 2:
		return "1 minute ago"
	case minutes < 60:
		return fmt.Sprintf("%d minutes ago", minutes)
	case hours < 2:
		return "1 hour ago"
	case hours < 24:
		return fmt.Sprintf("%d hours ago", hours)
	case days < 2:
		return "1 day ago"
	case days < 7:
		return fmt.Sprintf("%d days ago", days)
	case days < 14:
		return "1 week ago"
	case days < 30:
		return fmt.Sprintf("%d weeks ago", days/7)
	case days < 60:
		return "1 month ago"
	case days < 365:
		return fmt.Sprintf("%d months ago", days/30)
	case days < 730:
		return "1 year ago"
	default:
		return fmt.Sprintf("%d years ago", days/365)
	}
}

// TimeUntil returns a human-readable string representing the time difference between the current time and the specified time.
// It uses the local time zone for calculations and returns a relative time string (e.g., "in a few seconds", "in 5 minutes", "in 2 weeks").
//
// Parameters:
//   - `t`: A pointer to the Timex instance containing the time to compare.
//
// Returns:
//   - A string representing the relative time difference.
//
// Example:
//
//	t := Timex{Time: time.Now().Add(5 * time.Minute)}
//	fmt.Println(t.TimeUntil()) // Output: "in 5 minutes"
func (t *Timex) TimeUntil() string {
	diff := time.Until(t.Time)

	if diff < 0 {
		return "in the past"
	}

	seconds := int(diff.Seconds())
	minutes := int(diff.Minutes())
	hours := int(diff.Hours())
	days := int(diff.Hours() / 24)

	switch {
	case seconds < 60:
		return "in a few seconds"
	case minutes < 2:
		return "in 1 minute"
	case minutes < 60:
		return fmt.Sprintf("in %d minutes", minutes)
	case hours < 2:
		return "in 1 hour"
	case hours < 24:
		return fmt.Sprintf("in %d hours", hours)
	case days < 2:
		return "in 1 day"
	case days < 7:
		return fmt.Sprintf("in %d days", days)
	case days < 14:
		return "in 1 week"
	case days < 30:
		return fmt.Sprintf("in %d weeks", days/7)
	case days < 60:
		return "in 1 month"
	case days < 365:
		return fmt.Sprintf("in %d months", days/30)
	case days < 730:
		return "in 1 year"
	default:
		return fmt.Sprintf("in %d years", days/365)
	}
}
