// Package main demonstrates multi-timezone global team coordination using the
// timefy library.
//
// Real-world scenario: a distributed engineering organization spans five cities
// (Ho Chi Minh City, London, New York, Tokyo, Sydney).  This program:
//
//  1. Finds the overlap window where all offices are in business hours.
//  2. Renders a "world clock" snapshot for any given UTC instant.
//  3. Detects which team members are outside their business hours for a
//     proposed meeting time.
//  4. Generates the next 5 weekdays on which a cross-team standup can occur
//     without falling on any office's public holiday.
//
// Run with:
//
//	cd _example/global_team && go run main.go
package main

import (
	"fmt"
	"time"

	"github.com/sivaosorg/timefy"
)

// Office describes a single regional office.
type Office struct {
	Name            string
	Zone            timefy.ZoneRFC
	BusinessStart   int // hour of day (local time) when office opens, e.g. 9
	BusinessEnd     int // hour of day (local time) when office closes, e.g. 18
	PublicHolidays  []time.Time // simplified: a few fixed holiday dates in UTC
}

// loc is a convenience method that loads the *time.Location for the office.
func (o Office) loc() *time.Location {
	l, err := time.LoadLocation(string(o.Zone))
	if err != nil {
		return time.UTC
	}
	return l
}

// IsBusinessHour reports whether `t` (in any timezone) falls within the
// office's local business hours on a weekday.
//
// Implementation note:
//   - We first convert `t` to the office's local timezone with
//     timefy.AdjustTimezone.
//   - We then use timefy.New(localT).IsWeekday() to exclude weekends.
//   - Hour comparison uses the raw Hour() component of the converted time.
func (o Office) IsBusinessHour(t time.Time) bool {
	localT := timefy.AdjustTimezone(t, string(o.Zone))
	if !timefy.New(localT).IsWeekday() {
		return false
	}
	h := localT.Hour()
	return h >= o.BusinessStart && h < o.BusinessEnd
}

// LocalTime converts a UTC time to the office's local time string.
func (o Office) LocalTime(utc time.Time) string {
	local := timefy.AdjustTimezone(utc, string(o.Zone))
	return local.Format("Mon 2006-01-02 15:04 MST")
}

// IsHoliday reports whether `t` falls on a public holiday for this office.
// The comparison uses IsSameDay so the exact hour does not matter.
func (o Office) IsHoliday(t time.Time) bool {
	local := timefy.AdjustTimezone(t, string(o.Zone))
	tx := timefy.New(local)
	for _, h := range o.PublicHolidays {
		if tx.IsSameDay(h) {
			return true
		}
	}
	return false
}

// ─────────────────────────────────────────────────────────────────────────────
// globalOverlapSlots scans `numHours` consecutive one-hour UTC slots starting
// from `fromUTC` and returns all slots where EVERY office is in business hours.
//
// Why scan hour-by-hour?
//
//	Timezone offsets can be as fine as 30 minutes (e.g., India +5:30, Nepal
//	+5:45), so a 1-hour granularity is the correct minimal unit that gives
//	deterministic results without sub-hour complexity.
// ─────────────────────────────────────────────────────────────────────────────
func globalOverlapSlots(offices []Office, fromUTC time.Time, numHours int) []time.Time {
	var slots []time.Time
	for i := 0; i < numHours; i++ {
		candidate := timefy.AddHour(fromUTC, i)
		allOpen := true
		for _, o := range offices {
			if !o.IsBusinessHour(candidate) {
				allOpen = false
				break
			}
		}
		if allOpen {
			slots = append(slots, candidate)
		}
	}
	return slots
}

// ─────────────────────────────────────────────────────────────────────────────
// nextWorkingDays returns the next `n` calendar days (starting from `from`)
// that are weekdays in ALL provided offices and are not a public holiday for
// any office.
//
// Uses timefy.AddDay to advance the cursor one day at a time and
// timefy.New(day).IsWeekday() for the weekend guard.
// ─────────────────────────────────────────────────────────────────────────────
func nextWorkingDays(offices []Office, from time.Time, n int) []time.Time {
	var days []time.Time
	cursor := timefy.BeginOfDay(from)

	for len(days) < n {
		cursor = timefy.AddDay(cursor, 1)
		tx := timefy.New(cursor)
		if !tx.IsWeekday() {
			continue
		}
		holiday := false
		for _, o := range offices {
			if o.IsHoliday(cursor) {
				holiday = true
				break
			}
		}
		if !holiday {
			days = append(days, cursor)
		}
	}
	return days
}

// ─────────────────────────────────────────────────────────────────────────────
// worldClockSnapshot prints each office's local time for a given UTC instant.
// ─────────────────────────────────────────────────────────────────────────────
func worldClockSnapshot(offices []Office, utc time.Time) {
	fmt.Printf("\n  World clock @ %s\n", utc.Format("2006-01-02 15:04 UTC"))
	fmt.Println("  ─────────────────────────────────────────────────────────")
	for _, o := range offices {
		status := "✓ open"
		if !o.IsBusinessHour(utc) {
			status = "✗ closed"
		}
		if o.IsHoliday(utc) {
			status = "⊘ holiday"
		}
		fmt.Printf("  %-18s  %s  [%s]\n", o.Name, o.LocalTime(utc), status)
	}
}

func main() {
	// ── Office roster ─────────────────────────────────────────────────────────
	offices := []Office{
		{
			Name:          "Ho Chi Minh City",
			Zone:          timefy.DefaultTimezoneVietnam,
			BusinessStart: 9,
			BusinessEnd:   18,
			// Hung Kings Commemoration 2025 (public holiday in Vietnam)
			PublicHolidays: []time.Time{
				time.Date(2025, time.April, 7, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			Name:          "London",
			Zone:          timefy.DefaultTimezoneLondon,
			BusinessStart: 9,
			BusinessEnd:   18,
			// Early May Bank Holiday 2025
			PublicHolidays: []time.Time{
				time.Date(2025, time.May, 5, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			Name:          "New York",
			Zone:          timefy.DefaultTimezoneNewYork,
			BusinessStart: 9,
			BusinessEnd:   18,
			// Memorial Day 2025
			PublicHolidays: []time.Time{
				time.Date(2025, time.May, 26, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			Name:          "Tokyo",
			Zone:          timefy.DefaultTimezoneTokyo,
			BusinessStart: 9,
			BusinessEnd:   18,
			// Showa Day 2025
			PublicHolidays: []time.Time{
				time.Date(2025, time.April, 29, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			Name:          "Sydney",
			Zone:          timefy.DefaultTimezoneSydney,
			BusinessStart: 9,
			BusinessEnd:   18,
			PublicHolidays: []time.Time{
				time.Date(2025, time.April, 25, 0, 0, 0, 0, time.UTC), // Anzac Day
			},
		},
	}

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("     TIMEFY — GLOBAL TEAM COORDINATOR DEMO               ")
	fmt.Println("═══════════════════════════════════════════════════════════")

	// ── Snapshot 1: Monday morning UTC ───────────────────────────────────────
	mondayMorningUTC := time.Date(2025, time.March, 24, 8, 0, 0, 0, time.UTC)
	worldClockSnapshot(offices, mondayMorningUTC)

	// ── Snapshot 2: A time when both EU and APAC are open ────────────────────
	// 08:00 UTC = 15:00 HCM, 08:00 London (spring, no DST yet), 17:00 Tokyo
	tuesdayAM := time.Date(2025, time.March, 25, 8, 0, 0, 0, time.UTC)
	worldClockSnapshot(offices, tuesdayAM)

	// ── Find all 1-hour overlap windows in the next 48 hours ─────────────────
	from := time.Date(2025, time.March, 24, 0, 0, 0, 0, time.UTC)
	overlaps := globalOverlapSlots(offices, from, 48)

	fmt.Printf("\n  Global overlap windows (all 5 offices open) in next 48h:\n")
	fmt.Println("  ─────────────────────────────────────────────────────────")
	if len(overlaps) == 0 {
		fmt.Println("  ✗ No overlap found — the world is too spread out!")
	} else {
		for _, slot := range overlaps {
			fmt.Printf("  • %s UTC\n", slot.Format("Mon 2006-01-02 15:04"))
		}
	}

	// ── Next 5 no-holiday working days across all offices ────────────────────
	startScan := time.Date(2025, time.April, 6, 0, 0, 0, 0, time.UTC)
	workdays := nextWorkingDays(offices, startScan, 5)

	fmt.Printf("\n  Next 5 universal working days (no holidays in any office):\n")
	fmt.Println("  ─────────────────────────────────────────────────────────")
	for i, d := range workdays {
		fmt.Printf("  %d. %s\n", i+1, d.Format("Monday 2006-01-02"))
	}

	// ── Meeting proposal check for a specific UTC time ────────────────────────
	proposed := time.Date(2025, time.March, 27, 7, 0, 0, 0, time.UTC)
	fmt.Printf("\n  Meeting feasibility check: %s\n", proposed.Format("2006-01-02 15:04 UTC"))
	fmt.Println("  ─────────────────────────────────────────────────────────")
	for _, o := range offices {
		available := o.IsBusinessHour(proposed) && !o.IsHoliday(proposed)
		icon := "✓"
		if !available {
			icon = "✗"
		}
		local := o.LocalTime(proposed)
		fmt.Printf("  %s %-18s  %s\n", icon, o.Name, local)
	}

	// ── Quarter summary: which quarter are we in right now? ───────────────────
	fmt.Printf("\n  Current quarter (UTC clock): Q%d of %d\n",
		timefy.Quarter(), time.Now().Year())
	qStart := timefy.New(time.Now()).BeginningOfQuarter()
	qEnd := timefy.New(time.Now()).EndOfQuarter()
	elapsed := timefy.New(time.Now()).DurationInDays(qStart)
	total := timefy.New(qEnd).DurationInDays(qStart) + 1
	fmt.Printf("  Quarter span: %s → %s (%d days total, %d elapsed)\n",
		qStart.Format("Jan 2"), qEnd.Format("Jan 2"),
		total, elapsed)

	// ── Relative time display for the last major release ─────────────────────
	lastRelease := time.Date(2025, time.January, 15, 10, 0, 0, 0, time.UTC)
	fmt.Printf("\n  Last platform release was: %s\n",
		timefy.New(lastRelease).TimeAgo())

	nextDeployment := time.Date(2025, time.April, 1, 12, 0, 0, 0, time.UTC)
	fmt.Printf("  Next scheduled deployment: %s\n",
		timefy.New(nextDeployment).TimeUntil())
}
